package store

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

type DashboardStore struct {
	dashboards []model.Dashboard
}

// NewDashboardStore pattern parameter is root filepath pattern for dashboard yaml files. ex) etc/dashboards/**/*.yaml
func NewDashboardStore(pattern string) (*DashboardStore, error) {
	log.Println("Loading dashboards...")
	files, err := filepath.Glob("etc/dashboards/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("error on glob: %w", err)
	}
	var dashboards []model.Dashboard
	for _, filename := range files {
		dashboard, err := loadDashboardFromFile(filename)
		if err != nil {
			log.Printf("Warning: error on loadDashboardFromFile(skipped): %s", err)
			continue
		}
		dashboards = append(dashboards, *dashboard)
	}
	return &DashboardStore{dashboards: dashboards}, nil
}

func loadDashboardFromFile(filename string) (*model.Dashboard, error) {
	log.Printf("load dashboard file: %s\n", filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error on Open: %w", err)
	}
	var dashboard *model.Dashboard
	yamlBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error on ReadAll: %w", err)
	}
	if err := yaml.UnmarshalStrict(yamlBytes, &dashboard); err != nil {
		return nil, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}
	return dashboard, nil
}

func (s *DashboardStore) Dashboards() []model.Dashboard {
	return s.dashboards
}
