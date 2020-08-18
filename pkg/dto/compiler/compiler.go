package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type Imports interface {
	GetAlias(string) string
	RegisterPrefix(shortcut string, path string) error
}

type Validator interface {
	Validate(input.DTO) error
}

type Compiler struct {
	imports Imports
}

type compilerError struct {
	error
}

func panicIfNeed(err error) {
	if err != nil {
		panic(compilerError{err})
	}
}

func (c Compiler) Compile(i input.DTO) (result compiled.DTO, err error) {
	defer func() {
		recoverErr := recover()
		if recoverErr == nil {
			return
		}

		if cErr, ok := recoverErr.(compilerError); ok {
			result = compiled.DTO{}
			err = cErr
			return
		}

		panic(recoverErr)
	}()

	c.handleMeta(i, &result)

	return
}

func (c Compiler) handleMeta(i input.DTO, result *compiled.DTO) {
	c.handleImport(i)
}

func (c Compiler) handleImport(input input.DTO) {
	for a, p := range input.Meta.Imports {
		panicIfNeed(c.imports.RegisterPrefix(a, p))
	}
}
