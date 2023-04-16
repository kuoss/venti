package alerter

// type ValueType string

// const (
// 	ValueTypeNone   ValueType = "none"
// 	ValueTypeVector ValueType = "vector"
// 	ValueTypeScalar ValueType = "scalar"
// 	ValueTypeMatrix ValueType = "matrix"
// 	ValueTypeString ValueType = "string"
// )

// type AlertRuleGroupList struct {
// 	Groups []AlertRuleGroup `json:"groups"`
// }

// type AlertRuleGroup struct {
// 	Name           string             `json:"name"`
// 	Rules          []AlertRule        `json:"rules"`
// 	DatasourceType pkg.DatasourceType `json:"datasource" yaml:"datasource"`
// 	CommonLabels   map[string]string  `json:"commonLabels,omitempty" yaml:"commonLabels,omitempty"`
// }

// type AlertRule struct {
// 	Alert       string            `json:"alert,omitempty"`
// 	Expr        string            `json:"expr"`
// 	For         time.Duration     `json:"for,omitempty"`
// 	Labels      map[string]string `json:"labels,omitempty"`
// 	Annotations map[string]string `json:"annotations,omitempty"`
// 	State       AlertState        `json:"state,omitempty"`
// 	ActiveAt    time.Time         `json:"activeStartTime,omitempty"`
// }

// type AlertState string

// const (
// 	AlertStateInactive AlertState = "inactive"
// 	AlertStatePending  AlertState = "pending"
// 	AlertStateFiring   AlertState = "firing"
// )

// type Sample struct {
// 	Vaue   []interface{}     `json:"value"`
// 	Metric map[string]string `json:"metric"`
// }

// type Vector []Sample

// // {"data":{"result":[],"resultType":"logs"},"status":"success"}
// // {"status":"success","data":{"resultType":"vector","result":[]}}
// type QueryResult struct {
// 	Data   QueryData `json:"data"`
// 	Status string    `json:"status"`
// }
// type QueryData struct {
// 	ResultType ValueType `json:"resultType"`
// 	Result     Vector    `json:"result"`
// }

// type Alert struct {
// 	Status       string            `json:"status"`
// 	Labels       map[string]string `json:"labels,omitempty"`
// 	Annotations  map[string]string `json:"annotations,omitempty"`
// 	GeneratorURL string            `json:"generatorURL,omitempty"`
// }
