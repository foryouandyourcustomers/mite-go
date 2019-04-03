package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config interface {
	GetApiUrl() string
	GetApiKey() string
	Get(key string) string
	Set(key string, value string)
	PrintAll()
}

type config struct {
	fileName     string
	filePath     string
	fileType     string
	fileFullPath string
}

func NewConfig(fileName, filePath, fileType string) Config {
	viper.AddConfigPath("$HOME")
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)
	ffp := fmt.Sprintf("%s/%s.%s", filePath, fileName, fileType)
	return &config{fileName: fileName, filePath: filePath, fileType: fileType, fileFullPath: ffp}
}

func (c *config) GetApiUrl() string {
	return c.Get("api.url")
}

func (c *config) GetApiKey() string {
	return c.Get("api.key")
}

func (c *config) Get(key string) string {
	err := viper.ReadInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	return viper.GetString(key)
}

func (c *config) Set(key string, value string) {
	err := viper.ReadInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	viper.Set(key, value)
	err = viper.MergeInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	err = viper.WriteConfigAs(c.fileFullPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func (c *config) PrintAll() {
	err := viper.ReadInConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	wholeConfig := viper.AllSettings()
	fmt.Printf("%v\n", wholeConfig)
}
