package main

import (
	"fmt"
	"os"

	"github.com/gomponents/gontainer/cmd"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "gontainer",
		Short: "Gontainer is a missing DI container builder for GO",
		Long: `Gontainer allows you to build DI container based on provided configuration.
Re-use dependencies whenever you need and forget about dependency hell in main.go.`,
		Version: version,
	}

	rootCmd.AddCommand(
		cmd.NewBuildCmd(),
		cmd.NewDumpParamsCmd(),
		cmd.NewVersionCmd(version, commit, date),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
