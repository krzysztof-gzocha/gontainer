package arguments

import (
	"fmt"
)

type ArgumentKind uint

const (
	ArgumentKindService ArgumentKind = iota
	ArgumentKindCode
)

type ServiceMetadata struct {
	ID          string
	Import      string
	Type        string
	PointerType bool
}

type Argument struct {
	Kind            ArgumentKind
	Code            string
	ServiceMetadata ServiceMetadata
}

type Resolver interface {
	Resolve(string) (Argument, error)
}

type Subresolver interface {
	Resolver
	Supports(string) bool
}

type SimpleResolver struct {
	subresolvers []Subresolver
}

func NewSimpleResolver(subresolvers []Subresolver) *SimpleResolver {
	return &SimpleResolver{subresolvers: subresolvers}
}

func (s SimpleResolver) Resolve(expr string) (Argument, error) {
	for _, r := range s.subresolvers {
		if r.Supports(expr) {
			return r.Resolve(expr)
		}
	}

	return Argument{}, fmt.Errorf("cannot resolve argument `%d`", expr)
}
