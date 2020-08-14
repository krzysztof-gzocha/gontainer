package arguments

import (
	"github.com/gomponents/gontainer/pkg/dto"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/syntax"
)

type ServiceResolver struct {
	imports         imports.Imports
	patternResolver syntax.ServiceResolver
}

func NewServiceResolver(imports imports.Imports, patternResolver syntax.ServiceResolver) *ServiceResolver {
	return &ServiceResolver{imports: imports, patternResolver: patternResolver}
}

func (s ServiceResolver) Resolve(expr string) (dto.CompiledArg, error) {
	service, type_, err := s.patternResolver.ResolveService(expr)
	if err != nil {
		return dto.CompiledArg{}, err
	}
	return dto.CompiledArg{
		ServiceLink: &dto.ServiceLink{
			Name: service,
			Type: type_,
		},
	}, nil
}

func (s ServiceResolver) Supports(expr string) bool {
	return len(expr) >= 1 && []rune(expr)[0] == '@'
}
