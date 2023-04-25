package alerter

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"github.com/kuoss/venti/pkg/store/alertrule"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/kuoss/venti/pkg/store/remote"
	commonModel "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

var (
	stores   *store.Stores
	alerter1 *alerter
	servers  *ms.Servers
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	_ = os.Chdir("../..")
	servers = ms.New(ms.Requirements{
		{Type: ms.TypeAlertmanager, Port: 0, Name: "alertmanager1", IsMain: false},
		{Type: ms.TypeLethe, Port: 0, Name: "lethe1", IsMain: true},
		{Type: ms.TypeLethe, Port: 0, Name: "lethe2", IsMain: false},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus1", IsMain: true},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus2", IsMain: false},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus3", IsMain: false},
	})
	datasourceConfig := model.DatasourceConfig{
		Datasources: servers.GetDatasources(),
	}

	var discoverer discovery.Discoverer
	datasourceStore, _ := store.NewDatasourceStore(&datasourceConfig, discoverer)
	remoteStore := remote.New(&http.Client{}, 30*time.Second)
	alertRuleStore, _ := alertrule.New("etc/alertrules/*.yaml")
	stores = &store.Stores{
		AlertRuleStore:  alertRuleStore,
		DatasourceStore: datasourceStore,
		RemoteStore:     remoteStore,
	}

	alerter1 = New(stores)
	alerter1.SetAlertmanagerURL(servers.GetServersByType(ms.TypeAlertmanager)[0].URL)
}

func shutdown() {
	servers.Close()
}

func TestNew(t *testing.T) {
	alertRuleFiles := []model.RuleFile{
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
	assert.Equal(t, 1, len(alertRuleFiles))
	assert.Equal(t, 1, len(alertRuleFiles[0].RuleGroups))
	assert.Equal(t, 3, len(alertRuleFiles[0].RuleGroups[0].Rules))
	assert.Equal(t, model.DatasourceSelector{Type: "prometheus"}, alertRuleFiles[0].DatasourceSelector)

	assert.Equal(t, false, alerter1.repeat)
	assert.Equal(t, 3, len(alerter1.alertFiles))
	assert.Equal(t, "prometheus1", alerter1.alertFiles[0].Datasource.Name)
	assert.Equal(t, "prometheus2", alerter1.alertFiles[1].Datasource.Name)
	assert.Equal(t, "prometheus3", alerter1.alertFiles[2].Datasource.Name)
}

func TestGetAlertFiles(t *testing.T) {
	var discoverer discovery.Discoverer
	alertRuleStore, _ := alertrule.New("etc/alertrules/*.yaml")

	// stores1
	datasourceStore1, _ := store.NewDatasourceStore(&model.DatasourceConfig{}, discoverer)
	stores1 := &store.Stores{AlertRuleStore: alertRuleStore, DatasourceStore: datasourceStore1}

	// stores2
	datasourceStore2, _ := store.NewDatasourceStore(&model.DatasourceConfig{Datasources: []*model.Datasource{
		{Type: model.DatasourceTypePrometheus, URL: servers.Svrs[0].Server.URL},
	}}, discoverer)
	stores2 := &store.Stores{AlertRuleStore: alertRuleStore, DatasourceStore: datasourceStore2}

	// stores3
	datasourceStore3, _ := store.NewDatasourceStore(&model.DatasourceConfig{Datasources: []*model.Datasource{
		{Type: model.DatasourceTypePrometheus, URL: servers.Svrs[0].Server.URL},
		{Type: model.DatasourceTypePrometheus, URL: servers.Svrs[1].Server.URL},
	}}, discoverer)
	stores3 := &store.Stores{AlertRuleStore: alertRuleStore, DatasourceStore: datasourceStore3}

	testCases := []struct {
		stores *store.Stores
		want   []model.AlertFile
	}{
		{
			stores1,
			[]model.AlertFile{},
		},
		{
			stores2,
			[]model.AlertFile{{
				Datasource: model.Datasource{Type: "prometheus", Name: "", URL: "", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
				AlertGroups: []model.AlertGroup{
					{Alerts: []model.Alert{
						{State: 0, Name: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}, ActiveAt: 0},
						{State: 0, Name: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Monday"}, ActiveAt: 0},
						{State: 0, Name: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}, ActiveAt: 0}}}},
			}},
		},
		{
			stores3,
			[]model.AlertFile{
				{
					Datasource: model.Datasource{Type: "prometheus", Name: "", URL: "", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
					AlertGroups: []model.AlertGroup{
						{Alerts: []model.Alert{
							{State: 0, Name: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}, ActiveAt: 0},
							{State: 0, Name: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Monday"}, ActiveAt: 0},
							{State: 0, Name: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}, ActiveAt: 0}}}}},
				{
					Datasource: model.Datasource{Type: "prometheus", Name: "", URL: "", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
					AlertGroups: []model.AlertGroup{
						{Alerts: []model.Alert{
							{State: 0, Name: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}, ActiveAt: 0},
							{State: 0, Name: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "Monday"}, ActiveAt: 0},
							{State: 0, Name: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, Labels: map[string]string{"hello": "world", "rulefile": "sample-v3", "severity": "silence"}, Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}, ActiveAt: 0}}}}}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			alertFiles := getAlertFiles(tc.stores)
			for i := range alertFiles {
				alertFiles[i].Datasource.URL = ""
			}
			assert.Equal(tt, tc.want, alertFiles)
		})
	}
}

