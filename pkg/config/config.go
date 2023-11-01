package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

// Load EtcUser, DatasourceConfig files only.
// TODO: each Config filepath could be parameter.
func Load(version string) (*model.Config, error) {
	logger.Infof("loading configurations...")

	appInfo := model.AppInfo{
		Version: version,
	}

	globalConfig, err := loadGlobalConfigFile("etc/venti.yml")
	if err != nil {
		return nil, fmt.Errorf("loadGlobalConfigFile err: %w", err)
	}

	datasourceConfig, err := loadDatasourceConfigFile("etc/datasources.yml")
	if err != nil {
		return nil, fmt.Errorf("loadDatasourceConfigFile err: %w", err)
	}

	userConfig, err := loadUserConfigFile("etc/users.yml")
	if err != nil {
		return nil, fmt.Errorf("loadUserConfigFile err: %w", err)
	}

	alertingConfig, err := loadAlertingConfigFile("etc/alerting.yml")
	if err != nil {
		return nil, fmt.Errorf("loadAlertingConfigFile err: %w", err)
	}

	return &model.Config{
		AppInfo:          appInfo,
		GlobalConfig:     globalConfig,
		DatasourceConfig: *datasourceConfig,
		UserConfig:       *userConfig,
		AlertingConfig:   alertingConfig,
	}, nil
}

func loadGlobalConfigFile(file string) (model.GlobalConfig, error) {
	var globalConfig model.GlobalConfig
	logger.Infof("loading global config file: %s", file)
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return globalConfig, fmt.Errorf("error on ReadFile: %w", err)
	}
	if err := yaml.UnmarshalStrict(yamlBytes, &globalConfig); err != nil {
		return globalConfig, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}

	// ginMode
	ginMode := globalConfig.GinMode
	switch ginMode {
	case gin.DebugMode:
	case gin.ReleaseMode:
	case gin.TestMode:
	default:
		logger.Warnf("gin mode unknown: %s (available mode: debug release test)", globalConfig.GinMode)
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	// logLevel
	logLevel, err := logger.ParseLevel(globalConfig.LogLevel)
	if err != nil {
		logger.Warnf("ParseLevel err: %s", err)
		logLevel = logger.InfoLevel
	}
	logger.SetLevel(logLevel)
	return globalConfig, nil
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

func loadAlertingConfigFile(file string) (model.AlertingConfig, error) {
	logger.Infof("loading alerting config file: %s", file)
	yamlBytes, err := os.ReadFile(file)
	if err != nil {
		return model.AlertingConfig{}, fmt.Errorf("error on ReadFile: %w", err)
	}
	var alertingConfigFile *model.AlertingConfigFile
	if err := yaml.UnmarshalStrict(yamlBytes, &alertingConfigFile); err != nil {
		return model.AlertingConfig{}, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}

	// default
	alertingConfig := alertingConfigFile.AlertingConfig
	if alertingConfig.EvaluationInterval == 0 {
		alertingConfig.EvaluationInterval = 5 * time.Second
	}
	return alertingConfig, nil
}
