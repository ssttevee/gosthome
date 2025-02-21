package psutil

import (
	"github.com/gosthome/gosthome/components/sensor"
	"github.com/gosthome/gosthome/core/component"
)

type SensorConfig struct {
	sensor.BaseSensorConfig[Sensor, *Sensor]
}

type Sensor struct {
	sensor.BaseSensor[Sensor, *Sensor]
	component.WithInitializationPriorityProcessor
}

// Setup implements component.Component.
func (cpu *Sensor) Setup() {

}

func (cpu *Sensor) Close() error {
	return nil
}
