package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	datasourceservice "github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/remote"
	commonModel "github.com/prometheus/common/model"
	promWebAPI "github.com/prometheus/prometheus/web/api/v1"
	"github.com/valyala/fastjson"
)

type IAlertingService interface {
	DoAlert() error
}

type AlertingService struct {
	alertingRuleGroups  []AlertingRuleGroup
	globalLabels        map[string]string
	datasourceService   datasourceservice.IDatasourceService
	datasourceReload    bool
	remoteService       *remote.RemoteService
	alertmanagerConfigs model.AlertmanagerConfigs
	alertmanagerURL     string
	client              http.Client
}

var (
	fakeErr1 bool = false
	fakeErr2 bool = false
)

func New(cfg *model.Config, alertRuleFiles []model.RuleFile, datasourceService datasourceservice.IDatasourceService, remoteService *remote.RemoteService) *AlertingService {
	var alertmanagerConfigs model.AlertmanagerConfigs
	var alertmanagerURL string
	if len(cfg.AlertingConfig.AlertmanagerConfigs) > 0 {
		alertmanagerConfigs = cfg.AlertingConfig.AlertmanagerConfigs
		if len(cfg.AlertingConfig.AlertmanagerConfigs[0].StaticConfig) > 0 && len(cfg.AlertingConfig.AlertmanagerConfigs[0].StaticConfig[0].Targets) > 0 {
			alertmanagerURL = cfg.AlertingConfig.AlertmanagerConfigs[0].StaticConfig[0].Targets[0]
		}
	}
	var alertingRuleGroups = []AlertingRuleGroup{}
	for _, alertRuleFile := range alertRuleFiles {
		var alertingRules = []AlertingRule{}
		for _, group := range alertRuleFile.RuleGroups {
			for _, rule := range group.Rules {
				alertingRules = append(alertingRules, AlertingRule{
					Rule:   rule,
					Active: map[uint64]*Alert{},
				})
			}
		}
		alertingRuleGroups = append(alertingRuleGroups, AlertingRuleGroup{
			DatasourceSelector: alertRuleFile.DatasourceSelector,
			GroupLabels:        alertRuleFile.CommonLabels,
			AlertingRules:      alertingRules,
		})
	}
	return &AlertingService{
		alertingRuleGroups:  alertingRuleGroups,
		globalLabels:        cfg.AlertingConfig.GlobalLabels,
		datasourceService:   datasourceService,
		datasourceReload:    cfg.DatasourceConfig.Discovery.Enabled,
		remoteService:       remoteService,
		alertmanagerConfigs: alertmanagerConfigs,
		alertmanagerURL:     alertmanagerURL,
		client:              http.Client{Timeout: 5 * time.Second},
	}
}

func (s *AlertingService) GetAlertingRuleGroups() []AlertingRuleGroup {
	return s.alertingRuleGroups
}

func (s *AlertingService) GetAlertmanagerDiscovery() promWebAPI.AlertmanagerDiscovery {
	var alertmanagers []*promWebAPI.AlertmanagerTarget
	for _, alertmanagerConfig := range s.alertmanagerConfigs {
		for _, staticConfig := range alertmanagerConfig.StaticConfig {
			for _, target := range staticConfig.Targets {
				alertmanagers = append(alertmanagers, &promWebAPI.AlertmanagerTarget{URL: target})
			}
		}
	}
	return promWebAPI.AlertmanagerDiscovery{
		ActiveAlertmanagers:  alertmanagers,
		DroppedAlertmanagers: []*promWebAPI.AlertmanagerTarget{},
	}
}

func (s *AlertingService) DoAlert() error {
	if s.datasourceReload {
		err := s.datasourceService.Reload()
		logger.Debugf("datasourceService.Reload") // 2023-09-19
		if err != nil {
			return fmt.Errorf("reload err: %w", err)
		}
	}
	fires := []Fire{}
	s.evalAlertingRuleGroups(&fires)
	err := s.sendFires(fires)
	if err != nil {
		return fmt.Errorf("sendFires err: %w", err)
	}
	return nil
}

