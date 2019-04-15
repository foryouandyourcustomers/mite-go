package executor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	ExecDirectory  string
	ExecFullPath   string
	buildDirectory string
}

func Executor(buildDirectory string) *Config {
	execDirectory, err := ioutil.TempDir("", "executor")
	if err != nil {
		panic(err)
	}
	execFullPath := filepath.Join(execDirectory, "mite")

	args := []string{"bash", "-c", fmt.Sprintf("pushd %s; go build -o %s .; popd", buildDirectory, execFullPath)}
	cmd := exec.Command("/usr/bin/env", args...)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	return &Config{
		ExecDirectory:  execDirectory,
		ExecFullPath:   execFullPath,
		buildDirectory: buildDirectory,
	}
}

func (c *Config) Execute(args string) ([]byte, error) {
	subCmd := strings.Split(args, " ")
	cmd := exec.Command(c.ExecFullPath, subCmd...)
	cmd.Dir = c.ExecDirectory
	return cmd.Output()
}

func (c *Config) Clean() error {
	if !strings.HasPrefix("/tmp/", c.ExecDirectory) {
		return os.RemoveAll(c.ExecDirectory)
	}
	return errors.New("tried to remove file without prefix /tmp")
}
