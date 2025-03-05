package alertrule

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ruleFiles = []model.RuleFile{{
		Kind:               "AlertRuleFile",
		CommonLabels:       map[string]string{"rulefile": "sample-v3", "severity": "silence"},
		DatasourceSelector: model.DatasourceSelector{System: "", Type: "prometheus"},
		RuleGroups: []model.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{
				{Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
				{Alert: "PodNotHealthy", Expr: "sum by (namespace, pod) (kube_pod_status_phase{phase=~\"Pending|Unknown|Failed\"}) > 0", For: 3000000000, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "{{ $labels.namespace }}/{{ $labels.pod }}"}},
			}}}}}
	ruleFiles1 = []model.RuleFile{
		{Kind: "AlertRuleFile", CommonLabels: map[string]string{"rulefile": "sample-v3", "severity": "silence"}, DatasourceSelector: model.DatasourceSelector{System: "", Type: "prometheus"}, RuleGroups: []model.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{
				{Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
				{Alert: "PodNotHealthy", Expr: "sum by (namespace, pod) (kube_pod_status_phase{phase=~\"Pending|Unknown|Failed\"}) > 0", For: 3000000000, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "{{ $labels.namespace }}/{{ $labels.pod }}"}},
			}}}}}
)

func init() {
	err := os.Chdir("../../..")
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		name      string
		pattern   string
		want      *AlertRuleService
		wantError string
	}{
		// ok
		{
			"ok",
			"etc/alertrules/*.y*ml",
			&AlertRuleService{AlertRuleFiles: ruleFiles1},
			"",
		},
		{
			"ok",
			"",
			&AlertRuleService{AlertRuleFiles: ruleFiles1},
			"",
		},
		// error
		{
			"error",
			"asdf",
			&AlertRuleService{AlertRuleFiles: []model.RuleFile(nil)},
			"",
		},
		{
			"error",
			"[]",
			(*AlertRuleService)(nil),
			"glob err: syntax error in pattern",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, err := New(tc.pattern)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, service)
		})
	}
}

func TestNew_tempFiles(t *testing.T) {
	testCases := []struct {
		filename  string
		content   string
		want      *AlertRuleService
		wantError string
	}{
		{
			"test.ok.yaml",
			`
groups:
- name: info
  rules:
  - alert: hello
    expr: greet > 90
    for: 1m
    annotations:
      summary: "hello world"
`,
			&AlertRuleService{AlertRuleFiles: []model.RuleFile{{RuleGroups: []model.RuleGroup{
				{Name: "info", Rules: []model.Rule{{
					Alert:       "hello",
					Expr:        "greet > 90",
					For:         60000000000,
					Annotations: map[string]string{"summary": "hello world"}}}}}}}},
			"",
		},
		{
			"test.err.yaml",
			`
groups:
- name: info
  rules:
  - alert: hello
    expr: greet > 90
    for: 0m
    annotations:
      summary: "hello" world"
`,
			nil,
			"loadAlertRuleFileFromFilename err: unmarshalStrict err: yaml: line 8: did not find expected key",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			_ = os.WriteFile(tc.filename, []byte(tc.content), 0660)
			defer func() {
				os.RemoveAll(tc.filename)
			}()
			service, err := New(tc.filename)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, service)
		})
	}
}

func TestAlertRuleFiles(t *testing.T) {
	testCases := []struct {
		pattern string
		want    []model.RuleFile
	}{
		{
			"",
			ruleFiles,
		},
		{
			"asdf",
			[]model.RuleFile(nil),
		},
		{
			"etc/alertrules/*.yml",
			ruleFiles,
		},
		{
			"etc/alertrules/*.yaml",
			[]model.RuleFile(nil),
		},
		{
			"etc/alertrules/*.y*ml",
			ruleFiles,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			service, err := New(tc.pattern)
			assert.NoError(t, err)
			ruleFiles := service.GetAlertRuleFiles()
			assert.Equal(t, tc.want, ruleFiles)
		})
	}
}

