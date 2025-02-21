package entity

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
	"github.com/gosthome/gosthome/core/guarded"
)

func clone[T Entity](in []T) []Entity {
	ret := make([]Entity, 0, len(in))
	for _, e := range in {
		ret = append(ret, e)
	}
	return ret
}

func findByKey[T Entity](entityStore []T, hash uint32) (ret T, found bool) {
	i, found := slices.BinarySearchFunc(entityStore, hash, func(e T, t uint32) int {
		return cmp.Compare(e.HashID(), t)
	})
	if found {
		ret = entityStore[i]
	}
	return
}

func registerEntity[T Entity](entityStore []T, ent T) ([]T, error) {
	if ent.ID() == "" {
		return entityStore, fmt.Errorf("trying to register %T without an ID", ent)
	}
	i, found := slices.BinarySearchFunc(entityStore, ent, func(e T, t T) int {
		return cmp.Compare(e.HashID(), t.HashID())
	})
	if found {
		return entityStore, fmt.Errorf("hash id of %T is already registered %s!", ent, ent.ID())
	}
	return slices.Insert(entityStore, i, ent), nil
}

var domainHashes = func() (dh map[DomainType]uint32) {
	dh = make(map[DomainType]uint32)
	for dt, name := range _DomainTypeMap {
		dh[dt] = cid.HashID(name)
	}
	return
}()

type BaseDomain[Domain any, EntityType EntityComponent, PD interface {
	DomainTyper
	*Domain
}] struct {
	entities guarded.RWValue[[]EntityType]
}

func (bd *BaseDomain[Domain, EntityType, PD]) ID() string {
	return PD(nil).DomainType().String()
}
func (bd *BaseDomain[Domain, EntityType, PD]) HashID() uint32 {
	return domainHashes[PD(nil).DomainType()]
}

func (bd *BaseDomain[Domain, EntityType, PD]) Setup() {
}

func (bd *BaseDomain[Domain, EntityType, PD]) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityBus
}

func (bd *BaseDomain[Domain, EntityType, PD]) Close() error {
	bd.entities.Write(func(entities *[]EntityType) {
		*entities = nil
	})
	return nil
}

func (bd *BaseDomain[Domain, EntityType, PD]) Clone() (cloned []Entity) {
	bd.entities.Read(func(et *[]EntityType) {
		cloned = clone(*et)
	})
	return
}

func (bd *BaseDomain[Domain, EntityType, PD]) FindByKey(key uint32) (ret EntityType, found bool) {
	bd.entities.Read(func(et *[]EntityType) {
		ret, found = findByKey(*et, key)
	})
	return
}

func (bd *BaseDomain[Domain, EntityType, PD]) Register(ent EntityType) (err error) {
	bd.entities.Write(func(entities *[]EntityType) {
		*entities, err = registerEntity(*entities, ent)
	})
	return err
}

func (bd *BaseDomain[Domain, EntityType, PD]) DomainType() DomainType {
	return (PD)(nil).DomainType()
}
