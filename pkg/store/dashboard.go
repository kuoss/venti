package store

// dashboard
type Dashboard struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows"`
}

type Row struct {
	Panels []Panel `json:"panels"`
}

type Panel struct {
	Title        string        `json:"title" yaml:"title"`
	Type         string        `json:"type" yaml:"type"`
	Headers      []string      `json:"headers,omitempty" yaml:"headers,omitempty"`
	Targets      []Target      `json:"targets" yaml:"targets"`
	ChartOptions *ChartOptions `json:"chartOptions,omitempty" yaml:"chartOptions,omitempty"`
}

type ChartOptions struct {
	YMax int `json:"yMax,omitempty" yaml:"yMax,omitempty"`
}

// todo: what Legend, Legends is for?
type Target struct {
	Expr       string      `json:"expr"`
	Legend     string      `json:"legend,omitempty" yaml:"legend,omitempty"`
	Legends    []string    `json:"legends,omitempty" yaml:"legends,omitempty"`
	Unit       string      `json:"unit,omitempty" yaml:"unit,omitempty"`
	Columns    []string    `json:"columns,omitempty" yaml:"columns,omitempty"`
	Headers    []string    `json:"headers,omitempty" yaml:"headers,omitempty"`
	Key        string      `json:"key,omitempty" yaml:"key,omitempty"`
	Thresholds []Threshold `json:"thresholds,omitempty" yaml:"thresholds,omitempty"`
}

type Threshold struct {
	Values []int `yaml:"values,omitempty" json:"values,omitempty"`
	Invert bool  `yaml:"invert,omitempty" json:"invert,omitempty"`
}
