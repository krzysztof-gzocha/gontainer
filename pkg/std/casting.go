package std

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/exporters"
)

func MustConvertToString(v interface{}) string {
	if r, ok := v.(string); ok {
		return r
	}

	e := exporters.NewChainExporter(
		&exporters.BoolExporter{},
		&exporters.NilExporter{},
		&exporters.NumericExporter{},
	)

	r, err := e.Export(v)

	if err != nil {
		panic(fmt.Sprintf("cannot cast parameter of type `%T` to string: %s", v, err.Error()))
	}

	return r
}
