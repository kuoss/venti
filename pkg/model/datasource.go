package model

type Datasource struct {
	Type              DatasourceType `json:"type" yaml:"type"`
	Name              string         `json:"name" yaml:"name"`
	URL               string         `json:"url" yaml:"url"`
	BasicAuth         bool           `json:"basicAuth" yaml:"basicAuth"`
	BasicAuthUser     string         `json:"basicAuthUser" yaml:"basicAuthUser"`
	BasicAuthPassword string         `json:"basicAuthPassword" yaml:"basicAuthPassword"`
	IsMain            bool           `json:"isMain,omitempty" yaml:"isMain,omitempty"`
	IsDiscovered      bool           `json:"isDiscovered,omitempty" yaml:"isDiscovered,omitempty"`
}

type DatasourceType string

const (
	DatasourceTypeNone       DatasourceType = ""
	DatasourceTypePrometheus DatasourceType = "prometheus"
	DatasourceTypeLethe      DatasourceType = "lethe"
)
