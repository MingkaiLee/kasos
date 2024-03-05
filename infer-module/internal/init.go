package internal

import (
	"context"

	"github.com/MingkaiLee/kasos/infer-module/client"
	"github.com/MingkaiLee/kasos/infer-module/util"
)

var InferCronJob *CronJob

func InitInternal() {
	resp, err := client.CallListHpaServices(context.TODO())
	if err != nil {
		util.LogErrorf("Error when calling list hpa services")
	}
}
