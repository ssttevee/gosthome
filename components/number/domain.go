package number

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

// SetValue is a service request for setting the (float) value of a Number entity.
// It mirrors the NumberCommandRequest coming from the API layer.
type SetValue struct {
	Key   uint32
	State float32
}

// ServiceType implements bus.ServiceRequestData.
func (n *SetValue) ServiceType() string {
	return "number.set"
}

var _ bus.ServiceRequestData = (*SetValue)(nil)

// Config holds the number domain configuration.
type Config struct {
	component.ConfigOf[entity.NumberDomain, *entity.NumberDomain]
	config.PlatformConfig
}

// ValidateWithContext implements component.Config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return c.PlatformConfig.ValidateWithContext(ctx)
}

// NewConfig returns a default number domain config.
func NewConfig() *Config {
	return &Config{
		PlatformConfig: config.PlatformConfig{
			DomainType: entity.DomainTypeNumber,
		},
	}
}

// RegisterServiceCallHandlers registers service call handlers for the number domain.
func RegisterServiceCallHandlers(ctx context.Context, domain *entity.NumberDomain, b *bus.Bus) {
	// Handle number.set service calls.
	b.HandleServiceCalls(bus.ServiceHandlerWithRespose(b, func(t *SetValue) error {
		num, ok := domain.FindByKey(t.Key)
		if !ok {
			slog.Error("Tried to set value on nonexisting number", "key", t.Key)
			return fmt.Errorf("tried to set value on nonexisting number %d", t.Key)
		}
		return num.SetValue(ctx, t.State)
	}))
}

// New initializes the number domain, registers number entities and sets up service handlers.
func New(ctx context.Context, c *Config) ([]component.Component, error) {
	node := core.GetNode(ctx)
	if node == nil {
		panic("No node in context during number initialization")
	}
	domain := &entity.NumberDomain{}
	ret := []component.Component{domain}

	for _, platformConfig := range c.Configs {
		cd, ok := node.Config.Registry.GetEntityComponent(entity.DomainTypeNumber, platformConfig.Platform)
		if !ok {
			panic("unregistered number platform in config " + platformConfig.Platform)
		}
		comp, err := cd.Component(ctx, platformConfig.Config.Config)
		if err != nil {
			return nil, err
		}
		for _, cc := range comp {
			domain.Register(cc.(entity.Number))
		}
		ret = append(ret, comp...)
	}
	slog.Info("Initialized number domain")
	if err := node.CreateDomain(entity.PublicDomain(domain)); err != nil {
		return nil, err
	}

	b := bus.Get(ctx)
	if b == nil {
		panic("No bus in context during number initialization")
	}

	RegisterServiceCallHandlers(ctx, domain, b)

	return ret, nil
}

var _ component.Config = (*Config)(nil)
