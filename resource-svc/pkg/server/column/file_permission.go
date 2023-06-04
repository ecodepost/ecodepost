package column

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	columnv1 "ecodepost/pb/column/v1"

	errcodev1 "ecodepost/pb/errcode/v1"
)

// ListPermission 列表的permission
func (GrpcServer) ListPermission(ctx context.Context, req *columnv1.ListPermissionReq) (resp *columnv1.ListPermissionRes, err error) {
	var fileGuids []string
	db := invoker.Db.WithContext(ctx)
	err = db.Model(mysql.Edge{}).Where("space_guid = ?", req.GetSpaceGuid()).Pluck("file_guid", &fileGuids).Error
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("ListPermission fail, err: " + err.Error())
	}
	// 文件信息
	fileList, err := mysql.FileListByField(db, "guid,created_by,close_comment_time,space_guid", fileGuids)
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("ListPermission fail2, err: " + err.Error())
	}

	var (
		isAllowWrite         bool
		isAllowDelete        bool
		isAllowSetComment    bool
		isAllowCreateComment bool
	)

	output := make([]*columnv1.PermissionRes, 0)
	for _, value := range fileList {
		// 如果是作者
		// 可以写文档、可以删除文档，还可以设置评论是否关闭
		// if value.CreatedBy == req.GetUid() {
		isAllowWrite = true
		isAllowDelete = true
		isAllowSetComment = true
		isAllowCreateComment = true
		// }
		// 空间允许评论，并且文章没有关闭评论，那么用户才可以评论
		// isAllowCreateComment, err = pmspolicy.Check(ctx, commonv1.PMS_FILE_CREATE_COMMENT, req.GetCmtGuid(), req.GetUid(), value.Guid)
		// if err != nil {
		//	return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		// }
		// isAllowDelete, err = pmspolicy.Check(ctx, commonv1.PMS_FILE_DELETE, req.GetCmtGuid(), req.GetUid(), value.Guid)
		// if err != nil {
		//	return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		// }

		output = append(output, &columnv1.PermissionRes{
			Guid:                 value.Guid,
			IsAllowWrite:         isAllowWrite,
			IsAllowDelete:        isAllowDelete,
			IsAllowSetComment:    isAllowSetComment,
			IsAllowCreateComment: isAllowCreateComment,
		})
	}
	return &columnv1.ListPermissionRes{List: output}, nil
}
