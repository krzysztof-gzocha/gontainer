package template

import (
	"github.com/gomponents/gontainer/pkg/dto"
)

type Service struct {
	dto.Service
	CompiledArgs []dto.CompiledArg
}

type Services map[string]Service
