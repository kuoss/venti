package datasource

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/discovery"
	"github.com/stretchr/testify/require"
)

var (
	service *DatasourceService

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
	service, err = New(&datasourceConfig, discovery.Discoverer(nil))
	if err != nil {
		service = &DatasourceService{}
	}
}

type discovererOkMock struct{}

func (m *discovererOkMock) Do(discovery model.Discovery) ([]model.Datasource, error) {
	return []model.Datasource{}, nil
}

type discovererErrorMock struct{}

func (m *discovererErrorMock) Do(discovery model.Discovery) ([]model.Datasource, error) {
	return nil, errors.New("mock error")
}

func TestNew(t *testing.T) {
	testCases := []struct {
		cfg        *model.DatasourceConfig
		discoverer discovery.Discoverer
		want       *DatasourceService
		wantError  string
	}{
		{
			&model.DatasourceConfig{},
			discovery.Discoverer(nil),
			&DatasourceService{datasources: []model.Datasource{}},
			"",
		},
		{
			&model.DatasourceConfig{Datasources: []model.Datasource{
				{Name: "mainPrometheus", Type: model.DatasourceTypePrometheus, URL: "http://prometheus:9090", IsMain: true}}},
			discovery.Discoverer(nil),
			&DatasourceService{
				config: model.DatasourceConfig{
					Datasources: []model.Datasource{{Type: "prometheus", Name: "mainPrometheus", URL: "http://prometheus:9090", IsMain: true}},
					Discovery:   model.Discovery{Enabled: false, MainNamespace: "", AnnotationKey: "", ByNamePrometheus: false, ByNameLethe: false}},
				datasources: []model.Datasource{{Type: "prometheus", Name: "mainPrometheus", URL: "http://prometheus:9090", IsMain: true}},
				discoverer:  discovery.Discoverer(nil),
			},
			"",
		},
		{
			&model.DatasourceConfig{Datasources: []model.Datasource{}, Discovery: model.Discovery{Enabled: true}},
			&discovererOkMock{},
			&DatasourceService{
				config: model.DatasourceConfig{
					Datasources: []model.Datasource{},
					Discovery:   model.Discovery{Enabled: true}},
				datasources: []model.Datasource{},
				discoverer:  &discovererOkMock{},
			},
			"",
		},
		{
			&model.DatasourceConfig{Datasources: []model.Datasource{}, Discovery: model.Discovery{Enabled: true}},
			&discovererErrorMock{},
			nil, "load err: discoverer.Do err: mock error",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := New(tc.cfg, tc.discoverer)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.want, got)
		})
	}
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
			dss := service.GetDatasourcesWithSelector(tc.selector)
			names := []string{}
			for _, ds := range dss {
				names = append(names, ds.Name)
			}
			require.Equal(tt, tc.want, names)
		})
	}
}
