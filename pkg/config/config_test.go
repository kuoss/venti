package config

import (
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Chdir("../..")
}

func TestLoad(t *testing.T) {
	cfg, err := Load("Unknown")
	assert.Nil(t, err)
	assert.Equal(t, cfg.Version, "Unknown")
	assert.ElementsMatch(t, []model.Datasource{
		{Type: model.DatasourceTypePrometheus, Name: "Prometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
		{Type: model.DatasourceTypeLethe, Name: "Lethe", URL: "http://lethe:3100", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
	}, cfg.DatasourceConfig.Datasources)
	assert.Equal(t, cfg.UserConfig, model.UserConfig{EtcUsers: []model.EtcUser{
		{Username: "admin", Hash: "$2a$12$VcCDgh2NDk07JGN0rjGbM.Ad41qVR/YFJcgHp0UGns5JDymv..TOG", IsAdmin: true},
	}})
}

func TestLoadDatasourceConfigFromFilepath(t *testing.T) {

	want := &model.DatasourceConfig{
		QueryTimeout: 0,
		Datasources: []*model.Datasource{
			{Type: "prometheus",
				Name:              "Prometheus",
				URL:               "http://prometheus:9090",
				BasicAuth:         false,
				BasicAuthUser:     "",
				BasicAuthPassword: "",
				IsMain:            false,
				IsDiscovered:      false,
			},
			{Type: "lethe",
				Name:              "Lethe",
				URL:               "http://lethe:3100",
				BasicAuth:         false,
				BasicAuthUser:     "",
				BasicAuthPassword: "",
				IsMain:            false,
				IsDiscovered:      false,
			},
		},
		Discovery: model.Discovery{
			Enabled:          false,
			MainNamespace:    "",
			AnnotationKey:    "",
			ByNamePrometheus: false,
			ByNameLethe:      false,
		},
	}
	datasourceConfig, err := loadDatasourceConfigFromFilepath("etc/datasource.checks.yaml")
	assert.Nil(t, err)
	assert.Equal(t, want, datasourceConfig)

}
