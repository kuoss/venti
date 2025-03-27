package alerting

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/config"
	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	datasourceservice "github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/discovery"
	remoteservice "github.com/kuoss/venti/pkg/service/remote"
	commonmodel "github.com/prometheus/common/model"
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
	cfg := &config.Config{
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

func TestGetAlertingRuleGroups(t *testing.T) {
	got := alertingService1.GetAlertingRuleGroups()
	require.NotEmpty(t, got)
}

func TestGetAlertmanagerDiscovery(t *testing.T) {
	got := alertingService1.GetAlertmanagerDiscovery()
	require.NotEmpty(t, got)
}

func TestDoAlert(t *testing.T) {
	alertingRuleGroups := alertingService1.alertingRuleGroups
	datasourceService := alertingService1.datasourceService
	datasourceReload := alertingService1.datasourceReload
	alertmanagerURL := alertingService1.alertmanagerURL
	defer func() {
		alertingService1.alertingRuleGroups = alertingRuleGroups
		alertingService1.datasourceService = datasourceService
		alertingService1.datasourceReload = datasourceReload
		alertingService1.alertmanagerURL = alertmanagerURL
	}()

	testCases := []struct {
		alertingRuleGroups []AlertingRuleGroup
		datasourceService  datasourceservice.IDatasourceService
		datasourceReload   bool
		alertmanagerURL    string
		wantError          string
	}{
		{
			alertingRuleGroups, datasourceService, datasourceReload, alertmanagerURL,
			"",
		},
		{
			[]AlertingRuleGroup{}, datasourceService, datasourceReload, alertmanagerURL,
			"",
		},
		{
			[]AlertingRuleGroup{}, &mockDatasourceService{}, true, alertmanagerURL,
			"reload err: mock reload err",
		},
		{
			[]AlertingRuleGroup{}, datasourceService, datasourceReload, "",
			`sendFires err: post err: Post "/api/v2/alerts": unsupported protocol scheme ""`,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			alertingService1.alertingRuleGroups = tc.alertingRuleGroups
			alertingService1.datasourceService = tc.datasourceService
			alertingService1.datasourceReload = tc.datasourceReload
			alertingService1.alertmanagerURL = tc.alertmanagerURL
			err := alertingService1.DoAlert()
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
		})
	}
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
		{Labels: map[string]string{"alertname": "Pod", "datasource": "lethe1", "global1": "label1", "hello": "world", "rulefile": "sample-v3", "severity": "silence", "time": "2009-11-10T22:59:00.000000Z", "namespace": "namespace01", "pod": "nginx-deployment-75675f5897-7ci7o", "container": "nginx", "log": "hello world"}, Annotations: map[string]string{"summary": "PodLogs value=2"}},
		{Labels: map[string]string{"alertname": "Pod", "datasource": "lethe2", "global1": "label1", "hello": "world", "rulefile": "sample-v3", "severity": "silence", "time": "2009-11-10T22:59:00.000000Z", "namespace": "namespace01", "pod": "nginx-deployment-75675f5897-7ci7o", "container": "nginx", "log": "hello world"}, Annotations: map[string]string{"summary": "PodLogs value=2"}}}
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
		Rule: model.Rule{
			Annotations: map[string]string{"severity": "info"},
		},
	}
	sample := commonmodel.Sample{}
	labels := map[string]string{}
	evalTime := time.Now()

	want := AlertingRule{
		Rule: model.Rule{
			Labels:      map[string]string(nil),
			Annotations: map[string]string{"severity": "info"},
		},
		Active: active,
	}
	alertingService1.evalAlertingRuleSample(&ar, sample, labels, evalTime)
	require.Equal(t, want, ar)
}

