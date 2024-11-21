export interface Dashboard {
  title: string;
}

export interface ChartOptions {
  yMax: Number;
}

export interface PanelConfig {
  title: string;
  type: string;
  targets: Target[];
  chartOptions: ChartOptions;
  // Headers      []string
  // ChartOptions *ChartOptions
}

export default interface Target {
  expr: string;
  legend: string;
  // Legends    []string
  // Unit       string
  // Columns    []string
  // Headers    []string
  // Key        string
  // Thresholds []Threshold
}