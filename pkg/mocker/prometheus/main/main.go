package main

import (
	"log"
	"time"

	"github.com/kuoss/venti/pkg/mocker/prometheus"
)

var port = 9091

func main() {
	log.Printf("Staring prometheus...")
	s := prometheus.NewWithPort(port)
	if err := s.Start(); err != nil {
		log.Fatalf("Start err: %s", err.Error())
	}
	log.Println("Started...")
	for {
		time.Sleep(time.Second)
	}
}
