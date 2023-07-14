package alerting

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/mocker"
	"github.com/kuoss/venti/pkg/mocker/alertmanager"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/discovery"
	"github.com/stretchr/testify/require"
)

var (
	alertmanagerMock  *mocker.Server
	datasourceService *datasource.DatasourceService
	ruleFiles         = []model.RuleFile{{
		Kind:               "AlertRuleFile",
		CommonLabels:       map[string]string{"rulefile": "sample-v3", "severity": "silence"},
		DatasourceSelector: model.DatasourceSelector{System: "", Type: "prometheus"},
		RuleGroups: []model.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{
				{Record: "", Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Record: "", Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Record: "", Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
			}}}}}
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	var err error
	alertmanagerMock, err = alertmanager.New(0)
	if err != nil {
		panic(err)
	}

	err = os.Chdir("../../..")
	if err != nil {
		panic(err)
	}
	err = setDatasourceService()
	if err != nil {
		panic(err)
	}
	fmt.Printf("datasourceService=%#v\n", datasourceService)
}

func shutdown() {
	alertmanagerMock.Close()
}

func setDatasourceService() error {
	datasourceConfig := model.DatasourceConfig{
		Datasources: []model.Datasource{
			{Name: "mainPrometheus", Type: model.DatasourceTypePrometheus, URL: "http://prometheus:9090", IsMain: true},
			{Name: "subPrometheus1", Type: model.DatasourceTypePrometheus, URL: "http://prometheus1:9090", IsMain: false},
			{Name: "subPrometheus2", Type: model.DatasourceTypePrometheus, URL: "http://prometheus2:9090", IsMain: false},
			{Name: "mainLethe", Type: model.DatasourceTypeLethe, URL: "http://lethe:3100", IsMain: true},
			{Name: "subLethe1", Type: model.DatasourceTypeLethe, URL: "http://lethe1:3100", IsMain: false},
			{Name: "subLethe2", Type: model.DatasourceTypeLethe, URL: "http://lethe2:3100", IsMain: false},
		},
		Discovery: model.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}
	var err error
	datasourceService, err = datasource.New(&datasourceConfig, discovery.Discoverer(nil))
	if err != nil {
		return fmt.Errorf("datasource.New err: %w", err)
	}
	return nil
}

func TestNew(t *testing.T) {
	testCases := []struct {
		file string
		want model.AlertingFile
	}{
		{
			file: "",
			want: model.AlertingFile{Alertings: []model.Alerting{{Name: "alertmanager", Type: "alertmanager", URL: "http://localhost:9093"}}},
		},
		{
			file: "asdf",
			want: model.AlertingFile{Alertings: []model.Alerting(nil)},
		},
		{
			file: "etc/alerting.yml",
			want: model.AlertingFile{Alertings: []model.Alerting{{Name: "alertmanager", Type: "alertmanager", URL: "http://localhost:9093"}}},
		},
		{
			file: "etc/alerting.yaml",
			want: model.AlertingFile{Alertings: []model.Alerting(nil)},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := New(tc.file, ruleFiles, datasourceService)
			require.NotEmpty(t, got)
			require.NotEmpty(t, got.AlertFiles)
			require.Equal(t, tc.want, got.AlertingFile)
		})
	}

}

func TestLoadAlertingFile(t *testing.T) {
	testCases := []struct {
		file      string
		want      *model.AlertingFile
		wantError string
	}{
		{
			"",
			&model.AlertingFile{Alertings: []model.Alerting{{Name: "alertmanager", Type: model.AlertingTypeAlertmanager, URL: "http://localhost:9093"}}},
			"",
		},
		{
			"asdfasdf",
			&model.AlertingFile{Alertings: []model.Alerting(nil)},
			"readFile err: open asdfasdf: no such file or directory",
		},
		{
			"etc/alerting.yml",
			&model.AlertingFile{Alertings: []model.Alerting{{Name: "alertmanager", Type: model.AlertingTypeAlertmanager, URL: "http://localhost:9093"}}},
			"",
		},
		{
			"etc/alerting.yaml",
			&model.AlertingFile{Alertings: []model.Alerting(nil)},
			"readFile err: open etc/alerting.yaml: no such file or directory",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := loadAlertingFile(tc.file)
			fmt.Println("tc.wantError=", tc.wantError)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetAlertmanagerURL(t *testing.T) {
	testCases := []struct {
		file string
		want string
	}{
		{
			"",
			"http://localhost:9093",
		},
		{
			"asdf",
			"",
		},
		{
			"etc/alerting.yml",
			"http://localhost:9093",
		},
		{
			"etc/alerting.yaml",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			service := New(tc.file, ruleFiles, datasourceService)
			require.Equal(t, tc.want, service.GetAlertmanagerURL())
		})
	}
}

func TestSendTestAlert(t *testing.T) {
	service := New("etc/alerting.yml", ruleFiles, datasourceService)
	service.AlertingFile.Alertings[0].URL = alertmanagerMock.URL
	err := service.SendTestAlert()
	require.NoError(t, err)
}
