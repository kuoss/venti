package status

import (
	"runtime"

	"github.com/kuoss/venti/pkg/model"
)

type StatusStore struct {
	ventiVersion model.VentiVersion
}

func New(cfg *model.Config) *StatusStore {
	return &StatusStore{
		ventiVersion: model.VentiVersion{
			Version:   cfg.Version,
			GoVersion: runtime.Version(),
		},
	}
}

func (s *StatusStore) BuildInfo() model.VentiVersion {
	return s.ventiVersion
}
