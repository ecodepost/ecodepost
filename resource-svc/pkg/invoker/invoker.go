package invoker

import (
	communityv1 "ecodepost/pb/community/v1"
	loggerv1 "ecodepost/pb/logger/v1"
	notifyv1 "ecodepost/pb/notify/v1"
	"ecodepost/sdk/oss"
	"github.com/ego-component/eguid"

	statv1 "ecodepost/pb/stat/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/ego-component/egorm"
	"github.com/ego-component/ek8s"
	k8sregistry "github.com/ego-component/ek8s/registry"
	"github.com/ego-component/eredis"
	"github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/elog"
)

var (
	Logger        *elog.Component
	Db            *egorm.Component
	Redis         *eredis.Component
	AliSts        *oss.Component
	Guid          *eguid.Component
	GrpcUser      userv1.UserClient
	GrpcCommunity communityv1.CommunityClient
	GrpcStat      statv1.StatClient
	GrpcLogger    loggerv1.LoggerClient
	GrpcNotify    notifyv1.NotifyClient
)

func Init() error {
	Logger = elog.DefaultLogger
	Db = egorm.Load("mysql").Build()
	Redis = eredis.Load("redis").Build()
	Guid = eguid.Load("resource-svc.guid").Build()
	k8sregistry.Load("resource-svc.registry").Build(k8sregistry.WithClient(ek8s.Load("resource-svc.k8s").Build()))
	userConn := egrpc.Load("grpc.user").Build()
	GrpcUser = userv1.NewUserClient(userConn)
	GrpcStat = statv1.NewStatClient(userConn)
	GrpcLogger = loggerv1.NewLoggerClient(userConn)
	GrpcCommunity = communityv1.NewCommunityClient(userConn)
	GrpcNotify = notifyv1.NewNotifyClient(userConn)
	AliSts = oss.Load("oss").Build()
	return nil
}
