package entity

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
)

//go:generate go-enum --names --values --marshal

// ENUM(
// binary_sensor,
// cover,
// fan,
// light,
// sensor,
// switch,
// button,
// text_sensor,
// service,
// camera,
// climate,
// number,
// datetime_date,
// datetime_time,
// datetime_datetime,
// text,
// select,
// lock,
// valve,
// media_player,
// alarm_control_panel,
// siren,
// event,
// update,
// )
type DomainType uint

const (
	DomainTypeStart DomainType = DomainTypeBinarySensor
	DomainTypeEnd   DomainType = DomainTypeUpdate
)

type DomainTyper interface {
	DomainType() DomainType
}

// ==================	BinarySensor		=============================================

type BinarySensorDomain struct {
	BaseDomain[BinarySensorDomain, BinarySensor, *BinarySensorDomain]
}

func (*BinarySensorDomain) DomainType() DomainType {
	return DomainTypeBinarySensor
}

type BinarySensorState struct {
	State   bool
	Missing bool
}

// ENUM(
// battery,
// battery_charging,
// carbon_monoxide,
// cold,
// connectivity,
// door,
// empty,
// garage_door,
// gas,
// heat,
// light,
// lock,
// moisture,
// motion,
// moving,
// occupancy,
// opening,
// plug,
// power,
// presence,
// problem,
// running,
// safety,
// smoke,
// sound,
// tamper,
// update,
// vibration,
// window,
// )
type BinarySensorDeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *BinarySensorDeviceClass) DeviceClassValues() []string {
	return BinarySensorDeviceClassNames()
}

var _ (DeviceClassValues) = (*BinarySensorDeviceClass)(nil)

type BinarySensor interface {
	EntityComponent
	WithState[BinarySensorState]
	WithIcon
	WithDeviceClass[BinarySensorDeviceClass, *BinarySensorDeviceClass]
	IsStatusBinarySensor() bool
}

// ==================	Cover		=============================================

type CoverDomain struct {
	BaseDomain[CoverDomain, Cover, *CoverDomain]
}

func (*CoverDomain) DomainType() DomainType {
	return DomainTypeCover
}

// ENUM(open,closed)
type LegacyCoverState int32

type CoverState struct {
	LegacyState LegacyCoverState
	Position    float32
	Tilt        float32
}

// ENUM(
// awning,	//Control of an awning, such as an exterior retractible window, door, or patio cover.
// blind,	//Control of blinds, which are linked slats that expand or collapse to cover an opening or may be tilted to partially cover an opening, such as window blinds.
// curtain,	//Control of curtains or drapes, which is often fabric hung above a window or door that can be drawn open.
// damper,	//Control of a mechanical damper that reduces air flow, sound, or light.
// door,	//Control of a door that provides access to an area which is typically part of a structure.
// garage,	//Control of a garage door that provides access to a garage.
// gate,	//Control of a gate that provides access to a driveway or other area. Gates are found outside of a structure and are typically part of a fence.
// shade,	//Control of shades, which are a continuous plane of material or connected cells that expanded or collapsed over an opening, such as window shades.
// shutter,	//Control of shutters, which are linked slats that swing out/in to cover an opening or may be tilted to partially cover an opening, such as indoor or exterior window shutters.
// window,	//Control of a physical window that opens and closes or may tilt.
// )
type CoverDeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *CoverDeviceClass) DeviceClassValues() []string {
	return CoverDeviceClassNames()
}

var _ (DeviceClassValues) = (*CoverDeviceClass)(nil)

type Cover interface {
	EntityComponent
	WithState[CoverState]
	WithIcon
	WithDeviceClass[CoverDeviceClass, *CoverDeviceClass]
}

// ==================	Fan		=============================================

type FanDomain struct {
	BaseDomain[FanDomain, Fan, *FanDomain]
}

func (*FanDomain) DomainType() DomainType {
	return DomainTypeFan
}

// ENUM(low,medium,high)
type FanSpeed int32

// ENUM(forward.reverse)
type FanDirection int32

type FanState struct {
	State       bool
	Oscillating bool

	Speed      FanSpeed
	Direction  FanDirection
	SpeedLevel int32
	PresetMode string
}

type Fan interface {
	EntityComponent
	WithState[FanState]
	WithIcon
}

// ==================	Light		=============================================

type LightDomain struct {
	BaseDomain[LightDomain, Light, *LightDomain]
}

