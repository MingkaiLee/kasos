package service

import (
	"context"
	"encoding/json"

	"github.com/MingkaiLee/kasos/server/client"
	"github.com/MingkaiLee/kasos/server/model"
)

type RegisterServiceRequest struct {
	Name      string            `json:"name"`
	Tags      map[string]string `json:"tags"`
	ModelName string            `json:"model_name"`
}

type RegisterServiceResponse struct {
	Accepted bool   `json:"accepted"`
	Message  string `json:"message"`
}

func RegisterService(ctx context.Context, content []byte) (response *RegisterServiceResponse, err error) {
	var req RegisterServiceRequest
	err = json.Unmarshal(content, &req)
	if err != nil {
		response.Message = err.Error()
		return
	}

	// look up model information
	modelID, err := model.HpaModelGetID(req.ModelName)
	if err != nil {
		response.Message = err.Error()
		return
	}

	// register service
	err = model.HpaServiceCreate(req.Name, req.Tags, modelID)
	if err != nil {
		response.Message = err.Error()
		return
	}

	// TODO: send stress test request to hpa-executor

	// TODO: create a Prometheus ServiceMonitor
	err = client.CreateMonitorService(ctx, req.Name, req.Tags)
	if err != nil {
		response.Message = err.Error()
		return
	}

	response.Accepted = true
	return
}
