package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
)

type QueryRegisterResultResponse struct {
	Status    string `json:"status"`
	ErrorInfo string `json:"error_info"`
}

func QueryRegisterResult(ctx context.Context, serviceName string) (response *QueryRegisterResultResponse, err error) {
	svc, err := model.HpaServiceGet(serviceName)
	if err != nil {
		return
	}

	response.Status = svc.Status
	response.ErrorInfo = svc.ErrorInfo

	return
}
