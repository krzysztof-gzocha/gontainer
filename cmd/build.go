package cmd

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg"
	"github.com/gomponents/gontainer/pkg/template2"
	"io/ioutil"

	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/dto"
	"github.com/gomponents/gontainer/pkg/exporters"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/syntax"
	"github.com/gomponents/gontainer/pkg/tokens"
	"github.com/spf13/cobra"
)

func NewBuildCmd() *cobra.Command {
	var (
		inputFiles []string
		outputFile string
		cmd        *cobra.Command
	)

	handleErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	write := func(s string) {
		_, err := cmd.OutOrStdout().Write([]byte(s))
		if err != nil {
			panic(err)
		}
	}

	callback := func(cmd *cobra.Command, args []string) {
		input := dto.Input{}

		reader := pkg.NewDefaultConfigReader(func(s string) {
			write(fmt.Sprintf("Parsing file: `%s`\n", s))
		})
		input, err := reader.Read(inputFiles)
		handleErr(err)

		imps := imports.NewSimpleImports("_gontainer")
		compiler := createCompiler(imps)

		compiledInput, ciErr := compiler.Compile(input)
		handleErr(ciErr)

		if tpl, err := template2.NewSimpleBuilder(imps).Build(compiledInput); err != nil {
			panic(fmt.Sprintf("Unexpected error has occurred during building container: %s", err.Error()))
		} else {
			if fileErr := ioutil.WriteFile(outputFile, []byte(tpl), 0644); fileErr != nil {
				panic(fmt.Sprintf("Error has occurred during saving file: %s", fileErr.Error()))
			}

			_, _ = cmd.OutOrStdout().Write([]byte(fmt.Sprintf("Container has built [%s]\n", outputFile)))
		}
	}

	cmd = &cobra.Command{
		Use:   "build",
		Short: "build-container",
		Long:  "build-long",
		Run:   callback,
	}

	cmd.Flags().StringArrayVarP(&inputFiles, "input", "i", nil, "input name (required)")
	_ = cmd.MarkFlagRequired("input")

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "output name (required)")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func createCompiler(imports imports.Imports) *dto.BaseCompiler {
	tokenizer := tokens.NewPatternTokenizer([]tokens.TokenFactoryStrategy{
		tokens.NewTokenSimpleFunction(imports, "env", "os", "Getenv"),
		tokens.NewTokenSimpleFunction(imports, "envInt", "github.com/gomponents/gontainer/pkg/std", "MustGetIntEnv"),
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
