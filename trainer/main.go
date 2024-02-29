package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/MingkaiLee/kasos/trainer/config"
	"github.com/MingkaiLee/kasos/trainer/internal"
	"github.com/MingkaiLee/kasos/trainer/util"
)

func main() {
	util.LogInfof("start training process, current time: %s", time.Now().Format(time.DateTime))
	// 获取所需要训练的服务
	services, err := internal.ListServices()
	if err != nil {
		util.LogErrorf("process stopped, error: %v", err)
		panic(err)
	}
	// 获取当前时间
	currentTime := time.Now()
	// 创建拉取数据的目标文件夹
	year, month, day := currentTime.Date()
	fileDirName := fmt.Sprintf("{%s}/{%d}-{%d}-{%d}", config.DataDirectory, year, month, day)
	err = os.MkdirAll(fileDirName, 0755)
	if err != nil {
		util.LogErrorf("process stopped, error: %v", err)
		panic(err)
	}
	// 拉取数据
	// 对于拉取数据有误的模型, 取消本轮训练
	disableTrains := make([]bool, len(services))
	for idx := range services {
		worker := internal.NewFetchDataWorker(runtime.NumCPU(), *services[idx].Name, currentTime)
		// 内部并行拉取数据
		err = worker.Run()
		if err != nil {
			util.LogErrorf("fetch data error, service: %s", *services[idx].Name)
			// 拉取数据有误取消训练标志置位
			disableTrains[idx] = true
		}
	}
	// 开始训练
	for idx := range disableTrains {
		// 训练迭代数据拉取成功的模型
		// 由于训练时CPU密集型操作, 不再开启goroutine处理
		if !disableTrains[idx] {
			trainer := internal.NewTrainer(*services[idx].Name, *services[idx].ModelName, currentTime)
			err = trainer.Train()
			if err != nil {
				util.LogErrorf("train error, service: %s, error: %v", *services[idx].Name, err)
			}
		}
	}
	util.LogInfof("training process finished, current time: %s", time.Now().Format(time.DateTime))
}
