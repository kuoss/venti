package dashboard

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

type DashboardStore struct {
	dashboards []model.Dashboard
}

func New(pattern string) (*DashboardStore, error) {
	logger.Debugf("NewDashboardStore...")
	if pattern == "" {
		pattern = "etc/dashboards/*.yml"
	}
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob err: %w", err)
	}
	if len(files) < 1 {
		return nil, fmt.Errorf("no dashboard file: pattern: %s", pattern)
	}
	var dashboards []model.Dashboard
	for _, filename := range files {
		dashboard, err := loadDashboardFromFile(filename)
		if err != nil {
			logger.Warnf("Warning: error on loadDashboardFromFile(skipped): %s", err)
			continue
		}
		dashboards = append(dashboards, *dashboard)
	}
	return &DashboardStore{dashboards: dashboards}, nil
}

func loadDashboardFromFile(filename string) (*model.Dashboard, error) {
	logger.Infof("load dashboard file: %s", filename)
	yamlBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error on ReadFile: %w", err)
	}
	var dashboard *model.Dashboard
	if err := yaml.UnmarshalStrict(yamlBytes, &dashboard); err != nil {
		return nil, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}
	return dashboard, nil
}

func (s *DashboardStore) Dashboards() []model.Dashboard {
	return s.dashboards
}
