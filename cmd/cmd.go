package cmd

import (
	"github.com/leanovate/mite-go/config"
	"github.com/spf13/cobra"
)

type dependencies struct {
	conf config.Config
}

var deps dependencies

func HandleCommands(c config.Config) error {
	deps = dependencies{conf: c}
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "mite-go",
	Short: "cli client for mite time tracking",
	Run: func(cmd *cobra.Command, args []string) {
		// list entries for last 7 days
	},
}
