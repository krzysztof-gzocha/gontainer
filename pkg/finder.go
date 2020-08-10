package pkg

import (
	"path/filepath"
	"sort"
)

type finder interface {
	find([]string) ([]string, error)
}

type simpleFinder struct {
}

func (s simpleFinder) find(patterns []string) ([]string, error) {
	result := make([]string, 0)
	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err != nil {
			return nil, err
		}
		sort.Strings(matches)
		result = append(result, matches...)
	}
	return result, nil
}
