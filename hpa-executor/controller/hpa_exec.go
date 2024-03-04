package controller

func initHpaExec() {
	grp := H.Group("/hpa-exec")
	grp.POST("/report-qps")
}
