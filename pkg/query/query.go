package query

type Query struct {
	//ctx context.Context

	Path   string
	Params map[string]string
	Start  string
	End    string
	Step   string
}

/*
type Query interface {
	execute() (string, error)
}
*/

/*
func (iq InstantQuery) execute() (string, error) {
	return query{
		DatasourceType: iq.DatasourceType,
		Path:           "/api/v1/query",
		Params: map[string]string{
			"query": iq.Expr,
			"time":  iq.Time,
		},
	}.execute()
}

// https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries
func (rq RangeQuery) execute() (string, error) {
	return query{
		DatasourceType: rq.DatasourceType,
		Path:           "/api/v1/query_range",
		Params: map[string]string{
			"query": rq.Expr,
			"start": rq.Start,
			"end":   rq.End,
			"step":  rq.Step,
		},
	}.execute()
}
*/
