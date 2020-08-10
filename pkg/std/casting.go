package std

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/exporters"
)

var (
	defaultExporter = exporters.NewChainExporter(
		&exporters.BoolExporter{},
		&exporters.NilExporter{},
		&exporters.NumericExporter{},
	)
)

// TODO remove
// use exporters.MustToString
func MustConvertToString(v interface{}) string {
	if r, ok := v.(string); ok {
		return r
	}

	r, err := defaultExporter.Export(v)

	if err != nil {
		panic(fmt.Sprintf("cannot cast parameter of type `%T` to string: %s", v, err.Error()))
	}

	return r
}
