package alertmanager

import (
	"fmt"

	"github.com/kuoss/venti/pkg/mocker"
)

func New(port int) (*mocker.Server, error) {
	s := mocker.New()
	s.GET("/api/v1/alerts", func(c *mocker.Context) {
		// 200 {"status":"success","data":{}}`
		c.JSON(200, mocker.H{
			"status": "success",
			"data":   mocker.H{},
		})
	})
	err := s.Start(port)
	if err != nil {
		err = fmt.Errorf("error on Start: %w", err)
	}
	return s, err
}
