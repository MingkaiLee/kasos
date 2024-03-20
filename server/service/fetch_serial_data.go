package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MingkaiLee/kasos/server/client"
	"github.com/MingkaiLee/kasos/server/util"
	jsoniter "github.com/json-iterator/go"
)

type FetchSerialDataRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Tags      string `json:"tags"`
}

type FetchSerialDataResponse struct {
	Message    string              `json:"message"`
	SerialData map[string][]string `json:"serial_data"`
}

func FetchSerialData(ctx context.Context, content []byte) (response *FetchSerialDataResponse, err error) {
	var req FetchSerialDataRequest
	err = jsoniter.Unmarshal(content, &req)
	response = new(FetchSerialDataResponse)
	util.LogInfof("fetch serial data request: %+v", req)
	if err != nil {
		util.LogErrorf("failed to unmarshal request, error: %v", err)
		response.Message = err.Error()
		return
	}
	startTime, err := time.Parse(req.StartTime, time.DateTime)
	if err != nil {
		util.LogErrorf("failed to parse start time, error: %v", err)
		response.Message = err.Error()
		return
	}
	endTime, err := time.Parse(req.EndTime, time.DateTime)
	if err != nil {
		util.LogErrorf("failed to parse end time, error: %v", err)
		response.Message = err.Error()
		return
	}
	tags := util.RevertTags(req.Tags)
	util.LogInfof("start: %v, end: %v, tags: %v", startTime, endTime, tags)
	data, err := client.FetchSerialData(ctx, startTime, endTime, tags)
	if err != nil {
		util.LogErrorf("failed to fetch serial data, error: %v", err)
		response.Message = err.Error()
		return
	}
	response.SerialData = make(map[string][]string)
	for k, v := range data {
		s := make([]string, len(v))
		for idx := range v {
			s[idx] = fmt.Sprintf("%s\t%.4f", v[idx].Timestamp, v[idx].Value)
		}
		response.SerialData[k] = s
	}
	response.Message = "success"
	return
}
