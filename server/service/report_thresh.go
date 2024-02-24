package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
	jsoniter "github.com/json-iterator/go"
)

type ReportThreshRequest struct {
	ServiceName string `json:"service_name"`
	OK          bool   `json:"ok"`
	ErrorInfo   string `json:"error_info"`
	ThreshQPS   uint   `json:"thresh_qps"`
}

func ReportThresh(ctx context.Context, content []byte) (err error) {
	var req ReportThreshRequest
	err = jsoniter.Unmarshal(content, &req)
	if err != nil {
		util.LogErrorf("unmarshal error: %v", err)
		return
	}

	if req.OK {
		err = model.HpaServiceRecordThreshQPS(req.ServiceName, req.ThreshQPS)
		if err != nil {
			util.LogErrorf("report thresh error: %v", err)
			return
		}
	} else {
		err = model.HpaServiceRecordError(req.ServiceName, req.ErrorInfo)
		if err != nil {
			util.LogErrorf("report error error: %v", err)
			return
		}
	}

	return
}
