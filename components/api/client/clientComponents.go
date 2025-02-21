package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"weak"

	ehp "github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/entity/info"
	"github.com/majfault/signal"
)

var ErrClientGone = errors.New("client is destroyed")

type ComponentBase struct {
	c weak.Pointer[Client]
	component.WithInitializationPriorityIo
}

func (b *ComponentBase) Internal() bool {
	return false
}

type state[T any] struct {
	state       *T
	stateChange signal.Signal1[T]
}

// State implements entity.BinarySensor.
func (s *state[T]) State() T {
	if s.state == nil {
		return *new(T)
	}
	return *s.state
}

func (s *state[T]) SetState(t T) {
	s.state = &t
	s.stateChange.Emit(t)
}

func (s *state[T]) StateChange() *signal.Signal1[T] {
	return &s.stateChange
}

func (s *state[T]) Close() error {
	return s.stateChange.Close()
}

type BinarySensorComponent struct {
	ComponentBase
	state[entity.BinarySensorState]

	i info.BinarySensor
}

// IsStatusBinarySensor implements entity.BinarySensor.
func (b *BinarySensorComponent) IsStatusBinarySensor() bool {
	return b.i.IsStatusBinarySensor
}

// DeviceClass implements entity.BinarySensor.
func (b *BinarySensorComponent) DeviceClass() entity.BinarySensorDeviceClass {
	return entity.BinarySensorDeviceClass(b.i.DeviceClass)
}

// DisabledByDefault implements entity.BinarySensor.
func (b *BinarySensorComponent) DisabledByDefault() bool {
	return b.i.DisabledByDefault
}

// EntityCategory implements entity.BinarySensor.
func (b *BinarySensorComponent) EntityCategory() entity.Category {
	return b.i.EntityCategory
}

// HashID implements entity.BinarySensor.
func (b *BinarySensorComponent) HashID() uint32 {
	return b.i.Key
}

// ID implements entity.BinarySensor.
func (b *BinarySensorComponent) ID() string {
	return b.i.ObjectId
}

// Icon implements entity.BinarySensor.
func (b *BinarySensorComponent) Icon() string {
	return b.i.Icon
}

// Name implements entity.BinarySensor.
func (b *BinarySensorComponent) Name() string {
	return b.i.Name
}

// Setup implements entity.BinarySensor.
func (b *BinarySensorComponent) Setup() {}

func (b *BinarySensorComponent) UniqueID() string {
	return b.i.UniqueId
}

var _ (entity.BinarySensor) = (*BinarySensorComponent)(nil)

type CoverComponent struct {
	ComponentBase
	state[entity.CoverState]

	i info.Cover
}

// DeviceClass implements entity.Cover.
func (c *CoverComponent) DeviceClass() entity.CoverDeviceClass {
	return entity.CoverDeviceClass(c.i.DeviceClass)
}

// DisabledByDefault implements entity.Cover.
func (c *CoverComponent) DisabledByDefault() bool {
	return c.i.DisabledByDefault
}

// EntityCategory implements entity.Cover.
func (c *CoverComponent) EntityCategory() entity.Category {
	return c.i.EntityCategory
}

// HashID implements entity.Cover.
func (c *CoverComponent) HashID() uint32 {
	return c.i.Key
}

// ID implements entity.Cover.
func (c *CoverComponent) ID() string {
	return c.i.ObjectId
}

// Icon implements entity.Cover.
func (c *CoverComponent) Icon() string {
	return c.i.Icon
}

// Name implements entity.Cover.
func (c *CoverComponent) Name() string {
	return c.i.Name
}

// Setup implements entity.Cover.
func (c *CoverComponent) Setup() {}

func (c *CoverComponent) UniqueID() string {
	return c.i.UniqueId
}

var _ (entity.Cover) = (*CoverComponent)(nil)

type FanComponent struct {
	ComponentBase
	state[entity.FanState]

	i info.Fan
}

// DisabledByDefault implements entity.Fan.
func (f *FanComponent) DisabledByDefault() bool {
	return f.i.DisabledByDefault
}

// EntityCategory implements entity.Fan.
func (f *FanComponent) EntityCategory() entity.Category {
	return f.i.EntityCategory
}

// HashID implements entity.Fan.
func (f *FanComponent) HashID() uint32 {
	return f.i.Key
}

// ID implements entity.Fan.
func (f *FanComponent) ID() string {
	return f.i.ObjectId
}

