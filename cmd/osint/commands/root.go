package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/config"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/logger"
)

var (
	cfgFile string
	log     = logger.NewLogger()
	cfg     *config.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "osint",
	Short: "Production-grade OSINT intelligence gathering system",
	Long: `One-Liner-OSINT: A comprehensive, production-grade Go system for
Open-Source Intelligence (OSINT) gathering. Provides powerful tools for
domain reconnaissance, social media analysis, metadata extraction, and more.

Built with efficiency, scalability, and reliability for security researchers,
bug bounty hunters, and ethical hackers.`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.osint.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().String("output", "text", "output format (text, json, csv)")
	rootCmd.PersistentFlags().Int("workers", 10, "number of concurrent workers")
	rootCmd.PersistentFlags().Int("timeout", 30, "request timeout in seconds")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("workers", rootCmd.PersistentFlags().Lookup("workers"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".osint" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".osint")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Load configuration
	cfg = config.LoadConfig()
}
