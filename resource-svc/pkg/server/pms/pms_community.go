package pms

import (
	"context"

	userv1 "ecodepost/pb/user/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	pmsv1 "ecodepost/pb/pms/v1"
)

func (GrpcServer) CommunityPermission(ctx context.Context, req *pmsv1.CommunityPermissionReq) (*pmsv1.CommunityPermissionRes, error) {
	resp := &pmsv1.CommunityPermissionRes{}
	isAllowSpaceGroupCreate, err := pmspolicy.Check(ctx, commonv1.PMS_SPACE_GROUP_CREATE, req.GetUid(), "")
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("CommunityPermission fail, err: " + err.Error())
	}
	resp.IsAllowCreateSpaceGroup = isAllowSpaceGroupCreate

	isAllowSpaceCreate, err := pmspolicy.Check(ctx, commonv1.PMS_SPACE_CREATE, req.GetUid(), "")
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("CommunityPermission fail, err: " + err.Error())
	}
	resp.IsAllowCreateSpace = isAllowSpaceCreate
	isSuper, err := pmspolicy.IsCommunitySuperAdmin(ctx, req.GetUid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("CommunityPermission fail, err: " + err.Error())
	}
	if isSuper {
		resp.IsAllowManageCommunity = true
		resp.IsAllowUpgradeEdition = true
	}

	return resp, nil
}

// GetManagerMemberList 获取管理员的成员列表
func (GrpcServer) GetManagerMemberList(ctx context.Context, req *pmsv1.GetManagerMemberListReq) (*pmsv1.GetManagerMemberListRes, error) {
	superAdminMembers, err := mysql.GetPmsSuperAdminMembers(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetManagerMemberList fail, err: " + err.Error())
	}
	if len(superAdminMembers) == 0 {
		return nil, errcodev1.ErrNotFound().WithMessage("GetManagerMemberList not found")

	}

	uids := make([]int64, 0)
	for _, value := range superAdminMembers {
		uids = append(uids, value.Uid)
	}
	uMap, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{
		UidList: uids,
	})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail, err: " + err.Error())
	}

	return &pmsv1.GetManagerMemberListRes{
		List: superAdminMembers.ToPb(uMap.GetUserMap()),
	}, nil
}
