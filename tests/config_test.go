package tests_test

import (
	"errors"
	"flag"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/leanovate/mite-go/tests/executor"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const buildDirectory = "../"

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	c := cmdTest{
		executor: executor.Executor(buildDirectory),
	}

	s.AfterScenario(c.reset)
	s.AfterSuite(func() {
		err := c.executor.Clean()
		if err != nil {
			panic(err)
		}
	})
	s.Step(`^an empty config file called "([^"]*)"$`, c.anEmptyConfigFileCalled)
	s.Step(`^I execute "([^"]*)"$`, c.iExecute)
	s.Step(`^"([^"]*)" should return "([^"]*)"$`, c.shouldReturn)
}

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

type cmdTest struct {
	executor *executor.Config
}

func (c *cmdTest) reset(interface{}, error) {
	err := c.executor.Clean()
	if err != nil {
		panic(err)
	}
	c.executor = executor.Executor(buildDirectory)
}

func (c *cmdTest) anEmptyConfigFileCalled(arg1 string) error {
	_, err := os.Stat(filepath.Join(c.executor.ExecDirectory, arg1))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *cmdTest) iExecute(subCommand string) error {
	subCmd := strings.Split(subCommand, " ")
	_, err := c.executor.Execute(subCmd...)
	return err
}

func (c *cmdTest) shouldReturn(subCommand, output string) error {
	subCmd := strings.Split(subCommand, " ")
	stdout, err := c.executor.Execute(subCmd...)
	outputWithoutSpace := strings.TrimSpace(string(stdout))
	if strings.Compare(outputWithoutSpace, output) != 0 {
		return errors.New("output is not expected output")
	}
	return err
}
