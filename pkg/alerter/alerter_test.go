package alerter

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/config"
	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	alertingservice "github.com/kuoss/venti/pkg/service/alerting"
	datasourceservice "github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/discovery"
	remoteservice "github.com/kuoss/venti/pkg/service/remote"
	"github.com/stretchr/testify/require"
)

type mockAlertingService struct{}

func (m *mockAlertingService) DoAlert() error {
	return fmt.Errorf("mock doAlert err")
}

var (
	alerter1         *Alerter
	servers          *ms.Servers
	alertingService1 *alertingservice.AlertingService
	ruleFiles1       []model.RuleFile = []model.RuleFile{
		{
			Kind:               "AlertRuleFile",
			CommonLabels:       map[string]string{"rulefile": "sample-v3", "severity": "silence"},
			DatasourceSelector: model.DatasourceSelector{Type: model.DatasourceTypePrometheus},
			RuleGroups: []model.RuleGroup{{
				Name:     "sample",
				Interval: 0,
				Limit:    0,
				Rules: []model.Rule{
					{Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
					{Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, Annotations: map[string]string{"summary": "Monday"}},
					{Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
				}},
			},
		},
	}
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func shutdown() {
	servers.Close()
}

func setup() {
	err := os.Chdir("../..")
	if err != nil {
		panic(err)
	}
	servers = ms.New(ms.Requirements{
		{Type: ms.TypeAlertmanager, Name: "alertmanager1", IsMain: false},
		{Type: ms.TypeLethe, Name: "lethe1", IsMain: true},
		{Type: ms.TypeLethe, Name: "lethe2", IsMain: false},
		{Type: ms.TypePrometheus, Name: "prometheus1", IsMain: true},
		{Type: ms.TypePrometheus, Name: "prometheus2", IsMain: false},
		{Type: ms.TypePrometheus, Name: "prometheus3", IsMain: false},
	})
	cfg := &config.Config{
		AlertingConfig: model.AlertingConfig{
			AlertmanagerConfigs: model.AlertmanagerConfigs{
				&model.AlertmanagerConfig{
					StaticConfig: []*model.TargetGroup{
						{Targets: []string{servers.GetServersByType(ms.TypeAlertmanager)[0].URL}},
					},
				},
			},
		},
	}
	datasourceConfig := &model.DatasourceConfig{
		Datasources: servers.GetDatasources(),
	}
	datasourceService, err := datasourceservice.New(datasourceConfig, discovery.Discoverer(nil))
	if err != nil {
		panic(err)
	}
	remoteService := remoteservice.New(&http.Client{}, 30*time.Second)
	alertingService1 = alertingservice.New(cfg, ruleFiles1, datasourceService, remoteService)
	alerter1 = New(cfg, alertingService1)
}

func TestNew(t *testing.T) {
	require.Equal(t, 1, len(ruleFiles1))
	require.Equal(t, 1, len(ruleFiles1[0].RuleGroups))
	require.Equal(t, 3, len(ruleFiles1[0].RuleGroups[0].Rules))
	require.Equal(t, model.DatasourceSelector{Type: "prometheus"}, ruleFiles1[0].DatasourceSelector)
	require.Equal(t, false, alerter1.isRunning)
}

func TestStartAndStop(t *testing.T) {
	tempAlerter := New(&config.Config{}, alertingService1)
	tempAlerter.evaluationInterval = 1000 * time.Millisecond

	// start(ok) start(error)
	err := tempAlerter.Start()
	require.NoError(t, err)
	err = tempAlerter.Start()
	require.EqualError(t, err, "already running")

	time.Sleep(time.Second)
	// stop(ok) stop(error)
	err = tempAlerter.Stop()
	require.NoError(t, err)
	err = tempAlerter.Stop()
	require.EqualError(t, err, "already stopped")
}

func TestOnce(t *testing.T) {
	// DoAlert ok
	alerter1.Once()

	// DoAlert error
	temp := alerter1.alertingService
	alerter1.alertingService = &mockAlertingService{}
	alerter1.Once()
	alerter1.alertingService = temp
}