func (s *AlertingService) evalAlertingRuleGroups(fires *[]Fire) {
	evalTime := time.Now()
	for _, group := range s.alertingRuleGroups {
		s.evalAlertingRuleGroup(&group, evalTime, fires)
	}
}

func (s *AlertingService) evalAlertingRuleGroup(group *AlertingRuleGroup, evalTime time.Time, fires *[]Fire) {
	datasources := s.datasourceService.GetDatasourcesWithSelector(group.DatasourceSelector)
	logger.Debugf("datasources(%d): %v", len(datasources), datasources) // 2023-09-19
	labels := map[string]string{}
	for k, v := range s.globalLabels {
		labels[k] = v
	}
	for k, v := range group.GroupLabels {
		labels[k] = v
	}
	for _, ar := range group.AlertingRules {
		s.evalAlertingRule(&ar, datasources, labels, evalTime, fires)
	}
}

func (s *AlertingService) evalAlertingRule(ar *AlertingRule, datasources []model.Datasource, commonLabels map[string]string, evalTime time.Time, fires *[]Fire) {
	labels := map[string]string{}
	for k, v := range commonLabels {
		labels[k] = v
	}
	for k, v := range ar.Rule.Labels {
		labels[k] = v
	}
	labels["alertname"] = ar.Rule.Alert

	for _, datasource := range datasources {
		err := s.evalAlertingRuleDatasource(ar, datasource, labels, evalTime)
		if err != nil {
			logger.Warnf("evalAlertingRuleDatasource err: %s", err)
		}
	}
	for key, alert := range ar.Active {
		// remove old alerts
		if alert.UpdatedAt != evalTime {
			delete(ar.Active, key)
			continue
		}
		// add to fires
		if alert.State == StateFiring {
			*fires = append(*fires, Fire{
				Annotations: alert.Annotations,
				Labels:      alert.Labels,
			})
		}
	}
}

func (s *AlertingService) evalAlertingRuleDatasource(ar *AlertingRule, datasource model.Datasource, commonLabels map[string]string, evalTime time.Time) error {
	labels := map[string]string{}
	for k, v := range commonLabels {
		labels[k] = v
	}
	for k, v := range ar.Rule.Labels {
		labels[k] = v
	}
	labels["datasource"] = datasource.Name

	samples, err := s.queryRule(ar.Rule, datasource)
	if err != nil {
		return fmt.Errorf("queryRule err: %w", err)
	}
	for _, sample := range samples {
		s.evalAlertingRuleSample(ar, sample, labels, evalTime)
	}
	return nil
}

func (s *AlertingService) evalAlertingRuleSample(ar *AlertingRule, sample commonModel.Sample, commonLabels map[string]string, evalTime time.Time) {
	// labels & signature
	labels := map[string]string{}
	for k, v := range commonLabels {
		labels[k] = v
	}
	for k, v := range sample.Metric {
		labels[string(k)] = string(v)
	}
	signature := commonModel.LabelsToSignature(labels)

	// annotations & summary
	annotations := map[string]string{}
	for k, v := range ar.Rule.Annotations {
		annotations[k] = v
	}
	err := renderSummaryAnnotaion(annotations, labels, sample.Value.String())
	if err != nil {
		logger.Warnf("renderSummaryAnnotaion(%s) err: %s", ar.Rule.Alert, err)
	}

	// others
	createdAt := evalTime
	state := StatePending

	temp, exists := ar.Active[signature]
	if exists {
		createdAt = temp.CreatedAt
	}
	elapsed := evalTime.Sub(createdAt.Add(ar.Rule.For))
	if elapsed >= 0 {
		state = StateFiring
	}
	alert := &Alert{
		State:       state,
		CreatedAt:   createdAt,
		UpdatedAt:   evalTime,
		Labels:      labels,
		Annotations: annotations,
	}
	ar.Active[signature] = alert

	// show log if severity exists and not silence
	severity, ok := alert.Annotations["severity"]
	if ok && severity != "silence" {
		logger.Infof("%s(%s): %s: %s", alert.State.String(), elapsed.Round(time.Second), labels["alertname"], annotations["summary"])
	}
}

