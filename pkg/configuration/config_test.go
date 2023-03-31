package configuration

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

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
