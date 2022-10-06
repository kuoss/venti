package server

import (
	"log"
	"os"

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

	// web routes
	if os.Getenv("API_ONLY") != "1" {
		r.Static("/assets", "./web/dist/assets")
		r.StaticFile("/favicon.ico", "./web/public/favicon.ico")
		r.LoadHTMLGlob("web/dist/index.html")
		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})
	}

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

	log.Println("venti started...💨💨💨")
	r.Run() // 0.0.0.0:8080
}
