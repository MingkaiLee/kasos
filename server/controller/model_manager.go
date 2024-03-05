package controller

import (
	"context"
	"net/http"

	"github.com/MingkaiLee/kasos/server/service"
	"github.com/MingkaiLee/kasos/server/util"
	"github.com/cloudwego/hertz/pkg/app"
)

func initModelManager() {
	// model-manager api
	modelManager := H.Group("/model-manager")
	modelManager.GET("/list", ListModels)
	modelManager.GET("/find", FindModel)
	modelManager.POST("/register", RegisterModel)
	modelManager.GET("/register-result", RegisterModelResult)
	modelManager.POST("/report-valid", ReportModelValid)
}

func ListModels(ctx context.Context, c *app.RequestContext) {
	r, err := service.ListHpaModels(ctx)
	if err != nil {
		util.LogErrorf("list model error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, r)
}

func FindModel(ctx context.Context, c *app.RequestContext) {
	modelName, ok := c.GetQuery("name")
	if !ok {
		util.LogErrorf("lack model name")
		c.NotFound()
		return
	}
	r, err := service.FindHpaModel(ctx, modelName)
	if err != nil {
		util.LogErrorf("find model error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	if r == nil {
		util.LogErrorf("service not found")
		c.NotFound()
		return
	}
	c.JSON(http.StatusOK, r)
}

func RegisterModel(ctx context.Context, c *app.RequestContext) {
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
	r, err := service.RegisterModel(ctx, content)
	if err != nil {
		util.LogErrorf("register model error: %v", err)
		c.JSON(http.StatusBadRequest, r)
		return
	}
	c.JSON(http.StatusOK, r)
}

func RegisterModelResult(ctx context.Context, c *app.RequestContext) {
	modelName, ok := c.GetQuery("name")
	if !ok {
		util.LogErrorf("lack model name")
		c.NotFound()
		return
	}
	r, err := service.QueryModelResult(ctx, modelName)
	if err != nil {
		util.LogErrorf("query model result error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, r)
}

func ReportModelValid(ctx context.Context, c *app.RequestContext) {
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
	err = service.ReportModelValid(ctx, content)
	if err != nil {
		util.LogErrorf("report model valid error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.SetStatusCode(http.StatusOK)
}
