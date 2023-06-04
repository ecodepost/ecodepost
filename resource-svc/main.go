package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/task/ejob"

	"ecodepost/resource-svc/pkg/job"
	"ecodepost/resource-svc/pkg/server"
)

func main() {
	err := ego.New().Job(
		ejob.Job("init", job.RunInitData),
	).Cron().OrderServe(
		server.ServeGRPC(),
	).
		Run()
	if err != nil {
		elog.Panic(err.Error())
	}
}