func (*LightDomain) DomainType() DomainType {
	return DomainTypeLight
}

type LightState struct {
	State            bool
	Brightness       float32
	ColorMode        ColorMode
	ColorBrightness  float32
	Red              float32
	Green            float32
	Blue             float32
	White            float32
	ColorTemperature float32
	ColdWhite        float32
	WarmWhite        float32
	Effect           string
}

// Color capabilities are the various outputs that a light has and that can be independently controlled by the user.
// ENUM(
// )
type ColorCapability int32

const (
	// Light can be turned on/off.
	ColorCapabilityOnOff ColorCapability = 1 << iota
	// Master brightness of the light can be controlled.
	ColorCapabilityBrightness
	// Brightness of white channel can be controlled separately from other channels.
	ColorCapabilityWhite
	// Color temperature can be controlled.
	ColorCapabilityColorTemperature
	// Brightness of cold and warm white output can be controlled.
	ColorCapabilityColdWarmWhite
	// Color can be controlled using RGB format (includes a brightness control for the color).
	ColorCapabilityRgb
)

type ColorMode ColorCapability

const (
	// No color mode configured (cannot be a supported mode, only active when light is off).
	ColorModeUnknown = 0
	// Only on/off control.
	ColorModeOnOff = ColorCapabilityOnOff
	// Dimmable light.
	ColorModeBrightness = (ColorCapabilityOnOff | ColorCapabilityBrightness)
	// White output only (use only if the light also has another color mode such as RGB).
	ColorModeWhite = (ColorCapabilityOnOff | ColorCapabilityBrightness | ColorCapabilityWhite)
	// Controllable color temperature output.
	ColorTemperature = (ColorCapabilityOnOff | ColorCapabilityBrightness | ColorCapabilityColorTemperature)
	// Cold and warm white output with individually controllable brightness.
	ColdWarmWhite = (ColorCapabilityOnOff | ColorCapabilityBrightness | ColorCapabilityColdWarmWhite)
	// RGB color output.
	ColorModeRgb = (ColorCapabilityOnOff | ColorCapabilityBrightness | ColorCapabilityRgb)
	// RGB color output and a separate white output.
	RgbWhite = ColorCapabilityOnOff | ColorCapabilityBrightness | ColorCapabilityRgb | ColorCapabilityWhite
	// RGB color output and a separate white output with controllable color temperature.
	ColorModeRgbColorTemperature = (ColorCapabilityOnOff | ColorCapabilityBrightness | ColorCapabilityRgb |
		ColorCapabilityWhite | ColorCapabilityColorTemperature)
	// RGB color output and separate cold and warm white outputs.
	ColorModeRgbColdWarmWhite = (ColorCapabilityOnOff | ColorCapabilityBrightness | ColorCapabilityRgb |
		ColorCapabilityColdWarmWhite)
)

type Light interface {
	EntityComponent
	WithState[LightState]
	WithIcon
	SupportedColorModes() []ColorMode
	Effects() []string
	MinMireds() float32
	MaxMireds() float32

	Command(LightCommand) error
}

// ==================	Sensor		=============================================

type SensorDomain struct {
	BaseDomain[SensorDomain, Sensor, *SensorDomain]
}

func (*SensorDomain) DomainType() DomainType {
	return DomainTypeSensor
}

type SensorState struct {
	State        float32
	MissingState bool
}

// ENUM(
// none,
// measurement, // The state represents a measurement in present time, not a historical aggregation such as statistics or a prediction of the future. Examples of what should be classified measurement are: current temperature, humidity or electric power. Examples of what should not be classified as measurement: Forecasted temperature for tomorrow, yesterday's energy consumption or anything else that doesn't include the current measurement. For supported sensors, statistics of hourly min, max and average sensor readings is updated every 5 minutes.
// total, //The state represents a total amount that can both increase and decrease, e.g. a net energy meter. Statistics of the accumulated growth or decline of the sensor's value since it was first added is updated every 5 minutes. This state class should not be used for sensors where the absolute value is interesting instead of the accumulated growth or decline, for example remaining battery capacity or CPU load; in such cases state class measurement should be used instead.
// total_increasing, //Similar to total, with the restriction that the state represents a monotonically increasing positive total which periodically restarts counting from 0, e.g. a daily amount of consumed gas, weekly water consumption or lifetime energy consumption. Statistics of the accumulated growth of the sensor's value since it was first added is updated every 5 minutes. A decreasing value is interpreted as the start of a new meter cycle or the replacement of the meter.
// )
type SensorStateClass int32

