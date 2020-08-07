package template

import (
	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/dto"
)

type Service struct {
	dto.Service
	CompiledArgs []arguments.Argument
}

type Services map[string]Service
