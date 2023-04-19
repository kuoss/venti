package mock

import (
	"time"

	"github.com/kuoss/venti/pkg/model"
)

func DatasourceConfigFromDatasources(datasources []model.Datasource) *model.DatasourceConfig {
	var datasourcePointers []*model.Datasource
	for i := range datasources {
		datasourcePointers = append(datasourcePointers, &datasources[i])
	}
	return &model.DatasourceConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  datasourcePointers,
		Discovery:    model.Discovery{},
	}
}
