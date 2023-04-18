package main

import (
	"github.com/kuoss/venti/pkg/configuration"
	"github.com/kuoss/venti/pkg/handler"
	"github.com/kuoss/venti/pkg/logger"
	"github.com/kuoss/venti/pkg/store"
)

// Version will be overwrited by ldflags
var (
	Version = "unknown"
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
	router := handler.SetupRouter(cfg, stores)

	// run
	log.Infof("💨 venti running.... version %s", Version)
	_ = router.Run() // :8080
}
