package dto

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/syntax"
)

type serviceValidator func(n string, s Service) error

func ValidateServices(i Input) error {
	validators := []serviceValidator{
		ValidateServicesNames,
		ValidateServicesGetters,
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
	serviceGetterRegex = regexp.MustCompile(`^[A-Z][A-Za-z0-9_]$`)
)

func ValidateServicesNames(n string, _ Service) error {
	if !serviceNameRegex.MatchString(n) {
		return fmt.Errorf(
			"service name must match pattern `%s`, `%s` given",
			serviceNameRegex.String(),
			n,
		)
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

	return nil
}
