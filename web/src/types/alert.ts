export interface Annotations {
  summary: string
}

export interface Labels {
  severity: string
}

export interface AlertRule {
  name: string;
  expr: string;
  for: string;
  state: string;
  labels: Labels;
  annotations: Annotations;
  alert: string;
}

export interface AlertGroup {
  name: string;
  datasource: string;
  rules: AlertRule[];
}