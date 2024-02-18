package controller

import (
	"github.com/MingkaiLee/kasos/hpa-executor/conf"
	"github.com/cloudwego/hertz/pkg/app/server"
)

var Server *server.Hertz

func init() {
	Server = server.Default(server.WithHostPorts(conf.ServerConf.Port))
}
