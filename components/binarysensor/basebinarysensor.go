package binarysensor

import (
	"context"

	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/state"
)

type BaseBinarySensorConfig[T any, PT interface {
	*T
	component.Component
	entity.BinarySensor
}] struct {
	component.ConfigOf[T, PT]
	entity.EntityConfig                                                                            `yaml:",inline"`
	entity.DeviceClassMixinConfig[entity.BinarySensorDeviceClass, *entity.BinarySensorDeviceClass] `yaml:",inline"`
	entity.IconMixinConfig                                                                         `yaml:",inline"`
}

func (bsc *BaseBinarySensorConfig[T, PT]) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		bsc.EntityConfig.ValidateWithContext(ctx),
		bsc.DeviceClassMixinConfig.ValidateWithContext(ctx),
		bsc.IconMixinConfig.ValidateWithContext(ctx),
	)
}

type BaseBinarySensor[T any, PT interface {
	*T
	component.Component
	entity.BinarySensor
}] struct {
	entity.BaseEntity
	entity.DeviceClassMixin[entity.BinarySensorDeviceClass, *entity.BinarySensorDeviceClass]
	entity.IconMixin
	state.State_[entity.BinarySensorState]
	isStatusBinarySensor bool
}

// IsStatusBinarySensor implements entity.BinarySensor.
func (t *BaseBinarySensor[T, PT]) IsStatusBinarySensor() bool {
	return t.isStatusBinarySensor
}

func (t *BaseBinarySensor[T, PT]) SetStatusBinarySensor(b bool) {
	t.isStatusBinarySensor = b
}

func NewBaseBinarySensor[T any, PT interface {
	*T
	component.Component
	entity.BinarySensor
}](ctx context.Context, t PT, cfg *BaseBinarySensorConfig[T, PT]) (ret BaseBinarySensor[T, PT], err error) {
	ret.BaseEntity = entity.NewBaseEntity(entity.DomainTypeBinarySensor, &cfg.EntityConfig)
	ret.DeviceClassMixin = entity.NewDeviceClassMixin(&cfg.DeviceClassMixinConfig)
	ret.IconMixin = entity.NewIconMixin(&cfg.IconMixinConfig)
	ret.State_, err = state.NewState(ctx, t, entity.BinarySensorState{
		State:   false,
		Missing: true,
	})
	return
}
