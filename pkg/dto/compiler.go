// +build ignore

package dto

import (
	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/syntax"
)

type Compiler interface {
	Compile(Input) (CompiledInput, error)
}

type BaseCompiler struct {
	imports           imports.Imports
	validator         Validator
	compiledValidator CompiledValidator
	paramsResolver    parameters.BagFactory
	typeResolver      syntax.TypeResolver
	argumentResolver  arguments.Resolver
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

	for n, s := range i.Services {
		t, tErr := c.typeResolver.ResolveType(s.Type)
		if tErr != nil {
			return CompiledInput{}, tErr
		}

		// TODO resolve args

		r.Services[n] = CompiledService{
			Getter:      s.Getter,
			Type:        t,
			Constructor: s.Constructor,
			WithError:   s.WithError,
			Disposable:  s.Disposable,
			Args:        nil,
			Tags:        append(s.Args),
		}
	}

	if err := c.compiledValidator.Validate(r); err != nil {
		return CompiledInput{}, err
	}

	return r, nil
}
