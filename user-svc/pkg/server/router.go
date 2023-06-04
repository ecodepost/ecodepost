package server

import (
	job2 "ecodepost/job"
	communityv1 "ecodepost/pb/community/v1"
	countv1 "ecodepost/pb/count/v1"
	notifyv1 "ecodepost/pb/notify/v1"
	ssov1 "ecodepost/pb/sso/v1"
	"ecodepost/sdk/validator"
	"ecodepost/user-svc/pkg/cron"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/server/community"
	"ecodepost/user-svc/pkg/server/count"
	countpkg "ecodepost/user-svc/pkg/server/count/count"
	"ecodepost/user-svc/pkg/server/logger"
	"ecodepost/user-svc/pkg/server/notify"
	"ecodepost/user-svc/pkg/server/sso"
	"ecodepost/user-svc/pkg/server/stat"
	"ecodepost/user-svc/pkg/server/user"
	"ecodepost/user-svc/pkg/service"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/task/ejob"

	loggerv1 "ecodepost/pb/logger/v1"
	statv1 "ecodepost/pb/stat/v1"
	userv1 "ecodepost/pb/user/v1"
	"github.com/gotomicro/ego/server/egrpc"
)

func ServeGRPC(app *ego.Ego) *egrpc.Component {
	// 因为没有gRPC依赖，所以可以提前初始化invoker
	srv := egrpc.Load("server.user").Build(egrpc.WithUnaryInterceptor(validator.UnaryServerInterceptor()))
	srv.Invoker(
		invoker.Init,
		service.Init,
		func() error {
			app.Job(
				ejob.Job("init_sso", job2.InitSsoData),
				ejob.Job("init_system", job2.InitSystem),
				ejob.Job("update_sso", job2.UpdateSsoData),
			)
			return nil
		},
		func() error {
			// 这里有gRPC依赖于invoker里的启动项，需要放到这里
			// grpc 一个应用，一个连接
			// grpc 多个service
			userv1.RegisterUserServer(srv.Server, &user.GrpcServer{})
			statv1.RegisterStatServer(srv.Server, &stat.GrpcServer{})
			loggerv1.RegisterLoggerServer(srv.Server, &logger.GrpcServer{})
			ssov1.RegisterSsoServer(srv.Server, &sso.GrpcServer{})
			communityv1.RegisterCommunityServer(srv.Server, &community.GrpcServer{})
			notifyv1.RegisterNotifyServer(srv.Server, &notify.GrpcServer{})
			countv1.RegisterCountServer(srv.Server, &count.GrpcServer{
				Eng: countpkg.New(),
			})
			return nil
		},
		func() error {
			app.Cron(cron.NotifyCron())
			return nil
		},
	)

	return srv
}
