package parameters

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/tokens"
	"github.com/stretchr/testify/assert"
)

func TestBagFactory_Create(t *testing.T) {
	scenarios := []struct {
		input  map[string]interface{}
		output map[string]compiled.Param
		error  string
	}{
		{
			input: map[string]interface{}{
				"name":      "Jane",
				"full_name": "%name%",
			},
			output: map[string]compiled.Param{
				"full_name": {
					Code: `result.MustGetParam("name")`,
					Raw:  "%name%",
				},
				"name": {
					Code: `"Jane"`,
					Raw:  "Jane",
				},
			},
		},
		{
			input: map[string]interface{}{
				"name":      "%full_name%",
				"full_name": "%name%",
			},
			output: nil,
			error:  "cannot solve param `full_name`, circular dependencies: full_name, name",
		},
	}

	factory := NewBagFactory(
		tokens.NewPatternTokenizer([]tokens.TokenFactoryStrategy{
			tokens.TokenPercentSign{},
			tokens.TokenReference{},
			tokens.TokenString{},
		}),
		exporters.NewDefaultExporter(),
		imports.NewSimpleImports(),
	)

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			result, err := factory.Create(s.input)
			if s.error != "" {
				if assert.Error(t, err) {
					assert.Equal(t, s.error, err.Error())
				}
				assert.Nil(t, result)
				return
			}

			assert.Equal(t, s.output, result)
		})
	}
}
