package main

import (
	"fmt"

	"github.com/MingkaiLee/kasos/infer-module/internal"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	a := new(internal.ListHpaServicesResponse)
	a.HpaServices = make([]internal.HpaService, 0)
	// config.InitConf()
	r, e := jsoniter.Marshal(a)
	fmt.Printf("%s, error: %v", string(r), e)
}
