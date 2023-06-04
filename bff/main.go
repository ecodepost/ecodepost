package main

import (
	"ecodepost/bff/pkg/server"
	// _ "ecodepost/pb/activity/v1"
	// _ "ecodepost/pb/charge/v1"
	_ "ecodepost/pb/common/v1"
	_ "ecodepost/pb/community/v1"
	_ "ecodepost/pb/file/v1"
	// _ "ecodepost/pb/good/v1"
	_ "ecodepost/pb/notify/v1"
	// _ "ecodepost/pb/order/v1"
	_ "ecodepost/pb/pms/v1"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
)

func main() {
	err := ego.New().OrderServe(
		//egovernor.Load("server.governor").Build(),
		server.ServeHttp(),
	).Run()
	if err != nil {
		elog.Panic(err.Error())
	}
}
