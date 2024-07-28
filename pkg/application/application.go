package application

import (
	"fmt"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/alerter"
	"github.com/kuoss/venti/pkg/config"
	"github.com/kuoss/venti/pkg/handler"
	"github.com/kuoss/venti/pkg/service"
)

type IApp interface {
	Run(version string, addr ...string) error
}

type App struct{}

func (a App) Run(version string, addr ...string) error {
	logger.Infof("Starting Venti 💨 version=%s", version)

	// Load configuration
	cfg, err := new(config.ConfigProvider).New(version)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize services
	services, err := service.NewServices(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Start alerter
	alerter := alerter.New(cfg, services.AlertingService)
	if err = alerter.Start(); err != nil {
		// test unreachable
		return fmt.Errorf("failed to start alerter: %w", err)
	}

	// Start server
	router := handler.NewRouter(services)
	logger.Infof("listen %v", addr)
	if err = router.Run(addr...); err != nil {
		// test unreachable
		return fmt.Errorf("failed to run router: %w", err)
	}
	return nil
}
