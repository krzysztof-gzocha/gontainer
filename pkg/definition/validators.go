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
