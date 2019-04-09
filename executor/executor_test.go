package executor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCliBdd(t *testing.T) {
	// given
	tmpDir, err := ioutil.TempDir("", "clibdd")
	assert.Nil(t, err)
	efp, err := filepath.Abs(filepath.Join(tmpDir, "mite"))
	assert.Nil(t, err)
	bd, err := filepath.Abs("..")
	assert.Nil(t, err)
	config := Executor(efp, bd)

	// when
	stdOutAndErr, err := config.Execute("-c", "testdata/.mite.toml", "config")

	// then
	assert.Nil(t, err)
	fmt.Println(string(stdOutAndErr))
}
