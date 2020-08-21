package parameters

import (
	"strings"
)

type dependenciesBag []string

func newDependenciesBag() dependenciesBag {
	return make(dependenciesBag, 0)
}

func (b *dependenciesBag) Has(id string) bool {
	for _, s := range *b {
		if s == id {
			return true
		}
	}
	return false
}

func (b *dependenciesBag) Append(id string) {
	tmp := append(*b, id)
	*b = tmp
}

func (b dependenciesBag) Clone() dependenciesBag {
	return append(b)
}

func (b dependenciesBag) String() string {
	return strings.Join(b, ", ")
}
