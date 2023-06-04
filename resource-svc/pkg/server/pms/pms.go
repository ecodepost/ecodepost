package pms

import (
	"context"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service/pmsaction"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	pmsv1 "ecodepost/pb/pms/v1"
	userv1 "ecodepost/pb/user/v1"
)

type GrpcServer struct{}

var _ pmsv1.PmsServer = (*GrpcServer)(nil)

// Check 校验是否有权限
func (GrpcServer) Check(ctx context.Context, req *pmsv1.CheckReq) (*pmsv1.CheckRes, error) {
	if req.GetActionName() == commonv1.PMS_INVALID {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("action name empty")
	}
	flag, err := pmspolicy.Check(ctx, req.GetActionName(), req.GetUid(), req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrInvalidArgument()
	}
	return &pmsv1.CheckRes{
		Flag: flag,
	}, nil
}

func (GrpcServer) GetRoleMemberList(ctx context.Context, req *pmsv1.GetRoleMemberListReq) (*pmsv1.GetRoleMemberListRes, error) {
	list, err := mysql.GetRoleMembers(invoker.Db.WithContext(ctx), req.GetRoleId())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetRoleMember fail, err: " + err.Error())
	}
	// 如果为0，那么直接给出一个slice
	if len(list) == 0 {
		return &pmsv1.GetRoleMemberListRes{
			List: make([]*pmsv1.MemberInfo, 0),
		}, nil
	}
	uids := make([]int64, 0)
	for _, value := range list {
		uids = append(uids, value.Uid)
	}
	uMap, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{
		UidList: uids,
	})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail, err: " + err.Error())
	}

	return &pmsv1.GetRoleMemberListRes{
		List: list.ToPb(uMap.GetUserMap()),
	}, nil
}

// GetRolePermission 获取角色的权限
func (GrpcServer) GetRolePermission(ctx context.Context, req *pmsv1.GetRolePermissionReq) (*pmsv1.GetRolePermissionRes, error) {
	rolePolicies, spaceGroupPolicies, spacePolicies, err := mysql.PolicyPmsList(invoker.Db.WithContext(ctx), req.GetRoleId())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetRolePermission fail2, err: " + err.Error())
	}

	actionList, err := pmsaction.ListActionByRoleId(rolePolicies)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetRolePermission fail1, err: " + err.Error())
	}

	groupList, err := pmsaction.ListSpaceGroupActionByRoleId(invoker.Db.WithContext(ctx), req.GetRoleId(), spaceGroupPolicies)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("fail3, err: " + err.Error())
	}

	spaceList, err := pmsaction.ListSpaceActionByRoleId(invoker.Db.WithContext(ctx), req.GetRoleId(), spacePolicies)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("fail4, err: " + err.Error())
	}

	return &pmsv1.GetRolePermissionRes{
		List:           actionList,
		SpaceList:      spaceList,
		SpaceGroupList: groupList,
	}, nil
}

// GetInitActionOptionPermission 获取某种类型的初始权限列表
func (GrpcServer) GetInitActionOptionPermission(ctx context.Context, req *pmsv1.GetInitActionOptionPermissionReq) (resp *pmsv1.GetInitActionOptionPermissionRes, err error) {

	ctx, _ = context.WithTimeout(ctx, time.Second)

	switch req.GetType() {
	case commonv1.CMN_GUID_SPACE_GROUP:
		spaceGroupActions, err := pmsaction.ListActionOptionByType(commonv1.PMS_ACT_TYPE_SPACE_GROUP)
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("fail1 " + err.Error())
		}
		resp = &pmsv1.GetInitActionOptionPermissionRes{
			List: spaceGroupActions.ToPb(map[string]mysql.PmsPolicy{}),
		}
		return resp, nil
	case commonv1.CMN_GUID_SPACE:
		info, err := mysql.SpaceGetInfoByGuid(invoker.Db.WithContext(ctx), "type", req.GetGuid())
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("fail3 " + err.Error())
		}
		spaceActions, err := pmsaction.ListActionOptionBySpaceType(info.Type)
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("fail4 " + err.Error())
		}
		resp = &pmsv1.GetInitActionOptionPermissionRes{
			List: spaceActions.ToPb(map[string]mysql.PmsPolicy{}),
		}
		return resp, nil
	}
	return nil, errcodev1.ErrInvalidArgument().WithMessage("not exist guid type")

}

// TotalInfo 统计信息
func (GrpcServer) TotalInfo(ctx context.Context, req *pmsv1.TotalInfoReq) (*pmsv1.TotalInfoRes, error) {
	superAdminCnt, err := mysql.SuperAdminMemberCnt(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("community info fail1, err: " + err.Error())
	}
	roleMemberCnt, err := mysql.RoleMemberCnt(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("community info fail1, err: " + err.Error())
	}
	return &pmsv1.TotalInfoRes{
		SuperAdminCnt: superAdminCnt,
		RoleMemberCnt: roleMemberCnt,
	}, nil
}
