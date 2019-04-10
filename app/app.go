package app

import (
	"fmt"
	"github.com/leanovate/mite-go/config"
	"github.com/leanovate/mite-go/mite"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

type Application struct {
	Conf    config.Config
	MiteApi mite.Api
	Version string
	Commit  string
	Date    string
}

const defaultConfigFileName = ".mite.toml"

// these flags will be overwritten during the build process by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "none"
)

func NewApplication(fullConfigPath string) (*Application, error) {
	homeDirectory, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	if fullConfigPath == "" {
		fullConfigPath = filepath.Join(homeDirectory, defaultConfigFileName)
	}

	c := config.NewConfig(fullConfigPath)
	api, err := mite.NewApi(c.GetApiUrl(), c.GetApiKey(), version)
	if err != nil {
		return nil, err
	}

	if c.GetApiUrl() == "" {
		_, _ = fmt.Fprintln(os.Stderr, "please configure your API url by executing: 'mite config api.url=<your mite api url>'")
	}

	if c.GetApiKey() == "" {
		_, _ = fmt.Fprintln(os.Stderr, "please configure your API key by executing: 'mite config api.key=<your mite api key>'")
	}
	return &Application{Conf: c, MiteApi: api, Version: version, Commit: commit, Date: date}, nil
}
