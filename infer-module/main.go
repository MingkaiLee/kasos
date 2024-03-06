package main

import (
	"github.com/MingkaiLee/kasos/infer-module/client"
	"github.com/MingkaiLee/kasos/infer-module/config"
	"github.com/MingkaiLee/kasos/infer-module/controller"
	"github.com/MingkaiLee/kasos/infer-module/internal"
)

func main() {
	config.InitConf()
	client.InitClient()
	internal.InitInternal()
	controller.InitController()

	controller.H.Spin()
}
