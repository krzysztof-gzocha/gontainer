package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/definition"
	"github.com/gomponents/gontainer/pkg/exporters"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
	"github.com/gomponents/gontainer/pkg/template"
	"github.com/gomponents/gontainer/pkg/tokens"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewBuildCmd() *cobra.Command {
	var (
		inputFile  string
		outputFile string
		cmd        *cobra.Command
	)

	callback := func(cmd *cobra.Command, args []string) {
		input := struct {
			Meta struct {
				Pkg string `yaml:"pkg"`
			} `yaml:"meta"`
			Params   parameters.RawParameters `yaml:"parameters"`
			Services definition.Services      `yaml:"services"`
		}{}

		if file, err := ioutil.ReadFile(inputFile); err != nil {
			panic(fmt.Sprintf("Error has occurred during opening file: %s", err.Error()))
		} else {
			if yamlErr := yaml.Unmarshal(file, &input); yamlErr != nil {
				panic(fmt.Sprintf("Error has occurred during parsing yaml file: %s", yamlErr.Error()))
			}
		}

		imps := imports.NewSimpleImports("gontainer_")
		exporter := exporters.NewDefaultExporter()

		tokenFactories := []tokens.TokenFactoryStrategy{
			tokens.NewTokenSimpleFunction(imps, "env", "os", "Getenv"),
			tokens.TokenPercentSign{},
			tokens.TokenReference{},
			tokens.TokenString{},
		}

		tokenizer := tokens.NewPatternTokenizer(tokenFactories)

		bagFactory := parameters.NewSimpleBagFactory(tokenizer, exporter, imps)
		// todo
		resolvedParams, _ := bagFactory.Create(input.Params)

		argumentResolver := arguments.NewSimpleResolver(
			[]arguments.Subresolver{
				arguments.ServiceResolver{},
				arguments.NewPatternResolver(
					tokenizer,
					exporter,
					imps,
					resolvedParams,
				),
			},
		)

		services := make(map[string]template.Service)
		for id, s := range input.Services {
			compiledArgs := make([]arguments.Argument, 0)
			for _, a := range s.Args {
				ra, err := argumentResolver.Resolve(a)
				if err != nil {
					panic(err)
				}
				compiledArgs = append(compiledArgs, ra)
			}
			service := template.Service{}
			service.Service = s
			service.CompiledArgs = compiledArgs
			services[id] = service
		}

		tplBuilder, _ := template.NewSimpleBuilder(
			template.WithPackage(input.Meta.Pkg),
			template.WithImports(imps),
			template.WithParams(resolvedParams),
			template.WithServices(services),
		)

		if tpl, err := tplBuilder.Build(); err != nil {
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

	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "input name (required)")
	_ = cmd.MarkFlagRequired("input")

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "output name (required)")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
