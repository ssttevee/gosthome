package textsensor

import (
	"context"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/registry"
)

type testTextSensorConfig struct {
	BaseTextSensorConfig[testTextSensor, *testTextSensor] `yaml:",inline"`

	Abc string `yaml:"abc"`
}

func (t *testTextSensorConfig) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		t.BaseTextSensorConfig.ValidateWithContext(ctx),
		validation.ValidateStructWithContext(ctx, t,
			validation.Field(t.Abc, cv.String(cv.OneOf("yay", "nay"))),
		),
	)
}

type testTextSensor struct {
	BaseTextSensor[testTextSensor, *testTextSensor]
}

// Close implements component.Component.
func (t *testTextSensor) Close() error {
	return nil
}

// InitializationPriority implements component.Component.
func (t *testTextSensor) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityProcessor
}

// Setup implements component.Component.
func (t *testTextSensor) Setup() {
}

func newTestTextSensor(ctx context.Context, cfg *testTextSensorConfig) (ret *testTextSensor, err error) {
	ret = &testTextSensor{}
	ret.BaseTextSensor, err = NewBaseTextSensor(ctx, ret, &cfg.BaseTextSensorConfig)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

var _ component.Component = (*testTextSensor)(nil)
var _ entity.TextSensor = (*testTextSensor)(nil)

type testTextSensorDeclaration struct {
}

// Component implements component.Declaration.
func (t *testTextSensorDeclaration) Component(ctx context.Context, cfg component.Config) ([]component.Component, error) {
	panic("unimplemented")
}

// Config implements component.Declaration.
func (t *testTextSensorDeclaration) Config() *component.ConfigDecoder {
	return &component.ConfigDecoder{
		Config:    &testTextSensorConfig{},
		Marshal:   component.Marshal[testTextSensorConfig],
		Unmarshal: component.Unmarshal[testTextSensorConfig],
	}
}

var _ component.Declaration = (*testTextSensorDeclaration)(nil)

func TestBaseTextSensorComponent(t *testing.T) {
	reg := registry.NewRegistry()
	reg.RegisterEntityComponent(entity.DomainTypeTextSensor, "test", &testTextSensorDeclaration{})
}
