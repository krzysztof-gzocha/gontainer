package dto

import (
	"fmt"
	"reflect"
	"regexp"
)

func ValidateParams(i Input) error {
	validateType := func(val interface{}) bool {
		if val == nil {
			return true
		}

		allowedKinds := []reflect.Kind{
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

		t := reflect.TypeOf(val)

		for _, k := range allowedKinds {
			if k == t.Kind() {
				return true
			}
		}

		return false
	}

	nameP := "^[A-Za-z]([_.]?[A-Za-z0-9])*$"
	nameR := regexp.MustCompile(nameP)

	for k, v := range i.Params {
		if !nameR.MatchString(k) {
			return fmt.Errorf("parameter name should match `%s`, `%s` given", nameP, k)
		}

		if !validateType(v) {
			return fmt.Errorf("unsupported type `%T` of parameter `%s`", v, k)
		}
	}

	return nil
}
