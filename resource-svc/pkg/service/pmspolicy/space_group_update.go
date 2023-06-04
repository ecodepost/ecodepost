package pmspolicy

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
)

type SpaceGroupUpdate struct {
}

func init() {
	Register(NewActPolicy(&SpaceGroupUpdate{}, superAdminCheckFn))
}

func (*SpaceGroupUpdate) Scheme() commonv1.PMS_ACT {
	return commonv1.PMS_SPACE_GROUP_UPDATE
}

// Check 检查
// 1 超级管理员
// 2 用户组中是否设置里权限
func (s *SpaceGroupUpdate) Check(ctx context.Context, uid int64, guid string) (bool, error) {
	// 没有打开，查看下用户有哪些角色
	roleIds, err := mysql.RoleIdsByUid(invoker.Db.WithContext(ctx), uid)
	if err != nil {
		return false, nil
	}

	if len(roleIds) == 0 {
		return false, nil
	}
	return mysql.CheckRolePolicy(invoker.Db.WithContext(ctx), roleIds, s.Scheme())
}
