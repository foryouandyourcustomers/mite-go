package executor

import (
	"fmt"
	"github.com/progrium/go-shell"
	"os/exec"
)

type Config struct {
	execFullPath   string
	buildDirectory string
}

func Executor(execFullPath, buildDirectory string) *Config {
	shell.Shell = []string{"/bin/bash", "-c"}
	shell.Run()
	defer shell.ErrExit()

	shell.Run(fmt.Sprintf("pushd %s; go build -o %s .; popd", buildDirectory, execFullPath))

	return &Config{
		execFullPath:   execFullPath,
		buildDirectory: buildDirectory,
	}
}

func (c *Config) Execute(args ...string) ([]byte, error) {
	fmt.Println(c.execFullPath)
	cmd := exec.Command(c.execFullPath, args...)
	return cmd.CombinedOutput()
}
