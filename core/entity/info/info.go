package info

import "github.com/gosthome/gosthome/core/entity"

type BinarySensor struct {
	ObjectId             string
	Key                  uint32
	Name                 string
	UniqueId             string
	DeviceClass          string
	IsStatusBinarySensor bool
	DisabledByDefault    bool
	Icon                 string
	EntityCategory       entity.Category
}
type Cover struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	AssumedState      bool
	SupportsPosition  bool
	SupportsTilt      bool
	DeviceClass       string
	DisabledByDefault bool
	Icon              string
	EntityCategory    entity.Category
	SupportsStop      bool
}
type Fan struct {
	ObjectId             string
	Key                  uint32
	Name                 string
	UniqueId             string
	SupportsOscillation  bool
	SupportsSpeed        bool
	SupportsDirection    bool
	SupportedSpeedLevels int32
	DisabledByDefault    bool
	Icon                 string
	EntityCategory       entity.Category
	SupportedPresetModes []string
}
type Light struct {
	ObjectId            string
	Key                 uint32
	Name                string
	UniqueId            string
	SupportedColorModes []entity.ColorMode
	MinMireds           float32
	MaxMireds           float32
	Effects             []string
	DisabledByDefault   bool
	Icon                string
	EntityCategory      entity.Category
}
type Sensor struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	UnitOfMeasurement string
	AccuracyDecimals  int32
	ForceUpdate       bool
	DeviceClass       string
	StateClass        entity.SensorStateClass
	LastResetType     entity.SensorLastResetType
	DisabledByDefault bool
	EntityCategory    entity.Category
}
type Switch struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	AssumedState      bool
	DisabledByDefault bool
	EntityCategory    entity.Category
	DeviceClass       string
}
type TextSensor struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	DeviceClass       string
}

//	type Services struct {
//		Name string
//		Type entity.ServiceArgType
//	}
type Services struct {
	Name string
	Key  uint32
	//Args []*entity.ServicesArgument
	// TODO: Handle Args
}
type Camera struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	DisabledByDefault bool
	Icon              string
	EntityCategory    entity.Category
}
type Climate struct {
	ObjectId                          string
	Key                               uint32
	Name                              string
	UniqueId                          string
	SupportsCurrentTemperature        bool
	SupportsTwoPointTargetTemperature bool
	SupportedModes                    []entity.ClimateMode
	VisualMinTemperature              float32
	VisualMaxTemperature              float32
	VisualTargetTemperatureStep       float32
	// for older peer versions - in new system this
	// is if CLIMATE_PRESET_AWAY exists is supported_presets
	LegacySupportsAway           bool
	SupportsAction               bool
	SupportedFanModes            []entity.ClimateFanMode
	SupportedSwingModes          []entity.ClimateSwingMode
	SupportedCustomFanModes      []string
	SupportedPresets             []entity.ClimatePreset
	SupportedCustomPresets       []string
	DisabledByDefault            bool
	Icon                         string
	EntityCategory               entity.Category
	VisualCurrentTemperatureStep float32
	SupportsCurrentHumidity      bool
	SupportsTargetHumidity       bool
	VisualMinHumidity            float32
	VisualMaxHumidity            float32
}
type Number struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	MinValue          float32
	MaxValue          float32
	Step              float32
	DisabledByDefault bool
	EntityCategory    entity.Category
	UnitOfMeasurement string
	Mode              entity.NumberMode
	DeviceClass       string
}
type Select struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	Options           []string
	DisabledByDefault bool
	EntityCategory    entity.Category
}
type Siren struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	Tones             []string
	SupportsDuration  bool
	SupportsVolume    bool
	EntityCategory    entity.Category
}
type Lock struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	AssumedState      bool
	SupportsOpen      bool
	RequiresCode      bool
	CodeFormat        string
}
type Button struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	DeviceClass       string
}
type MediaPlayer struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	SupportsPause     bool
	SupportedFormats  []*entity.MediaPlayerSupportedFormat
}
type AlarmControlPanel struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	SupportedFeatures uint32
	RequiresCode      bool
	RequiresCodeToArm bool
}
type Text struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	MinLength         uint32
	MaxLength         uint32
	Pattern           string
	Mode              entity.TextMode
}
type Date struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
}
type Time struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
}
type Event struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	DeviceClass       string
	EventTypes        []string
}
type Valve struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	DeviceClass       string
	AssumedState      bool
	SupportsPosition  bool
	SupportsStop      bool
}
type DateTime struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
}
type Update struct {
	ObjectId          string
	Key               uint32
	Name              string
	UniqueId          string
	Icon              string
	DisabledByDefault bool
	EntityCategory    entity.Category
	DeviceClass       string
}
