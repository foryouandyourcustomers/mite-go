package cmd

import (
	"fmt"
	"github.com/leanovate/mite-go/config"
	"github.com/leanovate/mite-go/mite"
	"github.com/spf13/cobra"
)

type dependencies struct {
	conf    config.Config
	miteApi mite.Api
	version Version
}

type Version struct {
	Version string
	Commit  string
	Date    string
}

var deps dependencies

func HandleCommands(c config.Config, m mite.Api, v Version) error {
	deps = dependencies{conf: c, miteApi: m, version: v}
	rootCmd.Flags().BoolP("version", "v", false, "prints version")
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "mite",
	Short: "cli client for mite time tracking",
	RunE: func(cmd *cobra.Command, args []string) error {
		printShortVersion, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		}

		if printShortVersion {
			fmt.Printf("%s\n", deps.version.Version)
		}
		return nil
	},
}
