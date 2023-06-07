package datasource

import (
	"errors"
	"fmt"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/discovery"
)

// DatasourceService
type DatasourceService struct {
	config      model.DatasourceConfig
	datasources []model.Datasource
	discoverer  discovery.Discoverer
}

// NewDatasourceService return *DatasourceService after service discovery (with k8s service)
func New(cfg *model.DatasourceConfig, discoverer discovery.Discoverer) (*DatasourceService, error) {
	service := &DatasourceService{*cfg, nil, discoverer}
	err := service.load()
	if err != nil {
		return nil, fmt.Errorf("load err: %w", err)
	}
	return service, nil
}

func (s *DatasourceService) load() error {
	// load from config
	s.datasources = append(s.datasources, s.config.Datasources...)

	// load from discovery
	if s.config.Discovery.Enabled {

		discoveredDatasources, err := s.discoverer.Do(s.config.Discovery)
		if err != nil {
			logger.Errorf("error on discoverDatasources: %s", err)
		}
		s.datasources = append(s.datasources, discoveredDatasources...)
	}
	// if len(s.datasources) < 1 {
	// 	return fmt.Errorf("no datasource")
	// }
	// set main datasources
	s.setMainDatasources()
	return nil
}

// get datasources from k8s servicesWithoutAnnotation
// recognize as a datasource by annotation or name

// ensure that there is one main datasource for each type
func (s *DatasourceService) setMainDatasources() {

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
		for i, ds := range s.datasources {
			if ds.Type == model.DatasourceTypePrometheus {
				s.datasources[i].IsMain = true
				break
			}
		}
	}

	// fallback for main lethe datasource
	// If there is no main lethe, the first lethe will be a main lethe.
	if !existsMainLethe {
		for i, ds := range s.datasources {
			if ds.Type == model.DatasourceTypeLethe {
				s.datasources[i].IsMain = true
				break
			}
		}
	}
}

// return single datasource
func (s *DatasourceService) GetMainDatasourceByType(typ model.DatasourceType) (model.Datasource, error) {
	for _, ds := range s.datasources {
		if ds.Type == typ && ds.IsMain {
			return ds, nil
		}
	}
	return model.Datasource{}, fmt.Errorf("datasource of type %s not found", typ)
}

func (s *DatasourceService) GetDatasourceByIndex(index int) (model.Datasource, error) {
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
func (s *DatasourceService) GetDatasources() []model.Datasource {
	return s.datasources
}

func (s *DatasourceService) GetDatasourcesWithSelector(selector model.DatasourceSelector) []model.Datasource {
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
