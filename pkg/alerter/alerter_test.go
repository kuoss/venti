package alerter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/mock"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"github.com/kuoss/venti/pkg/store/discovery"
	commonModel "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

var (
	prometheus   *httptest.Server
	alertmanager *httptest.Server
	stores       *store.Stores
	alrtr        *alerter
)

func init() {
	_ = os.Chdir("../..")
	prometheus = mock.Prometheus()
	alertmanager = mock.Alertmanager()
	datasources := mock.Datasources()
	datasources[0].URL = prometheus.URL
	datasources[1].URL = prometheus.URL
	datasources[2].URL = prometheus.URL
	datasourceConfig := mock.DatasourceConfigFromDatasources(datasources)

	// init stores
	var discoverer discovery.Discoverer
	datasourceStore, _ := store.NewDatasourceStore(datasourceConfig, discoverer)
	remoteStore := store.NewRemoteStore(&http.Client{}, 30*time.Second)
	alertRuleStore, _ := store.NewAlertRuleStore("etc/alertrules/*.yaml")
	stores = &store.Stores{
		AlertRuleStore:  alertRuleStore,
		DatasourceStore: datasourceStore,
		RemoteStore:     remoteStore,
	}
	alrtr = NewAlerter(stores)
	alrtr.SetAlertmanagerURL(alertmanager.URL)
}

func TestNewAlerter(t *testing.T) {
	alertRuleFiles := stores.AlertRuleStore.AlertRuleFiles()
	assert.Equal(t, mock.AlertRuleFiles(), alertRuleFiles)
	assert.Equal(t, 1, len(alertRuleFiles))
	assert.Equal(t, 1, len(alertRuleFiles[0].RuleGroups))
	assert.Equal(t, 3, len(alertRuleFiles[0].RuleGroups[0].Rules))
	assert.Equal(t, model.DatasourceSelector{Type: "prometheus"}, alertRuleFiles[0].DatasourceSelector)

	assert.Equal(t, false, alrtr.repeat)
	assert.Equal(t, 3, len(alrtr.alertFiles))
	assert.Equal(t, "mainPrometheus", alrtr.alertFiles[0].Datasource.Name)
	assert.Equal(t, "subPrometheus1", alrtr.alertFiles[1].Datasource.Name)
	assert.Equal(t, "subPrometheus2", alrtr.alertFiles[2].Datasource.Name)
}

func TestOnce(t *testing.T) {
	alrtr.Once()

}

func TestEvaluateAlert(t *testing.T) {
	testCases := []struct {
		alert     *model.Alert
		queryData model.QueryData
		want      []model.Fire
		err       error
	}{
		{
			&model.Alert{}, model.QueryData{},
			[]model.Fire(nil), nil,
		},
	}
	for i, c := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires, err := evaluateAlert(c.alert, c.queryData)
			assert.Equal(tt, c.want, fires)
			assert.Equal(tt, c.err, err)
		})
	}
}

func TestGetFires(t *testing.T) {
	testCases := []struct {
		alert     *model.Alert
		queryData model.QueryData
		want      []model.Fire
	}{
		{
			&model.Alert{},
			model.QueryData{},
			[]model.Fire{{State: "firing", Labels: map[string]string{"alertname": "", "firer": "venti"}, Annotations: map[string]string{"summary": "dummy summary from venti"}}},
		},
		{
			&model.Alert{
				Annotations: map[string]string{},
				Labels:      map[string]string{"lemon": "orange"},
			},
			model.QueryData{},
			[]model.Fire{{State: "firing", Labels: map[string]string{"alertname": "", "firer": "venti", "lemon": "orange"}, Annotations: map[string]string{"summary": "dummy summary from venti"}}},
		},
		{
			&model.Alert{
				Annotations: map[string]string{"apple": "banana"},
				Labels:      map[string]string{},
			},
			model.QueryData{},
			[]model.Fire{{State: "firing", Labels: map[string]string{"alertname": "", "firer": "venti"}, Annotations: map[string]string{"apple": "banana", "summary": "dummy summary from venti"}}},
		},
		{
			&model.Alert{
				Annotations: map[string]string{"apple": "banana"},
				Labels:      map[string]string{"lemon": "orange"},
			},
			model.QueryData{},
			[]model.Fire{{State: "firing", Labels: map[string]string{"alertname": "", "firer": "venti", "lemon": "orange"}, Annotations: map[string]string{"apple": "banana", "summary": "dummy summary from venti"}}},
		},
	}
	for i, c := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			fires := getFires(c.alert, c.queryData)
			assert.Equal(tt, c.want, fires)
		})
	}
}

func TestRenderSummary(t *testing.T) {
	testCases := []struct {
		tmplString string
		sample     *commonModel.Sample
		want       string
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
			"Monday",
			&commonModel.Sample{Value: 42},
			"Monday",
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
	for i, c := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			ret := renderSummary(c.tmplString, c.sample)
			assert.Equal(tt, c.want, ret)
		})
	}
}
