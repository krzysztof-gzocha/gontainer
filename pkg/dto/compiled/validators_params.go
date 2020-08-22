package compiled

import (
	"fmt"
	"sort"
	"strings"
)

func DefaultParamsValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateParams,
	}
}

func ValidateParams(d DTO) error {
	mapping := make(map[string]Param)
	params := make([]Param, 0)
	for _, p := range d.Params {
		mapping[p.Name] = p
		params = append(params, p)
	}
	sort.SliceStable(params, func(i, j int) bool {
		return params[i].Name < params[j].Name
	})

	finder := newCircularDepFinder(func(id string) []string {
		param, _ := mapping[id]
		deps := param.DependsOn
		sort.Strings(deps)
		return deps
	})

	for _, p := range params {
		if deps := finder.find(p.Name); deps != nil {
			return fmt.Errorf("circular dependency in params: %s", strings.Join(deps, " -> "))
		}
	}

	return nil
}
