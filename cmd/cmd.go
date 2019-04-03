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
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
