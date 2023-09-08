package datasource

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/discovery"
	"github.com/stretchr/testify/assert"
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
			&DatasourceService{},
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

func TestReload(t *testing.T) {
	cfg := &model.DatasourceConfig{
		Datasources: []model.Datasource{
			{Type: "prometheus", Name: "mainPrometheus", URL: "http://prometheus:9090", IsMain: true},
		},
		Discovery: model.Discovery{Enabled: true},
	}
	service, err := New(cfg, &discovererOkMock{})
	require.NoError(t, err)
	require.NotZero(t, service)
	service.discoverer = &discovererErrorMock{}
	err = service.Reload()
	require.EqualError(t, err, "Reload err: discoverer.Do err: mock error")
}

func TestSetMainDatasources(t *testing.T) {
	testCases := []struct {
		datasources []model.Datasource
		want        []model.Datasource
	}{
		{ // empty
			[]model.Datasource{},
			[]model.Datasource{},
		},
		{ // 1 prom
			[]model.Datasource{{Type: model.DatasourceTypePrometheus, Name: "prom1"}},
			[]model.Datasource{{Type: model.DatasourceTypePrometheus, Name: "prom1", IsMain: true}},
		},
		{ // 1 lethe
			[]model.Datasource{{Type: model.DatasourceTypeLethe, Name: "lethe1"}},
			[]model.Datasource{{Type: model.DatasourceTypeLethe, Name: "lethe1", IsMain: true}},
		},
		{ // 2 proms
			[]model.Datasource{
				{Type: model.DatasourceTypePrometheus, Name: "prom1"},
				{Type: model.DatasourceTypePrometheus, Name: "prom2"},
			},
			[]model.Datasource{
				{Type: model.DatasourceTypePrometheus, Name: "prom1", IsMain: true},
				{Type: model.DatasourceTypePrometheus, Name: "prom2"},
			},
		},
		{ // 2 lethes
			[]model.Datasource{
				{Type: model.DatasourceTypeLethe, Name: "lethe1"},
				{Type: model.DatasourceTypeLethe, Name: "lethe2"},
			},
			[]model.Datasource{
				{Type: model.DatasourceTypeLethe, Name: "lethe1", IsMain: true},
				{Type: model.DatasourceTypeLethe, Name: "lethe2"},
			},
		},
		{ // 1 prom & 1 lethe
			[]model.Datasource{
				{Type: model.DatasourceTypePrometheus, Name: "prom1"},
				{Type: model.DatasourceTypeLethe, Name: "lethe1"},
			},
			[]model.Datasource{
				{Type: model.DatasourceTypePrometheus, Name: "prom1", IsMain: true},
				{Type: model.DatasourceTypeLethe, Name: "lethe1", IsMain: true},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			datasources = tc.datasources
			setMainDatasources(datasources)
			require.ElementsMatch(t, datasources, tc.want)
		})
	}
}

func TestGetDatasources(t *testing.T) {
	err := service.Reload()
	require.NoError(t, err)
	want := []model.Datasource{
		{Type: "prometheus", Name: "mainPrometheus", URL: "http://prometheus:9090", IsMain: true, IsDiscovered: false},
		{Type: "prometheus", Name: "subPrometheus1", URL: "http://prometheus1:9090", IsMain: false, IsDiscovered: false},
		{Type: "prometheus", Name: "subPrometheus2", URL: "http://prometheus2:9090", IsMain: false, IsDiscovered: false},
		{Type: "lethe", Name: "mainLethe", URL: "http://lethe:3100", IsMain: true, IsDiscovered: false},
		{Type: "lethe", Name: "subLethe1", URL: "http://lethe1:3100", IsMain: false, IsDiscovered: false},
		{Type: "lethe", Name: "subLethe2", URL: "http://lethe2:3100", IsMain: false, IsDiscovered: false},
	}
	got := service.GetDatasources()
	require.Equal(t, want, got)
}

func TestGetMainDatasourceByType(t *testing.T) {
	testCases := []struct {
		typ       model.DatasourceType
		want      model.Datasource
		wantError string
	}{
		{
			model.DatasourceTypeNone,
			model.Datasource{},
			"datasource of type  not found",
		},
		{
			model.DatasourceTypePrometheus,
			model.Datasource{Type: "prometheus", Name: "mainPrometheus", URL: "http://prometheus:9090", IsMain: true, IsDiscovered: false},
			"",
		},
		{
			model.DatasourceTypeLethe,
			model.Datasource{Type: "lethe", Name: "mainLethe", URL: "http://lethe:3100", IsMain: true, IsDiscovered: false},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := service.GetMainDatasourceByType(tc.typ)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetDatasourceByIndex(t *testing.T) {
	testCases := []struct {
		idx       int
		want      model.Datasource
		wantError string
	}{
		{
			-1,
			model.Datasource{},
			"datasource index[-1] not exists",
		},
		{
			0,
			model.Datasource{Type: "prometheus", Name: "mainPrometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: true, IsDiscovered: false},
			"",
		},
		{
			1,
			model.Datasource{Type: "prometheus", Name: "subPrometheus1", URL: "http://prometheus1:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
			"",
		},
		{
			2,
			model.Datasource{Type: "prometheus", Name: "subPrometheus2", URL: "http://prometheus2:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
			"",
		},
		{
			9999,
			model.Datasource{},
			"datasource index[9999] not exists",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := service.GetDatasourceByIndex(tc.idx)
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
			assert.Equal(tt, tc.want, names)
		})
	}
}
