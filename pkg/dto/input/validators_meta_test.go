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

func TestValidateMetaImports(t *testing.T) {
	// todo
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
