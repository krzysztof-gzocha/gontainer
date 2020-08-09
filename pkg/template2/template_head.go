package template2

const TemplateHead = `package {{.Input.Meta.Pkg}}

{{range $import := .Imports.GetImports -}}
import {{$import.Alias}} "{{$import.Path}}"
{{end}}
`
