package template

const TemplateHead = `package {{.Pkg}}

{{range $import := .Imports.GetImports -}}
import {{$import.Alias}} "{{$import.Path}}"
{{end}}
`
