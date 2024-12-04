package main

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	var originalPort = port
	port = 0
	defer func() {
		port = originalPort
	}()
	go main()
	time.Sleep(time.Duration(1) * time.Second)
}
