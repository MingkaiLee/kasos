package main

import (
	"github.com/MingkaiLee/kasos/hpa-executor/client"
	"github.com/MingkaiLee/kasos/hpa-executor/config"
	"github.com/MingkaiLee/kasos/hpa-executor/controller"
)

func main() {
	config.InitConfig()
	client.InitClient()

	controller.H.Spin()
}
