package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/EmmmmDeee/One-Liner-OSINT/internal/engine"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/osint"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/output"
)

var emailCmd = &cobra.Command{
	Use:   "email [email-address]",
	Short: "Perform OSINT on email addresses",
	Long: `Perform comprehensive OSINT on email addresses including:
- Email format validation
- Domain analysis
- Breach checking (requires Have I Been Pwned API key)
- Provider identification`,
	Args: cobra.ExactArgs(1),
	RunE: runEmail,
}

func init() {
	rootCmd.AddCommand(emailCmd)
	emailCmd.Flags().Bool("check-breaches", true, "Check for data breaches")
	emailCmd.Flags().Bool("extract-domain", true, "Extract domain information")
}

func runEmail(cmd *cobra.Command, args []string) error {
	email := args[0]

	// Create engine
	eng, err := engine.NewEngine(cfg)
	if err != nil {
		return fmt.Errorf("failed to create engine: %v", err)
	}
	defer eng.Shutdown()

	// Register email module
	emailModule := osint.NewEmailModule(cfg)
	eng.RegisterModule(emailModule)

	// Execute
	log.Infof("Analyzing email: %s", email)
	ctx := context.Background()
	results, err := eng.Execute(ctx, "email", email)
	if err != nil {
		return fmt.Errorf("email analysis failed: %v", err)
	}

	// Format and display results
	formatter := output.NewFormatter(
		cfg.Output,
		os.Stdout,
		cfg.OutputConfig.Pretty,
		cfg.OutputConfig.Timestamp,
		cfg.NoColor,
	)

	if err := formatter.Format(results); err != nil {
		return fmt.Errorf("failed to format results: %v", err)
	}

	log.Infof("Analysis completed: %d results", len(results))
	return nil
}
