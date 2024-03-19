package internal

import (
	"context"
	"time"

	"github.com/MingkaiLee/kasos/infer-module/util"
	jsoniter "github.com/json-iterator/go"
)

type CronJob struct {
	ticker  *time.Ticker
	Inferer *ParallelInferer
}

func NewCronJob(step int64, inferer *ParallelInferer) *CronJob {
	return &CronJob{
		ticker:  time.NewTicker(time.Second * time.Duration(step)),
		Inferer: inferer,
	}
}

func (c *CronJob) Start() {
	go func() {
		util.LogInfof("start infer loop, time: %s", time.Now().Format(time.DateTime))
		for range c.ticker.C {
			c.Inferer.Infer()
		}
	}()
}

func (c *CronJob) Stop() {
	c.ticker.Stop()
	util.LogInfof("stop infer loop, time: %s", time.Now().Format(time.DateTime))
}

type AddServiceRequest struct {
	ServiceName string `json:"service_name"`
	ModelName   string `json:"model_name"`
	Tags        string `json:"tags"`
}

func AddService(ctx context.Context, content []byte) (err error) {
	var req AddServiceRequest
	err = jsoniter.Unmarshal(content, &req)
	if err != nil {
		util.LogErrorf("unmarshal AddServiceRequest error: %v", err)
		return
	}
	InferCronJob.Inferer.AddService(req.ServiceName, req.ModelName, req.Tags)
	return
}
