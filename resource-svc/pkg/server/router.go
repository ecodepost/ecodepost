package server

import (
	articlev1 "ecodepost/pb/article/v1"
	columnv1 "ecodepost/pb/column/v1"
	commentv1 "ecodepost/pb/comment/v1"
	filev1 "ecodepost/pb/file/v1"
	pmsv1 "ecodepost/pb/pms/v1"
	questionv1 "ecodepost/pb/question/v1"
	spacev1 "ecodepost/pb/space/v1"
	uploadv1 "ecodepost/pb/upload/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/dao"
	"ecodepost/resource-svc/pkg/server/article"
	"ecodepost/resource-svc/pkg/server/column"
	"ecodepost/resource-svc/pkg/server/comment"
	"ecodepost/resource-svc/pkg/server/file"
	"ecodepost/resource-svc/pkg/server/pms"
	"ecodepost/resource-svc/pkg/server/question"
	"ecodepost/resource-svc/pkg/server/space"
	"ecodepost/resource-svc/pkg/server/upload"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/sdk/validator"

	"github.com/gotomicro/ego/server/egrpc"
)

func ServeGRPC() *egrpc.Component {
	srv := egrpc.Load("server.resource").Build(egrpc.WithUnaryInterceptor(validator.UnaryServerInterceptor()))
	srv.Invoker(
		invoker.Init,
		dao.InitGen,
		service.Init,
	)
	commentv1.RegisterCommentServer(srv.Server, &comment.GrpcServer{})
	uploadv1.RegisterUploadServer(srv.Server, &upload.GrpcServer{})
	spacev1.RegisterSpaceServer(srv.Server, &space.GrpcServer{})
	filev1.RegisterFileServer(srv.Server, &file.GrpcServer{})
	columnv1.RegisterColumnServer(srv.Server, &column.GrpcServer{})
	articlev1.RegisterArticleServer(srv.Server, &article.GrpcServer{})
	questionv1.RegisterQuestionServer(srv.Server, &question.GrpcServer{})
	pmsv1.RegisterPmsServer(srv.Server, &pms.GrpcServer{})
	return srv
}
