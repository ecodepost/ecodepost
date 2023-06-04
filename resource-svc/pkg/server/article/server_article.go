package article

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/cache"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	articlev1 "ecodepost/pb/article/v1"

	"gorm.io/gorm"

	commonv1 "ecodepost/pb/common/v1"

	errcodev1 "ecodepost/pb/errcode/v1"
)

// SetDocumentSpaceTop 将文章置顶
func (GrpcServer) SetDocumentSpaceTop(ctx context.Context, req *articlev1.SetDocumentSpaceTopReq) (*articlev1.SetDocumentSpaceTopRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_SET_SITE_TOP, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}

	spaceGuid, err := mysql.FileSpaceGuidByGuid(invoker.Db.WithContext(ctx), req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileSpaceGuidByGuid fail, err: " + err.Error())
	}

	cnt, err := mysql.FileSpaceTopCnt(invoker.Db.WithContext(ctx), spaceGuid)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileSpaceTopCnt fail, err: " + err.Error())
	}
	// 目前不让大于1
	if cnt > 1 {
		return nil, errcodev1.ErrOutOfRange().WithMessage("out of range")
	}

	id, err := mysql.FileSpaceTopId(invoker.Db.WithContext(ctx), req.GetGuid(), spaceGuid)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileSpaceTopId fail, err: " + err.Error())
	}
	if id > 0 {
		return nil, errcodev1.ErrInternal().WithMessage("exist id")
	}
	err = invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = mysql.SetFileSpaceTop(tx, req.GetUid(), req.GetGuid(), spaceGuid)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("FileSpaceTopCreate fail, err: " + err.Error())
		}
		err = cache.FileCache.SetSiteTopField(ctx, tx, req.GetGuid(), 1)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("FileSpaceTopCreate fail2, err: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &articlev1.SetDocumentSpaceTopRes{}, nil
}

// CancelDocumentSpaceTop 取消文章置顶
func (GrpcServer) CancelDocumentSpaceTop(ctx context.Context, req *articlev1.CancelDocumentSpaceTopReq) (*articlev1.CancelDocumentSpaceTopRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_SET_SITE_TOP, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}
	// spaceGuid, err := mysql.FileSpaceGuidByGuid(invoker.Db.WithContext(ctx), req.GetCmtGuid(), req.GetGuid())
	// if err != nil {
	//	return nil, errcodev1.ErrDbError().WithMessage("FileSpaceGuidByGuid fail, err: " + err.Error())
	// }
	err := invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := mysql.CancelFileSpaceTop(tx, req.GetUid(), req.GetGuid())
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("FileSpaceTopDelete fail, err: " + err.Error())
		}
		err = cache.FileCache.SetSiteTopField(ctx, tx, req.GetGuid(), 0)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("FileSpaceTopDelete fail2, err: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &articlev1.CancelDocumentSpaceTopRes{}, nil
}

// DocumentSpaceTopList 文章置顶列表
func (GrpcServer) DocumentSpaceTopList(ctx context.Context, req *articlev1.DocumentSpaceTopListReq) (*articlev1.DocumentSpaceTopListRes, error) {
	fileGuids, err := mysql.FileSpaceTopGuids(invoker.Db.WithContext(ctx), req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileSpaceTopGuids fail, err: " + err.Error())
	}
	list, err := service.File.ListCacheByGuids(ctx, invoker.Db.WithContext(ctx), fileGuids)
	return &articlev1.DocumentSpaceTopListRes{
		List: list.ToArticlePb(commonv1.SPC_LAYOUT_ARTICLE_CARD),
	}, nil
}

// SetDocumentRecommend 将文章推荐
func (GrpcServer) SetDocumentRecommend(ctx context.Context, req *articlev1.SetDocumentRecommendReq) (*articlev1.SetDocumentRecommendRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_SET_RECOMMEND, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}

	id, err := mysql.FileRecommendId(invoker.Db.WithContext(ctx), req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileRecommendId fail, err: " + err.Error())
	}
	if id > 0 {
		return nil, errcodev1.ErrAlreadyExist().WithMessage("exist id")
	}
	tx := invoker.Db.WithContext(ctx).Begin()
	err = mysql.FileSuggest(tx, req.GetUid(), req.GetGuid())
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("file suggest fail, err: " + err.Error())
	}
	err = cache.FileCache.SetRecommend(ctx, tx, req.GetGuid(), 1)
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("SetDocumentRecommend fail2, err: " + err.Error())
	}
	tx.Commit()
	return &articlev1.SetDocumentRecommendRes{}, nil
}

