package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

func NewVersionCmd(version, commit, date string) *cobra.Command {
	var (
		shortened  = false
		output     = "json"
		versionCmd = &cobra.Command{
			Use:   "version",
			Short: "Version will output the current build information",
			Long:  ``,
			Run: func(_ *cobra.Command, _ []string) {
				resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
				fmt.Print(resp)
				return
			},
		}
	)

	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
	versionCmd.Flags().StringVarP(&output, "output", "o", "json", "Output format. One of 'yaml' or 'json'.")

	return versionCmd
}
