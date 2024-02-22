package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type QueryRegisterResultResponse struct {
	Status    string `json:"status"`
	ErrorInfo string `json:"error_info"`
}

func QueryRegisterResult(ctx context.Context, serviceName string) (response *QueryRegisterResultResponse, err error) {
	svc, err := model.HpaServiceGet(serviceName)
	if err != nil {
		util.LogErrorf("failed to get register result, error: %v", err)
		return
	}

	response.Status = svc.Status
	response.ErrorInfo = svc.ErrorInfo

	return
}
