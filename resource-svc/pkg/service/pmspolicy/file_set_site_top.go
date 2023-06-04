package pmspolicy

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
)

type FileSetSiteTop struct {
}

func init() {
	Register(NewActPolicy(&FileSetSiteTop{}, superAdminCheckFn))
}

func (*FileSetSiteTop) Scheme() commonv1.PMS_ACT {
	return commonv1.PMS_FILE_SET_SITE_TOP
}

// Check 检查
// 1 超级管理员
// 2 用户组中是否设置里权限
func (s *FileSetSiteTop) Check(ctx context.Context, uid int64, fileGuid string) (bool, error) {
	// 文件信息
	fileInfo, err := mysql.FileInfoByFieldMustExistsEerror(invoker.Db.WithContext(ctx), "id,space_guid", fileGuid)
	if err != nil {
		return false, err
	}
	// 没有打开，查看下用户有哪些角色
	roleIds, err := mysql.RoleIdsByUid(invoker.Db.WithContext(ctx), uid)
	if err != nil {
		return false, fmt.Errorf("FileSetSiteTop get roleIds fail, err: %w", err)
	}
	if len(roleIds) == 0 {
		return false, nil
	}
	return mysql.CheckPolicy(invoker.Db.WithContext(ctx), roleIds, s.Scheme(), fileInfo.SpaceGuid, commonv1.PMS_RSC_SPACE)
}
