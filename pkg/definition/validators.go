package definition

import (
	"fmt"
	"reflect"
	"regexp"
)

type Validator interface {
	Validate(Definition) error
}

type ChainValidator struct {
	validators []func(Definition) error
}

func NewChainValidator(validators []func(Definition) error) *ChainValidator {
	return &ChainValidator{validators: validators}
}

func (c ChainValidator) Validate(d Definition) error {
	for _, v := range c.validators {
		err := v(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewDefaultValidator() Validator {
	return NewChainValidator([]func(Definition) error{
		ValidateMetaPkg,
		ValidateMetaImports,
		ValidateMetaContainerType,
		ValidateParams,
	})
}

func ValidateMetaPkg(d Definition) error {
	if d.Meta.Pkg == "" {
		return fmt.Errorf("meta.pkg cannot be empty")
	}

	r := "^[a-z][a-zA-Z0-9_]*$"

	if !regexp.MustCompile(r).MatchString(d.Meta.Pkg) {
		return fmt.Errorf("meta.pkg must match %s, `%s` given", r, d.Meta.Pkg)
	}

	return nil
}

// TODO improve impP regex e.g. "/aa" it's not valid import
func ValidateMetaImports(d Definition) error {
	aliasP := "^[a-zA-Z0-9_]+$"
	aliasR := regexp.MustCompile(aliasP)

	impP := "^[a-zA-Z0-9_./]+$"
	impR := regexp.MustCompile(impP)

	for alias, imp := range d.Meta.Imports {
		if !aliasR.MatchString(alias) {
			return fmt.Errorf("invalid import alias `%s`, must match `%s`", alias, aliasP)
		}

		if !impR.MatchString(imp) {
			return fmt.Errorf("invalid import `%s`, must match `%s`", imp, impP)
		}
	}

	return nil
}

func ValidateMetaContainerType(d Definition) error {
	p := "^[A-Za-z][A-Za-z0-9_]*$"

	if !regexp.MustCompile(p).MatchString(d.Meta.ContainerType) {
		return fmt.Errorf("meta.container_type must match %s, `%s` given", p, d.Meta.ContainerType)
	}

	return nil
}

func ValidateParams(d Definition) error {
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

	for k, v := range d.Params {
		if !nameR.MatchString(k) {
			return fmt.Errorf("parameter name should match `%s`, `%s` given", nameP, k)
		}

		if !validateType(v) {
			return fmt.Errorf("unsupported type `%T` of parameter `%s`", v, k)
		}
	}

	return nil
}
