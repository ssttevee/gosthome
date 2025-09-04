package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/gosthome/gosthome/components/api/common"
	ehp "github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/components/button"
	"github.com/gosthome/gosthome/components/climate"
	"github.com/gosthome/gosthome/components/number"
	"github.com/gosthome/gosthome/components/switchcomp"
	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/entity"
)

var (
	defaultHandlers = safeMessageHandlers{m: map[ehp.MessageType]AnyMessageHandler{}}
)

func dH(mh MessageHandler) byte {
	defaultHandlers.locked(func(m map[ehp.MessageType]AnyMessageHandler) {
		m[mh.Type] = mh.Handler
	})
	return 0
}

func WithAuth(mh MessageHandler) MessageHandler {
	oh := mh.Handler
	mh.Handler = func(ctx context.Context, c *Connection, m ehp.EsphomeMessageTyper) ([]ehp.EsphomeMessageTyper, error) {
		if !c.authenticated {
			return nil, errors.New("unauthenticated access")
		}
		return oh(ctx, c, m)
	}
	return mh
}

var (
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.HelloRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Info("Client connected", "clientApiVersionMajor", msg.ApiVersionMajor, "clientApiVersionMinor", msg.ApiVersionMinor, "clientInfo", msg.ClientInfo)
		c.clientInfo = msg.ClientInfo
		cfg := core.GetNode(ctx).Config
		return []ehp.EsphomeMessageTyper{
			&ehp.HelloResponse{
				ApiVersionMajor: common.ApiVersionMajor,
				ApiVersionMinor: common.ApiVersionMinor,

				Name:       cfg.Gosthome.Name,
				ServerInfo: fmt.Sprintf("gosthome %s based on aioesphomeapi %s", core.Version(), ehp.ESPHOME_VERSION),
			},
		}, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.ConnectRequest) ([]ehp.EsphomeMessageTyper, error) {
		pass := c.server.config.Password
		valid := !pass.Valid() || pass.Check(msg.Password)
		slog.Info("Connect request", "valid", valid)
		if valid {
			c.authenticated = true
		}
		c.canAuthenticate = false
		return []ehp.EsphomeMessageTyper{
			&ehp.ConnectResponse{
				InvalidPassword: !valid,
			},
		}, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.DisconnectRequest) ([]ehp.EsphomeMessageTyper, error) {
		// option (needs_setup_connection) = false;
		// option (needs_authentication) = false;
		return []ehp.EsphomeMessageTyper{&ehp.DisconnectResponse{}}, frameshakers.ErrCloseConnection
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.PingRequest) ([]ehp.EsphomeMessageTyper, error) {

		// option (needs_setup_connection) = false;
		// option (needs_authentication) = false;
		return []ehp.EsphomeMessageTyper{&ehp.PingResponse{}}, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.DeviceInfoRequest) ([]ehp.EsphomeMessageTyper, error) {
		cfg := core.GetNode(ctx).Config
		serverCfg := c.server.config
		return []ehp.EsphomeMessageTyper{&ehp.DeviceInfoResponse{
			UsesPassword:    serverCfg.Password.Valid(),
			Name:            cfg.Gosthome.Name,
			FriendlyName:    cfg.Gosthome.FriendlyName,
			SuggestedArea:   cfg.Gosthome.Area,
			MacAddress:      cfg.Gosthome.MAC.String(),
			EsphomeVersion:  ehp.ESPHOME_VERSION,
			CompilationTime: "2022",
			Manufacturer:    "gosthome",
			Model:           runtime.GOOS + "/" + runtime.GOARCH,
			HasDeepSleep:    false,
			ProjectName:     cfg.Gosthome.Project.Name,
			ProjectVersion:  cfg.Gosthome.Project.Version,
			// WebserverPort:               0,
			LegacyBluetoothProxyVersion: 0,
			BluetoothProxyFeatureFlags:  0,
			LegacyVoiceAssistantVersion: 0,
			VoiceAssistantFeatureFlags:  0,
		}}, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.GetTimeRequest) ([]ehp.EsphomeMessageTyper, error) {
		return []ehp.EsphomeMessageTyper{&ehp.GetTimeResponse{
			EpochSeconds: uint32(time.Now().Unix()),
		}}, nil
	}))
	_ = dH(WithAuth(Handler[ehp.ListEntitiesRequest](func(ctx context.Context, c *Connection, msg *ehp.ListEntitiesRequest) ([]ehp.EsphomeMessageTyper, error) {
		ret := []ehp.EsphomeMessageTyper{}
		node := core.GetNode(ctx)
		for t, ent := range entity.IterateRegistry(node.Registry) {
			if ent.Internal() {
				continue
			}
			switch typed := ent.(type) {
			case entity.BinarySensor:
				ret = append(ret, &ehp.ListEntitiesBinarySensorResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
				})
			case entity.Cover:
				ret = append(ret, &ehp.ListEntitiesCoverResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
				})
			case entity.Fan:
				ret = append(ret, &ehp.ListEntitiesFanResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Light:
				ret = append(ret, &ehp.ListEntitiesLightResponse{
					ObjectId:            typed.ID(),
					Key:                 typed.HashID(),
					DisabledByDefault:   typed.DisabledByDefault(),
					EntityCategory:      common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:                typed.Name(),
					UniqueId:            node.DefaultUniqueId(t, typed),
					Icon:                typed.Icon(),
					SupportedColorModes: common.Enums[int32](typed.SupportedColorModes()),
					Effects:             typed.Effects(),
					MinMireds:           typed.MinMireds(),
					MaxMireds:           typed.MaxMireds(),
				})
			case entity.Sensor:
				ret = append(ret, &ehp.ListEntitiesSensorResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					UnitOfMeasurement: typed.UnitOfMeasurement(),
					AccuracyDecimals:  typed.AccuracyDecimals(), //int32
					ForceUpdate:       typed.ForceUpdate(),      //bool
					DeviceClass:       string(typed.DeviceClass()),
					StateClass:        common.Enum[ehp.SensorStateClass](typed.StateClass()),
					LastResetType:     common.Enum[ehp.SensorLastResetType](typed.LastResetType()),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
				})
			case entity.Switch:
				ret = append(ret, &ehp.ListEntitiesSwitchResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
				})
			case entity.Button:
				ret = append(ret, &ehp.ListEntitiesButtonResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
				})
			case entity.TextSensor:
				ret = append(ret, &ehp.ListEntitiesTextSensorResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
				})
			case entity.Camera:
				ret = append(ret, &ehp.ListEntitiesCameraResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Climate:
				res := ehp.ListEntitiesClimateResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				}

				if c, ok := typed.(interface{ SupportsCurrentTemperature() bool }); ok {
					res.SupportsCurrentTemperature = c.SupportsCurrentTemperature()
				}

				if c, ok := typed.(interface{ SupportsTwoPointTargetTemperature() bool }); ok {
					res.SupportsTwoPointTargetTemperature = c.SupportsTwoPointTargetTemperature()
				}

				if c, ok := typed.(interface{ SupportedModes() []entity.ClimateMode }); ok {
					for _, mode := range c.SupportedModes() {
						res.SupportedModes = append(res.SupportedModes, common.Enum[ehp.ClimateMode](mode))
					}
				}

				if c, ok := typed.(interface{ VisualMinTemperature() float32 }); ok {
					res.VisualMinTemperature = c.VisualMinTemperature()
				}

				if c, ok := typed.(interface{ VisualMaxTemperature() float32 }); ok {
					res.VisualMaxTemperature = c.VisualMaxTemperature()
				}

				if c, ok := typed.(interface{ VisualTargetTemperatureStep() float32 }); ok {
					res.VisualTargetTemperatureStep = c.VisualTargetTemperatureStep()
				}

				if c, ok := typed.(interface{ LegacySupportsAway() bool }); ok {
					res.LegacySupportsAway = c.LegacySupportsAway()
				}

				if c, ok := typed.(interface{ SupportsAction() bool }); ok {
					res.SupportsAction = c.SupportsAction()
				}

				if c, ok := typed.(interface {
					SupportedFanModes() []entity.ClimateFanMode
				}); ok {
					for _, mode := range c.SupportedFanModes() {
						res.SupportedFanModes = append(res.SupportedFanModes, common.Enum[ehp.ClimateFanMode](mode))
					}
				}

				if c, ok := typed.(interface {
					SupportedSwingModes() []entity.ClimateSwingMode
				}); ok {
					for _, mode := range c.SupportedSwingModes() {
						res.SupportedSwingModes = append(res.SupportedSwingModes, common.Enum[ehp.ClimateSwingMode](mode))
					}
				}

				if c, ok := typed.(interface{ SupportedCustomFanModes() []string }); ok {
					res.SupportedCustomFanModes = c.SupportedCustomFanModes()
				}

				if c, ok := typed.(interface{ SupportedPresets() []entity.ClimatePreset }); ok {
					for _, preset := range c.SupportedPresets() {
						res.SupportedPresets = append(res.SupportedPresets, common.Enum[ehp.ClimatePreset](preset))
					}
				}

				if c, ok := typed.(interface{ SupportedCustomPresets() []string }); ok {
					res.SupportedCustomPresets = c.SupportedCustomPresets()
				}

				if c, ok := typed.(interface{ VisualCurrentTemperatureStep() float32 }); ok {
					res.VisualCurrentTemperatureStep = c.VisualCurrentTemperatureStep()
				}

				if c, ok := typed.(interface{ SupportsCurrentHumidity() bool }); ok {
					res.SupportsCurrentHumidity = c.SupportsCurrentHumidity()
				}

				if c, ok := typed.(interface{ SupportsTargetHumidity() bool }); ok {
					res.SupportsTargetHumidity = c.SupportsTargetHumidity()
				}

				if c, ok := typed.(interface{ VisualMinHumidity() float32 }); ok {
					res.VisualMinHumidity = c.VisualMinHumidity()
				}

				if c, ok := typed.(interface{ VisualMaxHumidity() float32 }); ok {
					res.VisualMaxHumidity = c.VisualMaxHumidity()
				}

				ret = append(ret, &res)
			case entity.Number:
				ret = append(ret, &ehp.ListEntitiesNumberResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
					UnitOfMeasurement: typed.UnitOfMeasurement(),
					MinValue:          typed.MinValue(),
					MaxValue:          typed.MaxValue(),
					Step:              typed.Step(),
					Mode:              ehp.NumberMode(typed.NumberMode()),
				})
			case entity.Date:
				ret = append(ret, &ehp.ListEntitiesDateResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Time:
				ret = append(ret, &ehp.ListEntitiesTimeResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Datetime:
				ret = append(ret, &ehp.ListEntitiesDateTimeResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Text:
				ret = append(ret, &ehp.ListEntitiesTextResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Select:
				ret = append(ret, &ehp.ListEntitiesSelectResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Siren:
				ret = append(ret, &ehp.ListEntitiesSirenResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					Tones:             typed.Tones(),
					SupportsDuration:  typed.SupportsDuration(),
					SupportsVolume:    typed.SupportsVolume(),
				})
			case entity.Lock:
				ret = append(ret, &ehp.ListEntitiesLockResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Valve:
				ret = append(ret, &ehp.ListEntitiesValveResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
				})
			case entity.MediaPlayer:
				ret = append(ret, &ehp.ListEntitiesMediaPlayerResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.AlarmControlPanel:
				ret = append(ret, &ehp.ListEntitiesAlarmControlPanelResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
				})
			case entity.Event:
				ret = append(ret, &ehp.ListEntitiesEventResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
				})
			case entity.Update:
				ret = append(ret, &ehp.ListEntitiesUpdateResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    common.Enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					UniqueId:          node.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       string(typed.DeviceClass()),
				})
			case entity.Service:
				ret = append(ret, &ehp.ListEntitiesServicesResponse{
					Key:  typed.HashID(),
					Name: typed.Name(),
					Args: nil, // TODO: fill
				})
			default:
			}
		}
		return append(ret, &ehp.ListEntitiesDoneResponse{}), nil
	})))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SubscribeStatesRequest) ([]ehp.EsphomeMessageTyper, error) {
		c.subscribed = true
		b := bus.Get(ctx)
		sub := b.HandleEvents(bus.EventHandler(func(t *bus.StateChangeEvent) {
			r := stateResponse(t.Key, t.NewState)
			if r != nil {
				slog.Debug("Sending state change", "key", t.Key, "state", t.NewState, "to", c.clientInfo)
				err := c.SendMessages([]ehp.EsphomeMessageTyper{r})
				if err != nil {
					slog.Error("Failed to send state update", "err", err)
				}
			}
		}))
		c.busEvents = append(c.busEvents, sub)
		ret := []ehp.EsphomeMessageTyper{}
		for _, ent := range entity.IterateRegistry(core.GetNode(ctx).Registry) {
			es := entityState(ent)
			if es == nil {
				continue
			}
			ret = append(ret, es)
		}

		return ret, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SubscribeLogsRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_logs, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SubscribeHomeassistantServicesRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_homeassistant_services, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SubscribeHomeAssistantStatesRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_home_assistant_states, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.ExecuteServiceRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command execute_service, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.CoverCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command cover_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.FanCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command fan_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.LightCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command light_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SwitchCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		core.GetNode(ctx).Bus.CallService(&switchcomp.SetState{
			Key:   msg.Key,
			State: msg.State,
		})
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.CameraImageRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command camera_image, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.ClimateCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		core.GetNode(ctx).Bus.CallService(&climate.SetState{
			Key:                      msg.Key,
			HasMode:                  msg.HasMode,
			Mode:                     entity.ClimateMode(msg.Mode),
			HasTargetTemperature:     msg.HasTargetTemperature,
			TargetTemperature:        msg.TargetTemperature,
			HasTargetTemperatureLow:  msg.HasTargetTemperatureLow,
			TargetTemperatureLow:     msg.TargetTemperatureLow,
			HasTargetTemperatureHigh: msg.HasTargetTemperatureHigh,
			TargetTemperatureHigh:    msg.TargetTemperatureHigh,
			HasLegacyAway:            msg.HasLegacyAway,
			LegacyAway:               msg.LegacyAway,
			HasFanMode:               msg.HasFanMode,
			FanMode:                  entity.ClimateFanMode(msg.FanMode),
			HasSwingMode:             msg.HasSwingMode,
			SwingMode:                entity.ClimateSwingMode(msg.SwingMode),
			HasCustomFanMode:         msg.HasCustomFanMode,
			CustomFanMode:            msg.CustomFanMode,
			HasPreset:                msg.HasPreset,
			Preset:                   entity.ClimatePreset(msg.Preset),
			HasCustomPreset:          msg.HasCustomPreset,
			CustomPreset:             msg.CustomPreset,
			HasTargetHumidity:        msg.HasTargetHumidity,
			TargetHumidity:           msg.TargetHumidity,
		})
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.NumberCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		core.GetNode(ctx).Bus.CallService(&number.SetValue{
			Key:   msg.Key,
			State: msg.State,
		})
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SelectCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command select_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.TextCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command text_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SirenCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command siren_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.ButtonCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		core.GetNode(ctx).Bus.CallService(&button.ButtonPress{
			Key: msg.Key,
		})
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.LockCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command lock_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.ValveCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command valve_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.MediaPlayerCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command media_player_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.DateCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command date_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.TimeCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command time_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.DateTimeCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command datetime_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SubscribeBluetoothLEAdvertisementsRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_bluetooth_le_advertisements, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.BluetoothDeviceRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_device_request, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTGetServicesRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_get_services, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTReadRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_read, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTWriteRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_write, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTReadDescriptorRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_read_descriptor, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTWriteDescriptorRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_write_descriptor, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTNotifyRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_notify, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.UnsubscribeBluetoothLEAdvertisementsRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command unsubscribe_bluetooth_le_advertisements, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.SubscribeVoiceAssistantRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_voice_assistant, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler(func(ctx context.Context, c *Connection, msg *ehp.AlarmControlPanelCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command alarm_control_panel_command, doing nothing")
		return nil, nil
	}))
)
