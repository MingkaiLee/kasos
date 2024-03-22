package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/MingkaiLee/kasos/trainer/client"
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
		dataRow := fmt.Sprintf("%s\t%.4f\n", serialData[i].Timestamp, serialData[i].Value)
		_, err = fp.WriteString(dataRow)
		if err != nil {
			util.LogErrorf("internal.serialDataSave error: %v", err)
			return
		}
	}

	return
}

// 文件按照名称字符排序后, 统一合并到第一个文件中
func mergeFiles(subFiles []string) (err error) {
	if len(subFiles) == 0 {
		err = fmt.Errorf("subFiles is empty")
		util.LogErrorf("internal.mergeFiles error: %v", err)
		return
	}
	// 排序子文件名, 保证数据是有序的
	sort.Strings(subFiles)
	// 将子文件合并, 其他子文件都放入第一个子文件的末尾
	dst, err := os.OpenFile(subFiles[0], os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		util.LogErrorf("internal.mergeFiles error: %v", err)
		return
	}
	defer dst.Close()
	for i := 1; i < len(subFiles); i++ {
		src, e := os.Open(subFiles[i])
		if e != nil {
			err = e
			util.LogErrorf("internal.mergeFiles error: %v", err)
			return
		}
		defer src.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			util.LogErrorf("internal.mergeFiles error: %v", err)
			return
		}
	}
	return
}

// 合并数据至一个文件中
func mergeData(finalFileName string, subFiles []string) (err error) {
	// 先调用函数合并至第一个文件中
	err = mergeFiles(subFiles)
	if err != nil {
		util.LogErrorf("internal.mergeData error: %v", err)
		return
	}
	// 修改第一个文件的名称
	err = os.Rename(subFiles[0], finalFileName)
	if err != nil {
		util.LogErrorf("internal.mergeData error: %v", err)
		return
	}
	// 删除其他文件
	for i := 1; i < len(subFiles); i++ {
		err = os.Remove(subFiles[i])
		if err != nil {
			util.LogErrorf("internal.mergeData error: %v", err)
			return
		}
	}
	return
}

func (worker *FetchDataWorker) Run(dirName string) (err error) {
	util.LogInfof("internal.FetchDataWorker.Run: start to fetch data of service %s", worker.serviceName)
	// 测试版的训练周期是1小时
	ctx := context.Background()
	finalFileName := fmt.Sprintf("%s/%s.csv", dirName, worker.serviceName)
	// 时间段的中止时间为当前时间截断到小时
	endTime := worker.date.Truncate(time.Hour)
	// 时间段的开始时间为结束时间减去1小时
	startTime := endTime.Add(-time.Hour)
	// 开始拉取数据
	serialData, err := client.FetchSerialData(ctx, startTime, endTime, worker.serviceName)
	if err != nil {
		util.LogErrorf("internal.FetchDataWorker.Run error: %v", err)
		return
	}
	// 保存数据
	err = serialDataSave(serialData, finalFileName)

	return
	// timeInterval := 24 / worker.workerNumber
	// if (timeInterval % worker.workerNumber) > 0 {
	// 	timeInterval++
	// }
	// // 开始时间, 以小时为单位, 0-24
	// startInHour := 0
	// year, month, day := worker.date.Date()
	// finalFileName := fmt.Sprintf("%s/%d-%d-%d/%s.csv",
	// 	config.DataDirectory,
	// 	year, month, day,
	// 	worker.serviceName,
	// )
	// subFiles := make([]string, 0, worker.workerNumber)
	// var g sync.WaitGroup
	// var mu sync.Mutex
	// for i := 0; startInHour < 24; i++ {
	// 	g.Add(1)
	// 	st := startInHour
	// 	et := startInHour + timeInterval
	// 	if et > 24 {
	// 		et = 24
	// 	}
	// 	startInHour = et
	// 	go func(start, end, idx int) {
	// 		defer g.Done()
	// 		// 并发查询与写文件
	// 		startTime := worker.date.Add(time.Duration(start) * time.Hour)
	// 		endTime := worker.date.Add(time.Duration(end) * time.Hour)
	// 		serialData, err := client.FetchSerialData(ctx, startTime, endTime, worker.serviceName)
	// 		if err != nil {
	// 			util.LogErrorf("internal.FetchDataWorker.Run goroutine error: %v", err)
	// 			return
	// 		}
	// 		fileName := fmt.Sprintf("%s/%d-%d-%d/%s_%d.csv",
	// 			config.DataDirectory,
	// 			year, month, day,
	// 			worker.serviceName,
	// 			idx,
	// 		)
	// 		err = serialDataSave(serialData, fileName)
	// 		if err != nil {
	// 			util.LogErrorf("internal.FetchDataWorker.Run goroutine error: %v", err)
	// 			return
	// 		}
	// 		// 防止写冲突
	// 		mu.Lock()
	// 		subFiles = append(subFiles, fileName)
	// 		mu.Unlock()
	// 	}(st, et, i)
	// }
	// g.Wait()
	// // 合并文件
	// err = mergeData(finalFileName, subFiles)
	// if err != nil {
	// 	util.LogErrorf("internal.FetchDataWorker.Run error: %v", err)
	// 	return
	// }
	// util.LogInfof("internal.FetchDataWorker.Run: fetch data of service %s done, saved to %s", worker.serviceName, finalFileName)
	// return
}