// ENUM(
// none,
// never,
// auto,
// )
type SensorLastResetType int32

// ENUM(
// apparent_power, // Apparent power
// aqi, // Air Quality Index
// area, // Area
// atmospheric_pressure, // Atmospheric pressure
// battery, // Percentage of battery that is left
// blood_glucose_concentration, // Blood glucose concentration```
// co2, // Concentration of carbon dioxide.
// co, // Concentration of carbon monoxide.
// conductivity, // Conductivity
// current, // Current
// data_rate, // Data rate
// data_size, // Data size
// date, // Date. Requires native_value to be a Python datetime.date object, or None.
// distance, // Generic distance
// duration, // Time period. Should not update only due to time passing. The device or service needs to give a new data point to update.
// energy, // Energy, this device class should be used for sensors representing energy consumption, for example an electricity meter. Represents power over time. Not to be confused with power.
// energy_distance, // Energy per distance, this device class should be used to represent energy consumption by distance, for example the amount of electric energy consumed by an electric car.
// energy_storage, // Stored energy, this device class should be used for sensors representing stored energy, for example the amount of electric energy currently stored in a battery or the capacity of a battery. Represents power over time. Not to be confused with power.
// enum, // The sensor has a limited set of (non-numeric) states. The options property must be set to a list of possible states when using this device class.
// frequency, // Frequency
// gas, // Volume of gas. Gas consumption measured as energy in kWh instead of a volume should be classified as energy.
// humidity, // Relative humidity
// illuminance, // Light level
// irradiance, // Irradiance
// moisture, // Moisture
// monetary, // Monetary value with a currency.
// nitrogen_dioxide, // Concentration of nitrogen dioxide
// nitrogen_monoxide, // Concentration of nitrogen monoxide
// nitrous_oxide, // Concentration of nitrous oxide
// ozone, // Concentration of ozone
// ph, // Potential hydrogen (pH) of an aqueous solution
// pm1, // Concentration of particulate matter less than 1 micrometer
// pm25, // Concentration of particulate matter less than 2.5 micrometers
// pm10, // Concentration of particulate matter less than 10 micrometers
// power, // Power.
// power_factor, // Power Factor
// precipitation, // Accumulated precipitation
// precipitation_intensity, // Precipitation intensity
// pressure, // Pressure.
// reactive_power, // Reactive power
// signal_strength, // Signal strength
// sound_pressure, // Sound pressure
// speed, // Generic speed
// sulphur_dioxide, // Concentration of sulphure dioxide
// temperature, // Temperature.
// timestamp, // Timestamp. Requires native_value to return a Python datetime.datetime object, with time zone information, or None.
// volatile_organic_compounds, // Concentration of volatile organic compounds
// volatile_organic_compounds_parts, // Ratio of volatile organic compounds
// voltage, // Voltage
// volume, // Generic volume, this device class should be used for sensors representing a consumption, for example the amount of fuel consumed by a vehicle.
// volume_flow_rate, // Volume flow rate, this device class should be used for sensors representing a flow of some volume, for example the amount of water consumed momentarily.
// volume_storage, // Generic stored volume, this device class should be used for sensors representing a stored volume, for example the amount of fuel in a fuel tank.
// water, // Water consumption
// weight, // Generic mass; weight is used instead of mass to fit with every day language.
// wind_direction, // Wind direction
// wind_speed, // Wind speed
// )
type SensorDeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *SensorDeviceClass) DeviceClassValues() []string {
	return SensorDeviceClassNames()
}

var _ (DeviceClassValues) = (*SensorDeviceClass)(nil)

type Sensor interface {
	EntityComponent
	WithState[SensorState]
	WithIcon
	WithDeviceClass[SensorDeviceClass, *SensorDeviceClass]
	WithUnitOfMeasurement
	AccuracyDecimals() int32
	ForceUpdate() bool
	StateClass() SensorStateClass
	LastResetType() SensorLastResetType
}

// ==================	Switch		=============================================

type SwitchDomain struct {
	BaseDomain[SwitchDomain, Switch, *SwitchDomain]
}

func (*SwitchDomain) DomainType() DomainType {
	return DomainTypeSwitch
}

type SwitchState struct {
	State bool
}

