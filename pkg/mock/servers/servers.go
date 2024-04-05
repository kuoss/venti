package servers

import (
	"fmt"

	"github.com/kuoss/venti/pkg/mocker"
	"github.com/kuoss/venti/pkg/mocker/alertmanager"
	"github.com/kuoss/venti/pkg/mocker/lethe"
	"github.com/kuoss/venti/pkg/mocker/prometheus"
	"github.com/kuoss/venti/pkg/model"
)

type Servers struct {
	Svrs []svr
}

type svr struct {
	Server    *mocker.Server
	Type      Type
	Name      string
	IsMain    bool
	BasicAuth bool
}

type Type int

const (
	TypeAlertmanager Type = iota
	TypeLethe
	TypePrometheus
)

type Requirements []Requirement

type Requirement struct {
	Type      Type
	Name      string
	IsMain    bool
	BasicAuth bool
}

func New(requirements Requirements) *Servers {
	s := &Servers{}
	for _, r := range requirements {
		var server *mocker.Server
		var err error
		switch r.Type {
		case TypeAlertmanager:
			server, err = alertmanager.New()
		case TypeLethe:
			server, err = lethe.New()
		case TypePrometheus:
			server, err = prometheus.New()
		}
		if err != nil {
			panic(err)
		}
		if r.BasicAuth {
			server.SetBasicAuth("abc", "123")
		}
		s.Svrs = append(s.Svrs, svr{
			Server:    server,
			Type:      r.Type,
			Name:      r.Name,
			IsMain:    r.IsMain,
			BasicAuth: r.BasicAuth,
		})
	}
	return s
}

func (s *Servers) GetDatasources() []model.Datasource {
	datasources := []model.Datasource{}
	for _, svr := range s.Svrs {
		var typ model.DatasourceType
		switch svr.Type {
		case TypeAlertmanager:
			continue
		case TypeLethe:
			typ = model.DatasourceTypeLethe
		case TypePrometheus:
			typ = model.DatasourceTypePrometheus
		}
		datasources = append(datasources, model.Datasource{
			Type:   typ,
			Name:   svr.Name,
			URL:    svr.Server.URL(),
			IsMain: svr.IsMain,
		})
	}
	return datasources
}

func (s *Servers) GetServersByType(typ Type) (servers []*mocker.Server) {
	for _, svr := range s.Svrs {
		if svr.Type == typ {
			servers = append(servers, svr.Server)
		}
	}
	return
}

func (s *Servers) Close() {
	for _, svr := range s.Svrs {
		svr.Server.Close()
	}
	fmt.Println("servers closed...")
}
