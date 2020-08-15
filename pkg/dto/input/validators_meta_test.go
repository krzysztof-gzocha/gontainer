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
			error: "Meta.Pkg cannot be empty",
		},
		{
			pkg: "main",
		},
		{
			pkg:   "123",
			error: "meta.pkg must match ^[a-z][a-zA-Z0-9_]*$, `123` given",
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
