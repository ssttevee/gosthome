package climate

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

// SetState is a service request used to change (parts of) the state of a Climate entity.
// Each Has* flag indicates whether the corresponding value should be applied. This mirrors
// the structure of the ESPHome API ClimateCommandRequest so the API layer can translate
// incoming commands directly to this service request.
type SetState struct {
	Key uint32

	HasMode              bool
	Mode                 entity.ClimateMode
	HasTargetTemperature bool
	TargetTemperature    float32

	HasTargetTemperatureLow  bool
	TargetTemperatureLow     float32
	HasTargetTemperatureHigh bool
	TargetTemperatureHigh    float32

	HasLegacyAway bool
	LegacyAway    bool

	HasFanMode   bool
	FanMode      entity.ClimateFanMode
	HasSwingMode bool
	SwingMode    entity.ClimateSwingMode

	HasCustomFanMode bool
	CustomFanMode    string

	HasPreset       bool
	Preset          entity.ClimatePreset
	HasCustomPreset bool
	CustomPreset    string

	HasTargetHumidity bool
	TargetHumidity    float32
}

// ServiceType implements bus.ServiceRequestData.
func (c *SetState) ServiceType() string {
	return "climate.set_state"
}

var _ bus.ServiceRequestData = (*SetState)(nil)

// Config for the climate domain.
type Config struct {
	component.ConfigOf[entity.ClimateDomain, *entity.ClimateDomain]
	config.PlatformConfig
}

// ValidateWithContext implements component.Config.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return c.PlatformConfig.ValidateWithContext(ctx)
}

// NewConfig returns a default climate domain config.
func NewConfig() *Config {
	return &Config{
		PlatformConfig: config.PlatformConfig{
			DomainType: entity.DomainTypeClimate,
		},
	}
}

// RegisterServiceCallHandlers registers service call handlers for the climate domain.
func RegisterServiceCallHandlers(ctx context.Context, domain *entity.ClimateDomain, b *bus.Bus) {
	// Handle climate.set_state service calls.
	b.HandleServiceCalls(bus.ServiceHandlerWithRespose(b, func(t *SetState) error {
		cl, ok := domain.FindByKey(t.Key)
		if !ok {
			slog.Error("Tried to set state on nonexisting climate entity", "key", t.Key)
			return fmt.Errorf("tried to set state on nonexisting climate entity %d", t.Key)
		}
		cur := cl.State()

		if t.HasMode {
			cur.Mode = t.Mode
		}
		if t.HasTargetTemperature {
			cur.TargetTemperature = t.TargetTemperature
		}
		if t.HasTargetTemperatureLow {
			cur.TargetTemperatureLow = t.TargetTemperatureLow
		}
		if t.HasTargetTemperatureHigh {
			cur.TargetTemperatureHigh = t.TargetTemperatureHigh
		}
		if t.HasFanMode {
			cur.FanMode = t.FanMode
		}
		if t.HasSwingMode {
			cur.SwingMode = t.SwingMode
		}
		if t.HasCustomFanMode {
			cur.CustomFanMode = t.CustomFanMode
		}
		if t.HasLegacyAway {
			if t.LegacyAway {
				cur.Preset = entity.ClimatePresetAway
			} else if cur.Preset == entity.ClimatePresetAway {
				cur.Preset = entity.ClimatePresetNone
			}
		}
		if t.HasPreset {
			cur.Preset = t.Preset
		}
		if t.HasCustomPreset {
			cur.CustomPreset = t.CustomPreset
		}
		if t.HasTargetHumidity {
			cur.TargetHumidity = t.TargetHumidity
		}
		return cl.SetState(ctx, cur)
	}))
}

// New initializes the climate domain, registers entities, and hooks up service handlers.
func New(ctx context.Context, c *Config) ([]component.Component, error) {
	node := core.GetNode(ctx)
	if node == nil {
		panic("No node in context during climate initialization")
	}
	domain := &entity.ClimateDomain{}
	ret := []component.Component{domain}

	for _, platformConfig := range c.Configs {
		cd, ok := node.Config.Registry.GetEntityComponent(entity.DomainTypeClimate, platformConfig.Platform)
		if !ok {
			panic("unregistered climate platform in config " + platformConfig.Platform)
		}
		comp, err := cd.Component(ctx, platformConfig.Config.Config)
		if err != nil {
			return nil, err
		}
		for _, cc := range comp {
			domain.Register(cc.(entity.Climate))
		}
		ret = append(ret, comp...)
	}

	slog.Info("Initialized climate domain")
	if err := node.CreateDomain(entity.PublicDomain(domain)); err != nil {
		return nil, err
	}

	b := bus.Get(ctx)
	if b == nil {
		panic("No bus in context during climate initialization")
	}

	RegisterServiceCallHandlers(ctx, domain, b)

	return ret, nil
}

var _ component.Config = (*Config)(nil)
