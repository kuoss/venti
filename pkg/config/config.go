package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

// Load EtcUser, DatasourceConfig files only.
// TODO: each Config filepath could be parameter.
func Load(version string) (*model.Config, error) {
	log.Println("loading configurations...")

	datasourceConfig, err := loadDatasourceConfigFromFilepath("etc/datasources.yaml")
	if err != nil {
		return nil, fmt.Errorf("error on loadDatasourceConfigFromFilepath: %w", err)
	}

	userConfig, err := loadUserConfigFromFilepath("etc/users.yaml")
	if err != nil {
		return nil, fmt.Errorf("error on loadUserConfigFromFilepath: %w", err)
	}

	return &model.Config{
		Version:          version,
		DatasourceConfig: datasourceConfig,
		UserConfig:       userConfig,
	}, nil
}

func loadDatasourceConfigFromFilepath(filepath string) (*model.DatasourceConfig, error) {
	log.Println("loading datasource config file:", filepath)
	yamlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error on ReadFile: %s", err)
	}
	var datasourceConfig *model.DatasourceConfig
	if err := yaml.UnmarshalStrict(yamlBytes, &datasourceConfig); err != nil {
		return nil, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}
	// default
	if datasourceConfig.QueryTimeout == 0 {
		datasourceConfig.QueryTimeout = 30 * time.Second
	}
	if datasourceConfig.Discovery.AnnotationKey == "" {
		datasourceConfig.Discovery.AnnotationKey = "kuoss.org/datasource"
	}
	return datasourceConfig, nil
}

func loadUserConfigFromFilepath(filepath string) (*model.UserConfig, error) {
	log.Println("loading user config file:", filepath)
	yamlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error on ReadFile: %s", err)
	}
	var userConfig *model.UserConfig
	if err := yaml.UnmarshalStrict(yamlBytes, &userConfig); err != nil {
		return nil, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}
	return userConfig, nil
}