// Icon implements entity.Fan.
func (f *FanComponent) Icon() string {
	return f.i.Icon
}

// Name implements entity.Fan.
func (f *FanComponent) Name() string {
	return f.i.Name
}

// Setup implements entity.Fan.
func (f *FanComponent) Setup() {}

func (f *FanComponent) UniqueID() string {
	return f.i.UniqueId
}

var _ (entity.Fan) = (*FanComponent)(nil)

type LightComponent struct {
	ComponentBase
	state[entity.LightState]

	i info.Light
}

func (l *LightComponent) Command(cmd entity.LightCommand) error {
	client := l.c.Value()
	if client == nil {
		return ErrClientGone
	}
	effect := ""
	if cmd.Effect.Has {
		ei := int(cmd.Effect.Value)
		if ei >= 0 && ei < len(l.i.Effects) {
			effect = l.i.Effects[ei]
		} else {
			cmd.Effect.Has = false
		}
	}
	slog.Info("LightCommandRequest",
		"HasState", cmd.State.Has,
		"State", cmd.State.Value,
		"HasBrightness", cmd.Brightness.Has,
		"Brightness", cmd.Brightness.Value,
		"HasColorMode", cmd.ColorMode.Has,
		"ColorMode", int32(cmd.ColorMode.Value),
		"HasColorBrightness", cmd.ColorBrightness.Has,
		"ColorBrightness", cmd.ColorBrightness.Value,
		"HasRgb", cmd.Rgb.Has,
		"Red", cmd.Rgb.Value.Red,
		"Green", cmd.Rgb.Value.Green,
		"Blue", cmd.Rgb.Value.Blue,
		"HasWhite", cmd.White.Has,
		"White", cmd.White.Value,
		"HasColorTemperature", cmd.ColorTemperature.Has,
		"ColorTemperature", cmd.ColorTemperature.Value,
		"HasColdWhite", cmd.ColdWhite.Has,
		"ColdWhite", cmd.ColdWhite.Value,
		"HasWarmWhite", cmd.WarmWhite.Has,
		"WarmWhite", cmd.WarmWhite.Value,
		"HasTransitionLength", cmd.TransitionLength.Has,
		"TransitionLength", uint32(cmd.TransitionLength.Value.Seconds()),
		"HasFlashLength", cmd.FlashLength.Has,
		"FlashLength", uint32(cmd.FlashLength.Value.Seconds()),
		"HasEffect", cmd.Effect.Has,
		"Effect", effect,
	)
	return client.sendMessages(&ehp.LightCommandRequest{
		Key:                 l.i.Key,
		HasState:            cmd.State.Has,
		State:               cmd.State.Value,
		HasBrightness:       cmd.Brightness.Has,
		Brightness:          cmd.Brightness.Value,
		HasColorMode:        cmd.ColorMode.Has,
		ColorMode:           int32(cmd.ColorMode.Value),
		HasColorBrightness:  cmd.ColorBrightness.Has,
		ColorBrightness:     cmd.ColorBrightness.Value,
		HasRgb:              cmd.Rgb.Has,
		Red:                 cmd.Rgb.Value.Red,
		Green:               cmd.Rgb.Value.Green,
		Blue:                cmd.Rgb.Value.Blue,
		HasWhite:            cmd.White.Has,
		White:               cmd.White.Value,
		HasColorTemperature: cmd.ColorTemperature.Has,
		ColorTemperature:    cmd.ColorTemperature.Value,
		HasColdWhite:        cmd.ColdWhite.Has,
		ColdWhite:           cmd.ColdWhite.Value,
		HasWarmWhite:        cmd.WarmWhite.Has,
		WarmWhite:           cmd.WarmWhite.Value,
		HasTransitionLength: cmd.TransitionLength.Has,
		TransitionLength:    uint32(cmd.TransitionLength.Value.Seconds()),
		HasFlashLength:      cmd.FlashLength.Has,
		FlashLength:         uint32(cmd.FlashLength.Value.Seconds()),
		HasEffect:           cmd.Effect.Has,
		Effect:              effect,
	})
}

// Effects implements entity.Light.
func (l *LightComponent) Effects() []string {
	return l.i.Effects
}
func (l *LightComponent) MinMireds() float32 {
	return l.i.MinMireds
}
func (l *LightComponent) MaxMireds() float32 {
	return l.i.MaxMireds
}

// Close implements entity.Light.
// Subtle: this method shadows the method (state).Close of LightComponent.state.
func (l *LightComponent) Close() error {
	return nil
}

