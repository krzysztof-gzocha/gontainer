package dto

import (
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
)

type Compiler interface {
	Compile(Input) (CompiledInput, error)
}

type BaseCompiler struct {
	imports           imports.Imports
	validator         Validator
	compiledValidator CompiledValidator
	paramsResolver    parameters.BagFactory
}

func (c *BaseCompiler) Compile(i Input) (CompiledInput, error) {
	if err := c.validator.Validate(i); err != nil {
		return CompiledInput{}, err
	}

	r := CompiledInput{}
	r.Meta.Pkg = i.Meta.Pkg
	r.Meta.ContainerType = i.Meta.ContainerType

	params, paramsErr := c.paramsResolver.Create(i.Params)
	if paramsErr != nil {
		return CompiledInput{}, paramsErr
	}
	r.Params = params

	if err := c.compiledValidator.Validate(r); err != nil {
		return CompiledInput{}, err
	}

	return r, nil
}
