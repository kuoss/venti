package config

import (
	"fmt"
	"os"
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
	cfg, err := Load("Unknown")
	assert.Nil(t, err)
	assert.Equal(t, cfg.Version, "Unknown")
	assert.Equal(t, []model.Datasource{
		{Type: model.DatasourceTypePrometheus, Name: "prometheus", URL: "http://localhost:9090"},
		{Type: model.DatasourceTypeLethe, Name: "lethe", URL: "http://localhost:6060"},
	}, cfg.DatasourceConfig.Datasources)
	assert.Equal(t, model.UserConfig{EtcUsers: []model.EtcUser{
		{Username: "admin", Hash: "$2a$12$VcCDgh2NDk07JGN0rjGbM.Ad41qVR/YFJcgHp0UGns5JDymv..TOG", IsAdmin: true},
	}}, cfg.UserConfig)
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
			"etc/datasources.yml",
			&model.DatasourceConfig{
				QueryTimeout: 30000000000,
				Datasources: []model.Datasource{
					{Type: "prometheus", Name: "prometheus", URL: "http://localhost:9090"},
					{Type: "lethe", Name: "lethe", URL: "http://localhost:6060"},
				},
				Discovery: model.Discovery{Enabled: false, MainNamespace: "", AnnotationKey: "kuoss.org/datasource-type"},
			},
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
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
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
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
