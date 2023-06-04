package article

import (
	"context"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"

	articlev1 "ecodepost/pb/article/v1"
	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"

	"github.com/ego-component/egorm"
	"github.com/samber/lo"
)

// HomeArticlePageList 首页的文章列表根据不同方式展示
func (GrpcServer) HomeArticlePageList(ctx context.Context, req *articlev1.HomeArticlePageListReq) (*articlev1.HomeArticlePageListRes, error) {
	var spaceList mysql.Spaces
	var err error
	// 是不是登录用户
	// 如果是未登录用户找到所有公开的space数据
	if req.Uid == nil || *req.Uid == 0 {
		spaceList, err = mysql.SpacePublicList(invoker.Db.WithContext(ctx))
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("home article page list fail, err: " + err.Error())
		}
	} else {
		// 如果是登录用户，找到自己有权限的space数据
		spaceList, err = mysql.SpaceListByUser(invoker.Db.WithContext(ctx), req.GetUid())
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("home article page list fail2, err: " + err.Error())
		}
	}

	conds := egorm.Conds{}
	conds["space_guid"] = spaceList.ToGuids()
	// conds["space_top_time"] = 0
	conds["file_type"] = int(commonv1.FILE_TYPE_DOCUMENT)
	reqList := req.GetPagination()
	list, err := service.File.ListPageCache(ctx, invoker.Db.WithContext(ctx), conds, reqList, req.Sort)
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("home article page list fail3, err: " + err.Error())
	}

	return &articlev1.HomeArticlePageListRes{List: list.ToFileShowPb(), Pagination: reqList}, nil
}

// HomeArticleHotList 首页的热门文章列表
func (GrpcServer) HomeArticleHotList(ctx context.Context, req *articlev1.HomeArticleHotListReq) (*articlev1.HomeArticleHotListRes, error) {
	var spaceList mysql.Spaces
	var err error
	// 是不是登录用户
	// 如果是未登录用户找到所有公开的space数据
	if req.Uid == nil || *req.Uid == 0 {
		spaceList, err = mysql.SpacePublicList(invoker.Db.WithContext(ctx))
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("home article page list fail, err: " + err.Error())
		}
	} else {
		// 如果是登录用户，找到自己有权限的space数据
		spaceList, err = mysql.SpaceListByUser(invoker.Db.WithContext(ctx), req.GetUid())
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("home article page list fail2, err: " + err.Error())
		}
	}

	conds := egorm.Conds{}
	conds["space_guid"] = spaceList.ToGuids()
	// conds["space_top_time"] = 0
	conds["file_type"] = int(commonv1.FILE_TYPE_DOCUMENT)
	conds["dtime"] = 0
	// 创建时间大于多少天之前的显示
	conds["ctime"] = egorm.Cond{
		Op:  ">",
		Val: time.Now().Unix() - 86400*int64(req.GetLatestTime()),
	}
	sql, binds := egorm.BuildQuery(conds)

	var fileGuids []mysql.FileGuid
	err = invoker.Db.WithContext(ctx).Table("file").Select("guid").Where(sql, binds...).Order("recommend_score desc").Limit(int(req.GetLimit())).Find(&fileGuids).Error
	if err != nil {
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("home article page list fail3, err: " + err.Error())
		}
	}
	guids := lo.Map(fileGuids, func(v mysql.FileGuid, _ int) string { return v.Guid })

	list, err := service.File.ListCacheByGuids(ctx, invoker.Db.WithContext(ctx), guids)
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("home article page list fail4, err: " + err.Error())
	}

	return &articlev1.HomeArticleHotListRes{
		List: list.ToFileShowPb(),
	}, nil
}
