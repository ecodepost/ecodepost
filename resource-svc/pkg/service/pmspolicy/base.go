package pmspolicy

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"

	"github.com/gotomicro/ego/core/econf"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

var registry = make(map[commonv1.PMS_ACT]*ActPolicy)

func Register(p *ActPolicy) {
	registry[p.mainPolicyAct] = p
}

// Policy 权限策略
type Policy interface {
	Scheme() commonv1.PMS_ACT
	Check(ctx context.Context, uid int64, guid string) (pass bool, err error)
}

// ActPolicy 执行Act动作的具体权限判断策略
type ActPolicy struct {
	// beforeChecks, 前置policy, 执行主policy前, 会执行beforeChecks, 任意beforeCheck无权限，则判定无权限
	beforeChecks []BeforeCheck
	// mainPolicy 执行beforeChecks后, 再执行主 policy, 主policy返回最终权限判定结果
	mainPolicy Policy
	// mainPolicyAct 主policy关联的PMS_ACT名称
	mainPolicyAct commonv1.PMS_ACT
}

// BeforeCheck 前置pms检查函数, 接收uid、guid等常规入参
// 返回 pass: 是否权限检查通过；ignoreNext:是否可忽略下个 BeforeCheck；err:是否执行报错等参数
type BeforeCheck = func(ctx context.Context, uid int64, guid string) (pass bool, ignoreNext bool, err error)

// NewActPolicy 构造ActPolicy，后续再改造参数为option模式
func NewActPolicy(mainPolicy Policy, beforePolicies ...BeforeCheck) *ActPolicy {
	return &ActPolicy{beforeChecks: beforePolicies, mainPolicy: mainPolicy, mainPolicyAct: mainPolicy.Scheme()}
}

// Check 进行权限检查
func (b *ActPolicy) Check(ctx context.Context, uid int64, guid string) (flag bool, err error) {
	// 执行提前检查policy
	for _, v := range b.beforeChecks {
		flag, ignoreNext, err := v(ctx, uid, guid)
		// 如果执行失败则直接判定无权限，并抛出错误
		if err != nil {
			return false, err
		}
		// 如果可以忽略后续beforeChecks和policy，则直接返回当前权限判定结果; 否则继续执行
		if ignoreNext {
			return flag, nil
		}
	}
	// 执行主检查policy
	return b.mainPolicy.Check(ctx, uid, guid)
}

// superAdminCheckFn 是否是超级管理员，如果是超级管理员，则直接pass
var superAdminCheckFn BeforeCheck = func(ctx context.Context, uid int64, guid string) (bool, bool, error) {
	flag, err := IsCommunitySuperAdmin(ctx, uid)
	if err != nil {
		return false, false, fmt.Errorf("get IsExistRoleMember fail, err: %w", err)
	}
	// 只有Check成功，才能忽略后续检查
	if flag {
		return true, true, nil
	}
	return false, false, nil
}

// platformAdminCheckFn 是否是平台管理员，如果是平台管理员，则直接pass
var platformAdminCheckFn BeforeCheck = func(ctx context.Context, uid int64, guid string) (bool, bool, error) {
	flag, err := IsPlatformAdmin(uid)
	if err != nil {
		return false, false, fmt.Errorf("get IsExistRoleMember fail, err: %w", err)
	}
	// 只有Check成功，才能忽略后续检查
	if flag {
		return true, true, nil
	}
	return false, false, nil
}

// Check 进行权限检查
func Check(ctx context.Context, scheme commonv1.PMS_ACT, uid int64, guid string) (bool, error) {
	policy, exist := registry[scheme]
	if !exist {
		return false, fmt.Errorf("not exist scheme: " + scheme.String())
	}

	return policy.Check(ctx, uid, guid)
}

// CheckEerror 进行权限检查, 并直接返回errors.Error, 供外部grpc使用
func CheckEerror(ctx context.Context, scheme commonv1.PMS_ACT, uid int64, guid string) error {
	isPass, err := Check(ctx, scheme, uid, guid)
	if err != nil {
		return errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	if !isPass {
		return errcodev1.ErrAlreadyPermissionDenied().WithMessage("no auth")
	}
	return nil
}

// IsCommunitySuperAdmin 查询该用户是不是超级管理员
func IsCommunitySuperAdmin(ctx context.Context, uid int64) (flag bool, err error) {
	info := mysql.PmsSuperAdminMember{}
	err = invoker.Db.WithContext(ctx).Select("id").Where("uid = ?", uid).Find(&info).Error
	if err != nil {
		err = fmt.Errorf("IsCommunitySuperAdmin fail, err: %w", err)
		return
	}
	if info.Id > 0 {
		flag = true
		return
	}
	return
}

// IsPlatformAdmin 查询用户是不是平台管理员
func IsPlatformAdmin(uid int64) (flag bool, err error) {
	adminUids := cast.ToIntSlice(econf.Get("adminUids"))
	isAdmin := false
	if len(adminUids) > 0 {
		isAdmin = lo.Contains(adminUids, int(uid))
	}
	return isAdmin, nil
}
