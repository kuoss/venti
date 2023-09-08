package main

import (
	"os"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/assert"
)

// The main function doesn't exit when the router runs,
// so we only test for errors here.
func Test_main(t *testing.T) {
	err := os.Chdir("./docs")
	assert.NoError(t, err)

	stdout, stderr, err := tester.CaptureChildTest(func() {
		main()
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "exit status 1")

	assert.Equal(t, "", stdout)

	assert.Contains(t, stderr, `level=info msg="loading configurations..."`)
	assert.Contains(t, stderr, `level=info msg="loading datasource config file: etc/datasources.yml"`)
	assert.Contains(t, stderr, `level=fatal msg="config.Load err: loadDatasourceConfigFile err: error on ReadFile: open etc/datasources.yml: no such file or directory"`)
}
