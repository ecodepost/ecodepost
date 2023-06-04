package pms

import (
	"context"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"github.com/spf13/cast"

	commonv1 "ecodepost/pb/common/v1"
	// editionv1 "ecodepost/pb/edition/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	loggerv1 "ecodepost/pb/logger/v1"
	pmsv1 "ecodepost/pb/pms/v1"
)

// CreateManagerMember 批量添加管理员的成员
func (GrpcServer) CreateManagerMember(ctx context.Context, req *pmsv1.CreateManagerMemberReq) (*pmsv1.CreateManagerMemberRes, error) {
	operateInfo, err := mysql.GetSuperAdminMemberId(invoker.Db.WithContext(ctx), req.GetOperateUid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateManagerMember fail1, err: " + err.Error())
	}
	// 说明没有权限
	if operateInfo.Id == 0 {
		return nil, errcodev1.ErrPmsPermissionDenied().WithMessage("没有修改权限，该用户不是超级管理员, err: " + cast.ToString(req.GetOperateUid()))

	}

	nowTime := time.Now().Unix()
	output := make([]mysql.PmsSuperAdminMember, 0)
	for _, uid := range req.GetUids() {
		roleMemberInfo, err := mysql.GetSuperAdminMemberId(invoker.Db.WithContext(ctx), uid)
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("CollectionCreate fail2, err: " + err.Error())
		}
		if roleMemberInfo.Id > 0 {
			continue
		}
		output = append(output, mysql.PmsSuperAdminMember{
			Uid:        uid,
			CreatedUid: req.GetOperateUid(),
			Ctime:      nowTime,
		})
	}

	err = mysql.CreateSuperAdminMemberInBatches(invoker.Db.WithContext(ctx), output)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateManagerMember Fail, err: " + err.Error())
	}

	// 记录日志
	invoker.GrpcLogger.BatchCreateByTargetUids(ctx, &loggerv1.BatchCreateByTargetUidsReq{
		Event:      commonv1.LOG_EVENT_PMS_SUPER_ADMIN_CREATE_MEMBER,
		Group:      commonv1.LOG_GROUP_PMS,
		TargetUids: req.GetUids(),
		OperateUid: req.GetOperateUid(),
	})
	return &pmsv1.CreateManagerMemberRes{}, nil
}

// DeleteManagerMember 删除管理员的成员
func (GrpcServer) DeleteManagerMember(ctx context.Context, req *pmsv1.DeleteManagerMemberReq) (*pmsv1.DeleteManagerMemberRes, error) {
	err := mysql.DeleteManagerMembers(invoker.Db.WithContext(ctx), req)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateManagerMember Fail, err: " + err.Error())
	}

	// 记录日志
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:      commonv1.LOG_EVENT_PMS_SUPER_ADMIN_DELETE_MEMBER,
		Group:      commonv1.LOG_GROUP_PMS,
		TargetUid:  req.GetUid(),
		OperateUid: req.GetOperateUid(),
	})
	return &pmsv1.DeleteManagerMemberRes{}, nil
}
