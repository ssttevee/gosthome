package psutil

import (
	"github.com/gosthome/gosthome/components/textsensor"
	"github.com/gosthome/gosthome/core/component"
)

type TextSensorConfig struct {
	textsensor.BaseTextSensorConfig[TextSensor, *TextSensor]
}

type TextSensor struct {
	textsensor.BaseTextSensor[TextSensor, *TextSensor]
	component.WithInitializationPriorityProcessor
}

// Setup implements component.Component.
func (cpu *TextSensor) Setup() {
}

func (cpu *TextSensor) Close() error {
	return nil
}
