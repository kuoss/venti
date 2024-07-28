package main

import (
	"log"
	"os"

	"github.com/kuoss/venti/pkg/application"
)

var (
	Version = "development" // Version will be overwritten by ldflags
	addr    = ":3030"

	app  application.IApp = application.App{}
	exit                  = os.Exit
)

func main() {
	if err := app.Run(Version, addr); err != nil {
		log.Printf("application error: %v", err)
		exit(1)
	} else {
		exit(0)
	}
}
