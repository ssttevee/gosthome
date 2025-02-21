package common

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"sync"

	ws "github.com/coder/websocket"
	"github.com/gosthome/gosthome/core/guarded"
)

type tcpWsConn struct {
	net.Conn
	la, ra websocketAddr
}

// LocalAddr implements net.Conn.
func (t *tcpWsConn) LocalAddr() net.Addr {
	return t.la
}

// RemoteAddr implements net.Conn.
func (t *tcpWsConn) RemoteAddr() net.Addr {
	return t.ra
}

var _ (net.Conn) = (*tcpWsConn)(nil)

type websocketAddr string

func (a websocketAddr) Network() string {
	return "tcp/websocket"
}

func (a websocketAddr) String() string {
	return string(a)
}

type WSProxyClient struct {
	Addr string
}

func (proxy *WSProxyClient) DialWS(ctx context.Context, addr string) (net.Conn, error) {
	wsaddr := proxy.Addr + "/" + addr
	c, _, err := ws.Dial(ctx, wsaddr, &ws.DialOptions{})
	if err != nil {
		return nil, err
	}
	return &tcpWsConn{
		Conn: ws.NetConn(ctx, c, ws.MessageBinary),
		la:   websocketAddr("unknown"),
		ra:   websocketAddr(addr),
	}, nil
}

type WSProxy struct {
}

func (l *WSProxy) RegisterHandler(mux *http.ServeMux) {
	mux.Handle("/{addr}/{port}", l)
}
func (l *WSProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addr := r.PathValue("addr")
	if addr == "" {
		http.Error(w, "no addr in path", http.StatusMisdirectedRequest)
		return
	}

	portStr := r.PathValue("port")
	if portStr == "" {
		http.Error(w, "no port in path", http.StatusMisdirectedRequest)
		return
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMisdirectedRequest)
		return
	}
	wc, err := ws.Accept(w, r, nil)
	if err != nil {
		slog.Error("Error accepting ", "err", err)
		return
	}
	defer wc.CloseNow()
	dc, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if addr == "" {
		slog.Error("Failed connecting to destination", "addr", addr, "err", err)
		return
	}
	defer dc.Close()
	nc := ws.NetConn(r.Context(), wc, ws.MessageBinary)
	var streamWait sync.WaitGroup
	streamWait.Add(2)

	streamConn := func(dst io.Writer, src io.Reader) {
		io.Copy(dst, src)
		streamWait.Done()
	}

	go streamConn(dc, nc)
	go streamConn(nc, dc)

	streamWait.Wait()
}

type wsServerConn struct {
	net.Conn
	c chan<- struct{}
}

// Close implements net.Conn.
func (w *wsServerConn) Close() (err error) {
	if w.Conn != nil {
		err = w.Conn.Close()
		w.Conn = nil
		close(w.c)
		w.c = nil
	}
	return
}

type WSServer struct {
	addr string
	ncr  <-chan *wsServerConn
	ncw  *guarded.Value[chan<- *wsServerConn]
}

func NewWSServer() *WSServer {
	nc := make(chan *wsServerConn)
	return &WSServer{
		ncr: nc,
		ncw: guarded.New((chan<- *wsServerConn)(nc)),
	}
}

func (s *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wc, err := ws.Accept(w, r, nil)
	if err != nil {
		slog.Error("Error accepting ", "err", err)
		return
	}
	defer wc.CloseNow()
	nc := make(chan struct{})
	ctx := r.Context()
	conn := &wsServerConn{
		Conn: ws.NetConn(ctx, wc, ws.MessageBinary),
		c:    nc,
	}
	done := ctx.Done()
	s.ncw.DoErr(func(c *chan<- *wsServerConn) error {
		if *c == nil {
			return errors.New("trying to connect to closed server")
		}
		select {
		case *c <- conn:
			return nil
		case <-done:
			return ctx.Err()
		}
	})
	select {
	case <-nc:
		return
	case <-done:
		return
	}
}

func (s *WSServer) ListenWs(ctx context.Context, addr string) (net.Listener, error) {
	s.addr = addr
	return s, nil
}

// Accept implements net.Listener.
func (s *WSServer) Accept() (net.Conn, error) {
	c, ok := <-s.ncr
	if !ok {
		return nil, errors.New("accepting on closed socket")
	}
	return c, nil
}

// Addr implements net.Listener.
func (s *WSServer) Addr() net.Addr {
	return websocketAddr(s.addr)
}

// Close implements net.Listener.
func (s *WSServer) Close() error {
	s.ncw.Do(func(c *chan<- *wsServerConn) {
		if *c == nil {
			return
		}
		close(*c)
		*c = nil
	})
	return nil
}

var _ net.Listener = (*WSServer)(nil)
var _ net.Conn = (*wsServerConn)(nil)
