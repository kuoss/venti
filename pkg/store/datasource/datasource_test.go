package datasource

import (
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	store *DatasourceStore

	datasourceConfig = model.DatasourceConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  datasources,
		Discovery: model.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}

	datasources = []model.Datasource{
		{Name: "mainPrometheus", Type: model.DatasourceTypePrometheus, URL: "http://prometheus:9090", IsMain: true},
		{Name: "subPrometheus1", Type: model.DatasourceTypePrometheus, URL: "http://prometheus1:9090", IsMain: false},
		{Name: "subPrometheus2", Type: model.DatasourceTypePrometheus, URL: "http://prometheus2:9090", IsMain: false},
		{Name: "mainLethe", Type: model.DatasourceTypeLethe, URL: "http://lethe:3100", IsMain: true},
		{Name: "subLethe1", Type: model.DatasourceTypeLethe, URL: "http://lethe1:3100", IsMain: false},
		{Name: "subLethe2", Type: model.DatasourceTypeLethe, URL: "http://lethe2:3100", IsMain: false},
	}
)

func init() {
	var err error
	store, err = New(&datasourceConfig, discovery.Discoverer(nil))
	if err != nil {
		store = &DatasourceStore{}
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		cfg  *model.DatasourceConfig
		want *DatasourceStore
	}{
		{
			&model.DatasourceConfig{},
			&DatasourceStore{},
		},
		{
			&model.DatasourceConfig{Datasources: []model.Datasource{
				{Name: "mainPrometheus", Type: model.DatasourceTypePrometheus, URL: "http://prometheus:9090", IsMain: true}}},
			&DatasourceStore{config: model.DatasourceConfig{QueryTimeout: 0, Datasources: []model.Datasource{
				{Type: "prometheus", Name: "mainPrometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: true, IsDiscovered: false}}, Discovery: model.Discovery{Enabled: false, MainNamespace: "", AnnotationKey: "", ByNamePrometheus: false, ByNameLethe: false}}, datasources: []model.Datasource{model.Datasource{Type: "prometheus", Name: "mainPrometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: true, IsDiscovered: false}}, discoverer: discovery.Discoverer(nil)},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := New(tc.cfg, discovery.Discoverer(nil))
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
	// assert.Equal(t, store.config, datasourceConfig)
	// assert.ElementsMatch(t, store.datasources, datasources)
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
