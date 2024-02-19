package controller

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

var H *server.Hertz

func InitController() {
	// init controller
	H = server.Default()

	// service-manager api
	serviceManager := H.Group("/service-manager")
	serviceManager.GET("/list")
	serviceManager.GET("/find", FindServices)

	// model-manager api
	modelManager := H.Group("/model-manager")
	modelManager.GET("/list")
}
