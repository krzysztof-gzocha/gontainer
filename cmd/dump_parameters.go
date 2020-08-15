package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gomponents/gontainer/pkg"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/spf13/cobra"
)

// todo https://github.com/olekukonko/tablewriter

type mockImports struct {
}

func (m mockImports) GetAlias(i string) string {
	// todo limit to X characters
	return "\"" + i + "\""
}

func (m mockImports) GetImports() []imports.Import {
	return nil
}

func (m mockImports) RegisterPrefix(shortcut string, path string) error {
	return nil
}

func NewDumpParamsCmd() *cobra.Command {
	var (
		inputFiles []string
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
			write(fmt.Sprintf("    %s\n", s))
		})
		write("Reading files...\n")
		input, err := reader.Read(inputFiles)
		handleErr(err, "Configuration error")

		imps := mockImports{}
		compiler := pkg.NewDefaultCompiler(imps)

		compiledInput, ciErr := compiler.Compile(input)
		handleErr(ciErr, "Cannot build container")

		write("Params\n")
		max := 0
		for n, _ := range compiledInput.Params {
			if len(n) > max {
				max = len(n)
			}
		}

		for n, p := range compiledInput.Params {
			spaces := strings.Repeat(" ", max-len(n))
			//raw := p.Raw
			//if strings.Contains(raw, "\n") {
			//	raw, _ = exporters.Export(raw)
			//}
			write(fmt.Sprintf("    %s: %s%s\n", n, spaces, p.Code))
		}
	}

	cmd = &cobra.Command{
		Use:   "dump-params",
		Short: "Dump parameters",
		Long:  "todo",
		Run:   callback,
	}

	cmd.Flags().StringArrayVarP(&inputFiles, "input", "i", nil, "input name (required)")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}
