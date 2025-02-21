package cid

import (
	"fmt"
	"hash/fnv"

	"github.com/gosthome/gosthome/core/guarded"
)

type Identifier interface {
	ID() string
	HashID() uint32
}

type CID interface {
	Identifier
	ensureID()
}

type cid struct {
	id   string
	hash uint32
}

func (cid) ensureID() {}

func NewID(id string) CID {
	return cid{
		id:   id,
		hash: HashID(id),
	}
}

var idgenMap = guarded.New(make(map[string]uint64))

func MakeStringID(prefix string) (ret string) {
	idgenMap.Do(func(idgen *map[string]uint64) {
		i, ok := (*idgen)[prefix]
		if !ok {
			i = 0
		}
		ret = prefix + "_" + fmt.Sprintf("%06x", i)
		(*idgen)[prefix] = i + 1
	})
	return
}

func MakeID(prefix string) CID {
	return NewID(MakeStringID(prefix))
}

// ID implements Entity.
func (i cid) ID() string {
	return i.id
}

func (i cid) SetID(s string) {
	i.id = s
	i.hash = HashID(i.id)
}

func HashID(id string) uint32 {
	h := fnv.New32()
	h.Write([]byte(id))
	return h.Sum32()
}

// HashID implements Entity.
func (i cid) HashID() uint32 {
	return i.hash
}
