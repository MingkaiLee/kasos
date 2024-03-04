package controller

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

var Server *server.Hertz

func init() {
	initStressTester()
	initHpaExec()
}
