package compiler

import (
	"fmt"
	"testing"
)

func TestCompiler_handleServiceType(t *testing.T) {
	scenarios := []inputOutputScenario{
		{
			input:  "MyStruct",
			output: "MyStruct",
		},
		{
			input:  "my/import.MyStruct",
			output: "alias.MyStruct",
		},
		{
			input:  `"my/import".MyStruct`,
			output: "alias.MyStruct",
		},
	}

	doTestInputOutput(
		t,
		Compiler{
			imports: mockImports{alias: "alias"},
		}.handleServiceType,
		scenarios...,
	)
}

func TestCompiler_handleServiceValue(t *testing.T) {
	scenarios := []inputOutputScenario{
		{
			input:  "MyValue",
			output: "MyValue",
		},
		{
			input:  "my/import.MyValue",
			output: "alias.MyValue",
		},
		{
			input:  `"my/import".MyValue`,
			output: "alias.MyValue",
		},
		{
			input:  `"my/import".MyStruct{}.MyMethod`,
			output: "alias.MyStruct{}.MyMethod",
		},
	}

	doTestInputOutput(
		t,
		Compiler{
			imports: mockImports{alias: "alias"},
		}.handleServiceValue,
		scenarios...,
	)
}

func TestCompiler_handleServiceConstructor(t *testing.T) {
	scenarios := []inputOutputScenario{
		{
			input:  "my/import.NewFoo",
			output: "alias.NewFoo",
		},
		{
			input:  `"my/import".NewBar`,
			output: "alias.NewBar",
		},
		{
			input:  "NewFoo",
			output: "NewFoo",
		},
	}

	doTestInputOutput(
		t,
		Compiler{
			imports: mockImports{alias: "alias"},
		}.handleServiceConstructor,
		scenarios...,
	)
}

type mockImports struct {
	alias string
}

func (m mockImports) GetAlias(string) string {
	return m.alias
}

func (m mockImports) RegisterPrefix(shortcut string, path string) error {
	return fmt.Errorf("error")
}
