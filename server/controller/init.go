package controller

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

var H *server.Hertz

func InitController() {
	// init controller
	H = server.Default()

	initServiceManager()

	// model-manager api
	modelManager := H.Group("/model-manager")
	modelManager.GET("/list")
}
