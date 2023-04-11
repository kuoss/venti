package store

import (
	"errors"
	"fmt"
	"github.com/kuoss/venti/pkg/store/discovery"
	"log"

	"github.com/kuoss/venti/pkg/model"
)

// DatasourceStore
type DatasourceStore struct {
	config      *model.DatasourcesConfig
	datasources []model.Datasource
	discoverer  discovery.Discoverer
}

// NewDatasourceStore return *DatasourceStore after service discovery (with k8s service)
func NewDatasourceStore(cfg *model.DatasourcesConfig, discoverer discovery.Discoverer) (*DatasourceStore, error) {
	store := &DatasourceStore{cfg, nil, discoverer}
	err := store.load()
	if err != nil {
		return nil, fmt.Errorf("error on load: %w", err)
	}
	return store, nil
}

func (s *DatasourceStore) load() error {
	// load from config
	for _, datasource := range s.config.Datasources {
		s.datasources = append(s.datasources, *datasource)
	}

	// load from discovery
	if s.config.Discovery.Enabled {

		discoveredDatasources, err := s.discoverer.Do(s.config.Discovery)
		if err != nil {
			log.Fatalf("error on discoverDatasources: %s", err)
		}
		s.datasources = append(s.datasources, discoveredDatasources...)
	}
	// set main datasources
	s.setMainDatasources()

	return nil
}

// get datasources from k8s servicesWithoutAnnotation
// recognize as a datasource by annotation or name

// ensure that there is one main datasource for each type
func (s *DatasourceStore) setMainDatasources() {

	existsMainPrometheus := false
	existsMainLethe := false

	for _, ds := range s.datasources {
		if !ds.IsMain {
			continue
		}
		switch ds.Type {
		case model.DatasourceTypePrometheus:
			existsMainPrometheus = true
			continue
		case model.DatasourceTypeLethe:
			existsMainLethe = true
			continue
		}
	}

	// fallback for main prometheus datasource
	// If there is no main prometheus, the first prometheus will be a main prometheus.
	if !existsMainPrometheus {
		for _, ds := range s.datasources {
			if ds.Type == model.DatasourceTypePrometheus {
				ds.IsMain = true
				break
			}
		}
	}

	// fallback for main lethe datasource
	// If there is no main lethe, the first lethe will be a main lethe.
	if !existsMainLethe {
		for _, ds := range s.datasources {
			if ds.Type == model.DatasourceTypeLethe {
				ds.IsMain = true
				break
			}
		}
	}
}

func (s *DatasourceStore) GetDatasources() []model.Datasource {
	return s.datasources
}

func (s *DatasourceStore) GetMainDatasourceByType(typ model.DatasourceType) (model.Datasource, error) {
	for _, ds := range s.datasources {
		if ds.Type == typ && ds.IsMain {
			return ds, nil
		}
	}
	return model.Datasource{}, fmt.Errorf("datasource of type %s not found", typ)
}

func (s *DatasourceStore) GetDatasourceByIndex(index int) (model.Datasource, error) {
	cnt := len(s.datasources)
	if cnt < 1 {
		return model.Datasource{}, errors.New("no datasource")
	}
	if index >= len(s.datasources) {
		return model.Datasource{}, fmt.Errorf("datasource index[%d] not exists", index)
	}
	return s.datasources[index], nil
}
