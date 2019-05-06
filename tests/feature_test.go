package tests_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/leanovate/mite-go/tests/executor"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
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
	s.Step(`^A local mock server is setup for the http method "([^"]*)" and path "([^"]*)" which returns:$`, c.aLocalMockServerIsSetupForTheHttpMethodAndPathWhichReturns)
	s.Step(`^Mite is setup to connect to this mock server$`, c.miteIsSetupToConnectToThisMockServer)
	s.Step(`^"([^"]*)" should return the following:$`, c.shouldReturnTheFollowing)
	s.Step(`^A local mock server is setup for the http method "([^"]*)" and path "([^"]*)" which expects a body of:$`, c.aLocalMockServerIsSetupForTheHttpMethodAndPathWhichExpectsABodyOf)
	s.Step(`^The mock server returns the following if the expectation is met:$`, c.theMockServerReturnsTheFollowingIfTheExpectationIsMet)
}

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

type cmdTest struct {
	executor   *executor.Config
	mockServer MockServer
}

type MockServer struct {
	Server   *httptest.Server
	Handlers []*MockHandler
}

type ReplyGenerator func() string

type MockHandler struct {
	ExpectedPath   string
	ExpectedMethod string
	ExpectedBody   string
	ReplyGenerator ReplyGenerator
}

func (c *cmdTest) reset(interface{}, error) {
	err := c.executor.Clean()
	if err != nil {
		panic(err)
	}
	c.executor = executor.Executor(buildDirectory)

	if c.mockServer.Server != nil {
		c.mockServer.Server.Close()
	}
}

func (c *cmdTest) anEmptyConfigFileCalled(arg1 string) error {
	_, err := os.Stat(filepath.Join(c.executor.ExecDirectory, arg1))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *cmdTest) iExecute(subCommand string) error {
	output, err := c.executor.Execute(subCommand)
	fmt.Println(string(output))
	return err
}

func (c *cmdTest) shouldReturn(subCommand, output string) error {
	stdout, err := c.executor.Execute(subCommand)
	if err != nil {
		return err
	}
	outputWithoutSpace := strings.TrimSpace(string(stdout))
	return assertEqual(output, outputWithoutSpace)
}

func (c *cmdTest) aLocalMockServerIsSetupForTheHttpMethodAndPathWhichReturns(arg1, arg2 string, arg3 *gherkin.DocString) error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != arg1 {
			w.WriteHeader(400)
		}
		if r.URL.Path != arg2 {
			w.WriteHeader(400)
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err := w.Write([]byte(arg3.Content))
		if err != nil {
			panic(err)
		}
	}
	handlerFunc := http.HandlerFunc(handler)
	c.mockServer.Server = httptest.NewServer(handlerFunc)
	return nil
}

func (c *cmdTest) miteIsSetupToConnectToThisMockServer() error {
	err := c.iExecute("-c .mite.toml config api.key=foo")
	if err != nil {
		return err
	}
	err = c.iExecute(fmt.Sprintf("-c .mite.toml config api.url=%s", c.mockServer.Server.URL))
	if err != nil {
		return err
	}
	return nil
}

func (c *cmdTest) shouldReturnTheFollowing(arg1 string, arg2 *gherkin.DocString) error {
	actualOutput, err := c.executor.Execute(arg1)
	if err != nil {
		return err
	}
	return assertEqual(strings.TrimSpace(arg2.Content), strings.TrimSpace(string(actualOutput)))
}

func (c *cmdTest) aLocalMockServerIsSetupForTheHttpMethodAndPathWhichExpectsABodyOf(method, path string, expectedBody *gherkin.DocString) error {
	mockHandler := MockHandler{
		ExpectedPath:   path,
		ExpectedMethod: method,
		ExpectedBody:   strings.TrimSpace(expectedBody.Content),
		ReplyGenerator: func() string {
			return ""
		},
	}

	c.mockServer.Handlers = append(c.mockServer.Handlers, &mockHandler)

	handler := func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			panic(err)
		}
		body := buf.String()
		mockHandler := c.mockServer.getMockHandlerFor(r.Method, r.URL.Path, strings.TrimSpace(body))
		if mockHandler == nil {
			w.WriteHeader(400)
			return
		}
		err = assertEqualJson(strings.TrimSpace(body), mockHandler.ExpectedBody)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		replyBody := mockHandler.ReplyGenerator()
		_, err = w.Write([]byte(replyBody))
		if err != nil {
			panic(err)
		}
	}

	handlerFunc := http.HandlerFunc(handler)
	c.mockServer.Server = httptest.NewServer(handlerFunc)
	return nil
}

func (ms MockServer) getMockHandlerFor(method, path, expectedBody string) *MockHandler {
	for _, handler := range ms.Handlers {
		if handler.ExpectedMethod == method &&
			handler.ExpectedPath == path &&
			isEqualJson(handler.ExpectedBody, expectedBody) {
			return handler
		}
	}
	return nil
}

func (c *cmdTest) theMockServerReturnsTheFollowingIfTheExpectationIsMet(replyBody *gherkin.DocString) error {
	c.mockServer.Handlers[len(c.mockServer.Handlers)-1].ReplyGenerator = func() string {
		return replyBody.Content
	}

	return nil
}

func assertEqual(expected, actual string) error {
	if strings.Compare(expected, actual) != 0 {
		return errors.New(fmt.Sprintf("expected: \n%s\n to equal: \n%s", expected, actual))
	}
	return nil
}

func assertEqualJson(s1, s2 string) error {
	var o1 interface{}
	var o2 interface{}

	if len(s1) == 0 && len(s2) == 0 {
		return nil
	}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(o1, o2) {
		return errors.New("json not deep equal")
	}

	return nil
}

func isEqualJson(s1, s2 string) bool {
	var o1 interface{}
	var o2 interface{}

	if len(s1) == 0 && len(s2) == 0 {
		return true
	}

	err := json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}

	if !reflect.DeepEqual(o1, o2) {
		return false
	}

	return true
}
