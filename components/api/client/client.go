package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"reflect"
	"sync"

	"github.com/gosthome/gosthome/components/api/common"
	ehp "github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/core/component/logger"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/guarded"
	"github.com/majfault/signal"
)

type connectState int

const (
	connectStateInvalid connectState = iota
	connectStateStart
	connectStateHandshake
	connectStateConnecting
	connectStateReady
	connectStateDisconnecting
	connectStateError
)

func (s connectState) Error() string {
	switch s {
	case connectStateStart:
		return "state was start"
	case connectStateHandshake:
		return "state was handshake"
	case connectStateConnecting:
		return "state was connecting"
	case connectStateReady:
		return "state was ready"
	case connectStateDisconnecting:
		return "state was disconnecting"
	case connectStateError:
		return "state was error"
	}
	return "unknown state"
}

type Client struct {
	reg entity.Registry

	dialer   common.Dialer
	shaker   frameshakers.ClientShaker
	ctx      context.Context
	cancel   context.CancelFunc
	address  string
	port     uint16
	password string
	psk      *frameshakers.ConfigNoisePSK

	conn       io.Closer
	wg         sync.WaitGroup
	sendFrames frameshakers.FrameSenderFunc

	stateRead         <-chan error
	stateWrite        guarded.Value[chan<- error]
	listEntitiesState guarded.Value[chan<- struct{}]
	logs              LogsSignal

	OnClose func()
}

type ClientOpt func(*Client)

func WithDialer(d common.Dialer) ClientOpt {
	return func(c *Client) {
		c.dialer = d
	}
}

func WithPassword(pwd string) ClientOpt {
	return func(c *Client) {
		c.password = pwd
	}
}

func WithNoisePSK(pwd *frameshakers.ConfigNoisePSK) ClientOpt {
	return func(c *Client) {
		c.shaker = frameshakers.NoiseClient
		c.psk = pwd
		c.ctx = frameshakers.ContextWithValue(c.ctx, "noisePSK", c.psk)
	}
}

func New(ctx context.Context, address string, port uint16, opts ...ClientOpt) *Client {
	ctx, cancel := context.WithCancel(ctx)
	c := &Client{
		dialer:  common.DialTCP,
		shaker:  frameshakers.PlaintextClient,
		ctx:     ctx,
		cancel:  cancel,
		address: address,
		port:    port,
	}
	for _, o := range opts {
		o(c)
	}
	c.reg.CreateDomain(entity.PublicDomain(&entity.BinarySensorDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.AlarmControlPanelDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.CoverDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.FanDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.LightDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.SensorDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.SwitchDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.TextSensorDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.ServiceDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.CameraDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.ClimateDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.NumberDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.SelectDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.SirenDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.LockDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.ButtonDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.MediaPlayerDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.TextDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.DateDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.TimeDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.EventDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.ValveDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.DatetimeDomain{}))
	c.reg.CreateDomain(entity.PublicDomain(&entity.UpdateDomain{}))
	return c
}

func (c *Client) Connect() error {
	c.stateWrite.Do(func(ch *chan<- error) {
		rwc := make(chan error, 1)
		*ch = rwc
		c.stateRead = rwc
	})
	slog.Debug("Connecting", "address", c.address, "port", c.port)
	conn, err := c.dialer(c.ctx, fmt.Sprintf("%s:%d", c.address, c.port))
	if err != nil {
		return err
	}
	slog.Debug("Handshaking", "address", c.address, "port", c.port)
	c.conn = conn
	err = c.handshake()
	if err != nil {
		return err
	}

	err = c.sendMessages(&ehp.HelloRequest{
		ClientInfo:      "gosthome client",
		ApiVersionMajor: common.ApiVersionMajor,
		ApiVersionMinor: common.ApiVersionMinor,
	}, &ehp.ConnectRequest{
		Password: c.password,
	})
	if err != nil {
		return err
	}
	slog.Debug("Wainting for connect to succeed")
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case err = <-c.stateRead:
			if state, ok := err.(connectState); ok {
				if state == connectStateConnecting {
					slog.Debug("Recieved hello")
					continue
				}
				if state == connectStateReady {
					slog.Debug("Recieved connected")
					return nil
				}
			}
			return fmt.Errorf("error druring handshake %w", err)
		}
	}
}

