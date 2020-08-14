package dto

import (
	"fmt"
	"sort"

	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/syntax"
)

type Compiler interface {
	Compile(Input) (CompiledInput, error)
}

// todo remove and use arguments.Resolver, now it's impossible, because of import cycle
type ArgResolver interface {
	Resolve(string) (CompiledArg, error)
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

	var names []string
	for n, _ := range i.Services {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, n := range names {
		s := i.Services[n]
		args := make([]CompiledArg, 0)
		for argN, a := range s.Args {
			arg, argErr := c.argumentResolver.Resolve(a)
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

		compiled := CompiledService{
			Name:        n,
			Getter:      s.Getter,
			Type:        "",
			Constructor: "",
			Disposable:  s.Disposable,
			Args:        args,
			Tags:        append(s.Args),
		}

		if s.Type != "" {
			type_, tErr := c.typeResolver.ResolveType(s.Type)
			if tErr != nil {
				return CompiledInput{}, tErr
			}
			compiled.Type = type_
		}

		if s.Constructor != "" {
			constructor, cErr := c.fnResolver.ResolveFunction(s.Constructor)
			if cErr != nil {
				return CompiledInput{}, cErr
			}
			compiled.Constructor = constructor
		}

		r.Services = append(r.Services, compiled)
	}

	if err := c.compiledValidator.Validate(r); err != nil {
		return CompiledInput{}, err
	}

	return r, nil
}
