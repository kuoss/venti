package status

import (
	"runtime"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var store *StatusStore

func init() {
	store = New(&model.Config{
		Version: "test",
	})
}

func TestNew(t *testing.T) {
	assert.NotEmpty(t, store)

	assert.Equal(t, "test", store.ventiVersion.Version)
	assert.Equal(t, runtime.Version(), store.ventiVersion.GoVersion)
}

func TestBuildInfo(t *testing.T) {
	got := store.BuildInfo()
	assert.Equal(t, "test", got.Version)
	assert.Equal(t, runtime.Version(), got.GoVersion)
}
