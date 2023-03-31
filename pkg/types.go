package pkg

type Annotation struct {
	Type DatasourceType `json:"type,omitempty"` // default: "kuoss.org/datasource"
	Port string         `json:"port,omitempty"` // default: "kuoss.org/port"
}
