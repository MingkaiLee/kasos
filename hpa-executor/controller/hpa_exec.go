package controller

func initHpaExec() {
	grp := Server.Group("/hpa-exec")
	grp.POST("/report-qps")
}
