package pmspolicy

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
)

type FileManage struct {
}

func init() {
	Register(NewActPolicy(&FileManage{}, superAdminCheckFn))
}

func (*FileManage) Scheme() commonv1.PMS_ACT {
	return commonv1.PMS_FILE_MANAGE
}

// Check 检查
// 1 超级管理员
// 2 空间设置里是否可以创建文件
// 3 用户组中是否设置里权限
func (s *FileManage) Check(ctx context.Context, uid int64, spaceGuid string) (bool, error) {
	// 如果没有打开，查看下用户有哪些角色
	// 没有打开，查看下用户有哪些角色
	roleIds, err := mysql.RoleIdsByUid(invoker.Db.WithContext(ctx), uid)
	if err != nil {
		return false, fmt.Errorf("SpaceArticleCreate get roleIds fail, err: %w", err)
	}

	if len(roleIds) == 0 {
		return false, nil
	}
	return mysql.CheckPolicy(invoker.Db.WithContext(ctx), roleIds, s.Scheme(), spaceGuid, commonv1.PMS_RSC_SPACE)
}
