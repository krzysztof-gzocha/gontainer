package template

const TemplateBody = `
{{- range $name, $param := .Params -}}
// {{$name}}
// Raw: {{ export $param.Raw }}
// GO:  {{$param.Code}}
// -------------------
{{end}}
type {{.ContainerType}} {{.RootImportAlias}}.Container

func CreateParamContainer() {{.RootImportAlias}}.ParamContainer {
	params := make(map[string]interface{})
{{range $name, $param := .Params}}	params["{{$name}}"] = {{$param.Code}}
{{end}}
	return {{.RootImportAlias}}.NewBaseParamContainer(params)
}
{{ $RootImportAlias := .RootImportAlias -}}
{{- $Imports := .Imports }}
func CreateContainer() {{.ContainerType}} {
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

// {{ $ContainerType := .ContainerType -}}
// {{range $service := .Services}}
// {{- if ne $service.Service.Getter "" }}
// func (c *{{$ContainerType}}) {{ $service.Service.Getter }} {{ $service.Service.Type }} {
// }
// {{ end }}
// {{ end }}
`
