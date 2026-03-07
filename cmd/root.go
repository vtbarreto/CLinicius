package cmd

import (
	"github.com/spf13/cobra"
)

// LangFlag holds the value of the --lang persistent flag.
// Subcommands read it to initialize the correct localizer.
var LangFlag string

var rootCmd = &cobra.Command{
	Use:   "clinicius",
	Short: "Architectural governance CLI for Go projects",
	Long: `CLinicius is a static analysis tool that enforces architectural integrity
in Go codebases by validating layer boundaries and detecting cyclic dependencies.

It reads a clinicius.yaml config file that declares layers and their forbidden
dependencies, then checks your packages against those rules.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&LangFlag, "lang", "",
		`Language for output messages (e.g. en-US, pt-BR).
Overrides CLINICIUS_LANG and LANG environment variables.`,
	)
}
