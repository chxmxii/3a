package cli

import (
	"github.com/spf13/cobra"
)

var dbPath string

// NewRootCmd creates the root 3A command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "a3",
		Short: "3A - Agnostic Account Assessment",
		Long:  "3A assesses cloud accounts (AWS/OCI) using Steampipe for resource discovery, evaluates security posture, estimates costs, and generates reports.",
	}

	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "", "path to SQLite database (default: ~/.a3/a3.db)")

	rootCmd.AddCommand(newAssessCmd())
	rootCmd.AddCommand(newProfilesCmd())
	rootCmd.AddCommand(newReportCmd())
	rootCmd.AddCommand(newConfigureCmd())

	return rootCmd
}

func getDBPath() string {
	if dbPath != "" {
		return dbPath
	}
	return "~/.a3/a3.db"
}
