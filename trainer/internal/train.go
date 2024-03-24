package internal

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/MingkaiLee/kasos/trainer/config"
	"github.com/MingkaiLee/kasos/trainer/util"
)

type Trainer struct {
	serviceName string    // 服务名
	modelName   string    // 模型名
	date        time.Time // 日期
}

func NewTrainer(serviceName, modelName string, date time.Time) *Trainer {
	return &Trainer{
		serviceName: serviceName,
		modelName:   modelName,
		date:        date,
	}
}

func (t *Trainer) Train(dataDirName string) (err error) {
	// 路径准备
	dataPath := fmt.Sprintf("%s/%s.csv", dataDirName, t.serviceName)
	scriptPath := fmt.Sprintf("%s/train/%s.py",
		config.ScriptDirectory, t.modelName)
	modelDir := fmt.Sprintf("%s/%s", config.ModelDirectory, t.serviceName)
	// 先尝试创建服务专用模型目录
	err = os.MkdirAll(modelDir, 0777)
	if err != nil {
		util.LogErrorf("train.Trainer.Train, create model dir error: %v", err)
		return
	}
	modelPath := fmt.Sprintf("%s/%s/%s",
		config.ModelDirectory, t.serviceName, t.modelName)
	// 检查数据是否已准备好
	dataFileStat, err := os.Stat(dataPath)
	if err != nil {
		util.LogErrorf("train.Trainer.Train, data file not ready, error: %v", err)
		return
	}
	if dataFileStat.IsDir() {
		util.LogErrorf("train.Trainer.Train, data is a directory, expected a file")
		return
	}
	// 检查训练脚本是否已准备好
	scriptFileStat, err := os.Stat(scriptPath)
	if err != nil {
		util.LogErrorf("train.Trainer.Train, script file not ready, error: %v", err)
		return
	}
	if scriptFileStat.IsDir() {
		util.LogErrorf("train.Trainer.Train, script is a directory, expected a file")
		return
	}
	// 检查模型是否已存在
	modelFileStat, err := os.Stat(modelPath)
	if os.IsNotExist(err) {
		// 模型不存在, 则首次训练
		util.LogInfof("train.Trainer.Train, train command: python3 %s --new -d %s -m %s", scriptPath, dataPath, modelPath)
		cmd := exec.Command("python3", scriptPath, "--new", "-d", dataPath, "-m", modelPath)
		err = cmd.Run()
		if err != nil {
			util.LogErrorf("train.Trainer.Train, train failed, error: %v", err)
			return
		}
	} else {
		// 模型文件若为目录, 则遇到了严重错误
		if modelFileStat.IsDir() {
			util.LogErrorf("train.Trainer.Train, model is a directory, expected a file")
			return
		}
		// 模型已存在, 迭代训练
		util.LogInfof("train.Trainer.Train, train command: python3 %s -d %s -m %s", scriptPath, dataPath, modelPath)
		cmd := exec.Command("python3", scriptPath, "-d", dataPath, "-m", modelPath)
		err = cmd.Run()
		if err != nil {
			util.LogErrorf("train.Trainer.Train, train failed, error: %v", err)
			return
		}
	}
	util.LogInfof("train.Trainer.Train, train success, service: %s, model: %s", t.serviceName, t.modelName)
	return
}
