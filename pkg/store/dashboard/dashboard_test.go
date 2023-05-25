package dashboard

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	store1          *DashboardStore
	sampleDashboard *model.Dashboard
)

func init() {
	_ = os.Chdir("../../..")
	sampleDashboard = &model.Dashboard{Title: "Sample", Rows: []model.Row{
		{Panels: []model.Panel{
			{Title: "time", Type: "stat", Headers: []string(nil), Targets: []model.Target{{Expr: "time()", Legend: "", Legends: []string(nil), Unit: "dateTimeAsLocal", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "apiserver%", Type: "stat", Headers: []string(nil), Targets: []model.Target{{Expr: "100 * up{job=\"kubernetes-apiservers\"}", Legend: "", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold{{Values: []int{80, 100}, Invert: true}}}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "k8s version", Type: "stat", Headers: []string(nil), Targets: []model.Target{{Expr: "kubernetes_build_info{job=\"kubernetes-apiservers\"}", Legend: "{{git_version}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)}}},
		{Panels: []model.Panel{
			{Title: "node", Type: "piechart", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(kube_node_status_condition{status='true'}) by (condition) > 0", Legend: "{{condition}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "namespace", Type: "piechart", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(kube_namespace_status_phase) by (phase) > 0", Legend: "{{phase}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "pod", Type: "piechart", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(kube_pod_status_phase{namespace=~\"$namespace\",node=~\"$node\"}) by (phase) > 0", Legend: "{{phase}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "job", Type: "piechart", Headers: []string(nil), Targets: []model.Target{
				{Expr: "sum(kube_job_status_active{namespace=~\"$namespace\"}) > 0", Legend: "Active", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)},
				{Expr: "sum(kube_job_status_failed{namespace=~\"$namespace\"}) > 0", Legend: "Failed", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)},
				{Expr: "sum(kube_job_status_succeeded{namespace=~\"$namespace\"}) > 0", Legend: "Succeeded", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "pvc", Type: "piechart", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(kube_persistentvolumeclaim_status_phase{namespace=~\"$namespace\"}) by (phase) > 0", Legend: "{{phase}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "pv", Type: "piechart", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(kube_persistentvolume_status_phase) by (phase) > 0", Legend: "{{phase}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)}}},
		{Panels: []model.Panel{
			{Title: "node load", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "node_load1", Legend: "{{node}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "node cpu%", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "100 * sum(rate(node_cpu_seconds_total{mode!=\"idle\",mode!=\"iowait\"}[3m])) by (node) / sum(kube_node_status_allocatable{resource=\"cpu\"}) by (node)", Legend: "{{node}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: &model.ChartOptions{YMax: 100}},
			{Title: "node mem%", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "100 * (1 - ( node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes ))", Legend: "{{node}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: &model.ChartOptions{YMax: 100}},
			{Title: "node pods", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(kubelet_running_pods) by (instance)", Legend: "{{instance}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: &model.ChartOptions{YMax: 120}}}},
		{Panels: []model.Panel{
			{Title: "node receive (Ki/m)", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(rate(node_network_receive_bytes_total[3m])) by (node) / 1024", Legend: "{{node}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "node transmit (Ki/m)", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(rate(node_network_transmit_bytes_total[3m])) by (node) / 1024", Legend: "{{node}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "node disk read (Ki/m)", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(rate(node_disk_read_bytes_total[3m])) by (node) / 1024", Legend: "{{node}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "node disk write (Ki/m)", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(rate(node_disk_written_bytes_total[3m])) by (node) / 1024", Legend: "{{node}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "node root fs%", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "100 * sum( 1-(node_filesystem_avail_bytes{mountpoint=\"/\"} / node_filesystem_size_bytes{mountpoint=\"/\"}) ) by (node)", Legend: "{{node}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: &model.ChartOptions{YMax: 100}}}},
		{Panels: []model.Panel{
			{Title: "pvc%", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "100 * max( 1 - kubelet_volume_stats_available_bytes / kubelet_volume_stats_capacity_bytes) by (namespace, persistentvolumeclaim)", Legend: "{{namespace}}/{{persistentvolumeclaim}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: &model.ChartOptions{YMax: 100}},
			{Title: "pvc inodes%", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "kubelet_volume_stats_inodes_used / kubelet_volume_stats_inodes * 100", Legend: "{{namespace}}/{{persistentvolumeclaim}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: &model.ChartOptions{YMax: 100}},
			{Title: "pod cpu(mcores)", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(rate(container_cpu_usage_seconds_total{namespace=~\"$namespace\", instance=~\"$node\", container!=\"\"}[5m])) by (namespace, pod) * 1000", Legend: "{{namespace}}/{{pod}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)},
			{Title: "pod mem(Mi)", Type: "time_series", Headers: []string(nil), Targets: []model.Target{{Expr: "sum(container_memory_working_set_bytes{namespace=~\"$namespace\", instance=~\"$node\", container!=\"\"}) by (namespace, pod) / 1024 / 1024", Legend: "{{namespace}}/{{pod}}", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)}}},
		{Panels: []model.Panel{
			{Title: "event", Type: "logs", Headers: []string(nil), Targets: []model.Target{{Expr: "pod{namespace=\"kube-system\",pod=\"eventrouter-.*\"}", Legend: "", Legends: []string(nil), Unit: "", Columns: []string(nil), Headers: []string(nil), Key: "", Thresholds: []model.Threshold(nil)}}, ChartOptions: (*model.ChartOptions)(nil)}}}}}

	var err error
	store1, err = New("etc/dashboards")
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		pattern   string
		want      *DashboardStore
		wantError string
	}{
		{
			"",
			&DashboardStore{dashboards: []model.Dashboard{*sampleDashboard}},
			"",
		},
		{
			"etc/hello",
			nil,
			"getDashboardFilesFromPath err: no dashboard file: dirpath: etc/hello",
		},
		{
			"etc/dashboards",
			&DashboardStore{dashboards: []model.Dashboard{*sampleDashboard}},
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			store, err := New(tc.pattern)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, store)
		})
	}
}

func TestLoadDashboardFromFile(t *testing.T) {
	testCases := []struct {
		pattern   string
		want      *model.Dashboard
		wantError string
	}{
		{
			"",
			nil,
			"error on ReadFile: open : no such file or directory",
		},
		{
			"etc/dashboards/sample.yml",
			sampleDashboard,
			"",
		},
	}
	for _, tc := range testCases {
		got, err := loadDashboardFromFile(tc.pattern)
		if tc.wantError == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tc.wantError)
		}
		assert.Equal(t, tc.want, got)
	}
}

func TestDashboards(t *testing.T) {
	want := []model.Dashboard{*sampleDashboard}
	got := store1.Dashboards()
	assert.Equal(t, want, got)
}
