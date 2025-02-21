package entity

import (
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
	"github.com/gosthome/gosthome/core/util"
)

//go:generate go-enum

// ENUM(ok,warning,error)
type EntityStatus uint8

type Entity interface {
	cid.Identifier
	Name() string
	Internal() bool
	DisabledByDefault() bool
	EntityCategory() Category
}

type EntityComponent interface {
	Entity
	component.Component
}

type WithIcon interface {
	Icon() string
}

type WithDeviceClass[Enum any, PE interface {
	DeviceClassValues
	*Enum
}] interface {
	DeviceClass() Enum
}

type WithUnitOfMeasurement interface {
	UnitOfMeasurement() string
}

type WithState[T any] interface {
	State() T
}

type BaseEntity struct {
	cid.CID
	idhash   uint32
	category Category
	name     string

	internal    bool
	idhashReady bool

	disabledByDefault bool
}

func NewBaseEntity(t DomainType, cfg *EntityConfig) BaseEntity {
	id := cfg.ID
	if id == "" {
		if cfg.Name != "" {
			id = util.CleanString(util.SnakeCase(cfg.Name))
		} else {
			id = cid.MakeStringID(t.String())
		}
	}
	b := BaseEntity{
		CID:               cid.NewID(id),
		name:              cfg.Name,
		internal:          cfg.Internal,
		category:          cfg.Category,
		disabledByDefault: cfg.DisabledByDefault,
	}
	return b
}

// Name implements Entity.
func (b *BaseEntity) Name() string {
	return b.name
}

func (b *BaseEntity) SetName(name string) {
	b.name = name
}

// Internal implements Entity.
func (b *BaseEntity) Internal() bool {
	return b.internal
}

func (b *BaseEntity) SetInternal(internal bool) {
	b.internal = internal
}

// DisabledByDefault implements Entity.
func (b *BaseEntity) DisabledByDefault() bool {
	return b.disabledByDefault
}

func (b *BaseEntity) SetDisabledByDefault(disabledByDefault bool) {
	b.disabledByDefault = disabledByDefault
}

// EntityCategory implements Entity.
func (b *BaseEntity) EntityCategory() Category {
	return b.category
}

func (b *BaseEntity) SetEntityCategory(entityCategory Category) {
	b.category = entityCategory
}

var _ Entity = (*BaseEntity)(nil)

type IconMixin struct {
	icon string
}

func NewIconMixin(cfg *IconMixinConfig) IconMixin {
	return IconMixin{
		icon: cfg.Icon,
	}
}

func (b *IconMixin) SetIcon(icon string) {
	b.icon = icon
}

func (b *IconMixin) Icon() string {
	return b.icon
}

var _ WithIcon = (*IconMixin)(nil)

type DeviceClassMixin[Enum any, PE interface {
	DeviceClassValues
	*Enum
}] struct {
	deviceClass Enum
}

func NewDeviceClassMixin[Enum any, PE interface {
	DeviceClassValues
	*Enum
}](cfg *DeviceClassMixinConfig[Enum, PE]) DeviceClassMixin[Enum, PE] {
	return DeviceClassMixin[Enum, PE]{
		deviceClass: cfg.DeviceClass,
	}
}

func (b *DeviceClassMixin[Enum, PE]) SetDeviceClass(deviceClass Enum) {
	b.deviceClass = deviceClass
}

func (b *DeviceClassMixin[Enum, PE]) DeviceClass() Enum {
	return b.deviceClass
}

type UnitOfMeasurementMixin struct {
	unitOfMeasurement string
}

func NewUnitOfMeasurementMixin(cfg *UnitOfMeasurementMixinConfig) UnitOfMeasurementMixin {
	return UnitOfMeasurementMixin{
		unitOfMeasurement: cfg.UnitOfMeasurement,
	}
}

func (b *UnitOfMeasurementMixin) SetUnitOfMeasurement(unitOfMeasurement string) {
	b.unitOfMeasurement = unitOfMeasurement
}

func (b *UnitOfMeasurementMixin) UnitOfMeasurement() string {
	return b.unitOfMeasurement
}

var _ WithUnitOfMeasurement = (*UnitOfMeasurementMixin)(nil)
