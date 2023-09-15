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
	"github.com/valyala/fastjson"
)

type IAlertingService interface {
	DoAlert() error
}

type AlertingService struct {
	alertingRules     []AlertingRule
	globalLabels      map[string]string
	datasourceService datasourceservice.IDatasourceService
	datasourceReload  bool
	remoteService     *remote.RemoteService
	alertmanagerURL   string
	client            http.Client
}

func New(cfg *model.Config, alertRuleFiles []model.RuleFile, datasourceService datasourceservice.IDatasourceService, remoteService *remote.RemoteService) *AlertingService {
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
		globalLabels:      cfg.AlertingConfig.GlobalLabels,
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
	// catch new alerts via query
	for _, datasource := range datasources {
		samples, err := s.queryRule(ar.rule, datasource)
		if err != nil {
			continue
		}

		commonlabels := map[string]string{
			"alertname":  ar.rule.Alert,
			"datasource": datasource.Name,
		}
		for k, v := range s.globalLabels {
			commonlabels[k] = v
		}
		for k, v := range ar.commonLabels {
			commonlabels[k] = v
		}
		for k, v := range ar.rule.Labels {
			commonlabels[k] = v
		}

		for _, sample := range samples {
			labels := map[string]string{}
			for k, v := range sample.Metric {
				labels[string(k)] = string(v)
			}
			for k, v := range commonlabels {
				labels[k] = v
			}
			signature := commonModel.LabelsToSignature(labels)

			createdAt := now
			state := StatePending
			annotations := map[string]string{}
			for k, v := range ar.rule.Annotations {
				annotations[k] = v
			}
			// render summary
			summary, err := renderSummary(annotations["summary"], sample)
			if err != nil {
				logger.Warnf("renderSummary err: %s", err)
			}

			temp, exists := ar.active[signature]
			if exists {
				createdAt = temp.CreatedAt
			}
			if !now.Before(createdAt.Add(time.Duration(ar.rule.For))) {
				state = StateFiring
			}
			annotations["summary"] = summary
			alert := &Alert{
				State:       state,
				CreatedAt:   createdAt,
				UpdatedAt:   now,
				Labels:      labels,
				Annotations: annotations,
			}
			ar.active[signature] = alert
			logger.Infof("[%s] labels:%v annotations:%v", alert.State.String(), alert.Labels, alert.Annotations)
		}
	}
	// remove old alerts
	for key, alert := range ar.active {
		if alert.UpdatedAt != now {
			delete(ar.active, key)
		}
	}
}

func renderSummary(input string, sample commonModel.Sample) (string, error) {
	// pre-render
	text := input
	text = strings.ReplaceAll(text, "$value", sample.Value.String())
	text = strings.ReplaceAll(text, "$labels.", ".")
	text = strings.ReplaceAll(text, "$labels", ".")

	// render
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return input, fmt.Errorf("parse err: %w", err)
	}
	labels := map[string]string{}
	for k, v := range sample.Metric {
		labels[string(k)] = string(v)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, labels)
	if err != nil {
		return input, fmt.Errorf("tmpl.execute err: %w", err)
	}
	return buf.String(), nil
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
	logger.Infof("sending %d fires...", len(fires))
	pbytes, err := json.Marshal(fires)
	if err != nil {
		// unreachable
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
		{Labels: map[string]string{"alertname": "pizza", "severity": "info", "pizza": "🍕", "time": time.Now().String()}},
	}
	err := s.sendFires(fires)
	if err != nil {
		return fmt.Errorf("sendFires err: %w", err)
	}
	return nil
}
