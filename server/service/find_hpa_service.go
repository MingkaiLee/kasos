package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type HpaService struct {
	Name      *string           `json:"name"`
	Tags      map[string]string `json:"tags"`
	ThreshQPS *uint             `json:"thresh_qps"`
	ModelName *string           `json:"model_name"`
}

func FindHpaService(ctx context.Context, serviceName string) (response *HpaService, err error) {
	svc, err := model.HpaServiceGet(serviceName)
	if err != nil {
		util.LogErrorf("failed to get hpa service: %s, error: %v", serviceName, err)
		return nil, err
	}
	response.Name = &svc.ServiceName
	response.Tags = util.RevertTags(svc.Tags)
	response.ThreshQPS = &svc.ThreshQPS
	response.ModelName = &svc.HpaModel.ModelName

	return
}
