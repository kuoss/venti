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
	"github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/remote"
	commonModel "github.com/prometheus/common/model"
	"github.com/valyala/fastjson"
)

type AlertingService struct {
	alertingRules     []AlertingRule
	datasourceService *datasource.DatasourceService
	datasourceReload  bool
	remoteService     *remote.RemoteService
	alertmanagerURL   string
	client            http.Client
}

func New(cfg *model.Config, alertRuleFiles []model.RuleFile, datasourceService *datasource.DatasourceService, remoteService *remote.RemoteService) *AlertingService {
	var alertmanagerURL string
	if len(cfg.AlertingConfig.AlertmanagerConfigs) > 0 && len(cfg.AlertingConfig.AlertmanagerConfigs[0].StaticConfig) > 0 && len(cfg.AlertingConfig.AlertmanagerConfigs[0].StaticConfig[0].Targets) > 0 {
		alertmanagerURL = cfg.AlertingConfig.AlertmanagerConfigs[0].StaticConfig[0].Targets[0]
	}
	var alertingRules = []AlertingRule{}
	for _, alertRuleFile := range alertRuleFiles {
		for _, group := range alertRuleFile.RuleGroups {
			for _, rule := range group.Rules {
				alertingRules = append(alertingRules, AlertingRule{
					datasourceSelector: alertRuleFile.DatasourceSelector,
					commonLabels:       alertRuleFile.CommonLabels,
					rule:               rule,
					active:             map[uint64]*Alert{},
				})
			}
		}
	}
	return &AlertingService{
		alertingRules:     alertingRules,
		datasourceService: datasourceService,
		datasourceReload:  cfg.DatasourceConfig.Discovery.Enabled,
		remoteService:     remoteService,
		alertmanagerURL:   alertmanagerURL,
		client:            http.Client{Timeout: 5 * time.Second},
	}
}

func (s *AlertingService) DoAlert() error {
	if s.alertingRules == nil || len(s.alertingRules) < 1 {
		return fmt.Errorf("no alertingRules")
	}
	if s.datasourceReload {
		err := s.datasourceService.Reload()
		if err != nil {
			return fmt.Errorf("reload err: %w", err)
		}
	}
	now := time.Now()
	fires := []Fire{}
	for _, ar := range s.alertingRules {
		s.updateAlertingRule(&ar, now)
		fires = append(fires, getFiresFromAlertingRule(&ar)...)
	}
	err := s.sendFires(fires)
	if err != nil {
		return fmt.Errorf("sendFires err: %w", err)
	}
	return nil
}

func getFiresFromAlertingRule(ar *AlertingRule) []Fire {
	fires := []Fire{}
	for _, alert := range ar.active {
		if alert.State == StateFiring {
			fire := Fire{
				Annotations: alert.Annotations,
				Labels:      alert.Labels,
			}
			fires = append(fires, fire)
		}
	}
	return fires
}

func (s *AlertingService) updateAlertingRule(ar *AlertingRule, now time.Time) {
	datasources := s.datasourceService.GetDatasourcesWithSelector(ar.datasourceSelector)
	maxState := StateInactive
	// catch new alerts via query
	for _, datasource := range datasources {
		samples, err := s.queryRule(ar.rule, datasource)
		if err != nil {
			continue
		}
		ruleLabels := ar.commonLabels
		ruleLabels["datasource"] = datasource.Name
		for k, v := range ar.rule.Labels {
			ruleLabels[k] = v
		}
		for _, sample := range samples {
			// signature
			tempLabels := map[string]string{"fingerprint": sample.Metric.Fingerprint().String()}
			for k, v := range ruleLabels {
				tempLabels[k] = v
			}
			tempLabels["alertname"] = ar.rule.Alert
			signature := commonModel.LabelsToSignature(tempLabels)

			_, exists := ar.active[signature]
			if !exists {
				ar.active[signature] = &Alert{
					State:       StatePending,
					CreatedAt:   now,
					UpdatedAt:   now,
					Labels:      ruleLabels,
					Annotations: ar.rule.Annotations,
				}
			} else {
				ar.active[signature].UpdatedAt = now
			}

			// render summary
			err := renderSummary(ar.active[signature], sample)
			if err != nil {
				logger.Warnf("renderSummary err: %s", err)
			}

			// update state
			maxState = StatePending
			// if ( created + for >= now ) firing
			if !ar.active[signature].CreatedAt.Add(time.Duration(ar.rule.For)).After(now) {
				ar.active[signature].State = StateFiring
				maxState = StateFiring
			}
		}
	}
	ar.state = maxState
	// remove old alerts
	for key, alert := range ar.active {
		if alert.UpdatedAt != now {
			delete(ar.active, key)
		}
	}
}

func renderSummary(alert *Alert, sample commonModel.Sample) error {
	input, exists := alert.Annotations["summary"]
	if !exists {
		input = "placeholder summary"
	}
	// pre-render
	text := input
	text = strings.ReplaceAll(text, "$value", sample.Value.String())
	text = strings.ReplaceAll(text, "$labels.", ".")
	text = strings.ReplaceAll(text, "$labels", ".")

	// render
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return fmt.Errorf("parse err: %w", err)
	}
	labels := map[string]string{}
	for k, v := range sample.Metric {
		labels[string(k)] = string(v)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, labels)
	if err != nil {
		return fmt.Errorf("tmpl.execute err: %w", err)
	}
	alert.Annotations["summary"] = buf.String()
	return nil
}

func (s *AlertingService) queryRule(rule model.Rule, ds model.Datasource) ([]commonModel.Sample, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	code, body, err := s.remoteService.GET(ctx, &ds, remote.ActionQuery, "query="+url.QueryEscape(rule.Expr))
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
	if err != nil {
		return []commonModel.Sample{}, fmt.Errorf("unmarshal err: %w", err)
	}
	return body.Data.Result, nil
}

func (s *AlertingService) sendFires(fires []Fire) error {
	pbytes, err := json.Marshal(fires)
	if err != nil {
		// test not reachable: memory full?
		return fmt.Errorf("error on Marshal: %w", err)
	}
	buff := bytes.NewBuffer(pbytes)
	resp, err := s.client.Post(s.alertmanagerURL+"/api/v1/alerts", "application/json", buff)
	if err != nil {
		return fmt.Errorf("error on Post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statusCode is not ok(200)")
	}
	return nil
}

func (s *AlertingService) SendTestAlert() error {
	fires := []Fire{
		{Labels: map[string]string{"test": "test", "severity": "info", "pizza": "🍕", "time": time.Now().String()}},
	}
	err := s.sendFires(fires)
	if err != nil {
		return fmt.Errorf("sendFires err: %w", err)
	}
	return nil
}
