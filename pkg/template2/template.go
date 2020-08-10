package template2

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gomponents/gontainer/pkg/dto"
	"github.com/gomponents/gontainer/pkg/exporters"
	"github.com/gomponents/gontainer/pkg/imports"
)

type Builder interface {
	Build(dto.CompiledInput) (string, error)
}

type SimpleBuilder struct {
	imports imports.Imports
}

func NewSimpleBuilder(imports imports.Imports) *SimpleBuilder {
	return &SimpleBuilder{imports: imports}
}

func (s SimpleBuilder) Build(i dto.CompiledInput) (string, error) {
	data := map[string]interface{}{
		"Imports": s.imports,
		"Input":   i,
	}

	fncs := template.FuncMap{
		"export": func(input interface{}) string {
			r, err := exporters.NewDefaultExporter().Export(input)
			if err != nil {
				panic(err)
			}
			return r
		},
		"importAlias": func(input interface{}) string {
			alias, ok := input.(string)
			if !ok {
				panic(fmt.Errorf("func `importAlias` expects `%T`, `%T` given", "", input))
			}

			return s.imports.GetAlias(alias)
		},
		"replace": func(input, from, to string) string {
			return strings.Replace(input, from, to, -1)
		},
	}

	exec := func(name string, tplBody string) (string, error) {
		tpl, newErr := template.New("gontainer_" + name).Funcs(fncs).Parse(tplBody)
		if newErr != nil {
			return "", newErr
		}
		var b bytes.Buffer
		tplErr := tpl.Execute(&b, data)
		return b.String(), tplErr
	}

	var (
		body, head string
		err        error
	)

	if body, err = exec("body", TemplateBody); err != nil {
		return "", err
	}

	if head, err = exec("head", TemplateHead); err != nil {
		return "", err
	}

	return head + body, nil
}
