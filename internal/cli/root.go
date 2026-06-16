package cli

import (
	"github.com/spf13/cobra"
)

var dbPath string

// NewRootCmd creates the root 3A command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "3a",
		Short: "3A - Agnostic Account Assessment",
		Long:  "3A assesses cloud accounts (AWS/OCI) using Steampipe for resource discovery, evaluates security posture, estimates costs, and generates reports.",
	}

	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "", "path to SQLite database (default: ~/.3a/3a.db)")

	rootCmd.AddCommand(newAssessCmd())
	rootCmd.AddCommand(newProfilesCmd())
	rootCmd.AddCommand(newReportCmd())

	return rootCmd
}

func getDBPath() string {
	if dbPath != "" {
		return dbPath
	}
	return "~/.3a/3a.db"
}
