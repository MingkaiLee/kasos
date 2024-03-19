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
	NextIndex   int          `json:"next_index"`
}

type HpaService struct {
	Name      *string           `json:"name"`
	Tags      map[string]string `json:"tags"`
	ThreshQPS *uint             `json:"thresh_qps"`
	ModelName *string           `json:"model_name"`
}

func ListServices() (services []HpaService, err error) {
	services = make([]HpaService, 0)
	ctx := context.Background()
	index := 0
	for index >= 0 {
		response, e := client.CallListHpaServices(ctx, index)
		if e != nil {
			err = e
			util.LogErrorf("internal.ListServices error: %v", err)
			return
		}
		defer response.Body.Close()
		body, e := io.ReadAll(response.Body)
		if e != nil {
			err = e
			util.LogErrorf("internal.ListServices error: %v", err)
			return
		}
		var servicesList ListHpaServicesResponse
		err = jsoniter.Unmarshal(body, &servicesList)
		if err != nil {
			util.LogErrorf("internal.ListServices error: %v", err)
			return
		}
		services = append(services, servicesList.HpaServices...)
		index = servicesList.NextIndex
	}
	return
}
