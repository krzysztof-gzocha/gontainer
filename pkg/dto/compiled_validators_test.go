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

func Test_validateCircularDependency(t *testing.T) {
	scenarios := []struct {
		services []CompiledService
		error    string
	}{
		{
			services: []CompiledService{
				{
					Name: "foo",
					Args: []CompiledArg{
						{ServiceLink: &ServiceLink{Name: "bar"}},
					},
				},
				{
					Name: "bar",
					Args: []CompiledArg{
						{ServiceLink: &ServiceLink{Name: "foo"}},
					},
				},
			},
			error: "found circular dependency bar -> foo -> bar",
		},
		{
			services: []CompiledService{
				{
					Name: "s1",
					Args: []CompiledArg{
						{ServiceLink: &ServiceLink{Name: "s2"}},
					},
				},
				{
					Name: "s2",
					Args: []CompiledArg{
						{ServiceLink: &ServiceLink{Name: "s3"}},
					},
				},
				{
					Name: "s3",
					Args: []CompiledArg{
						{ServiceLink: &ServiceLink{Name: "s4"}},
					},
				},
				{
					Name: "s4",
					Args: []CompiledArg{
						{ServiceLink: &ServiceLink{Name: "s2"}},
					},
				},
			},
			error: "found circular dependency s1 -> s2 -> s3 -> s4 -> s2",
		},
		{
			services: []CompiledService{
				{
					Name: "foo",
					Args: []CompiledArg{
						{ServiceLink: &ServiceLink{Name: "foo"}},
					},
				},
			},
			error: "found circular dependency foo -> foo",
		},
		{
			services: []CompiledService{
				{
					Name: "foo",
					Args: []CompiledArg{
						{ServiceLink: &ServiceLink{Name: "bar"}},
					},
				},
			},
			error: "",
		},
	}

	for id, s := range scenarios {
		t.Run(fmt.Sprintf("scenario #%d", id), func(t *testing.T) {
			err := validateCircularDependency(CompiledInput{Services: s.services})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, s.error)
		})
	}
}
