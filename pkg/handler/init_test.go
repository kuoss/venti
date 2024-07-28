package handler

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/config"
	"github.com/kuoss/venti/pkg/mocker"
	"github.com/kuoss/venti/pkg/mocker/alertmanager"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service"
)

var (
	services         *service.Services
	handlers         *Handlers
	alertmanagerMock *mocker.Server
	cfg              = &config.Config{
		AppInfo:    model.AppInfo{Version: "Unknown"},
		UserConfig: model.UserConfig{},
		DatasourceConfig: model.DatasourceConfig{
			Datasources: []model.Datasource{
				{Type: model.DatasourceTypePrometheus, Name: "prometheus", IsMain: true},
			},
		},
	}
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	var err error
	alertmanagerMock, err = alertmanager.New()
	if err != nil {
		panic(err)
	}
	err = os.Chdir("../..") // project root
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("working directory:", wd)
	cfg.AlertingConfig = model.AlertingConfig{
		AlertmanagerConfigs: model.AlertmanagerConfigs{
			&model.AlertmanagerConfig{
				StaticConfig: []*model.TargetGroup{
					{Targets: []string{alertmanagerMock.URL}},
				},
			},
		},
	}
	services, err = service.NewServices(cfg)
	if err != nil {
		panic(err)
	}
	handlers = loadHandlers(services)
}

func shutdown() {
	alertmanagerMock.Close()
}
