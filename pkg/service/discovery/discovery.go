package discovery

import "github.com/kuoss/venti/pkg/model"

type Discoverer interface {
	Do(discovery model.Discovery) ([]model.Datasource, error)
}
