package main

import (
	"os"

	"github.com/chxmxii/a3/internal/cli"
)

// Set via ldflags at build time.
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	rootCmd := cli.NewRootCmd()
	rootCmd.Version = version + " (" + commit + ") built " + buildTime
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
