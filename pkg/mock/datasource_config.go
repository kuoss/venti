package mock

import (
	"time"

	"github.com/kuoss/venti/pkg/model"
)

func DatasourceConfigFromDatasources(datasources []model.Datasource) *model.DatasourceConfig {
	return &model.DatasourceConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  datasources,
		Discovery:    model.Discovery{},
	}
}
