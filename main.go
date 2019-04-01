package main

import (
	"github.com/spf13/viper"
)

const configFileName = ".mite-go.toml"
const configPath = "$HOME"

func main() {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(".mite-go")
	viper.SetConfigType("toml")
	cmdLineHandler()
}
