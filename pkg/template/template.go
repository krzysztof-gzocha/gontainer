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
	ContainerType   string
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

func WithContainerType(t string) SimpleBuilderOpt {
	return func(builder *SimpleBuilder) error {
		builder.data.ContainerType = t
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
		// TODO remove given function
		// *gontainer/example/pkg.Person
		// registers import "gontainer/example/pkg"
		// returns *Person
		"decorateType": func(t string) string {
			ptr := []rune(t)[0] == '*'
			t = strings.TrimLeft(t, "*")

			parts := strings.Split(t, "/")
			result := parts[len(parts)-1]
			path := parts[:len(parts)-1]

			if strings.Contains(result, ".") {
				subparts := strings.Split(result, ".")
				result = subparts[len(subparts)-1]
				path = append(path, strings.Join(subparts[:len(subparts)-1], "."))
			}

			// TODO move imports somewhere else
			if len(path) > 0 {
				imp := strings.Join(path, "/")
				result = r.data.Imports.GetAlias(imp) + "." + result
			}

			if ptr {
				result = "*" + result
			}

			return result
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
