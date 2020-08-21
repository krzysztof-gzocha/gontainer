package arguments

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
)

type ServiceResolver struct {
}

func NewServiceResolver() *ServiceResolver {
	return &ServiceResolver{}
}

func (s ServiceResolver) Resolve(v interface{}) (compiled.Arg, error) {
	expr, _ := v.(string)
	service := expr[1:]
	return compiled.Arg{
		Code:              fmt.Sprintf("result.MustGet(%+q)", service),
		Raw:               expr,
		DependsOnServices: []string{service},
	}, nil
}

func (s ServiceResolver) Supports(v interface{}) bool {
	expr, ok := v.(string)
	if !ok {
		return false
	}
	return len(expr) > 1 && []rune(expr)[0] == '@'
}
