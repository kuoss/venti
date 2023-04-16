package alerter

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/mock"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/stretchr/testify/assert"
)

var (
	prometheus *httptest.Server
	stores     *store.Stores
	alrtr      *alerter
)

func init() {
	_ = os.Chdir("../..")
	prometheus = mock.PrometheusServer()
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
