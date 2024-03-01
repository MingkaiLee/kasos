package controller

import "sync"

type ParallelInferer struct {
	mu       sync.Mutex
	services []Service
}

type Service struct {
	ServiceName string
	ModelName   string
	Tags        string
}
