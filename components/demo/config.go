package demo

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gosthome/gosthome/components/binarysensor"
	"github.com/gosthome/gosthome/components/button"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/util"
)

type Config struct {
	component.ConfigOf[Demo, *Demo]

	Seeds [2]uint64

	BinarySensors []DemoBinarySensorConfig `yaml:"binary_sensors"`
	Buttons       []DemoButtonConfig       `yaml:"buttons"`
	Sensors       []DemoSensorConfig       `yaml:"buttons"`
}

func NewConfig() *Config {
	return &Config{
		Seeds: [2]uint64{1, 2},
		BinarySensors: []DemoBinarySensorConfig{
			util.Modify(NewDemoBinarySensorConfig(), func(c *DemoBinarySensorConfig) {
				c.Name = "Demo Basement Floor Wet"
				c.DeviceClass = entity.BinarySensorDeviceClassMoisture
			}),
			util.Modify(NewDemoBinarySensorConfig(), func(c *DemoBinarySensorConfig) {
				c.Name = "Demo Movement Backyard"
				c.DeviceClass = entity.BinarySensorDeviceClassMotion
			}),
		},
		Buttons: []DemoButtonConfig{
			util.Modify(NewDemoButtonConfig(), func(c *DemoButtonConfig) {
				c.Name = "Demo Regenerate Seed"
				c.DeviceClass = entity.ButtonDeviceClassRestart
				c.Category = entity.CategoryConfig
			}),
			util.Modify(NewDemoButtonConfig(), func(c *DemoButtonConfig) {
				c.Name = "Demo Regenerate Seed Config"
				c.DeviceClass = entity.ButtonDeviceClassRestart
				c.Category = entity.CategoryConfig
			}),
		},
		Sensors: []DemoSensorConfig{
			util.Modify(NewDemoSensorConfig(), func(c *DemoSensorConfig) {
				c.Name = "Demo Temperature Sensor"
				c.DeviceClass = entity.SensorDeviceClassTemperature
				c.UnitOfMeasurement = "*C"
			}),
		},
	}
}

// Validate implements validation.Validatable.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c, validation.Field(&c.BinarySensors))
}

// AutoLoad implements component.AutoLoader.
func (c *Config) AutoLoad() component.Dependencies {
	return component.Depends(
		binarysensor.COMPONENT_KEY,
		button.COMPONENT_KEY,
	)
}

var _ component.Config = (*Config)(nil)
var _ component.AutoLoader = (*Config)(nil)
