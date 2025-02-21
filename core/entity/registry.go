package entity

import (
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/guarded"
)

type Registry struct {
	binarySensorDomain      atomic.Pointer[BinarySensorDomain]
	coverDomain             atomic.Pointer[CoverDomain]
	fanDomain               atomic.Pointer[FanDomain]
	lightDomain             atomic.Pointer[LightDomain]
	sensorDomain            atomic.Pointer[SensorDomain]
	switchDomain            atomic.Pointer[SwitchDomain]
	buttonDomain            atomic.Pointer[ButtonDomain]
	textSensorDomain        atomic.Pointer[TextSensorDomain]
	serviceDomain           atomic.Pointer[ServiceDomain]
	cameraDomain            atomic.Pointer[CameraDomain]
	climateDomain           atomic.Pointer[ClimateDomain]
	numberDomain            atomic.Pointer[NumberDomain]
	dateDomain              atomic.Pointer[DateDomain]
	timeDomain              atomic.Pointer[TimeDomain]
	datetimeDomain          atomic.Pointer[DatetimeDomain]
	textDomain              atomic.Pointer[TextDomain]
	selectDomain            atomic.Pointer[SelectDomain]
	sirenDomain             atomic.Pointer[SirenDomain]
	lockDomain              atomic.Pointer[LockDomain]
	valveDomain             atomic.Pointer[ValveDomain]
	mediaPlayerDomain       atomic.Pointer[MediaPlayerDomain]
	alarmControlPanelDomain atomic.Pointer[AlarmControlPanelDomain]
	eventDomain             atomic.Pointer[EventDomain]
	updateDomain            atomic.Pointer[UpdateDomain]

	external guarded.RWValue[map[string]component.Component]
}

type ErrAlreadyRegistered DomainType

func (e ErrAlreadyRegistered) Error() string {
	return fmt.Sprintf("domain already registered: %s", DomainType(e))
}

type DomainDefinition interface {
	create(reg *Registry) error
}

type publicDomainDefinition struct {
	d DomainTyper
}

func PublicDomain[T any, PT interface {
	*T
	DomainTyper
}](dom PT) publicDomainDefinition {
	return publicDomainDefinition{
		d: dom,
	}
}

func createPD[T any, PT interface {
	*T
	DomainTyper
}](ptr *atomic.Pointer[T], t *T) error {
	if !ptr.CompareAndSwap(nil, t) {
		return ErrAlreadyRegistered(PT(nil).DomainType())
	}
	return nil
}

