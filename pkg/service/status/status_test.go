package status

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/config"
	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/require"
)

var (
	service1     *StatusService
	goVersion1   string      = runtime.Version()
	buildInfo1   BuildInfo   = BuildInfo{Version: "", Revision: "(TBD)", Branch: "(TBD)", BuildUser: "(TBD)", BuildDate: "(TBD)", GoVersion: goVersion1}
	gomaxprocs1  int         = runtime.GOMAXPROCS(0)
	runtimeInfo1 RuntimeInfo = RuntimeInfo{StartTime: time.Time{}, CWD: "", ReloadConfigSuccess: true, LastConfigTime: time.Time{}, CorruptionCount: -1, GoroutineCount: -1, GOMAXPROCS: gomaxprocs1, GOMEMLIMIT: -1, GOGC: "", GODEBUG: "", StorageRetention: "N/A"}
)

func init() {
	cwd, _ := os.Getwd()
	runtimeInfo1.CWD = cwd

	var err error
	service1, err = New(&config.Config{
		AppInfo:          model.AppInfo{Version: "test"},
		GlobalConfig:     model.GlobalConfig{},
		DatasourceConfig: model.DatasourceConfig{},
		UserConfig:       model.UserConfig{},
		AlertingConfig:   model.AlertingConfig{},
	})
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	buildInfo101 := buildInfo1
	buildInfo101.Version = "hello"

	testCases := []struct {
		cfg  *config.Config
		want *StatusService
	}{
		{
			&config.Config{},
			&StatusService{
				buildInfo:   buildInfo1,
				runtimeInfo: runtimeInfo1,
			},
		},
		{
			&config.Config{AppInfo: model.AppInfo{Version: "hello"}},
			&StatusService{
				buildInfo:   buildInfo101,
				runtimeInfo: runtimeInfo1,
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := New(tc.cfg)
			require.NoError(t, err)

			tc.want.buildInfo.BuildDate = got.buildInfo.BuildDate
			tc.want.runtimeInfo.StartTime = got.runtimeInfo.StartTime
			tc.want.runtimeInfo.LastConfigTime = got.runtimeInfo.LastConfigTime
			require.Equal(t, tc.want, got)
		})
	}
}

func TestBuildInfo(t *testing.T) {
	got := service1.BuildInfo()
	require.Equal(t, "test", got.Version)
	require.Equal(t, goVersion1, got.GoVersion)
}

func TestRuntimeInfo(t *testing.T) {
	got := service1.RuntimeInfo()

	runtimeInfo101 := runtimeInfo1
	runtimeInfo101.StartTime = got.StartTime
	runtimeInfo101.LastConfigTime = got.LastConfigTime
	runtimeInfo101.GoroutineCount = got.GoroutineCount
	runtimeInfo101.GOMEMLIMIT = got.GOMEMLIMIT
	require.Equal(t, runtimeInfo101, got)
}
