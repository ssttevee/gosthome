package demo

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/gosthome/gosthome/components/sensor"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/entity"
)

type DemoSensorConfig struct {
	sensor.BaseSensorConfig[DemoSensor, *DemoSensor] `yaml:",inline"`

	Poll component.PollingComponentConfig `yaml:",inline"`
}

func NewDemoSensorConfig() DemoSensorConfig {
	return DemoSensorConfig{
		Poll: component.PollingComponentConfig{
			UpdateInterval: 10 * time.Second,
		},
	}
}

type DemoSensor struct {
	sensor.BaseSensor[DemoSensor, *DemoSensor]
	poll *component.PollingComponent[DemoSensor, *DemoSensor]
	rand *rand.Rand
}

func NewDemoSensor(ctx context.Context, src rand.Source, cfg *DemoSensorConfig) (d *DemoSensor, err error) {
	d = &DemoSensor{
		rand: rand.New(src),
	}
	d.BaseSensor, err = sensor.NewBaseSensor(ctx, d, &cfg.BaseSensorConfig)
	if err != nil {
		return nil, err
	}
	d.poll, err = component.NewPollingComponent(ctx, d, &cfg.Poll)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// InitializationPriority implements component.Poller.
func (d *DemoSensor) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityData
}

// Setup implements component.Poller.
func (d *DemoSensor) Setup() {
	d.Poll()
	d.poll.Setup()
}

// Poll implements component.Poller.
func (d *DemoSensor) Poll() {
	d.SetState(entity.SensorState{
		State:        d.rand.Float32(),
		MissingState: false,
	})
}

// Close implements component.Poller.
func (d *DemoSensor) Close() error {
	return d.poll.Close()
}

var _ (entity.Sensor) = (*DemoSensor)(nil)
var _ (component.Poller) = (*DemoSensor)(nil)
