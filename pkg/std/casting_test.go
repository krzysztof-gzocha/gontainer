package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustConvertToString(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		output string
		panic  bool
	}{
		{
			name:   "strings",
			input:  "hello",
			output: "hello",
			panic:  false,
		},
		{
			name:  "panic",
			input: struct{}{},
			panic: true,
		},
		{
			name:   "nil",
			input:  nil,
			output: "nil",
			panic:  false,
		},
		{
			name:   "bool",
			input:  true,
			output: "true",
			panic:  false,
		},
		{
			name:   "uint32(534)",
			input:  uint32(534),
			output: "534",
			panic:  false,
		},
		//{
		//	name:   "complex128(534)",
		//	input:  complex128(534),
		//	output: "534",
		//	panic:  false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Panics(t, func() {
					MustConvertToString(tt.input)
				})
				return
			}
			assert.Equal(t, tt.output, MustConvertToString(tt.input))
		})
	}
}
