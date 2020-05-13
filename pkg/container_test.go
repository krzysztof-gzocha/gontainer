package pkg

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseParamContainer(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		hasNot []string
	}{
		{
			name:   "empty container",
			params: nil,
			hasNot: []string{"param"},
		},
		{
			name:   "with parameters",
			params: map[string]interface{}{"one": 1, "pi": math.Pi, "name": "gontainer"},
			hasNot: []string{"e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewBaseParamContainer(tt.params)
			for id, val := range tt.params {
				assert.True(t, container.HasParam(id))
				assert.Equal(t, val, container.MustGetParam(id))

				realValue, err := container.GetParam(id)
				assert.NoError(t, err)
				assert.Equal(t, val, realValue)
			}
			for _, id := range tt.hasNot {
				assert.False(t, container.HasParam(id))
				assert.Panics(t, func() {
					container.MustGetParam(id)
				})

				realValue, err := container.GetParam(id)
				assert.Error(t, err)
				assert.Nil(t, realValue)
			}
		})
	}
}
