package compiler

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/regex"
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
	RegisterFunction(goImport string, goFunc string, tokenFunc string)
}

type ArgResolver interface {
	Resolve(interface{}) (compiled.Arg, error)
}

type Compiler struct {
	inputValidator    InputValidator
	compiledValidator CompiledValidator
	imports           Imports
	tokenizer         Tokenizer
	paramResolver     parameters.Resolver
	argResolver       ArgResolver
}

func NewCompiler(
	inputValidator InputValidator,
	compiledValidator CompiledValidator,
	imports Imports,
	tokenizer Tokenizer,
	paramResolver parameters.Resolver,
	argResolver ArgResolver,
) *Compiler {
	return &Compiler{
		inputValidator:    inputValidator,
		compiledValidator: compiledValidator,
		imports:           imports,
		tokenizer:         tokenizer,
		paramResolver:     paramResolver,
		argResolver:       argResolver,
	}
}

type compilerError struct {
	error
}

func throwCompilerError(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			err = fmt.Errorf("%s: %s", msg, err.Error())
		}
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
		throwCompilerError(
			c.imports.RegisterPrefix(a, sanitizeImport(p)),
			"cannot register alias",
		)
	}
}

var (
	regexMetaGoFn = regexp.MustCompile("^" + regex.MetaGoFn + "$")
)

func (c Compiler) handleMetaFuncs(funcs map[string]string) {
	for fn, goFn := range funcs {
		_, m := regex.Match(regexMetaGoFn, goFn)
		c.tokenizer.RegisterFunction(sanitizeImport(m["import"]), m["fn"], fn)
	}
}

func (c Compiler) handleParams(i input.DTO, result *compiled.DTO) {
	for n, v := range i.Params {
		param, err := c.paramResolver.Resolve(v)
		if err != nil {
			throwCompilerError(
				err,
				fmt.Sprintf("cannot resolve param `%s`", n),
			)
		}
		result.Params = append(
			result.Params,
			compiled.Param{
				Name:      n,
				Code:      param.Code,
				Raw:       param.Raw,
				DependsOn: param.DependsOn,
			},
		)
	}

	sort.SliceStable(result.Params, func(i, j int) bool {
		return result.Params[i].Name < result.Params[j].Name
	})
}

func (c Compiler) handleServices(i input.DTO, result *compiled.DTO) {
	for n, s := range i.Services {
		result.Services = append(result.Services, c.handleService(n, s))
	}
	sort.SliceStable(
		result.Services,
		func(i, j int) bool {
			return result.Services[i].Name < result.Services[j].Name
		},
	)
}

var (
	regexServiceType        = regexp.MustCompile("^" + regex.ServiceType + "$")
	regexServiceValue       = regexp.MustCompile("^" + regex.ServiceValue + "$")
	regexServiceConstructor = regexp.MustCompile("^" + regex.ServiceConstructor + "$")
)

func (c Compiler) handleService(name string, s input.Service) compiled.Service {
	if s.Todo {
		return compiled.Service{
			Name: name,
			Todo: true,
		}
	}

	r := compiled.Service{
		Name:        name,
		Getter:      s.Getter,
		Type:        c.handleServiceType(s.Type),
		Value:       c.handleServiceValue(s.Value),
		Constructor: c.handleServiceConstructor(s.Constructor),
		Args:        c.handleServiceArgs(fmt.Sprintf("service `%s`"), s.Args),
		Calls:       c.handleServiceCalls(name, s.Calls),
		Fields:      c.handleServiceFields(name, s.Fields),
		Tags:        s.Tags,
		Todo:        false,
	}

	return r
}

func (c Compiler) handleServiceType(serviceType string) string {
	_, m := regex.Match(regexServiceType, serviceType)
	t := m["type"]
	if m["import"] != "" {
		t = c.imports.GetAlias(sanitizeImport(m["import"])) + "." + t
	}
	return m["ptr"] + t
}

func (c Compiler) handleServiceValue(serviceValue string) string {
	_, m := regex.Match(regexServiceValue, serviceValue)
	parts := make([]string, 0)
	if m["import"] != "" {
		parts = append(parts, c.imports.GetAlias(sanitizeImport(m["import"])))
	}
	if m["struct"] != "" {
		parts = append(parts, m["struct"]+"{}")
	}
	return strings.Join(append(parts, m["value"]), ".")
}

func (c Compiler) handleServiceConstructor(serviceConstructor string) string {
	_, m := regex.Match(regexServiceConstructor, serviceConstructor)
	r := ""
	if m["import"] != "" {
		r = c.imports.GetAlias(sanitizeImport(m["import"])) + "."
	}
	return r + m["fn"]
}

func (c Compiler) handleServiceArgs(errorPrefix string, args []interface{}) (res []compiled.Arg) {
	for i, a := range args {
		arg, err := c.argResolver.Resolve(a)
		throwCompilerError(
			err,
			fmt.Sprintf("%s: cannot solve arg%d", errorPrefix, i),
		)
		res = append(res, arg)
	}
	return
}

func (c Compiler) handleServiceCalls(serviceName string, calls []input.Call) (res []compiled.Call) {
	for _, raw := range calls {
		call := compiled.Call{
			Method: raw.Method,
			Args: c.handleServiceArgs(
				fmt.Sprintf("service: `%s`: call `%s`", serviceName, raw.Method),
				raw.Args,
			),
			Immutable: raw.Immutable,
		}
		res = append(res, call)
	}
	return
}

func (c Compiler) handleServiceFields(serviceName string, fields map[string]interface{}) (res []compiled.Field) {
	for n, f := range fields {
		arg, err := c.argResolver.Resolve(f)
		throwCompilerError(
			err,
			fmt.Sprintf("service `%s`: field `%s`", serviceName, n),
		)
		field := compiled.Field{
			Name:  n,
			Value: arg,
		}
		res = append(res, field)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return
}
