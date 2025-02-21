package psutil

import (
	"context"
	"errors"

	"github.com/gosthome/gosthome/components/binarysensor"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
)

type Config struct {
	component.ConfigOf[PSUtil, *PSUtil]
	CPU     CPUConfig     `yaml:"cpu"`
	Host    HostConfig    `yaml:"host"`
	Sensors SensorsConfig `yaml:"sensors"`
}

func NewConfig() *Config {
	return &Config{
		CPU:     NewCPUConfig(),
		Host:    NewHostConfig(),
		Sensors: NewSensorsConfig(),
	}
}

// Validate implements validation.Validatable.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	// return validation.ValidateStruct(c, validation.Field(&c.BinarySensors))
	return nil
}

// AutoLoad implements component.AutoLoader.
func (c *Config) AutoLoad() component.Dependencies {
	return component.Depends(
		binarysensor.COMPONENT_KEY,
	)
}

var _ component.Config = (*Config)(nil)
var _ component.AutoLoader = (*Config)(nil)

type PSUtil struct {
	cid.CID
	ctx context.Context
	component.WithInitializationPriorityProcessor
	cpu     *CPU
	host    *Host
	sensors *Sensors
}

func New(ctx context.Context, cfg *Config) (cmp []component.Component, err error) {
	ret := &PSUtil{
		CID: cid.NewID("psutil"),
		ctx: ctx,
	}
	ret.cpu, err = NewCPU(ctx, &cfg.CPU)
	if err != nil {
		return nil, err
	}
	ret.host, err = NewHost(ctx, &cfg.Host)
	if err != nil {
		return nil, err
	}
	ret.sensors, err = NewSensors(ctx, &cfg.Sensors)
	if err != nil {
		return nil, err
	}
	return []component.Component{ret}, nil
}

// Setup implements component.Component.
func (ps *PSUtil) Setup() {
	ps.cpu.Setup()
	ps.host.Setup()
	ps.sensors.Setup()
}

// Close implements component.Component.
func (ps *PSUtil) Close() error {
	return errors.Join(
		ps.cpu.Close(),
		ps.host.Close(),
		ps.sensors.Close(),
	)
}

var _ component.Component = (*PSUtil)(nil)
