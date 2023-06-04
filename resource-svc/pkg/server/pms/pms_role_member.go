package pms

import (
	"context"
	"encoding/json"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	loggerv1 "ecodepost/pb/logger/v1"
	pmsv1 "ecodepost/pb/pms/v1"
)

func (GrpcServer) CreateRoleMember(ctx context.Context, req *pmsv1.CreateRoleMemberReq) (*pmsv1.CreateRoleMemberRes, error) {
	nowTime := time.Now().Unix()
	output := make([]mysql.PmsRoleMember, 0)
	for _, uid := range req.GetUids() {
		roleMemberInfo, err := mysql.GetRoleMemberId(invoker.Db.WithContext(ctx), uid, req.GetRoleId())
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("CollectionCreate fail2, err: " + err.Error())
		}
		if roleMemberInfo.Id > 0 {
			continue
		}
		output = append(output, mysql.PmsRoleMember{
			RoleId:    req.GetRoleId(),
			Uid:       uid,
			CreatedBy: req.GetOperateUid(),
			Ctime:     nowTime,
		})
	}

	err := mysql.CreateRoleMemberInBatches(invoker.Db.WithContext(ctx), output)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateRoleMember Fail, err: " + err.Error())
	}

	metadata := make(map[string]any)
	metadata["role_id"] = req.GetRoleId()
	metadataBytes, _ := json.Marshal(metadata)
	// 记录日志
	invoker.GrpcLogger.BatchCreateByTargetUids(ctx, &loggerv1.BatchCreateByTargetUidsReq{
		Event:      commonv1.LOG_EVENT_PMS_ROLE_CREATE,
		Group:      commonv1.LOG_GROUP_PMS,
		TargetUids: req.GetUids(),
		OperateUid: req.GetOperateUid(),
		Metadata:   string(metadataBytes),
	})
	return &pmsv1.CreateRoleMemberRes{}, nil
}

func (GrpcServer) DeleteRoleMember(ctx context.Context, req *pmsv1.DeleteRoleMemberReq) (*pmsv1.DeleteRoleMemberRes, error) {
	tx := invoker.Db.WithContext(ctx).Begin()
	for _, uid := range req.GetUids() {
		err := mysql.DeleteRoleMember(tx, req.GetRoleId(), uid)
		if err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("DeleteRoleMember Fail, err: " + err.Error())
		}
	}
	tx.Commit()
	metadata := make(map[string]any)
	metadata["role_id"] = req.GetRoleId()
	metadataBytes, _ := json.Marshal(metadata)
	// 记录日志
	invoker.GrpcLogger.BatchCreateByTargetUids(ctx, &loggerv1.BatchCreateByTargetUidsReq{
		Event:      commonv1.LOG_EVENT_PMS_ROLE_DELETE,
		Group:      commonv1.LOG_GROUP_PMS,
		TargetUids: req.GetUids(),
		OperateUid: req.GetOperateUid(),
		Metadata:   string(metadataBytes),
	})
	return &pmsv1.DeleteRoleMemberRes{}, nil
}
