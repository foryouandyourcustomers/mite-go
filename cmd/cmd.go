package cmd

import (
	"fmt"
	"github.com/leanovate/mite-go/app"
	"github.com/spf13/cobra"
	"path/filepath"
)

var application *app.Application

func HandleCommands() error {
	rootCmd.Flags().BoolP("version", "v", false, "prints version")
	rootCmd.PersistentFlags().StringP("config", "c", "", "alternative config file location")

	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "mite",
	Short: "cli client for mite time tracking",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var fullConfigPath string
		configArg, err := cmd.Flags().GetString("config")
		if err == nil && configArg != "" {
			fullConfigPath, err = filepath.Abs(configArg)
			if err != nil {
				return err
			}
		}

		application, err = app.NewApplication(fullConfigPath)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		printShortVersion, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		}

		if printShortVersion {
			fmt.Printf("%s\n", application.Version)
		}
		return nil
	},
}
