package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/util"
)

type Config struct {
	AppInfo          model.AppInfo
	GlobalConfig     model.GlobalConfig
	DatasourceConfig model.DatasourceConfig
	UserConfig       model.UserConfig
	AlertingConfig   model.AlertingConfig
}

type IConfigProvider interface {
	New(version string) (*Config, error)
}

type ConfigProvider struct{}

// TODO: each Config filepath could be parameter.
func (p *ConfigProvider) New(version string) (*Config, error) {
	logger.Infof("loading configurations...")

	cfg := &Config{
		AppInfo: model.AppInfo{
			Version: version,
		},
	}
	if err := cfg.loadGlobalConfigFile("etc/venti.yml"); err != nil {
		return nil, fmt.Errorf("loadGlobalConfigFile err: %w", err)
	}
	if err := cfg.loadDatasourceConfigFile("etc/datasources.yml"); err != nil {
		return nil, fmt.Errorf("loadDatasourceConfigFile err: %w", err)
	}
	if err := cfg.loadUserConfigFile("etc/users.yml"); err != nil {
		return nil, fmt.Errorf("loadUserConfigFile err: %w", err)
	}
	if err := cfg.loadAlertingConfigFile("etc/alerting.yml"); err != nil {
		return nil, fmt.Errorf("loadAlertingConfigFile err: %w", err)
	}
	return cfg, nil
}

func (c *Config) loadGlobalConfigFile(file string) error {
	var cfg model.GlobalConfig
	logger.Infof("loading global config file: %s", file)
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error on ReadFile: %w", err)
	}
	if err := util.UnmarshalStrict(yamlBytes, &cfg); err != nil {
		return fmt.Errorf("unmarshalStrict err: %w", err)
	}

	// set gin mode
	gin.SetMode(cfg.GinMode)
	// set log level
	logLevel, err := logger.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("ParseLevel err: %s", err)
	}
	logger.SetLevel(logLevel)

	c.GlobalConfig = cfg
	return nil
}

func (c *Config) loadDatasourceConfigFile(file string) error {
	logger.Infof("loading datasource config file: %s", file)
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error on ReadFile: %w", err)
	}
	var cfg model.DatasourceConfig
	if err := util.UnmarshalStrict(yamlBytes, &cfg); err != nil {
		return fmt.Errorf("unmarshalStrict err: %w", err)
	}

	// default
	if cfg.QueryTimeout == 0 {
		cfg.QueryTimeout = 30 * time.Second
	}
	if cfg.Discovery.AnnotationKey == "" {
		cfg.Discovery.AnnotationKey = "kuoss.org/datasource-type"
	}

	c.DatasourceConfig = cfg
	return nil
}

func (c *Config) loadUserConfigFile(file string) error {
	logger.Infof("loading user config file: %s", file)
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error on ReadFile: %w", err)
	}
	var cfg model.UserConfig
	if err := util.UnmarshalStrict(yamlBytes, &cfg); err != nil {
		return fmt.Errorf("unmarshalStrict err: %w", err)
	}
	c.UserConfig = cfg
	return nil
}

func (c *Config) loadAlertingConfigFile(file string) error {
	logger.Infof("loading alerting config file: %s", file)
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error on ReadFile: %w", err)
	}
	var cfg *model.AlertingConfigFile
	if err := util.UnmarshalStrict(yamlBytes, &cfg); err != nil {
		return fmt.Errorf("unmarshalStrict err: %w", err)
	}

	// default
	alertingConfig := cfg.AlertingConfig
	if alertingConfig.EvaluationInterval == 0 {
		alertingConfig.EvaluationInterval = 5 * time.Second
	}
	c.AlertingConfig = alertingConfig
	return nil
}
