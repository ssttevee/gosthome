package psutil

import (
	"context"
	"log/slog"
	"time"

	"github.com/gosthome/gosthome/components/sensor"
	"github.com/gosthome/gosthome/components/textsensor"
	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/util"
	psutilHost "github.com/shirou/gopsutil/v4/host"
)

type HostConfig struct {
	component.PollingComponentConfig
	Hostname             bool `yaml:"hostname"`
	Uptime               bool `yaml:"uptime"`
	BootTime             bool `yaml:"boot_time"`
	Procs                bool `yaml:"procs"`
	OS                   bool `yaml:"os"`
	Platform             bool `yaml:"platform"`
	PlatformFamily       bool `yaml:"platform_family"`
	PlatformVersion      bool `yaml:"platform_version"`
	KernelVersion        bool `yaml:"kernel_version"`
	KernelArch           bool `yaml:"kernel_arch"`
	VirtualizationSystem bool `yaml:"virtualization_system"`
	VirtualizationRole   bool `yaml:"virtualization_role"`
	HostID               bool `yaml:"host_id"`
}

func NewHostConfig() HostConfig {
	return HostConfig{
		PollingComponentConfig: component.PollingComponentConfig{
			UpdateInterval: 10 * time.Second,
		},
		Hostname:             true,
		Uptime:               true,
		BootTime:             true,
		Procs:                true,
		OS:                   true,
		Platform:             true,
		PlatformFamily:       true,
		PlatformVersion:      true,
		KernelVersion:        true,
		KernelArch:           true,
		VirtualizationSystem: true,
		VirtualizationRole:   true,
		HostID:               true,
	}
}

// Validate implements validation.Validatable.
func (c *HostConfig) ValidateWithContext(ctx context.Context) error {
	return nil
}

// AutoLoad implements component.AutoLoader.
func (c *HostConfig) AutoLoad() []string {
	return []string{
		sensor.COMPONENT_KEY,
		textsensor.COMPONENT_KEY,
	}
}

var _ component.Config = (*Config)(nil)
var _ component.AutoLoader = (*Config)(nil)

type Host struct {
	cid.CID
	*component.PollingComponent[Host, *Host]
	component.WithInitializationPriorityProcessor
	ctx context.Context
	cfg *HostConfig

	sensors     map[string]*Sensor
	textSensors map[string]*TextSensor
}

func NewHost(ctx context.Context, cfg *HostConfig) (ret *Host, err error) {
	ret = &Host{
		CID:         cid.NewID("host"),
		ctx:         ctx,
		cfg:         cfg,
		sensors:     make(map[string]*Sensor),
		textSensors: make(map[string]*TextSensor),
	}
	ret.PollingComponent, err = component.NewPollingComponent[Host, *Host](ctx, ret, &cfg.PollingComponentConfig)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func str(s string) func() (string, uint64, bool) {
	return func() (string, uint64, bool) { return s, 0, true }
}

func uint(u uint64) func() (string, uint64, bool) {
	return func() (string, uint64, bool) { return "", u, false }
}

func timestamp(ts uint64) func() (string, uint64, bool) {
	return func() (string, uint64, bool) { return time.Unix(int64(ts), 0).Format(time.RFC3339), 0, true }
}

// Setup implements component.Component.
func (host *Host) Setup() {
	host.Poll()
	host.PollingComponent.Setup()
}

func (host *Host) getInfo(node *core.Node) {
	type hostSensor struct {
		enabled bool
		poll    bool
		name    string
		val     func() (string, uint64, bool)
	}
	info, err := psutilHost.Info()
	if err != nil {
		slog.Error("host info", "err", err)
	}
	for _, hs := range []hostSensor{
		{enabled: host.cfg.Hostname, name: "host_hostname", val: str(info.Hostname)},
		{enabled: host.cfg.Uptime, name: "host_uptime", val: timestamp(info.Uptime)},
		{enabled: host.cfg.BootTime, name: "host_boot_time", val: timestamp(info.BootTime)},
		{enabled: host.cfg.Procs, name: "host_procs", val: uint(info.Procs)},
		{enabled: host.cfg.OS, name: "host_os", val: str(info.OS)},
		{enabled: host.cfg.Platform, name: "host_platform", val: str(info.Platform)},
		{enabled: host.cfg.PlatformFamily, name: "host_platform_family", val: str(info.PlatformFamily)},
		{enabled: host.cfg.PlatformVersion, name: "host_platform_version", val: str(info.PlatformVersion)},
		{enabled: host.cfg.KernelVersion, name: "host_kernel_version", val: str(info.KernelVersion)},
		{enabled: host.cfg.KernelArch, name: "host_kernel_arch", val: str(info.KernelArch)},
		{enabled: host.cfg.VirtualizationSystem, name: "host_virtualization_system", val: str(info.VirtualizationSystem)},
		{enabled: host.cfg.VirtualizationRole, name: "host_virtualization_role", val: str(info.VirtualizationRole)},
		{enabled: host.cfg.HostID, name: "host_host_id", val: str(info.HostID)},
	} {
		if !hs.enabled {
			continue
		}
		sval, ival, isString := hs.val()
		if isString {
			ns, ok := host.textSensors[hs.name]
			if !ok {
				cfg := util.Modify(TextSensorConfig{}, func(c *TextSensorConfig) {
					c.Name = hs.name
				})
				ns = &TextSensor{}
				ns.BaseTextSensor, err = textsensor.NewBaseTextSensor(host.ctx, ns, &cfg.BaseTextSensorConfig)
				if err != nil {
					slog.Error("failed to create text sensor", "name", hs.name, "err", err)
					continue
				}
				node.RegisterTextSensor(ns)
				host.textSensors[hs.name] = ns
			}
			ns.SetState(entity.TextSensorState{
				State:        sval,
				MissingState: false,
			})
		} else {
			cfg := util.Modify(SensorConfig{}, func(c *SensorConfig) {
				c.Name = hs.name
			})
			ns, ok := host.sensors[hs.name]
			if !ok {
				ns = &Sensor{}
				ns.BaseSensor, err = sensor.NewBaseSensor(host.ctx, ns, &cfg.BaseSensorConfig)
				if err != nil {
					slog.Error("failed to create text sensor", "name", hs.name, "err", err)
					continue
				}
				node.RegisterSensor(ns)
				host.sensors[hs.name] = ns
			}
			ns.SetState(entity.SensorState{
				State:        float32(ival),
				MissingState: false,
			})
		}
	}
}

func (host *Host) Poll() {
	node := core.GetNode(host.ctx)
	host.getInfo(node)
}

// Close implements component.Component.
func (host *Host) Close() error {
	return host.PollingComponent.Close()
}
