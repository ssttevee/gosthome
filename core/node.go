package core

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"

	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/config"
	"github.com/gosthome/gosthome/core/entity"
)

type Node struct {
	*entity.Registry
	*bus.Bus
	Config *config.Config

	cmp []component.Component
	ctx context.Context
}

type nodeCtxKey struct{}

func GetNode(ctx context.Context) *Node {
	v := ctx.Value(nodeCtxKey{})
	if v == nil {
		return nil
	}
	n, ok := v.(*Node)
	if !ok {
		return nil
	}
	return n
}

func createFromConfig(ctx context.Context, cfg *config.Config, name string, config component.Config) ([]component.Component, error) {
	decl, ok := cfg.Registry.Get(name)
	if !ok {
		panic("unregistered component in config " + name)
	}
	cmpA, err := decl.Component(ctx, config)
	if err != nil {
		return nil, err
	}
	return cmpA, nil
}

func NewNode(ctx context.Context, cfg *config.Config) (*Node, error) {
	requiredComponents := component.Depends()
	for _, componentConfig := range cfg.Components {
		al, ok := componentConfig.Config.(component.AutoLoader)
		if ok {
			requiredComponents.Join(al.AutoLoad())
		}
	}
	ret := &Node{
		Config:   cfg,
		Bus:      bus.New(),
		Registry: &entity.Registry{},
		cmp:      []component.Component{},
	}
	ctx = context.WithValue(ctx, nodeCtxKey{}, ret)
	ctx = bus.Context(ctx, ret.Bus)
	for c := range requiredComponents {
		componentConfig, ok := cfg.Components[c]
		if !ok {
			slog.Debug("required component is not in config, using default", "component", c)
			cd, ok := cfg.Get(c)
			if !ok {
				return nil, fmt.Errorf("required component %s is not registered!", c)
			}
			cfg.Components[c] = cd.Config()
			componentConfig = cfg.Components[c]
		}
		cmpA, err := createFromConfig(ctx, cfg, c, componentConfig.Config)
		if err != nil {
			return nil, err
		}
		ret.cmp = append(ret.cmp, cmpA...)
	}
	for name, componentConfig := range cfg.Components {
		if _, ok := requiredComponents[name]; ok {
			continue
		}
		cmpA, err := createFromConfig(ctx, cfg, name, componentConfig.Config)
		if err != nil {
			return nil, err
		}
		ret.cmp = append(ret.cmp, cmpA...)
	}
	return ret, nil
}

func (n *Node) Start() {
	for _, cmp := range slices.SortedFunc(slices.Values(n.cmp), func(l component.Component, r component.Component) int {
		return cmp.Compare(int(l.InitializationPriority()), int(r.InitializationPriority()))
	}) {
		slog.Info("Setting up component", "cmp", reflect.TypeOf(cmp).String())
		cmp.Setup()
		slog.Info("Done setting up", "cmp", reflect.TypeOf(cmp).String())
	}
}

type ComponentPredicate func(c component.Component) bool

func (n *Node) GetComponent(predicate ComponentPredicate) (component.Component, bool) {
	i := slices.IndexFunc(n.cmp, predicate)
	if i < 0 {
		return nil, false
	}
	return n.cmp[i], true
}

func (n *Node) Close() error {
	errs := []error{}
	for _, c := range n.cmp {
		err := c.Close()
		if err != nil {
			errs = append(errs, err)
			slog.Error("Failed to stop component", "err", err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (n *Node) DefaultUniqueId(type_ entity.DomainType, entity entity.Entity) string {
	return n.Config.Gosthome.Name + type_.String() + entity.ID()
}
