package alerting

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	datasourceservice "github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/discovery"
	remoteservice "github.com/kuoss/venti/pkg/service/remote"
	commonModel "github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
)

type mockDatasourceService struct{}

func (m *mockDatasourceService) Reload() error {
	return fmt.Errorf("mock reload err")
}

func (m *mockDatasourceService) GetDatasourcesWithSelector(selector model.DatasourceSelector) []model.Datasource {
	return []model.Datasource{}
}

var (
	servers          *ms.Servers
	alertingService1 *AlertingService
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
					{Alert: "Up", Expr: "up", For: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "Up value={{ $value }}"}},
					{Alert: "AlwaysOn", Expr: "vector(1234)", For: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
					{Alert: "Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "Monday"}},
					{Alert: "NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
				}},
			},
		},
		{
			Kind:               "AlertRuleFile",
			CommonLabels:       map[string]string{"rulefile": "sample-v3", "severity": "silence"},
			DatasourceSelector: model.DatasourceSelector{Type: model.DatasourceTypeLethe},
			RuleGroups: []model.RuleGroup{{
				Name:     "sample2",
				Interval: 0,
				Limit:    0,
				Rules: []model.Rule{
					{Alert: "Pod", Expr: `pod{namespace="namespace01"}`, For: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "PodLogs value={{ $value }}"}},
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
	cfg := &model.Config{
		AlertingConfig: model.AlertingConfig{
			GlobalLabels: map[string]string{
				"global1": "label1",
			},
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
	alertingService1 = New(cfg, ruleFiles1, datasourceService, remoteService)
}

func TestNew(t *testing.T) {
	require.NotZero(t, alertingService1)
}

func TestDoAlert(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := alertingService1.DoAlert()
		require.NoError(t, err)
	})

	t.Run("no_alertingRules", func(t *testing.T) {
		temp := alertingService1.alertingRuleGroups
		alertingService1.alertingRuleGroups = []AlertingRuleGroup{}
		err := alertingService1.DoAlert()
		require.NoError(t, err)
		alertingService1.alertingRuleGroups = temp
	})

	t.Run("reload err", func(t *testing.T) {
		temp := alertingService1.datasourceService
		alertingService1.datasourceService = &mockDatasourceService{}
		alertingService1.datasourceReload = true
		err := alertingService1.DoAlert()
		require.EqualError(t, err, "reload err: mock reload err")
		alertingService1.datasourceService = temp
	})

	t.Run("sendFire err", func(t *testing.T) {
		temp := alertingService1.alertmanagerURL
		alertingService1.alertmanagerURL = ""
		err := alertingService1.DoAlert()
		require.EqualError(t, err, `sendFires err: error on Post: Post "/api/v1/alerts": unsupported protocol scheme ""`)
		alertingService1.alertmanagerURL = temp
	})
}

func TestEvalAlertingRuleGroups(t *testing.T) {
	fires := []Fire{}
	alertingService1.evalAlertingRuleGroups(&fires)
	want := []Fire{
		{Labels: map[string]string{"__name__": "up", "alertname": "Up", "datasource": "prometheus1", "global1": "label1", "hello": "world", "instance": "localhost:9090", "job": "prometheus", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Up value=1"}},
		{Labels: map[string]string{"__name__": "up", "alertname": "Up", "datasource": "prometheus1", "global1": "label1", "hello": "world", "instance2": "localhost:9092", "job": "prometheus2", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Up value=1"}},
		{Labels: map[string]string{"__name__": "up", "alertname": "Up", "datasource": "prometheus2", "global1": "label1", "hello": "world", "instance": "localhost:9090", "job": "prometheus", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Up value=1"}},
		{Labels: map[string]string{"__name__": "up", "alertname": "Up", "datasource": "prometheus2", "global1": "label1", "hello": "world", "instance2": "localhost:9092", "job": "prometheus2", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Up value=1"}},
		{Labels: map[string]string{"__name__": "up", "alertname": "Up", "datasource": "prometheus3", "global1": "label1", "hello": "world", "instance": "localhost:9090", "job": "prometheus", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Up value=1"}},
		{Labels: map[string]string{"__name__": "up", "alertname": "Up", "datasource": "prometheus3", "global1": "label1", "hello": "world", "instance2": "localhost:9092", "job": "prometheus2", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Up value=1"}},
		{Labels: map[string]string{"alertname": "Pod", "datasource": "lethe1", "global1": "label1", "hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "PodLogs value=2"}},
		{Labels: map[string]string{"alertname": "Pod", "datasource": "lethe2", "global1": "label1", "hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "PodLogs value=2"}}}
	require.ElementsMatch(t, fires, want)
}

func TestEvalAlertingRuleGroup(t *testing.T) {
	group := AlertingRuleGroup{}
	evalTime := time.Now()
	fires := []Fire{}
	alertingService1.evalAlertingRuleGroup(&group, evalTime, &fires)
	want := []Fire{}
	require.Equal(t, want, fires)
}

func TestEvalAlertingRule(t *testing.T) {
	evalTime := time.Now()
	testCases := []struct {
		active       map[uint64]*Alert
		commonLabels map[string]string
		want         []Fire
	}{
		{
			map[uint64]*Alert{},
			map[string]string{},
			[]Fire{},
		},
		{
			map[uint64]*Alert{
				1: {},
				2: {},
			},
			map[string]string{},
			[]Fire{},
		},
		{
			map[uint64]*Alert{
				1: {UpdatedAt: evalTime, State: StateFiring},
				2: {UpdatedAt: evalTime, State: StateFiring},
			},
			map[string]string{},
			[]Fire{
				{Labels: map[string]string(nil), Annotations: map[string]string(nil)},
				{Labels: map[string]string(nil), Annotations: map[string]string(nil)},
			},
		},
		{
			map[uint64]*Alert{
				1: {UpdatedAt: evalTime, Labels: map[string]string{"hello": "world"}, State: StateFiring},
				2: {UpdatedAt: evalTime, Labels: map[string]string{"hello": "world"}, State: StateFiring},
			},
			map[string]string{},
			[]Fire{
				{Labels: map[string]string{"hello": "world"}, Annotations: map[string]string(nil)},
				{Labels: map[string]string{"hello": "world"}, Annotations: map[string]string(nil)},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			ar := &AlertingRule{Active: tc.active}
			fires := []Fire{}
			alertingService1.evalAlertingRule(ar, servers.GetDatasources(), tc.commonLabels, evalTime, &fires)
			require.Equal(t, tc.want, fires)
		})
	}
}

func TestEvalAlertingRuleSample(t *testing.T) {
	active := map[uint64]*Alert{}
	ar := AlertingRule{
		Active: active,
	}
	sample := commonModel.Sample{}
	labels := map[string]string{}
	evalTime := time.Now()
	alertingService1.evalAlertingRuleSample(&ar, sample, labels, evalTime)
	want := AlertingRule{
		Rule: model.Rule{
			Labels:      map[string]string(nil),
			Annotations: map[string]string(nil),
		},
		Active: active,
	}
	require.Equal(t, want, ar)
}

func TestRenderSummaryAnnotation(t *testing.T) {
	testCases := []struct {
		annotations map[string]string
		value       string
		labels      map[string]string
		want        map[string]string
		wantError   string
	}{
		// error
		{
			map[string]string{},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "placeholder summary"},
			`no summary annotation`,
		},
		{
			map[string]string{"summary": "{{$xxx}}"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "{{$xxx}}"},
			`parse err: template: :1: undefined variable "$xxx"`,
		},
		// ok
		{
			map[string]string{"summary": ""},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": ""},
			"",
		},
		{
			map[string]string{"summary": "hello"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "hello"},
			"",
		},
		{
			map[string]string{"summary": "{{$value}}"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "100"},
			"",
		},
		{
			map[string]string{"summary": "{{$labels}}"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "map[datasource:datasource1 foo:bar hello:world]"},
			"",
		},
		{
			map[string]string{"summary": "{{$labels.hello}}"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "world"},
			"",
		},
		{
			map[string]string{"summary": "{{$labels.xxx}}"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "<no value>"},
			"",
		},
		{
			map[string]string{"summary": "{{$}}"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "map[datasource:datasource1 foo:bar hello:world]"},
			"",
		},
		{
			map[string]string{"summary": "{{$.foo}}"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "bar"},
			"",
		},
		{
			map[string]string{"summary": "{{.}}"},
			"100", map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			map[string]string{"summary": "map[datasource:datasource1 foo:bar hello:world]"},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			annotations := tc.annotations
			err := renderSummaryAnnotaion(annotations, tc.labels, tc.value)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.want, annotations)
		})
	}
}

func TestQueryRule(t *testing.T) {
	testCases := []struct {
		rule      model.Rule
		ds        model.Datasource
		want      []commonModel.Sample
		wantError string
	}{
		{
			model.Rule{},
			model.Datasource{},
			[]commonModel.Sample{},
			`GET err: error on Do: Get "/api/v1/query?query=": unsupported protocol scheme ""`,
		},
		{
			ruleFiles1[0].RuleGroups[0].Rules[0],
			servers.GetDatasources()[3],
			[]commonModel.Sample{
				{Metric: commonModel.Metric{"__name__": "up", "instance": "localhost:9090", "job": "prometheus"}, Value: 1, Timestamp: 1435781451781},
				{Metric: commonModel.Metric{"__name__": "up", "instance2": "localhost:9092", "job": "prometheus2"}, Value: 1, Timestamp: 1435781451781}},
			``,
		},
		{
			ruleFiles1[0].RuleGroups[0].Rules[0],
			servers.GetDatasources()[1],
			[]commonModel.Sample{
				{Metric: commonModel.Metric{"__name__": "up", "instance": "localhost:6060", "job": "lethe"}, Value: 1, Timestamp: 1435781451781}},
			``,
		},
		{
			ruleFiles1[1].RuleGroups[0].Rules[0],
			servers.GetDatasources()[1],
			[]commonModel.Sample{{Value: 2}},
			``,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := alertingService1.queryRule(tc.rule, tc.ds)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, `GET err: error on Do: Get "/api/v1/query?query=": unsupported protocol scheme ""`)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetDataFromLogs(t *testing.T) {
	testCases := []struct {
		body string
		want []commonModel.Sample
	}{
		{
			`{}`,
			[]commonModel.Sample{{Value: 0}},
		},
		{
			`{"status":"success","data":{"resultType":"logs", "result":[
				{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"lerom ipsum"},
				{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"hello world"}]}}`,
			[]commonModel.Sample{{Value: 2}},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := getDataFromLogs([]byte(tc.body))
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetDataFromVector(t *testing.T) {
	// ok
	body := `{"status":"success","data":{"resultType":"vector","result":[
		{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]},
		{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`
	want := []commonModel.Sample{
		{Metric: commonModel.Metric{"__name__": "up", "instance": "localhost:9090", "job": "prometheus"}, Value: 1, Timestamp: 1435781451781},
		{Metric: commonModel.Metric{"__name__": "up", "instance": "localhost:9090", "job": "prometheus"}, Value: 1, Timestamp: 1435781451781}}
	got, err := getDataFromVector([]byte(body))
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestSendFires(t *testing.T) {
	fires := []Fire{
		{Labels: map[string]string{"test": "test", "severity": "info", "pizza": "🍕", "time": time.Now().String()}},
	}
	// ok
	err := alertingService1.sendFires(fires)
	require.NoError(t, err)
	// error
	temp := alertingService1.alertmanagerURL
	alertingService1.alertmanagerURL = ""
	err = alertingService1.sendFires(fires)
	require.EqualError(t, err, `error on Post: Post "/api/v1/alerts": unsupported protocol scheme ""`)
	alertingService1.alertmanagerURL = temp
}

func TestSendTestAlert(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := alertingService1.SendTestAlert()
		require.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		temp := alertingService1.alertmanagerURL
		alertingService1.alertmanagerURL = ""
		err := alertingService1.SendTestAlert()
		require.EqualError(t, err, `sendFires err: error on Post: Post "/api/v1/alerts": unsupported protocol scheme ""`)
		alertingService1.alertmanagerURL = temp
	})
}
