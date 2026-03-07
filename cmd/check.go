package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vtbarreto/CLinicius/internal/analyzer"
	"github.com/vtbarreto/CLinicius/internal/config"
	"github.com/vtbarreto/CLinicius/internal/i18n"
	"github.com/vtbarreto/CLinicius/internal/reporter"
	"github.com/vtbarreto/CLinicius/internal/rules"
)

var (
	configFile string
	ciMode     bool
	jsonMode   bool
)

var checkCmd = &cobra.Command{
	Use:   "check [patterns...]",
	Short: "Check Go packages for architectural violations",
	Long: `Loads the specified Go packages (e.g. ./...), builds a dependency graph,
and runs all configured rules against it.

Examples:
  clinicius check ./...
  clinicius check ./... --ci
  clinicius check ./... --json
  clinicius check ./... --config path/to/clinicius.yaml`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&configFile, "config", "c", "clinicius.yaml", "Path to the configuration file")
	checkCmd.Flags().BoolVar(&ciMode, "ci", false, "Exit with code 1 if any violations are found (CI mode)")
	checkCmd.Flags().BoolVar(&jsonMode, "json", false, "Output violations as JSON")
}

func runCheck(_ *cobra.Command, args []string) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	graph, err := analyzer.LoadPackages(args)
	if err != nil {
		return fmt.Errorf("analyzing packages: %w", err)
	}

	engine := rules.NewEngine(
		rules.NewLayerRule(),
		rules.NewCycleRule(),
	)

	violations := engine.Run(graph, cfg)

	if jsonMode {
		r := reporter.NewJSONReporter(os.Stdout)
		if err := r.Report(violations); err != nil {
			return err
		}
	} else {
		loc := i18n.NewLocalizer(LangFlag)
		reporter.NewConsoleReporter(os.Stdout, loc).Report(violations)
	}

	if ciMode && len(violations) > 0 {
		os.Exit(1)
	}

	return nil
}
