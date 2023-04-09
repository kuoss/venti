package store

import (
	"github.com/kuoss/venti/pkg/model"
	"log"
	"os"
	"path/filepath"
)

type DashboardStore struct {
	dashboards []model.Dashboard
}

// NewDashboardStore pattern parameter is root filepath pattern for dashboard yaml files. ex) etc/dashboards/**/*.yaml
func NewDashboardStore(pattern string) (*DashboardStore, error) {
	log.Println("Loading dashboards...")
	// todo : why using custom glob function instead using filepath.Glob
	files, err := filepath.Glob("etc/dashboards/**/*.yaml")
	if err != nil {
		return nil, err
	}
	dashboards := make([]model.Dashboard, 0)
	for _, filename := range files {
		log.Printf("dashboard file: %s\n", filename)
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("error open dashboard file: %s\n", err.Error())
		}
		var d *model.Dashboard
		err = loadYaml(f, d)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, *d)
	}
	return &DashboardStore{dashboards: dashboards}, nil
}

func (dbs *DashboardStore) Dashboards() []model.Dashboard {
	return dbs.dashboards
}
