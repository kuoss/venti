package main

import (
	"os"

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
		logger.Errorf("config load failed: %s", err)
		os.Exit(1)
	}
	stores, err := store.LoadStores(cfg)
	if err != nil {
		logger.Errorf("load store failed: %s", err)
		os.Exit(2)
	}

	// starting
	logger.Infof("💨 venti starting.... version %s", Version)

	alerter := alerter.NewAlerter(stores)
	alerter.Start()

	router := handler.NewRouter(cfg, stores)
	_ = router.Run() // :8080
}
