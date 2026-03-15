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

var ipCmd = &cobra.Command{
	Use:   "ip [ip-address]",
	Short: "Perform OSINT on IP addresses",
	Long: `Perform comprehensive OSINT on IP addresses including:
- Basic IP information (IPv4/IPv6, private/public)
- Reverse DNS lookup
- Geolocation lookup
- ASN information (if available)`,
	Args: cobra.ExactArgs(1),
	RunE: runIP,
}

func init() {
	rootCmd.AddCommand(ipCmd)
	ipCmd.Flags().Bool("geolocation", true, "Perform geolocation lookup")
	ipCmd.Flags().Bool("reverse-dns", true, "Perform reverse DNS lookup")
}

func runIP(cmd *cobra.Command, args []string) error {
	ip := args[0]

	// Create engine
	eng, err := engine.NewEngine(cfg)
	if err != nil {
		return fmt.Errorf("failed to create engine: %v", err)
	}
	defer eng.Shutdown()

	// Register IP module
	ipModule := osint.NewIPModule(cfg)
	eng.RegisterModule(ipModule)

	// Execute
	log.Infof("Analyzing IP: %s", ip)
	ctx := context.Background()
	results, err := eng.Execute(ctx, "ip", ip)
	if err != nil {
		return fmt.Errorf("IP analysis failed: %v", err)
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
