package imports

import (
	"fmt"
	"strconv"
	"strings"
)

type Import struct {
	Path  string
	Alias string
}

type Imports interface {
	GetAlias(string) string
	GetImports() []Import
	RegisterPrefix(shortcut string, path string) error
}

type SimpleImports struct {
	prefix    string
	counter   int64
	imports   map[string]Import
	shortcuts map[string]string
}

func (s *SimpleImports) RegisterPrefix(shortcut string, path string) error {
	if _, ok := s.shortcuts[shortcut]; ok {
		return fmt.Errorf("shortcut `%s` is already registered", shortcut)
	}

	s.shortcuts[shortcut] = path
	return nil
}

func (s *SimpleImports) GetAlias(path string) string {
	if i, ok := s.imports[path]; ok {
		return i.Alias
	}

	parts := strings.Split(path, "/")

	i := Import{
		Path:  path,
		Alias: strings.ReplaceAll(parts[len(parts)-1], ".", "_") + "_" + s.prefix + strconv.FormatInt(s.counter, 16),
	}
	s.imports[path] = i
	s.counter++

	return i.Alias
}

func (s *SimpleImports) GetImports() []Import {
	r := make([]Import, 0)
	for _, i := range s.imports {
		r = append(r, i)
	}

	return r
}

func NewSimpleImports(prefix string) *SimpleImports {
	return &SimpleImports{
		prefix:  prefix,
		imports: make(map[string]Import),
	}
}
