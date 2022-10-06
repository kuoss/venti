package main

import (
	"os"

	"github.com/kuoss/venti/server"
)

var ventiVersion string

func main() {
	if ventiVersion == "" {
		ventiVersion = os.Getenv("VENTI_VERSION")
	}
	server.Run(ventiVersion)
}
