package component

import (
	"context"
	"maps"
	"reflect"

	"github.com/gosthome/gosthome/core/component/cid"
	cv "github.com/gosthome/gosthome/core/configvalidation"
)

type Component interface {
	cid.Identifier
	Setup()
	InitializationPriority() InitializationPriority
	Close() error
}

type Config interface {
	ComponentType() reflect.Type
	cv.Validatable
}

type Dependencies map[string]struct{}

func Depends(ds ...string) Dependencies {
	d := make(Dependencies)
	for _, s := range ds {
		d.Add(s)
	}
	return d
}

func (d Dependencies) Join(others ...Dependencies) {
	for _, o := range others {
		maps.Copy(d, o)
	}
}

func (d Dependencies) Add(s string) {
	d[s] = struct{}{}
}

// AutoLoader is a component config showing need for another component during startup
// E.g. component is registering EntityComponents
type AutoLoader interface {
	AutoLoad() Dependencies
}

type ConfigOf[T any, PT interface {
	*T
	Component
}] struct{}

func (*ConfigOf[T, PT]) ComponentType() reflect.Type {
	return reflect.TypeOf(PT(nil))
}

type Declaration interface {
	Config() *ConfigDecoder
	Component(ctx context.Context, cfg Config) ([]Component, error)
}
