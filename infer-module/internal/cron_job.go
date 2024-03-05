package internal

import (
	"time"

	"github.com/MingkaiLee/kasos/infer-module/util"
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
			go c.Inferer.Infer()
		}
	}()
}

func (c *CronJob) Stop() {
	c.ticker.Stop()
	util.LogInfof("stop infer loop, time: %s", time.Now().Format(time.DateTime))
}
