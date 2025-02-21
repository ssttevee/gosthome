package sensor

import (
	"context"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/registry"
)

type testSensorConfig struct {
	BaseSensorConfig[testSensor, *testSensor] `yaml:",inline"`

	Abc string `yaml:"abc"`
}

func (t *testSensorConfig) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		t.BaseSensorConfig.ValidateWithContext(ctx),
		validation.ValidateStructWithContext(ctx, t,
			validation.Field(t.Abc, cv.String(cv.OneOf("yay", "nay"))),
		),
	)
}

type testSensor struct {
	BaseSensor[testSensor, *testSensor]
}

// Close implements component.Component.
func (t *testSensor) Close() error {
	return nil
}

// InitializationPriority implements component.Component.
func (t *testSensor) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityProcessor
}

// Setup implements component.Component.
func (t *testSensor) Setup() {
}

func newTestSensor(ctx context.Context, cfg *testSensorConfig) (ret *testSensor, err error) {
	ret = &testSensor{}
	ret.BaseSensor, err = NewBaseSensor(ctx, ret, &cfg.BaseSensorConfig)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

var _ component.Component = (*testSensor)(nil)
var _ entity.Sensor = (*testSensor)(nil)

type testSensorDeclaration struct {
}

// Component implements component.Declaration.
func (t *testSensorDeclaration) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	panic("unimplemented")
}

// Config implements component.Declaration.
func (t *testSensorDeclaration) Config() *component.ConfigDecoder {
	return &component.ConfigDecoder{
		Config:    &testSensorConfig{},
		Marshal:   component.Marshal[testSensorConfig],
		Unmarshal: component.Unmarshal[testSensorConfig],
	}
}

var _ component.Declaration = (*testSensorDeclaration)(nil)

func TestBaseSensorComponent(t *testing.T) {
	reg := registry.NewRegistry()
	reg.RegisterEntityComponent(entity.DomainTypeSensor, "test", &testSensorDeclaration{})
}
