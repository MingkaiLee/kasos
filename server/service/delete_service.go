package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type DeleteServiceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func DeleteService(ctx context.Context, serviceName string) (response *DeleteServiceResponse, err error) {
	err = model.HpaServiceDelete(serviceName)
	response = new(DeleteServiceResponse)
	if err != nil {
		util.LogErrorf("delete service error: %v", err)
		response.Message = err.Error()
		return
	}
	response.Success = true
	return
}
