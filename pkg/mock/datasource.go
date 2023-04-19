package mock

import "github.com/kuoss/venti/pkg/model"

func Datasources() []model.Datasource {
	return []model.Datasource{
		{Name: "mainPrometheus", Type: model.DatasourceTypePrometheus, URL: "http://prometheus:9090", IsMain: true},
		{Name: "subPrometheus1", Type: model.DatasourceTypePrometheus, URL: "http://prometheus1:9090", IsMain: false},
		{Name: "subPrometheus2", Type: model.DatasourceTypePrometheus, URL: "http://prometheus2:9090", IsMain: false},
		{Name: "mainLethe", Type: model.DatasourceTypeLethe, URL: "http://lethe:3100", IsMain: true},
		{Name: "subLethe1", Type: model.DatasourceTypeLethe, URL: "http://lethe1:3100", IsMain: false},
		{Name: "subLethe2", Type: model.DatasourceTypeLethe, URL: "http://lethe2:3100", IsMain: false},
	}
}
