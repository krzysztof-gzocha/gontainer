package main

import (
	"fmt"
	"os"

	"github.com/gomponents/gontainer/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "gontainer",
		Short: "Gontainer is a missing DI container builder for GO",
		Long: `Gontainer allows you to build DI container based on provided configuration.
Re-use dependencies whenever you need and forget about dependency hell in main.go.`,
		Version: "0.1",
	}

	rootCmd.AddCommand(
		cmd.NewBuildCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
