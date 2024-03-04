package controller

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func initHpaExec() {
	grp := H.Group("/hpa-exec")
	grp.POST("/report-qps")
}

// 上报QPS预测值辅助自动扩缩容的接口
func ReportQPS(ctx context.Context, c *app.RequestContext)
