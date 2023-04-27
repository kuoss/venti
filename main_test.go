package main

import (
	"os"
	"testing"

	"github.com/kuoss/common/tester"
	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	err := os.Chdir("/tmp")
	assert.NoError(t, err)

	stdout, stderr, err := tester.CaptureChildTest(func() {
		main()
	})
	assert.Error(t, err)
	assert.EqualError(t, err, "exit status 1")
	assert.Equal(t, "", stdout)
	assert.Contains(t, stderr, "config load failed: error on loadDatasourceConfigFile: error on ReadFile: open etc/datasources.yml: no such file or directory")
}
