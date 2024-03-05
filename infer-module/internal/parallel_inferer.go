package internal

import (
	"context"
	"sync"

	"github.com/MingkaiLee/kasos/infer-module/util"
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

type ParallelInferer struct {
	mu       sync.Mutex
	services []*Service
}

func NewParallerInferer(r *ListHpaServicesResponse) *ParallelInferer {
	inferer := &ParallelInferer{}
	for _, hpaService := range r.HpaServices {
		inferer.AddService(*hpaService.Name, *hpaService.ModelName, util.ConvertTags(hpaService.Tags))
	}
	return inferer
}

func (p *ParallelInferer) Infer() {
	p.mu.Lock()
	defer p.mu.Unlock()
	// 触发推理
	var g sync.WaitGroup
	for idx := range p.services {
		g.Add(1)
		go func(i int) {
			defer g.Done()
			err := p.services[i].Run(context.TODO())
			if err != nil {
				util.LogErrorf("infer service %s failed, error: %v", p.services[i].GetName(), err)
				return
			}
			util.LogInfof("infer service %s success", p.services[i].GetName())
		}(idx)
	}
	g.Wait()
}

// 添加新服务
func (p *ParallelInferer) AddService(serviceName, modelName, tags string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.services = append(p.services, NewService(serviceName, modelName, tags))
}

// 删除服务
func (p *ParallelInferer) RemoveService(serviceName string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for idx := range p.services {
		if p.services[idx].GetName() == serviceName {
			p.services = append(p.services[:idx], p.services[idx+1:]...)
			return
		}
	}
}
