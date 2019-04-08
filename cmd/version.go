package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "prints version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", deps.version.Version, deps.version.Commit, deps.version.Date)
		return nil
	},
}
