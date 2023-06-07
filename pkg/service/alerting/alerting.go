package alerting

import (
	"fmt"
	"os"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/datasource"
	"gopkg.in/yaml.v2"
)

type AlertingService struct {
	AlertingFile model.AlertingFile
	AlertFiles   []model.AlertFile
}

func New(file string, ruleFiles []model.RuleFile, datasourceService *datasource.DatasourceService) (alertingService *AlertingService) {
	logger.Infof("initializing alerting service...")
	alertingFile, err := loadAlertingFile(file)
	if err != nil {
		logger.Warnf("loadAlertingFile err: %s", err.Error())
	}

	var alertFiles []model.AlertFile
	for _, ruleFile := range ruleFiles {
		var alertGroups []model.AlertGroup
		datasources := datasourceService.GetDatasourcesWithSelector(ruleFile.DatasourceSelector)
		for _, ruleGroup := range ruleFile.RuleGroups {
			var ruleAlerts []model.RuleAlert
			for _, rule := range ruleGroup.Rules {
				var alerts []model.Alert
				for i := range datasources {
					alerts = append(alerts, model.Alert{
						Datasource: &datasources[i],
					})
				}
				ruleAlerts = append(ruleAlerts, model.RuleAlert{
					Rule:   rule,
					Alerts: alerts,
				})
			}
			alertGroups = append(alertGroups, model.AlertGroup{
				Name:       ruleGroup.Name,
				Interval:   ruleGroup.Interval,
				RuleAlerts: ruleAlerts,
			})
		}
		alertFiles = append(alertFiles, model.AlertFile{
			CommonLabels:       ruleFile.CommonLabels,
			DatasourceSelector: ruleFile.DatasourceSelector,
			AlertGroups:        alertGroups,
		})
	}
	return &AlertingService{
		AlertingFile: *alertingFile,
		AlertFiles:   alertFiles,
	}
}

func loadAlertingFile(file string) (*model.AlertingFile, error) {
	logger.Infof("load alerting file: %s", file)
	if file == "" {
		file = "etc/alerting.yml"
	}
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return new(model.AlertingFile), fmt.Errorf("readFile err: %w", err)
	}
	var alertingFile *model.AlertingFile
	if err := yaml.UnmarshalStrict(yamlBytes, &alertingFile); err != nil {
		return new(model.AlertingFile), fmt.Errorf("unmarshalStrict err: %w", err)
	}
	return alertingFile, nil
}

func (s *AlertingService) GetAlertmanagerURL() string {
	if len(s.AlertingFile.Alertings) > 0 {
		return s.AlertingFile.Alertings[0].URL
	}
	return ""
}
