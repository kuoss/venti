export interface RuntimeInfo {
  startTime: string;
  CWD: string;
  reloadConfigSuccess: boolean;
  lastConfigTime: string;
  goroutineCount: number;
  GOMAXPROCS: number;
  GOMEMLIMIT: number;
  GOGC: string;
  GODEBUG: string;
}

export interface BuildInfo {
  version: string;
  revision: string;
  branch: string;
  buildUser: string;
  buildDate: string;
  goVersion: string;
}

interface Alertmanager {
  url: string;
}

export interface Alertmanagers {
  activeAlertmanagers: Alertmanager[];
  droppedAlertmanagers: Alertmanager[];
}