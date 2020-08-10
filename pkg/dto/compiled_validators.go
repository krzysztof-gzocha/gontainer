package dto

import (
	"fmt"
	"sort"
	"strings"
)

type CompiledValidator interface {
	Validate(CompiledInput) error
}

type ChainCompiledValidator struct {
	validators []func(CompiledInput) error
}

func NewChainCompiledValidator(validators []func(CompiledInput) error) *ChainCompiledValidator {
	return &ChainCompiledValidator{validators: validators}
}

func (c *ChainCompiledValidator) Validate(i CompiledInput) error {
	for _, v := range c.validators {
		if err := v(i); err != nil {
			return err
		}
	}
	return nil
}

func NewDefaultCompiledValidator() CompiledValidator {
	return NewChainCompiledValidator([]func(CompiledInput) error{
		validateCircularDependency,
	})
}

type seekableStringSlice []string

func (sss seekableStringSlice) contains(s string) bool {
	for _, n := range sss {
		if n == s {
			return true
		}
	}

	return false
}

func findCircularDep(n string, deps map[string][]string, path seekableStringSlice) []string {
	r := append(path, n)
	for _, subdep := range deps[n] {
		if r.contains(subdep) {
			return append(r, subdep)
		}

		if found := findCircularDep(subdep, deps, append(r)); found != nil {
			return found
		}
	}
	return nil
}

func validateCircularDependency(i CompiledInput) error {
	services := make([]string, 0)
	deps := make(map[string][]string)
	for _, s := range i.Services {
		services = append(services, s.Name)
		deps[s.Name] = make([]string, 0)
		for _, a := range s.Args {
			if !a.IsService() {
				continue
			}
			deps[s.Name] = append(deps[s.Name], a.ServiceLink.Name)
		}
	}

	sort.Strings(services)

	for _, n := range services {
		found := findCircularDep(n, deps, nil)
		if found == nil {
			continue
		}
		return fmt.Errorf("found circular dependency %s", strings.Join(found, " -> "))
	}

	return nil
}
