package std

import (
	"math"
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
			name:   "new line",
			input:  "hello\nworld",
			output: "hello\nworld",
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
		{
			name:   "complex128(534)",
			input:  complex128(534 + 0i),
			output: "(534+0i)",
			panic:  false,
		},
		{
			name:   "complex128(math.Pi)",
			input:  complex128(math.Pi),
			output: "(3.141592653589793+0i)",
			panic:  false,
		},
		{
			name:   "complex64(math.Pi)",
			input:  complex64(math.Pi),
			output: "(3.1415927+0i)",
			panic:  false,
		},
		{
			name:   "float32(math.Pi)",
			input:  float32(math.Pi),
			output: "3.1415927",
			panic:  false,
		},
		{
			name:   "float32(0.34)",
			input:  float32(0.34),
			output: "0.34",
			panic:  false,
		},
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
