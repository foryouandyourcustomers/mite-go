package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(completionCommand)
}

var completionCommand = &cobra.Command{
	Use:   "completion",
	Short: "generates a bash autocomplete script",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletion(os.Stdout)
	},
}
