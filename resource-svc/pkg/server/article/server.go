package article

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/service"

	articlev1 "ecodepost/pb/article/v1"
	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"

	"github.com/ego-component/egorm"
)

type GrpcServer struct{}

var _ articlev1.ArticleServer = (*GrpcServer)(nil)

func (GrpcServer) ListDocumentByGuids(ctx context.Context, req *articlev1.ListDocumentByGuidsReq) (*articlev1.ListDocumentByGuidsRes, error) {
	if req.Guids == nil || len(req.Guids) == 0 {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("invalid guids")
	}
	var conds = make(egorm.Conds)
	conds["guid"] = req.Guids
	list, err := service.File.ListPageCache(ctx, invoker.Db.WithContext(ctx), conds, req.Pagination, commonv1.CMN_SORT_CREATE_TIME)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("ListDocumentByGuids document fail, err: " + err.Error())
	}
	return &articlev1.ListDocumentByGuidsRes{
		List:       list.ToArticlePb(commonv1.SPC_LAYOUT_ARTICLE_CARD),
		Pagination: req.GetPagination(),
	}, nil
}
