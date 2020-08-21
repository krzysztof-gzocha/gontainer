package input

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexParamName = regexp.MustCompile("^" + regex.ParamName + "$")
)

// DefaultParamsValidators returns validators for DTO.Params.
func DefaultParamsValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateParams,
	}
}

func ValidateParams(d DTO) error {
	for k, v := range d.Params {
		if !regexParamName.MatchString(k) {
			return fmt.Errorf("parameter name should match `%s`, `%s` given", regexParamName.String(), k)
		}

		if !isPrimitiveType(v) {
			return fmt.Errorf("unsupported type `%T` of parameter `%s`", v, k)
		}
	}
	return nil
}
