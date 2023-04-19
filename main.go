package main

import (
	"log"

	"github.com/kuoss/venti/pkg/alerter"
	"github.com/kuoss/venti/pkg/config"
	"github.com/kuoss/venti/pkg/handler"
	"github.com/kuoss/venti/pkg/logger"
	"github.com/kuoss/venti/pkg/store"
)

var (
	Version = "unknown" // Version will be overwrited by ldflags
)

func main() {
	log := logger.GetLogger()
	//load configuration
	cfg, err := configuration.Load(Version)
	if err != nil {
		log.Errorf("config load failed: %s", err.Error())
	}
	stores, err := store.LoadStores(cfg)
	if err != nil {
		log.Errorf("load store failed: %s", err.Error())
	}

	// run
	log.Infof("💨 venti running.... version %s", Version)

  alerter := alerter.NewAlerter(stores)
	alerter.Start()

	router := handler.NewRouter(cfg, stores)
	_ = router.Run() // :8080
}