// ColorModes implements entity.Light.
func (l *LightComponent) SupportedColorModes() []entity.ColorMode {
	return l.i.SupportedColorModes
}

func (l *LightComponent) Info() info.Light {
	return l.i
}

// DisabledByDefault implements entity.Light.
func (l *LightComponent) DisabledByDefault() bool {
	return l.i.DisabledByDefault
}

// EntityCategory implements entity.Light.
func (l *LightComponent) EntityCategory() entity.Category {
	return l.i.EntityCategory
}

// HashID implements entity.Light.
func (l *LightComponent) HashID() uint32 {
	return l.i.Key
}

// ID implements entity.Light.
func (l *LightComponent) ID() string {
	return l.i.ObjectId
}

// Icon implements entity.Light.
func (l *LightComponent) Icon() string {
	return l.i.Icon
}

// Name implements entity.Light.
func (l *LightComponent) Name() string {
	return l.i.Name
}

// Setup implements entity.Light.
func (l *LightComponent) Setup() {}

func (l *LightComponent) UniqueID() string {
	return l.i.UniqueId
}

var _ (entity.Light) = (*LightComponent)(nil)

type SensorComponent struct {
	ComponentBase
	state[entity.SensorState]

	i info.Sensor
}

// AccuracyDecimals implements entity.Sensor.
func (s *SensorComponent) AccuracyDecimals() int32 {
	return s.i.AccuracyDecimals
}

// DeviceClass implements entity.Sensor.
func (s *SensorComponent) DeviceClass() entity.SensorDeviceClass {
	return entity.SensorDeviceClass(s.i.DeviceClass)
}

// DisabledByDefault implements entity.Sensor.
func (s *SensorComponent) DisabledByDefault() bool {
	return s.i.DisabledByDefault
}

// EntityCategory implements entity.Sensor.
func (s *SensorComponent) EntityCategory() entity.Category {
	return s.i.EntityCategory
}

// ForceUpdate implements entity.Sensor.
func (s *SensorComponent) ForceUpdate() bool {
	return s.i.ForceUpdate
}

// HashID implements entity.Sensor.
func (s *SensorComponent) HashID() uint32 {
	return s.i.Key
}

// ID implements entity.Sensor.
func (s *SensorComponent) ID() string {
	return s.i.ObjectId
}

// Icon implements entity.Sensor.
func (s *SensorComponent) Icon() string {
	return s.i.Icon
}

// LastResetType implements entity.Sensor.
func (s *SensorComponent) LastResetType() entity.SensorLastResetType {
	return s.i.LastResetType
}

// Name implements entity.Sensor.
func (s *SensorComponent) Name() string {
	return s.i.Name
}

// Setup implements entity.Sensor.
func (s *SensorComponent) Setup() {}

func (s *SensorComponent) UniqueID() string {
	return s.i.UniqueId
}

// StateClass implements entity.Sensor.
func (s *SensorComponent) StateClass() entity.SensorStateClass {
	return s.i.StateClass
}

// UnitOfMeasurement implements entity.Sensor.
func (s *SensorComponent) UnitOfMeasurement() string {
	return s.i.UnitOfMeasurement
}

var _ (entity.Sensor) = (*SensorComponent)(nil)

type SwitchComponent struct {
	ComponentBase
	state[entity.SwitchState]

	i info.Switch
}

// DeviceClass implements entity.Switch.
func (s *SwitchComponent) DeviceClass() entity.SwitchDeviceClass {
	return entity.SwitchDeviceClass(s.i.DeviceClass)
}

// DisabledByDefault implements entity.Switch.
func (s *SwitchComponent) DisabledByDefault() bool {
	return s.i.DisabledByDefault
}

// EntityCategory implements entity.Switch.
func (s *SwitchComponent) EntityCategory() entity.Category {
	return s.i.EntityCategory
}

// HashID implements entity.Switch.
func (s *SwitchComponent) HashID() uint32 {
	return s.i.Key
}

// ID implements entity.Switch.
func (s *SwitchComponent) ID() string {
	return s.i.ObjectId
}

// Icon implements entity.Switch.
func (s *SwitchComponent) Icon() string {
	return s.i.Icon
}

// Name implements entity.Switch.
func (s *SwitchComponent) Name() string {
	return s.i.Name
}

// Setup implements entity.Switch.
func (s *SwitchComponent) Setup() {}

func (s *SwitchComponent) UniqueID() string {
	return s.i.UniqueId
}

