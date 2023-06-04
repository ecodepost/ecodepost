package space

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"
	"ecodepost/resource-svc/pkg/service/spaceoption"

	commonv1 "ecodepost/pb/common/v1"
	communityv1 "ecodepost/pb/community/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	spacev1 "ecodepost/pb/space/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/samber/lo"
)

type GrpcServer struct {
}

var _ spacev1.SpaceServer = (*GrpcServer)(nil)

func (GrpcServer) ListSpaceAndGroup(ctx context.Context, req *spacev1.ListSpaceAndGroupReq) (*spacev1.ListSpaceAndGroupRes, error) {
	spaceList, groupList, err := service.Space.List(ctx, invoker.Db.WithContext(ctx), req.OperateUid)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("list space tree fail, err: " + err.Error())
	}
	return &spacev1.ListSpaceAndGroupRes{
		SpaceGroupList: groupList,
		SpaceList:      spaceList,
	}, nil
}

// ListTree 空间的树形结构
//func (GrpcServer) ListTree(ctx context.Context, req *spacev1.ListTreeReq) (resp *spacev1.ListTreeRes, err error) {
//	list, err := service.Space.Tree(ctx, invoker.Db.WithContext(ctx), req.GetCmtGuid(), req.GetOperateUid())
//	if err != nil {
//		return nil, errcodev1.ErrDbError().WithMessage("list space tree fail, err: " + err.Error())
//	}
//	return &spacev1.ListTreeRes{Tree: list}, nil
//}

// ListPublicSpace 获取社区公开的space，社区跳转默认首页只允许公开的space
func (GrpcServer) ListPublicSpace(ctx context.Context, req *spacev1.ListPublicSpaceReq) (*spacev1.ListPublicSpaceRes, error) {
	spaces, err := mysql.SpacePublicList(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("list space tree fail, err: " + err.Error())
	}
	output := lo.Map(spaces, func(t *mysql.Space, i int) *spacev1.SimpleSpaceInfo {
		return &spacev1.SimpleSpaceInfo{
			SpaceGuid: t.Guid,
			Name:      t.Name,
		}
	})
	return &spacev1.ListPublicSpaceRes{
		SimpleSpaceInfo: output,
	}, nil
}

//// ListSpaceGuidsByUid 根据用户uid，获取到用户的所有空间
//func (GrpcServer) ListSpaceGuidsByUid(ctx context.Context, req *spacev1.ListSpaceGuidsByUidReq) (*spacev1.ListSpaceGuidsByUidRes, error) {
//	list, err := mysql.ListSpaceGuidsByUser(invoker.Db.WithContext(ctx), req.GetUid())
//	if err != nil {
//		return nil, errcodev1.ErrDbError().WithMessage("list space guids by uid, err: " + err.Error())
//	}
//	return &spacev1.ListSpaceGuidsByUidRes{
//		SpaceGuids: list.ToGuids(),
//	}, nil
//}

// SpaceInfo 空间基本信息
func (GrpcServer) SpaceInfo(ctx context.Context, req *spacev1.SpaceInfoReq) (*spacev1.SpaceInfoRes, error) {
	fields := "guid,name,icon_type,icon,`type`,layout,visibility,is_allow_read_member_list,head_image,`desc`,cover,charge_type,origin_price,price,access,link"
	info, err := mysql.GetSpaceInfo(invoker.Db.WithContext(ctx), fields, req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("space info, err: " + err.Error())
	}
	cnt, err := mysql.SpaceMemberCnt(invoker.Db.WithContext(ctx), req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("space info2, err: " + err.Error())
	}
	// 根据space type，获取space的配置信息
	optionList, err := spaceoption.GetListBySpaceType(info.Type)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("space info3, err: " + err.Error())
	}

	optionValueList, err := mysql.GetSpaceOptionList(invoker.Db.WithContext(ctx), req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("space info4, err: " + err.Error())
	}

	ouputOptions := make([]*commonv1.SpaceOption, 0)
	for _, value := range optionList {
		outputOption := &commonv1.SpaceOption{
			Name:            value.Name,
			SpaceOptionId:   value.Option,
			SpaceOptionType: value.Type,
		}
		for _, optionValue := range optionValueList {
			if optionValue.OptionName == value.Option.String() {
				outputOption.Value = optionValue.OptionValue
			}
		}
		ouputOptions = append(ouputOptions, outputOption)
	}
	si := &commonv1.SpaceInfo{
		Guid:         info.Guid,
		Name:         info.Name,
		IconType:     info.IconType,
		Icon:         info.Icon,
		SpaceType:    info.Type,
		SpaceLayout:  info.Layout,
		Visibility:   info.Visibility,
		MemberCnt:    cnt,
		SpaceOptions: ouputOptions,
		ChargeType:   info.ChargeType,
		OriginPrice:  info.OriginPrice,
		Price:        info.Price,
		Desc:         info.Desc,
		HeadImage:    info.HeadImage,
		Cover:        info.Cover,
		Access:       info.Access,
		Link:         info.Link,
	}
	// 根据空间类型，决定是否查询其他额外属性
	switch info.Type {
	}

	return &spacev1.SpaceInfoRes{SpaceInfo: si}, nil
}

