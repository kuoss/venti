package main

import (
	"log"

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
		log.Fatalf("config load failed: %s", err.Error())
	}
	stores, err := store.LoadStores(cfg)
	if err != nil {
		log.Fatalf("load store failed: %s", err.Error())
	}

	// starting
	log.Printf("venti starting.... version %s", Version)

	alerter := alerter.NewAlerter(stores)
	alerter.Start()

	router := handler.NewRouter(cfg, stores)
	_ = router.Run() // :8080
}
