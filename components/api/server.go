package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"sync"

	"maps"

	"github.com/gosthome/gosthome/components/api/common"
	ehp "github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
)

type Connection struct {
	server          *Server
	sendFrames      frameshakers.FrameSenderFunc
	authenticated   bool
	canAuthenticate bool
	subscribed      bool
	clientInfo      string
	asyncHandlers   asyncHandlers
	busEvents       []bus.EventSubsciption
}

func (c *Connection) SendMessages(msgs []ehp.EsphomeMessageTyper) error {
	frames, err := common.EncodeFrames(msgs)
	if err != nil {
		return err
	}
	return c.sendFrames(frames)
}

func (c *Connection) Server() *Server {
	return c.server
}

func (c *Connection) Authenticated() bool {
	return c.authenticated
}

type AnyMessageHandler func(ctx context.Context, c *Connection, m ehp.EsphomeMessageTyper) ([]ehp.EsphomeMessageTyper, error)

type MessageHandler struct {
	Type    ehp.MessageType
	Handler AnyMessageHandler
}

type HandlerFunc[T any, PT interface {
	ehp.EsphomeMessageTyper
	*T
}] func(ctx context.Context, c *Connection, msg PT) ([]ehp.EsphomeMessageTyper, error)

func Handler[T any, PT interface {
	ehp.EsphomeMessageTyper
	*T
}](hf HandlerFunc[T, PT]) MessageHandler {
	return MessageHandler{
		Type: (PT)(nil).EsphomeMessageType(),
		Handler: func(ctx context.Context, c *Connection, m ehp.EsphomeMessageTyper) ([]ehp.EsphomeMessageTyper, error) {
			return hf(ctx, c, m.(PT))
		},
	}
}

// Handle implements frameshakers.FramesHandler.
func (c *Connection) Handle(ctx context.Context, input []frameshakers.Frame) (retFrames []frameshakers.Frame, err error) {
	closing := false
	for _, frame := range input {
		_, msg, err := common.DecodeFrame(frame)
		if err != nil {
			return nil, err
		}
		if slog.Default().Enabled(ctx, slog.LevelDebug) {
			slog.Default().Debug("handleRPC", "msg", reflect.TypeOf(msg))
		}
		var h AnyMessageHandler
		var ok bool
		c.server.handlers.rlocked(func(m map[ehp.MessageType]AnyMessageHandler) {
			h, ok = m[ehp.MessageType(frame.Type)]
		})
		if !ok {
			slog.Warn("Unhandled message", "msg", reflect.TypeOf(msg))
			continue
		}

		retMsgs, err := h(ctx, c, msg)
		if err != nil {
			if !errors.Is(err, frameshakers.ErrCloseConnection) {
				return nil, err
			}
			closing = true
		}
		encMsgFrames, err := common.EncodeFrames(retMsgs)
		if err != nil {
			return nil, err
		}
		retFrames = append(retFrames, encMsgFrames...)
	}
	if closing {
		err = frameshakers.ErrCloseConnection
	}
	return
}

func (c *Connection) Close() error {
	for _, sub := range c.busEvents {
		sub.Close()
	}
	return c.asyncHandlers.Close()
}

type asyncHandlers struct {
	cancels []context.CancelFunc
	wg      sync.WaitGroup
}

func (ah *asyncHandlers) Close() error {
	for _, c := range ah.cancels {
		c()
	}
	ah.wg.Wait()
	return nil
}

func AsyncHandler(mh MessageHandler) MessageHandler {
	oh := mh.Handler
	mh.Handler = func(ctx context.Context, c *Connection, m ehp.EsphomeMessageTyper) ([]ehp.EsphomeMessageTyper, error) {
		ctx, cancel := context.WithCancel(ctx)
		c.asyncHandlers.cancels = append(c.asyncHandlers.cancels, cancel)
		c.asyncHandlers.wg.Add(1)
		go func() {
			defer c.asyncHandlers.wg.Done()
			frames, herr := oh(ctx, c, m)
			if herr != nil {
				slog.Error("error handling async message", "type", fmt.Sprintf("%T", m), "err", herr)
				return
			}
			herr = c.SendMessages(frames)
			if herr != nil {
				slog.Error("error sending handled frames", "type", fmt.Sprintf("%T", m), "err", herr)
				return
			}
		}()
		return nil, nil
	}
	return mh
}

