package psutil

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gosthome/gosthome/components/sensor"
	"github.com/gosthome/gosthome/components/textsensor"
	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
	"github.com/gosthome/gosthome/core/entity"
	"github.com/gosthome/gosthome/core/util"
	pstilCPU "github.com/shirou/gopsutil/v4/cpu"
)

type CPUConfig struct {
	component.PollingComponentConfig
	Count   CPUCountConfig       `yaml:"count"`
	Info    CPUInfoConfig        `yaml:"info"`
	Times   CPUPerformanceConfig `yaml:"times"`
	Percent CPUPerformanceConfig `yaml:"percent"`
}

type CPUCountConfig struct {
	Enabled bool `yaml:"enabled"`
	SensorConfig
	IncludeLogical bool `yaml:"include_logical"`
}

type CPUInfoConfig struct {
	Enabled bool `yaml:"enabled"`
}

type CPUPerformanceConfig struct {
	Enabled bool `yaml:"enabled"`
	Total   bool `yaml:"total"`
	PerCpu  bool `yaml:"per_cpu"`
}

func NewCPUConfig() CPUConfig {
	return CPUConfig{
		PollingComponentConfig: component.PollingComponentConfig{
			UpdateInterval: 10 * time.Second,
		},
		Count: util.Modify(CPUCountConfig{}, func(c *CPUCountConfig) {
			c.Name = "CPU count"
			c.Icon = "mdi:chip"
		}),
		Info: CPUInfoConfig{
			Enabled: true,
		},
		Times: CPUPerformanceConfig{
			Enabled: true,
			Total:   true,
		},
		Percent: CPUPerformanceConfig{
			Enabled: true,
			Total:   true,
		},
	}
}

// Validate implements validation.Validatable.
func (c *CPUConfig) ValidateWithContext(ctx context.Context) error {
	return nil
}

// AutoLoad implements component.AutoLoader.
func (c *CPUConfig) AutoLoad() component.Dependencies {
	return component.Depends(
		sensor.COMPONENT_KEY,
	)
}

var _ component.Config = (*Config)(nil)
var _ component.AutoLoader = (*Config)(nil)

type CPU struct {
	cid.CID
	*component.PollingComponent[CPU, *CPU]
	component.WithInitializationPriorityProcessor
	ctx context.Context
	cfg *CPUConfig

	times    map[string]map[string]*Sensor
	percents map[string]*Sensor
}

