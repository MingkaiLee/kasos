package config

import (
	"fmt"
	"os"

	"github.com/MingkaiLee/kasos/trainer/util"
	jsoniter "github.com/json-iterator/go"
)

var (
	DataDirectory   string
	ModelDirectory  string
	ScriptDirectory string
)

const volumeConfFile = "/etc/config/volume.json"

type VolumeConf struct {
	MountPath       string `json:"mount_path"`
	DataDirectory   string `json:"data_directory"`
	ModelDirectory  string `json:"model_directory"`
	ScriptDirectory string `json:"script_directory"`
}

func initVolumeConf() {
	var conf VolumeConf
	var err error

	d, err := os.ReadFile(volumeConfFile)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	err = jsoniter.Unmarshal(d, &conf)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	DataDirectory = fmt.Sprintf("%s/%s", conf.MountPath, conf.DataDirectory)
	ModelDirectory = fmt.Sprintf("%s/%s", conf.MountPath, conf.ModelDirectory)
	ScriptDirectory = fmt.Sprintf("%s/%s", conf.MountPath, conf.ScriptDirectory)
}
