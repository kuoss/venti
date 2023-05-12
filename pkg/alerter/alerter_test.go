package alerter

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store/alerting"
	dsStore "github.com/kuoss/venti/pkg/store/datasource"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/kuoss/venti/pkg/store/remote"
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
		{Type: ms.TypeAlertmanager, Port: 0, Name: "alertmanager1", IsMain: false},
		{Type: ms.TypeLethe, Port: 0, Name: "lethe1", IsMain: true},
		{Type: ms.TypeLethe, Port: 0, Name: "lethe2", IsMain: false},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus1", IsMain: true},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus2", IsMain: false},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus3", IsMain: false},
	})
	datasourceConfig := &model.DatasourceConfig{
		Datasources: servers.GetDatasources(),
	}

	datasourceStore, err := dsStore.New(datasourceConfig, discovery.Discoverer(nil))
	if err != nil {
		panic(err)
	}

	remoteStore := remote.New(&http.Client{}, 30*time.Second)
	alertingStore := alerting.New("", ruleFiles1, datasourceStore)
	alerter1 = New(alertingStore, remoteStore)
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
	datasourceStore, err := dsStore.New(&model.DatasourceConfig{}, discovery.Discoverer(nil))
	require.NoError(t, err)
	alertingStore := alerting.New("", ruleFiles1, datasourceStore)
	remoteStore := remote.New(&http.Client{}, 30*time.Second)
	tempAlerter := New(alertingStore, remoteStore)

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
	remoteStore := remote.New(&http.Client{}, 30*time.Second)
	// tempAlerter_ok2.alertingStore.AlertFiles[0].AlertGroups[0].RuleAlerts[0].Rule = model.Rule{Alert: "alert1", Expr: "unmarshalable"}

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

			alertingStore := alerting.New("", ruleFiles1, &dsStore.DatasourceStore{})
			alerter := New(alertingStore, remoteStore)
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
	datasource5 := &model.Datasource{Type: model.DatasourceTypePrometheus, URL: servers.Svrs[5].Server.URL}
	testCases := []struct {
		ruleAlert *model.RuleAlert
		want      []model.Fire
	}{
		// ok
		{
			&model.RuleAlert{Rule: model.Rule{Alert: "alert1", Expr: "up"}, Alerts: []model.Alert{{Datasource: datasource5}}},
			[]model.Fire{{Labels: map[string]string{"alertname": "alert1", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}}},
		},
		// error
		{
			&model.RuleAlert{Rule: model.Rule{Alert: "alert1", Expr: "unmarshalable"}, Alerts: []model.Alert{{Datasource: datasource5}}},
			nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := alerter1.processRuleAlert(tc.ruleAlert)
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
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			queryData, err := alerter1.queryAlert(tc.rule, tc.alert)
			require.Equal(tt, tc.want, queryData)
			if tc.wantError == "" {
				require.NoError(tt, err)
			} else {
				require.EqualError(tt, err, tc.wantError)
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
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
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
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := evaluateAlert(tc.queryData, tc.rule, tc.alert)
			require.Equal(tt, tc.want, fires)
		})
	}
}

func TestGetFires_zero_QueryData(t *testing.T) {
	queryData := model.QueryData{}
	testCases := []struct {
		rule *model.Rule
		want []model.Fire
	}{
		{
			&model.Rule{},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
			},
		},
		{
			&model.Rule{
				Annotations: map[string]string{"hello": "world"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}},
			},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{ $value }}"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "alert1", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{ $value }}"}},
			},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{}}"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "alert1", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{}}"}},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := getFires(tc.rule, queryData)
			require.Equal(tt, tc.want, fires)
		})
	}
}

func TestGetFires_vector_zero_Result(t *testing.T) {
	queryData := model.QueryData{ResultType: commonModel.ValVector, Result: []commonModel.Sample{}}
	testCases := []struct {
		rule *model.Rule
		want []model.Fire
	}{
		{
			&model.Rule{},
			[]model.Fire{},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"apple": "banana"},
				Labels:      map[string]string{"lemon": "orange"},
			},
			[]model.Fire{},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := getFires(tc.rule, queryData)
			require.Equal(tt, tc.want, fires)
		})
	}
}

func TestGetFires_vector_two_Result(t *testing.T) {
	queryData := model.QueryData{ResultType: commonModel.ValVector, Result: []commonModel.Sample{
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
		rule *model.Rule
		want []model.Fire
	}{
		{
			&model.Rule{},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
			},
		},
		{
			&model.Rule{
				Annotations: map[string]string{"hello": "world"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}},
				{Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}},
			},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{ $value }}"},
				Labels:      map[string]string{"rule": "pod-v1"},
			},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "alert1", "firer": "venti", "rule": "pod-v1"}, Annotations: map[string]string{"summary": "pod=pod1 value=1111"}},
				{Labels: map[string]string{"alertname": "alert1", "firer": "venti", "rule": "pod-v1"}, Annotations: map[string]string{"summary": "pod=pod2 value=2222"}},
			},
		},
		{
			&model.Rule{
				Alert:       "alert1",
				Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}, // parse error
				Labels:      map[string]string{"rule": "pod-v1"},
			},
			[]model.Fire{
				{Labels: map[string]string{"alertname": "alert1", "firer": "venti", "rule": "pod-v1"}, Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}},
				{Labels: map[string]string{"alertname": "alert1", "firer": "venti", "rule": "pod-v1"}, Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := getFires(tc.rule, queryData)
			require.Equal(tt, tc.want, fires)
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
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			ret, err := renderSummary(tc.input, tc.sample)
			require.Nil(tt, err)
			require.Equal(tt, tc.want, ret)
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
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			output, err := renderSummary(tc.input, tc.sample)
			require.NotNil(tt, err)
			require.Error(tt, err, tc.wantError)
			require.Equal(tt, tc.input, output)
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
			`error on Parse: template: :1: missing value for command`,
		},
		{
			"AlwaysOn value={{ }",
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
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			output, err := renderSummary(tc.input, tc.sample)
			require.NotNil(tt, err)
			require.Error(tt, err, tc.wantError)
			require.Equal(tt, tc.input, output)
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
	require.Nil(t, err)
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
