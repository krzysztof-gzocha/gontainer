package compiler

import (
	"regexp"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexMetaGoFn = regexp.MustCompile("^" + regex.MetaGoFn + "$")
)

type Imports interface {
	GetAlias(string) string
	RegisterPrefix(shortcut string, path string) error
}

type InputValidator interface {
	Validate(input.DTO) error
}

type CompiledValidator interface {
	Validate(compiled.DTO) error
}

type Tokenizer interface {
	RegisterFunction(goImport string, goFunc string, tokenFun string)
}

type ParamBagFactory interface {
	Create(map[string]interface{}) (map[string]compiled.Param, error)
}

type Compiler struct {
	inputValidator    InputValidator
	compiledValidator CompiledValidator
	imports           Imports
	tokenizer         Tokenizer
	paramBagFactory   ParamBagFactory
}

type compilerError struct {
	error
}

func throwCompilerError(err error) {
	if err != nil {
		panic(compilerError{err})
	}
}

func (c Compiler) Compile(i input.DTO) (result compiled.DTO, err error) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		if cErr, ok := recovered.(compilerError); ok {
			result = compiled.DTO{}
			err = cErr
			return
		}

		panic(recovered)
	}()

	c.validateInput(i)
	c.handleMeta(i, &result)
	c.handleParams(i, &result)
	c.handleServices(i, &result)
	c.validateCompiled(result)

	return
}

func (c Compiler) validateInput(i input.DTO) {
	throwCompilerError(c.inputValidator.Validate(i))
}

func (c Compiler) validateCompiled(o compiled.DTO) {
	throwCompilerError(c.compiledValidator.Validate(o))
}

func (c Compiler) handleMeta(i input.DTO, result *compiled.DTO) {
	result.Meta.Pkg = i.Meta.Pkg
	result.Meta.ContainerType = i.Meta.ContainerType
	c.handleMetaImport(i.Meta.Imports)
	c.handleMetaFuncs(i.Meta.Functions)
}

func (c Compiler) handleMetaImport(imports map[string]string) {
	for a, p := range imports {
		throwCompilerError(c.imports.RegisterPrefix(a, sanitizeImport(p)))
	}
}

func (c Compiler) handleMetaFuncs(funcs map[string]string) {
	for fn, goFn := range funcs {
		_, m := regex.Match(regexMetaGoFn, goFn)
		c.tokenizer.RegisterFunction(sanitizeImport(m["import"]), m["fn"], fn)
	}
}

func (c Compiler) handleParams(i input.DTO, result *compiled.DTO) {
	var err error
	result.Params, err = c.paramBagFactory.Create(i.Params)
	throwCompilerError(err)
}

func (c Compiler) handleServices(i input.DTO, result *compiled.DTO) {

}