func NewCPU(ctx context.Context, cfg *CPUConfig) (ret *CPU, err error) {
	ret = &CPU{
		CID: cid.NewID("cpu"),
		ctx: ctx,
		cfg: cfg,
	}
	ret.PollingComponent, err = component.NewPollingComponent[CPU, *CPU](ctx, ret, &cfg.PollingComponentConfig)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Setup implements component.Component.
func (cpu *CPU) Setup() {
	cpu.times = make(map[string]map[string]*Sensor)
	cpu.percents = make(map[string]*Sensor)
	node := core.GetNode(cpu.ctx)
	cpus, err := pstilCPU.CountsWithContext(cpu.ctx, cpu.cfg.Count.IncludeLogical)
	if err != nil {
		slog.Error("failed to initialize", "err", err)
		return
	}
	counts := &Sensor{}
	counts.BaseSensor, err = sensor.NewBaseSensor(cpu.ctx, counts, &cpu.cfg.Count.BaseSensorConfig)
	counts.SetState(entity.SensorState{
		State:        float32(cpus),
		MissingState: false,
	})
	node.RegisterSensor(counts)

	if cpu.cfg.Info.Enabled {
		cpu.setupInfo(node)
	}

	if cpu.cfg.Times.Enabled {
		if cpu.cfg.Times.Total {
			cpu.getTimes(node, false)
		}
		if cpu.cfg.Times.PerCpu {
			cpu.getTimes(node, true)
		}
	}

	if cpu.cfg.Percent.Enabled {
		if cpu.cfg.Percent.Total {
			cpu.getPercent(node, false)
		}
		if cpu.cfg.Percent.PerCpu {
			cpu.getPercent(node, true)
		}
	}
	cpu.PollingComponent.Setup()
}

func (cpu *CPU) setupInfo(node *core.Node) {
	cpuInfos, err := pstilCPU.Info()
	if err != nil {
		slog.Error("failed to get cpu info", "err", err)
		return
	}
	for _, cpuInfo := range cpuInfos {
		for id, val := range map[string]string{
			"vendor_id":   cpuInfo.VendorID,
			"family":      cpuInfo.Family,
			"model":       cpuInfo.Model,
			"physical_id": cpuInfo.PhysicalID,
			"core_id":     cpuInfo.CoreID,
			"model_name":  cpuInfo.ModelName,
			"flags":       strings.Join(cpuInfo.Flags, ","),
			"microcode":   cpuInfo.Microcode,
		} {
			cfg := util.Modify(TextSensorConfig{}, func(c *TextSensorConfig) {
				c.Name = fmt.Sprintf("CPU %d %s", cpuInfo.CPU, id)
			})
			ns := &TextSensor{}
			ns.BaseTextSensor, err = textsensor.NewBaseTextSensor(cpu.ctx, ns, &cfg.BaseTextSensorConfig)
			ns.SetState(entity.TextSensorState{
				State:        val,
				MissingState: false,
			})
			node.RegisterTextSensor(ns)
		}
		for id, val := range map[string]float32{
			"stepping":   float32(cpuInfo.Stepping),
			"cores":      float32(cpuInfo.Cores),
			"mhz":        float32(cpuInfo.Mhz),
			"cache_size": float32(cpuInfo.CacheSize),
		} {
			cfg := util.Modify(SensorConfig{}, func(c *SensorConfig) {
				c.Name = fmt.Sprintf("CPU %d %s", cpuInfo.CPU, id)
			})
			ns := &Sensor{}
			ns.BaseSensor, err = sensor.NewBaseSensor(cpu.ctx, ns, &cfg.BaseSensorConfig)
			ns.SetState(entity.SensorState{
				State:        float32(val),
				MissingState: false,
			})
			node.RegisterSensor(ns)
		}
	}
}

func (cpu *CPU) getTimes(node *core.Node, perCPU bool) {
	cpuInfos, err := pstilCPU.Times(perCPU)
	if err != nil {
		slog.Error("failed to get cpu info", "err", err)
		return
	}
	for _, ctimes := range cpuInfos {
		sns, ok := cpu.times[ctimes.CPU]
		if !ok {
			sns = make(map[string]*Sensor)
		}

		for id, val := range map[string]float64{
			"user":       ctimes.User,
			"system":     ctimes.System,
			"idle":       ctimes.Idle,
			"nice":       ctimes.Nice,
			"iowait":     ctimes.Iowait,
			"irq":        ctimes.Irq,
			"softirq":    ctimes.Softirq,
			"steal":      ctimes.Steal,
			"guest":      ctimes.Guest,
			"guest_nice": ctimes.GuestNice,
		} {
			ns, ok := sns[id]
			if !ok {
				cfg := util.Modify(SensorConfig{}, func(c *SensorConfig) {
					c.Name = fmt.Sprintf("CPU %s time for %s", id, ctimes.CPU)
				})
				ns = &Sensor{}
				ns.BaseSensor, err = sensor.NewBaseSensor(cpu.ctx, ns, &cfg.BaseSensorConfig)
				node.RegisterSensor(ns)
				sns[id] = ns
			}
			ns.SetState(entity.SensorState{
				State:        float32(val),
				MissingState: false,
			})
		}
		cpu.times[ctimes.CPU] = sns
	}
}

func (cpu *CPU) getPercent(node *core.Node, perCPU bool) {
	percents, err := pstilCPU.Percent(0, perCPU)
	if err != nil {
		slog.Error("failed to get cpu info", "err", err)
		return
	}
	for i, percent := range percents {
		id := "total"
		if perCPU {
			id = fmt.Sprintf("cpu-%d", i)
		}
		cfg := util.Modify(SensorConfig{}, func(c *SensorConfig) {
			c.Name = fmt.Sprintf("Load %s", id)
		})
		ns, ok := cpu.percents[id]
		if !ok {
			ns = &Sensor{}
			ns.BaseSensor, err = sensor.NewBaseSensor(cpu.ctx, ns, &cfg.BaseSensorConfig)
			node.RegisterSensor(ns)
			cpu.percents[id] = ns
		}
		ns.SetState(entity.SensorState{
			State:        float32(percent),
			MissingState: false,
		})
	}
}

func (cpu *CPU) Poll() {
	node := core.GetNode(cpu.ctx)
	if cpu.cfg.Times.Enabled {
		if cpu.cfg.Times.Total {
			cpu.getTimes(node, false)
		}
		if cpu.cfg.Times.PerCpu {
			cpu.getTimes(node, true)
		}
	}

	if cpu.cfg.Percent.Enabled {
		if cpu.cfg.Percent.Total {
			cpu.getPercent(node, false)
		}
		if cpu.cfg.Percent.PerCpu {
			cpu.getPercent(node, true)
		}
	}
}

// Close implements component.Component.
func (cpu *CPU) Close() error {
	return cpu.PollingComponent.Close()
}
