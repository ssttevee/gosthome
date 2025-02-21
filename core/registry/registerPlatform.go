package registry

import (
	"errors"

	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/entity"
)

func (r *Registry) tryRegisterPlatforms(name string, cd component.Declaration) error {
	errs := []error{}
	if p, ok := cd.(entity.BinarySensorPlatformer); ok {
		ec := p.BinarySensorPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeBinarySensor, name, ec))
		}
	}
	if p, ok := cd.(entity.CoverPlatformer); ok {
		ec := p.CoverPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeCover, name, ec))
		}
	}
	if p, ok := cd.(entity.FanPlatformer); ok {
		ec := p.FanPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeFan, name, ec))
		}
	}
	if p, ok := cd.(entity.LightPlatformer); ok {
		ec := p.LightPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeLight, name, ec))
		}
	}
	if p, ok := cd.(entity.SensorPlatformer); ok {
		ec := p.SensorPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeSensor, name, ec))
		}
	}
	if p, ok := cd.(entity.SwitchPlatformer); ok {
		ec := p.SwitchPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeSwitch, name, ec))
		}
	}
	if p, ok := cd.(entity.ButtonPlatformer); ok {
		ec := p.ButtonPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeButton, name, ec))
		}
	}
	if p, ok := cd.(entity.TextSensorPlatformer); ok {
		ec := p.TextSensorPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeTextSensor, name, ec))
		}
	}
	if p, ok := cd.(entity.ServicePlatformer); ok {
		ec := p.ServicePlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeService, name, ec))
		}
	}
	if p, ok := cd.(entity.CameraPlatformer); ok {
		ec := p.CameraPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeCamera, name, ec))
		}
	}
	if p, ok := cd.(entity.ClimatePlatformer); ok {
		ec := p.ClimatePlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeClimate, name, ec))
		}
	}
	if p, ok := cd.(entity.NumberPlatformer); ok {
		ec := p.NumberPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeNumber, name, ec))
		}
	}
	if p, ok := cd.(entity.DatePlatformer); ok {
		ec := p.DatePlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeDatetimeDate, name, ec))
		}
	}
	if p, ok := cd.(entity.TimePlatformer); ok {
		ec := p.TimePlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeDatetimeTime, name, ec))
		}
	}
	if p, ok := cd.(entity.DatetimePlatformer); ok {
		ec := p.DatetimePlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeDatetimeDatetime, name, ec))
		}
	}
	if p, ok := cd.(entity.TextPlatformer); ok {
		ec := p.TextPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeText, name, ec))
		}
	}
	if p, ok := cd.(entity.SelectPlatformer); ok {
		ec := p.SelectPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeSelect, name, ec))
		}
	}
	if p, ok := cd.(entity.LockPlatformer); ok {
		ec := p.LockPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeLock, name, ec))
		}
	}
	if p, ok := cd.(entity.ValvePlatformer); ok {
		ec := p.ValvePlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeValve, name, ec))
		}
	}
	if p, ok := cd.(entity.MediaPlayerPlatformer); ok {
		ec := p.MediaPlayerPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeMediaPlayer, name, ec))
		}
	}
	if p, ok := cd.(entity.AlarmControlPanelPlatformer); ok {
		ec := p.AlarmControlPanelPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeAlarmControlPanel, name, ec))
		}
	}
	if p, ok := cd.(entity.EventPlatformer); ok {
		ec := p.EventPlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeEvent, name, ec))
		}
	}
	if p, ok := cd.(entity.UpdatePlatformer); ok {
		ec := p.UpdatePlatform()
		if ec != nil {
			errs = append(errs, r.RegisterEntityComponent(entity.DomainTypeUpdate, name, ec))
		}
	}
	return errors.Join(errs...)
}
