package input

import (
	"fmt"
	"regexp"
)

func DefaultMetaValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateMetaPkg,
		ValidateMetaImports,
		ValidateMetaContainerType,
	}
}

func ValidateMetaPkg(d DTO) error {
	if d.Meta.Pkg == "" {
		return fmt.Errorf("Meta.Pkg cannot be empty")
	}

	r := "^[a-z][a-zA-Z0-9_]*$"

	if !regexp.MustCompile(r).MatchString(d.Meta.Pkg) {
		return fmt.Errorf("meta.pkg must match %s, `%s` given", r, d.Meta.Pkg)
	}

	return nil
}

func ValidateMetaImports(d DTO) error {
	// todo
	return nil
}

func ValidateMetaContainerType(d DTO) error {
	// todo
	return nil
}
