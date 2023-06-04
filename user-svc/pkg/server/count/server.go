package count

import (
	"context"

	countv1 "ecodepost/pb/count/v1"
	"ecodepost/sdk/validator"
	"ecodepost/user-svc/pkg/server/count/count"
	"github.com/gotomicro/ego/server/egrpc"
)

var _ countv1.CountServer = (*GrpcServer)(nil)

type GrpcServer struct {
	Eng *count.Engine
}

func ServeGRPC() *egrpc.Component {
	srv := egrpc.Load("server.grpc").Build(egrpc.WithUnaryInterceptor(validator.UnaryServerInterceptor()))
	countv1.RegisterCountServer(srv.Server, &GrpcServer{
		Eng: count.New(),
	})
	return srv
}

// Set 增加计数关系 add、subtract、update
func (g *GrpcServer) Set(ctx context.Context, req *countv1.SetReq) (*countv1.SetRes, error) {
	return g.Eng.Set(ctx, req)
}

// GetTdetailsByTids 批量获取业务目标ID的数量(比如:评论的赞踩数量,资讯的推荐数量) 以及用户对业务的点赞或者收藏状态
func (g *GrpcServer) GetTdetailsByTids(ctx context.Context, req *countv1.GetTdetailsByTidsReq) (*countv1.GetTdetailsByTidsRes, error) {
	return g.Eng.GetTdetailsByTids(ctx, req)
}

// GetTdetailsByFid 通过业务和类型获取业务来源的操作列表(比如:用户的推荐列表记录,点赞列表记录,收藏列表记录)
func (g *GrpcServer) GetTdetailsByFid(ctx context.Context, req *countv1.GetTdetailsByFidReq) (*countv1.GetTdetailsByFidRes, error) {
	return g.Eng.GetTdetailsByFid(ctx, req)
}

// GetFidsTdetailByTid 获取一条资讯或者评论的推荐总数与推荐来源列表(比如:一条资讯被哪些人推荐过的一个列表及总数)
func (g *GrpcServer) GetFidsTdetailByTid(ctx context.Context, req *countv1.GetFidsTdetailByTidReq) (*countv1.GetFidsTdetailByTidRes, error) {
	return g.Eng.GetFidsTdetailByTid(ctx, req)
}

// GetTnumByBaksAndTids 批量获取目标ID计数总数值
func (g *GrpcServer) GetTnumByBaksAndTids(ctx context.Context, req *countv1.GetTnumByBaksAndTidsReq) (*countv1.GetTnumByBaksAndTidsRes, error) {
	return g.Eng.GetTnumByBaksAndTids(ctx, req)
}

// DBGetTdetailsByFid 类似 GetTdetailsByFid，不过是从数据库查询，可以查询完整列表
func (g *GrpcServer) DBGetTdetailsByFid(ctx context.Context, req *countv1.DBGetTdetailsByFidReq) (*countv1.DBGetTdetailsByFidRes, error) {
	return g.Eng.DBGetTdetailsByFid(ctx, req)
}

// DBGetFidsTdetailByTid 类似 GetFidsTdetailByTid，不过是从数据库查询，可以查询完整列表
func (g *GrpcServer) DBGetFidsTdetailByTid(ctx context.Context, req *countv1.DBGetFidsTdetailByTidReq) (*countv1.DBGetFidsTdetailByTidRes, error) {
	return g.Eng.DBGetFidsTdetailByTid(ctx, req)
}

// GetTnumByFids 使用场景：查询多个用户(A\A1\A2...)的TargetNum
func (g *GrpcServer) GetTnumByFids(ctx context.Context, req *countv1.GetTnumByFidsReq) (*countv1.GetTnumByFidsRes, error) {
	return g.Eng.GetTnumByFids(ctx, req)
}
