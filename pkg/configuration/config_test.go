package configuration

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Chdir("../..")
}

func TestLoad(t *testing.T) {
	cfg, err := Load("Unknown")
	assert.Nil(t, err)
	assert.Equal(t, cfg.Version, "Unknown")
	assert.Equal(t, cfg.UserConfig, UsersConfig{EtcUsers: []EtcUser{
		{Username: "admin", Hash: "$2a$12$VcCDgh2NDk07JGN0rjGbM.Ad41qVR/YFJcgHp0UGns5JDymv..TOG", IsAdmin: true},
	}})
	assert.ElementsMatch(t, cfg.DatasourcesConfig.Datasources, []*Datasource{
		{Type: DatasourceTypePrometheus, Name: "Prometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: false},
		{Type: DatasourceTypeLethe, Name: "Lethe", URL: "http://lethe:3100", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: false},
	})
}

func TestLoadDatasourcesConfig(t *testing.T) {

	tests := map[string]struct {
		input io.Reader
		want  *DatasourcesConfig
	}{
		"default case": {
			input: bytes.NewReader([]byte(`
datasources:
- name: Prometheus
  type: prometheus
  url: http://prometheus:9090
- name: Lethe
  type: lethe
  url: http://lethe:3100
`)),
			want: &DatasourcesConfig{
				QueryTimeout: 0,
				Datasources: []*Datasource{
					{Type: "prometheus",
						Name:              "Prometheus",
						URL:               "http://prometheus:9090",
						BasicAuth:         false,
						BasicAuthUser:     "",
						BasicAuthPassword: "",
						IsDefault:         false,
						IsDiscovered:      false,
					},
					{Type: "lethe",
						Name:              "Lethe",
						URL:               "http://lethe:3100",
						BasicAuth:         false,
						BasicAuthUser:     "",
						BasicAuthPassword: "",
						IsDefault:         false,
						IsDiscovered:      false,
					},
				},
				Discovery: Discovery{
					Enabled:          false,
					DefaultNamespace: "",
					AnnotationKey:    "",
					ByNamePrometheus: false,
					ByNameLethe:      false,
				},
			},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(subt *testing.T) {
			var dataSourceConfig *DatasourcesConfig
			err := loadConfig(testcase.input, &dataSourceConfig)
			if err != nil {
				subt.Fatalf("error on loading Datasources Config: %s", err.Error())
			}
			assert.Equal(t, dataSourceConfig, testcase.want)

		})
	}
}
