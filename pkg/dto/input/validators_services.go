package input

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexServiceName = regexp.MustCompile("^" + regex.MetaServiceName + "$")
)

type ValidateService func(string, Service) error

func DefaultServicesValidators() []func(DTO) error {
	validators := []ValidateService{
		func(s string, service Service) error {
			// todo remove
		},
	}

	return []func(DTO) error{
		func(d DTO) error {
			for n, s := range d.Services {
				if err := ValidateServiceName(n, d); err != nil {
					return err
				}
				if s.Todo {
					continue
				}
				for _, v := range validators {
					err := v(n, s)
					if err == nil {
						continue
					}
					return fmt.Errorf("service `%s`: %s", n, err.Error())
				}
			}
		},
	}
}

func ValidateServiceName(n string, _ Service) error {
	if !regexServiceName.MatchString(n) {
		return fmt.Errorf(
			"service name must match pattern `%s`, `%s` given",
			regexServiceName.String(),
			n,
		)
	}
	return nil
}
