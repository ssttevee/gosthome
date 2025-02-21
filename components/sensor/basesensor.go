package sensor

import (
	"context"

	"github.com/gosthome/gosthome/core/component"
	cv "github.com/gosthome/gosthome/core/configvalidation"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/state"
)

func unitsOfMeasurementOf(dc entity.SensorDeviceClass) []string {
	switch dc {
	case entity.SensorDeviceClassApparentPower:
		return []string{"VA"}
	case entity.SensorDeviceClassAqi:
		return []string{}
	case entity.SensorDeviceClassArea:
		return []string{"m²", "cm²", "km²", "mm²", "in²", "ft²", "yd²", "mi²", "ac", "ha"}
	case entity.SensorDeviceClassAtmosphericPressure:
		return []string{"cbar", "bar", "hPa", "mmHG", "inHg", "kPa", "mbar", "Pa", "psi"}
	case entity.SensorDeviceClassBattery:
		return []string{"%"}
	case entity.SensorDeviceClassBloodGlucoseConcentration:
		return []string{"mg/dL", "mmol/L"}
	case entity.SensorDeviceClassCo2:
		return []string{"ppm"}
	case entity.SensorDeviceClassCo:
		return []string{"ppm"}
	case entity.SensorDeviceClassConductivity:
		return []string{"S/cm", "mS/cm", "µS/cm"}
	case entity.SensorDeviceClassCurrent:
		return []string{"A", "mA"}
	case entity.SensorDeviceClassDataRate:
		return []string{"bit/s", "kbit/s", "Mbit/s", "Gbit/s", "B/s", "kB/s", "MB/s", "GB/s", "KiB/s", "MiB/s", "GiB/s"}
	case entity.SensorDeviceClassDataSize:
		return []string{"bit", "kbit", "Mbit", "Gbit", "B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	case entity.SensorDeviceClassDate:
		return []string{}
	case entity.SensorDeviceClassDistance:
		return []string{"km", "m", "cm", "mm", "mi", "nmi", "yd", "in"}
	case entity.SensorDeviceClassDuration:
		return []string{"d", "h", "min", "s", "ms"}
	case entity.SensorDeviceClassEnergy:
		return []string{"J", "kJ", "MJ", "GJ", "mWh", "Wh", "kWh", "MWh", "GWh", "TWh", "cal", "kcal", "Mcal", "Gcal"}
	case entity.SensorDeviceClassEnergyDistance:
		return []string{"kWh/100km", "mi/kWh", "km/kWh"}
	case entity.SensorDeviceClassEnergyStorage:
		return []string{"J", "kJ", "MJ", "GJ", "mWh", "Wh", "kWh", "MWh", "GWh", "TWh", "cal", "kcal", "Mcal", "Gcal"}
	case entity.SensorDeviceClassEnum:
		return []string{}
	case entity.SensorDeviceClassFrequency:
		return []string{"Hz", "kHz", "MHz", "GHz"}
	case entity.SensorDeviceClassGas:
		return []string{"m³", "ft³", "CCF"}
	case entity.SensorDeviceClassHumidity:
		return []string{"%"}
	case entity.SensorDeviceClassIlluminance:
		return []string{"lx"}
	case entity.SensorDeviceClassIrradiance:
		return []string{"W/m²", "BTU/(h⋅ft²)"}
	case entity.SensorDeviceClassMoisture:
		return []string{"%"}
	case entity.SensorDeviceClassMonetary:
		return []string{"ISO 4217"}
	case entity.SensorDeviceClassNitrogenDioxide:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassNitrogenMonoxide:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassNitrousOxide:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassOzone:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassPh:
		return []string{"None"}
	case entity.SensorDeviceClassPm1:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassPm25:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassPm10:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassPower:
		return []string{"mW", "W", "kW", "MW", "GW", "TW"}
	case entity.SensorDeviceClassPowerFactor:
		return []string{"%", "None"}
	case entity.SensorDeviceClassPrecipitation:
		return []string{"cm", "in", "mm"}
	case entity.SensorDeviceClassPrecipitationIntensity:
		return []string{"in/d", "in/h", "mm/d", "mm/h"}
	case entity.SensorDeviceClassPressure:
		return []string{"cbar", "bar", "hPa", "mmHg", "inHg", "kPa", "mbar", "Pa", "psi"}
	case entity.SensorDeviceClassReactivePower:
		return []string{"var"}
	case entity.SensorDeviceClassSignalStrength:
		return []string{"dB", "dBm"}
	case entity.SensorDeviceClassSoundPressure:
		return []string{"dB", "dBA"}
	case entity.SensorDeviceClassSpeed:
		return []string{"ft/s", "in/d", "in/h", "in/s", "km/h", "kn", "m/s", "mph", "mm/d", "mm/s"}
	case entity.SensorDeviceClassSulphurDioxide:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassTemperature:
		return []string{"°C", "degC", "C", "°F", "degF", "F", "K"}
	case entity.SensorDeviceClassTimestamp:
		return []string{}
	case entity.SensorDeviceClassVolatileOrganicCompounds:
		return []string{"µg/m³"}
	case entity.SensorDeviceClassVolatileOrganicCompoundsParts:
		return []string{"ppm", "ppb"}
	case entity.SensorDeviceClassVoltage:
		return []string{"V", "mV", "µV", "kV", "MV"}
	case entity.SensorDeviceClassVolume:
		return []string{"L", "mL", "gal", "fl. oz.", "m³", "ft³", "CCF"}
	case entity.SensorDeviceClassVolumeFlowRate:
		return []string{"m³/h", "ft³/min", "L/min", "gal/min", "mL/s"}
	case entity.SensorDeviceClassVolumeStorage:
		return []string{"L", "mL", "gal", "fl. oz.", "m³", "ft³", "CCF"}
	case entity.SensorDeviceClassWater:
		return []string{"L", "gal", "m³", "ft³", "CCF"}
	case entity.SensorDeviceClassWeight:
		return []string{"kg", "g", "mg", "µg", "oz", "lb", "st"}
	case entity.SensorDeviceClassWindDirection:
		return []string{"°", "deg"}
	case entity.SensorDeviceClassWindSpeed:
		return []string{"ft/s", "km/h", "kn", "m/s", "mph"}
	default:
		return []string{}
	}
}

type DeviceClassMixinConfig = entity.DeviceClassMixinConfig[entity.SensorDeviceClass, *entity.SensorDeviceClass]

type BaseSensorConfig[T any, PT interface {
	*T
	component.Component
	entity.Sensor
}] struct {
	component.ConfigOf[T, PT]
	entity.EntityConfig                 `yaml:",inline"`
	DeviceClassMixinConfig              `yaml:",inline"`
	entity.IconMixinConfig              `yaml:",inline"`
	entity.UnitOfMeasurementMixinConfig `yaml:",inline"`
}

func (bsc *BaseSensorConfig[T, PT]) ValidateWithContext(ctx context.Context) error {
	return cv.ValidateEmbedded(
		bsc.EntityConfig.ValidateWithContext(ctx),
		bsc.DeviceClassMixinConfig.ValidateWithContext(ctx),
		bsc.IconMixinConfig.ValidateWithContext(ctx),
		bsc.UnitOfMeasurementMixinConfig.ValidateWithContext(ctx),
	)
}

type BaseSensor[T any, PT interface {
	*T
	component.Component
	entity.Sensor
}] struct {
	entity.BaseEntity
	entity.DeviceClassMixin[entity.SensorDeviceClass, *entity.SensorDeviceClass]
	entity.IconMixin
	state.State_[entity.SensorState]

	accuracy          int32
	forceUpdate       bool
	lastResetType     entity.SensorLastResetType
	stateClass        entity.SensorStateClass
	unitOfMeasurement string
}

func NewBaseSensor[T any, PT interface {
	*T
	component.Component
	entity.Sensor
}](ctx context.Context, t PT, cfg *BaseSensorConfig[T, PT]) (ret BaseSensor[T, PT], err error) {
	ret.BaseEntity = entity.NewBaseEntity(entity.DomainTypeSensor, &cfg.EntityConfig)
	ret.DeviceClassMixin = entity.NewDeviceClassMixin(&cfg.DeviceClassMixinConfig)
	ret.IconMixin = entity.NewIconMixin(&cfg.IconMixinConfig)
	ret.State_, err = state.NewState(ctx, t, entity.SensorState{
		State:        0,
		MissingState: true,
	})
	if err != nil {
		return
	}
	return
}

// AccuracyDecimals implements entity.Sensor.
func (t *BaseSensor[T, PT]) AccuracyDecimals() int32 {
	return t.accuracy
}

// ForceUpdate implements entity.Sensor.
func (t *BaseSensor[T, PT]) ForceUpdate() bool {
	return t.forceUpdate
}

// LastResetType implements entity.Sensor.
func (t *BaseSensor[T, PT]) LastResetType() entity.SensorLastResetType {
	return t.lastResetType
}

// StateClass implements entity.Sensor.
func (t *BaseSensor[T, PT]) StateClass() entity.SensorStateClass {
	return t.stateClass
}

// UnitOfMeasurement implements entity.Sensor.
func (t *BaseSensor[T, PT]) UnitOfMeasurement() string {
	return t.unitOfMeasurement
}
