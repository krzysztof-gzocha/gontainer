package pkg

import (
	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/dto"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/syntax"
	"github.com/gomponents/gontainer/pkg/tokens"
)

func NewDefaultCompiler(imports imports.Imports) dto.Compiler {
	tokenizer := tokens.NewPatternTokenizer([]tokens.TokenFactoryStrategy{
		tokens.NewTokenSimpleFunction(imports, "env", "os", "Getenv"),
		tokens.NewTokenSimpleFunction(imports, "envInt", "github.com/gomponents/gontainer-helpers/env", "MustGetInt"),
		tokens.NewTokenSimpleFunction(imports, "todo", "github.com/gomponents/gontainer-helpers/std", "MustGetMissingParameter"),
		tokens.TokenPercentSign{},
		tokens.TokenReference{},
		tokens.TokenString{},
	})

	exporter := exporters.NewDefaultExporter()

	argumentResolver := arguments.NewSimpleResolver(
		[]arguments.Subresolver{
			arguments.NewServiceResolver(imports, syntax.NewSimpleServiceResolver(syntax.NewSimpleTypeResolver(imports))),
			arguments.NewPatternResolver(
				tokenizer,
				exporter,
				imports,
			),
		},
	)

	return dto.NewBaseCompiler(
		imports,
		dto.NewDefaultValidator(),
		dto.NewDefaultCompiledValidator(),
		parameters.NewSimpleBagFactory(tokenizer, exporter, imports),
		syntax.NewSimpleTypeResolver(imports),
		argumentResolver,
		syntax.NewSimpleFunctionResolver(imports),
	)
}
