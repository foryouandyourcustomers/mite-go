package main

import (
	"fmt"
	"github.com/leanovate/mite-go/cmd"
	"github.com/leanovate/mite-go/config"
	"github.com/leanovate/mite-go/mite"
	"github.com/mitchellh/go-homedir"
	"os"
)

const configFileName = ".mite-go"
const configType = "toml"

func main() {
	homeDirectory, err := homedir.Dir()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	c := config.NewConfig(configFileName, homeDirectory, configType)
	api := mite.NewApi(c.GetApiUrl(), c.GetApiKey())

	err = cmd.HandleCommands(c, api)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
