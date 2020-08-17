package input

import (
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestCall_UnmarshalYAML(t *testing.T) {
	// todo more tests
	y := `
- ["SetName", ["Mary"], true]
`
	calls := make([]Call, 0)
	err := yaml.Unmarshal([]byte(y), &calls)
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]Call{
			{
				Method:    "SetName",
				Args:      []interface{}{"Mary"},
				Immutable: true,
			},
		},
		calls,
	)
}

func TestCreateDefaultDTO(t *testing.T) {
	assert.NoError(
		t,
		NewDefaultValidator().Validate(CreateDefaultDTO()),
	)
}
