package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type ListHpaModelsResponse struct {
	HpaModels []HpaModel `json:"hpa_models"`
}

func ListHpaModels(ctx context.Context) (response *ListHpaModelsResponse, err error) {
	mds, err := model.HpaModelList()
	if err != nil {
		util.LogErrorf("failed to list hpa models, error: %v", err)
		return
	}
	response.HpaModels = make([]HpaModel, len(mds))
	for i := range mds {
		response.HpaModels[i] = HpaModel{
			Name:        &mds[i].ModelName,
			TrainScript: &mds[i].TrainScript,
			InferScript: &mds[i].InferScript,
		}
	}
	return
}
