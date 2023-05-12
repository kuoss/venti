package main

import (
	"fmt"
	"time"

	ms "github.com/kuoss/venti/pkg/mock/servers"
)

func main() {
	_ = ms.New(ms.Requirements{
		{Type: ms.TypeAlertmanager, Port: 9093, Name: "alertmanager1"},
		{Type: ms.TypeLethe, Port: 6060, Name: "lethe1"},
		{Type: ms.TypePrometheus, Port: 9090, Name: "prometheus1"},
	})
	fmt.Println("starting servers...")
	for {
		time.Sleep(10 * 365 * 24 * time.Hour)
	}
}
