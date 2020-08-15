package input

import (
	"fmt"
	"regexp"
)

var (
	regexpMetaPkg           = regexp.MustCompile("^[a-z][A-Za-z0-9_]*$")
	regexpMetaContainerType = regexpMetaPkg
)

func DefaultMetaValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateMetaPkg,
		ValidateMetaImports,
		ValidateMetaContainerType,
		ValidateMetaFunctions,
	}
}

func ValidateMetaPkg(d DTO) error {
	if d.Meta.Pkg == "" {
		return fmt.Errorf("meta.pkg cannot be empty")
	}

	if !regexpMetaPkg.MatchString(d.Meta.Pkg) {
		return fmt.Errorf(
			"meta.pkg must match %s, `%s` given",
			regexpMetaPkg.String(),
			d.Meta.Pkg,
		)
	}

	return nil
}

func ValidateMetaContainerType(d DTO) error {
	if !regexpMetaContainerType.MatchString(d.Meta.ContainerType) {
		return fmt.Errorf(
			"meta.container_type must match %s, `%s` given",
			regexpMetaContainerType,
			d.Meta.ContainerType,
		)
	}
	return nil
}

func ValidateMetaImports(d DTO) error {
	// todo
	return nil
}

func ValidateMetaFunctions(d DTO) error {
	// todo
	return nil
}
