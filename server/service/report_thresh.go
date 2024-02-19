package service

import (
	"context"
	"encoding/json"

	"github.com/MingkaiLee/kasos/server/model"
)

type ReportThreshRequest struct {
	ServiceName string `json:"service_name"`
	OK          bool   `json:"ok"`
	ErrorInfo   string `json:"error_info"`
	ThreshQPS   uint   `json:"thresh_qps"`
}

func ReportThresh(ctx context.Context, content []byte) (err error) {
	var req ReportThreshRequest
	err = json.Unmarshal(content, &req)
	if err != nil {
		return
	}

	if req.OK {
		err = model.HpaServiceRecordThreshQPS(req.ServiceName, req.ThreshQPS)
		if err != nil {
			return
		}
	} else {
		err = model.HpaServiceRecordError(req.ServiceName, req.ErrorInfo)
		if err != nil {
			return
		}
	}

	return
}
