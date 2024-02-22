package controller

import (
	"context"
	"net/http"

	"github.com/MingkaiLee/kasos/server/service"
	"github.com/MingkaiLee/kasos/server/util"
	"github.com/cloudwego/hertz/pkg/app"
)

func initServiceManager() {
	// service-manager api
	serviceManager := H.Group("/service-manager")
	serviceManager.GET("/list", ListServices)
	serviceManager.GET("/find", FindService)
	serviceManager.POST("/register", RegisterService)
	serviceManager.GET("/register/result", RegisterServiceResult)
	serviceManager.POST("/delete", DeleteService)
	serviceManager.POST("/report/thresh", ReportThresh)
}

func ListServices(ctx context.Context, c *app.RequestContext) {
	r, err := service.ListHpaServices(ctx)
	if err != nil {
		util.LogErrorf("list service error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, r)
}

func FindService(ctx context.Context, c *app.RequestContext) {
	serviceName, ok := c.GetQuery("name")
	if !ok {
		util.LogErrorf("lack service name")
		c.NotFound()
		return
	}
	r, err := service.FindHpaService(ctx, serviceName)
	if err != nil {
		util.LogErrorf("find service error: %v", err)
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

func RegisterService(ctx context.Context, c *app.RequestContext) {
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
	r, err := service.RegisterService(ctx, content)
	if err != nil {
		util.LogErrorf("register service error: %v", err)
		c.JSON(http.StatusBadRequest, r)
		return
	}
	c.JSON(http.StatusOK, r)
}

func RegisterServiceResult(ctx context.Context, c *app.RequestContext) {
	serviceName, ok := c.GetQuery("name")
	if !ok {
		util.LogErrorf("lack service name")
		c.NotFound()
		return
	}
	r, err := service.QueryRegisterResult(ctx, serviceName)
	if err != nil {
		util.LogErrorf("query register result error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, r)
}

func DeleteService(ctx context.Context, c *app.RequestContext) {
	serviceName, ok := c.GetQuery("name")
	if !ok {
		util.LogErrorf("lack service name")
		c.NotFound()
		return
	}
	r, err := service.DeleteService(ctx, serviceName)
	if err != nil {
		util.LogErrorf("delete service error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, r)
}

func ReportThresh(ctx context.Context, c *app.RequestContext) {
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
	err = service.ReportThresh(ctx, content)
	if err != nil {
		util.LogErrorf("report thresh error: %v", err)
		c.SetStatusCode(http.StatusInternalServerError)
		return
	}
	c.SetStatusCode(http.StatusOK)
}
