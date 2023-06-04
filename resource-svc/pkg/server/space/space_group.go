package space

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	commonv1 "ecodepost/pb/common/v1"
	loggerv1 "ecodepost/pb/logger/v1"

	communityv1 "ecodepost/pb/community/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	spacev1 "ecodepost/pb/space/v1"
	userv1 "ecodepost/pb/user/v1"
)

func (GrpcServer) CreateSpaceGroup(ctx context.Context, req *spacev1.CreateSpaceGroupReq) (resp *spacev1.CreateSpaceGroupRes, err error) {
	if err = pmspolicy.CheckEerror(ctx, commonv1.PMS_SPACE_GROUP_CREATE, req.GetOperateUid(), ""); err != nil {
		return nil, err
	}

	var data *spacev1.SpaceGroupInfo
	data, err = service.Space.CreateGroup(ctx, req)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("create group fail, err: " + err.Error())
	}
	return &spacev1.CreateSpaceGroupRes{
		Info: data,
	}, nil
}

// SpaceGroupInfo 空间分组基本信息
func (GrpcServer) SpaceGroupInfo(ctx context.Context, req *spacev1.SpaceGroupInfoReq) (*spacev1.SpaceGroupInfoRes, error) {
	info, err := mysql.GetSpaceGroupInfo(invoker.Db.WithContext(ctx), "guid,name,icon_type,icon,visibility,is_allow_read_member_list", req.GetSpaceGroupGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("space info, err: " + err.Error())
	}

	cnt, err := mysql.SpaceGroupMemberCnt(invoker.Db.WithContext(ctx), req.GetSpaceGroupGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("space info2, err: " + err.Error())
	}

	spaces, err := mysql.SpaceListByUserAndGroups(invoker.Db.WithContext(ctx), req.GetUid(), []string{req.GetSpaceGroupGuid()})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("space group info3, err: " + err.Error())
	}
	return &spacev1.SpaceGroupInfoRes{
		Guid:       info.Guid,
		Name:       info.Name,
		IconType:   info.IconType,
		Icon:       info.Icon,
		List:       spaces.ToPbWithCnt(ctx),
		Visibility: info.Visibility,
		MemberCnt:  cnt,
	}, nil
}

func (GrpcServer) DeleteSpaceGroup(ctx context.Context, req *spacev1.DeleteSpaceGroupReq) (*spacev1.DeleteSpaceGroupRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_SPACE_GROUP_DELETE, req.GetOperateUid(), ""); err != nil {
		return nil, err
	}
	if err := service.Space.DeleteGroup(invoker.Db.WithContext(ctx), req.GetOperateUid(), req.GetSpaceGroupGuid()); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("delete space group fail, err: " + err.Error())
	}
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:          commonv1.LOG_EVENT_SPACE_GROUP_DELETE,
		Group:          commonv1.LOG_GROUP_SPACE_GROUP,
		OperateUid:     req.GetOperateUid(),
		SpaceGroupGuid: req.GetSpaceGroupGuid(),
	})
	return &spacev1.DeleteSpaceGroupRes{}, nil
}

func (GrpcServer) UpdateSpaceGroup(ctx context.Context, req *spacev1.UpdateSpaceGroupReq) (*spacev1.UpdateSpaceGroupRes, error) {
	// TODO 部分更新
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_SPACE_GROUP_UPDATE, req.GetOperateUid(), ""); err != nil {
		return nil, err
	}

	err := service.Space.UpdateGroup(invoker.Db.WithContext(ctx), req.GetOperateUid(), req.GetSpaceGroupGuid(), map[string]any{
		"name":                      req.GetName(),
		"icon":                      req.GetIcon(),
		"icon_type":                 req.GetIconType().Number(),
		"visibility":                req.GetVisibility().Number(),
		"is_allow_read_member_list": req.GetIsAllowReadMemberList(),
	})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("delete space group fail, err: " + err.Error())
	}
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:          commonv1.LOG_EVENT_SPACE_GROUP_UPDATE,
		Group:          commonv1.LOG_GROUP_SPACE_GROUP,
		OperateUid:     req.GetOperateUid(),
		SpaceGroupGuid: req.GetSpaceGroupGuid(),
	})
	return &spacev1.UpdateSpaceGroupRes{}, nil
}

func (GrpcServer) ChangeSpaceGroupSort(ctx context.Context, req *spacev1.ChangeSpaceGroupSortReq) (*spacev1.ChangeSpaceGroupSortRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_SPACE_GROUP_UPDATE, req.GetOperateUid(), ""); err != nil {
		return nil, err
	}
	if err := service.Space.ChangeGroupSort(ctx, req.GetSpaceGroupGuid(), req.GetTargetSpaceGroupGuid(), req.GetDropPosition()); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("change space sort fail, err: " + err.Error())
	}
	invoker.GrpcLogger.Create(ctx, &loggerv1.CreateReq{
		Event:          commonv1.LOG_EVENT_SPACE_GROUP_CHANGE_SORT,
		Group:          commonv1.LOG_GROUP_SPACE_GROUP,
		OperateUid:     req.GetOperateUid(),
		SpaceGroupGuid: req.GetSpaceGroupGuid(),
	})
	return &spacev1.ChangeSpaceGroupSortRes{}, nil
}

func (GrpcServer) SpaceGroupMemberList(ctx context.Context, req *spacev1.SpaceGroupMemberListReq) (*spacev1.SpaceGroupMemberListRes, error) {
	list, err := mysql.SpaceGroupMemberList(invoker.Db.WithContext(ctx), req.SpaceGroupGuid, req.Uids, req.Pagination)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("SpaceGroupMemberList fail, err: " + err.Error())
	}
	uids := make([]int64, 0)
	for _, value := range list {
		uids = append(uids, value.Uid)
	}
	uMap, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{UidList: uids})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail, err: " + err.Error())
	}
	managerMembers, err := mysql.GetPmsSuperAdminMembers(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("list by role type, err: " + err.Error())
	}

	cmtInfo, err := invoker.GrpcCommunity.Info(ctx, &communityv1.InfoReq{})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetAllManagerMemberList fail2, err: " + err.Error())
	}

	return &spacev1.SpaceGroupMemberListRes{
		List:       list.ToMemberPb(uMap.GetUserMap(), cmtInfo.GetCommunity().GetUid(), managerMembers),
		Pagination: req.GetPagination(),
	}, nil
}
