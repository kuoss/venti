package status

import (
	"runtime"

	"github.com/kuoss/venti/pkg/model"
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
