package template

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gomponents/gontainer/pkg/exporters"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
)

const rootPkg = "github.com/gomponents/gontainer/pkg"

const TemplateHead = `package {{.Pkg}}

{{range $import := .Imports.GetImports -}}
import {{$import.Alias}} "{{$import.Path}}"
{{end}}
`

const TemplateBody = `
{{- range $name, $param := .Params -}}
// {{$name}}
// Raw: {{ export $param.Raw }}
// GO:  {{$param.Code}}
// -------------------
{{end}}
func CreateParamContainer() {{.RootImportAlias}}.ParamContainer {
	params := make(map[string]interface{})
{{range $name, $param := .Params}}	params["{{$name}}"] = {{$param.Code}}
{{end}}
	return {{.RootImportAlias}}.NewBaseParamContainer(params)
}
{{ $RootImportAlias := .RootImportAlias -}}
{{- $Imports := .Imports }}
func CreateContainer() {{.RootImportAlias}}.Container {
	var result *{{.RootImportAlias}}.BaseContainer

	getters := make(map[string]{{.RootImportAlias}}.GetterDefinition)
{{range $name, $service := .Services}}	getters[{{ export $name }}] = {{$RootImportAlias}}.GetterDefinition{
		Getter: func() (interface{}, error) {
{{- range $argIndex, $arg := $service.Service.CompiledArgs -}}
{{- if eq $arg.Kind 0 }}
			arg{{ $argIndex }}, errGet{{ $argIndex }} := result.Get({{ export $arg.ServiceMetadata.ID }})
			if errGet{{ $argIndex }} != nil {
				return nil, {{ importAlias "fmt" }}.Errorf("cannot create %s due to: %s", {{ export $name }}, errGet{{ $argIndex }}.Error())
			}
{{- if ne $arg.ServiceMetadata.Import "" }}
			val{{ $argIndex }}, ok{{ $argIndex }} := arg{{ $argIndex }}.({{ if $arg.ServiceMetadata.PointerType }}*{{ end }}{{ importAlias $arg.ServiceMetadata.Import }}.{{ $arg.ServiceMetadata.Type }})
			if !ok{{ $argIndex }} {
				return nil, {{ importAlias "fmt" }}.Errorf("service %s is not an instance of %s, %T given", {{ export $arg.ServiceMetadata.ID }}, {{ export $arg.ServiceMetadata.Type }}, arg{{ $argIndex }})
			}
{{ else }}
			val{{ $argIndex }}, := arg{{ $argIndex }}
{{- end -}}
{{- end -}}
{{ end }}
			return {{importAlias $service.Import}}.{{$service.Function}}(
				{{- range $argIndex, $arg := $service.Service.CompiledArgs -}}
					{{- if eq $arg.Kind 0 }}
				val{{ $argIndex }},
					{{- else }}
				{{ $arg.Code }},
					{{- end -}}
				{{- end }}
			){{if not $service.Service.WithError}}, nil{{end}}
		},
		Disposable: {{$service.Service.Disposable}},
	}
{{end}}
	result = {{.RootImportAlias}}.NewBaseContainer(getters)
	return result
}
`

type Builder interface {
	Build() (string, error)
}

type serviceMetaDefinition struct {
	Service  Service
	Import   string
	Function string
}

// todo
// http://technosophos.com/2013/11/23/using-custom-template-functions-in-go.html
// ExecuteTemplate with tplVars
type simpleBuilderData struct {
	Pkg             string
	Imports         imports.Imports
	Params          parameters.ResolvedParams
	Services        map[string]serviceMetaDefinition
	RootImportAlias string
}

type SimpleBuilderOpt func(*SimpleBuilder) error

func WithPackage(pkg string) SimpleBuilderOpt {
	return func(builder *SimpleBuilder) error {
		builder.data.Pkg = pkg
		return nil
	}
}

func WithImports(imports imports.Imports) SimpleBuilderOpt {
	return func(builder *SimpleBuilder) error {
		builder.data.Imports = imports
		// todo refactor, use fn importAlias
		builder.data.RootImportAlias = imports.GetAlias(rootPkg)
		return nil
	}
}

func WithParams(params parameters.ResolvedParams) SimpleBuilderOpt {
	return func(builder *SimpleBuilder) error {
		builder.data.Params = params
		return nil
	}
}

func WithServices(services Services) SimpleBuilderOpt {
	return func(builder *SimpleBuilder) error {
		result := make(map[string]serviceMetaDefinition)
		for n, s := range services {
			parts := strings.Split(s.Constructor, ".")

			result[n] = serviceMetaDefinition{
				Service:  s,
				Import:   strings.Join(parts[:len(parts)-1], "."),
				Function: parts[len(parts)-1],
			}
		}
		builder.data.Services = result
		return nil
	}
}

type SimpleBuilder struct {
	data    simpleBuilderData
	tplBody *template.Template
	tplHead *template.Template
}

func (s SimpleBuilder) Build() (string, error) {
	exec := func(tpl *template.Template, data interface{}) (string, error) {
		var b bytes.Buffer
		err := tpl.Execute(&b, data)
		return b.String(), err
	}
	var (
		body, head string
		err        error
	)

	if body, err = exec(s.tplBody, s.data); err != nil {
		return "", err
	}

	if head, err = exec(s.tplHead, s.data); err != nil {
		return "", err
	}

	return head + body, nil
}

func NewSimpleBuilder(opts ...SimpleBuilderOpt) (*SimpleBuilder, error) {
	var r *SimpleBuilder

	fncs := template.FuncMap{
		"export": func(input interface{}) string {
			r, err := exporters.NewDefaultExporter().Export(input)
			if err != nil {
				panic(err)
			}
			return r
		},
		"importAlias": func(input interface{}) string {
			// todo handle err
			alias, _ := input.(string)

			return r.data.Imports.GetAlias(alias)
		},
	}

	r = &SimpleBuilder{}

	if tpl, parseErr := newTemplate("gointainer_body", fncs, TemplateBody); parseErr != nil {
		return nil, parseErr
	} else {
		r.tplBody = tpl
	}

	if tpl, parseErr := newTemplate("gontainer_imports", fncs, TemplateHead); parseErr != nil {
		return nil, parseErr
	} else {
		r.tplHead = tpl
	}

	for _, o := range opts {
		if err := o(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func newTemplate(name string, fncs map[string]interface{}, tpl string) (*template.Template, error) {
	return template.New(name).Funcs(fncs).Parse(tpl)
}
