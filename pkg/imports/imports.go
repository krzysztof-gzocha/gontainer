package imports

import (
	"sort"
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
}

type SimpleImports struct {
	prefix  string
	counter int64
	imports map[string]Import
}

func (s *SimpleImports) GetAlias(path string) string {
	if i, ok := s.imports[path]; ok {
		return i.Alias
	}

	parts := strings.Split(path, "/")

	i := Import{
		Path:  path,
		Alias: strings.Replace(parts[len(parts)-1], ".", "_", 999) + "_" + s.prefix + strconv.FormatInt(s.counter, 16),
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

	sort.Slice(r, func(i, j int) bool {
		return strings.Compare(r[i].Path, r[j].Path) < 0
	})

	return r
}

func NewSimpleImports(prefix string) *SimpleImports {
	return &SimpleImports{
		prefix:  prefix,
		imports: make(map[string]Import),
	}
}
