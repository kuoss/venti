package config

import (
	"os"
	"os/exec"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := os.Chdir("../..")
	if err != nil {
		panic(err)
	}
}

func TestLoad(t *testing.T) {
	_ = exec.Command("cp", "docs/examples/datasources.dev1.yml", "etc/datasources.yml").Run()
	defer func() {
		os.RemoveAll("etc/datasources.yml")
	}()
	cfg, err := Load("Unknown")
	assert.NoError(t, err)
	assert.Equal(t, cfg.AppInfo.Version, "Unknown")
	assert.Equal(t, []model.Datasource{
		{Type: model.DatasourceTypePrometheus, Name: "prometheus", URL: "http://localhost:9090"},
		{Type: model.DatasourceTypeLethe, Name: "lethe", URL: "http://localhost:6060"},
	}, cfg.DatasourceConfig.Datasources)
	assert.Equal(t, model.UserConfig{EtcUsers: []model.EtcUser{
		{Username: "admin", Hash: "$2a$12$VcCDgh2NDk07JGN0rjGbM.Ad41qVR/YFJcgHp0UGns5JDymv..TOG", IsAdmin: true},
	}}, cfg.UserConfig)
}

func TestLoadGlobalConfigFile(t *testing.T) {
	testCases := []struct {
		file      string
		want      model.GlobalConfig
		wantError string
	}{
		{
			"",
			model.GlobalConfig{LogLevel: ""},
			"error on ReadFile: open : no such file or directory",
		},
		{
			"etc/datasources.yml",
			model.GlobalConfig{LogLevel: ""},
			//"error on UnmarshalStrict: yaml: unmarshal errors:\n  line 1: field datasources not found in type model.GlobalConfig",
			"error on ReadFile: open etc/datasources.yml: no such file or directory",
		},
		{
			"etc/venti.yml",
			model.GlobalConfig{GinMode: "release", LogLevel: "info"},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := loadGlobalConfigFile(tc.file)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}

}

func TestLoadDatasourceConfigFile(t *testing.T) {
	testCases := []struct {
		file      string
		want      *model.DatasourceConfig
		wantError string
	}{
		{
			"",
			nil,
			"error on ReadFile: open : no such file or directory",
		},
		{
			"docs/examples/datasources.dev1.yml",
			&model.DatasourceConfig{
				QueryTimeout: 30000000000,
				Datasources: []model.Datasource{
					{Type: "prometheus", Name: "prometheus", URL: "http://localhost:9090"},
					{Type: "lethe", Name: "lethe", URL: "http://localhost:6060"},
				},
				Discovery: model.Discovery{AnnotationKey: "kuoss.org/datasource-type"},
			},
			"",
		},
		{
			"docs/examples/datasources.dev2.yml",
			&model.DatasourceConfig{
				QueryTimeout: 30000000000,
				Datasources: []model.Datasource{
					{Type: "prometheus", Name: "prometheus1", URL: "http://vs-prometheus-server"},
					{Type: "prometheus", Name: "prometheus2", URL: "http://vs-prometheus-server"},
					{Type: "lethe", Name: "lethe", URL: "http://vs-lethe"},
				},
				Discovery: model.Discovery{AnnotationKey: "kuoss.org/datasource-type"},
			},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := loadDatasourceConfigFile(tc.file)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestLoadUserConfigFile(t *testing.T) {
	testCases := []struct {
		file      string
		want      *model.UserConfig
		wantError string
	}{
		{
			"",
			nil,
			"error on ReadFile: open : no such file or directory",
		},
		{
			"etc/users.yml",
			&model.UserConfig{EtcUsers: []model.EtcUser{
				{Username: "admin", Hash: "$2a$12$VcCDgh2NDk07JGN0rjGbM.Ad41qVR/YFJcgHp0UGns5JDymv..TOG", IsAdmin: true},
			}},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := loadUserConfigFile(tc.file)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestLoadAlertingConfigFile(t *testing.T) {
	testCases := []struct {
		file      string
		want      model.AlertingConfig
		wantError string
	}{
		{
			"",
			model.AlertingConfig{},
			"error on ReadFile: open : no such file or directory",
		},
		{
			"etc/alerting.yml",
			model.AlertingConfig{
				EvaluationInterval:  5000000000,
				AlertRelabelConfigs: nil,
				AlertmanagerConfigs: model.AlertmanagerConfigs{
					{StaticConfig: []*model.TargetGroup{
						{Targets: []string{"http://vs-alertmanager:9093"}},
					}},
				},
				GlobalLabels: map[string]string{"venti": "development"},
			},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := loadAlertingConfigFile(tc.file)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
