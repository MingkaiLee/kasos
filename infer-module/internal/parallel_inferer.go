package internal

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/MingkaiLee/kasos/infer-module/util"
)

type ParallelInferer struct {
	mu       sync.Mutex
	services []*Service
}

func NewParallelInferer(services []HpaService) *ParallelInferer {
	inferer := &ParallelInferer{
		mu:       sync.Mutex{},
		services: make([]*Service, 0),
	}
	for _, hpaService := range services {
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
	defer func() {
		if r := recover(); r != nil {
			util.LogErrorf("AddService panic: %v", r)
			util.LogErrorf("Stack trace: %s", string(debug.Stack()))
		}
	}()
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
