package file

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	filev1 "ecodepost/pb/file/v1"
)

func (GrpcServer) ListPage(ctx context.Context, req *filev1.ListPageReq) (*filev1.ListPageRes, error) {
	list, err := service.File.FileListPage(ctx, invoker.Db.WithContext(ctx), req.GetSpaceGuid(), req.GetPagination(), req.GetSort())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("ListPage fail, err: " + err.Error())
	}

	return &filev1.ListPageRes{
		List:       list.ToFileShowPb(),
		Pagination: req.GetPagination(),
	}, nil
}

// SpaceTopList 文章置顶列表
func (GrpcServer) SpaceTopList(ctx context.Context, req *filev1.SpaceTopListReq) (*filev1.SpaceTopListRes, error) {
	fileGuids, err := mysql.FileSpaceTopGuids(invoker.Db.WithContext(ctx), req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileSpaceTopGuids fail, err: " + err.Error())
	}
	list, err := service.File.ListCacheByGuids(ctx, invoker.Db.WithContext(ctx), fileGuids)
	return &filev1.SpaceTopListRes{
		List: list.ToFileShowPb(),
	}, nil
}

// RecommendList 文章推荐列表
// 免费5个
// 收费10个
func (GrpcServer) RecommendList(ctx context.Context, req *filev1.RecommendListReq) (*filev1.RecommendListRes, error) {
	fileGuids, err := mysql.FileRecommendGuids(invoker.Db.WithContext(ctx), req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileRecommendGuids fail, err: " + err.Error())
	}
	list, err := service.File.ListCacheByGuids(ctx, invoker.Db.WithContext(ctx), fileGuids)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileListByGuids fail, err: " + err.Error())
	}
	return &filev1.RecommendListRes{
		List: list.ToFileShowPb(),
	}, nil
}

// PermissionList 列表的permission
func (GrpcServer) PermissionList(ctx context.Context, req *filev1.PermissionListReq) (*filev1.PermissionListRes, error) {
	// 文件信息
	fileList, err := mysql.FileListByField(invoker.Db.WithContext(ctx), "guid,created_by,close_comment_time,space_guid", req.GetFileGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("IsFileAuthor fail, err: " + err.Error())
	}

	var (
		isAllowWrite         bool
		isAllowDelete        bool
		isAllowSetComment    bool
		isAllowRecommend     bool
		isAllowSiteTop       bool
		isAllowCreateComment bool
	)

	isAdmin, _ := pmspolicy.IsPlatformAdmin(req.GetUid())
	fmt.Printf("isAdmin--------------->"+"%+v\n", isAdmin)
	output := make([]*filev1.PermissionRes, 0)

	for _, value := range fileList {
		if isAdmin {
			output = append(output, &filev1.PermissionRes{
				Guid:                 value.Guid,
				IsAllowWrite:         true,
				IsAllowDelete:        true,
				IsAllowSiteTop:       true,
				IsAllowRecommend:     true,
				IsAllowSetComment:    true,
				IsAllowCreateComment: true,
			})
			continue
		}

		// 如果是作者
		// 可以写文档、可以删除文档，还可以设置评论是否关闭
		if value.CreatedBy == req.GetUid() {
			isAllowWrite = true
			isAllowDelete = true
			isAllowSetComment = true
		}
		// 空间允许评论，并且文章没有关闭评论，那么用户才可以评论
		isAllowCreateComment, err = pmspolicy.Check(ctx, commonv1.PMS_FILE_CREATE_COMMENT, req.GetUid(), value.Guid)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}

		isAllowWrite, err = pmspolicy.Check(ctx, commonv1.PMS_FILE_UPDATE, req.GetUid(), value.Guid)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}

		isAllowDelete, err = pmspolicy.Check(ctx, commonv1.PMS_FILE_DELETE, req.GetUid(), value.Guid)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}

		isAllowRecommend, err = pmspolicy.Check(ctx, commonv1.PMS_FILE_SET_RECOMMEND, req.GetUid(), value.Guid)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}

		isAllowSiteTop, err = pmspolicy.Check(ctx, commonv1.PMS_FILE_SET_SITE_TOP, req.GetUid(), value.Guid)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}

		output = append(output, &filev1.PermissionRes{
			Guid:                 value.Guid,
			IsAllowWrite:         isAllowWrite,
			IsAllowDelete:        isAllowDelete,
			IsAllowSiteTop:       isAllowSiteTop,
			IsAllowRecommend:     isAllowRecommend,
			IsAllowSetComment:    isAllowSetComment,
			IsAllowCreateComment: isAllowCreateComment,
		})
	}
	return &filev1.PermissionListRes{List: output}, nil
}

// ListPageByParent List Answer
func (GrpcServer) ListPageByParent(ctx context.Context, req *filev1.ListPageByParentReq) (*filev1.ListPageByParentRes, error) {
	list, err := service.File.FileListPageByParent(ctx, invoker.Db.WithContext(ctx), req.GetParentGuid(), req.GetPagination(), req.GetSort())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("ListQuestion fail, err: " + err.Error())
	}

	return &filev1.ListPageByParentRes{
		List:       list.ToFileShowPb(),
		Pagination: req.GetPagination(),
	}, nil
}
