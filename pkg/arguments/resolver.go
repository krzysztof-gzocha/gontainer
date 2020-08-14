package arguments

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/dto"
)

type Resolver interface {
	Resolve(string) (dto.CompiledArg, error)
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

func (s SimpleResolver) Resolve(expr string) (dto.CompiledArg, error) {
	for _, r := range s.subresolvers {
		if r.Supports(expr) {
			result, err := r.Resolve(expr)
			if err == nil {
				result.Raw = expr
			}

			return result, err
		}
	}

	return dto.CompiledArg{}, fmt.Errorf("cannot resolve argument `%s`", expr)
}
