package main

import (
	"github.com/kuoss/venti/server"
)

// Version will be overwrited by ldflags
var (
	Version = "unknown"
)

func main() {
	server.Run(Version)
}
