package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type HpaModel struct {
	Name        *string `json:"name"`
	TrainScript *string `json:"train_script"`
	InferScript *string `json:"infer_script"`
}

func FindHpaModel(ctx context.Context, modelName string) (response *HpaModel, err error) {
	md, err := model.HpaModelGet(modelName)
	if err != nil {
		util.LogErrorf("failed to get hpa model: %s, error: %s", modelName, err)
		return
	}
	response.Name = &md.ModelName
	response.TrainScript = &md.TrainScript
	response.InferScript = &md.InferScript
	return
}
