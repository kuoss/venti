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
	assert.Equal(t, "", stdout)
	assert.Contains(t, stderr, "config load failed: error on loadDatasourceConfigFromFilepath: error on ReadFile: open etc/datasources.yaml: no such file or directory")
	assert.Error(t, err)
	assert.EqualError(t, err, "exit status 1")
}
