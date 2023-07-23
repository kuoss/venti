package alertmanager

import (
	"fmt"

	"github.com/kuoss/venti/pkg/mocker"
)

func New() (*mocker.Server, error) {
	s := mocker.New()
	s.GET("/api/v1/status/buildinfo", handleBuildInfo)
	s.GET("/api/v1/alerts", handleAlerts)

	err := s.Start()
	if err != nil {
		err = fmt.Errorf("error on Start: %w", err)
	}
	return s, err
}

func handleBuildInfo(c *mocker.Context) {
	c.JSONString(200, `{"status":"success","data":{"version":"2.41.0-alertmanager","revision":"c0d8a56c69014279464c0e15d8bfb0e153af0dab","branch":"HEAD","buildUser":"root@d20a03e77067","buildDate":"20221220-10:40:45","goVersion":"go1.19.4"}}`)
}

func handleAlerts(c *mocker.Context) {
	// 200 {"status":"success","data":{}}`
	c.JSON(200, mocker.H{
		"status": "success",
		"data":   mocker.H{},
	})
}