func renderSummaryAnnotaion(annotations map[string]string, labels map[string]string, value string) error {
	summary, exists := annotations["summary"]
	if !exists {
		annotations["summary"] = "placeholder summary"
		return fmt.Errorf("no summary annotation")
	}

	// pre-render
	summary = strings.ReplaceAll(summary, "$value", value)
	summary = strings.ReplaceAll(summary, "$labels.", ".")
	summary = strings.ReplaceAll(summary, "$labels", ".")

	// render
	tmpl, err := template.New("").Parse(summary)
	if err != nil {
		return fmt.Errorf("parse err: %w", err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, labels)
	if err != nil || fakeErr1 {
		return fmt.Errorf("tmpl.Execute err: %w", err)
	}
	annotations["summary"] = buf.String()
	return nil
}

func (s *AlertingService) queryRule(rule model.Rule, datasource model.Datasource) ([]commonModel.Sample, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	code, body, err := s.remoteService.GET(ctx, &datasource, remote.ActionQuery, "query="+url.QueryEscape(rule.Expr))
	if err != nil {
		return []commonModel.Sample{}, fmt.Errorf("GET err: %w", err)
	}
	if code != http.StatusOK {
		return []commonModel.Sample{}, fmt.Errorf("not successful code=%d", code)
	}
	bodyBytes := []byte(body)
	status := fastjson.GetString(bodyBytes, "status")
	if status != "success" {
		return []commonModel.Sample{}, fmt.Errorf("not successful status=%s", status)
	}
	resultType := fastjson.GetString(bodyBytes, "data", "resultType")
	if resultType == "logs" {
		samples, err := getDataFromLogs(bodyBytes)
		if err != nil {
			return []commonModel.Sample{}, fmt.Errorf("getDataFromLogs err: %w", err)
		}
		return samples, nil
	}
	samples, err := getDataFromVector(bodyBytes)
	if err != nil {
		return []commonModel.Sample{}, fmt.Errorf("getDataFromVector err: %w", err)
	}
	return samples, nil
}

func getDataFromLogs(bodyBytes []byte) ([]commonModel.Sample, error) {
	type Data struct {
		ResultType string              `json:"resultType"`
		Result     []map[string]string `json:"result"`
	}
	type Body struct {
		Status string `json:"status"`
		Data   Data   `json:"data"`
	}
	var body Body
	err := json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return []commonModel.Sample{}, fmt.Errorf("unmarshal err: %w", err)
	}
	return []commonModel.Sample{{Value: commonModel.SampleValue(len(body.Data.Result))}}, nil
}

func getDataFromVector(bodyBytes []byte) ([]commonModel.Sample, error) {
	type Data struct {
		ResultType string               `json:"resultType"`
		Result     []commonModel.Sample `json:"result"`
	}
	type Body struct {
		Status string `json:"status"`
		Data   Data   `json:"data"`
	}
	var body Body
	err := json.Unmarshal(bodyBytes, &body)
	if err != nil || fakeErr1 {
		return []commonModel.Sample{}, fmt.Errorf("unmarshal err: %w", err)
	}
	return body.Data.Result, nil
}

func (s *AlertingService) sendFires(fires []Fire) error {
	logger.Infof("sending %d fires...", len(fires))
	pbytes, err := json.Marshal(fires)
	if err != nil || fakeErr1 {
		return fmt.Errorf("marshal err: %w", err)
	}
	buff := bytes.NewBuffer(pbytes)
	resp, err := s.client.Post(s.alertmanagerURL+"/api/v2/alerts", "application/json", buff)
	if err != nil {
		return fmt.Errorf("post err: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || fakeErr2 {
		return fmt.Errorf("statusCode is not ok(200)")
	}
	return nil
}

func (s *AlertingService) SendTestAlert() error {
	fires := []Fire{
		{Labels: map[string]string{"alertname": "pizza", "severity": "info", "pizza": "🍕", "time": time.Now().String()}},
	}
	err := s.sendFires(fires)
	if err != nil {
		return fmt.Errorf("sendFires err: %w", err)
	}
	return nil
}
