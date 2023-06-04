package pmspolicy

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
)

type FileUpdate struct {
}

func init() {
	Register(NewActPolicy(&FileUpdate{}, superAdminCheckFn, platformAdminCheckFn))

}

func (*FileUpdate) Scheme() commonv1.PMS_ACT {
	return commonv1.PMS_FILE_UPDATE
}

// Check 检查
// 1 是不是平台管理员
// 2 是不是该文章的作者
func (s *FileUpdate) Check(ctx context.Context, uid int64, fileGuid string) (bool, error) {
	// 文件信息
	fileInfo, err := mysql.FileInfoByFieldMustExistsEerror(invoker.Db.WithContext(ctx), "id,created_by,space_guid", fileGuid)
	if err != nil {
		return false, err
	}
	if fileInfo.CreatedBy == uid {
		return true, nil
	}
	// 没有打开，查看下用户有哪些角色
	roleIds, err := mysql.RoleIdsByUid(invoker.Db.WithContext(ctx), uid)
	if err != nil {
		return false, fmt.Errorf("SpaceArticleDelete get roleIds fail, err: %w", err)
	}

	if len(roleIds) == 0 {
		return false, nil
	}
	return mysql.CheckPolicy(invoker.Db.WithContext(ctx), roleIds, s.Scheme(), fileInfo.SpaceGuid, commonv1.PMS_RSC_SPACE)
}
