package alerter

import (
	"fmt"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/alerting"
)

type alerter struct {
	alertingService    *alerting.AlertingService
	evaluationInterval time.Duration
	isRunning          bool
	quitCh             chan bool
}

func New(cfg *model.Config, alertingService *alerting.AlertingService) *alerter {
	return &alerter{
		alertingService:    alertingService,
		evaluationInterval: cfg.GlobalConfig.EvaluationInterval,
	}
}

func (a *alerter) Start() error {
	if a.isRunning {
		return fmt.Errorf("already running")
	}
	a.isRunning = true
	logger.Infof("starting alerter...")
	a.quitCh = make(chan bool)
	go a.loop(a.quitCh)
	return nil
}

func (a *alerter) Stop() error {
	if !a.isRunning {
		return fmt.Errorf("already stopped")
	}
	logger.Infof("stopping alerter...")
	a.quitCh <- true
	a.isRunning = false
	return nil
}

func (a *alerter) loop(quitCh chan bool) {
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

func (a *alerter) Once() {
	err := a.alertingService.DoAlert()
	if err != nil {
		logger.Errorf("CatchAlertingRule err: %s", err)
	}
}
