package input

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexServiceName = regexp.MustCompile("^" + regex.MetaServiceName + "$")
)

type ValidateService func(Service) error

func DefaultServicesValidators() []func(DTO) error {
	validators := []ValidateService{
		ValidateConstructorType,
	}

	return []func(DTO) error{
		func(d DTO) error {
			for n, s := range d.Services {
				if err := ValidateServiceName(n); err != nil {
					return err
				}
				if s.Todo {
					continue
				}
				for _, v := range validators {
					err := v(s)
					if err == nil {
						continue
					}
					return fmt.Errorf("service `%s`: %s", n, err.Error())
				}
			}
			return nil
		},
	}
}

func ValidateServiceName(n string) error {
	if !regexServiceName.MatchString(n) {
		return fmt.Errorf(
			"service name must match pattern `%s`, `%s` given",
			regexServiceName.String(),
			n,
		)
	}
	return nil
}

func ValidateConstructorType(s Service) error {
	if s.Constructor == "" && s.Value == "" && s.Type == "" {
		return fmt.Errorf("missing constructor, value or type")
	}

	if s.Constructor != "" && s.Value != "" {
		return fmt.Errorf("cannot define constructor and value together")
	}

	// e.g.
	// Service{
	//		Getter:      "",
	//		Type:        "MyType",
	//		Value:       "",
	//		Constructor: "NewService",
	//	}
	if (s.Constructor != "" || s.Value != "") && (s.Getter == "" && s.Type != "") {
		return fmt.Errorf("defined type will not be used, provide getter")
	}

	if len(s.Args) > 0 && s.Constructor == "" {
		return fmt.Errorf("arguments are not empty, but constructor is missing")
	}

	return nil
}