var _ (entity.Switch) = (*SwitchComponent)(nil)

type TextSensorComponent struct {
	ComponentBase
	state[entity.TextSensorState]

	i info.TextSensor
}

// DeviceClass implements entity.TextSensor.
func (t *TextSensorComponent) DeviceClass() entity.TextSensorDeviceClass {
	return entity.TextSensorDeviceClass(t.i.DeviceClass)
}

// DisabledByDefault implements entity.TextSensor.
func (t *TextSensorComponent) DisabledByDefault() bool {
	return t.i.DisabledByDefault
}

// EntityCategory implements entity.TextSensor.
func (t *TextSensorComponent) EntityCategory() entity.Category {
	return t.i.EntityCategory
}

// HashID implements entity.TextSensor.
func (t *TextSensorComponent) HashID() uint32 {
	return t.i.Key
}

// ID implements entity.TextSensor.
func (t *TextSensorComponent) ID() string {
	return t.i.ObjectId
}

// Icon implements entity.TextSensor.
func (t *TextSensorComponent) Icon() string {
	return t.i.Icon
}

// Name implements entity.TextSensor.
func (t *TextSensorComponent) Name() string {
	return t.i.Name
}

// Setup implements entity.TextSensor.
func (t *TextSensorComponent) Setup() {}

func (t *TextSensorComponent) UniqueID() string {
	return t.i.UniqueId
}

var _ (entity.TextSensor) = (*TextSensorComponent)(nil)

type ServiceComponent struct {
	ComponentBase

	i info.Services
}

// DisabledByDefault implements entity.Service.
func (s *ServiceComponent) DisabledByDefault() bool {
	return false
}

// EntityCategory implements entity.Service.
func (s *ServiceComponent) EntityCategory() entity.Category {
	return entity.CategoryNone
}

// HashID implements entity.Service.
func (s *ServiceComponent) HashID() uint32 {
	return s.i.Key
}

// ID implements entity.Service.
func (s *ServiceComponent) ID() string {
	return fmt.Sprintf("service %x", s.i.Key)
}

// Name implements entity.Service.
func (s *ServiceComponent) Name() string {
	return s.i.Name
}

// Setup implements entity.Service.
func (s *ServiceComponent) Setup() {}

func (u *ServiceComponent) Close() error {
	return nil
}

var _ (entity.Service) = (*ServiceComponent)(nil)

type CameraComponent struct {
	ComponentBase

	i info.Camera
}

// DisabledByDefault implements entity.Camera.
func (c *CameraComponent) DisabledByDefault() bool {
	return c.i.DisabledByDefault
}

// EntityCategory implements entity.Camera.
func (c *CameraComponent) EntityCategory() entity.Category {
	return c.i.EntityCategory
}

// HashID implements entity.Camera.
func (c *CameraComponent) HashID() uint32 {
	return c.i.Key
}

// ID implements entity.Camera.
func (c *CameraComponent) ID() string {
	return c.i.ObjectId
}

// Icon implements entity.Camera.
func (c *CameraComponent) Icon() string {
	return c.i.Icon
}

// Name implements entity.Camera.
func (c *CameraComponent) Name() string {
	return c.i.Name
}

// Setup implements entity.Camera.
func (c *CameraComponent) Setup() {}

func (c *CameraComponent) UniqueID() string {
	return c.i.UniqueId
}

func (u *CameraComponent) Close() error {
	return nil
}

var _ (entity.Camera) = (*CameraComponent)(nil)

type ClimateComponent struct {
	ComponentBase
	state[entity.ClimateState]

	i info.Climate
}

// DisabledByDefault implements entity.Climate.
func (c *ClimateComponent) DisabledByDefault() bool {
	return c.i.DisabledByDefault
}

// EntityCategory implements entity.Climate.
func (c *ClimateComponent) EntityCategory() entity.Category {
	return c.i.EntityCategory
}

// HashID implements entity.Climate.
func (c *ClimateComponent) HashID() uint32 {
	return c.i.Key
}

// ID implements entity.Climate.
func (c *ClimateComponent) ID() string {
	return c.i.ObjectId
}

// Icon implements entity.Climate.
func (c *ClimateComponent) Icon() string {
	return c.i.Icon
}

// Name implements entity.Climate.
func (c *ClimateComponent) Name() string {
	return c.i.Name
}

// Setup implements entity.Climate.
func (c *ClimateComponent) Setup() {}

