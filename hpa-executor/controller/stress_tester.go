package controller

import (
	"context"
	"net/http"

	"github.com/MingkaiLee/kasos/hpa-executor/util"
	"github.com/cloudwego/hertz/pkg/app"
)

func initStressTester() {
	grp := Server.Group("/stress-test")
	grp.POST("/normal-test", NormalTest)
}

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
}