type safeMessageHandlers struct {
	sync.RWMutex
	m map[ehp.MessageType]AnyMessageHandler
}

func (h *safeMessageHandlers) locked(f func(m map[ehp.MessageType]AnyMessageHandler)) {
	h.Lock()
	defer h.Unlock()
	f(h.m)
}

func (h *safeMessageHandlers) rlocked(f func(m map[ehp.MessageType]AnyMessageHandler)) {
	h.RLock()
	defer h.RUnlock()
	f(h.m)
}

type Server struct {
	cid.CID
	component.WithInitializationPriorityAfterWifi

	baseCtx context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup

	listenerFactory common.ListenerFactory
	listener        net.Listener

	shaker   frameshakers.ServerShaker
	handlers safeMessageHandlers

	config *Config
}

func New(ctx context.Context, cfg *Config) (n []component.Component, err error) {
	s, err := NewServer(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return []component.Component{s}, nil
}

type ServerOpt func(*Server)

func WithListenerFactory(lf common.ListenerFactory) ServerOpt {
	return func(s *Server) {
		s.listenerFactory = lf
	}
}

func NewServer(ctx context.Context, cfg *Config, opts ...ServerOpt) (n *Server, err error) {
	n = &Server{
		CID:             cid.NewID(cfg.ID),
		shaker:          frameshakers.PlaintextServer,
		handlers:        safeMessageHandlers{m: map[ehp.MessageType]AnyMessageHandler{}},
		listenerFactory: common.ListenTCP,
		config:          cfg,
	}
	for _, opt := range opts {
		opt(n)
	}
	n.baseCtx, n.cancel = context.WithCancel(ctx)
	// n.baseCtx = frameshakers.ContextWithValue(n.baseCtx, "serverName", cfg.Gosthome.Name)
	if n.config.Encryption.Key.Valid() {
		n.baseCtx = frameshakers.ContextWithValue(n.baseCtx, "noisePSK", n.config.Encryption.Key)
		n.shaker = frameshakers.NoiseServer
		slog.Debug("Api is starting with noise frame shaker")
	} else {
		slog.Debug("Api is starting with plaintext frame shaker")
	}
	n.handlers.locked(func(m map[ehp.MessageType]AnyMessageHandler) {
		defaultHandlers.rlocked(func(dm map[ehp.MessageType]AnyMessageHandler) {
			maps.Copy(m, dm)
		})
	})
	return n, nil
}

// Setup implements component.Component.
func (n *Server) Setup() {
	var err error
	n.listener, err = n.listenerFactory(n.baseCtx, fmt.Sprintf("%s:%d", n.config.Address, n.config.Port))
	if err != nil {
		slog.Error("Failed to initialize api", "err", err)
		return
	}
	go n.run()
}

func (n *Server) run() {
	for {
		nconn, err := n.listener.Accept()
		if err != nil {
			slog.Error("gosthome.Node got accept error", "err", err)
			return
		}
		slog.Info("Accepting connection", "from", nconn.RemoteAddr())
		n.wg.Add(1)
		go func() {
			defer n.wg.Done()
			defer nconn.Close()
			r, w := frameshakers.SplitConnection(nconn)
			rerr := n.shaker(n.baseCtx, r, w, n.connection)
			if rerr != nil {
				slog.Error("handling connection failed", "err", rerr)
			}
			slog.Debug("Done serving connection", "from", nconn.RemoteAddr())
		}()
	}
}

func (n *Server) connection(sendFrames frameshakers.FrameSenderFunc) (handler frameshakers.FramesHandler, err error) {
	c := &Connection{
		server: n,

		authenticated: !n.config.Password.Valid(),
		sendFrames:    sendFrames,
	}
	return c, nil
}

func (n *Server) Close() error {
	if n.listener != nil {
		n.listener.Close()
		n.cancel()
		n.wg.Wait()
	}
	return nil
}

var _ component.Component = (*Server)(nil)
