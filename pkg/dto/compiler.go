package dto

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/syntax"
)

type Compiler interface {
	Compile(Input) (CompiledInput, error)
}

// todo remove and use arguments.Resolver, now it's impossible, because of import cycle
type ArgResolver interface {
	Resolve(string, parameters.ResolvedParams) (CompiledArg, error)
}

type BaseCompiler struct {
	imports           imports.Imports
	validator         Validator
	compiledValidator CompiledValidator
	paramsResolver    parameters.BagFactory
	typeResolver      syntax.TypeResolver
	argumentResolver  ArgResolver
	fnResolver        syntax.FunctionResolver
}

func NewBaseCompiler(
	imports imports.Imports,
	validator Validator,
	compiledValidator CompiledValidator,
	paramsResolver parameters.BagFactory,
	typeResolver syntax.TypeResolver,
	argumentResolver ArgResolver,
	fnResolver syntax.FunctionResolver,
) *BaseCompiler {
	return &BaseCompiler{
		imports:           imports,
		validator:         validator,
		compiledValidator: compiledValidator,
		paramsResolver:    paramsResolver,
		typeResolver:      typeResolver,
		argumentResolver:  argumentResolver,
		fnResolver:        fnResolver,
	}
}

func (c *BaseCompiler) Compile(i Input) (CompiledInput, error) {
	if err := c.validator.Validate(i); err != nil {
		return CompiledInput{}, err
	}

	r := CompiledInput{}
	r.Meta.Pkg = i.Meta.Pkg
	r.Meta.ContainerType = i.Meta.ContainerType

	for short, path := range i.Meta.Imports {
		if err := c.imports.RegisterPrefix(short, path); err != nil {
			return CompiledInput{}, err
		}
	}

	params, paramsErr := c.paramsResolver.Create(i.Params)
	if paramsErr != nil {
		return CompiledInput{}, paramsErr
	}
	r.Params = params
	r.Services = make(map[string]CompiledService)

	for n, s := range i.Services {
		type_, tErr := c.typeResolver.ResolveType(s.Type)
		if tErr != nil {
			return CompiledInput{}, tErr
		}

		constructor, cErr := c.fnResolver.ResolveFunction(s.Constructor)
		if cErr != nil {
			return CompiledInput{}, cErr
		}

		args := make([]CompiledArg, 0)
		for argN, a := range s.Args {
			arg, argErr := c.argumentResolver.Resolve(a, r.Params)
			if argErr != nil {
				return CompiledInput{}, fmt.Errorf(
					"cannot build `%s` service, error during building arg%d: %s",
					n,
					argN,
					argErr.Error(),
				)
			}
			args = append(args, arg)
		}

		r.Services[n] = CompiledService{
			Getter:      s.Getter,
			Type:        type_,
			Constructor: constructor,
			WithError:   s.WithError,
			Disposable:  s.Disposable,
			Args:        args,
			Tags:        append(s.Args),
		}
	}

	if err := c.compiledValidator.Validate(r); err != nil {
		return CompiledInput{}, err
	}

	return r, nil
}
