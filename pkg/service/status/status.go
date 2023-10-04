package status

import (
	"runtime"

	"github.com/kuoss/venti/pkg/model"
	apiV1 "github.com/prometheus/prometheus/web/api/v1"
)

type StatusService struct {
	ventiVersion model.VentiVersion
}

func New(cfg *model.Config) *StatusService {
	return &StatusService{
		ventiVersion: model.VentiVersion{
			Version:   cfg.Version,
			GoVersion: runtime.Version(),
		},
	}
}

func (s *StatusService) BuildInfo() model.VentiVersion {
	return s.ventiVersion
}

func (s *StatusService) RuntimeInfo() apiV1.RuntimeInfo {
	return apiV1.RuntimeInfo{}
}
