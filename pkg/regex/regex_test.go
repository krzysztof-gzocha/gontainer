package regex

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	scenarios := []struct {
		regex  string
		input  string
		match  bool
		result map[string]string
	}{
		{
			regex:  "^(?P<firstname>[A-Z][a-z]+) (?P<lastname>[A-Z][a-z]+)$",
			input:  "Mary Jane",
			match:  true,
			result: map[string]string{"firstname": "Mary", "lastname": "Jane"},
		},
		{
			regex:  "^(?P<fullname>(?P<firstname>[A-Z][a-z]+) (?P<lastname>[A-Z][a-z]+))$",
			input:  "Mary Jane",
			match:  true,
			result: map[string]string{"firstname": "Mary", "lastname": "Jane", "fullname": "Mary Jane"},
		},
		{
			regex:  "^(?P<firstname>[A-Z][a-z]+) (?P<lastname>[A-Z][a-z]+)$",
			input:  "Mary Jane-Jane",
			match:  false,
			result: nil,
		},
	}

	for id, s := range scenarios {
		t.Run(fmt.Sprintf("scenario #%d", id), func(t *testing.T) {
			m, r := Match(regexp.MustCompile(s.regex), s.input)
			assert.Equal(t, s.match, m)
			assert.Equal(t, s.result, r)
		})
	}
}
