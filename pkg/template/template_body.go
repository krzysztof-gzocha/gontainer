package template

const TemplateBody = `
{{- range $name, $param := .Params -}}
// {{$name}}
// Raw: {{ export $param.Raw }}
// GO:  {{$param.Code}}
// -----------------------------------------------------------------------------
{{end}}
type {{.ContainerType}} struct {
	container {{.RootImportAlias}}.Container
}

func (c *{{.ContainerType}}) Get(id string) (interface{}, error) {
	return c.container.Get(id)
}

func (c *{{.ContainerType}}) MustGet(id string) interface{} {
	return c.container.MustGet(id)
}

func (c *{{.ContainerType}}) Has(id string) bool {
	return c.container.Has(id)
}

func CreateParamContainer() {{.RootImportAlias}}.ParamContainer {
	params := make(map[string]interface{})
{{range $name, $param := .Params}}	params["{{$name}}"] = {{$param.Code}}
{{end}}
	return {{.RootImportAlias}}.NewBaseParamContainer(params)
}

{{ $RootImportAlias := .RootImportAlias -}}
func CreateContainer() *{{.ContainerType}} {
	result := &{{.ContainerType}}{}

	getters := make(map[string]{{.RootImportAlias}}.GetterDefinition)
{{range $name, $service := .Services}}	getters[{{ export $name }}] = {{$RootImportAlias}}.GetterDefinition{
		Getter: func() (interface{}, error) {
{{- range $argIndex, $arg := $service.Service.CompiledArgs -}}
{{- if $arg.IsService }}
			arg{{ $argIndex }}, errGet{{ $argIndex }} := result.Get({{ export $arg.ServiceLink.Name }})
			if errGet{{ $argIndex }} != nil {
				return nil, {{ importAlias "fmt" }}.Errorf("cannot create %s due to: %s", {{ export $name }}, errGet{{ $argIndex }}.Error())
			}
{{- if ne $arg.ServiceLink.Name "" }}
			val{{ $argIndex }}, ok{{ $argIndex }} := arg{{ $argIndex }}.({{ $arg.ServiceLink.Type }})
			if !ok{{ $argIndex }} {
				return nil, {{ importAlias "fmt" }}.Errorf("service %s is not an instance of %s, %T given", "todo", {{ export $arg.ServiceLink.Type }}, arg{{ $argIndex }})
			}
{{ else }}
// todo
			val{{ $argIndex }}, := arg{{ $argIndex }}
{{- end -}}
{{- end -}}
{{ end }}
			return {{importAlias $service.Import}}.{{$service.Function}}(
				{{- range $argIndex, $arg := $service.Service.CompiledArgs -}}
					{{- if $arg.IsService }}
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
	result.container = {{.RootImportAlias}}.NewBaseContainer(getters)
	return result
}

{{- $ContainerType := .ContainerType -}}
{{- range $name, $service := .Services -}}
{{- if ne $service.Service.Getter "" }}

func (c *{{$ContainerType}}) {{ $service.Service.Getter }}() (result {{ $service.Service.Type }}, err error) {
	var object interface{}
	var ok bool

	object, err = c.Get({{ export $name }})

	if err != nil {
		return
	}

	if result, ok = object.({{ $service.Service.Type }}); !ok {
		err = {{ importAlias "fmt" }}.Errorf("cannot create %s, because cannot cast %T to %T", {{ export $name }}, object, result)
	}

	return
}
{{- end -}}
{{- end }}
`
