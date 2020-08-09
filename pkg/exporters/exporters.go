package exporters

import (
	"errors"
	"fmt"
	"reflect"
)

type Exporter interface {
	Export(interface{}) (string, error)
}

type Subexporter interface {
	Exporter
	Supports(interface{}) bool
}

type ChainExporter struct {
	exporters []Subexporter
}

func (c ChainExporter) Export(v interface{}) (string, error) {
	for _, e := range c.exporters {
		if e.Supports(v) {
			return e.Export(v)
		}
	}

	return "", errors.New(fmt.Sprintf("parameter of type `%T` is not supported", v))
}

func NewDefaultExporter() Exporter {
	return NewChainExporter(
		&BoolExporter{},
		&NilExporter{},
		&NumericExporter{},
		&StringExporter{},
	)
}

func NewChainExporter(exporters ...Subexporter) *ChainExporter {
	return &ChainExporter{exporters: exporters}
}

type BoolExporter struct{}

func (b BoolExporter) Export(v interface{}) (string, error) {
	if v == true {
		return "true", nil
	}

	return "false", nil
}

func (b BoolExporter) Supports(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}

type NilExporter struct{}

func (n NilExporter) Export(interface{}) (string, error) {
	return "nil", nil
}

func (n NilExporter) Supports(v interface{}) bool {
	return v == nil
}

type NumericExporter struct{}

func (n NumericExporter) Export(v interface{}) (string, error) {
	switch reflect.TypeOf(v).Kind() {
	case
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		return fmt.Sprintf("%#v", v), nil
	}
	return fmt.Sprintf("%d", v), nil
}

func (n NumericExporter) Supports(v interface{}) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}

	allowedKinds := []reflect.Kind{
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

	for _, k := range allowedKinds {
		if k == t.Kind() {
			return true
		}
	}

	return false
}

type StringExporter struct{}

func (s StringExporter) Export(v interface{}) (string, error) {
	return fmt.Sprintf("%+q", v), nil
}

func (s StringExporter) Supports(v interface{}) bool {
	_, ok := v.(string)
	return ok
}
