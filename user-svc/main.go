package main

import (
	"ecodepost/user-svc/pkg/server"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
)

func main() {
	app := ego.New()
	err := app.
		Serve(
			server.ServeGRPC(app),
		).
		Run()
	if err != nil {
		elog.Panic(err.Error())
	}
}
