package uart

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gosthome/gosthome/components/button"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
	cv "github.com/gosthome/gosthome/core/configvalidation"
)

type ButtonConfig struct {
	button.BaseButtonConfig[Button, *Button] `yaml:",inline"`

	UARTID string    `yaml:"uart_id"`
	Data   *cv.Bytes `yaml:"data"`
}

func NewButtonConfig() *ButtonConfig {
	return &ButtonConfig{}
}

// Validate implements validation.Validatable.
func (c *ButtonConfig) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		c.BaseButtonConfig.ValidateWithContext(ctx),
		validation.ValidateStructWithContext(
			ctx, c, validation.Field(&c.UARTID, cv.String(cv.Optional(cv.Name()))),
		),
	)
}

var _ component.Config = (*ButtonConfig)(nil)

type Button struct {
	button.BaseButton[Button, *Button]
	b      *bus.Bus
	uartID cid.CID
	data   []byte
}

func NewButton(ctx context.Context, cfg *ButtonConfig) (retc []component.Component, err error) {
	uid := cfg.UARTID
	if uid != "" {
		uid = COMPONENT_KEY
	}
	ret := &Button{
		uartID: cid.NewID(uid),
	}
	ret.BaseButton, err = button.NewBaseButton(ctx, ret, &cfg.BaseButtonConfig)
	if err != nil {
		return nil, err
	}
	ret.b = bus.Get(ctx)
	ret.data = cfg.Data.Data
	return []component.Component{ret}, nil
}

// Setup implements component.Component.
func (b *Button) Setup() {

}

func (b *Button) Press(ctx context.Context) error {
	b.b.CallService(&UARTWrite{
		Key:  b.uartID.HashID(),
		Data: b.data,
	})
	return nil
}

// Close implements component.Component.
func (b *Button) Close() error {

	return nil
}

// InitializationPriority implements component.Component.
func (c *Button) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityHardware
}

var _ component.Component = (*Button)(nil)
