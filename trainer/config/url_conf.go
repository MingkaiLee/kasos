package config

import (
	"os"

	"github.com/MingkaiLee/kasos/server/util"
	jsoniter "github.com/json-iterator/go"
)

var (
	ServerUrl string
)

const urlConfFile = "/etc/config/url.json"

type UrlConf struct {
	HpaExecutorUrl string `json:"hpa_executor_url"`
	TrainerUrl     string `json:"trainer_url"`
	InferModuleUrl string `json:"infer_module_url"`
	ServerUrl      string `json:"server_url"`
}

func initUrlConf() {
	var conf UrlConf
	var err error

	d, err := os.ReadFile(urlConfFile)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	err = jsoniter.Unmarshal(d, &conf)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	ServerUrl = conf.ServerUrl
}