func (c *ClimateComponent) UniqueID() string {
	return c.i.UniqueId
}

var _ (entity.Climate) = (*ClimateComponent)(nil)

type NumberComponent struct {
	ComponentBase
	state[entity.NumberState]

	i info.Number
}

// DeviceClass implements entity.Number.
func (n *NumberComponent) DeviceClass() entity.NumberDeviceClass {
	return entity.NumberDeviceClass(n.i.DeviceClass)
}

// DisabledByDefault implements entity.Number.
func (n *NumberComponent) DisabledByDefault() bool {
	return n.i.DisabledByDefault
}

// EntityCategory implements entity.Number.
func (n *NumberComponent) EntityCategory() entity.Category {
	return n.i.EntityCategory
}

// HashID implements entity.Number.
func (n *NumberComponent) HashID() uint32 {
	return n.i.Key
}

// ID implements entity.Number.
func (n *NumberComponent) ID() string {
	return n.i.ObjectId
}

// Icon implements entity.Number.
func (n *NumberComponent) Icon() string {
	return n.i.Icon
}

// Name implements entity.Number.
func (n *NumberComponent) Name() string {
	return n.i.Name
}

// NumberMode implements entity.Number.
func (n *NumberComponent) NumberMode() entity.NumberMode {
	return n.i.Mode
}

// Setup implements entity.Number.
func (n *NumberComponent) Setup() {}

func (n *NumberComponent) UniqueID() string {
	return n.i.UniqueId
}

// UnitOfMeasurement implements entity.Number.
func (n *NumberComponent) UnitOfMeasurement() string {
	return n.i.UnitOfMeasurement
}

var _ (entity.Number) = (*NumberComponent)(nil)

type SelectComponent struct {
	ComponentBase
	state[entity.SelectState]

	i info.Select
}

// Values implements entity.Select.
func (s *SelectComponent) Values() []string {
	return s.i.Options
}

// DisabledByDefault implements entity.Select.
func (s *SelectComponent) DisabledByDefault() bool {
	return s.i.DisabledByDefault
}

// EntityCategory implements entity.Select.
func (s *SelectComponent) EntityCategory() entity.Category {
	return s.i.EntityCategory
}

// HashID implements entity.Select.
func (s *SelectComponent) HashID() uint32 {
	return s.i.Key
}

// ID implements entity.Select.
func (s *SelectComponent) ID() string {
	return s.i.ObjectId
}

// Icon implements entity.Select.
func (s *SelectComponent) Icon() string {
	return s.i.Icon
}

// Name implements entity.Select.
func (s *SelectComponent) Name() string {
	return s.i.Name
}

// Setup implements entity.Select.
func (s *SelectComponent) Setup() {}

func (s *SelectComponent) UniqueID() string {
	return s.i.UniqueId
}

func (s *SelectComponent) Command(value string) error {
	client := s.c.Value()
	if client == nil {
		return ErrClientGone
	}
	return client.sendMessages(&ehp.SelectCommandRequest{
		Key:   s.i.Key,
		State: value,
	})
}

var _ (entity.Select) = (*SelectComponent)(nil)

type SirenComponent struct {
	ComponentBase
	state[entity.SirenState]

	i info.Siren
}

// Setup implements entity.Siren.
func (s *SirenComponent) Setup() {
}

// DisabledByDefault implements entity.Siren.
func (s *SirenComponent) DisabledByDefault() bool {
	return s.i.DisabledByDefault
}

// EntityCategory implements entity.Siren.
func (s *SirenComponent) EntityCategory() entity.Category {
	return s.i.EntityCategory
}

// HashID implements entity.Siren.
func (s *SirenComponent) HashID() uint32 {
	return s.i.Key
}

// ID implements entity.Siren.
func (s *SirenComponent) ID() string {
	return s.i.ObjectId
}

// Icon implements entity.Siren.
func (s *SirenComponent) Icon() string {
	return s.i.Icon
}

// Name implements entity.Siren.
func (s *SirenComponent) Name() string {
	return s.i.Name
}

// SupportsDuration implements entity.Siren.
func (s *SirenComponent) SupportsDuration() bool {
	return s.i.SupportsDuration
}

// SupportsVolume implements entity.Siren.
func (s *SirenComponent) SupportsVolume() bool {
	return s.i.SupportsVolume
}

// Tones implements entity.Siren.
func (s *SirenComponent) Tones() []string {
	return s.i.Tones
}

var _ (entity.Siren) = (*SirenComponent)(nil)

