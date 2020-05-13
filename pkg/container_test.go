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
	}{
		{
			name:   "empty container",
			params: nil,
		},
		{
			name:   "with parameters",
			params: map[string]interface{}{"one": 1, "pi": math.Pi, "name": "gontainer"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewBaseParamContainer(tt.params)
			for id, val := range tt.params {
				assert.Equal(t, val, container.MustGetParam(id))
			}
		})
	}
}
