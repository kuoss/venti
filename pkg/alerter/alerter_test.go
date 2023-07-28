package alerter

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/alerting"
	dsService "github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/discovery"
	"github.com/kuoss/venti/pkg/service/remote"
	commonModel "github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
)

var (
	alerter1   *alerter
	servers    *ms.Servers
	ruleFiles1 []model.RuleFile = []model.RuleFile{
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
	datasourceConfig := &model.DatasourceConfig{
		Datasources: servers.GetDatasources(),
	}

	datasourceService, err := dsService.New(datasourceConfig, discovery.Discoverer(nil))
	if err != nil {
		panic(err)
	}

	remoteService := remote.New(&http.Client{}, 30*time.Second)
	alertingService := alerting.New("", ruleFiles1, datasourceService)
	alerter1 = New(alertingService, remoteService)
	alerter1.SetAlertmanagerURL(servers.GetServersByType(ms.TypeAlertmanager)[0].URL)
}

func TestNew(t *testing.T) {

	require.Equal(t, 1, len(ruleFiles1))
	require.Equal(t, 1, len(ruleFiles1[0].RuleGroups))
	require.Equal(t, 3, len(ruleFiles1[0].RuleGroups[0].Rules))
	require.Equal(t, model.DatasourceSelector{Type: "prometheus"}, ruleFiles1[0].DatasourceSelector)
	require.Equal(t, false, alerter1.repeat)
}

func TestSetAlertmanagerURL(t *testing.T) {
	temp := alerter1.alertmanagerURL
	alerter1.SetAlertmanagerURL("hello")
	require.Equal(t, "hello", alerter1.alertmanagerURL)
	alerter1.SetAlertmanagerURL(temp)
}

func TestStartAndStop(t *testing.T) {
	datasourceService, err := dsService.New(&model.DatasourceConfig{}, discovery.Discoverer(nil))
	require.NoError(t, err)
	alertingService := alerting.New("", ruleFiles1, datasourceService)
	remoteService := remote.New(&http.Client{}, 30*time.Second)
	tempAlerter := New(alertingService, remoteService)

	tempAlerter.SetAlertmanagerURL(servers.Svrs[0].Server.URL)
	tempAlerter.evaluationInterval = 1000 * time.Millisecond

	err = tempAlerter.Start()
	require.NoError(t, err)

	// 'Stop' works fine. but panic in race condition test ( go test ./... -v -failfast -race )
	// tempAlerter.Stop()

}

func TestOnce(t *testing.T) {
	// alerter1
	alerter1.Once()

	// empty alerter
	tempAlerter := alerter{}
	tempAlerter.Once()
}

func TestProcessAlertFiles(t *testing.T) {
	remoteService := remote.New(&http.Client{}, 30*time.Second)
	// tempAlerter_ok2.alertingService.AlertFiles[0].AlertGroups[0].RuleAlerts[0].Rule = model.Rule{Alert: "alert1", Expr: "unmarshalable"}

	testCases := []struct {
		alertmanagerURL   string
		wantErrorContains string
	}{
		// ok
		{
			servers.Svrs[0].Server.URL,
			"",
		},
		{
			servers.Svrs[0].Server.URL,
			"",
		},
		// error
		{
			"",
			`sendFires err: error on Post: Post "/api/v1/alerts": unsupported protocol scheme ""`,
		},
		{
			"http://alertmanager:9093",
			"sendFires err: error on Post: Post \"http://alertmanager:9093/api/v1/alerts\": dial tcp: lookup alertmanager on",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {

			alertingService := alerting.New("", ruleFiles1, &dsService.DatasourceService{})
			alerter := New(alertingService, remoteService)
			alerter.SetAlertmanagerURL(tc.alertmanagerURL)

			err := alerter.processAlertFiles()
			if tc.wantErrorContains == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErrorContains)
			}
		})
	}
}

func TestProcessRuleAlert(t *testing.T) {
	datasource5 := &model.Datasource{Name: "prom", Type: model.DatasourceTypePrometheus, URL: servers.Svrs[5].Server.URL}
	testCases := []struct {
		ruleAlert *model.RuleAlert
		want      []model.Fire
	}{
		// ok
		{
			&model.RuleAlert{Rule: model.Rule{Alert: "alert1", Expr: "up"}, Alerts: []model.Alert{{Datasource: datasource5}}},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "alert1", "datasource": "prom", "firer": "venti", "severity": "info"}, Annotations: map[string]string{"summary": "placeholder summary"}}},
		},
		// error
		{
			&model.RuleAlert{Rule: model.Rule{Alert: "alert1", Expr: "unmarshalable"}, Alerts: []model.Alert{{Datasource: datasource5}}},
			nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := alerter1.processRuleAlert(tc.ruleAlert, &map[string]string{"severity": "info"})
			require.Equal(t, tc.want, got)
		})
	}
}