type LockComponent struct {
	ComponentBase
	state[entity.LockState]

	i info.Lock
}

// DisabledByDefault implements entity.Lock.
func (l *LockComponent) DisabledByDefault() bool {
	return l.i.DisabledByDefault
}

// EntityCategory implements entity.Lock.
func (l *LockComponent) EntityCategory() entity.Category {
	return l.i.EntityCategory
}

// HashID implements entity.Lock.
func (l *LockComponent) HashID() uint32 {
	return l.i.Key
}

// ID implements entity.Lock.
func (l *LockComponent) ID() string {
	return l.i.ObjectId
}

// Icon implements entity.Lock.
func (l *LockComponent) Icon() string {
	return l.i.Icon
}

// Name implements entity.Lock.
func (l *LockComponent) Name() string {
	return l.i.Name
}

// Setup implements entity.Lock.
func (l *LockComponent) Setup() {}

func (l *LockComponent) UniqueID() string {
	return l.i.UniqueId
}

var _ (entity.Lock) = (*LockComponent)(nil)

type ButtonComponent struct {
	ComponentBase

	i info.Button
}

// DeviceClass implements entity.Button.
func (b *ButtonComponent) DeviceClass() entity.ButtonDeviceClass {
	return entity.ButtonDeviceClass(b.i.DeviceClass)
}

// DisabledByDefault implements entity.Button.
func (b *ButtonComponent) DisabledByDefault() bool {
	return b.i.DisabledByDefault
}

// EntityCategory implements entity.Button.
func (b *ButtonComponent) EntityCategory() entity.Category {
	return b.i.EntityCategory
}

// HashID implements entity.Button.
func (b *ButtonComponent) HashID() uint32 {
	return b.i.Key
}

// ID implements entity.Button.
func (b *ButtonComponent) ID() string {
	return b.i.ObjectId
}

// Icon implements entity.Button.
func (b *ButtonComponent) Icon() string {
	return b.i.Icon
}

// Name implements entity.Button.
func (b *ButtonComponent) Name() string {
	return b.i.Name
}

// Press implements entity.Button.
func (b *ButtonComponent) Press(ctx context.Context) error {
	client := b.c.Value()
	if client == nil {
		return ErrClientGone
	}
	return client.sendMessages(&ehp.ButtonCommandRequest{
		Key: b.i.Key,
	})
}

// Setup implements entity.Button.
func (b *ButtonComponent) Setup() {}

func (b *ButtonComponent) UniqueID() string {
	return b.i.UniqueId
}

func (u *ButtonComponent) Close() error {
	return nil
}

var _ (entity.Button) = (*ButtonComponent)(nil)

type MediaPlayerComponent struct {
	ComponentBase
	state[entity.MediaPlayerState]

	i info.MediaPlayer
}

// DisabledByDefault implements entity.MediaPlayer.
func (m *MediaPlayerComponent) DisabledByDefault() bool {
	return m.i.DisabledByDefault
}

// EntityCategory implements entity.MediaPlayer.
func (m *MediaPlayerComponent) EntityCategory() entity.Category {
	return m.i.EntityCategory
}

// HashID implements entity.MediaPlayer.
func (m *MediaPlayerComponent) HashID() uint32 {
	return m.i.Key
}

// ID implements entity.MediaPlayer.
func (m *MediaPlayerComponent) ID() string {
	return m.i.ObjectId
}

// Icon implements entity.MediaPlayer.
func (m *MediaPlayerComponent) Icon() string {
	return m.i.Icon
}

// Name implements entity.MediaPlayer.
func (m *MediaPlayerComponent) Name() string {
	return m.i.Name
}

// Setup implements entity.MediaPlayer.
func (m *MediaPlayerComponent) Setup() {}

func (m *MediaPlayerComponent) UniqueID() string {
	return m.i.UniqueId
}

var _ (entity.MediaPlayer) = (*MediaPlayerComponent)(nil)

type AlarmControlPanelComponent struct {
	ComponentBase
	state[entity.AlarmControlPanelState]

	i info.AlarmControlPanel
}

// Setup implements entity.AlarmControlPanel.
func (a *AlarmControlPanelComponent) Setup() {
}

// DisabledByDefault implements entity.AlarmControlPanel.
func (a *AlarmControlPanelComponent) DisabledByDefault() bool {
	return a.i.DisabledByDefault
}

