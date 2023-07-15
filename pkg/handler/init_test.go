package handler

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/mocker"
	"github.com/kuoss/venti/pkg/mocker/alertmanager"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service"
)

var (
	services         *service.Services
	handlers         *Handlers
	alertmanagerMock *mocker.Server
	cfg              = &model.Config{
		Version:    "Unknown",
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
	alertmanagerMock, err = alertmanager.New(0)
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
	services, err = service.NewServices(cfg)
	if err != nil {
		panic(err)
	}
	services.AlertingService.AlertingFile.Alertings[0].URL = alertmanagerMock.URL
	handlers = loadHandlers(cfg, services)
}

func shutdown() {
	alertmanagerMock.Close()
}
