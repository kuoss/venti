export interface Dashboard {
  title: string;
}

export interface PanelConfig {
  title: string;
  type: string;
// Headers      []string
  targets: Target[];
// ChartOptions *ChartOptions
}

export default interface Target {
expr: string;
  // Legend     string
// Legends    []string
// Unit       string
// Columns    []string
// Headers    []string
// Key        string
// Thresholds []Threshold
}