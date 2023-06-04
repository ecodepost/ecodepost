package pmspolicy

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
)

type FileCreate struct{}

func init() {
	Register(NewActPolicy(&FileCreate{}, superAdminCheckFn, platformAdminCheckFn))
}

func (*FileCreate) Scheme() commonv1.PMS_ACT {
	return commonv1.PMS_FILE_CREATE
}

// Check 检查
// 1 超级管理员
// 2 平台管理员
// 3 空间设置里是否可以创建文件
// 4 用户组中是否设置里权限
// 5 是不是空间成员
func (s *FileCreate) Check(ctx context.Context, uid int64, spaceGuid string) (bool, error) {
	// todo 暂时关闭，用于跑脚本
	info, err := mysql.SpaceMemberInfo(invoker.Db.WithContext(ctx), "id", spaceGuid, uid)
	if err != nil {
		return false, fmt.Errorf("SpaceCreateArticle Check fail, err: %w", err)
	}
	if info.Id == 0 {
		return false, NotSpaceMemberError
	}
	spaceOption, err := mysql.GetSpaceOptionInfo(invoker.Db.WithContext(ctx), spaceGuid, commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_FILE)
	if err != nil {
		return false, fmt.Errorf("SpaceCreateArticle Check fail, err: %w", err)
	}

	if spaceOption.OptionValue > 0 {
		return true, nil
	}
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
