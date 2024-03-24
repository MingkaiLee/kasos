package controller

import (
	"context"
	"net/http"

	"github.com/MingkaiLee/kasos/hpa-executor/internal"
	"github.com/MingkaiLee/kasos/hpa-executor/util"
	"github.com/cloudwego/hertz/pkg/app"
)

func initHpaExec() {
	grp := H.Group("/hpa-exec")
	grp.POST("/report-qps", ReportQPS)
}

// 上报QPS预测值辅助自动扩缩容的接口
func ReportQPS(ctx context.Context, c *app.RequestContext) {
	contentType := string(c.Request.Header.ContentType())
	if contentType != "application/json" {
		util.LogErrorf("content type not supported: %s", contentType)
		c.SetStatusCode(http.StatusUnsupportedMediaType)
		return
	}
	content, err := c.Body()
	if err != nil {
		util.LogErrorf("read body error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	err = internal.HpaExec(ctx, content)
	if err != nil {
		util.LogErrorf("report qps error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.SetStatusCode(http.StatusOK)
}
