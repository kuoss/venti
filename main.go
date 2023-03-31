package main

import (
	"github.com/kuoss/venti/pkg"
)

// Version will be overwrited by ldflags
var (
	Version = "unknown"
)

func main() {
	pkg.Run(Version)
}
