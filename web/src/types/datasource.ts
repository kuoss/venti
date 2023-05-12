export interface Datasource {
  name: string;
  type: string;
  url: string;
  isDefault: boolean;
  isDiscovered: boolean;
  health: number;
  targets: Target[];
}

export interface Target {
  age: string;
  discoveredLabels: DiscoveredLabels,
  health: string;
  icon: string;
  job: string;
  name: string;
  // Legend     string      `json:"legend,omitempty" yaml:"legend,omitempty"`
// Legends    []string    `json:"legends,omitempty" yaml:"legends,omitempty"`
// Unit       string      `json:"unit,omitempty" yaml:"unit,omitempty"`
// Columns    []string    `json:"columns,omitempty" yaml:"columns,omitempty"`
// Headers    []string    `json:"headers,omitempty" yaml:"headers,omitempty"`
// Key        string      `json:"key,omitempty" yaml:"key,omitempty"`
// Thresholds []Threshold `json:"thresholds,omitempty" yaml:"thresholds,omitempty"`
}

export interface DiscoveredLabels {
  job: string;
  __address__: string;
}