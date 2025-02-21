package client

import (
	"context"
	"errors"
	"iter"
	"log/slog"
	"time"
	"weak"

	"github.com/gosthome/gosthome/components/api/common"
	ehp "github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/entity/info"
)

var ErrAlreadyInProgress = errors.New("already in progress")

func (c *Client) SubscribeStates() error {
	err := c.sendMessages(&ehp.SubscribeStatesRequest{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ListEntities(timeout time.Duration) error {
	var ctx context.Context
	var canc context.CancelFunc
	var listEndChan chan struct{}
	err := c.listEntitiesState.DoErr(func(state *chan<- struct{}) error {
		if *state != nil {
			return ErrAlreadyInProgress
		}
		ctx, canc = context.WithTimeout(c.ctx, timeout)
		listEndChan = make(chan struct{})
		*state = listEndChan
		err := c.sendMessages(&ehp.ListEntitiesRequest{})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	defer canc()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-listEndChan:
		return nil
	}
}

func (c *Client) AllEntities() iter.Seq2[entity.DomainType, entity.Entity] {
	return entity.IterateRegistry(&c.reg)
}

func (c *Client) componentRegistration(err error) error {
	if err != nil {
		slog.Warn("This host does not generate unique ids!", "err", err)
	}
	return nil
}

func (c *Client) listEntitiesResponse(msg ehp.EsphomeMessageTyper) error {
	var ok bool
	c.listEntitiesState.Do(func(state *chan<- struct{}) { ok = *state == nil })
	if ok {
		slog.Error("Unexpected list entities message", "msg", msg)
		return nil
	}
	switch list := msg.(type) {
	case *ehp.ListEntitiesBinarySensorResponse:
		err := c.reg.RegisterBinarySensor(&BinarySensorComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.BinarySensor{
				ObjectId:             list.ObjectId,
				Key:                  list.Key,
				Name:                 list.Name,
				UniqueId:             list.UniqueId,
				DeviceClass:          list.DeviceClass,
				IsStatusBinarySensor: list.IsStatusBinarySensor,
				DisabledByDefault:    list.DisabledByDefault,
				Icon:                 list.Icon,
				EntityCategory:       common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesCoverResponse:
		err := c.reg.RegisterCover(&CoverComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Cover{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				AssumedState:      list.AssumedState,
				SupportsPosition:  list.SupportsPosition,
				SupportsTilt:      list.SupportsTilt,
				DeviceClass:       list.DeviceClass,
				DisabledByDefault: list.DisabledByDefault,
				Icon:              list.Icon,
				SupportsStop:      list.SupportsStop,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesFanResponse:
		err := c.reg.RegisterFan(&FanComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Fan{
				ObjectId:             list.ObjectId,
				Key:                  list.Key,
				Name:                 list.Name,
				UniqueId:             list.UniqueId,
				SupportsOscillation:  list.SupportsOscillation,
				SupportsSpeed:        list.SupportsSpeed,
				SupportsDirection:    list.SupportsDirection,
				SupportedSpeedLevels: list.SupportedSpeedLevels,
				DisabledByDefault:    list.DisabledByDefault,
				Icon:                 list.Icon,
				EntityCategory:       common.Enum[entity.Category](list.EntityCategory),
				SupportedPresetModes: list.SupportedPresetModes,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesLightResponse:
		err := c.reg.RegisterLight(&LightComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Light{
				ObjectId:            list.ObjectId,
				Key:                 list.Key,
				Name:                list.Name,
				UniqueId:            list.UniqueId,
				SupportedColorModes: common.Enums[entity.ColorMode](list.SupportedColorModes),
				MinMireds:           list.MinMireds,
				MaxMireds:           list.MaxMireds,
				Effects:             list.Effects,
				DisabledByDefault:   list.DisabledByDefault,
				Icon:                list.Icon,
				EntityCategory:      common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesSensorResponse:
		err := c.reg.RegisterSensor(&SensorComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Sensor{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				UnitOfMeasurement: list.UnitOfMeasurement,
				AccuracyDecimals:  list.AccuracyDecimals,
				ForceUpdate:       list.ForceUpdate,
				DeviceClass:       list.DeviceClass,
				StateClass:        common.Enum[entity.SensorStateClass](list.StateClass),
				LastResetType:     common.Enum[entity.SensorLastResetType](list.LastResetType),
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesSwitchResponse:
		err := c.reg.RegisterSwitch(&SwitchComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Switch{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				AssumedState:      list.AssumedState,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				DeviceClass:       list.DeviceClass,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesTextSensorResponse:
		err := c.reg.RegisterTextSensor(&TextSensorComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.TextSensor{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				DeviceClass:       list.DeviceClass,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesServicesResponse:
		err := c.reg.RegisterService(&ServiceComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Services{
				Name: list.Name,
				Key:  list.Key,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesCameraResponse:
		err := c.reg.RegisterCamera(&CameraComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Camera{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				DisabledByDefault: list.DisabledByDefault,
				Icon:              list.Icon,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesClimateResponse:
		err := c.reg.RegisterClimate(&ClimateComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Climate{
				ObjectId:                          list.ObjectId,
				Key:                               list.Key,
				Name:                              list.Name,
				UniqueId:                          list.UniqueId,
				SupportsCurrentTemperature:        list.SupportsCurrentTemperature,
				SupportsTwoPointTargetTemperature: list.SupportsTwoPointTargetTemperature,
				SupportedModes:                    common.Enums[entity.ClimateMode](list.SupportedModes),
				VisualMinTemperature:              list.VisualMinTemperature,
				VisualMaxTemperature:              list.VisualMaxTemperature,
				VisualTargetTemperatureStep:       list.VisualTargetTemperatureStep,
				// for older peer versions - in new system this
				// is if CLIMATE_PRESET_AWAY exists is supported_presets
				LegacySupportsAway:           list.LegacySupportsAway,
				SupportsAction:               list.SupportsAction,
				SupportedFanModes:            common.Enums[entity.ClimateFanMode](list.SupportedFanModes),
				SupportedSwingModes:          common.Enums[entity.ClimateSwingMode](list.SupportedSwingModes),
				SupportedCustomFanModes:      list.SupportedCustomFanModes,
				SupportedPresets:             common.Enums[entity.ClimatePreset](list.SupportedPresets),
				SupportedCustomPresets:       list.SupportedCustomPresets,
				DisabledByDefault:            list.DisabledByDefault,
				Icon:                         list.Icon,
				EntityCategory:               common.Enum[entity.Category](list.EntityCategory),
				VisualCurrentTemperatureStep: list.VisualCurrentTemperatureStep,
				SupportsCurrentHumidity:      list.SupportsCurrentHumidity,
				SupportsTargetHumidity:       list.SupportsTargetHumidity,
				VisualMinHumidity:            list.VisualMinHumidity,
				VisualMaxHumidity:            list.VisualMaxHumidity,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesNumberResponse:
		err := c.reg.RegisterNumber(&NumberComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Number{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				MinValue:          list.MinValue,
				MaxValue:          list.MaxValue,
				Step:              list.Step,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				UnitOfMeasurement: list.UnitOfMeasurement,
				Mode:              common.Enum[entity.NumberMode](list.Mode),
				DeviceClass:       list.DeviceClass,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesSelectResponse:
		err := c.reg.RegisterSelect(&SelectComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Select{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				Options:           list.Options,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesSirenResponse:
		c.reg.RegisterSiren(&SirenComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Siren{},
		})
		return nil
	case *ehp.ListEntitiesLockResponse:
		err := c.reg.RegisterLock(&LockComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Lock{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				AssumedState:      list.AssumedState,
				SupportsOpen:      list.SupportsOpen,
				RequiresCode:      list.RequiresCode,
				CodeFormat:        list.CodeFormat,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesButtonResponse:
		err := c.reg.RegisterButton(&ButtonComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Button{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				DeviceClass:       list.DeviceClass,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesMediaPlayerResponse:
		err := c.reg.RegisterMediaPlayer(&MediaPlayerComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.MediaPlayer{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				SupportsPause:     list.SupportsPause,
				SupportedFormats: func() []*entity.MediaPlayerSupportedFormat {
					ret := make([]*entity.MediaPlayerSupportedFormat, len(list.SupportedFormats))
					for i, sf := range list.SupportedFormats {
						ret[i] = &entity.MediaPlayerSupportedFormat{
							Format:      sf.Format,
							SampleRate:  sf.SampleRate,
							NumChannels: sf.NumChannels,
							Purpose:     common.Enum[entity.MediaPlayerFormatPurpose](sf.Purpose),
							SampleBytes: sf.SampleBytes,
						}
					}
					return ret
				}(),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesAlarmControlPanelResponse:
		err := c.reg.RegisterAlarmControlPanel(&AlarmControlPanelComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.AlarmControlPanel{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				SupportedFeatures: list.SupportedFeatures,
				RequiresCode:      list.RequiresCode,
				RequiresCodeToArm: list.RequiresCodeToArm,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesTextResponse:
		err := c.reg.RegisterText(&TextComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Text{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				MinLength:         list.MinLength,
				MaxLength:         list.MaxLength,
				Pattern:           list.Pattern,
				Mode:              common.Enum[entity.TextMode](list.Mode),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesDateResponse:
		err := c.reg.RegisterDate(&DateComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Date{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesTimeResponse:
		err := c.reg.RegisterTime(&TimeComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Time{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesEventResponse:
		err := c.reg.RegisterEvent(&EventComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Event{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				DeviceClass:       list.DeviceClass,
				EventTypes:        list.EventTypes,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesValveResponse:
		err := c.reg.RegisterValve(&ValveComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Valve{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				DeviceClass:       list.DeviceClass,
				AssumedState:      list.AssumedState,
				SupportsPosition:  list.SupportsPosition,
				SupportsStop:      list.SupportsStop,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesDateTimeResponse:
		err := c.reg.RegisterDatetime(&DatetimeComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.DateTime{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesUpdateResponse:
		err := c.reg.RegisterUpdate(&UpdateComponent{
			ComponentBase: ComponentBase{
				c: weak.Make(c),
			},
			i: info.Update{
				ObjectId:          list.ObjectId,
				Key:               list.Key,
				Name:              list.Name,
				UniqueId:          list.UniqueId,
				Icon:              list.Icon,
				DisabledByDefault: list.DisabledByDefault,
				EntityCategory:    common.Enum[entity.Category](list.EntityCategory),
				DeviceClass:       list.DeviceClass,
			},
		})
		return c.componentRegistration(err)
	case *ehp.ListEntitiesDoneResponse:
		c.listEntitiesState.Do(func(state *chan<- struct{}) {
			close(*state)
			*state = nil
		})
		return nil
	default:
		slog.Error("unexpected message ", "msg", msg)
		return nil
	}
}

func (c *Client) stateChangeResponse(msg ehp.EsphomeMessageTyper) error {
	switch state := msg.(type) {
	case *ehp.BinarySensorStateResponse:
		comp, ok := c.BinarySensorByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this BinarySensor, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.BinarySensorState{
			State:   state.State,
			Missing: state.MissingState,
		})
	case *ehp.CoverStateResponse:
		comp, ok := c.CoverByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Cover, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.CoverState{
			LegacyState: common.Enum[entity.LegacyCoverState](state.LegacyState),
			Position:    state.Position,
			Tilt:        state.Tilt,
		})
	case *ehp.FanStateResponse:
		comp, ok := c.FanByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Fan, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.FanState{
			State:       state.State,
			Oscillating: state.Oscillating,
			Speed:       common.Enum[entity.FanSpeed](state.Speed),
			Direction:   common.Enum[entity.FanDirection](state.Direction),
			SpeedLevel:  state.SpeedLevel,
			PresetMode:  state.PresetMode,
		})
	case *ehp.LightStateResponse:
		comp, ok := c.LightByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Light, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.LightState{
			State:            state.State,
			Brightness:       state.Brightness,
			ColorMode:        common.Enum[entity.ColorMode](state.ColorMode),
			ColorBrightness:  state.ColorBrightness,
			Red:              state.Red,
			Green:            state.Green,
			Blue:             state.Blue,
			White:            state.White,
			ColorTemperature: state.ColorTemperature,
			ColdWhite:        state.ColdWhite,
			WarmWhite:        state.WarmWhite,
			Effect:           state.Effect,
		})

	case *ehp.SensorStateResponse:
		comp, ok := c.SensorByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Sensor, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.SensorState{
			State:        state.State,
			MissingState: state.MissingState,
		})
	case *ehp.SwitchStateResponse:
		comp, ok := c.SwitchByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Switch, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.SwitchState{
			State: state.State,
		})
	case *ehp.TextSensorStateResponse:
		comp, ok := c.TextSensorByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this TextSensor, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.TextSensorState{
			State:        state.State,
			MissingState: state.MissingState,
		})
	case *ehp.ClimateStateResponse:
		comp, ok := c.ClimateByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Climate, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.ClimateState{
			Mode:                  common.Enum[entity.ClimateMode](state.Mode),
			CurrentTemperature:    state.CurrentTemperature,
			TargetTemperature:     state.TargetTemperature,
			TargetTemperatureLow:  state.TargetTemperatureLow,
			TargetTemperatureHigh: state.TargetTemperatureHigh,
			Action:                common.Enum[entity.ClimateAction](state.Action),
			FanMode:               common.Enum[entity.ClimateFanMode](state.FanMode),
			SwingMode:             common.Enum[entity.ClimateSwingMode](state.SwingMode),
			CustomFanMode:         state.CustomFanMode,
			Preset:                common.Enum[entity.ClimatePreset](state.Preset),
			CustomPreset:          state.CustomPreset,
			CurrentHumidity:       state.CurrentHumidity,
			TargetHumidity:        state.TargetHumidity,
		})
	case *ehp.NumberStateResponse:
		comp, ok := c.NumberByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Number, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.NumberState{
			State:        state.State,
			MissingState: state.MissingState,
		})
	case *ehp.SelectStateResponse:
		comp, ok := c.SelectByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Select, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.SelectState{
			State:        state.State,
			MissingState: state.MissingState,
		})
	case *ehp.SirenStateResponse:
		comp, ok := c.SirenByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Siren, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.SirenState(state.State))
	case *ehp.LockStateResponse:
		comp, ok := c.LockByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Lock, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(common.Enum[entity.LockState](state.State))
	case *ehp.MediaPlayerStateResponse:
		comp, ok := c.MediaPlayerByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this MediaPlayer, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.MediaPlayerState{
			State:  common.Enum[entity.MediaPlayingState](state.State),
			Volume: state.Volume,
			Muted:  state.Muted,
		})
	case *ehp.AlarmControlPanelStateResponse:
		comp, ok := c.AlarmControlPanelByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this AlarmControlPanel, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(common.Enum[entity.AlarmControlPanelState](state.State))
	case *ehp.TextStateResponse:
		comp, ok := c.TextByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Text, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.TextState{
			State:        state.State,
			MissingState: state.MissingState,
		})
	case *ehp.DateStateResponse:
		comp, ok := c.DateByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Date, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.DateState{
			MissingState: state.MissingState,
			Year:         state.Year,
			Month:        state.Month,
			Day:          state.Day,
		})
	case *ehp.TimeStateResponse:
		comp, ok := c.TimeByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Time, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.TimeState{
			MissingState: state.MissingState,
			Hour:         state.Hour,
			Minute:       state.Minute,
			Second:       state.Second,
		})
	case *ehp.ValveStateResponse:
		comp, ok := c.ValveByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Valve, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.ValveState{
			Position:         state.Position,
			CurrentOperation: common.Enum[entity.ValveOperation](state.CurrentOperation),
		})
	case *ehp.DateTimeStateResponse:
		comp, ok := c.DatetimeByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this DateTime, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.DatetimeState{
			MissingState: state.MissingState,
			EpochSeconds: state.EpochSeconds,
		})
	case *ehp.UpdateStateResponse:
		comp, ok := c.UpdateByKey(state.Key)
		if !ok {
			slog.Warn("Client does not know about this Update, did you subscribed to state changes before lising entities?", "key", state.Key)
			return nil
		}
		comp.SetState(entity.UpdateState{
			MissingState:   state.MissingState,
			InProgress:     state.InProgress,
			HasProgress:    state.HasProgress,
			Progress:       state.Progress,
			CurrentVersion: state.CurrentVersion,
			LatestVersion:  state.LatestVersion,
			Title:          state.Title,
			ReleaseSummary: state.ReleaseSummary,
			ReleaseUrl:     state.ReleaseUrl,
		})
	}
	return nil
}