// ENUM(
// outlet, // Device is an outlet for power.
// switch, // Device is switch for some type of entity.
// )
type SwitchDeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *SwitchDeviceClass) DeviceClassValues() []string {
	return SwitchDeviceClassNames()
}

var _ (DeviceClassValues) = (*SwitchDeviceClass)(nil)

type Switch interface {
	EntityComponent
	WithState[SwitchState]
	WithIcon
	WithDeviceClass[SwitchDeviceClass, *SwitchDeviceClass]
}

// ==================	Button		=============================================

type ButtonDomain struct {
	BaseDomain[ButtonDomain, Button, *ButtonDomain]
}

func (*ButtonDomain) DomainType() DomainType {
	return DomainTypeButton
}

// ENUM(
// identify,//The button is used to identify a device.
// restart,//The button restarts the device.
// update,//The button updates the software of the device.
// )
type ButtonDeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *ButtonDeviceClass) DeviceClassValues() []string {
	return ButtonDeviceClassNames()
}

var _ (DeviceClassValues) = (*ButtonDeviceClass)(nil)

type Button interface {
	EntityComponent
	WithIcon
	WithDeviceClass[ButtonDeviceClass, *ButtonDeviceClass]
	Press(ctx context.Context) error
}

// ==================	TextSensor		=============================================

type TextSensorDomain struct {
	BaseDomain[TextSensorDomain, TextSensor, *TextSensorDomain]
}

func (*TextSensorDomain) DomainType() DomainType {
	return DomainTypeTextSensor
}

type TextSensorState struct {
	State        string
	MissingState bool
}

type TextSensorDeviceClass = SensorDeviceClass

type TextSensor interface {
	EntityComponent
	WithState[TextSensorState]
	WithIcon
	WithDeviceClass[TextSensorDeviceClass, *TextSensorDeviceClass]
}

// ==================	Service		=============================================

type ServiceDomain struct {
	BaseDomain[ServiceDomain, Service, *ServiceDomain]
}

func (*ServiceDomain) DomainType() DomainType {
	return DomainTypeService
}

type Service interface {
	EntityComponent
}

// ==================	Camera		=============================================

type CameraDomain struct {
	BaseDomain[CameraDomain, Camera, *CameraDomain]
}

func (*CameraDomain) DomainType() DomainType {
	return DomainTypeCamera
}

type Camera interface {
	EntityComponent
	WithIcon
}

// ==================	Climate		=============================================

type ClimateDomain struct {
	BaseDomain[ClimateDomain, Climate, *ClimateDomain]
}

func (*ClimateDomain) DomainType() DomainType {
	return DomainTypeClimate
}

// ENUM(off,heat_cool,cool,heat,fan_only,dry,auto)
type ClimateMode int32

// ENUM(on,off,auto,low,medium,high,middle,focus,diffuse,quiet)
type ClimateFanMode int32

// ENUM(off,both,vertical,horizontal)
type ClimateSwingMode int32

// ENUM(off,cooling,heating,idle,drying,fan)
type ClimateAction int32

// ENUM(none,home,away,boost,comfort,eco,sleep,activity)
type ClimatePreset int32

type ClimateState struct {
	Mode                  ClimateMode
	CurrentTemperature    float32
	TargetTemperature     float32
	TargetTemperatureLow  float32
	TargetTemperatureHigh float32
	Action                ClimateAction
	FanMode               ClimateFanMode
	SwingMode             ClimateSwingMode
	CustomFanMode         string
	Preset                ClimatePreset
	CustomPreset          string
	CurrentHumidity       float32
	TargetHumidity        float32
}

type Climate interface {
	EntityComponent
	WithState[ClimateState]
	WithIcon
}

// ENUM(auto,box,slider)
type NumberMode int32

// ==================	Number		=============================================

type NumberDomain struct {
	BaseDomain[NumberDomain, Number, *NumberDomain]
}

func (*NumberDomain) DomainType() DomainType {
	return DomainTypeNumber
}

type NumberState struct {
	State        float32
	MissingState bool
}

type NumberDeviceClass = SensorDeviceClass

type Number interface {
	EntityComponent
	WithState[NumberState]
	WithIcon
	WithDeviceClass[NumberDeviceClass, *NumberDeviceClass]
	WithUnitOfMeasurement
	NumberMode() NumberMode
}

// ==================	Date		=============================================

type DateDomain struct {
	BaseDomain[DateDomain, Date, *DateDomain]
}

