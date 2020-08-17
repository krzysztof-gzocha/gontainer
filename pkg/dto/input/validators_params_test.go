package input

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/stretchr/testify/assert"
)

func TestValidateParams(t *testing.T) {
	scenarios := []struct {
		name  string
		val   interface{}
		error string
	}{
		{
			name: "test",
			val:  5,
		},
		{
			name:  "@",
			error: "parameter name should match `^" + regexParamName.String() + "`, `@` given",
		},
		{
			name:  "param",
			val:   struct{}{},
			error: "unsupported type `struct {}` of parameter `param`",
		},
	}
	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			d := DTO{}
			d.Params = parameters.RawParameters{
				s.name: s.val,
			}
			err := ValidateParams(d)

			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, s.error, err.Error())
		})
	}
}
