package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestTag_UnmarshalYAML(t *testing.T) {
	y1 := `
name: hello
priority: 100
`
	t1 := Tag{}
	err1 := yaml.Unmarshal([]byte(y1), &t1)
	assert.NoError(t, err1)
	assert.Equal(t, Tag{Name: "hello", Priority: 100}, t1)

	y2 := `tag`
	t2 := Tag{}
	err2 := yaml.Unmarshal([]byte(y2), &t2)
	assert.NoError(t, err2)
	assert.Equal(t, Tag{Name: "tag", Priority: 0}, t2)
}
