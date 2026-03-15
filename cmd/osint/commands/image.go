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

var imageCmd = &cobra.Command{
	Use:   "image [image-path-or-url]",
	Short: "Extract metadata from images",
	Long: `Extract EXIF metadata from images including:
- Camera information (make, model)
- Timestamp information
- GPS coordinates (if available)
- Software used
- Image dimensions`,
	Args: cobra.ExactArgs(1),
	RunE: runImage,
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.Flags().Bool("extract-gps", true, "Extract GPS coordinates")
	imageCmd.Flags().Bool("camera-info", true, "Extract camera information")
}

func runImage(cmd *cobra.Command, args []string) error {
	imagePath := args[0]

	// Create engine
	eng, err := engine.NewEngine(cfg)
	if err != nil {
		return fmt.Errorf("failed to create engine: %v", err)
	}
	defer eng.Shutdown()

	// Register image module
	imageModule := osint.NewImageModule(cfg)
	eng.RegisterModule(imageModule)

	// Execute
	log.Infof("Analyzing image: %s", imagePath)
	ctx := context.Background()
	results, err := eng.Execute(ctx, "image", imagePath)
	if err != nil {
		return fmt.Errorf("image analysis failed: %v", err)
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
