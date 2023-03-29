package server

import (
	"github.com/kuoss/venti/server/alert"
	"github.com/kuoss/venti/server/configuration"
	"github.com/kuoss/venti/server/handler"
	"github.com/kuoss/venti/server/service"
	"log"

	"github.com/gin-gonic/gin"
)

// var secret = []byte("secret")

func Run(version string) {
	configuration.Load(version)
	service.InitDB()
	alert.StartAlertDaemon()
	StartServer()
}

func StartServer() {
	log.Println("Starting Venti Server...")

	router := gin.New()

	// token not required
	user := router.Group("/user")
	{
		user.POST("/login", handler.login)
		user.POST("logout", handler.logout)
	}

	// routerGroup routes
	routerGroup := router.Group("/handler")

	// TODO: add to routerGroup routes
	// routerGroup.Use(tokenRequired)
	handler.routesAPIConfig(routerGroup)
	handler.routesAPIDatasources(routerGroup)
	handler.routesAPILethe(routerGroup)
	handler.routesAPIPrometheus(routerGroup)

	if len(configuration.config.AlertRuleGroups) < 1 {
		log.Println("No AlertRuleGroups...")
	} else {
		handler.routesAPIAlerts(routerGroup)
	}

	router.Use(handleSPA())

	log.Println("venti started...💨💨💨")
	_ = router.Run() // :8080
}
