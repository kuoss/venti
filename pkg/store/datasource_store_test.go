package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/store/discovery"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	datasourceConfig   *model.DatasourceConfig
	datasources        []model.Datasource
	datasourcePointers []*model.Datasource
	store              *DatasourceStore
)

func init() {
	datasources = []model.Datasource{
		{Name: "mainPrometheus", Type: model.DatasourceTypePrometheus, URL: "http://prometheus:9090", IsMain: true},
		{Name: "subPrometheus1", Type: model.DatasourceTypePrometheus, URL: "http://prometheus1:9090", IsMain: false},
		{Name: "subPrometheus2", Type: model.DatasourceTypePrometheus, URL: "http://prometheus2:9090", IsMain: false},
		{Name: "mainLethe", Type: model.DatasourceTypeLethe, URL: "http://lethe:3100", IsMain: true},
		{Name: "subLethe1", Type: model.DatasourceTypeLethe, URL: "http://lethe1:3100", IsMain: false},
		{Name: "subLethe2", Type: model.DatasourceTypeLethe, URL: "http://lethe2:3100", IsMain: false},
	}
	for i := range datasources {
		datasourcePointers = append(datasourcePointers, &datasources[i])
	}
	datasourceConfig = &model.DatasourceConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  datasourcePointers,
		Discovery: model.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}
	var discoverer discovery.Discoverer
	var err error
	store, err = NewDatasourceStore(datasourceConfig, discoverer)
	if err != nil {
		store = &DatasourceStore{}
	}
}

func TestNewDatasourceStore(t *testing.T) {
	assert.Equal(t, store.config, datasourceConfig)
	assert.ElementsMatch(t, store.datasources, datasources)
}

func TestGetDatasourcesWithSelector(t *testing.T) {
	testCases := []struct {
		selector model.DatasourceSelector
		want     []string
	}{
		{
			model.DatasourceSelector{},
			[]string{"mainPrometheus", "subPrometheus1", "subPrometheus2", "mainLethe", "subLethe1", "subLethe2"},
		},
		{
			model.DatasourceSelector{System: model.DatasourceSystemMain},
			[]string{"mainPrometheus", "mainLethe"},
		},
		{
			model.DatasourceSelector{System: model.DatasourceSystemSub},
			[]string{"subPrometheus1", "subPrometheus2", "subLethe1", "subLethe2"},
		},
		{
			model.DatasourceSelector{Type: model.DatasourceTypeLethe},
			[]string{"mainLethe", "subLethe1", "subLethe2"},
		},
		{
			model.DatasourceSelector{Type: model.DatasourceTypePrometheus},
			[]string{"mainPrometheus", "subPrometheus1", "subPrometheus2"},
		},
		{
			model.DatasourceSelector{System: model.DatasourceSystemMain, Type: model.DatasourceTypeLethe},
			[]string{"mainLethe"},
		},
		{
			model.DatasourceSelector{System: model.DatasourceSystemMain, Type: model.DatasourceTypePrometheus},
			[]string{"mainPrometheus"},
		},
		{
			model.DatasourceSelector{System: model.DatasourceSystemSub, Type: model.DatasourceTypeLethe},
			[]string{"subLethe1", "subLethe2"},
		},
		{
			model.DatasourceSelector{System: model.DatasourceSystemSub, Type: model.DatasourceTypePrometheus},
			[]string{"subPrometheus1", "subPrometheus2"},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d %v", i, tc.selector), func(tt *testing.T) {
			dss := store.GetDatasourcesWithSelector(tc.selector)
			names := []string{}
			for _, ds := range dss {
				names = append(names, ds.Name)
			}
			assert.Equal(tt, tc.want, names)
		})
	}
}
