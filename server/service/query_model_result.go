package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type QueryModelResultResponse struct {
	Status    string `json:"status"`
	ErrorInfo string `json:"error_info"`
}

func QueryModelResult(ctx context.Context, modelName string) (response *QueryModelResultResponse, err error) {
	md, err := model.HpaModelGet(modelName)
	if err != nil {
		util.LogErrorf("failed to get model result, error: %v", err)
		return
	}
	response = new(QueryModelResultResponse)
	if md == nil {
		response.ErrorInfo = "unknown model name"
		return
	}

	response.Status = md.Status
	response.ErrorInfo = md.ErrorInfo

	return
}
