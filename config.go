package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func configGetApiUrl() string {
	return configGet("api.url")
}

func configGetApiKey() string {
	return configGet("api.key")
}

func configGet(key string) string {
	err := viper.ReadInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	return viper.GetString(key)
}

func configSet(key string, value string) {
	err := viper.ReadInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	viper.Set(key, value)
	dir, err := homedir.Dir()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	err = viper.MergeInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	err = viper.WriteConfigAs(filepath.Join(dir, configFileName))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func configPrintAll() {
	err := viper.ReadInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	c := viper.AllSettings()
	fmt.Printf("%v\n", c)
}
