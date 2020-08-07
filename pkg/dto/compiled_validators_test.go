package dto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_seekableStringSlice_contains(t *testing.T) {
	scenarios := []struct {
		slice  seekableStringSlice
		input  string
		output bool
	}{
		{
			slice:  seekableStringSlice{"foo", "bar"},
			input:  "bar",
			output: true,
		},
		{
			slice:  seekableStringSlice{"foo", "bar"},
			input:  "foobar",
			output: false,
		},
		{
			slice:  seekableStringSlice{},
			input:  "bar",
			output: false,
		},
		{
			slice:  nil,
			input:  "bar",
			output: false,
		},
	}

	for id, s := range scenarios {
		t.Run(fmt.Sprintf("scenario #%d", id), func(t *testing.T) {
			assert.Equal(
				t,
				s.output,
				s.slice.contains(s.input),
			)
		})
	}
}
