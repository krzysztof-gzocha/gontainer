package input

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexParamName     = regexp.MustCompile("^" + regex.MetaParamName + "$")
	allowedParamsKinds = []reflect.Kind{
		reflect.String,
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
	}
)

func DefaultParamsValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateParams,
	}
}

func validateParamType(val interface{}) bool {
	if val == nil {
		return true
	}

	t := reflect.TypeOf(val)

	for _, k := range allowedParamsKinds {
		if k == t.Kind() {
			return true
		}
	}

	return false
}

func ValidateParams(d DTO) error {
	for k, v := range d.Params {
		if !regexParamName.MatchString(k) {
			return fmt.Errorf("parameter name should match `%s`, `%s` given", regexParamName.String(), k)
		}

		if !validateParamType(v) {
			return fmt.Errorf("unsupported type `%T` of parameter `%s`", v, k)
		}
	}
	return nil
}
