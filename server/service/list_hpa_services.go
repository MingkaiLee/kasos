package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type ListHpaServicesResponse struct {
	HpaServices []HpaService `json:"hpa_services"`
	NextIndex   int          `json:"next_index"`
}

func ListHpaServices(ctx context.Context, startIndex uint) (response *ListHpaServicesResponse, err error) {
	svcs, err := model.HpaServiceList(startIndex)
	if err != nil {
		util.LogErrorf("failed to list hpa services, error: %v", err)
		return
	}
	response = new(ListHpaServicesResponse)
	if svcs == nil {
		response.HpaServices = make([]HpaService, 0)
		response.NextIndex = -1
		return
	}
	if len(svcs) < model.PageSize {
		response.NextIndex = -1
	} else {
		response.NextIndex = int(svcs[9].ID) + 1
	}
	response.HpaServices = make([]HpaService, len(svcs))
	for i := range svcs {
		response.HpaServices[i] = HpaService{
			Name:      &svcs[i].ServiceName,
			Tags:      util.RevertTags(svcs[i].Tags),
			ThreshQPS: &svcs[i].ThreshQPS,
			ModelName: &svcs[i].HpaModel.ModelName,
		}
	}
	return
}
