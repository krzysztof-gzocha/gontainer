package exporters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainExporter_Export(t *testing.T) {
	exporter := NewDefaultExporter()

	t.Run("Given valid scenarios", func(t *testing.T) {
		scenarios := map[string]struct {
			input  interface{}
			output string
		}{
			"nil": {
				input:  nil,
				output: "nil",
			},
			"false": {
				input:  false,
				output: "false",
			},
			"true": {
				input:  true,
				output: "true",
			},
			"123": {
				input:  123,
				output: "123",
			},
			"`hello world`": {
				input:  "hello world",
				output: `"hello world"`,
			},
		}

		for k, s := range scenarios {
			t.Run(k, func(t *testing.T) {
				output, err := exporter.Export(s.input)
				assert.NoError(t, err)
				assert.Equal(t, s.output, output)
			})
		}
	})

	t.Run("Given invalid scenarios", func(t *testing.T) {
		scenarios := map[string]struct {
			input interface{}
			error string
		}{
			"struct {}": {
				input: struct{}{},
				error: "parameter of type `struct {}` is not supported",
			},
			"*testing.T": {
				input: t,
				error: "parameter of type `*testing.T` is not supported",
			},
		}

		for k, s := range scenarios {
			t.Run(k, func(t *testing.T) {
				output, err := exporter.Export(s.input)
				assert.EqualError(t, err, s.error)
				assert.Equal(t, "", output)
			})
		}
	})
}
