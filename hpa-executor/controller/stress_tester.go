package controller

import (
	"context"
	"net/http"

	"github.com/MingkaiLee/kasos/hpa-executor/internal"
	"github.com/MingkaiLee/kasos/hpa-executor/util"
	"github.com/cloudwego/hertz/pkg/app"
)

func initStressTester() {
	grp := H.Group("/stress-test")
	grp.POST("/normal-test", NormalTest)
}

// 压测接口
func NormalTest(ctx context.Context, c *app.RequestContext) {
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
	err = internal.NormalTest(ctx, content)
	if err != nil {
		util.LogErrorf("normal test error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.SetStatusCode(http.StatusOK)
}
