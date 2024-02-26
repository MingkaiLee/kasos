package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MingkaiLee/kasos/trainer/client"
	"github.com/MingkaiLee/kasos/trainer/config"
)

type FetchDataWorker struct {
	workerNumber int    // 工作协程数量
	serviceName  string // 服务名
	date         string // 日期
}

func NewFetchDataWorker(workerNumber int, serviceName string, date time.Time) *FetchDataWorker {
	return &FetchDataWorker{
		workerNumber: workerNumber,
		serviceName:  serviceName,
		date:         date,
	}
}

func (worker *FetchDataWorker) Run() {
	ctx := context.Background()
	timeInterval := 24 / worker.workerNumber
	var g sync.WaitGroup
	for i := 0; i < worker.workerNumber; i++ {
		g.Add(1)
		go func(idx int) {
			startTime := time.Date()
			client.FetchSerialData(ctx, startTime, endTime, worker.serviceName)
			fileName := fmt.Sprint("{%s}/{%s}/{%s}_{%d}.csv",
				config.DataDirectory,
				worker.date,
				worker.serviceName,
				idx,
			)
			defer g.Done()
		}(i)
	}
	g.Wait()
}

// 合并数据至一个文件中
func MergeData(fileName string) {

}