func TestRenderSummaryAnnotation(t *testing.T) {
	defer func() {
		fakeErr1 = false
	}()
	testCases := []struct {
		annotations map[string]string
		labels      map[string]string
		value       string
		fakeErr1    bool
		want        map[string]string
		wantError   string
	}{
		// error
		{
			map[string]string{},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "placeholder summary"},
			`no summary annotation`,
		},
		{
			map[string]string{"summary": "{{$xxx}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "{{$xxx}}"},
			`parse err: template: :1: undefined variable "$xxx"`,
		},
		{
			map[string]string{"summary": "{{.}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", true,
			map[string]string{"summary": "{{.}}"},
			`tmpl.Execute err: %!w(<nil>)`,
		},
		// ok
		{
			map[string]string{"summary": ""},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": ""},
			``,
		},
		{
			map[string]string{"summary": "hello"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "hello"},
			``,
		},
		{
			map[string]string{"summary": "{{$value}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "100"},
			``,
		},
		{
			map[string]string{"summary": "{{$labels}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "map[datasource:datasource1 foo:bar hello:world]"},
			``,
		},
		{
			map[string]string{"summary": "{{$labels.hello}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "world"},
			``,
		},
		{
			map[string]string{"summary": "{{$labels.xxx}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "<no value>"},
			``,
		},
		{
			map[string]string{"summary": "{{$}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "map[datasource:datasource1 foo:bar hello:world]"},
			``,
		},
		{
			map[string]string{"summary": "{{$.foo}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "bar"},
			``,
		},
		{
			map[string]string{"summary": "{{.}}"},
			map[string]string{"datasource": "datasource1", "hello": "world", "foo": "bar"},
			"100", false,
			map[string]string{"summary": "map[datasource:datasource1 foo:bar hello:world]"},
			``,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			annotations := tc.annotations
			fakeErr1 = tc.fakeErr1
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
		want      []commonmodel.Sample
		wantError string
	}{
		{
			model.Rule{},
			model.Datasource{},
			[]commonmodel.Sample{},
			`GET err: error on Do: Get "/api/v1/query?query=": unsupported protocol scheme ""`,
		},
		{
			ruleFiles1[0].RuleGroups[0].Rules[0],
			servers.GetDatasources()[3],
			[]commonmodel.Sample{
				{Metric: commonmodel.Metric{"__name__": "up", "instance": "localhost:9090", "job": "prometheus"}, Value: 1, Timestamp: 1435781451781},
				{Metric: commonmodel.Metric{"__name__": "up", "instance2": "localhost:9092", "job": "prometheus2"}, Value: 1, Timestamp: 1435781451781}},
			``,
		},
		{
			ruleFiles1[0].RuleGroups[0].Rules[0],
			servers.GetDatasources()[1],
			[]commonmodel.Sample{
				{Metric: commonmodel.Metric{"__name__": "up", "instance": "localhost:6060", "job": "lethe"}, Value: 1, Timestamp: 1435781451781}},
			``,
		},
		{
			ruleFiles1[1].RuleGroups[0].Rules[0],
			servers.GetDatasources()[1],
			[]commonmodel.Sample{{Metric: commonmodel.Metric{"container": "nginx", "log": "hello world", "namespace": "namespace01", "pod": "nginx-deployment-75675f5897-7ci7o", "time": "2009-11-10T22:59:00.000000Z"}, Value: 2}},
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
		want []commonmodel.Sample
	}{
		{
			`{"status":"success","data":{"resultType":"logs", "result":[]}}`,
			[]commonmodel.Sample{},
		},
		{
			`{"status":"success","data":{"resultType":"logs", "result":[
				{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"lerom ipsum"},
				{"time":"2009-11-10T22:59:00.000000Z","namespace":"namespace01","pod":"nginx-deployment-75675f5897-7ci7o","container":"nginx","log":"hello world"}]}}`,
			[]commonmodel.Sample{{Metric: commonmodel.Metric{"container": "nginx", "log": "hello world", "namespace": "namespace01", "pod": "nginx-deployment-75675f5897-7ci7o", "time": "2009-11-10T22:59:00.000000Z"}, Value: 2}},
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
	defer func() {
		fakeErr1 = false
	}()

	body := `{"status":"success","data":{"resultType":"vector","result":[
		{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]},
		{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`
	testCases := []struct {
		body      string
		fakeErr1  bool
		want      []commonmodel.Sample
		wantError string
	}{
		{
			body, false,
			[]commonmodel.Sample{
				{Metric: commonmodel.Metric{"__name__": "up", "instance": "localhost:9090", "job": "prometheus"}, Value: 1, Timestamp: 1435781451781},
				{Metric: commonmodel.Metric{"__name__": "up", "instance": "localhost:9090", "job": "prometheus"}, Value: 1, Timestamp: 1435781451781}},
			``,
		},
		{
			body, true,
			[]commonmodel.Sample{},
			`unmarshal err: %!w(<nil>)`,
		},
	}
	for _, tc := range testCases {
		fakeErr1 = tc.fakeErr1
		got, err := getDataFromVector([]byte(body))
		if tc.wantError == "" {
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		} else {
			require.EqualError(t, err, tc.wantError)
			require.Equal(t, tc.want, got)
		}
	}

}

func TestSendFires(t *testing.T) {
	alertmanagerURL := alertingService1.alertmanagerURL
	defer func() {
		alertingService1.alertmanagerURL = alertmanagerURL
		fakeErr1 = false
		fakeErr2 = false
	}()
	fires1 := []Fire{
		{Labels: map[string]string{"test": "test", "severity": "info", "pizza": "🍕", "time": time.Now().String()}},
	}
	testCases := []struct {
		fires           []Fire
		alertmanagerURL string
		fakeErr1        bool
		fakeErr2        bool
		wantError       string
	}{
		{
			fires1, alertmanagerURL, false, false,
			``,
		},
		{
			fires1, "", false, false,
			`post err: Post "/api/v2/alerts": unsupported protocol scheme ""`,
		},
		{
			fires1, alertmanagerURL, true, false,
			`marshal err: %!w(<nil>)`,
		},
		{
			fires1, alertmanagerURL, false, true,
			`statusCode is not ok(200)`,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			alertingService1.alertmanagerURL = tc.alertmanagerURL
			fakeErr1 = tc.fakeErr1
			fakeErr2 = tc.fakeErr2
			err := alertingService1.sendFires(tc.fires)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
		})
	}
}

func TestSendTestAlert(t *testing.T) {
	alertmanagerURL := alertingService1.alertmanagerURL
	defer func() {
		alertingService1.alertmanagerURL = alertmanagerURL
	}()
	testCases := []struct {
		alertmanagerURL string
		wantError       string
	}{
		{alertmanagerURL, ``},
		{"", `sendFires err: post err: Post "/api/v2/alerts": unsupported protocol scheme ""`},
	}
	for _, tc := range testCases {
		alertingService1.alertmanagerURL = tc.alertmanagerURL
		err := alertingService1.SendTestAlert()
		if tc.wantError == "" {
			require.NoError(t, err)
		} else {
			require.EqualError(t, err, tc.wantError)
		}
	}
}
