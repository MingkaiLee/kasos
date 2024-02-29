package controller

import (
	"context"
	"net/http"

	"github.com/MingkaiLee/kasos/server/service"
	"github.com/MingkaiLee/kasos/server/util"
	"github.com/cloudwego/hertz/pkg/app"
)

func initDataManager() {
	// data-manager api
	dataManager := H.Group("/data-manager")
	dataManager.POST("/fetch")
}

func FetchSerialData(ctx context.Context, c *app.RequestContext) {
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
	r, err := service.FetchSerialData(ctx, content)
	if err != nil {
		util.LogErrorf("fetch serial data error: %v", err)
		c.JSON(http.StatusBadRequest, r)
		return
	}
	c.JSON(http.StatusOK, r)
}
