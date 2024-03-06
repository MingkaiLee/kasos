package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type ListHpaModelsResponse struct {
	HpaModels []HpaModel `json:"hpa_models"`
	NextIndex int        `json:"next_index"`
}

func ListHpaModels(ctx context.Context, startIndex uint) (response *ListHpaModelsResponse, err error) {
	mds, err := model.HpaModelList(startIndex)
	if err != nil {
		util.LogErrorf("failed to list hpa models, error: %v", err)
		return
	}
	response = new(ListHpaModelsResponse)
	if mds == nil {
		response.HpaModels = make([]HpaModel, 0)
		response.NextIndex = -1
		return
	}
	if len(mds) < model.PageSize {
		response.NextIndex = -1
	} else {
		response.NextIndex = int(mds[9].ID) + 1
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
