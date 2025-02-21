package component

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"time"
	"weak"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gosthome/gosthome/core/guarded"
)

type PollingComponentConfig struct {
	UpdateInterval time.Duration `yaml:"update_interval"`
}

// Validate implements Config.
func (p *PollingComponentConfig) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(
			&p.UpdateInterval,
			validation.Required,
			validation.Min(0*time.Second).Exclusive()),
	)
}

type Poller interface {
	Component
	Poll()
}

type poller struct {
	basectx        context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	updateInterval time.Duration
	handle         func()
}

func (p *poller) start() {
	if p.cancel != nil {
		return
	}
	ctx, cancel := context.WithCancel(p.basectx)
	p.cancel = cancel
	p.wg.Add(1)
	done := ctx.Done()
	go func() {
		defer p.wg.Done()
		timer := time.NewTicker(p.updateInterval)
		for {
			select {
			case <-timer.C:
				p.handle()
			case <-done:
				timer.Stop()
				return
			}
		}
	}()
}

func (p *poller) stop() {
	if p.cancel != nil {
		p.cancel()
		p.wg.Wait()
		p.cancel = nil
		p.wg = sync.WaitGroup{}
	}
}

type PollingComponent[T any, PT interface {
	*T
	Poller
}] struct {
	poller *guarded.Value[poller]
	type_  string
	poll   weak.Pointer[T]
}

func NewPollingComponent[T any, PT interface {
	*T
	Poller
}](ctx context.Context, poll PT, cfg *PollingComponentConfig) (*PollingComponent[T, PT], error) {
	p := &PollingComponent[T, PT]{
		type_: reflect.TypeOf(PT(nil)).String(),
		poll:  weak.Make(poll),
	}
	p.poller = guarded.New(poller{
		basectx:        ctx,
		cancel:         nil,
		wg:             sync.WaitGroup{},
		updateInterval: cfg.UpdateInterval,
		handle: func() {
			t := p.poll.Value()
			if t != nil {
				slog.Debug("Poll", "type", p.type_)
				PT(t).Poll()
			} else {
				slog.Warn("Unable to poll for expired", "type", p.type_)
			}
		},
	})
	return p, nil
}

// Setup implements Component.
func (p *PollingComponent[T, PT]) Setup() {
	p.Start()
}

func (p *PollingComponent[T, PT]) Start() {
	slog.Debug("Starting polling for", "type", p.type_)
	p.poller.Do((*poller).start)
}

func (p *PollingComponent[T, PT]) Stop() {
	p.poller.Do((*poller).stop)
	slog.Debug("Stopped polling for", "type", p.type_)
}

// Close implements Component.
func (p *PollingComponent[T, PT]) Close() error {
	p.Stop()
	return nil
}
