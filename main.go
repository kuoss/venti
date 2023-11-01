package main

import (
	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/alerter"
	"github.com/kuoss/venti/pkg/config"
	"github.com/kuoss/venti/pkg/handler"
	"github.com/kuoss/venti/pkg/service"
)

var (
	Version = "development" // Version will be overwrited by ldflags
)

func main() {
	logger.Infof("Starting Venti 💨 version=%s", Version)

	// load configuration
	cfg, err := config.Load(Version)
	if err != nil {
		logger.Fatalf("config.Load err: %s", err)
	}

	// init stores
	services, err := service.NewServices(cfg)
	if err != nil {
		logger.Fatalf("NewServices err: %s", err)
	}

	// alerter start
	alerter := alerter.New(cfg, services.AlertingService)
	err = alerter.Start()
	if err != nil {
		logger.Fatalf("alerter start err: %s", err)
	}

	// router run
	router := handler.NewRouter(cfg, services)
	logger.Infof("listen :3030")
	err = router.Run(":3030")
	if err != nil {
		logger.Fatalf("router run err: %s", err)
	}
}
