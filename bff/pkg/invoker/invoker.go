package invoker

import (
	"log"
	"strings"

	"ecodepost/bff/pkg/server/bffcore/bffsso"
	articlev1 "ecodepost/pb/article/v1"
	columnv1 "ecodepost/pb/column/v1"
	commentv1 "ecodepost/pb/comment/v1"
	communityv1 "ecodepost/pb/community/v1"
	countv1 "ecodepost/pb/count/v1"
	filev1 "ecodepost/pb/file/v1"
	loggerv1 "ecodepost/pb/logger/v1"
	notifyv1 "ecodepost/pb/notify/v1"
	pmsv1 "ecodepost/pb/pms/v1"
	questionv1 "ecodepost/pb/question/v1"
	spacev1 "ecodepost/pb/space/v1"
	ssov1 "ecodepost/pb/sso/v1"
	statv1 "ecodepost/pb/stat/v1"
	uploadv1 "ecodepost/pb/upload/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/BurntSushi/toml"
	"github.com/ego-component/ek8s"
	k8sregistry "github.com/ego-component/ek8s/registry"
	"github.com/ego-component/eredis"
	"github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/econf"
)

var (
	Sso           *bffsso.Component
	Redis         *eredis.Component
	GrpcCommunity communityv1.CommunityClient
	GrpcComment   commentv1.CommentClient
	GrpcSpace     spacev1.SpaceClient
	GrpcFile      filev1.FileClient
	GrpcArticle   articlev1.ArticleClient
	GrpcColumn    columnv1.ColumnClient
	GrpcNotify    notifyv1.NotifyClient
	GrpcQuestion  questionv1.QuestionClient
	GrpcUpload    uploadv1.UploadClient
	GrpcUser      userv1.UserClient
	GrpcCount     countv1.CountClient
	GrpcStat      statv1.StatClient
	GrpcPms       pmsv1.PmsClient
	GrpcLogger    loggerv1.LoggerClient
	GrpcSso       ssov1.SsoClient
)

func Init() error {
	// 必须注册在grpc前面
	// 本地环境才启用，k8s resolver
	k8sregistry.Load("bff.registry").Build(k8sregistry.WithClient(ek8s.Load("bff.k8s").Build()))
	userConn := egrpc.Load("bff.grpc.user").Build()
	resourceConn := egrpc.Load("bff.grpc.resource").Build()
	GrpcCommunity = communityv1.NewCommunityClient(userConn)
	GrpcUser = userv1.NewUserClient(userConn)
	GrpcCount = countv1.NewCountClient(userConn)
	GrpcSso = ssov1.NewSsoClient(userConn)
	GrpcStat = statv1.NewStatClient(userConn)
	GrpcNotify = notifyv1.NewNotifyClient(userConn)
	GrpcSpace = spacev1.NewSpaceClient(resourceConn)
	GrpcFile = filev1.NewFileClient(resourceConn)
	GrpcArticle = articlev1.NewArticleClient(resourceConn)
	GrpcUpload = uploadv1.NewUploadClient(resourceConn)
	GrpcComment = commentv1.NewCommentClient(resourceConn)
	GrpcQuestion = questionv1.NewQuestionClient(resourceConn)
	GrpcPms = pmsv1.NewPmsClient(resourceConn)
	GrpcColumn = columnv1.NewColumnClient(resourceConn)
	GrpcLogger = loggerv1.NewLoggerClient(userConn)
	Redis = eredis.Load("redis").Build()
	Sso = bffsso.Load("bff.sso").Build()
	return nil
}

// InitByConf 测试用
func InitByConf(conf string) error {
	err := econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	if err != nil {
		log.Fatalf("init exited with error: %v", err)
	}
	return Init()
}
