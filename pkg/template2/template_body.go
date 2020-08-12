package template2

const TemplateBody = `
{{- $ContainerPathAlias := importAlias "github.com/gomponents/gontainer-helpers/container" -}}
{{- $ContainerType := .Input.Meta.ContainerType -}}

{{- range $name, $param := .Input.Params -}}
// {{$name}}
// Raw: {{ export $param.Raw }}
// GO:  {{$param.Code}}
// -----------------------------------------------------------------------------
{{end}}
type {{$ContainerType}} struct {
	container {{$ContainerPathAlias}}.Container
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
{{range $service := .Input.Services -}} {{ "		" }} {{- export $service.Name }},{{ "\n" }}{{ end -}}
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

func Hey(provider interface{}, params ...interface{}) (interface{}, error) {
	return {{ importAlias "github.com/gomponents/gontainer-helpers/caller" }}.CallProvider(provider, params...)
}

func CreateParamContainer() {{$ContainerPathAlias}}.ParamContainer {
	params := make(map[string]interface{})
{{range $name, $param := .Input.Params}}	params["{{$name}}"] = {{$param.Code}}
{{end}}
	return {{$ContainerPathAlias}}.NewBaseParamContainer(params)
}

func CreateContainer() *{{$ContainerType}} {
	result := &{{$ContainerType}}{}

	getters := make(map[string]{{$ContainerPathAlias}}.GetterDefinition)
	getters["serviceContainer"] = {{$ContainerPathAlias}}.GetterDefinition{
		Getter: func() (interface{}, error) {
			return result, nil
		},
		Disposable: false,
	}
{{range $service := .Input.Services}}	getters[{{ export $service.Name }}] = {{$ContainerPathAlias}}.GetterDefinition{
		Getter: func() (service interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					service = nil
					err = {{ importAlias "fmt" }}.Errorf("%s", r)
				}
			}()
{{- range $argIndex, $arg := $service.Args -}}
{{- if $arg.IsService }}
			arg{{ $argIndex }}, errGet{{ $argIndex }} := result.Get({{ export $arg.ServiceLink.Name }})
			if errGet{{ $argIndex }} != nil {
				return nil, {{ importAlias "fmt" }}.Errorf("cannot create %s due to: %s", {{ export $service.Name }}, errGet{{ $argIndex }}.Error())
			}
{{- if ne $arg.ServiceLink.Type "" }}
			val{{ $argIndex }}, ok{{ $argIndex }} := arg{{ $argIndex }}.({{ $arg.ServiceLink.Type }})
			if !ok{{ $argIndex }} {
				return nil, {{ importAlias "fmt" }}.Errorf("service %s is not an instance of %T, %T given", {{ export $service.Name }}, val{{ $argIndex }}, arg{{ $argIndex }})
			}
{{ else }}
			val{{ $argIndex }} := arg{{ $argIndex }}
{{- end -}}
{{- end -}}
{{ end }}
			{{ if ne $service.Constructor "" }}return {{$service.Constructor}}(
				{{- range $argIndex, $arg := $service.Args }}
				// {{ export $arg.Raw }}
					{{- if $arg.IsService }}
				val{{ $argIndex }},
					{{- else }}
				{{ $arg.Code }},
					{{- end -}}
				{{- end }}
			){{if not $service.WithError}}, nil{{end}}{{ else }}return {{ replace $service.Type "*" "&" }}{}, nil{{ end }}
		},
		Disposable: {{$service.Disposable}},
	}
{{end}}
	result.container = {{$ContainerPathAlias}}.NewBaseContainer(getters)
	return result
}

{{- range $service := .Input.Services -}}
{{- if ne $service.Getter "" }}

func (c *{{$ContainerType}}) {{ $service.Getter }}() (result {{ $service.Type }}, err error) {
	var object interface{}
	var ok bool

	object, err = c.Get({{ export $service.Name }})

	if err != nil {
		return
	}

	if result, ok = object.({{ $service.Type }}); !ok {
		err = {{ importAlias "fmt" }}.Errorf("cannot create %s, because cannot cast %T to %T", {{ export $service.Name }}, object, result)
	}

	return
}
{{- end -}}
{{- end }}
`