func TestQueryAlert(t *testing.T) {
	testCases := []struct {
		rule      *model.Rule
		alert     *model.Alert
		want      model.QueryData
		wantError string
	}{
		// ok
		{
			&model.Rule{Alert: "alert1", Expr: "up"},
			&model.Alert{Datasource: &model.Datasource{URL: servers.Svrs[5].Server.URL}},
			model.QueryData{ResultType: 2, Result: []commonModel.Sample{
				{Metric: commonModel.Metric{"__name__": "up", "instance": "localhost:9090", "job": "prometheus"}, Value: 1, Timestamp: 1435781451781, Histogram: (*commonModel.SampleHistogram)(nil)},
			}},
			"",
		},
		// error
		{
			&model.Rule{Alert: "alert1", Expr: "unmarshalable"},
			&model.Alert{Datasource: &model.Datasource{URL: servers.Svrs[5].Server.URL}},
			model.QueryData{},
			"unmarshal err: invalid character ':' after object key:value pair, body: {\"status\":\"success\",\"data\":{\"a\":\"b\":\"c\"}}",
		},
		{
			&model.Rule{},
			&model.Alert{},
			model.QueryData{},
			"datasource is nil",
		},
		{
			&model.Rule{},
			&model.Alert{Datasource: &model.Datasource{URL: servers.Svrs[4].Server.URL}}, // 401 unauthroized (basicAuth)
			model.QueryData{},
			"not success status (status=error, code=405)",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := alerter1.queryAlert(tc.rule, tc.alert)
			require.Equal(t, tc.want, got)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
		})
	}
}