func (c *Client) handshake() error {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer c.close()
		r, w := frameshakers.SplitConnection(c.conn.(net.Conn))
		neerr := c.shaker(c.ctx, r, w, func(sendFrames frameshakers.FrameSenderFunc) (handler frameshakers.FrameSenderFunc, err error) {
			c.stateWrite.Do(func(r *chan<- error) {
				if *r == nil {
					err = errors.New("wrong connection state")
					return
				}
				slog.Debug("handshake done")
				*r <- connectStateHandshake
				slog.Debug("handshake done")
			})
			if err != nil {
				return nil, err
			}
			c.sendFrames = sendFrames
			return c.handleFrames, nil
		})
		if neerr != nil {
			c.stateWrite.Do(func(r *chan<- error) {
				if *r == nil {
					slog.Error("Failed to connect", "err", neerr)
					return
				}
				*r <- neerr
				close(*r)
				*r = nil
			})
		}
	}()
	slog.Debug("waiting for handshake")
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case err := <-c.stateRead:
		if state, ok := err.(connectState); ok && state == connectStateHandshake {
			return nil
		}
		return fmt.Errorf("error during handshake %w", err)
	}
}
func (c *Client) Close() (err error) {
	defer slog.Debug("Client closed", "err", err)
	c.cancel()
	if c.OnClose != nil {
		c.OnClose()
	}
	c.wg.Wait()
	errs := []error{}
	ok := false
	for {
		select {
		case err, ok = <-c.stateRead:
			if !ok {
				err = errors.Join(errs...)
				return err
			}
			if state, ok := err.(connectState); ok {
				if state != connectStateError {
					continue
				}
			}
			err = errors.Join(errs...)
			return err
		default:
			err = errors.Join(errs...)
			return err
		}
	}
}

func (c *Client) close() error {
	c.stateWrite.Do(func(r *chan<- error) {
		if *r != nil {
			close(*r)
			*r = nil
		}
	})
	c.listEntitiesState.Do(func(state *chan<- struct{}) {
		if *state != nil {
			close(*state)
			*state = nil
		}
	})
	c.logs.Close()
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) sendMessages(msgs ...ehp.EsphomeMessageTyper) error {
	if c.sendFrames == nil {
		return errors.New("Connection is not established yet")
	}
	frames, err := common.EncodeFrames(msgs)
	if err != nil {
		return err
	}
	return c.sendFrames(frames)
}

