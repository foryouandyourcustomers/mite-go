package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(configCommand)
}

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "sets or reads a config property",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			configPrintAll()
			return
		}

		firstArgument := args[0]
		configKey := firstArgument
		containsEquals := strings.Index(firstArgument, "=") > 0
		err := viper.ReadInConfig()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		if containsEquals {
			// write to config
			configKeyValue := strings.Split(firstArgument, "=")
			configKey := configKeyValue[0]
			configValue := configKeyValue[1]
			configSet(configKey, configValue)
			return
		}
		fmt.Println(configGet(configKey))
	},
}
