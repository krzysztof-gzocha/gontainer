package arguments

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/dto"
	"github.com/gomponents/gontainer/pkg/parameters"
)

type Resolver interface {
	Resolve(string, parameters.ResolvedParams) (dto.CompiledArg, error)
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

func (s SimpleResolver) Resolve(expr string, params parameters.ResolvedParams) (dto.CompiledArg, error) {
	for _, r := range s.subresolvers {
		if r.Supports(expr) {
			return r.Resolve(expr, params)
		}
	}

	return dto.CompiledArg{}, fmt.Errorf("cannot resolve argument `%s`", expr)
}