func (c *Client) handleFrames(input []frameshakers.Frame) (err error) {
	defer func() {
		if err != nil {
			slog.Error("error handling frames", "err", err)
		}
	}()
	closing := false
	for _, frame := range input {
		mt, msg, err := common.DecodeFrame(frame)
		if err != nil {
			return err
		}
		if slog.Default().Enabled(c.ctx, slog.LevelDebug) {
			slog.Default().Debug("client handles message from server", "msg", reflect.TypeOf(msg))
		}
		src := msg.EsphomeSource()
		if src == ehp.APISourceType_SOURCE_CLIENT {
			return errors.New("unexpected client message")
		}
		slog.Info("client recieved", "frame", mt, "msg", msg)
		switch mt {
		case ehp.MessageTypeHelloResponse:
			c.stateWrite.Do(func(r *chan<- error) {
				*r <- connectStateConnecting
			})
		case ehp.MessageTypeConnectResponse:
			resp := msg.(*ehp.ConnectResponse)
			if resp.InvalidPassword {
				return errors.New("unauthrized")
			}
			c.stateWrite.Do(func(r *chan<- error) {
				*r <- connectStateReady
			})
			continue
		case ehp.MessageTypeDeviceInfoResponse:
		case ehp.MessageTypeDisconnectRequest:
			c.sendMessages(&ehp.DisconnectResponse{})

			break
		case ehp.MessageTypeDisconnectResponse:
			slog.Error("Server is disconnecting")
			closing = true
			break
		case ehp.MessageTypePingRequest:
			c.sendMessages(&ehp.PingResponse{})
			continue
		case ehp.MessageTypePingResponse:
			// 	slog.Error("Dont know how to handle PingResponse")
			continue
		case ehp.MessageTypeGetTimeRequest:
		case ehp.MessageTypeGetTimeResponse:
		case ehp.MessageTypeSubscribeLogsResponse:
			log := msg.(*ehp.SubscribeLogsResponse)
			c.logs.Emit(logger.ParseLevelFromInt(log.Level), log.Message)

		// VoiceAssistantAudio
		case
			ehp.MessageTypeListEntitiesBinarySensorResponse,
			ehp.MessageTypeListEntitiesCoverResponse,
			ehp.MessageTypeListEntitiesFanResponse,
			ehp.MessageTypeListEntitiesLightResponse,
			ehp.MessageTypeListEntitiesSensorResponse,
			ehp.MessageTypeListEntitiesSwitchResponse,
			ehp.MessageTypeListEntitiesTextSensorResponse,
			ehp.MessageTypeListEntitiesServicesResponse,
			ehp.MessageTypeListEntitiesCameraResponse,
			ehp.MessageTypeListEntitiesClimateResponse,
			ehp.MessageTypeListEntitiesNumberResponse,
			ehp.MessageTypeListEntitiesSelectResponse,
			ehp.MessageTypeListEntitiesSirenResponse,
			ehp.MessageTypeListEntitiesLockResponse,
			ehp.MessageTypeListEntitiesButtonResponse,
			ehp.MessageTypeListEntitiesMediaPlayerResponse,
			ehp.MessageTypeListEntitiesTextResponse,
			ehp.MessageTypeListEntitiesDateResponse,
			ehp.MessageTypeListEntitiesTimeResponse,
			ehp.MessageTypeListEntitiesEventResponse,
			ehp.MessageTypeListEntitiesValveResponse,
			ehp.MessageTypeListEntitiesDateTimeResponse,
			ehp.MessageTypeListEntitiesUpdateResponse,
			ehp.MessageTypeListEntitiesAlarmControlPanelResponse,
			ehp.MessageTypeListEntitiesDoneResponse:
			err = c.listEntitiesResponse(msg)
			if err != nil {
				return err
			}
			continue
		case
			ehp.MessageTypeBinarySensorStateResponse,
			ehp.MessageTypeCoverStateResponse,
			ehp.MessageTypeFanStateResponse,
			ehp.MessageTypeLightStateResponse,
			ehp.MessageTypeSensorStateResponse,
			ehp.MessageTypeSwitchStateResponse,
			ehp.MessageTypeTextSensorStateResponse,
			ehp.MessageTypeSubscribeHomeAssistantStateResponse,
			ehp.MessageTypeClimateStateResponse,
			ehp.MessageTypeNumberStateResponse,
			ehp.MessageTypeSelectStateResponse,
			ehp.MessageTypeSirenStateResponse,
			ehp.MessageTypeLockStateResponse,
			ehp.MessageTypeMediaPlayerStateResponse,
			ehp.MessageTypeAlarmControlPanelStateResponse,
			ehp.MessageTypeTextStateResponse,
			ehp.MessageTypeDateStateResponse,
			ehp.MessageTypeTimeStateResponse,
			ehp.MessageTypeValveStateResponse,
			ehp.MessageTypeDateTimeStateResponse,
			ehp.MessageTypeUpdateStateResponse:
			err = c.stateChangeResponse(msg)
			if err != nil {
				return err
			}
			continue
		default:
			slog.Error("Dont know how to handle unknown message", "type", fmt.Sprintf("%T", msg))
			// SubscribeLogsResponse
			// HomeassistantServiceResponse
			// CameraImageResponse
			// BluetoothLEAdvertisementResponse
			// BluetoothDeviceConnectionResponse
			// BluetoothGATTGetServicesResponse
			// BluetoothGATTGetServicesDoneResponse
			// BluetoothGATTReadResponse
			// BluetoothGATTNotifyDataResponse
			// BluetoothConnectionsFreeResponse
			// BluetoothGATTErrorResponse
			// BluetoothGATTWriteResponse
			// BluetoothGATTNotifyResponse
			// BluetoothDevicePairingResponse
			// BluetoothDeviceUnpairingResponse
			// BluetoothDeviceClearCacheResponse
			// VoiceAssistantRequest
			// BluetoothLERawAdvertisementsResponse
			// EventResponse
			// VoiceAssistantAnnounceFinished
			// VoiceAssistantConfigurationResponse
		}
	}
	if closing {
		err = frameshakers.ErrCloseConnection
	}
	return
}

type (
	LogsSignal = signal.Signal2[logger.Level, []byte]
	LogsSlot   = signal.Slot2[logger.Level, []byte]
)

func (c *Client) StartLogs() error {
	return c.sendMessages(&ehp.SubscribeLogsRequest{
		Level:      ehp.LogLevel_LOG_LEVEL_VERBOSE,
		DumpConfig: true,
	})
}

func (c *Client) Logs() *LogsSignal {
	return &c.logs
}
