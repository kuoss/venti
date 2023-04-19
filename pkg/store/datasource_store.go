package store

import (
	"errors"
	"fmt"
	"log"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store/discovery"
)

// DatasourceStore
type DatasourceStore struct {
	config      *model.DatasourceConfig
	datasources []model.Datasource
	discoverer  discovery.Discoverer
}

// NewDatasourceStore return *DatasourceStore after service discovery (with k8s service)
func NewDatasourceStore(cfg *model.DatasourceConfig, discoverer discovery.Discoverer) (*DatasourceStore, error) {
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

// return single datasource
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

// return multiple datasources
func (s *DatasourceStore) GetDatasources() []model.Datasource {
	return s.datasources
}

func (s *DatasourceStore) GetDatasourcesWithSelector(selector model.DatasourceSelector) []model.Datasource {
	outputs := s.datasources
	outputs = filterBySystem(outputs, selector.System)
	outputs = filterByType(outputs, selector.Type)
	return outputs
}

func filterBySystem(inputs []model.Datasource, system model.DatasourceSystem) []model.Datasource {
	if system == model.DatasourceSystemNone {
		return inputs
	}
	outputs := []model.Datasource{}
	for _, input := range inputs {
		if input.IsMain == (system == model.DatasourceSystemMain) {
			outputs = append(outputs, input)
		}
	}
	return outputs
}

func filterByType(inputs []model.Datasource, typ model.DatasourceType) []model.Datasource {
	if typ == model.DatasourceTypeNone {
		return inputs
	}
	outputs := []model.Datasource{}
	for _, input := range inputs {
		if input.Type == typ {
			outputs = append(outputs, input)
		}
	}
	return outputs
}
