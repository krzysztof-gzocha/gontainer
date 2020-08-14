package dto

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/gomponents/gontainer-helpers/container"
	"github.com/gomponents/gontainer/pkg/syntax"
)

type serviceValidator func(n string, s Service) error

func ValidateServices(i Input) error {
	validators := []serviceValidator{
		ValidateServicesNames,
		ValidateServicesGetters,
		ValidateBuildingMethod,
	}
	for _, v := range validators {
		for n, d := range i.Services {
			if err := v(n, d); err != nil {
				return err
			}
		}
	}
	return nil
}

var (
	serviceNameRegex   = regexp.MustCompile(`^` + syntax.ServiceNamePattern + `$`)
	serviceGetterRegex = regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*$`)
)

func ValidateServicesNames(n string, s Service) error {
	if !serviceNameRegex.MatchString(n) {
		return fmt.Errorf(
			"service name must match pattern `%s`, `%s` given",
			serviceNameRegex.String(),
			n,
		)
	}

	reserved := []string{"serviceContainer"}
	for _, w := range reserved {
		if w == n {
			return fmt.Errorf("service `%s`: getter `%s` is reserved", n, w)
		}
	}

	return nil
}

func ValidateServicesGetters(n string, s Service) error {
	if s.Getter != "" && !serviceGetterRegex.MatchString(s.Getter) {
		return fmt.Errorf(
			"service `%s`: getter must match `%s`, `%s` given",
			n,
			serviceGetterRegex.String(),
			s.Getter,
		)
	}

	var reserved []string
	r := reflect.TypeOf(struct {
		container.BaseContainer
		container.BaseParamContainer
	}{})
	for i := 0; i < r.NumMethod(); i++ {
		reserved = append(reserved, r.Method(i).Name)
	}
	reserved = append(reserved, "ValidateAllServices")

	for _, w := range reserved {
		if w == s.Getter {
			return fmt.Errorf("service `%s`: getter `%s` is reserved", n, w)
		}
	}

	return nil
}

func ValidateBuildingMethod(n string, s Service) error {
	if s.Constructor == "" && s.Type == "" {
		return fmt.Errorf("service `%s`: missing contructor or type", n)
	}

	if s.Constructor == "" && len(s.Args) > 0 {
		return fmt.Errorf("service `%s`: arguments are given, but constructor is missing", n)
	}

	if s.Getter != "" && s.Type == "" {
		return fmt.Errorf("service `%s`: getter is given, but type is missing", n)
	}

	return nil
}