// EntityCategory implements entity.AlarmControlPanel.
func (a *AlarmControlPanelComponent) EntityCategory() entity.Category {
	return a.i.EntityCategory
}

// HashID implements entity.AlarmControlPanel.
func (a *AlarmControlPanelComponent) HashID() uint32 {
	return a.i.Key
}

// ID implements entity.AlarmControlPanel.
func (a *AlarmControlPanelComponent) ID() string {
	return a.i.ObjectId
}

// Icon implements entity.AlarmControlPanel.
func (a *AlarmControlPanelComponent) Icon() string {
	return a.i.Icon
}

// Name implements entity.AlarmControlPanel.
func (a *AlarmControlPanelComponent) Name() string {
	return a.i.Name
}

var _ (entity.AlarmControlPanel) = (*AlarmControlPanelComponent)(nil)

type TextComponent struct {
	ComponentBase
	state[entity.TextState]

	i info.Text
}

// DisabledByDefault implements entity.Text.
func (t *TextComponent) DisabledByDefault() bool {
	return t.i.DisabledByDefault
}

// EntityCategory implements entity.Text.
func (t *TextComponent) EntityCategory() entity.Category {
	return t.i.EntityCategory
}

// HashID implements entity.Text.
func (t *TextComponent) HashID() uint32 {
	return t.i.Key
}

// ID implements entity.Text.
func (t *TextComponent) ID() string {
	return t.i.ObjectId
}

// Icon implements entity.Text.
func (t *TextComponent) Icon() string {
	return t.i.Icon
}

// Name implements entity.Text.
func (t *TextComponent) Name() string {
	return t.i.Name
}

// Setup implements entity.Text.
func (t *TextComponent) Setup() {}

func (t *TextComponent) UniqueID() string {
	return t.i.UniqueId
}

// TextMode implements entity.Text.
func (t *TextComponent) TextMode() entity.TextMode {
	return t.i.Mode
}

var _ (entity.Text) = (*TextComponent)(nil)

type DateComponent struct {
	ComponentBase
	state[entity.DateState]

	i info.Date
}

// DisabledByDefault implements entity.Date.
func (d *DateComponent) DisabledByDefault() bool {
	return d.i.DisabledByDefault
}

// EntityCategory implements entity.Date.
func (d *DateComponent) EntityCategory() entity.Category {
	return d.i.EntityCategory
}

// HashID implements entity.Date.
func (d *DateComponent) HashID() uint32 {
	return d.i.Key
}

// ID implements entity.Date.
func (d *DateComponent) ID() string {
	return d.i.ObjectId
}

// Icon implements entity.Date.
func (d *DateComponent) Icon() string {
	return d.i.Icon
}

// Name implements entity.Date.
func (d *DateComponent) Name() string {
	return d.i.Name
}

// Setup implements entity.Date.
func (d *DateComponent) Setup() {}

func (d *DateComponent) UniqueID() string {
	return d.i.UniqueId
}

var _ (entity.Date) = (*DateComponent)(nil)

type TimeComponent struct {
	ComponentBase
	state[entity.TimeState]

	i info.Time
}

// DisabledByDefault implements entity.Time.
func (t *TimeComponent) DisabledByDefault() bool {
	return t.i.DisabledByDefault
}

// EntityCategory implements entity.Time.
func (t *TimeComponent) EntityCategory() entity.Category {
	return t.i.EntityCategory
}

// HashID implements entity.Time.
func (t *TimeComponent) HashID() uint32 {
	return t.i.Key
}

// ID implements entity.Time.
func (t *TimeComponent) ID() string {
	return t.i.ObjectId
}

// Icon implements entity.Time.
func (t *TimeComponent) Icon() string {
	return t.i.Icon
}

// Name implements entity.Time.
func (t *TimeComponent) Name() string {
	return t.i.Name
}

// Setup implements entity.Time.
func (t *TimeComponent) Setup() {}

func (t *TimeComponent) UniqueID() string {
	return t.i.UniqueId
}

var _ (entity.Time) = (*TimeComponent)(nil)

type EventComponent struct {
	ComponentBase

	i info.Event
}

// DeviceClass implements entity.Event.
func (e *EventComponent) DeviceClass() entity.EventDeviceClass {
	return entity.EventDeviceClass(e.i.DeviceClass)
}

// DisabledByDefault implements entity.Event.
func (e *EventComponent) DisabledByDefault() bool {
	return e.i.DisabledByDefault
}

// EntityCategory implements entity.Event.
func (e *EventComponent) EntityCategory() entity.Category {
	return e.i.EntityCategory
}

