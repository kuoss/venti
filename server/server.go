package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

// var secret = []byte("secret")

func Run(ventiVersion string) {
	LoadConfig(ventiVersion)
	InitDB()
	StartAlertDaemon()
	StartServer()
}

func StartServer() {
	log.Println("Starting Venti Server...")
	r := gin.New()

	// token not required
	r.POST("/user/login", login)
	r.POST("/user/logout", logout)

	// api routes
	api := r.Group("/api")

	// TODO: add to api routes
	// api.Use(tokenRequired)
	routesAPIConfig(api)
	routesAPIDatasources(api)
	routesAPILethe(api)
	routesAPIPrometheus(api)
	if len(config.AlertRuleGroups) < 1 {
		log.Println("No AlertRuleGroups...")
	} else {
		routesAPIAlerts(api)
	}

	r.Use(handleSPA())

	log.Println("venti started...💨💨💨")
	_ = r.Run() // :8080
}
