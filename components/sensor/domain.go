package sensor

import (
	"context"
	"log/slog"

	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/config"
	"github.com/gosthome/gosthome/core/entity"
)

type Config struct {
	component.ConfigOf[entity.SensorDomain, *entity.SensorDomain]
	config.PlatformConfig
}

// Validate implements component.Config.
func (c *Config) Validate() error {
	return nil
}

func NewConfig() *Config {
	return &Config{
		PlatformConfig: config.PlatformConfig{
			DomainType: entity.DomainTypeSensor,
		},
	}
}

func New(ctx context.Context, c *Config) ([]component.Component, error) {
	node := core.GetNode(ctx)
	if node == nil {
		panic("No node in config during sensor initialization")
	}
	domain := &entity.SensorDomain{}
	ret := []component.Component{domain}
	for _, platformConfig := range c.Configs {
		cd, ok := node.Config.Registry.GetEntityComponent(entity.DomainTypeSensor, platformConfig.Platform)
		if !ok {
			panic("unregistered sensor platform in config " + platformConfig.Platform)
		}
		comp, err := cd.Component(ctx, platformConfig.Config.Config)
		if err != nil {
			return nil, err
		}
		for _, c := range comp {
			domain.Register(c.(entity.Sensor))
		}
		ret = append(ret, comp...)
	}
	slog.Info("Initialized sensor domain")
	err := node.CreateDomain(entity.PublicDomain(domain))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

var _ component.Config = (*Config)(nil)
