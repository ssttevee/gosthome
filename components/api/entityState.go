package api

import (
	"fmt"
	"log/slog"

	"github.com/gosthome/gosthome/components/api/common"
	ehp "github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/core/entity"
)

func entityState(ent entity.Entity) ehp.EsphomeMessageTyper {
	switch typed := ent.(type) {
	case entity.BinarySensor:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Cover:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Fan:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Light:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Sensor:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Switch:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.TextSensor:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Climate:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Number:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Date:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Time:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Datetime:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Text:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Select:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Lock:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Siren:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Valve:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.MediaPlayer:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.AlarmControlPanel:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	case entity.Update:
		state := typed.State()
		return stateResponse(typed.HashID(), &state)
	default:
		return nil
	}
}

func stateResponse(key uint32, astate any) ehp.EsphomeMessageTyper {
	switch state := astate.(type) {
	case *entity.BinarySensorState:
		return &ehp.BinarySensorStateResponse{
			Key:          key,
			State:        state.State,
			MissingState: state.Missing,
		}
	case *entity.CoverState:
		return &ehp.CoverStateResponse{
			Key:         key,
			LegacyState: common.Enum[ehp.LegacyCoverState](state.LegacyState),
		}
	case *entity.FanState:
		return &ehp.FanStateResponse{
			Key:         key,
			State:       state.State,
			Oscillating: state.Oscillating,
			Speed:       common.Enum[ehp.FanSpeed](state.Speed),
			Direction:   common.Enum[ehp.FanDirection](state.Direction),
			SpeedLevel:  state.SpeedLevel,
			PresetMode:  state.PresetMode,
		}
	case *entity.LightState:
		return &ehp.LightStateResponse{
			Key:              key,
			State:            state.State,
			Brightness:       state.Brightness,
			ColorMode:        common.Enum[int32](state.ColorMode),
			ColorBrightness:  state.ColorBrightness,
			Red:              state.Red,
			Green:            state.Green,
			Blue:             state.Blue,
			White:            state.White,
			ColorTemperature: state.ColorTemperature,
			ColdWhite:        state.ColdWhite,
			WarmWhite:        state.WarmWhite,
			Effect:           state.Effect,
		}
	case *entity.SensorState:
		return &ehp.SensorStateResponse{
			Key:          key,
			State:        state.State,
			MissingState: state.MissingState,
		}
	case *entity.SwitchState:
		return &ehp.SwitchStateResponse{
			Key:   key,
			State: state.State,
		}
	case *entity.TextSensorState:
		return &ehp.TextSensorStateResponse{
			Key:          key,
			State:        state.State,
			MissingState: state.MissingState,
		}
	case *entity.ClimateState:
		return &ehp.ClimateStateResponse{
			Key:                   key,
			Mode:                  common.Enum[ehp.ClimateMode](state.Mode),
			CurrentTemperature:    state.CurrentTemperature,
			TargetTemperature:     state.TargetTemperature,
			TargetTemperatureLow:  state.TargetTemperatureLow,
			TargetTemperatureHigh: state.TargetTemperatureHigh,
			LegacyAway:            state.Preset == entity.ClimatePresetAway,
			Action:                common.Enum[ehp.ClimateAction](state.Action),
			FanMode:               common.Enum[ehp.ClimateFanMode](state.FanMode),
			SwingMode:             common.Enum[ehp.ClimateSwingMode](state.SwingMode),
			CustomFanMode:         state.CustomFanMode,
			Preset:                common.Enum[ehp.ClimatePreset](state.Preset),
			CustomPreset:          state.CustomPreset,
			CurrentHumidity:       state.CurrentHumidity,
			TargetHumidity:        state.TargetHumidity,
		}
	case *entity.NumberState:
		return &ehp.NumberStateResponse{
			Key:          key,
			State:        state.State,
			MissingState: state.MissingState,
		}
	case *entity.DateState:
		return &ehp.DateStateResponse{
			Key:          key,
			MissingState: state.MissingState,
			Year:         state.Year,
			Month:        state.Month,
			Day:          state.Day,
		}
	case *entity.TimeState:
		return &ehp.TimeStateResponse{
			Key:          key,
			MissingState: state.MissingState,
			Hour:         state.Hour,
			Minute:       state.Minute,
			Second:       state.Second,
		}
	case *entity.DatetimeState:
		return &ehp.DateTimeStateResponse{
			Key:          key,
			MissingState: state.MissingState,
			EpochSeconds: state.EpochSeconds,
		}
	case *entity.TextState:
		return &ehp.TextStateResponse{
			Key:          key,
			State:        state.State,
			MissingState: state.MissingState,
		}
	case *entity.SelectState:
		return &ehp.SelectStateResponse{
			Key:          key,
			State:        state.State,
			MissingState: state.MissingState,
		}
	case *entity.SirenState:
		return &ehp.SirenStateResponse{
			Key:   key,
			State: bool(*state),
		}
	case *entity.LockState:
		return &ehp.LockStateResponse{
			Key:   key,
			State: common.Enum[ehp.LockState](*state),
		}
	case *entity.ValveState:
		return &ehp.ValveStateResponse{
			Key:              key,
			Position:         state.Position,
			CurrentOperation: common.Enum[ehp.ValveOperation](state.CurrentOperation),
		}
	case *entity.MediaPlayerState:
		return &ehp.MediaPlayerStateResponse{
			Key:    key,
			State:  common.Enum[ehp.MediaPlayerState](state.State),
			Volume: state.Volume,
			Muted:  state.Muted,
		}
	case *entity.AlarmControlPanelState:
		return &ehp.AlarmControlPanelStateResponse{
			Key:   key,
			State: common.Enum[ehp.AlarmControlPanelState](*state),
		}
	case *entity.UpdateState:
		return &ehp.UpdateStateResponse{
			Key:            key,
			MissingState:   state.MissingState,
			InProgress:     state.InProgress,
			HasProgress:    state.HasProgress,
			Progress:       state.Progress,
			CurrentVersion: state.CurrentVersion,
			LatestVersion:  state.LatestVersion,
			Title:          state.Title,
			ReleaseSummary: state.ReleaseSummary,
			ReleaseUrl:     state.ReleaseUrl,
		}
	default:
		slog.Error("Unknown state in entity state", "state", fmt.Sprintf("%#v", state))
		return nil
	}
}
