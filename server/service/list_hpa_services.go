package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
)

type ListHpaServicesResponse struct {
	HpaServices []HpaService `json:"hpa_services"`
}

func ListHpaServices(ctx context.Context) (response *ListHpaServicesResponse, err error) {
	svcs, err := model.HpaServiceList()
	if err != nil {
		util.LogErrorf("failed to list hpa services, error: %v", err)
		return
	}
	response = new(ListHpaServicesResponse)
	if svcs == nil {
		return
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
