package template

import (
	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/definition"
)

type Service struct {
	definition.Service
	CompiledArgs []arguments.Argument
}

type Services map[string]Service