func (*DateDomain) DomainType() DomainType {
	return DomainTypeDatetimeDate
}

type DateState struct {
	MissingState bool
	Year         uint32
	Month        uint32
	Day          uint32
}

type Date interface {
	EntityComponent
	WithState[DateState]
	WithIcon
}

// ==================	Time		=============================================

type TimeDomain struct {
	BaseDomain[TimeDomain, Time, *TimeDomain]
}

func (*TimeDomain) DomainType() DomainType {
	return DomainTypeDatetimeTime
}

type TimeState struct {
	MissingState bool
	Hour         uint32
	Minute       uint32
	Second       uint32
}

type Time interface {
	EntityComponent
	WithState[TimeState]
	WithIcon
}

// ==================	Datetime		=============================================

type DatetimeDomain struct {
	BaseDomain[DatetimeDomain, Datetime, *DatetimeDomain]
}

func (*DatetimeDomain) DomainType() DomainType {
	return DomainTypeDatetimeDatetime
}

type DatetimeState struct {
	MissingState bool
	EpochSeconds uint32
}

type Datetime interface {
	EntityComponent
	WithState[DatetimeState]
	WithIcon
}

// ==================	Text		=============================================

type TextDomain struct {
	BaseDomain[TextDomain, Text, *TextDomain]
}

func (*TextDomain) DomainType() DomainType {
	return DomainTypeText
}

type TextState struct {
	State        string
	MissingState bool
}

// ENUM(text,password)
type TextMode int32

type Text interface {
	EntityComponent
	WithState[TextState]
	WithIcon
	TextMode() TextMode
}

// ==================	Select		=============================================

type SelectDomain struct {
	BaseDomain[SelectDomain, Select, *SelectDomain]
}

func (*SelectDomain) DomainType() DomainType {
	return DomainTypeSelect
}

type SelectState struct {
	State        string
	MissingState bool
}

type Select interface {
	EntityComponent
	WithState[SelectState]
	WithIcon
	Values() []string

	Command(value string) error
}

// ==================	Siren		=============================================

type SirenDomain struct {
	BaseDomain[SirenDomain, Siren, *SirenDomain]
}

func (*SirenDomain) DomainType() DomainType {
	return DomainTypeSiren
}

type SirenState bool

type Siren interface {
	EntityComponent
	WithState[SirenState]
	WithIcon
	Tones() []string
	SupportsDuration() bool
	SupportsVolume() bool
}

// ==================	Lock		=============================================

type LockDomain struct {
	BaseDomain[LockDomain, Lock, *LockDomain]
}

func (*LockDomain) DomainType() DomainType {
	return DomainTypeLock
}

// ENUM(none,locked,unlocked,jammed,locking,unlocking)
type LockState int32

type Lock interface {
	EntityComponent
	WithState[LockState]
	WithIcon
}

// ==================	Valve		=============================================

type ValveDomain struct {
	BaseDomain[ValveDomain, Valve, *ValveDomain]
}

func (*ValveDomain) DomainType() DomainType {
	return DomainTypeValve
}

// ENUM(idle,is_opening,is_closing)
type ValveOperation int32

type ValveState struct {
	Position         float32
	CurrentOperation ValveOperation
}

// ENUM(
// water,	// Control of a water valve.
// gas,		// Control of a gas valve.
// )
type ValveDeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *ValveDeviceClass) DeviceClassValues() []string {
	return ValveDeviceClassNames()
}

var _ (DeviceClassValues) = (*ValveDeviceClass)(nil)

type Valve interface {
	EntityComponent
	WithState[ValveState]
	WithIcon
	WithDeviceClass[ValveDeviceClass, *ValveDeviceClass]
}

// ==================	MediaPlayer		=============================================

type MediaPlayerDomain struct {
	BaseDomain[MediaPlayerDomain, MediaPlayer, *MediaPlayerDomain]
}

func (*MediaPlayerDomain) DomainType() DomainType {
	return DomainTypeMediaPlayer
}

// ENUM(none,idle,playing,paused)
type MediaPlayingState int32

type MediaPlayerState struct {
	State  MediaPlayingState
	Volume float32
	Muted  bool
}

// ENUM(default,announcement)
type MediaPlayerFormatPurpose int32

type MediaPlayerSupportedFormat struct {
	Format      string
	SampleRate  uint32
	NumChannels uint32
	Purpose     MediaPlayerFormatPurpose
	SampleBytes uint32
}

