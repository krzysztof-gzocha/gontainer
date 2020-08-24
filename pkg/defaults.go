package pkg

import (
	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/compiler"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/tokens"
)

const gontainerHelperPath = "github.com/gomponents/gontainer-helpers"

func NewDefaultCompiler(imports imports.Imports) *compiler.Compiler {
	tokenizer := tokens.NewPatternTokenizer(
		[]tokens.TokenFactoryStrategy{
			tokens.NewTokenSimpleFunction(imports, "env", "os", "Getenv"),
			tokens.NewTokenSimpleFunction(imports, "envInt", gontainerHelperPath+"/env", "MustGetInt"),
			tokens.NewTokenSimpleFunction(imports, "todo", gontainerHelperPath+"/std", "GetMissingParameter"),
			tokens.TokenPercentSign{},
			tokens.TokenReference{},
			tokens.TokenString{},
		},
		imports,
	)

	paramResolver := parameters.NewSimpleResolver(
		tokenizer,
		exporters.NewDefaultExporter(),
		imports,
	)

	return compiler.NewCompiler(
		input.NewDefaultValidator(),
		compiled.NewDefaultValidator(),
		imports,
		tokenizer,
		paramResolver,
		arguments.NewDefaultResolver(paramResolver),
	)
}
