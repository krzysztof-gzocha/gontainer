package syntax

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/stretchr/testify/assert"
)

type importsMock struct {
}

func (i importsMock) GetAlias(string) string {
	return "alias"
}

func (i importsMock) GetImports() []imports.Import {
	return nil
}

func TestSimpleFunctionResolver_ResolveFunction(t *testing.T) {
	resolver := NewSimpleFunctionResolver(importsMock{})

	scenarios := []struct {
		input  string
		output string
		err    bool
	}{
		{
			input:  "foo/bar.NewFooBar",
			output: "alias.NewFooBar",
			err:    false,
		},
		{
			input:  "NewFooBar",
			output: "NewFooBar",
			err:    false,
		},
		{
			input:  "!@#",
			output: "",
			err:    true,
		},
	}

	for id, s := range scenarios {
		t.Run(fmt.Sprintf("scenario #%d", id), func(t *testing.T) {
			o, err := resolver.ResolveFunction(s.input)
			if s.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, s.err)
			}
			assert.Equal(t, s.output, o)
		})
	}
}

func TestSimpleTypeResolver_ResolveType(t *testing.T) {
	resolver := NewSimpleTypeResolver(importsMock{})

	scenarios := []struct {
		input  string
		output string
		err    bool
	}{
		{
			input:  "foo/bar.MyType",
			output: "alias.MyType",
			err:    false,
		},
		{
			input:  "*foo/bar.MyType",
			output: "*alias.MyType",
			err:    false,
		},
		{
			input:  "MyType",
			output: "MyType",
			err:    false,
		},
		{
			input:  "*MyType",
			output: "*MyType",
			err:    false,
		},
		{
			input:  "!@#",
			output: "",
			err:    true,
		},
	}

	for id, s := range scenarios {
		t.Run(fmt.Sprintf("scenario #%d", id), func(t *testing.T) {
			o, err := resolver.ResolveType(s.input)
			if s.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, s.err)
			}
			assert.Equal(t, s.output, o)
		})
	}
}

func TestSimpleServiceResolver_ResolveService(t *testing.T) {
	resolver := NewSimpleServiceResolver(
		NewSimpleTypeResolver(importsMock{}),
	)

	scenarios := []struct {
		input   string
		service string
		type_   string
		err     bool
	}{
		{
			input:   "@service",
			service: "service",
			type_:   "",
			err:     false,
		},
		{
			input:   "@myDb.(*my/import/db.Mysql)",
			service: "myDb",
			type_:   "*alias.Mysql",
			err:     false,
		},
		{
			input:   "@myDb.(my/import/db.Mysql)",
			service: "myDb",
			type_:   "alias.Mysql",
			err:     false,
		},
		{
			input:   "@123service",
			service: "",
			type_:   "",
			err:     true,
		},
	}

	for id, s := range scenarios {
		t.Run(fmt.Sprintf("scenario #%d", id), func(t *testing.T) {
			service, type_, err := resolver.ResolveService(s.input)
			if s.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, s.err)
			}
			assert.Equal(t, s.service, service)
			assert.Equal(t, s.type_, type_)
		})
	}
}
