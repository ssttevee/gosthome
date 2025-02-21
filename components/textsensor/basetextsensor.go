package textsensor

import (
	"context"

	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/state"
)

type BaseTextSensorConfig[T any, PT interface {
	*T
	component.Component
	entity.TextSensor
}] struct {
	component.ConfigOf[T, PT]
	entity.EntityConfig                                                                `yaml:",inline"`
	entity.DeviceClassMixinConfig[entity.SensorDeviceClass, *entity.SensorDeviceClass] `yaml:",inline"`
	entity.IconMixinConfig                                                             `yaml:",inline"`
}

func (bsc *BaseTextSensorConfig[T, PT]) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		bsc.EntityConfig.ValidateWithContext(ctx),
		bsc.DeviceClassMixinConfig.ValidateWithContext(ctx),
		bsc.IconMixinConfig.ValidateWithContext(ctx),
	)
}

type BaseTextSensor[T any, PT interface {
	*T
	component.Component
	entity.TextSensor
}] struct {
	entity.BaseEntity
	entity.DeviceClassMixin[entity.SensorDeviceClass, *entity.SensorDeviceClass]
	entity.IconMixin
	state.State_[entity.TextSensorState]
}

func NewBaseTextSensor[T any, PT interface {
	*T
	component.Component
	entity.TextSensor
}](ctx context.Context, t PT, cfg *BaseTextSensorConfig[T, PT]) (ret BaseTextSensor[T, PT], err error) {
	ret.BaseEntity = entity.NewBaseEntity(entity.DomainTypeTextSensor, &cfg.EntityConfig)
	ret.DeviceClassMixin = entity.NewDeviceClassMixin(&cfg.DeviceClassMixinConfig)
	ret.IconMixin = entity.NewIconMixin(&cfg.IconMixinConfig)
	ret.State_, err = state.NewState(ctx, t, entity.TextSensorState{
		State:        "",
		MissingState: true,
	})
	return
}
