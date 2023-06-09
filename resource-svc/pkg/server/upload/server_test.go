package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/invoker/alioss"
	"ecodepost/resource-svc/pkg/invoker/alists"

	uploadv1 "ecodepost/pb/upload/v1"

	"github.com/BurntSushi/toml"
	"github.com/ego-component/egorm"
	"github.com/ego-component/ek8s"
	k8sregistry "github.com/ego-component/ek8s/registry"
	"github.com/ego-component/eredis"
	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/server/egrpc"
	"github.com/stretchr/testify/assert"
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
	k8sregistry.Load("registry").Build(k8sregistry.WithClient(ek8s.Load("k8s").Build()))
	invoker.Db = egorm.Load("mysql").Build()
	invoker.Redis = eredis.Load("redis").Build()
	invoker.AliOss = alioss.Load("alioss").Build()
	invoker.AliSts = alists.Load("alists").Build()

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
	uploadv1.RegisterUploadServer(grpcCmp.Server, &GrpcServer{})
	return grpcCmp
}

func TestGetOssToken(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := uploadv1.NewUploadClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := cli.GetOssToken(ctx, &uploadv1.GetOssTokenReq{
		Uid:       999,
		ClientIp:  "0.0.0.0",
		Refer:     "",
		SpaceGuid: "spc-xxx",
	})
	resRes, _ := json.Marshal(res)
	fmt.Printf("res--------------->"+"%s\n", string(resRes))
	assert.NoError(t, err)
}