type MediaPlayer interface {
	EntityComponent
	WithState[MediaPlayerState]
	WithIcon
}

// ==================	AlarmControlPanel		=============================================

type AlarmControlPanelDomain struct {
	BaseDomain[AlarmControlPanelDomain, AlarmControlPanel, *AlarmControlPanelDomain]
}

func (*AlarmControlPanelDomain) DomainType() DomainType {
	return DomainTypeAlarmControlPanel
}

// ENUM(disarmed,armed_home,armed_away,armed_night,armed_vacation,armed_custom_bypass,pending,arming,disarming,triggered)
type AlarmControlPanelState int32

type AlarmControlPanel interface {
	EntityComponent
	WithState[AlarmControlPanelState]
	WithIcon
}

// ==================	Event		=============================================

type EventDomain struct {
	BaseDomain[EventDomain, Event, *EventDomain]
}

func (*EventDomain) DomainType() DomainType {
	return DomainTypeEvent
}

type EventDeviceClass string

var _EventDeviceClassNamesValues = sync.OnceValue(func() map[string]EventDeviceClass {
	a := slices.Clone(BinarySensorDeviceClassNames())
	a = append(a, CoverDeviceClassNames()...)
	a = append(a, SwitchDeviceClassNames()...)
	a = append(a, ButtonDeviceClassNames()...)
	a = append(a, SensorDeviceClassNames()...)
	a = append(a, ValveDeviceClassNames()...)
	return maps.Collect[string, EventDeviceClass](func(yield func(string, EventDeviceClass) bool) {
		for _, e := range a {
			if !yield(e, EventDeviceClass(e)) {
				return
			}
		}
	})
})

type errInvalidEventDeviceClass func() error

func (e errInvalidEventDeviceClass) Error() string {
	return e().Error()
}

var ErrInvalidEventDeviceClass = errInvalidEventDeviceClass(sync.OnceValue(func() error {
	return fmt.Errorf("not a valid EventDeviceClassNames, try [%s]", strings.Join(EventDeviceClassNames(), ", "))
}))

var EventDeviceClassNames = sync.OnceValue(func() []string {
	return slices.Collect(maps.Keys(_EventDeviceClassNamesValues()))
})

// ParseBinarySensorDeviceClass attempts to convert a string to a BinarySensorDeviceClass.
func ParseEventDeviceClass(name string) (EventDeviceClass, error) {
	if x, ok := _EventDeviceClassNamesValues()[name]; ok {
		return x, nil
	}
	return EventDeviceClass(""), fmt.Errorf("%s is %w", name, ErrInvalidEventDeviceClass())
}

func (x EventDeviceClass) String() string {
	return string(x)
}
func (x EventDeviceClass) IsValid() bool {
	_, err := ParseBinarySensorDeviceClass(string(x))
	return err == nil
}
func (x EventDeviceClass) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}
func (x *EventDeviceClass) UnmarshalText(text []byte) error {
	tmp, err := ParseEventDeviceClass(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// DeviceClassValues implements entity.DeviceClassValues.
func (x *EventDeviceClass) DeviceClassValues() []string {
	return EventDeviceClassNames()
}

var _ (DeviceClassValues) = (*EventDeviceClass)(nil)

type Event interface {
	EntityComponent
	WithIcon
	WithDeviceClass[EventDeviceClass, *EventDeviceClass]
}

// ==================	Update		=============================================

type UpdateDomain struct {
	BaseDomain[UpdateDomain, Update, *UpdateDomain]
}

func (*UpdateDomain) DomainType() DomainType {
	return DomainTypeUpdate
}

type UpdateState struct {
	MissingState   bool
	InProgress     bool
	HasProgress    bool
	Progress       float32
	CurrentVersion string
	LatestVersion  string
	Title          string
	ReleaseSummary string
	ReleaseUrl     string
}

// ENUM(,firmware)
type UpdateDeviceClass string

// DeviceClassValues implements entity.DeviceClassValues.
func (x *UpdateDeviceClass) DeviceClassValues() []string {
	return UpdateDeviceClassNames()
}

var _ (DeviceClassValues) = (*UpdateDeviceClass)(nil)

type Update interface {
	EntityComponent
	WithState[UpdateState]
	WithIcon
	WithDeviceClass[UpdateDeviceClass, *UpdateDeviceClass]
}
