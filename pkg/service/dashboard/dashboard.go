package dashboard

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

type DashboardService struct {
	dashboards []model.Dashboard
}

func getDashboardFilesFromPath(dirpath string) ([]string, error) {
	if dirpath == "" {
		dirpath = "etc/dashboards"
	}
	files, err := filepath.Glob(dirpath + "/*.y*ml")
	if err != nil {
		return nil, fmt.Errorf("glob err: %w", err)
	}
	files2, err := filepath.Glob(dirpath + "/*/*.y*ml")
	if err != nil {
		return nil, fmt.Errorf("glob err: %w", err)
	}
	files = append(files, files2...)
	if len(files) < 1 {
		return nil, fmt.Errorf("no dashboard file: dirpath: %s", dirpath)
	}
	return files, nil
}

func New(dirpath string) (*DashboardService, error) {
	logger.Debugf("NewDashboardService...")
	files, err := getDashboardFilesFromPath(dirpath)
	if err != nil {
		return nil, fmt.Errorf("getDashboardFilesFromPath err: %w", err)
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
	return &DashboardService{dashboards: dashboards}, nil
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

func (s *DashboardService) Dashboards() []model.Dashboard {
	return s.dashboards
}
