package parameters

import (
	"strings"
)

type dependenciesBag map[string]bool

func newDependenciesBag() dependenciesBag {
	return make(dependenciesBag)
}

func (b dependenciesBag) Has(id string) bool {
	_, ok := b[id]
	return ok
}

func (b dependenciesBag) Append(id string) {
	b[id] = true
}

func (b dependenciesBag) Clone() dependenciesBag {
	r := newDependenciesBag()
	for id, _ := range b {
		r.Append(id)
	}

	return r
}

func (b dependenciesBag) ToString() string {
	names := make([]string, 0)
	for id, _ := range b {
		names = append(names, id)
	}

	return strings.Join(names, ", ")
}
