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

var domainCmd = &cobra.Command{
	Use:   "domain [domain-name]",
	Short: "Perform OSINT on domains",
	Long: `Perform comprehensive OSINT on domains including:
- DNS record enumeration (A, MX, NS, TXT, CNAME)
- Subdomain discovery
- Mail server identification
- SPF/DMARC/DKIM record analysis`,
	Args: cobra.ExactArgs(1),
	RunE: runDomain,
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.Flags().Bool("enumerate-subdomains", false, "Enumerate common subdomains")
	domainCmd.Flags().Bool("dns-records", true, "Query DNS records")
}

func runDomain(cmd *cobra.Command, args []string) error {
	domain := args[0]

	// Create engine
	eng, err := engine.NewEngine(cfg)
	if err != nil {
		return fmt.Errorf("failed to create engine: %v", err)
	}
	defer eng.Shutdown()

	// Register domain module
	domainModule := osint.NewDomainModule(cfg)
	eng.RegisterModule(domainModule)

	// Execute
	log.Infof("Analyzing domain: %s", domain)
	ctx := context.Background()
	results, err := eng.Execute(ctx, "domain", domain)
	if err != nil {
		return fmt.Errorf("domain analysis failed: %v", err)
	}

	// Check for subdomain enumeration flag
	enumerateSubdomains, _ := cmd.Flags().GetBool("enumerate-subdomains")
	if enumerateSubdomains {
		log.Info("Enumerating subdomains...")
		subdomains, err := domainModule.LookupSubdomains(ctx, domain)
		if err != nil {
			log.Warnf("Subdomain enumeration failed: %v", err)
		} else {
			log.Infof("Found %d subdomains", len(subdomains))
			for _, sub := range subdomains {
				log.Infof("  - %s", sub)
			}
		}
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
