package controller

import (
	"context"
	"net/http"

	"github.com/MingkaiLee/kasos/infer-module/internal"
	"github.com/MingkaiLee/kasos/infer-module/util"
	"github.com/cloudwego/hertz/pkg/app"
)

func initModel() {
	grp := H.Group("/model")
	grp.POST("/validate", ModelValidate)
}

func ModelValidate(ctx context.Context, c *app.RequestContext) {
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
	err = internal.ScriptValidate(ctx, content)
	if err != nil {
		util.LogErrorf("script validate error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.SetStatusCode(http.StatusOK)
}
