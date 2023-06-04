package pmspolicy

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	errcodev1 "ecodepost/pb/errcode/v1"

	commonv1 "ecodepost/pb/common/v1"
)

type FileDeleteComment struct {
}

func init() {
	Register(NewActPolicy(&FileDeleteComment{}, superAdminCheckFn))
}

func (*FileDeleteComment) Scheme() commonv1.PMS_ACT {
	return commonv1.PMS_FILE_DELETE_COMMENT
}

// Check 检查
// 1 超级管理员
// 2 是不是该评论的作者
// 3 用户组中是否设置里权限
func (s *FileDeleteComment) Check(ctx context.Context, uid int64, fileGuid string) (bool, error) {
	// 文件信息
	commentInfo, err := mysql.CommentInfoByField(invoker.Db.WithContext(ctx), "id,uid,biz_guid", fileGuid)
	if err != nil {
		return false, errcodev1.ErrInternal().WithMessage("FileDeleteComment fail, err: " + err.Error())
	}
	if commentInfo.Id == 0 {
		return false, errcodev1.ErrDbError().WithMessage("FileDeleteComment comment info not exist")
	}
	if commentInfo.Uid == uid {
		return true, nil
	}
	// 文件信息
	fileInfo, err := mysql.FileInfoByFieldMustExistsEerror(invoker.Db.WithContext(ctx), "id,created_by,space_guid", fileGuid)
	if err != nil {
		return false, err
	}

	roleIds, err := mysql.RoleIdsByUid(invoker.Db.WithContext(ctx), uid)
	if err != nil {
		return false, fmt.Errorf("FileDeleteComment get roleIds fail, err: %w", err)
	}

	if len(roleIds) == 0 {
		return false, nil
	}
	return mysql.CheckPolicy(invoker.Db.WithContext(ctx), roleIds, s.Scheme(), fileInfo.SpaceGuid, commonv1.PMS_RSC_SPACE)
}
