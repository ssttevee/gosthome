package psutil

import (
	"context"
	"log/slog"
	"time"

	"github.com/gosthome/gosthome/components/sensor"
	"github.com/gosthome/gosthome/components/textsensor"
	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/util"
	psutilSensors "github.com/shirou/gopsutil/v4/sensors"
)

type SensorsConfig struct {
	component.PollingComponentConfig
}

func NewSensorsConfig() SensorsConfig {
	return SensorsConfig{
		PollingComponentConfig: component.PollingComponentConfig{
			UpdateInterval: 10 * time.Second,
		},
	}
}

// Validate implements validation.Validatable.
func (c *SensorsConfig) ValidateWithContext(ctx context.Context) error {
	return nil
}

// AutoLoad implements component.AutoLoader.
func (c *SensorsConfig) AutoLoad() []string {
	return []string{
		sensor.COMPONENT_KEY,
		textsensor.COMPONENT_KEY,
	}
}

var _ component.Config = (*Config)(nil)
var _ component.AutoLoader = (*Config)(nil)

type Sensors struct {
	cid.CID
	*component.PollingComponent[Sensors, *Sensors]
	component.WithInitializationPriorityProcessor

	ctx     context.Context
	cfg     *SensorsConfig
	sensors map[string]*Sensor
}

func NewSensors(ctx context.Context, cfg *SensorsConfig) (ret *Sensors, err error) {
	ret = &Sensors{
		CID:     cid.NewID("host"),
		ctx:     ctx,
		cfg:     cfg,
		sensors: make(map[string]*Sensor),
	}
	ret.PollingComponent, err = component.NewPollingComponent[Sensors, *Sensors](ctx, ret, &cfg.PollingComponentConfig)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Setup implements component.Component.
func (host *Sensors) Setup() {
	host.Poll()
	host.PollingComponent.Setup()
}

func (host *Sensors) getInfo(node *core.Node) {
	type hostSensor struct {
		enabled bool
		poll    bool
		name    string
		val     float64
	}
	infos, err := psutilSensors.SensorsTemperatures()
	if err != nil {
		slog.Error("SensorsTemperatures", "err", err)
		return
	}
	for _, info := range infos {
		for _, hs := range []hostSensor{
			{name: "tsensor_" + info.SensorKey + "_temperature", val: info.Temperature},
			{name: "tsensor_" + info.SensorKey + "_high", val: info.High},
			{name: "tsensor_" + info.SensorKey + "_critical", val: info.Critical},
		} {
			cfg := util.Modify(SensorConfig{}, func(c *SensorConfig) {
				c.Name = hs.name
			})
			ns, ok := host.sensors[hs.name]
			if !ok {
				ns = &Sensor{}
				ns.BaseSensor, err = sensor.NewBaseSensor(host.ctx, ns, &cfg.BaseSensorConfig)
				if err != nil {
					slog.Error("failed to create text sensor", "name", hs.name, "err", err)
					continue
				}
				node.RegisterSensor(ns)
				host.sensors[hs.name] = ns
			}
			ns.SetState(entity.SensorState{
				State:        float32(hs.val),
				MissingState: false,
			})
		}
	}
}

func (host *Sensors) Poll() {
	node := core.GetNode(host.ctx)
	host.getInfo(node)
}

// Close implements component.Component.
func (host *Sensors) Close() error {
	return host.PollingComponent.Close()
}
