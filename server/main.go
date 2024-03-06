package main

import (
	"github.com/MingkaiLee/kasos/server/client"
	"github.com/MingkaiLee/kasos/server/config"
	"github.com/MingkaiLee/kasos/server/controller"
	"github.com/MingkaiLee/kasos/server/model"
)

func main() {
	// Init by specified order
	config.InitConf()
	model.InitModel()
	client.InitClient()
	controller.InitController()

	controller.H.Spin()
}
