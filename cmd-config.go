package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
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
			printFullConfig()
			return
		}

		firstArgument := args[0]
		configKey := firstArgument
		containsEquals := strings.Index(firstArgument, "=") > 0
		if containsEquals {
			// write to config
			configKeyValue := strings.Split(firstArgument, "=")
			configKey := configKeyValue[0]
			configValue := configKeyValue[1]
			viper.Set(configKey, configValue)
			dir, err := homedir.Dir()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
			err = viper.WriteConfigAs(filepath.Join(dir, configFileName))
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
			return
		}
		err := viper.ReadInConfig()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		configValue := viper.Get(configKey)
		switch v := configValue.(type) {
		case string:
			fmt.Println(v)
		default:
			_, _ = fmt.Fprintf(os.Stderr, "unknown type %v\n", reflect.TypeOf(v))
		}
	},
}

func printFullConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	c := viper.AllSettings()
	fmt.Printf("%v\n", c)
}
