package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
	jsoniter "github.com/json-iterator/go"
)

type RegisterModelRequest struct {
	Name        string `json:"name"`
	TrainScript string `json:"train_script"`
	InferScript string `json:"infer_script"`
}

type RegisterModelResponse struct {
	Accepted bool   `json:"accepted"`
	Message  string `json:"message"`
}

func RegisterModel(ctx context.Context, content []byte) (response *RegisterModelResponse, err error) {
	var req RegisterModelRequest
	err = jsoniter.Unmarshal(content, &req)
	response = new(RegisterModelResponse)
	if err != nil {
		util.LogErrorf("failed to unmarshal request, error: %v", err)
		response.Message = err.Error()
		return
	}
	util.LogInfof("register service request: %+v", req)

	// register model
	err = model.HpaModelCreate(req.Name, req.TrainScript, req.InferScript)
	if err != nil {
		util.LogErrorf("failed to create model, error: %v", err)
		response.Message = err.Error()
	}

	// TODO: send a model validation request to infer-module

	response.Accepted = true
	return
}