// HashID implements entity.Event.
func (e *EventComponent) HashID() uint32 {
	return e.i.Key
}

// ID implements entity.Event.
func (e *EventComponent) ID() string {
	return e.i.ObjectId
}

// Icon implements entity.Event.
func (e *EventComponent) Icon() string {
	return e.i.Icon
}

// Name implements entity.Event.
func (e *EventComponent) Name() string {
	return e.i.Name
}

// Setup implements entity.Event.
func (e *EventComponent) Setup() {}

func (e *EventComponent) UniqueID() string {
	return e.i.UniqueId
}

func (u *EventComponent) Close() error {
	return nil
}

var _ (entity.Event) = (*EventComponent)(nil)

type ValveComponent struct {
	ComponentBase
	state[entity.ValveState]

	i info.Valve
}

// DeviceClass implements entity.Valve.
func (v *ValveComponent) DeviceClass() entity.ValveDeviceClass {
	return entity.ValveDeviceClass(v.i.DeviceClass)
}

// DisabledByDefault implements entity.Valve.
func (v *ValveComponent) DisabledByDefault() bool {
	return v.i.DisabledByDefault
}

// EntityCategory implements entity.Valve.
func (v *ValveComponent) EntityCategory() entity.Category {
	return v.i.EntityCategory
}

// HashID implements entity.Valve.
func (v *ValveComponent) HashID() uint32 {
	return v.i.Key
}

// ID implements entity.Valve.
func (v *ValveComponent) ID() string {
	return v.i.ObjectId
}

// Icon implements entity.Valve.
func (v *ValveComponent) Icon() string {
	return v.i.Icon
}

// Name implements entity.Valve.
func (v *ValveComponent) Name() string {
	return v.i.Name
}

// Setup implements entity.Valve.
func (v *ValveComponent) Setup() {}

func (v *ValveComponent) UniqueID() string {
	return v.i.UniqueId
}

var _ (entity.Valve) = (*ValveComponent)(nil)

type DatetimeComponent struct {
	ComponentBase
	state[entity.DatetimeState]

	i info.DateTime
}

// DisabledByDefault implements entity.Datetime.
func (d *DatetimeComponent) DisabledByDefault() bool {
	return d.i.DisabledByDefault
}

// EntityCategory implements entity.Datetime.
func (d *DatetimeComponent) EntityCategory() entity.Category {
	return d.i.EntityCategory
}

// HashID implements entity.Datetime.
func (d *DatetimeComponent) HashID() uint32 {
	return d.i.Key
}

// ID implements entity.Datetime.
func (d *DatetimeComponent) ID() string {
	return d.i.ObjectId
}

// Icon implements entity.Datetime.
func (d *DatetimeComponent) Icon() string {
	return d.i.Icon
}

// Name implements entity.Datetime.
func (d *DatetimeComponent) Name() string {
	return d.i.Name
}

// Setup implements entity.Datetime.
func (d *DatetimeComponent) Setup() {}

func (d *DatetimeComponent) UniqueID() string {
	return d.i.UniqueId
}

var _ (entity.Datetime) = (*DatetimeComponent)(nil)

type UpdateComponent struct {
	ComponentBase
	state[entity.UpdateState]

	i info.Update
}

// DeviceClass implements entity.Update.
func (u *UpdateComponent) DeviceClass() entity.UpdateDeviceClass {
	return entity.UpdateDeviceClass(u.i.DeviceClass)
}

// DisabledByDefault implements entity.Update.
func (u *UpdateComponent) DisabledByDefault() bool {
	return u.i.DisabledByDefault
}

// EntityCategory implements entity.Update.
func (u *UpdateComponent) EntityCategory() entity.Category {
	return u.i.EntityCategory
}

// HashID implements entity.Update.
func (u *UpdateComponent) HashID() uint32 {
	return u.i.Key
}

// ID implements entity.Update.
func (u *UpdateComponent) ID() string {
	return u.i.ObjectId
}

// Icon implements entity.Update.
func (u *UpdateComponent) Icon() string {
	return u.i.Icon
}

// Name implements entity.Update.
func (u *UpdateComponent) Name() string {
	return u.i.Name
}

// Setup implements entity.Update.
func (u *UpdateComponent) Setup() {}

func (u *UpdateComponent) UniqueID() string {
	return u.i.UniqueId
}

var _ (entity.Update) = (*UpdateComponent)(nil)
