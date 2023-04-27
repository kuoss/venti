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
		logger.Fatalf("config load failed: %s", err)
	}
	stores, err := store.NewStores(cfg)
	if err != nil {
		logger.Fatalf("load store failed: %s", err)
	}

	// starting
	logger.Infof("💨 venti starting.... version %s", Version)

	alerter := alerter.New(stores)
	err = alerter.Start()
	if err != nil {
		logger.Warnf("error on alerter.Start: %s", err)
	}

	router := handler.NewRouter(cfg, stores)
	logger.Infof("listen :8080")
	_ = router.Run() // :8080
}
