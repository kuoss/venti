package application

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRun_LoadGlobalConfigFileError(t *testing.T) {
	_, cleanup := testutil.SetupTest(t, map[string]string{})
	defer cleanup()

	err := new(App).Run("1.0.0")
	assert.EqualError(t, err, "failed to load configuration: loadGlobalConfigFile err: error on ReadFile: open etc/venti.yml: no such file or directory")
}

func TestRun_NewServicesError_DBOpenError(t *testing.T) {
	_, cleanup := testutil.SetupTest(t, map[string]string{
		"@/etc":                                "etc",
		"@/docs/examples/datasources.dev1.yml": "etc/datasources.yml",
	})
	defer cleanup()

	err := new(App).Run("1.0.0")
	assert.EqualError(t, err, "failed to initialize services: NewUserService err: DB open err: unable to open database file: no such file or directory")
}

func TestRun_NewServicesError_Alert(t *testing.T) {
	tempDir, cleanup := testutil.SetupTest(t, map[string]string{
		"@/etc":                                "etc",
		"@/docs/examples/datasources.dev1.yml": "etc/datasources.yml",
	})
	defer cleanup()

	invalidFileData := "kind: AlertRuleFile" +
		"\n" + "groups:" +
		"\n" + "  ..."
	err := os.WriteFile(tempDir+"/etc/alertrules/sample.yml", []byte(invalidFileData), os.ModePerm)
	assert.NoError(t, err)

	err = new(App).Run("1.0.0")
	assert.EqualError(t, err, "failed to initialize services: new alertRuleService err: loadAlertRuleFileFromFilename err: unmarshalStrict err: yaml: unmarshal errors:\n  line 3: cannot unmarshal !!str `...` into []model.RuleGroup")
}

func TestRun_SmokeTest(t *testing.T) {
	_, cleanup := testutil.SetupTest(t, map[string]string{
		"@/data":                               "data",
		"@/etc":                                "etc",
		"@/docs/examples/datasources.dev1.yml": "etc/datasources.yml",
	})
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	panicChan := make(chan interface{}, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan <- r
			}
		}()
		err := new(App).Run("1.0.0")
		assert.NoError(t, err)
		close(panicChan)
	}()

	var done bool
	select {
	case <-ctx.Done():
		done = true
	case p := <-panicChan:
		t.Fatalf("panic occurred: %v", p)
	}
	assert.True(t, done)
}
