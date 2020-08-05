package definition

import (
	"fmt"
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

// TODO improve impP regex
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
