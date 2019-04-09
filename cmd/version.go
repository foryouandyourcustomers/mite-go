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
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", application.Version, application.Commit, application.Date)
		return nil
	},
}
