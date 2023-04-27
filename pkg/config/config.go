package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

// Load EtcUser, DatasourceConfig files only.
// TODO: each Config filepath could be parameter.
func Load(version string) (*model.Config, error) {
	logger.Infof("loading configurations...")

	datasourceConfig, err := loadDatasourceConfigFile("etc/datasources.yml")
	if err != nil {
		return nil, fmt.Errorf("error on loadDatasourceConfigFile: %w", err)
	}

	userConfig, err := loadUserConfigFile("etc/users.yml")
	if err != nil {
		return nil, fmt.Errorf("error on loadUserConfigFile: %w", err)
	}

	return &model.Config{
		Version:          version,
		DatasourceConfig: datasourceConfig,
		UserConfig:       userConfig,
	}, nil
}

func loadDatasourceConfigFile(file string) (*model.DatasourceConfig, error) {
	logger.Infof("loading datasource config file: %s", file)
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error on ReadFile: %w", err)
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
		datasourceConfig.Discovery.AnnotationKey = "kuoss.org/datasource-type"
	}
	return datasourceConfig, nil
}

func loadUserConfigFile(file string) (*model.UserConfig, error) {
	logger.Infof("loading user config file: %s", file)
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error on ReadFile: %w", err)
	}
	var userConfig *model.UserConfig
	if err := yaml.UnmarshalStrict(yamlBytes, &userConfig); err != nil {
		return nil, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}
	return userConfig, nil
}
