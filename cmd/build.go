package cmd

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/template2"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

func NewBuildCmd() *cobra.Command {
	var (
		inputFiles []string
		outputFile string
		cmd        *cobra.Command
	)

	writeErr := func(s string) {
		_, err := cmd.OutOrStderr().Write([]byte(s))
		if err != nil {
			panic(err)
		}
	}

	handleErr := func(err error, msg string) {
		if err != nil {
			writeErr(fmt.Sprintf("%s: %s\n", msg, err.Error()))
			os.Exit(1)
		}
	}

	write := func(s string) {
		_, err := cmd.OutOrStdout().Write([]byte(s))
		if err != nil {
			panic(err)
		}
	}

	callback := func(cmd *cobra.Command, args []string) {
		reader := pkg.NewDefaultConfigReader(func(s string) {
			write(fmt.Sprintf("\t%s\n", s))
		})
		write("Reading files...\n")
		input, err := reader.Read(inputFiles)
		handleErr(err, "Configuration error")

		imps := imports.NewSimpleImports("_gontainer")
		compiler := pkg.NewDefaultCompiler(imps)

		compiledInput, ciErr := compiler.Compile(input)
		handleErr(ciErr, "Cannot build container")

		tpl, tplErr := template2.NewSimpleBuilder(imps).Build(compiledInput)
		handleErr(tplErr, "Unexpected error has occurred during building container")
		fileErr := ioutil.WriteFile(outputFile, []byte(tpl), 0644)
		handleErr(fileErr, "Error has occurred during saving file")
		write(fmt.Sprintf("Container has been built\n\t%s\n", outputFile))
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
