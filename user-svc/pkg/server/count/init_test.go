package count

import (
	"log"
	"strings"
	"time"

	countv1 "ecodepost/pb/count/v1"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/server/count/count"
	"github.com/BurntSushi/toml"
	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/server/egrpc"
	"google.golang.org/grpc"
)

// svc generated by protoc-gen-go-test, you should not edit it.
// @Override=true
var svc *egrpc.Component

// init generated by protoc-gen-go-test, you can fill initial logic by yourself.
// @Override=true
func init() {
	conf := `

`
	// 加载配置
	err := econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	if err != nil {
		log.Fatalf("init exited with error: %v", err)
	}
	invoker.Init()
	// 初始化bufnet gRPC的测试服务
	svc = Server()
	err = svc.Init()
	if err != nil {
		log.Fatalf("init server with error: %v", err)
	}

	go func() {
		// 启动服务
		if err = svc.Start(); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func Server() *egrpc.Component {
	grpcCmp := egrpc.Load("server.grpc").Build()
	countv1.RegisterCountServer(grpcCmp.Server, &GrpcServer{Eng: count.New()})
	return grpcCmp
}

func newCC() grpc.ClientConnInterface {
	return cegrpc.DefaultContainer().Build(
		cegrpc.WithBufnetServerListener(svc.Listener()),
		cegrpc.WithDialTimeout(5*time.Second),
		cegrpc.WithReadTimeout(15*time.Second),
	).ClientConn
}
