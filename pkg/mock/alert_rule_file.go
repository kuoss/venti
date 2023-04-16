package mock

import "github.com/kuoss/venti/pkg/model"

func AlertRuleFiles() []model.RuleFile {
	return []model.RuleFile{
		{
			Kind:               "AlertRuleFile",
			CommonLabels:       map[string]string{"severity": "silence"},
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

}
