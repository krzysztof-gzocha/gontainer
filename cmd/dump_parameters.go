package cmd

import (
	"fmt"
	"os"

	"github.com/gomponents/gontainer/pkg"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/landoop/tableprinter"
	"github.com/spf13/cobra"
)

type mockImports struct {
}

func (m mockImports) GetAlias(i string) string {
	const max = 30
	const preffix = "(...)"
	r := []rune(i)
	if len(r) > max {
		r = r[len(r)-(max-len([]rune(preffix))):]
		i = preffix + string(r)
	}
	return "\"" + i + "\""
}

func (m mockImports) GetImports() []imports.Import {
	return nil
}

func (m mockImports) RegisterPrefix(shortcut string, path string) error {
	return nil
}

type paramRow struct {
	Name    string `header:"Name"`
	Pattern string `header:"Param"`
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

	//write := func(s string) {
	//	_, err := cmd.OutOrStdout().Write([]byte(s))
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	callback := func(cmd *cobra.Command, args []string) {
		reader := pkg.NewDefaultConfigReader(func(s string) {
			//write(fmt.Sprintf("    %s\n", s))
		})
		//write("Reading files...\n")
		input, err := reader.Read(inputFiles)
		handleErr(err, "Configuration error")

		imps := &mockImports{}
		compiler := pkg.NewDefaultCompiler(imps)

		compiledInput, ciErr := compiler.Compile(input)
		handleErr(ciErr, "Cannot build container")

		var rows []paramRow
		for _, p := range compiledInput.Params {
			rows = append(rows, paramRow{
				Name:    p.Name,
				Pattern: p.Code,
			})
		}

		p := tableprinter.New(os.Stdout)
		p.Print(rows)
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
