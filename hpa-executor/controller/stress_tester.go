package controller

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func initStressTester() {
	grp := Server.Group("/stress-test")
	grp.POST("/normal-test", normalTest)
}

func normalTest(ctx context.Context, c *app.RequestContext) {}
