package input

import (
	"reflect"
)

var (
	primitiveKinds = []reflect.Kind{
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

func isPrimitiveType(val interface{}) bool {
	if val == nil {
		return true
	}

	t := reflect.TypeOf(val)

	for _, k := range primitiveKinds {
		if k == t.Kind() {
			return true
		}
	}

	return false
}
