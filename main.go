package main

import (
	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/alerter"
	"github.com/kuoss/venti/pkg/config"
	"github.com/kuoss/venti/pkg/handler"
	"github.com/kuoss/venti/pkg/store"
)

var (
	Version = "unknown" // Version will be overwrited by ldflags
)

func main() {

	// load configuration
	cfg, err := config.Load(Version)
	if err != nil {
		logger.Fatalf("config.Load err: %s", err)
	}

	// init stores
	stores, err := store.NewStores(cfg)
	if err != nil {
		logger.Fatalf("NewStores err: %s", err)
	}

	// show starting & version
	logger.Infof("💨 venti starting.... version %s", Version)

	// alerter start
	alerter := alerter.New(stores)
	err = alerter.Start()
	if err != nil {
		logger.Warnf("alerter.Start err: %s", err)
	}

	// router run
	router := handler.NewRouter(cfg, stores)
	logger.Infof("listen :3030")
	err = router.Run(":3030")
	if err != nil {
		logger.Fatalf("router.Run err: %s", err)
	}
}