// SpaceMemberList Space Member List
func (GrpcServer) SpaceMemberList(ctx context.Context, req *spacev1.SpaceMemberListReq) (*spacev1.SpaceMemberListRes, error) {
	list, err := mysql.SpaceMemberList(invoker.Db.WithContext(ctx), req.SpaceGuid, req.Uids, req.Pagination)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("SpaceMemberList fail, err: " + err.Error())
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
	return &spacev1.SpaceMemberListRes{
		List:       list.ToMemberPb(uMap.GetUserMap(), cmtInfo.GetCommunity().GetUid(), managerMembers),
		Pagination: req.GetPagination(),
	}, nil
}

// SearchSpaceGroupMember ...
func (GrpcServer) SearchSpaceGroupMember(ctx context.Context, req *spacev1.SearchSpaceGroupMemberReq) (*spacev1.SearchSpaceGroupMemberRes, error) {
	list, err := mysql.SpaceGroupMemberSearch(invoker.Db.WithContext(ctx), req.GetSpaceGroupGuid(), req.GetKeyword())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("SearchSpaceGroupMember Fail, err: " + err.Error())
	}
	uids := make([]int64, 0)
	for _, value := range list {
		uids = append(uids, value.Uid)
	}
	uMap, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{UidList: uids})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail, err: " + err.Error())
	}
	return &spacev1.SearchSpaceGroupMemberRes{List: list.ToPb(uMap.GetUserMap())}, nil
}

// SearchSpaceMember ...
func (GrpcServer) SearchSpaceMember(ctx context.Context, req *spacev1.SearchSpaceMemberReq) (*spacev1.SearchSpaceMemberRes, error) {
	list, err := mysql.SpaceMemberSearch(invoker.Db.WithContext(ctx), req.GetSpaceGuid(), req.GetKeyword())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("SearchSpaceMember Fail, err: " + err.Error())
	}
	uids := make([]int64, 0)
	for _, value := range list {
		uids = append(uids, value.Uid)
	}
	uMap, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{UidList: uids})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail, err: " + err.Error())
	}

	return &spacev1.SearchSpaceMemberRes{
		List: list.ToBasePb(uMap.GetUserMap()),
	}, nil
}

// TotalInfo 统计信息
func (GrpcServer) TotalInfo(ctx context.Context, req *spacev1.TotalInfoReq) (*spacev1.TotalInfoRes, error) {
	spaceCnt, err := service.Space.Count(invoker.Db.WithContext(ctx))
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("community info fail1, err: " + err.Error())
	}

	return &spacev1.TotalInfoRes{
		SpaceCnt: spaceCnt,
	}, nil
}

func (GrpcServer) ListSpaceInfo(ctx context.Context, req *spacev1.ListSpaceInfoReq) (*spacev1.ListSpaceInfoRes, error) {
	spaces, err := mysql.SpacePublicListByGuids(invoker.Db.WithContext(ctx), req.GetSpaceGuids())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("list space tree fail, err: " + err.Error())
	}
	return &spacev1.ListSpaceInfoRes{SpaceInfos: spaces.ToPb()}, nil
}

func (GrpcServer) GetMemberStatus(ctx context.Context, req *spacev1.GetMemberStatusReq) (*spacev1.GetMemberStatusRes, error) {
	spaceGuids := lo.Uniq(req.GetSpaceGuids())
	if len(spaceGuids) == 0 {
		return nil, errcodev1.ErrSpaceEmpty()
	}
	sms, err := mysql.GetSpaceMemberBySpaces(invoker.Db.WithContext(ctx), spaceGuids, req.Uid)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetSpaceMemberBySpaces fail, err: " + err.Error())
	}

	inSpaceGuids := lo.Associate(sms, func(ms mysql.SpaceMember) (string, bool) { return ms.Guid, true })
	list := make([]*spacev1.MemberStatus, 0, len(spaceGuids))
	for _, spaceGuid := range spaceGuids {
		isAllowManage, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_MANAGE, req.Uid, spaceGuid)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
		}
		isMember := inSpaceGuids[spaceGuid]
		list = append(list, &spacev1.MemberStatus{
			SpaceGuid:     spaceGuid,
			Uid:           req.Uid,
			IsMember:      isMember,
			IsAllowManage: isAllowManage,
		})
	}
	return &spacev1.GetMemberStatusRes{List: list}, nil
}
