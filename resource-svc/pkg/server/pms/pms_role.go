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

// CreateRole 创建角色
func (GrpcServer) CreateRole(ctx context.Context, req *pmsv1.CreateRoleReq) (*pmsv1.CreateRoleRes, error) {
	nowTime := time.Now().Unix()
	data := &mysql.PmsRole{
		Name:  req.GetName(),
		Ctime: nowTime,
		Utime: nowTime,
	}
	err := mysql.CreatePmsRole(invoker.Db.WithContext(ctx), data)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateRole fail, err: " + err.Error())
	}

	// 记录日志
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:      commonv1.LOG_EVENT_PMS_ROLE_CREATE,
		Group:      commonv1.LOG_GROUP_PMS,
		RoleId:     data.Id,
		OperateUid: req.GetOperateUid(),
	})

	return &pmsv1.CreateRoleRes{
		RoleId: data.Id,
		Name:   req.GetName(),
	}, nil
}

// UpdateRole 修改角色
func (GrpcServer) UpdateRole(ctx context.Context, req *pmsv1.UpdateRoleReq) (*pmsv1.UpdateRoleRes, error) {
	roleInfo, err := mysql.GetPmsRole(invoker.Db.WithContext(ctx), req.GetRoleId())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdateRole fail1, err: " + err.Error())
	}
	if roleInfo.Id == 0 {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("not exit role id")
	}

	err = mysql.UpdatePmsRole(invoker.Db.WithContext(ctx), req.GetRoleId(), req.GetName())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdateRole fail2, err: " + err.Error())
	}

	metadata := make(map[string]any)
	metadata["old_name"] = roleInfo.Name
	metadata["new_name"] = req.GetName()
	metadataByte, _ := json.Marshal(metadata)
	// 为了国际化没办法。只能记录这些数据
	// metadata["msg"] = "修改角色名称，由" + roleInfo.Name + "改为" + req.GetName()
	// 记录日志
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:      commonv1.LOG_EVENT_PMS_ROLE_UPDATE,
		Group:      commonv1.LOG_GROUP_PMS,
		RoleId:     req.GetRoleId(),
		OperateUid: req.GetOperateUid(),
		Metadata:   string(metadataByte),
	})
	return &pmsv1.UpdateRoleRes{}, nil
}

// DeleteRole 删除角色
func (GrpcServer) DeleteRole(ctx context.Context, req *pmsv1.DeleteRoleReq) (*pmsv1.DeleteRoleRes, error) {
	tx := invoker.Db.WithContext(ctx).Begin()
	err := mysql.DeletePmsRole(tx, req.GetRoleId())
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("delete role fail, err: " + err.Error())
	}
	tx.Commit()

	// 记录日志
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:      commonv1.LOG_EVENT_PMS_ROLE_DELETE,
		Group:      commonv1.LOG_GROUP_PMS,
		RoleId:     req.GetRoleId(),
		OperateUid: req.GetOperateUid(),
	})
	return &pmsv1.DeleteRoleRes{}, nil
}

func (GrpcServer) GetRoleList(ctx context.Context, req *pmsv1.GetRoleListReq) (*pmsv1.GetRoleListRes, error) {
	roles, err := mysql.GetPmsRoles(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetRoleList fail, err: " + err.Error())
	}

	return &pmsv1.GetRoleListRes{List: roles.ToPb()}, nil
}

// GetRoleIds 获取某个用户的role ids
func (GrpcServer) GetRoleIds(ctx context.Context, req *pmsv1.GetRoleIdsReq) (*pmsv1.GetRoleIdsRes, error) {
	roleIds, err := mysql.RoleIdsByUid(invoker.Db.WithContext(ctx), req.GetUid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetRoleList fail, err: " + err.Error())
	}
	return &pmsv1.GetRoleIdsRes{
		RoleIds: roleIds,
	}, nil
}

func (GrpcServer) PutRolePermission(ctx context.Context, req *pmsv1.PutRolePermissionReq) (*pmsv1.PutRolePermissionRes, error) {
	tx := invoker.Db.WithContext(ctx).Begin()
	oldPolicy, newPolicy, err := mysql.PutRolePolicy(tx, req)
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("PutRolePermission fail, err: " + err.Error())
	}
	tx.Commit()

	metadata := make(map[string]any)
	metadata["old_policy"] = oldPolicy
	metadata["new_policy"] = newPolicy
	metadataByte, _ := json.Marshal(metadata)
	// 为了国际化没办法。只能记录这些数据
	// metadata["msg"] = "修改角色名称，由" + roleInfo.Name + "改为" + req.GetName()
	// 记录日志
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:      commonv1.LOG_EVENT_PMS_ROLE_UPDATE,
		Group:      commonv1.LOG_GROUP_PMS,
		RoleId:     req.GetRoleId(),
		OperateUid: req.GetOperateUid(),
		Metadata:   string(metadataByte),
	})

	return &pmsv1.PutRolePermissionRes{}, nil
}

func (GrpcServer) PutRoleSpaceGroupPermission(ctx context.Context, req *pmsv1.PutRoleSpaceGroupPermissionReq) (*pmsv1.PutRoleSpaceGroupPermissionRes, error) {
	tx := invoker.Db.WithContext(ctx).Begin()
	oldPolicy, newPolicy, err := mysql.PutRoleSpaceGroupPolicy(tx, req)
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("PutRoleSpaceGroupPermission fail, err: " + err.Error())
	}
	tx.Commit()

	metadata := make(map[string]any)
	metadata["old_policy"] = oldPolicy
	metadata["new_policy"] = newPolicy
	metadataByte, _ := json.Marshal(metadata)
	// 为了国际化没办法。只能记录这些数据
	// metadata["msg"] = "修改角色名称，由" + roleInfo.Name + "改为" + req.GetName()
	// 记录日志
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:      commonv1.LOG_EVENT_PMS_ROLE_UPDATE,
		Group:      commonv1.LOG_GROUP_PMS,
		RoleId:     req.GetRoleId(),
		OperateUid: req.GetOperateUid(),
		Metadata:   string(metadataByte),
	})
	return &pmsv1.PutRoleSpaceGroupPermissionRes{}, nil
}

func (GrpcServer) PutRoleSpacePermission(ctx context.Context, req *pmsv1.PutRoleSpacePermissionReq) (*pmsv1.PutRoleSpacePermissionRes, error) {
	tx := invoker.Db.WithContext(ctx).Begin()
	oldPolicy, newPolicy, err := mysql.PutRoleSpacePolicy(tx, req)
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("PutRoleSpacePermission fail, err: " + err.Error())
	}
	tx.Commit()

	metadata := make(map[string]any)
	metadata["old_policy"] = oldPolicy
	metadata["new_policy"] = newPolicy
	metadataByte, _ := json.Marshal(metadata)
	// 为了国际化没办法。只能记录这些数据
	// metadata["msg"] = "修改角色名称，由" + roleInfo.Name + "改为" + req.GetName()
	// 记录日志
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:      commonv1.LOG_EVENT_PMS_ROLE_UPDATE,
		Group:      commonv1.LOG_GROUP_PMS,
		RoleId:     req.GetRoleId(),
		OperateUid: req.GetOperateUid(),
		Metadata:   string(metadataByte),
	})
	return &pmsv1.PutRoleSpacePermissionRes{}, nil
}
