package input

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetaPkg(t *testing.T) {
	scenarios := []struct {
		pkg   string
		error string
	}{
		{
			pkg:   "",
			error: "meta.pkg cannot be empty",
		},
		{
			pkg: "main",
		},
		{
			pkg:   "123",
			error: "meta.pkg must match " + regexpMetaPkg.String() + ", `123` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			dto := DTO{}
			dto.Meta.Pkg = s.pkg
			err := ValidateMetaPkg(dto)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, s.error, err.Error())
		})
	}
}

func TestValidateMetaContainerType(t *testing.T) {
	scenarios := []struct {
		containerType string
		error         string
	}{
		{
			containerType: "",
			error:         "meta.container_type must match " + regexpMetaContainerType.String() + ", `` given",
		},
		{
			containerType: "myContainer123",
		},
		{
			containerType: "0MyContainer",
			error:         "meta.container_type must match " + regexpMetaContainerType.String() + ", `0MyContainer` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			dto := DTO{}
			dto.Meta.ContainerType = s.containerType
			err := ValidateMetaContainerType(dto)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, s.error, err.Error())
		})
	}
}

func TestValidateMetaImports(t *testing.T) {
	scenarios := []struct {
		import_ string
		alias   string
		error   string
	}{
		{
			import_: "github.com/stretchr/testify/assert",
			alias:   "assert",
		},
		{
			import_: "oneTwoThree",
			alias:   "$123",
			error:   "invalid alias `$123`, must match `" + regexMetaImportAlias.String() + "`",
		},
		{
			import_: "!!!",
			alias:   "alias",
			error:   "invalid import `!!!`, must match `" + regexMetaImport.String() + "`",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			d := DTO{}
			d.Meta.Imports = map[string]string{
				s.alias: s.import_,
			}

			err := ValidateMetaImports(d)

			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, s.error, err.Error())
		})
	}
}

func TestValidateMetaFunctions(t *testing.T) {
	scenarios := []struct {
		alias string
		goFn  string
		error string
	}{
		{
			alias: "env",
			goFn:  "os.Getenv",
		},
		{
			alias: "env",
			goFn:  "env",
		},
		{
			alias: "$fn",
			goFn:  "os.Getenv",
			error: "invalid function `$fn`, must match `" + regexMetaFn.String() + "`",
		},
		{
			alias: "env",
			goFn:  "os.1Getenv",
			error: "invalid go function `os.1Getenv`, must match `" + regexMetaGoFn.String() + "`",
		},
	}
	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			d := DTO{}
			d.Meta.Functions = map[string]string{
				s.alias: s.goFn,
			}
			err := ValidateMetaFunctions(d)

			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, s.error, err.Error())
		})
	}
}
