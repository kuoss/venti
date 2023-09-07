package datasource

import (
	"fmt"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/discovery"
)

// DatasourceService
type DatasourceService struct {
	loaded      bool
	config      model.DatasourceConfig
	datasources []model.Datasource
	discoverer  discovery.Discoverer
}

// NewDatasourceService return *DatasourceService after service discovery (with k8s service)
func New(cfg *model.DatasourceConfig, discoverer discovery.Discoverer) (*DatasourceService, error) {
	service := &DatasourceService{false, *cfg, nil, discoverer}
	err := service.load()
	if err != nil {
		return nil, fmt.Errorf("load err: %w", err)
	}
	return service, nil
}

func (s *DatasourceService) load() error {
	// load from config
	datasources := s.config.Datasources

	// load from discovery
	if s.config.Discovery.Enabled {

		discoveredDatasources, err := s.discoverer.Do(s.config.Discovery)
		if err != nil {
			return fmt.Errorf("discoverer.Do err: %w", err)
		}
		datasources = append(datasources, discoveredDatasources...)
	}
	setMainDatasources(datasources)
	s.datasources = datasources
	return nil
}

func (s *DatasourceService) Reload() error {
	if err := s.load(); err != nil {
		return fmt.Errorf("Reload err: %w", err)
	}
	return nil
}

// get datasources from k8s servicesWithoutAnnotation
// recognize as a datasource by annotation or name

// ensure that there is one main datasource for each type
func setMainDatasources(datasources []model.Datasource) {

	existsMainPrometheus := false
	existsMainLethe := false

	for _, ds := range datasources {
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
		for i, ds := range datasources {
			if ds.Type == model.DatasourceTypePrometheus {
				datasources[i].IsMain = true
				break
			}
		}
	}

	// fallback for main lethe datasource
	// If there is no main lethe, the first lethe will be a main lethe.
	if !existsMainLethe {
		for i, ds := range datasources {
			if ds.Type == model.DatasourceTypeLethe {
				datasources[i].IsMain = true
				break
			}
		}
	}
}

// return deep copied datasources
func (s *DatasourceService) getDatasources() []model.Datasource {
	datasources := []model.Datasource{}
	datasources = append(datasources, s.datasources...)
	return datasources
}

// return single datasource
func (s *DatasourceService) GetMainDatasourceByType(typ model.DatasourceType) (model.Datasource, error) {
	for _, ds := range s.getDatasources() {
		if ds.Type == typ && ds.IsMain {
			return ds, nil
		}
	}
	return model.Datasource{}, fmt.Errorf("datasource of type %s not found", typ)
}

func (s *DatasourceService) GetDatasourceByIndex(index int) (model.Datasource, error) {
	datasources := s.getDatasources()
	if index < 0 || index >= len(datasources) {
		return model.Datasource{}, fmt.Errorf("datasource index[%d] not exists", index)
	}
	return datasources[index], nil
}

// return multiple datasources
func (s *DatasourceService) GetDatasources() []model.Datasource {
	return s.getDatasources()
}

func (s *DatasourceService) GetDatasourcesWithSelector(selector model.DatasourceSelector) []model.Datasource {
	outputs := s.getDatasources()
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
