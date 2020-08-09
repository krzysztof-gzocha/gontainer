package template2

const TemplateBody = `
{{- $RootImportAlias := importAlias "github.com/gomponents/gontainer/pkg" -}}
{{- $ContainerType := .Input.Meta.ContainerType -}}

{{- range $name, $param := .Input.Params -}}
// {{$name}}
// Raw: {{ export $param.Raw }}
// GO:  {{$param.Code}}
// -----------------------------------------------------------------------------
{{end}}
type {{$ContainerType}} struct {
	container {{$RootImportAlias}}.Container
}

func (c *{{$ContainerType}}) Get(id string) (interface{}, error) {
	return c.container.Get(id)
}

func (c *{{$ContainerType}}) MustGet(id string) interface{} {
	return c.container.MustGet(id)
}

func (c *{{$ContainerType}}) Has(id string) bool {
	return c.container.Has(id)
}

func (c *{{$ContainerType}}) ValidateAllServices() (errors map[string]error) {
	errors = make(map[string]error)
	for _, id := range []string{
{{range $name, $service := .Services -}} {{ "		" }} {{- export $name }},{{ "\n" }}{{ end -}}
{{- "	" -}} } {
		if _, err := c.Get(id); err != nil {
			errors[id] = err
		}
	}
	if len(errors) == 0 {
		errors = nil
	}
	return
}

func CreateParamContainer() {{$RootImportAlias}}.ParamContainer {
	params := make(map[string]interface{})
{{range $name, $param := .Input.Params}}	params["{{$name}}"] = {{$param.Code}}
{{end}}
	return {{$RootImportAlias}}.NewBaseParamContainer(params)
}

func CreateContainer() *{{$ContainerType}} {
	result := &{{$ContainerType}}{}

	getters := make(map[string]{{$RootImportAlias}}.GetterDefinition)
{{range $name, $service := .Input.Services}}	getters[{{ export $name }}] = {{$RootImportAlias}}.GetterDefinition{
		Getter: func() (interface{}, error) {
{{- range $argIndex, $arg := $service.Args -}}
{{- if $arg.IsService }}
			arg{{ $argIndex }}, errGet{{ $argIndex }} := result.Get({{ export $arg.ServiceLink.Name }})
			if errGet{{ $argIndex }} != nil {
				return nil, {{ importAlias "fmt" }}.Errorf("cannot create %s due to: %s", {{ export $name }}, errGet{{ $argIndex }}.Error())
			}
{{- if ne $arg.ServiceLink.Name "" }}
			val{{ $argIndex }}, ok{{ $argIndex }} := arg{{ $argIndex }}.({{ $arg.ServiceLink.Type }})
			if !ok{{ $argIndex }} {
				return nil, {{ importAlias "fmt" }}.Errorf("service %s is not an instance of %T, %T given", {{ export $name }}, val{{ $argIndex }}, arg{{ $argIndex }})
			}
{{ else }}
// todo
			val{{ $argIndex }}, := arg{{ $argIndex }}
{{- end -}}
{{- end -}}
{{ end }}
			return {{$service.Constructor}}(
				{{- range $argIndex, $arg := $service.Args -}}
					{{- if $arg.IsService }}
				val{{ $argIndex }},
					{{- else }}
				{{ $arg.Code }},
					{{- end -}}
				{{- end }}
			){{if not $service.WithError}}, nil{{end}}
		},
		Disposable: {{$service.Disposable}},
	}
{{end}}
	result.container = {{$RootImportAlias}}.NewBaseContainer(getters)
	return result
}

{{- range $name, $service := .Input.Services -}}
{{- if ne $service.Getter "" }}

func (c *{{$ContainerType}}) {{ $service.Getter }}() (result {{ $service.Type }}, err error) {
	var object interface{}
	var ok bool

	object, err = c.Get({{ export $name }})

	if err != nil {
		return
	}

	if result, ok = object.({{ $service.Type }}); !ok {
		err = {{ importAlias "fmt" }}.Errorf("cannot create %s, because cannot cast %T to %T", {{ export $name }}, object, result)
	}

	return
}
{{- end -}}
{{- end }}
`