func TestSetAlertmanagerURL(t *testing.T) {
	temp := alerter1.alertmanagerURL
	alerter1.SetAlertmanagerURL("hello")
	assert.Equal(t, "hello", alerter1.alertmanagerURL)
	alerter1.SetAlertmanagerURL(temp)
}

func TestStartAndStop(t *testing.T) {
	var discoverer discovery.Discoverer
	alertRuleStore, _ := alertrule.New("etc/alertrules/*.yaml")
	remoteStore := remote.New(&http.Client{}, 30*time.Second)

	datasourceStore, _ := store.NewDatasourceStore(&model.DatasourceConfig{}, discoverer)
	tempAlerter := New(&store.Stores{AlertRuleStore: alertRuleStore, DatasourceStore: datasourceStore, RemoteStore: remoteStore})
	tempAlerter.SetAlertmanagerURL(servers.Svrs[0].Server.URL)
	tempAlerter.evaluationInterval = 1000 * time.Millisecond

	tempAlerter.Start()

	// this works fine. but panic in race condition test ( go test ./... -v -failfast -race )
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
	var discoverer discovery.Discoverer
	alertRuleStore, _ := alertrule.New("etc/alertrules/*.yaml")
	remoteStore := remote.New(&http.Client{}, 30*time.Second)

	// tempAlerter_ok1
	datasourceStore_ok1, _ := store.NewDatasourceStore(&model.DatasourceConfig{
		Datasources: []*model.Datasource{{Type: model.DatasourceTypePrometheus, URL: servers.Svrs[2].Server.URL}}}, discoverer)
	tempAlerter_ok1 := New(&store.Stores{AlertRuleStore: alertRuleStore, DatasourceStore: datasourceStore_ok1, RemoteStore: remoteStore})
	tempAlerter_ok1.SetAlertmanagerURL(servers.Svrs[0].Server.URL)

	// tempAlerter_ok2
	datasourceStore_ok2, _ := store.NewDatasourceStore(&model.DatasourceConfig{
		Datasources: []*model.Datasource{{Type: model.DatasourceTypePrometheus, URL: servers.Svrs[2].Server.URL}}}, discoverer)
	tempAlerter_ok2 := New(&store.Stores{AlertRuleStore: alertRuleStore, DatasourceStore: datasourceStore_ok2, RemoteStore: remoteStore})
	tempAlerter_ok2.SetAlertmanagerURL(servers.Svrs[0].Server.URL)
	tempAlerter_ok2.alertFiles[0].AlertGroups[0].Alerts[0] = model.Alert{Name: "alert1", Expr: "unmarshalable"}

	// tempAlerter_error1
	tempAlerter_error1 := &alerter{}

	// tempAlerter_error2
	datasourceStore_error2, _ := store.NewDatasourceStore(&model.DatasourceConfig{}, discoverer)
	tempAlerter_error2 := New(&store.Stores{AlertRuleStore: alertRuleStore, DatasourceStore: datasourceStore_error2, RemoteStore: remoteStore})

	testCases := []struct {
		alerter         *alerter
		wantErrorRegexp string
	}{
		// ok
		{
			alerter1,
			"",
		},
		{
			tempAlerter_ok1,
			"",
		},
		{
			tempAlerter_ok2,
			"",
		},
		// error
		{
			tempAlerter_error1,
			`error on sendFires: error on Post: Post "/api/v1/alerts": unsupported protocol scheme ""`,
		},
		{
			tempAlerter_error2,
			"error on sendFires: error on Post: Post \"http://alertmanager:9093/api/v1/alerts\": dial tcp: lookup alertmanager on [0-9.]+:53: no such host",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			err := tc.alerter.processAlertFiles()
			if tc.wantErrorRegexp == "" {
				assert.NoError(tt, err)
			} else {
				assert.Regexp(tt, tc.wantErrorRegexp, err.Error())
			}
		})
	}
}

