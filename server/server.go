package server

import (
	"github.com/kuoss/venti/server/alert"
	"github.com/kuoss/venti/server/api"
	"github.com/kuoss/venti/server/configuration"
	"log"

	"github.com/gin-gonic/gin"
)

// var secret = []byte("secret")

func Run(ventiVersion string) {
	configuration.Load(ventiVersion)
	InitDB()
	alert.StartAlertDaemon()
	StartServer()
}

func StartServer() {
	log.Println("Starting Venti Server...")
	r := gin.New()

	// token not required
	r.POST("/user/login", api.login)
	r.POST("/user/logout", api.logout)

	// routerGroup routes
	routerGroup := r.Group("/api")

	// TODO: add to routerGroup routes
	// routerGroup.Use(tokenRequired)
	api.routesAPIConfig(routerGroup)
	api.routesAPIDatasources(routerGroup)
	api.routesAPILethe(routerGroup)
	api.routesAPIPrometheus(routerGroup)

	if len(configuration.config.AlertRuleGroups) < 1 {
		log.Println("No AlertRuleGroups...")
	} else {
		api.routesAPIAlerts(routerGroup)
	}

	r.Use(handleSPA())

	log.Println("venti started...💨💨💨")
	_ = r.Run() // :8080
}
