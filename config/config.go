package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

type Config interface {
	GetApiUrl() string
	GetApiKey() string
	GetActivity(activity string) Activity
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

type Activity struct {
	ProjectId string
	ServiceId string
}

func NewConfig(fileName, filePath, fileType string) Config {
	viper.AddConfigPath("$HOME")
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)
	ffp := fullConfigPath(filePath, fileName, fileType)
	createConfigFileIfNonExistent(ffp)

	return &config{fileName: fileName, filePath: filePath, fileType: fileType, fileFullPath: ffp}
}

func fullConfigPath(filePath string, fileName string, fileType string) string {
	ffp := fmt.Sprintf("%s/%s.%s", filePath, fileName, fileType)
	if runtime.GOOS == "windows" {
		ffp = fmt.Sprintf("%s\\%s.%s", filePath, fileName, fileType)
	}
	return ffp
}

func createConfigFileIfNonExistent(ffp string) {
	if _, err := os.Stat(ffp); os.IsExist(err) {
		return
	}

	_, err := os.Create(ffp)
	if err != nil {
		panic(fmt.Sprintf("config file does not exists and wasn't able to create it: %s\n", err))
	}
}

func (c *config) GetApiUrl() string {
	return c.Get("api.url")
}

func (c *config) GetApiKey() string {
	return c.Get("api.key")
}

func (c *config) GetActivity(activity string) Activity {
	projectId := c.Get(fmt.Sprintf("activity.%s.projectId", activity))
	serviceId := c.Get(fmt.Sprintf("activity.%s.serviceId", activity))
	return Activity{
		ProjectId: projectId,
		ServiceId: serviceId,
	}
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
