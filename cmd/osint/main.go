package main

import (
	"os"

	"github.com/EmmmmDeee/One-Liner-OSINT/cmd/osint/commands"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Execute root command
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
