export enum DatasourceType {
  Prometheus = 'persist',
  Lethe = 'demand',
}

export interface Datasource {
  name: string
  type: string
  url: string
  isMain: boolean
  isDiscovered: boolean
  health: number
  targets: Target[]
}

export interface Target {
  discoveredLabels: DiscoveredLabels
  globalUrl: string
  health: string
  labels: object
  lastError: string
  lastScrape: string
  lastScrapeDuration: number
  scrapeInterval: string
  scrapePool: string
  scrapeTimeout: string
  scrapeUrl: string
}

export interface DiscoveredLabels {
  job: string
  __address__: string
  __meta_kubernetes_namespace: string
  __meta_kubernetes_node_name: string
  __meta_kubernetes_pod_name: string
  __meta_kubernetes_service_name: string
}