// CancelDocumentRecommend 取消文章推荐
func (GrpcServer) CancelDocumentRecommend(ctx context.Context, req *articlev1.CancelDocumentRecommendReq) (*articlev1.CancelDocumentRecommendRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_SET_RECOMMEND, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}
	err := invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := mysql.FileCancelSuggest(tx, req.GetUid(), req.GetGuid()); err != nil {
			return errcodev1.ErrDbError().WithMessage("file suggest cancel fail, err: " + err.Error())
		}
		if err := cache.FileCache.SetRecommend(ctx, tx, req.GetGuid(), 0); err != nil {
			return errcodev1.ErrDbError().WithMessage("CancelDocumentRecommend fail2, err: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &articlev1.CancelDocumentRecommendRes{}, nil
}

// DocumentRecommendList 文章推荐列表
// 免费5个
// 收费10个
func (GrpcServer) DocumentRecommendList(ctx context.Context, req *articlev1.DocumentRecommendListReq) (*articlev1.DocumentRecommendListRes, error) {
	fileGuids, err := mysql.FileRecommendGuids(invoker.Db.WithContext(ctx), req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileRecommendGuids fail, err: " + err.Error())
	}
	list, err := service.File.ListCacheByGuids(ctx, invoker.Db.WithContext(ctx), fileGuids)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileListByGuids fail, err: " + err.Error())
	}
	return &articlev1.DocumentRecommendListRes{
		List: list.ToArticlePb(commonv1.SPC_LAYOUT_ARTICLE_LIST),
	}, nil
}

// CloseDocumentComment 关闭评论
func (GrpcServer) CloseDocumentComment(ctx context.Context, req *articlev1.CloseDocumentCommentReq) (*articlev1.CloseDocumentCommentRes, error) {
	// 权限检查
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_SET_COMMENT, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}
	// 更新db
	if err := mysql.FileCloseComment(invoker.Db.WithContext(ctx), req.GetUid(), req.GetGuid()); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("file suggest cancel fail, err: " + err.Error())
	}
	// 更新缓存
	if err := cache.FileCache.SetIsAllowCreateComment(ctx, invoker.Db.WithContext(ctx), req.GetGuid(), 0); err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("file suggest cancel fail2, err: " + err.Error())
	}
	return &articlev1.CloseDocumentCommentRes{}, nil
}

// OpenDocumentComment 打开评论
func (GrpcServer) OpenDocumentComment(ctx context.Context, req *articlev1.OpenDocumentCommentReq) (*articlev1.OpenDocumentCommentRes, error) {
	// 权限检查
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_SET_COMMENT, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}
	// 更新db
	if err := mysql.FileOpenComment(invoker.Db.WithContext(ctx), req.GetUid(), req.GetGuid()); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("OpenDocumentComment fail, err: " + err.Error())
	}
	// 更新缓存
	if err := cache.FileCache.SetIsAllowCreateComment(ctx, invoker.Db.WithContext(ctx), req.GetGuid(), 1); err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("file suggest cancel fail2, err: " + err.Error())
	}
	return &articlev1.OpenDocumentCommentRes{}, nil
}

// PublicListByUserCreated 用户公开创建的文章列表
func (GrpcServer) PublicListByUserCreated(ctx context.Context, req *articlev1.PublicListByUserCreatedReq) (*articlev1.PublicListByUserCreatedRes, error) {
	files, err := service.File.PublicUserArticleListPage(ctx, invoker.Db.WithContext(ctx), req.GetCreatedUid(), req.GetPagination())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("PublicDocumentListByUserCreated fail, err: " + err.Error())
	}

	return &articlev1.PublicListByUserCreatedRes{
		List:       files.ToArticlePb(commonv1.SPC_LAYOUT_ARTICLE_LIST),
		Pagination: req.GetPagination(),
	}, nil
}
