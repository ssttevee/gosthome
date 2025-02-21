package button

import (
	"context"

	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
)

type BaseButtonConfig[T any, PT interface {
	*T
	component.Component
	entity.Button
}] struct {
	component.ConfigOf[T, PT]
	entity.EntityConfig                                                                `yaml:",inline"`
	entity.DeviceClassMixinConfig[entity.ButtonDeviceClass, *entity.ButtonDeviceClass] `yaml:",inline"`
	entity.IconMixinConfig                                                             `yaml:",inline"`
}

func (bsc *BaseButtonConfig[T, PT]) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		bsc.EntityConfig.ValidateWithContext(ctx),
		bsc.DeviceClassMixinConfig.ValidateWithContext(ctx),
		bsc.IconMixinConfig.ValidateWithContext(ctx))
}

type BaseButton[T any, PT interface {
	*T
	component.Component
	entity.Button
}] struct {
	entity.BaseEntity
	entity.DeviceClassMixin[entity.ButtonDeviceClass, *entity.ButtonDeviceClass]
	entity.IconMixin
}

func NewBaseButton[T any, PT interface {
	*T
	component.Component
	entity.Button
}](ctx context.Context, t PT, cfg *BaseButtonConfig[T, PT]) (BaseButton[T, PT], error) {
	return BaseButton[T, PT]{
		BaseEntity:       entity.NewBaseEntity(entity.DomainTypeButton, &cfg.EntityConfig),
		DeviceClassMixin: entity.NewDeviceClassMixin(&cfg.DeviceClassMixinConfig),
		IconMixin:        entity.NewIconMixin(&cfg.IconMixinConfig),
	}, nil
}
