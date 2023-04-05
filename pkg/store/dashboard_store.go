package store

import (
	"log"
	"os"
	"path/filepath"
)

type DashboardStore struct {
	dashboards []Dashboard
}

// NewDashboardStore pattern parameter is root filepath pattern for dashboard yaml files. ex) etc/dashboards/**/*.yaml
func NewDashboardStore(pattern string) (*DashboardStore, error) {
	log.Println("Loading dashboards...")
	// todo : why using custom glob function instead using filepath.Glob
	files, err := filepath.Glob("etc/dashboards/**/*.yaml")
	if err != nil {
		return nil, err
	}
	dashboards := make([]Dashboard, 0)
	for _, filename := range files {
		log.Printf("dashboard file: %s\n", filename)
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("error open dashboard file: %s\n", err.Error())
		}
		var d *Dashboard
		err = loadYaml(f, d)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, *d)
	}
	return &DashboardStore{dashboards: dashboards}, nil
}

func (dbs *DashboardStore) Dashboards() []Dashboard {
	return dbs.dashboards
}
