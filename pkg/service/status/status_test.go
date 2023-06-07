package status

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	service1  *StatusService
	goVersion string = runtime.Version()
)

func init() {
	service1 = New(&model.Config{
		Version: "test",
	})
}

func TestNew(t *testing.T) {
	testCases := []struct {
		cfg  *model.Config
		want *StatusService
	}{
		{
			&model.Config{},
			&StatusService{ventiVersion: model.VentiVersion{Version: "", GoVersion: goVersion}},
		},
		{
			&model.Config{Version: "hello"},
			&StatusService{ventiVersion: model.VentiVersion{Version: "hello", GoVersion: goVersion}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := New(tc.cfg)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestBuildInfo(t *testing.T) {
	got := service1.BuildInfo()
	assert.Equal(t, "test", got.Version)
	assert.Equal(t, goVersion, got.GoVersion)
}
