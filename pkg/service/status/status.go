package status

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/kuoss/venti/pkg/config"
)

type StatusService struct {
	buildInfo   BuildInfo
	runtimeInfo RuntimeInfo
}

func New(cfg *config.Config) (*StatusService, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getwd err: %w", err)
	}
	return &StatusService{
		buildInfo: BuildInfo{
			Version:   cfg.AppInfo.Version,
			Revision:  "(TBD)",
			Branch:    "(TBD)",
			BuildUser: "(TBD)",
			BuildDate: "(TBD)",
			GoVersion: runtime.Version(),
		},
		runtimeInfo: RuntimeInfo{
			StartTime:           time.Now(),
			CWD:                 cwd,
			ReloadConfigSuccess: true,
			LastConfigTime:      time.Now(),
			CorruptionCount:     -1,
			GoroutineCount:      -1,
			GOMAXPROCS:          runtime.GOMAXPROCS(0),
			GOMEMLIMIT:          -1,
			GOGC:                os.Getenv("GOGC"),
			GODEBUG:             os.Getenv("GODEBUG"),
			StorageRetention:    "N/A",
		},
	}, nil
}

func (s *StatusService) BuildInfo() BuildInfo {
	return s.buildInfo
}

func (s *StatusService) RuntimeInfo() RuntimeInfo {
	s.runtimeInfo.GOMEMLIMIT = debug.SetMemoryLimit(-1)
	s.runtimeInfo.GoroutineCount = runtime.NumGoroutine()
	return s.runtimeInfo
}
