package config

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config interface {
	GetApiUrl() string
	GetApiKey() string
	GetActivity(activity string) Activity
	GetDisplayLocation() *time.Location
	GetVacation() Activity
	Get(key string) string
	Set(key string, value string)
	PrintAll()
}

type config struct {
	fileFullPath string
}

type Activity struct {
	ProjectId string
	ServiceId string
}

func NewConfig(fullPath string) Config {
	createConfigFileIfNonExistent(fullPath)
	viper.SetConfigFile(fullPath)

	return &config{fileFullPath: fullPath}
}

func createConfigFileIfNonExistent(ffp string) {
	if _, err := os.Stat(ffp); os.IsNotExist(err) {
		_, err := os.Create(ffp)
		if err != nil {
			panic(fmt.Sprintf("could not create config file: %s\n", err))
		}
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

func (c *config) GetVacation() Activity {
	projectId := c.Get("vacation.projectId")
	serviceId := c.Get("vacation.serviceId")
	return Activity{
		ProjectId: projectId,
		ServiceId: serviceId,
	}
}

func (c *config) GetDisplayLocation() *time.Location {
	s := c.Get("display.location")
	if s == "" {
		return time.Local
	}

	loc, err := time.LoadLocation(s)
	if err != nil {
		return time.Local
	}

	return loc
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
	file, err := os.Open(c.fileFullPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		panic(err)
	}
}
