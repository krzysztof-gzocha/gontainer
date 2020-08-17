package dto

import (
	"fmt"
	"regexp"
)

func ValidateMetaPkg(i Input) error {
	if i.Meta.Pkg == "" {
		return fmt.Errorf("meta.pkg cannot be empty")
	}

	r := "^[a-z][a-zA-Z0-9_]*$"

	if !regexp.MustCompile(r).MatchString(i.Meta.Pkg) {
		return fmt.Errorf("meta.pkg must match %s, `%s` given", r, i.Meta.Pkg)
	}

	return nil
}

// TODO improve impP regex e.g. "/aa" it's not valid import
func ValidateMetaImports(i Input) error {
	aliasP := "^[a-zA-Z0-9_]+$"
	aliasR := regexp.MustCompile(aliasP)

	impP := "^[a-zA-Z0-9_./]+$"
	impR := regexp.MustCompile(impP)

	for alias, imp := range i.Meta.Imports {
		if !aliasR.MatchString(alias) {
			return fmt.Errorf("invalid import alias `%s`, must match `%s`", alias, aliasP)
		}

		if !impR.MatchString(imp) {
			return fmt.Errorf("invalid import `%s`, must match `%s`", imp, impP)
		}
	}

	return nil
}

func ValidateMetaContainerType(i Input) error {
	p := "^[a-z][A-Za-z0-9_]*$"

	if !regexp.MustCompile(p).MatchString(i.Meta.ContainerType) {
		return fmt.Errorf("meta.container_type must match %s, `%s` given", p, i.Meta.ContainerType)
	}

	return nil
}
