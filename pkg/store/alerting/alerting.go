package alerting

import (
	"fmt"
	"os"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

type AlertingStore struct {
	AlertingFile *model.AlertingFile
}

func New(file string) (alertingStore *AlertingStore) {
	logger.Infof("loading alertrules...")
	alertingFile, err := loadAlertingFile(file)
	if err != nil {
		logger.Warnf("error on loadAlertingFile(skipped): %s", err)
		return &AlertingStore{AlertingFile: &model.AlertingFile{Alertings: []model.Alerting{}}}
	}
	return &AlertingStore{AlertingFile: alertingFile}
}

func loadAlertingFile(file string) (*model.AlertingFile, error) {
	logger.Infof("load alertrule file: %s", file)
	if file == "" {
		file = "etc/alerting.yml"
	}
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error on ReadFile: %w", err)
	}
	var alertingFile *model.AlertingFile
	if err := yaml.UnmarshalStrict(yamlBytes, &alertingFile); err != nil {
		return nil, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}
	return alertingFile, nil
}

func (s *AlertingStore) GetAlertmanagerURL() string {
	if len(s.AlertingFile.Alertings) > 0 {
		return s.AlertingFile.Alertings[0].URL
	}
	return ""
}
