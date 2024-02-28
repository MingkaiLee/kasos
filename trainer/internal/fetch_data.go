package internal

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/MingkaiLee/kasos/trainer/client"
	"github.com/MingkaiLee/kasos/trainer/config"
	"github.com/MingkaiLee/kasos/trainer/util"
)

type FetchDataWorker struct {
	workerNumber int       // 工作协程数量
	serviceName  string    // 服务名
	date         time.Time // 日期
}

func NewFetchDataWorker(workerNumber int, serviceName string, date time.Time) *FetchDataWorker {
	if workerNumber < 1 {
		workerNumber = 1
	}
	if workerNumber > 24 {
		workerNumber = 24
	}
	return &FetchDataWorker{
		workerNumber: workerNumber,
		serviceName:  serviceName,
		date:         date,
	}
}

func serialDataSave(serialData []client.SerialDataPoint, fileName string) (err error) {
	fp, err := os.Create(fileName)
	if err != nil {
		util.LogErrorf("internal.serialDataSave error: %v", err)
		return
	}
	defer fp.Close()

	for i := range serialData {
		dataRow := fmt.Sprintf("%s\t%.2f\n", serialData[i].Timestamp, serialData[i].Value)
		_, err = fp.WriteString(dataRow)
		if err != nil {
			util.LogErrorf("internal.serialDataSave error: %v", err)
			return
		}
	}

	return
}

// 合并数据至一个文件中
func mergeData(finalFileName string, subFiles []string) (err error) {
	// 排序文件名, 保证数据是有序的
	sort.Strings(subFiles)
	return
}

func (worker *FetchDataWorker) Run() {
	ctx := context.Background()
	timeInterval := 24 / worker.workerNumber
	if (timeInterval % worker.workerNumber) > 0 {
		timeInterval++
	}
	// 开始时间, 以小时为单位, 0-24
	startInHour := 0
	year, month, day := worker.date.Date()
	finalFileName := fmt.Sprintf("{%s}/{%d}-{%d}-{%d}/{%s}.csv",
		config.DataDirectory,
		year, month, day,
		worker.serviceName,
	)
	subFiles := make([]string, 0, worker.workerNumber)
	var g sync.WaitGroup
	var mu sync.Mutex
	for i := 0; startInHour < 24; i++ {
		g.Add(1)
		st := startInHour
		et := startInHour + timeInterval
		if et > 24 {
			et = 24
		}
		startInHour = et
		go func(start, end, idx int) {
			defer g.Done()
			// 并发查询与写文件
			startTime := worker.date.Add(time.Duration(start) * time.Hour)
			endTime := worker.date.Add(time.Duration(end) * time.Hour)
			serialData, err := client.FetchSerialData(ctx, startTime, endTime, worker.serviceName)
			if err != nil {
				util.LogErrorf("internal.FetchDataWorker.Run goroutine error: %v", err)
				return
			}
			fileName := fmt.Sprintf("{%s}/{%d}-{%d}-{%d}/{%s}_{%d}.csv",
				config.DataDirectory,
				year, month, day,
				worker.serviceName,
				idx,
			)
			err = serialDataSave(serialData, fileName)
			if err != nil {
				util.LogErrorf("internal.FetchDataWorker.Run goroutine error: %v", err)
				return
			}
			mu.Lock()
			subFiles = append(subFiles, fileName)
			mu.Unlock()
		}(st, et, i)
	}
	g.Wait()
	// 合并文件
	err := mergeData(finalFileName, subFiles)
	if err != nil {
		util.LogErrorf("internal.FetchDataWorker.Run error: %v", err)
		return
	}
}
