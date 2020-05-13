package parameters

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/gomponents/gontainer/pkg/tokens"
)

type Validator interface {
	Validate(ResolvedParams) []string
}

const (
	RegexParamName = tokens.RegexTokenReference
)

type SimpleValidator struct{}

func (s SimpleValidator) Validate(parameters ResolvedParams) []string {
	msgs := make([]string, 0)
	for id, v := range parameters {
		msgs = appendNotNilString(
			msgs,
			s.validateID(id),
			s.validateType(v),
		)
	}

	return msgs
}

func (SimpleValidator) validateID(id string) *string {
	if !regexp.MustCompile(RegexParamName).MatchString(id) {
		return stringPointer(fmt.Sprintf("invalid param name `%s`", id))
	}

	return nil
}

func (SimpleValidator) validateType(val interface{}) *string {
	if val == nil {
		return nil
	}

	allowedKinds := []reflect.Kind{
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
			return nil
		}
	}

	return stringPointer(fmt.Sprintf("unsupported type `%T`", val))
}

func stringPointer(v string) *string {
	return &v
}

func appendNotNilString(slice []string, elems ...*string) []string {
	for _, e := range elems {
		if e != nil {
			slice = append(slice, *e)
		}
	}

	return slice
}
