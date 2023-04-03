package main

import (
	"github.com/kuoss/venti/pkg"
	"github.com/kuoss/venti/pkg/configuration"
	"log"
)

// Version will be overwrited by ldflags
var (
	Version = "unknown"
)

func main() {
	//load configuration
	config, err := configuration.Load(Version)
	if err != nil {
		log.Fatalf("configuration load failed. %s", err.Error())
	}
	s, err := pkg.LoadStores(config)
	if err != nil {
		log.Fatalf("load store failed. %s", err.Error())
	}

	r := pkg.LoadRouter(s, config)

	// run
	log.Printf("venti running.... version %s", Version)
	_ = r.Run() // :8080
}