func TestProcessAlert(t *testing.T) {
	datasource5 := &model.Datasource{Type: model.DatasourceTypePrometheus, URL: servers.Svrs[5].Server.URL}
	testCases := []struct {
		alert      *model.Alert
		datasource *model.Datasource
		want       []model.Fire
		wantError  string
	}{
		// ok
		{
			&model.Alert{Name: "alert1", Expr: "up"},
			datasource5,
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "alert1", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
			},
			"",
		},
		// error
		{
			&model.Alert{Name: "alert1", Expr: "unmarshalable"},
			datasource5,
			nil,
			"error on queryAlert: error on Unmarshal: invalid character ':' after object key:value pair",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			queryData, err := alerter1.processAlert(tc.alert, tc.datasource)
			assert.Equal(tt, tc.want, queryData)
			if tc.wantError == "" {
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, tc.wantError)
			}
		})
	}
}

func TestQueryAlert(t *testing.T) {
	testCases := []struct {
		alert      *model.Alert
		datasource *model.Datasource
		want       model.QueryData
		wantError  string
	}{
		// ok
		{
			&model.Alert{Name: "alert1", Expr: "up"},
			&model.Datasource{URL: servers.Svrs[5].Server.URL},
			model.QueryData{ResultType: 2, Result: []commonModel.Sample{
				{Metric: commonModel.Metric{"__name__": "up", "instance": "localhost:9090", "job": "prometheus"}, Value: 1, Timestamp: 1435781451781, Histogram: (*commonModel.SampleHistogram)(nil)},
			}},
			"",
		},
		// error
		{
			&model.Alert{Name: "alert1", Expr: "unmarshalable"},
			&model.Datasource{URL: servers.Svrs[5].Server.URL},
			model.QueryData{},
			"error on Unmarshal: invalid character ':' after object key:value pair",
		},
		{
			&model.Alert{},
			&model.Datasource{},
			model.QueryData{},
			"error on GET: error on Do: Get \"/api/v1/query?query=\": unsupported protocol scheme \"\"",
		},
		{
			&model.Alert{},
			&model.Datasource{URL: servers.Svrs[4].Server.URL}, // 401 unauthroized (basicAuth)
			model.QueryData{},
			"not success status (status=error, code=405)",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			queryData, err := alerter1.queryAlert(tc.alert, tc.datasource)
			assert.Equal(tt, tc.want, queryData)
			if tc.wantError == "" {
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, tc.wantError)
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
		alert     *model.Alert
		queryData model.QueryData
		want      []model.Fire
	}{
		{
			&model.Alert{},
			model.QueryData{},
			[]model.Fire(nil),
		},
		// firing
		{
			&model.Alert{},
			queryData_two_samples,
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
				{State: "firing", Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
			},
		},
		// pending
		{
			&model.Alert{
				ActiveAt: commonModel.Time(commonModel.Now().Add(-1 * time.Minute)), // 1 min ago
				For:      commonModel.Duration(5 * time.Minute),                     // 5 min --> pending
			},
			queryData_two_samples,
			[]model.Fire(nil),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := evaluateAlert(tc.alert, tc.queryData)
			assert.Equal(tt, tc.want, fires)
		})
	}
}

func TestGetFires_zero_QueryData(t *testing.T) {
	queryData := model.QueryData{}
	testCases := []struct {
		alert *model.Alert
		want  []model.Fire
	}{
		{
			&model.Alert{},
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
			},
		},
		{
			&model.Alert{
				Annotations: map[string]string{"hello": "world"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}},
			},
		},
		{
			&model.Alert{
				Name:        "alert1",
				Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{ $value }}"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "alert1", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{ $value }}"}},
			},
		},
		{
			&model.Alert{
				Name:        "alert1",
				Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{}}"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "alert1", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "lorem={{ $labels.lorem }} value={{}}"}},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := getFires(tc.alert, queryData)
			assert.Equal(tt, tc.want, fires)
		})
	}
}

