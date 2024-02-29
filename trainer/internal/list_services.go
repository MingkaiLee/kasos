package internal

import (
	"context"
	"io"

	"github.com/MingkaiLee/kasos/trainer/client"
	"github.com/MingkaiLee/kasos/trainer/util"
	jsoniter "github.com/json-iterator/go"
)

type ListHpaServicesResponse struct {
	HpaServices []HpaService `json:"hpa_services"`
}

type HpaService struct {
	Name      *string           `json:"name"`
	Tags      map[string]string `json:"tags"`
	ThreshQPS *uint             `json:"thresh_qps"`
	ModelName *string           `json:"model_name"`
}

func ListServices() (services []HpaService, err error) {
	ctx := context.Background()
	response, err := client.CallListHpaServices(ctx)
	if err != nil {
		util.LogErrorf("internal.ListServices error: %v", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		util.LogErrorf("internal.ListServices error: %v", err)
		return
	}
	var servicesList ListHpaServicesResponse
	err = jsoniter.Unmarshal(body, &servicesList)
	if err != nil {
		util.LogErrorf("internal.ListServices error: %v", err)
		return
	}
	services = servicesList.HpaServices
	return
}
