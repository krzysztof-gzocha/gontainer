package syntax

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/regex"
)

type FunctionResolver interface {
	// my/import.Foobar => alias.Foobar
	// Foobar => Foobar
	ResolveFunction(string) (string, error)
}

type TypeResolver interface {
	// *my/import.Book => *alias.Book
	// my/import.Book => alis.Book
	ResolveType(string) (string, error)
}

type ServiceResolver interface {
	// @person.(*gontainer/example/pkg.Person) => person, *alias.Person
	// @person => person, ""
	ResolveService(string) (string, string, error)
}

type SimpleFunctionResolver struct {
	imports imports.Imports
}

type SimpleTypeResolver struct {
	imports imports.Imports
}

type SimpleServiceResolver struct {
	typeResolver TypeResolver
}

// todo improve regexps
var (
	importPart         = `(?P<import>([A-Za-z][A-Z-a-z0-9._/-]*\.)*)`
	namePart           = `(?P<name>[A-Za-z_][A-Za-z0-9_]*)`
	ServiceNamePattern = `(?P<service>[A-Za-z]([._]?[A-Za-z0-9])*)` // todo move to regex, see regex.MetaServiceName
	fnRegex            = regexp.MustCompile(`^` + importPart + namePart + `$`)
	typeRegex          = regexp.MustCompile(`^(?P<pointer>\*?)` + importPart + namePart + `$`)
	serviceRegex       = regexp.MustCompile(`^@` + ServiceNamePattern + `(\.\((?P<type>` + strings.Trim(typeRegex.String(), `^$`) + `)\))?$`)
)

func (s SimpleFunctionResolver) ResolveFunction(i string) (string, error) {
	match, params := regex.Match(fnRegex, i)
	if !match {
		return "", fmt.Errorf(
			"invalid syntax, function must follow pattern `%s`, `%s` given",
			fnRegex.String(),
			i,
		)
	}

	if params["import"] == "" {
		return params["name"], nil
	}

	imp := strings.TrimRight(params["import"], ".")

	return s.imports.GetAlias(imp) + "." + params["name"], nil
}

func NewSimpleFunctionResolver(imports imports.Imports) *SimpleFunctionResolver {
	return &SimpleFunctionResolver{imports: imports}
}

func (s SimpleTypeResolver) ResolveType(i string) (string, error) {
	match, params := regex.Match(typeRegex, i)
	if !match {
		return "", fmt.Errorf(
			"invalid syntax, type must follow pattern `%s`, `%s` given",
			typeRegex.String(),
			i,
		)
	}

	if params["import"] == "" {
		return params["pointer"] + params["name"], nil
	}

	imp := strings.TrimRight(params["import"], ".")

	return params["pointer"] + s.imports.GetAlias(imp) + "." + params["name"], nil
}

func NewSimpleTypeResolver(imports imports.Imports) *SimpleTypeResolver {
	return &SimpleTypeResolver{imports: imports}
}

func (s SimpleServiceResolver) ResolveService(i string) (string, string, error) {
	match, params := regex.Match(serviceRegex, i)
	if !match {
		return "", "", fmt.Errorf(
			"invalid syntax, service must follow pattern `%s`, `%s` given",
			serviceRegex.String(),
			i,
		)
	}

	if params["type"] == "" {
		return params["service"], "", nil
	}

	t, err := s.typeResolver.ResolveType(params["type"])

	if err != nil {
		return "", "", err
	}

	return params["service"], t, nil
}

func NewSimpleServiceResolver(typeResolver TypeResolver) *SimpleServiceResolver {
	return &SimpleServiceResolver{typeResolver: typeResolver}
}
