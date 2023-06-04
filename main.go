package main

import (
	bffserver "ecodepost/bff/pkg/server"
	"ecodepost/job"
	resourceserver "ecodepost/resource-svc/pkg/server"
	userserver "ecodepost/user-svc/pkg/server"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/task/ejob"
)

func main() {
	app := ego.New()
	err := app.OrderServe(
		userserver.ServeGRPC(app),
		resourceserver.ServeGRPC(),
		bffserver.ServeHttp(),
	).Job(
		ejob.Job("install", job.RunInstall),
		ejob.Job("init", job.InitData),
	).Run()
	if err != nil {
		elog.Panic(err.Error())
	}
}
