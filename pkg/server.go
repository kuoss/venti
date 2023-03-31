package pkg

import (
	"github.com/kuoss/venti/pkg/alert"
	"github.com/kuoss/venti/pkg/configuration"
	"github.com/kuoss/venti/pkg/handler"
	"github.com/kuoss/venti/pkg/service"
	"log"

	"github.com/gin-gonic/gin"
)

func Run(version string) {
	configuration.Load(version)
	service.InitDB()
	alert.StartAlertDaemon()
	StartServer()
}

// TODO: add to routerGroup routes
// routerGroup.Use(tokenRequired)

func StartServer() {
	log.Println("Starting Venti Server...")

	router := gin.New()

	// token not required

	user := router.Group("/user")
	{
		user.POST("/login", handler.login)
		user.POST("/logout", handler.logout)
	}

	// routerGroup routes
	api := router.Group("/api")

	api.Group("/config")
	{

	}
	api.Group("/datasource")
	{

	}
	api.Group("/lethe")
	{

	}
	api.Group("/prometheus")
	{

	}
	api.Group("/alert")
	{

	}

	if len(configuration.config.AlertRuleGroups) < 1 {
		log.Println("No AlertRuleGroups...")
	} else {
		handler.routesAPIAlerts(routerGroup)
	}

	router.Use(handleSPA())

	log.Println("venti started...💨💨💨")
	_ = router.Run() // :8080
}