func (pdd publicDomainDefinition) create(reg *Registry) error {
	switch d := pdd.d.(type) {
	case *BinarySensorDomain:
		return createPD(&reg.binarySensorDomain, d)
	case *CoverDomain:
		return createPD(&reg.coverDomain, d)
	case *FanDomain:
		return createPD(&reg.fanDomain, d)
	case *LightDomain:
		return createPD(&reg.lightDomain, d)
	case *SensorDomain:
		return createPD(&reg.sensorDomain, d)
	case *SwitchDomain:
		return createPD(&reg.switchDomain, d)
	case *ButtonDomain:
		return createPD(&reg.buttonDomain, d)
	case *TextSensorDomain:
		return createPD(&reg.textSensorDomain, d)
	case *ServiceDomain:
		return createPD(&reg.serviceDomain, d)
	case *CameraDomain:
		return createPD(&reg.cameraDomain, d)
	case *ClimateDomain:
		return createPD(&reg.climateDomain, d)
	case *NumberDomain:
		return createPD(&reg.numberDomain, d)
	case *DateDomain:
		return createPD(&reg.dateDomain, d)
	case *TimeDomain:
		return createPD(&reg.timeDomain, d)
	case *DatetimeDomain:
		return createPD(&reg.datetimeDomain, d)
	case *TextDomain:
		return createPD(&reg.textDomain, d)
	case *SelectDomain:
		return createPD(&reg.selectDomain, d)
	case *SirenDomain:
		return createPD(&reg.sirenDomain, d)
	case *LockDomain:
		return createPD(&reg.lockDomain, d)
	case *ValveDomain:
		return createPD(&reg.valveDomain, d)
	case *MediaPlayerDomain:
		return createPD(&reg.mediaPlayerDomain, d)
	case *AlarmControlPanelDomain:
		return createPD(&reg.alarmControlPanelDomain, d)
	case *EventDomain:
		return createPD(&reg.eventDomain, d)
	case *UpdateDomain:
		return createPD(&reg.updateDomain, d)
	default:
		panic("unknown domain")
	}
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (er *Registry) CreateDomain(f DomainDefinition) error {
	return f.create(er)
}

func (er *Registry) BinarySensorByKey(key uint32) (BinarySensor, bool) {
	d := er.binarySensorDomain.Load()
	if d == nil {
		slog.Error("BinarySensorDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) BinarySensors() []Entity {
	d := er.binarySensorDomain.Load()
	if d == nil {
		slog.Debug("BinarySensorDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) CoverByKey(key uint32) (Cover, bool) {
	d := er.coverDomain.Load()
	if d == nil {
		slog.Error("CoverDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Covers() []Entity {
	d := er.coverDomain.Load()
	if d == nil {
		slog.Debug("CoverDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) FanByKey(key uint32) (Fan, bool) {
	d := er.fanDomain.Load()
	if d == nil {
		slog.Error("FanDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Fans() []Entity {
	d := er.fanDomain.Load()
	if d == nil {
		slog.Debug("FanDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) LightByKey(key uint32) (Light, bool) {
	d := er.lightDomain.Load()
	if d == nil {
		slog.Error("LightDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Lights() []Entity {
	d := er.lightDomain.Load()
	if d == nil {
		slog.Debug("LightDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) SensorByKey(key uint32) (Sensor, bool) {
	d := er.sensorDomain.Load()
	if d == nil {
		slog.Error("SensorDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Sensors() []Entity {
	d := er.sensorDomain.Load()
	if d == nil {
		slog.Debug("SensorDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) SwitchByKey(key uint32) (Switch, bool) {
	d := er.switchDomain.Load()
	if d == nil {
		slog.Error("SwitcheDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Switches() []Entity {
	d := er.switchDomain.Load()
	if d == nil {
		slog.Debug("SwitcheDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) ButtonByKey(key uint32) (Button, bool) {
	d := er.buttonDomain.Load()
	if d == nil {
		slog.Error("ButtonDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Buttons() []Entity {
	d := er.buttonDomain.Load()
	if d == nil {
		slog.Debug("ButtonDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) TextSensorByKey(key uint32) (TextSensor, bool) {
	d := er.textSensorDomain.Load()
	if d == nil {
		slog.Error("TextSensorDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) TextSensors() []Entity {
	d := er.textSensorDomain.Load()
	if d == nil {
		slog.Debug("TextSensorDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) ServiceByKey(key uint32) (Service, bool) {
	d := er.serviceDomain.Load()
	if d == nil {
		slog.Error("ServiceDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Services() []Entity {
	d := er.serviceDomain.Load()
	if d == nil {
		slog.Debug("ServiceDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) CameraByKey(key uint32) (Camera, bool) {
	d := er.cameraDomain.Load()
	if d == nil {
		slog.Error("CameraDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Cameras() []Entity {
	d := er.cameraDomain.Load()
	if d == nil {
		slog.Debug("CameraDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) ClimateByKey(key uint32) (Climate, bool) {
	d := er.climateDomain.Load()
	if d == nil {
		slog.Error("ClimateDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Climates() []Entity {
	d := er.climateDomain.Load()
	if d == nil {
		slog.Debug("ClimateDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) NumberByKey(key uint32) (Number, bool) {
	d := er.numberDomain.Load()
	if d == nil {
		slog.Error("NumberDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Numbers() []Entity {
	d := er.numberDomain.Load()
	if d == nil {
		slog.Debug("NumberDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) DateByKey(key uint32) (Date, bool) {
	d := er.dateDomain.Load()
	if d == nil {
		slog.Error("DateDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Dates() []Entity {
	d := er.dateDomain.Load()
	if d == nil {
		slog.Debug("DateDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) TimeByKey(key uint32) (Time, bool) {
	d := er.timeDomain.Load()
	if d == nil {
		slog.Error("TimeDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Times() []Entity {
	d := er.timeDomain.Load()
	if d == nil {
		slog.Debug("TimeDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) DatetimeByKey(key uint32) (Datetime, bool) {
	d := er.datetimeDomain.Load()
	if d == nil {
		slog.Error("DatetimeDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Datetimes() []Entity {
	d := er.datetimeDomain.Load()
	if d == nil {
		slog.Debug("DatetimeDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) TextByKey(key uint32) (Text, bool) {
	d := er.textDomain.Load()
	if d == nil {
		slog.Error("TextDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Texts() []Entity {
	d := er.textDomain.Load()
	if d == nil {
		slog.Debug("TextDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) SelectByKey(key uint32) (Select, bool) {
	d := er.selectDomain.Load()
	if d == nil {
		slog.Error("SelectDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Selects() []Entity {
	d := er.selectDomain.Load()
	if d == nil {
		slog.Debug("SelectDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) SirenByKey(key uint32) (Siren, bool) {
	d := er.sirenDomain.Load()
	if d == nil {
		slog.Error("SirenDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Sirens() []Entity {
	d := er.sirenDomain.Load()
	if d == nil {
		slog.Debug("SirenDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) LockByKey(key uint32) (Lock, bool) {
	d := er.lockDomain.Load()
	if d == nil {
		slog.Error("LockDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Locks() []Entity {
	d := er.lockDomain.Load()
	if d == nil {
		slog.Debug("LockDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) ValveByKey(key uint32) (Valve, bool) {
	d := er.valveDomain.Load()
	if d == nil {
		slog.Error("ValveDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Valves() []Entity {
	d := er.valveDomain.Load()
	if d == nil {
		slog.Debug("ValveDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) MediaPlayerByKey(key uint32) (MediaPlayer, bool) {
	d := er.mediaPlayerDomain.Load()
	if d == nil {
		slog.Error("MediaPlayerDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) MediaPlayers() []Entity {
	d := er.mediaPlayerDomain.Load()
	if d == nil {
		slog.Debug("MediaPlayerDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) AlarmControlPanelByKey(key uint32) (AlarmControlPanel, bool) {
	d := er.alarmControlPanelDomain.Load()
	if d == nil {
		slog.Error("AlarmControlPanelDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) AlarmControlPanels() []Entity {
	d := er.alarmControlPanelDomain.Load()
	if d == nil {
		slog.Debug("AlarmControlPanelDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) EventByKey(key uint32) (Event, bool) {
	d := er.eventDomain.Load()
	if d == nil {
		slog.Error("EventDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Events() []Entity {
	d := er.eventDomain.Load()
	if d == nil {
		slog.Debug("EventDomain is not registered!")
		return nil
	}
	return d.Clone()
}
func (er *Registry) UpdateByKey(key uint32) (Update, bool) {
	d := er.updateDomain.Load()
	if d == nil {
		slog.Error("UpdateDomain is not registered!")
		return nil, false
	}
	return d.FindByKey(key)
}
func (er *Registry) Updates() []Entity {
	d := er.updateDomain.Load()
	if d == nil {
		slog.Debug("UpdateDomain is not registered!")
		return nil
	}
	return d.Clone()
}

func (er *Registry) RegisterBinarySensor(ent BinarySensor) (err error) {
	d := er.binarySensorDomain.Load()
	if d == nil {
		return errors.New("BinarySensorDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterCover(ent Cover) (err error) {
	d := er.coverDomain.Load()
	if d == nil {
		return errors.New("CoverDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterFan(ent Fan) (err error) {
	d := er.fanDomain.Load()
	if d == nil {
		return errors.New("FanDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterLight(ent Light) (err error) {
	d := er.lightDomain.Load()
	if d == nil {
		return errors.New("LightDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterSensor(ent Sensor) (err error) {
	d := er.sensorDomain.Load()
	if d == nil {
		return errors.New("SensorDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterSwitch(ent Switch) (err error) {
	d := er.switchDomain.Load()
	if d == nil {
		return errors.New("SwitcheDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterButton(ent Button) (err error) {
	d := er.buttonDomain.Load()
	if d == nil {
		return errors.New("ButtonDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterTextSensor(ent TextSensor) (err error) {
	d := er.textSensorDomain.Load()
	if d == nil {
		return errors.New("TextSensorDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterService(ent Service) (err error) {
	d := er.serviceDomain.Load()
	if d == nil {
		return errors.New("CameraDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterCamera(ent Camera) (err error) {
	d := er.cameraDomain.Load()
	if d == nil {
		return errors.New("CameraDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterClimate(ent Climate) (err error) {
	d := er.climateDomain.Load()
	if d == nil {
		return errors.New("ClimateDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterNumber(ent Number) (err error) {
	d := er.numberDomain.Load()
	if d == nil {
		return errors.New("NumberDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterDate(ent Date) (err error) {
	d := er.dateDomain.Load()
	if d == nil {
		return errors.New("DateDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterTime(ent Time) (err error) {
	d := er.timeDomain.Load()
	if d == nil {
		return errors.New("TimeDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterDatetime(ent Datetime) (err error) {
	d := er.datetimeDomain.Load()
	if d == nil {
		return errors.New("DatetimeDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterText(ent Text) (err error) {
	d := er.textDomain.Load()
	if d == nil {
		return errors.New("TextDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterSelect(ent Select) (err error) {
	d := er.selectDomain.Load()
	if d == nil {
		return errors.New("SelectDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterSiren(ent Siren) (err error) {
	d := er.sirenDomain.Load()
	if d == nil {
		return errors.New("SirenDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterLock(ent Lock) (err error) {
	d := er.lockDomain.Load()
	if d == nil {
		return errors.New("LockDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterValve(ent Valve) (err error) {
	d := er.valveDomain.Load()
	if d == nil {
		return errors.New("ValveDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterMediaPlayer(ent MediaPlayer) (err error) {
	d := er.mediaPlayerDomain.Load()
	if d == nil {
		return errors.New("MediaPlayerDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterAlarmControlPanel(ent AlarmControlPanel) (err error) {
	d := er.alarmControlPanelDomain.Load()
	if d == nil {
		return errors.New("AlarmControlPanelDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterEvent(ent Event) (err error) {
	d := er.eventDomain.Load()
	if d == nil {
		return errors.New("EventDomain is not registered!")
	}
	return d.Register(ent)
}

func (er *Registry) RegisterUpdate(ent Update) (err error) {
	d := er.updateDomain.Load()
	if d == nil {
		return errors.New("UpdateDomain is not registered!")
	}
	return d.Register(ent)
}
