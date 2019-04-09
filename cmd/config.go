package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	rootCmd.AddCommand(configCommand)
}

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "sets or reads a config property",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			application.Conf.PrintAll()
			return nil
		}

		firstArgument := args[0]
		configKey := firstArgument
		containsEquals := strings.Index(firstArgument, "=") > 0
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
		if containsEquals {
			// write listTo config
			configKeyValue := strings.Split(firstArgument, "=")
			configKey := configKeyValue[0]
			configValue := configKeyValue[1]
			application.Conf.Set(configKey, configValue)
			return nil
		}
		fmt.Println(application.Conf.Get(configKey))
		return nil
	},
}
