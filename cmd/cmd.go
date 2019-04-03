package cmd

import (
	"github.com/leanovate/mite-go/config"
	"github.com/leanovate/mite-go/mite"
	"github.com/spf13/cobra"
)

type dependencies struct {
	conf    config.Config
	miteApi mite.MiteApi
}

var deps dependencies

func HandleCommands(c config.Config, m mite.MiteApi) error {
	deps = dependencies{conf: c, miteApi: m}
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "mite-go",
	Short: "cli client for mite time tracking",
	Run: func(cmd *cobra.Command, args []string) {
		// list entries for last 7 days
	},
}