func TestEvaluateAlert(t *testing.T) {
	queryData_two_samples := model.QueryData{ResultType: commonModel.ValVector, Result: []commonModel.Sample{
		{
			Metric:    map[commonModel.LabelName]commonModel.LabelValue{"pod": "pod1"},
			Value:     1111,
			Timestamp: 0,
			Histogram: &commonModel.SampleHistogram{},
		},
		{
			Metric:    map[commonModel.LabelName]commonModel.LabelValue{"pod": "pod2"},
			Value:     2222,
			Timestamp: 0,
			Histogram: &commonModel.SampleHistogram{},
		},
	}}
	testCases := []struct {
		queryData model.QueryData
		rule      *model.Rule
		alert     *model.Alert
		want      []model.Fire
	}{
		{
			model.QueryData{},
			&model.Rule{},
			&model.Alert{},
			[]model.Fire(nil),
		},
		// firing
		{
			queryData_two_samples,
			&model.Rule{},
			&model.Alert{},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "severity": "info"}, Annotations: map[string]string{"summary": "placeholder summary"}},
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "severity": "info"}, Annotations: map[string]string{"summary": "placeholder summary"}},
			},
		},
		// pending
		{
			queryData_two_samples,
			&model.Rule{For: commonModel.Duration(5 * time.Minute)},                           // 5 min --> pending
			&model.Alert{ActiveAt: commonModel.Time(commonModel.Now().Add(-1 * time.Minute))}, // 1 min ago
			[]model.Fire(nil),
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got := evaluateAlert(tc.queryData, tc.rule, tc.alert, &map[string]string{"severity": "info"})
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetFires(t *testing.T) {
	queryData1 := model.QueryData{ResultType: commonModel.ValVector, Result: []commonModel.Sample{
		{
			Metric:    map[commonModel.LabelName]commonModel.LabelValue{"pod": "pod1"},
			Value:     1111,
			Timestamp: 0,
			Histogram: &commonModel.SampleHistogram{},
		},
		{
			Metric:    map[commonModel.LabelName]commonModel.LabelValue{"pod": "pod2"},
			Value:     2222,
			Timestamp: 0,
			Histogram: &commonModel.SampleHistogram{},
		},
	}}

	testCases := []struct {
		rule         *model.Rule
		queryData    model.QueryData
		commonLabels *map[string]string
		datasource   *model.Datasource
		want         []model.Fire
	}{
		{
			&model.Rule{},
			model.QueryData{ResultType: commonModel.ValVector, Result: []commonModel.Sample{}},
			&map[string]string{"severity": "info"},
			&model.Datasource{},
			[]model.Fire{},
		},
		{
			&model.Rule{},
			model.QueryData{},
			&map[string]string{"severity": "info"},
			nil,
			[]model.Fire{{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "severity": "info"}, Annotations: map[string]string{"summary": "placeholder summary"}}},
		},
		{
			&model.Rule{},
			model.QueryData{},
			&map[string]string{"severity": "info"},
			&model.Datasource{Name: "temp-datasource"},
			[]model.Fire{{Labels: map[string]string{"alertname": "placeholder name", "datasource": "temp-datasource", "firer": "venti", "severity": "info"}, Annotations: map[string]string{"summary": "placeholder summary"}}},
		},
		{
			&model.Rule{},
			queryData1,
			&map[string]string{"severity": "info"},
			&model.Datasource{Name: "temp-datasource"},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "placeholder name", "datasource": "temp-datasource", "firer": "venti", "severity": "info"}, Annotations: map[string]string{"summary": "placeholder summary"}},
				{Labels: map[string]string{"alertname": "placeholder name", "datasource": "temp-datasource", "firer": "venti", "severity": "info"}, Annotations: map[string]string{"summary": "placeholder summary"}},
			},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"apple": "banana"},
				Labels:      map[string]string{"lemon": "orange"},
			},
			model.QueryData{ResultType: commonModel.ValVector, Result: []commonModel.Sample{}},
			&map[string]string{"severity": "info"},
			&model.Datasource{},
			[]model.Fire{},
		},
		{
			&model.Rule{
				Annotations: map[string]string{"hello": "world"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			model.QueryData{},
			&map[string]string{"severity": "info"},
			&model.Datasource{Name: "temp-datasource"},
			[]model.Fire{{Labels: map[string]string{"alertname": "placeholder name", "datasource": "temp-datasource", "firer": "venti", "lorem": "ipsum", "severity": "info"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}}},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{ $value }}"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			model.QueryData{},
			&map[string]string{"severity": "info"},
			&model.Datasource{Name: "temp-datasource"},
			[]model.Fire{{Labels: map[string]string{"alertname": "alert1", "datasource": "temp-datasource", "firer": "venti", "lorem": "ipsum", "severity": "info"}, Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{ $value }}"}}}},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{}}"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			model.QueryData{},
			&map[string]string{"severity": "info"},
			&model.Datasource{Name: "temp-datasource"},
			[]model.Fire{{Labels: map[string]string{"alertname": "alert1", "datasource": "temp-datasource", "firer": "venti", "lorem": "ipsum", "severity": "info"}, Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{}}"}}},
		},
		{
			&model.Rule{
				Annotations: map[string]string{"hello": "world"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			queryData1,
			&map[string]string{"severity": "info"},
			&model.Datasource{Name: "temp-datasource"},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "placeholder name", "datasource": "temp-datasource", "firer": "venti", "lorem": "ipsum", "severity": "info"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}},
				{Labels: map[string]string{"alertname": "placeholder name", "datasource": "temp-datasource", "firer": "venti", "lorem": "ipsum", "severity": "info"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}},
			},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{ $value }}"},
				Labels:      map[string]string{"rule": "pod-v1"},
			},
			queryData1,
			&map[string]string{"severity": "info"},
			&model.Datasource{Name: "temp-datasource"},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "alert1", "datasource": "temp-datasource", "firer": "venti", "rule": "pod-v1", "severity": "info"}, Annotations: map[string]string{"summary": "pod=pod1 value=1111"}},
				{Labels: map[string]string{"alertname": "alert1", "datasource": "temp-datasource", "firer": "venti", "rule": "pod-v1", "severity": "info"}, Annotations: map[string]string{"summary": "pod=pod2 value=2222"}},
			},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}, // parse error
				Labels:      map[string]string{"rule": "pod-v1"},
			},
			queryData1,
			&map[string]string{"severity": "info"},
			&model.Datasource{Name: "temp-datasource"},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "alert1", "datasource": "temp-datasource", "firer": "venti", "rule": "pod-v1", "severity": "info"}, Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}},
				{Labels: map[string]string{"alertname": "alert1", "datasource": "temp-datasource", "firer": "venti", "rule": "pod-v1", "severity": "info"}, Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got := getFires(tc.rule, tc.queryData, tc.commonLabels, tc.datasource)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRenderSummary_ok(t *testing.T) {
	testCases := []struct {
		input  string
		sample *commonModel.Sample
		want   string
	}{
		{
			"AlwaysOn value={{ $value }}",
			&commonModel.Sample{},
			"AlwaysOn value=0",
		},
		{
			"AlwaysOn value={{ $value }}",
			&commonModel.Sample{Value: 42},
			"AlwaysOn value=42",
		},
		{
			"Monday",
			&commonModel.Sample{},
			"Monday",
		},
		{
			"namespace={{ $labels.namespace }} pizza={{ $labels.pizza }}",
			&commonModel.Sample{Value: 42, Metric: commonModel.Metric{"namespace": "ns1", "pizza": "🍕"}},
			"namespace=ns1 pizza=🍕",
		},
		{
			"labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}",
			&commonModel.Sample{Value: 42, Metric: commonModel.Metric{"hello": "world", "lorum": "ipsum", "namespace": "ns1"}},
			"labels=map[hello:world lorum:ipsum namespace:ns1] namespace=ns1 value=42",
		},
		{
			"labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}",
			&commonModel.Sample{Value: 42},
			"labels=map[] namespace= value=42",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := renderSummary(tc.input, tc.sample)
			require.Nil(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRenderSummary_error_on_Parse(t *testing.T) {
	testCases := []struct {
		input     string
		sample    *commonModel.Sample
		wantError string
	}{
		{
			"AlwaysOn value={{}}",
			&commonModel.Sample{},
			`error on Parse: template: :1: missing value for command`,
		},
		{
			"AlwaysOn value={{ }}",
			&commonModel.Sample{},
			`error on Parse: template: :1: missing value for command`,
		},
		{
			"AlwaysOn value={{ }}",
			&commonModel.Sample{Value: 111},
			`error on Parse: template: :1: missing value for command`,
		},
		{
			"AlwaysOn value={{ $notexist }}",
			&commonModel.Sample{},
			`error on Parse: template: :1: undefined variable "$notexist"`,
		},
		{
			"AlwaysOn value={{ $notexist }}",
			&commonModel.Sample{Value: 111},
			`error on Parse: template: :1: undefined variable "$notexist"`,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := renderSummary(tc.input, tc.sample)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.input, got)
		})
	}
}

func TestRenderSummary_error_on_Execute(t *testing.T) {
	testCases := []struct {
		input     string
		sample    *commonModel.Sample
		wantError string
	}{
		{
			"AlwaysOn value={{}",
			&commonModel.Sample{},
			`error on Parse: template: :1: unexpected "}" in command`,
		},
		{
			"AlwaysOn value={{ }",
			&commonModel.Sample{},
			`error on Parse: template: :1: unexpected "}" in command`,
		},
		{
			"AlwaysOn value={{ }}",
			&commonModel.Sample{Value: 111},
			`error on Parse: template: :1: missing value for command`,
		},
		{
			"AlwaysOn value={{ $notexist }}",
			&commonModel.Sample{},
			`error on Parse: template: :1: undefined variable "$notexist"`,
		},
		{
			"AlwaysOn value={{ $notexist }}",
			&commonModel.Sample{Value: 111},
			`error on Parse: template: :1: undefined variable "$notexist"`,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := renderSummary(tc.input, tc.sample)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.input, got)
		})
	}
}

func TestSendFires_ok(t *testing.T) {
	err := alerter1.sendFires([]model.Fire{
		{
			Labels:      map[string]string{"hello": "world"},
			Annotations: map[string]string{"lorem": "ipsum"},
		},
	})
	require.NoError(t, err)
}

func TestSendFires_error_on_Post(t *testing.T) {
	tempURL := alerter1.alertmanagerURL
	alerter1.alertmanagerURL = ""
	err := alerter1.sendFires([]model.Fire{
		{
			Labels:      map[string]string{"hello": "world"},
			Annotations: map[string]string{"lorem": "ipsum"},
		},
	})
	require.NotNil(t, err)
	require.Equal(t, `error on Post: Post "/api/v1/alerts": unsupported protocol scheme ""`, err.Error())
	alerter1.alertmanagerURL = tempURL
}

func TestSendFires_not_ok(t *testing.T) {
	tempURL := alerter1.alertmanagerURL
	// prometheus insted alertmanager
	alerter1.alertmanagerURL = servers.GetServersByType(ms.TypePrometheus)[0].URL
	err := alerter1.sendFires([]model.Fire{
		{
			Labels:      map[string]string{"hello": "world"},
			Annotations: map[string]string{"lorem": "ipsum"},
		},
	})
	require.NotNil(t, err)
	require.Equal(t, `statusCode is not ok(200)`, err.Error())
	alerter1.alertmanagerURL = tempURL

}
