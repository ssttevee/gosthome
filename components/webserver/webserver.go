package webserver

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"

	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
)

type Config struct {
	cid.IDConfig
	component.ConfigOf[WebServer, *WebServer]
	Address string `yaml:"address"`
	Port    uint16 `yaml:"port"`
}

func NewConfig() *Config {
	return &Config{
		Port: 6057,
	}
}

// Validate implements validation.Validatable.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	// return validation.ValidateStruct(c, validation.Field(&c.BinarySensors))

	return nil
}

var _ component.Config = (*Config)(nil)

type WebServer struct {
	cid.CID
	component.WithInitializationPriorityAfterWifi
	cfg *Config
	ctx context.Context

	server *http.Server
}

func New(ctx context.Context, cfg *Config) ([]component.Component, error) {
	id := cfg.ID
	if id == "" {
		id = "webserver"
	}
	return []component.Component{&WebServer{
		CID: cid.NewID(id),
		cfg: cfg,
		ctx: ctx,
	}}, nil
}

// Setup implements component.Component.
func (ws *WebServer) Setup() {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(ws.home))
	ws.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ws.cfg.Address, ws.cfg.Port),
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ws.ctx
		},
	}

	go func() {
		err := ws.server.ListenAndServe()
		slog.Error("webserver closed", "err", err)
	}()

	node := core.GetNode(ws.ctx)
	node.Bus.HandleEvents(bus.EventHandler[bus.StateChangeEvent](func(e *bus.StateChangeEvent) {
	}))
}

func (ws *WebServer) home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// Close implements component.Component.
func (ws *WebServer) Close() error {
	if ws.server != nil {
		ws.server.Close()
	}
	return nil
}

var _ component.Component = (*WebServer)(nil)
