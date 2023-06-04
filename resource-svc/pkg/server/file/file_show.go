package file

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	"gorm.io/gorm"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	filev1 "ecodepost/pb/file/v1"
)

func (GrpcServer) GetShowInfo(ctx context.Context, req *filev1.GetShowInfoReq) (*filev1.GetShowInfoRes, error) {
	fileInfo, err := mysql.GetFileGuidEerror(invoker.Db.WithContext(ctx), req.GetGuid())
	if err != nil {
		return nil, err
	}
	cacheInfo, err := service.File.Info(ctx, invoker.Db.WithContext(ctx), req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("document fail3, err: " + err.Error())
	}
	cacheInfo.CntCollect = fileInfo.CntCollect
	cacheInfo.CntComment = fileInfo.CntComment
	cacheInfo.CntView = fileInfo.CntView

	//aliSignUrl, err := invoker.AliOss.CdnAuthURL(cacheInfo.ContentKey, cacheInfo.Hash)
	//if err != nil {
	//	return nil, errcodev1.ErrInternal().WithMessage("document fail4, err: " + err.Error())
	//}

	go func() {
		invoker.Db.Model(mysql.File{}).Where("guid = ?", req.Guid).Update("cnt_view", gorm.Expr("`cnt_view` + 1"))
	}()

	return &filev1.GetShowInfoRes{
		File: cacheInfo.ToOneFileShowPb(),
	}, nil
}

// Permission 单个file的permission
func (GrpcServer) Permission(ctx context.Context, req *filev1.PermissionReq) (*filev1.PermissionRes, error) {
	resp := &filev1.PermissionRes{
		Guid: req.GetFileGuid(),
	}
	// 文件信息
	fileInfo, err := mysql.FileInfoByFieldMustExistsEerror(invoker.Db.WithContext(ctx), "id,created_by,close_comment_time,space_guid", req.GetFileGuid())
	if err != nil {
		return nil, err
	}

	// 如果是作者
	// 可以写文档、可以删除文档，还可以设置评论是否关闭
	if fileInfo.CreatedBy == req.GetUid() {
		resp.IsAllowWrite = true
		resp.IsAllowSetComment = true
	}

	// 空间允许评论，并且文章没有关闭评论，那么用户才可以评论
	isAllowCreateComment, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_CREATE_COMMENT, req.GetUid(), req.GetFileGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	resp.IsAllowCreateComment = isAllowCreateComment

	isAllowWrite, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_UPDATE, req.GetUid(), req.GetFileGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	resp.IsAllowWrite = isAllowWrite

	isArticleDelete, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_DELETE, req.GetUid(), req.GetFileGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	resp.IsAllowDelete = isArticleDelete

	isArticleSetRecommend, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_SET_RECOMMEND, req.GetUid(), req.GetFileGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	resp.IsAllowRecommend = isArticleSetRecommend

	isArticleSetSiteTop, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_SET_SITE_TOP, req.GetUid(), req.GetFileGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	resp.IsAllowSiteTop = isArticleSetSiteTop

	return resp, nil
}