func TestGetFires_vector_zero_Result(t *testing.T) {
	queryData := model.QueryData{ResultType: commonModel.ValVector, Result: []commonModel.Sample{}}
	testCases := []struct {
		alert *model.Alert
		want  []model.Fire
	}{
		{
			&model.Alert{},
			[]model.Fire{},
		},
		{
			&model.Alert{
				Name:        "alert1",
				Annotations: map[string]string{"apple": "banana"},
				Labels:      map[string]string{"lemon": "orange"},
			},
			[]model.Fire{},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := getFires(tc.alert, queryData)
			assert.Equal(tt, tc.want, fires)
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
		alert *model.Alert
		want  []model.Fire
	}{
		{
			&model.Alert{},
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
				{State: "firing", Labels: map[string]string{"alertname": "placeholder name", "firer": "venti"}, Annotations: map[string]string{"summary": "placeholder summary"}},
			},
		},
		{
			&model.Alert{
				Annotations: map[string]string{"hello": "world"},
				Labels:      map[string]string{"lorem": "ipsum"},
			},
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}},
				{State: "firing", Labels: map[string]string{"alertname": "placeholder name", "firer": "venti", "lorem": "ipsum"}, Annotations: map[string]string{"hello": "world", "summary": "placeholder summary"}},
			},
		},
		{
			&model.Alert{
				Name:        "alert1",
				Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{ $value }}"},
				Labels:      map[string]string{"rule": "pod-v1"},
			},
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "alert1", "firer": "venti", "rule": "pod-v1"}, Annotations: map[string]string{"summary": "pod=pod1 value=1111"}},
				{State: "firing", Labels: map[string]string{"alertname": "alert1", "firer": "venti", "rule": "pod-v1"}, Annotations: map[string]string{"summary": "pod=pod2 value=2222"}},
			},
		},
		{
			&model.Alert{
				Name:        "alert1",
				Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}, // parse error
				Labels:      map[string]string{"rule": "pod-v1"},
			},
			[]model.Fire{
				{State: "firing", Labels: map[string]string{"alertname": "alert1", "firer": "venti", "rule": "pod-v1"}, Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}},
				{State: "firing", Labels: map[string]string{"alertname": "alert1", "firer": "venti", "rule": "pod-v1"}, Annotations: map[string]string{"summary": "pod={{ $labels.pod }} value={{}}"}},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := getFires(tc.alert, queryData)
			assert.Equal(tt, tc.want, fires)
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
			assert.Nil(tt, err)
			assert.Equal(tt, tc.want, ret)
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
			assert.NotNil(tt, err)
			assert.Error(tt, err, tc.wantError)
			assert.Equal(tt, tc.input, output)
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
			assert.NotNil(tt, err)
			assert.Error(tt, err, tc.wantError)
			assert.Equal(tt, tc.input, output)
		})
	}
}

func TestSendFires_ok(t *testing.T) {
	err := alerter1.sendFires([]model.Fire{
		{
			State:       "firing",
			Labels:      map[string]string{"hello": "world"},
			Annotations: map[string]string{"lorem": "ipsum"},
		},
	})
	assert.Nil(t, err)
}

func TestSendFires_error_on_Post(t *testing.T) {
	tempURL := alerter1.alertmanagerURL
	alerter1.alertmanagerURL = ""
	err := alerter1.sendFires([]model.Fire{
		{
			State:       "firing",
			Labels:      map[string]string{"hello": "world"},
			Annotations: map[string]string{"lorem": "ipsum"},
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, `error on Post: Post "/api/v1/alerts": unsupported protocol scheme ""`, err.Error())
	alerter1.alertmanagerURL = tempURL
}

func TestSendFires_not_ok(t *testing.T) {
	tempURL := alerter1.alertmanagerURL
	// prometheus insted alertmanager
	alerter1.alertmanagerURL = servers.GetServersByType(ms.TypePrometheus)[0].URL
	err := alerter1.sendFires([]model.Fire{
		{
			State:       "firing",
			Labels:      map[string]string{"hello": "world"},
			Annotations: map[string]string{"lorem": "ipsum"},
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, `statusCode is not ok(200)`, err.Error())
	alerter1.alertmanagerURL = tempURL

}
