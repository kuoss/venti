package alerter

import (
	"fmt"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/config"
	"github.com/kuoss/venti/pkg/service/alerting"
)

type Alerter struct {
	alertingService    alerting.IAlertingService
	evaluationInterval time.Duration
	isRunning          bool
	quitCh             chan bool
}

func New(cfg *config.Config, alertingService alerting.IAlertingService) *Alerter {
	return &Alerter{
		alertingService:    alertingService,
		evaluationInterval: cfg.AlertingConfig.EvaluationInterval,
	}
}

func (a *Alerter) Start() error {
	if a.isRunning {
		return fmt.Errorf("already running")
	}
	a.isRunning = true
	logger.Infof("starting alerter...")
	a.quitCh = make(chan bool)
	go a.loop(a.quitCh)
	return nil
}

func (a *Alerter) Stop() error {
	if !a.isRunning {
		return fmt.Errorf("already stopped")
	}
	logger.Infof("stopping alerter...")
	a.quitCh <- true
	a.isRunning = false
	return nil
}

func (a *Alerter) loop(quitCh chan bool) {
	for {
		select {
		case <-quitCh:
			logger.Infof("alerter stopped")
			close(quitCh)
			return
		default:
			a.Once()
			logger.Infof("sleep: %s", a.evaluationInterval)
			time.Sleep(a.evaluationInterval)
		}
	}
}

func (a *Alerter) Once() {
	err := a.alertingService.DoAlert()
	if err != nil {
		logger.Errorf("CatchAlertingRule err: %s", err)
	}
}
