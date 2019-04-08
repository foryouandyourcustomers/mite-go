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

// these flags will be overwritten during the build process by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "none"
)

func main() {
	v := cmd.Version{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
	homeDirectory, err := homedir.Dir()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	c := config.NewConfig(configFileName, homeDirectory, configType)
	api := mite.NewApi(c.GetApiUrl(), c.GetApiKey(), v.Version)

	if c.GetApiUrl() == "" {
		_, _ = fmt.Fprintln(os.Stderr, "please configure your API url by executing: 'mite config api.url=<your mite api url>'")
	}

	if c.GetApiKey() == "" {
		_, _ = fmt.Fprintln(os.Stderr, "please configure your API key by executing: 'mite config api.key=<your mite api key>'")
	}

	err = cmd.HandleCommands(c, api, v)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
