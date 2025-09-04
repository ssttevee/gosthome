package switchcomp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/config"
	"github.com/gosthome/gosthome/core/entity"
)

// SetState is a service request for changing the state of a switch entity.
// It mirrors the pattern used by the button component's service request type.
type SetState struct {
	Key   uint32
	State bool
}

// ServiceType implements bus.ServiceRequestData.
func (s *SetState) ServiceType() string {
	return "switch.set_state"
}

var _ bus.ServiceRequestData = (*SetState)(nil)

// Config holds the switch domain configuration.
type Config struct {
	component.ConfigOf[entity.SwitchDomain, *entity.SwitchDomain]
	config.PlatformConfig
}

// ValidateWithContext implements component.Config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return c.PlatformConfig.ValidateWithContext(ctx)
}

// NewConfig returns a default switch domain config.
func NewConfig() *Config {
	return &Config{
		PlatformConfig: config.PlatformConfig{
			DomainType: entity.DomainTypeSwitch,
		},
	}
}

// RegisterServiceCallHandlers registers service call handlers for the switch domain.
func RegisterServiceCallHandlers(ctx context.Context, domain *entity.SwitchDomain, b *bus.Bus) {
	b.HandleServiceCalls(bus.ServiceHandlerWithRespose(b, func(t *SetState) error {
		sw, ok := domain.FindByKey(t.Key)
		if !ok {
			slog.Error("Tried to set state on nonexisting switch", "key", t.Key)
			return fmt.Errorf("tried to set state on nonexisting switch %d", t.Key)
		}
		return sw.SetState(ctx, t.State)
	}))
}

// New initializes the switch domain, registers switch entities and sets up service handlers.
func New(ctx context.Context, c *Config) ([]component.Component, error) {
	node := core.GetNode(ctx)
	if node == nil {
		panic("No node in context during switch initialization")
	}
	domain := &entity.SwitchDomain{}
	ret := []component.Component{domain}

	for _, platformConfig := range c.Configs {
		cd, ok := node.Config.Registry.GetEntityComponent(entity.DomainTypeSwitch, platformConfig.Platform)
		if !ok {
			panic("unregistered switch platform in config " + platformConfig.Platform)
		}
		comp, err := cd.Component(ctx, platformConfig.Config.Config)
		if err != nil {
			return nil, err
		}
		for _, cc := range comp {
			domain.Register(cc.(entity.Switch))
		}
		ret = append(ret, comp...)
	}
	slog.Info("Initialized switch domain")
	if err := node.CreateDomain(entity.PublicDomain(domain)); err != nil {
		return nil, err
	}

	b := bus.Get(ctx)
	if b == nil {
		panic("No bus in context during switch initialization")
	}

	RegisterServiceCallHandlers(ctx, domain, b)

	return ret, nil
}

var _ component.Config = (*Config)(nil)
