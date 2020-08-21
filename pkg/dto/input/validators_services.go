package input

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexServiceName        = regexp.MustCompile("^" + regex.ServiceName + "$")
	regexServiceGetter      = regexp.MustCompile("^" + regex.ServiceGetter + "$")
	regexServiceType        = regexp.MustCompile("^" + regex.ServiceType + "$")
	regexServiceValue       = regexp.MustCompile("^" + regex.ServiceValue + "$")
	regexServiceConstructor = regexp.MustCompile("^" + regex.ServiceConstructor + "$")
	regexServiceCallName    = regexp.MustCompile("^" + regex.ServiceCallName + "$")
	regexServiceFieldName   = regexp.MustCompile("^" + regex.ServiceFieldName + "$")
	regexServiceTag         = regexp.MustCompile("^" + regex.ServiceTag + "$")
)

type ValidateService func(Service) error

// DefaultServicesValidators returns validators for DTO.Services.
func DefaultServicesValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateServices,
	}
}

func ValidateServices(d DTO) error {
	validators := []ValidateService{
		ValidateConstructorType,
		ValidateServiceGetter,
		ValidateServiceType,
		ValidateServiceValue,
		ValidateServiceConstructor,
		ValidateServiceArgs,
		ValidateServiceCalls,
		ValidateServiceFields,
		ValidateServiceTags,
	}

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

func ValidateServiceGetter(s Service) error {
	return validateRegexField("getter", s.Getter, regexServiceGetter, true)
}

func ValidateServiceType(s Service) error {
	return validateRegexField("type", s.Type, regexServiceType, true)
}

func ValidateServiceValue(s Service) error {
	return validateRegexField("value", s.Value, regexServiceValue, true)
}

func ValidateServiceConstructor(s Service) error {
	return validateRegexField("constructor", s.Constructor, regexServiceConstructor, true)
}

func ValidateServiceArgs(s Service) error {
	for i, a := range s.Args {
		if !isPrimitiveType(a) {
			return fmt.Errorf("unsupported type `%T` of arg%d", a, i)
		}
	}
	return nil
}

func ValidateServiceCalls(s Service) error {
	for j, c := range s.Calls {
		if err := validateRegexField(fmt.Sprintf("call%d", j), c.Method, regexServiceCallName, false); err != nil {
			return err
		}
		for i, a := range c.Args {
			if !isPrimitiveType(a) {
				return fmt.Errorf("unsupported type `%T` of call%d.arg%d", a, j, i)
			}
		}
	}
	return nil
}

func ValidateServiceFields(s Service) error {
	for n, v := range s.Fields {
		if err := validateRegexField("field", n, regexServiceFieldName, false); err != nil {
			return err
		}
		if !isPrimitiveType(v) {
			return fmt.Errorf("unsupported type `%T` of field `%s`", v, n)
		}
	}
	return nil
}

func ValidateServiceTags(s Service) error {
	for _, t := range s.Tags {
		if err := validateRegexField("tag", t, regexServiceTag, false); err != nil {
			return err
		}
	}
	return nil
}

func validateRegexField(field string, value string, expr *regexp.Regexp, optional bool) error {
	if optional && value == "" {
		return nil
	}
	if !expr.MatchString(value) {
		return fmt.Errorf("%s must match `%s`, `%s` given", field, expr.String(), value)
	}
	return nil
}