func TestLoadAwesomeAlerts(t *testing.T) {
	// Set up test case files
	tests := []struct {
		file string
		want *model.RuleFile
	}{
		{
			"testdata/awesome-prometheus-alerts/kubestate-exporter.yml",
			&model.RuleFile{
				RuleGroups: []model.RuleGroup{{Name: "KubestateExporter", Rules: []model.Rule{
					{Alert: "KubernetesNodeNotReady", Expr: "kube_node_status_condition{condition=\"Ready\",status=\"true\"} == 0", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Node {{ $labels.node }} has been unready for a long time\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Node ready (node {{ $labels.node }})"}},
					{Alert: "KubernetesNodeMemoryPressure", Expr: "kube_node_status_condition{condition=\"MemoryPressure\",status=\"true\"} == 1", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Node {{ $labels.node }} has MemoryPressure condition\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes memory pressure (node {{ $labels.node }})"}},
					{Alert: "KubernetesNodeDiskPressure", Expr: "kube_node_status_condition{condition=\"DiskPressure\",status=\"true\"} == 1", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Node {{ $labels.node }} has DiskPressure condition\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes disk pressure (node {{ $labels.node }})"}},
					{Alert: "KubernetesNodeNetworkUnavailable", Expr: "kube_node_status_condition{condition=\"NetworkUnavailable\",status=\"true\"} == 1", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Node {{ $labels.node }} has NetworkUnavailable condition\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Node network unavailable (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesNodeOutOfPodCapacity", Expr: "sum by (node) ((kube_pod_status_phase{phase=\"Running\"} == 1) + on(uid, instance) group_left(node) (0 * kube_pod_info{pod_template_hash=\"\"})) / sum by (node) (kube_node_status_allocatable{resource=\"pods\"}) * 100 > 90", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Node {{ $labels.node }} is out of pod capacity\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Node out of pod capacity (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesContainerOomKiller", Expr: "(kube_pod_container_status_restarts_total - kube_pod_container_status_restarts_total offset 10m >= 1) and ignoring (reason) min_over_time(kube_pod_container_status_last_terminated_reason{reason=\"OOMKilled\"}[10m]) == 1", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Container {{ $labels.container }} in pod {{ $labels.namespace }}/{{ $labels.pod }} has been OOMKilled {{ $value }} times in the last 10 minutes.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes container oom killer ({{ $labels.namespace }}/{{ $labels.pod }}:{{ $labels.container }})"}},
					{Alert: "KubernetesJobFailed", Expr: "kube_job_status_failed > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Job {{ $labels.namespace }}/{{ $labels.job_name }} failed to complete\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Job failed ({{ $labels.namespace }}/{{ $labels.job_name }})"}},
					{Alert: "KubernetesJobNotStarting", Expr: "kube_job_status_active == 0 and kube_job_status_failed == 0 and kube_job_status_succeeded == 0 and (time() - kube_job_status_start_time) > 600", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Job {{ $labels.namespace }}/{{ $labels.job_name }} did not start for 10 minutes\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Job not starting ({{ $labels.namespace }}/{{ $labels.job_name }})"}},
					{Alert: "KubernetesCronjobSuspended", Expr: "kube_cronjob_spec_suspend != 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "CronJob {{ $labels.namespace }}/{{ $labels.cronjob }} is suspended\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes CronJob suspended ({{ $labels.namespace }}/{{ $labels.cronjob }})"}},
					{Alert: "KubernetesPersistentvolumeclaimPending", Expr: "kube_persistentvolumeclaim_status_phase{phase=\"Pending\"} == 1", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "PersistentVolumeClaim {{ $labels.namespace }}/{{ $labels.persistentvolumeclaim }} is pending\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes PersistentVolumeClaim pending ({{ $labels.namespace }}/{{ $labels.persistentvolumeclaim }})"}},
					{Alert: "KubernetesVolumeOutOfDiskSpace", Expr: "kubelet_volume_stats_available_bytes / kubelet_volume_stats_capacity_bytes * 100 < 10", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Volume is almost full (< 10% left)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Volume out of disk space (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesVolumeFullInFourDays", Expr: "predict_linear(kubelet_volume_stats_available_bytes[6h:5m], 4 * 24 * 3600) < 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Volume under {{ $labels.namespace }}/{{ $labels.persistentvolumeclaim }} is expected to fill up within four days. Currently {{ $value | humanize }}% is available.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Volume full in four days (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesPersistentvolumeError", Expr: "kube_persistentvolume_status_phase{phase=~\"Failed|Pending\", job=\"kube-state-metrics\"} > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Persistent volume {{ $labels.persistentvolume }} is in bad state\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes PersistentVolumeClaim pending ({{ $labels.namespace }}/{{ $labels.persistentvolumeclaim }})"}},
					{Alert: "KubernetesStatefulsetDown", Expr: "kube_statefulset_replicas != kube_statefulset_status_replicas_ready > 0", For: 60000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "StatefulSet {{ $labels.namespace }}/{{ $labels.statefulset }} went down\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes StatefulSet down ({{ $labels.namespace }}/{{ $labels.statefulset }})"}},
					{Alert: "KubernetesHpaScaleInability", Expr: "(kube_horizontalpodautoscaler_spec_max_replicas - kube_horizontalpodautoscaler_status_desired_replicas) * on (horizontalpodautoscaler,namespace) (kube_horizontalpodautoscaler_status_condition{condition=\"ScalingLimited\", status=\"true\"} == 1) == 0", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "HPA {{ $labels.namespace }}/{{ $labels.horizontalpodautoscaler }} is unable to scale\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes HPA scale inability (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesHpaMetricsUnavailability", Expr: "kube_horizontalpodautoscaler_status_condition{status=\"false\", condition=\"ScalingActive\"} == 1", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "HPA {{ $labels.namespace }}/{{ $labels.horizontalpodautoscaler }} is unable to collect metrics\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes HPA metrics unavailability (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesHpaScaleMaximum", Expr: "(kube_horizontalpodautoscaler_status_desired_replicas >= kube_horizontalpodautoscaler_spec_max_replicas) and (kube_horizontalpodautoscaler_spec_max_replicas > 1) and (kube_horizontalpodautoscaler_spec_min_replicas != kube_horizontalpodautoscaler_spec_max_replicas)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "info"}, Annotations: map[string]string{"description": "HPA {{ $labels.namespace }}/{{ $labels.horizontalpodautoscaler }} has hit maximum number of desired pods\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes HPA scale maximum (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesHpaUnderutilized", Expr: "max(quantile_over_time(0.5, kube_horizontalpodautoscaler_status_desired_replicas[1d]) == kube_horizontalpodautoscaler_spec_min_replicas) by (horizontalpodautoscaler) > 3", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "info"}, Annotations: map[string]string{"description": "HPA {{ $labels.namespace }}/{{ $labels.horizontalpodautoscaler }} is constantly at minimum replicas for 50% of the time. Potential cost saving here.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes HPA underutilized (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesPodNotHealthy", Expr: "sum by (namespace, pod) (kube_pod_status_phase{phase=~\"Pending|Unknown|Failed\"}) > 0", For: 900000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Pod {{ $labels.namespace }}/{{ $labels.pod }} has been in a non-running state for longer than 15 minutes.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Pod not healthy ({{ $labels.namespace }}/{{ $labels.pod }})"}},
					{Alert: "KubernetesPodCrashLooping", Expr: "increase(kube_pod_container_status_restarts_total[1m]) > 3", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Pod {{ $labels.namespace }}/{{ $labels.pod }} is crash looping\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes pod crash looping ({{ $labels.namespace }}/{{ $labels.pod }})"}},
					{Alert: "KubernetesReplicasetReplicasMismatch", Expr: "kube_replicaset_spec_replicas != kube_replicaset_status_ready_replicas", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "ReplicaSet {{ $labels.namespace }}/{{ $labels.replicaset }} replicas mismatch\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes ReplicasSet mismatch ({{ $labels.namespace }}/{{ $labels.replicaset }})"}},
					{Alert: "KubernetesDeploymentReplicasMismatch", Expr: "kube_deployment_spec_replicas != kube_deployment_status_replicas_available", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Deployment {{ $labels.namespace }}/{{ $labels.deployment }} replicas mismatch\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Deployment replicas mismatch ({{ $labels.namespace }}/{{ $labels.deployment }})"}},
					{Alert: "KubernetesStatefulsetReplicasMismatch", Expr: "kube_statefulset_status_replicas_ready != kube_statefulset_status_replicas", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "StatefulSet does not match the expected number of replicas.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes StatefulSet replicas mismatch (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesDeploymentGenerationMismatch", Expr: "kube_deployment_status_observed_generation != kube_deployment_metadata_generation", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Deployment {{ $labels.namespace }}/{{ $labels.deployment }} has failed but has not been rolled back.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes Deployment generation mismatch ({{ $labels.namespace }}/{{ $labels.deployment }})"}},
					{Alert: "KubernetesStatefulsetGenerationMismatch", Expr: "kube_statefulset_status_observed_generation != kube_statefulset_metadata_generation", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "StatefulSet {{ $labels.namespace }}/{{ $labels.statefulset }} has failed but has not been rolled back.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes StatefulSet generation mismatch ({{ $labels.namespace }}/{{ $labels.statefulset }})"}},
					{Alert: "KubernetesStatefulsetUpdateNotRolledOut", Expr: "max without (revision) (kube_statefulset_status_current_revision unless kube_statefulset_status_update_revision) * (kube_statefulset_replicas != kube_statefulset_status_replicas_updated)", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "StatefulSet {{ $labels.namespace }}/{{ $labels.statefulset }} update has not been rolled out.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes StatefulSet update not rolled out ({{ $labels.namespace }}/{{ $labels.statefulset }})"}},
					{Alert: "KubernetesDaemonsetRolloutStuck", Expr: "kube_daemonset_status_number_ready / kube_daemonset_status_desired_number_scheduled * 100 < 100 or kube_daemonset_status_desired_number_scheduled - kube_daemonset_status_current_number_scheduled > 0", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Some Pods of DaemonSet {{ $labels.namespace }}/{{ $labels.daemonset }} are not scheduled or not ready\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes DaemonSet rollout stuck ({{ $labels.namespace }}/{{ $labels.daemonset }})"}},
					{Alert: "KubernetesDaemonsetMisscheduled", Expr: "kube_daemonset_status_number_misscheduled > 0", For: 60000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Some Pods of DaemonSet {{ $labels.namespace }}/{{ $labels.daemonset }} are running where they are not supposed to run\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes DaemonSet misscheduled ({{ $labels.namespace }}/{{ $labels.daemonset }})"}},
					{Alert: "KubernetesCronjobTooLong", Expr: "time() - kube_cronjob_next_schedule_time > 3600", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "CronJob {{ $labels.namespace }}/{{ $labels.cronjob }} is taking more than 1h to complete.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes CronJob too long ({{ $labels.namespace }}/{{ $labels.cronjob }})"}},
					{Alert: "KubernetesJobSlowCompletion", Expr: "kube_job_spec_completions - kube_job_status_succeeded - kube_job_status_failed > 0", For: 43200000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Kubernetes Job {{ $labels.namespace }}/{{ $labels.job_name }} did not complete in time.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes job slow completion ({{ $labels.namespace }}/{{ $labels.job_name }})"}},
					{Alert: "KubernetesApiServerErrors", Expr: "sum(rate(apiserver_request_total{job=\"apiserver\",code=~\"(?:5..)\"}[1m])) by (instance, job) / sum(rate(apiserver_request_total{job=\"apiserver\"}[1m])) by (instance, job) * 100 > 3", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Kubernetes API server is experiencing high error rate\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes API server errors (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesApiClientErrors", Expr: "(sum(rate(rest_client_requests_total{code=~\"(4|5)..\"}[1m])) by (instance, job) / sum(rate(rest_client_requests_total[1m])) by (instance, job)) * 100 > 1", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Kubernetes API client is experiencing high error rate\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes API client errors (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesClientCertificateExpiresNextWeek", Expr: "apiserver_client_certificate_expiration_seconds_count{job=\"apiserver\"} > 0 and histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job=\"apiserver\"}[5m]))) < 7*24*60*60", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "A client certificate used to authenticate to the apiserver is expiring next week.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes client certificate expires next week (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesClientCertificateExpiresSoon", Expr: "apiserver_client_certificate_expiration_seconds_count{job=\"apiserver\"} > 0 and histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job=\"apiserver\"}[5m]))) < 24*60*60", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "A client certificate used to authenticate to the apiserver is expiring in less than 24.0 hours.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes client certificate expires soon (instance {{ $labels.instance }})"}},
					{Alert: "KubernetesApiServerLatency", Expr: "histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket{verb!~\"(?:CONNECT|WATCHLIST|WATCH|PROXY)\"} [10m])) WITHOUT (subresource)) > 1", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Kubernetes API server has a 99th percentile latency of {{ $value }} seconds for {{ $labels.verb }} {{ $labels.resource }}.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Kubernetes API server latency (instance {{ $labels.instance }})"}},
				}}},
			},
		},
		{
			"testdata/awesome-prometheus-alerts/node-exporter.yml",
			&model.RuleFile{
				RuleGroups: []model.RuleGroup{{Name: "NodeExporter", Interval: 0, Limit: 0, Rules: []model.Rule{
					{Alert: "HostOutOfMemory", Expr: "(node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < .10)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Node memory is filling up (< 10% left)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host out of memory (instance {{ $labels.instance }})"}},
					{Alert: "HostMemoryUnderMemoryPressure", Expr: "(rate(node_vmstat_pgmajfault[5m]) > 1000)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "The node is under heavy memory pressure. High rate of loading memory pages from disk.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host memory under memory pressure (instance {{ $labels.instance }})"}},
					{Alert: "HostMemoryIsUnderutilized", Expr: "min_over_time(node_memory_MemFree_bytes[1w]) > node_memory_MemTotal_bytes * .8", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "info"}, Annotations: map[string]string{"description": "Node memory usage is < 20% for 1 week. Consider reducing memory space. (instance {{ $labels.instance }})\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host Memory is underutilized (instance {{ $labels.instance }})"}},
					{Alert: "HostUnusualNetworkThroughputIn", Expr: "((rate(node_network_receive_bytes_total[5m]) / on(instance, device) node_network_speed_bytes) > .80)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Host receive bandwidth is high (>80%).\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host unusual network throughput in (instance {{ $labels.instance }})"}},
					{Alert: "HostUnusualNetworkThroughputOut", Expr: "((rate(node_network_transmit_bytes_total[5m]) / on(instance, device) node_network_speed_bytes) > .80)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Host transmit bandwidth is high (>80%)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host unusual network throughput out (instance {{ $labels.instance }})"}},
					{Alert: "HostUnusualDiskReadRate", Expr: "(rate(node_disk_io_time_seconds_total[5m]) > .80)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Disk is too busy (IO wait > 80%)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host unusual disk read rate (instance {{ $labels.instance }})"}},
					{Alert: "HostOutOfDiskSpace", Expr: "(node_filesystem_avail_bytes{fstype!~\"^(fuse.*|tmpfs|cifs|nfs)\"} / node_filesystem_size_bytes < .10 and on (instance, device, mountpoint) node_filesystem_readonly == 0)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Disk is almost full (< 10% left)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host out of disk space (instance {{ $labels.instance }})"}},
					{Alert: "HostDiskMayFillIn24Hours", Expr: "predict_linear(node_filesystem_avail_bytes{fstype!~\"^(fuse.*|tmpfs|cifs|nfs)\"}[1h], 86400) <= 0 and node_filesystem_avail_bytes > 0", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Filesystem will likely run out of space within the next 24 hours.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host disk may fill in 24 hours (instance {{ $labels.instance }})"}},
					{Alert: "HostOutOfInodes", Expr: "(node_filesystem_files_free / node_filesystem_files < .10 and ON (instance, device, mountpoint) node_filesystem_readonly == 0)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Disk is almost running out of available inodes (< 10% left)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host out of inodes (instance {{ $labels.instance }})"}},
					{Alert: "HostFilesystemDeviceError", Expr: "node_filesystem_device_error{fstype!~\"^(fuse.*|tmpfs|cifs|nfs)\"} == 1", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Error stat-ing the {{ $labels.mountpoint }} filesystem\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host filesystem device error (instance {{ $labels.instance }})"}},
					{Alert: "HostInodesMayFillIn24Hours", Expr: "predict_linear(node_filesystem_files_free{fstype!~\"^(fuse.*|tmpfs|cifs|nfs)\"}[1h], 86400) <= 0 and node_filesystem_files_free > 0", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Filesystem will likely run out of inodes within the next 24 hours at current write rate\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host inodes may fill in 24 hours (instance {{ $labels.instance }})"}},
					{Alert: "HostUnusualDiskReadLatency", Expr: "(rate(node_disk_read_time_seconds_total[1m]) / rate(node_disk_reads_completed_total[1m]) > 0.1 and rate(node_disk_reads_completed_total[1m]) > 0)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Disk latency is growing (read operations > 100ms)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host unusual disk read latency (instance {{ $labels.instance }})"}},
					{Alert: "HostUnusualDiskWriteLatency", Expr: "(rate(node_disk_write_time_seconds_total[1m]) / rate(node_disk_writes_completed_total[1m]) > 0.1 and rate(node_disk_writes_completed_total[1m]) > 0)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Disk latency is growing (write operations > 100ms)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host unusual disk write latency (instance {{ $labels.instance }})"}},
					{Alert: "HostHighCpuLoad", Expr: "(avg by (instance) (rate(node_cpu_seconds_total{mode!=\"idle\"}[2m]))) > .80", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "CPU load is > 80%\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host high CPU load (instance {{ $labels.instance }})"}},
					{Alert: "HostCpuIsUnderutilized", Expr: "(min by (instance) (rate(node_cpu_seconds_total{mode=\"idle\"}[1h]))) > 0.8", For: 604800000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "info"}, Annotations: map[string]string{"description": "CPU load has been < 20% for 1 week. Consider reducing the number of CPUs.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host CPU is underutilized (instance {{ $labels.instance }})"}},
					{Alert: "HostCpuStealNoisyNeighbor", Expr: "avg by(instance) (rate(node_cpu_seconds_total{mode=\"steal\"}[5m])) * 100 > 10", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "CPU steal is > 10%. A noisy neighbor is killing VM performances or a spot instance may be out of credit.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host CPU steal noisy neighbor (instance {{ $labels.instance }})"}},
					{Alert: "HostCpuHighIowait", Expr: "avg by (instance) (rate(node_cpu_seconds_total{mode=\"iowait\"}[5m])) > .10", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "CPU iowait > 10%. Your CPU is idling waiting for storage to respond.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host CPU high iowait (instance {{ $labels.instance }})"}},
					{Alert: "HostUnusualDiskIo", Expr: "rate(node_disk_io_time_seconds_total[5m]) > 0.8", For: 300000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Disk usage >80%. Check storage for issues or increase IOPS capabilities. Check storage for issues.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host unusual disk IO (instance {{ $labels.instance }})"}},
					{Alert: "HostContextSwitchingHigh", Expr: "(rate(node_context_switches_total[15m])/count without(mode,cpu) (node_cpu_seconds_total{mode=\"idle\"})) / (rate(node_context_switches_total[1d])/count without(mode,cpu) (node_cpu_seconds_total{mode=\"idle\"})) > 2", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Context switching is growing on the node (twice the daily average during the last 15m)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host context switching high (instance {{ $labels.instance }})"}},
					{Alert: "HostSwapIsFillingUp", Expr: "((1 - (node_memory_SwapFree_bytes / node_memory_SwapTotal_bytes)) * 100 > 80)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Swap is filling up (>80%)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host swap is filling up (instance {{ $labels.instance }})"}},
					{Alert: "HostSystemdServiceCrashed", Expr: "(node_systemd_unit_state{state=\"failed\"} == 1)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "systemd service crashed\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host systemd service crashed (instance {{ $labels.instance }})"}},
					{Alert: "HostPhysicalComponentTooHot", Expr: "node_hwmon_temp_celsius > node_hwmon_temp_max_celsius", For: 300000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Physical hardware component too hot\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host physical component too hot (instance {{ $labels.instance }})"}},
					{Alert: "HostNodeOvertemperatureAlarm", Expr: "((node_hwmon_temp_crit_alarm_celsius == 1) or (node_hwmon_temp_alarm == 1))", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Physical node temperature alarm triggered\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host node overtemperature alarm (instance {{ $labels.instance }})"}},
					{Alert: "HostSoftwareRaidInsufficientDrives", Expr: "((node_md_disks_required - on(device, instance) node_md_disks{state=\"active\"}) > 0)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "MD RAID array {{ $labels.device }} on {{ $labels.instance }} has insufficient drives remaining.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host software RAID insufficient drives (instance {{ $labels.instance }})"}},
					{Alert: "HostSoftwareRaidDiskFailure", Expr: "(node_md_disks{state=\"failed\"} > 0)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "MD RAID array {{ $labels.device }} on {{ $labels.instance }} needs attention.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host software RAID disk failure (instance {{ $labels.instance }})"}},
					{Alert: "HostKernelVersionDeviations", Expr: "changes(node_uname_info[1h]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "info"}, Annotations: map[string]string{"description": "Kernel version for {{ $labels.instance }} has changed.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host kernel version deviations (instance {{ $labels.instance }})"}},
					{Alert: "HostOomKillDetected", Expr: "(increase(node_vmstat_oom_kill[1m]) > 0)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "OOM kill detected\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host OOM kill detected (instance {{ $labels.instance }})"}},
					{Alert: "HostEdacCorrectableErrorsDetected", Expr: "(increase(node_edac_correctable_errors_total[1m]) > 0)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "info"}, Annotations: map[string]string{"description": "Host {{ $labels.instance }} has had {{ printf \"%.0f\" $value }} correctable memory errors reported by EDAC in the last 5 minutes.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host EDAC Correctable Errors detected (instance {{ $labels.instance }})"}},
					{Alert: "HostEdacUncorrectableErrorsDetected", Expr: "(node_edac_uncorrectable_errors_total > 0)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Host {{ $labels.instance }} has had {{ printf \"%.0f\" $value }} uncorrectable memory errors reported by EDAC in the last 5 minutes.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host EDAC Uncorrectable Errors detected (instance {{ $labels.instance }})"}},
					{Alert: "HostNetworkReceiveErrors", Expr: "(rate(node_network_receive_errs_total[2m]) / rate(node_network_receive_packets_total[2m]) > 0.01)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Host {{ $labels.instance }} interface {{ $labels.device }} has encountered {{ printf \"%.0f\" $value }} receive errors in the last two minutes.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host Network Receive Errors (instance {{ $labels.instance }})"}},
					{Alert: "HostNetworkTransmitErrors", Expr: "(rate(node_network_transmit_errs_total[2m]) / rate(node_network_transmit_packets_total[2m]) > 0.01)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Host {{ $labels.instance }} interface {{ $labels.device }} has encountered {{ printf \"%.0f\" $value }} transmit errors in the last two minutes.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host Network Transmit Errors (instance {{ $labels.instance }})"}},
					{Alert: "HostNetworkBondDegraded", Expr: "((node_bonding_active - node_bonding_slaves) != 0)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Bond \"{{ $labels.device }}\" degraded on \"{{ $labels.instance }}\".\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host Network Bond Degraded (instance {{ $labels.instance }})"}},
					{Alert: "HostConntrackLimit", Expr: "(node_nf_conntrack_entries / node_nf_conntrack_entries_limit > 0.8)", For: 300000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "The number of conntrack is approaching limit\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host conntrack limit (instance {{ $labels.instance }})"}},
					{Alert: "HostClockSkew", Expr: "((node_timex_offset_seconds > 0.05 and deriv(node_timex_offset_seconds[5m]) >= 0) or (node_timex_offset_seconds < -0.05 and deriv(node_timex_offset_seconds[5m]) <= 0))", For: 600000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Clock skew detected. Clock is out of sync. Ensure NTP is configured correctly on this host.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host clock skew (instance {{ $labels.instance }})"}},
					{Alert: "HostClockNotSynchronising", Expr: "(min_over_time(node_timex_sync_status[1m]) == 0 and node_timex_maxerror_seconds >= 16)", For: 120000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Clock not synchronising. Ensure NTP is configured on this host.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host clock not synchronising (instance {{ $labels.instance }})"}},
					{Alert: "HostRequiresReboot", Expr: "(node_reboot_required > 0)", For: 14400000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "info"}, Annotations: map[string]string{"description": "{{ $labels.instance }} requires a reboot.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Host requires reboot (instance {{ $labels.instance }})"}},
				}}},
			},
		},
		{
			"testdata/awesome-prometheus-alerts/embedded-exporter.yml",
			&model.RuleFile{
				RuleGroups: []model.RuleGroup{{Name: "EmbeddedExporter", Rules: []model.Rule{
					{Alert: "PrometheusJobMissing", Expr: "absent(up{job=\"prometheus\"})", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "A Prometheus job has disappeared\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus job missing (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTargetMissing", Expr: "up == 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "A Prometheus target has disappeared. An exporter might be crashed.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus target missing (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusAllTargetsMissing", Expr: "sum by (job) (up) == 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "A Prometheus job does not have living target anymore.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus all targets missing (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTargetMissingWithWarmupTime", Expr: "sum by (instance, job) ((up == 0) * on (instance) group_left(__name__) (node_time_seconds - node_boot_time_seconds > 600))", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Allow a job time to start up (10 minutes) before alerting that it's down.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus target missing with warmup time (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusConfigurationReloadFailure", Expr: "prometheus_config_last_reload_successful != 1", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Prometheus configuration reload error\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus configuration reload failure (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTooManyRestarts", Expr: "changes(process_start_time_seconds{job=~\"prometheus|pushgateway|alertmanager\"}[15m]) > 2", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Prometheus has restarted more than twice in the last 15 minutes. It might be crashlooping.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus too many restarts (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusAlertmanagerJobMissing", Expr: "absent(up{job=\"alertmanager\"})", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "A Prometheus AlertManager job has disappeared\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus AlertManager job missing (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusAlertmanagerConfigurationReloadFailure", Expr: "alertmanager_config_last_reload_successful != 1", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "AlertManager configuration reload error\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus AlertManager configuration reload failure (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusAlertmanagerConfigNotSynced", Expr: "count(count_values(\"config_hash\", alertmanager_config_hash)) > 1", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Configurations of AlertManager cluster instances are out of sync\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus AlertManager config not synced (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusAlertmanagerE2eDeadManSwitch", Expr: "vector(1)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus DeadManSwitch is an always-firing alert. It's used as an end-to-end test of Prometheus through the Alertmanager.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus AlertManager E2E dead man switch (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusNotConnectedToAlertmanager", Expr: "prometheus_notifications_alertmanagers_discovered < 1", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus cannot connect the alertmanager\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus not connected to alertmanager (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusRuleEvaluationFailures", Expr: "increase(prometheus_rule_evaluation_failures_total[3m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} rule evaluation failures, leading to potentially ignored alerts.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus rule evaluation failures (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTemplateTextExpansionFailures", Expr: "increase(prometheus_template_text_expansion_failures_total[3m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} template text expansion failures\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus template text expansion failures (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusRuleEvaluationSlow", Expr: "prometheus_rule_group_last_duration_seconds > prometheus_rule_group_interval_seconds", For: 300000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Prometheus rule evaluation took more time than the scheduled interval. It indicates a slower storage backend access or too complex query.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus rule evaluation slow (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusNotificationsBacklog", Expr: "min_over_time(prometheus_notifications_queue_length[10m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "The Prometheus notification queue has not been empty for 10 minutes\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus notifications backlog (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusAlertmanagerNotificationFailing", Expr: "rate(alertmanager_notifications_failed_total[1m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Alertmanager is failing sending notifications\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus AlertManager notification failing (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTargetEmpty", Expr: "prometheus_sd_discovered_targets == 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus has no target in service discovery\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus target empty (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTargetScrapingSlow", Expr: "prometheus_target_interval_length_seconds{quantile=\"0.9\"} / on (interval, instance, job) prometheus_target_interval_length_seconds{quantile=\"0.5\"} > 1.05", For: 300000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Prometheus is scraping exporters slowly since it exceeded the requested interval time. Your Prometheus server is under-provisioned.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus target scraping slow (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusLargeScrape", Expr: "increase(prometheus_target_scrapes_exceeded_sample_limit_total[10m]) > 10", For: 300000000000, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Prometheus has many scrapes that exceed the sample limit\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus large scrape (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTargetScrapeDuplicate", Expr: "increase(prometheus_target_scrapes_sample_duplicate_timestamp_total[5m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "Prometheus has many samples rejected due to duplicate timestamps but different values\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus target scrape duplicate (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTsdbCheckpointCreationFailures", Expr: "increase(prometheus_tsdb_checkpoint_creations_failed_total[1m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} checkpoint creation failures\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus TSDB checkpoint creation failures (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTsdbCheckpointDeletionFailures", Expr: "increase(prometheus_tsdb_checkpoint_deletions_failed_total[1m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} checkpoint deletion failures\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus TSDB checkpoint deletion failures (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTsdbCompactionsFailed", Expr: "increase(prometheus_tsdb_compactions_failed_total[1m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} TSDB compactions failures\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus TSDB compactions failed (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTsdbHeadTruncationsFailed", Expr: "increase(prometheus_tsdb_head_truncations_failed_total[1m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} TSDB head truncation failures\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus TSDB head truncations failed (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTsdbReloadFailures", Expr: "increase(prometheus_tsdb_reloads_failures_total[1m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} TSDB reload failures\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus TSDB reload failures (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTsdbWalCorruptions", Expr: "increase(prometheus_tsdb_wal_corruptions_total[1m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} TSDB WAL corruptions\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus TSDB WAL corruptions (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTsdbWalTruncationsFailed", Expr: "increase(prometheus_tsdb_wal_truncations_failed_total[1m]) > 0", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "critical"}, Annotations: map[string]string{"description": "Prometheus encountered {{ $value }} TSDB WAL truncation failures\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus TSDB WAL truncations failed (instance {{ $labels.instance }})"}},
					{Alert: "PrometheusTimeseriesCardinality", Expr: "label_replace(count by(__name__) ({__name__=~\".+\"}), \"name\", \"$1\", \"__name__\", \"(.+)\") > 10000", For: 0, KeepFiringFor: 0, Labels: map[string]string{"severity": "warning"}, Annotations: map[string]string{"description": "The \"{{ $labels.name }}\" timeseries cardinality is getting very high: {{ $value }}\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}", "summary": "Prometheus timeseries cardinality (instance {{ $labels.instance }})"}},
				}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			alertRuleFile, err := loadAlertRuleFileFromFilename(tt.file)
			require.NoError(t, err)
			require.Equal(t, tt.want, alertRuleFile)
		})
	}
}
