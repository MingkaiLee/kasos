package service

import (
	"context"

	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
	jsoniter "github.com/json-iterator/go"
)

type ReportModelValidRequest struct {
	ModelName string `json:"model_name"`
	Ok        bool   `json:"ok"`
	ErrorInfo string `json:"error_info"`
}

func ReportModelValid(ctx context.Context, content []byte) (err error) {
	var req ReportModelValidRequest
	err = jsoniter.Unmarshal(content, &req)
	if err != nil {
		util.LogErrorf("unmarshal error: %v", err)
		return
	}

	if req.Ok {
		err = model.HpaModelRecordOk(req.ModelName)
		if err != nil {
			util.LogErrorf("report model valid error: %v", err)
			return
		}
	} else {
		err = model.HpaModelRecordError(req.ModelName, req.ErrorInfo)
		if err != nil {
			util.LogErrorf("report model valid error: %v", err)
			return
		}
	}

	return
}
