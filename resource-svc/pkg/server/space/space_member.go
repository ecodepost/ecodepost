package space

import (
	"context"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	spacev1 "ecodepost/pb/space/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/gotomicro/ego/core/elog"
)

// AddSpaceMember ...
func (GrpcServer) AddSpaceMember(ctx context.Context, req *spacev1.AddSpaceMemberReq) (*spacev1.AddSpaceMemberRes, error) {
	// 1. 查询空间信息，如果空间不存在则报错
	spaceInfo, err := mysql.GetSpaceInfo(invoker.Db.WithContext(ctx), "`type`,`space_group_guid`", req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AddSpaceMember SpaceGetSpaceType fail, err: " + err.Error())
	}
	list := make(mysql.SpaceMembers, 0)

	// 2. 判断需要添加的用户是否已加入空间
	existList, err := mysql.GetSpaceMemberByBatch(invoker.Db.WithContext(ctx), req.GetSpaceGuid(), req.GetAddUids())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail2, err: " + err.Error())
	}

	// 3. 传过的uid, 可能已经是成员了, 找到不存在的uids
	notExistUids := make([]int64, 0)
	for _, uid := range req.GetAddUids() {
		existFlag := false
		for _, value := range existList {
			if value.Uid == uid {
				existFlag = true
				break
			}
		}
		if existFlag {
			continue
		}
		notExistUids = append(notExistUids, uid)
	}

	// 4. 查找不存在的用户详情
	uMap, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{UidList: notExistUids})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail, err: " + err.Error())
	}

	// 5. 插入到空间
	nowTime := time.Now().Unix()
	for _, uid := range notExistUids {
		var nickname string
		if uMap.GetUserMap() != nil {
			nickname = uMap.GetUserMap()[uid].GetNickname()
		}
		list = append(list, mysql.SpaceMember{
			Ctime:     nowTime,
			Utime:     nowTime,
			Uid:       uid,
			Nickname:  nickname,
			Guid:      req.GetSpaceGuid(),
			CreatedBy: req.GetOperateUid(),
			UpdatedBy: req.GetOperateUid(),
		})
	}
	err = service.Space.CreateBatchMember(ctx, invoker.Db.WithContext(ctx), list)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AddSpaceGroupMember Fail, err: " + err.Error())
	}

	// 6. 判断是否需要加入到社区
	for _, uid := range req.GetAddUids() {
		info, err := mysql.SpaceGroupMemberInfo(invoker.Db.WithContext(ctx), "id", spaceInfo.SpaceGroupGuid, uid)
		if err != nil {
			elog.Error("SpaceGroupMemberInfo fail", elog.FieldErr(err))
			continue
		}
		if info.Id == 0 {
			var nickname string
			if uMap.GetUserMap() != nil {
				nickname = uMap.GetUserMap()[uid].GetNickname()
			}
			err = mysql.SpaceGroupMemberCreate(invoker.Db.WithContext(ctx), &mysql.SpaceGroupMember{
				Ctime:     nowTime,
				Utime:     nowTime,
				Uid:       uid,
				Nickname:  nickname,
				Guid:      spaceInfo.SpaceGroupGuid,
				CreatedBy: req.GetOperateUid(),
			})
			if err != nil {
				elog.Error("SpaceGroupMemberCreate fail", elog.FieldErr(err))
				continue
			}
		}
	}
	return &spacev1.AddSpaceMemberRes{}, nil
}

// DeleteSpaceMember ...
func (GrpcServer) DeleteSpaceMember(ctx context.Context, req *spacev1.DeleteSpaceMemberReq) (*spacev1.DeleteSpaceMemberRes, error) {
	err := service.Space.DeleteBatchMember(ctx, invoker.Db.WithContext(ctx), req.GetSpaceGuid(), req.GetDeleteUids(), commonv1.TRACK_TOTAL_EVENT_SPACE_DELETE_MEMBER)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("DeleteSpaceMember Fail, err: " + err.Error())
	}
	return &spacev1.DeleteSpaceMemberRes{}, nil
}

func (GrpcServer) AddSpaceGroupMember(ctx context.Context, req *spacev1.AddSpaceGroupMemberReq) (*spacev1.AddSpaceGroupMemberRes, error) {
	list := make(mysql.SpaceGroupMembers, 0)
	existList, err := mysql.GetSpaceGroupMemberByBatch(invoker.Db.WithContext(ctx), req.GetSpaceGroupGuid(), req.GetAddUids())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail2, err: " + err.Error())
	}

	// 传过的uid，可能以及是成员了
	// 找到不存在的uids
	notExistUids := make([]int64, 0)
	for _, uid := range req.GetAddUids() {
		existFlag := false
		for _, value := range existList {
			if value.Uid == uid {
				existFlag = true
				break
			}
		}
		if existFlag {
			continue
		}
		notExistUids = append(notExistUids, uid)
	}

	uMap, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{
		UidList: notExistUids,
	})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("MemberList Fail, err: " + err.Error())
	}

	nowTime := time.Now().Unix()
	for _, uid := range notExistUids {
		var nickname string
		if uMap.GetUserMap() != nil {
			userInfo, flag := uMap.GetUserMap()[uid]
			if flag {
				nickname = userInfo.GetNickname()
			}
		}
		list = append(list, mysql.SpaceGroupMember{
			Ctime:     nowTime,
			Utime:     nowTime,
			Uid:       uid,
			Nickname:  nickname,
			Guid:      req.GetSpaceGroupGuid(),
			CreatedBy: req.GetOperateUid(),
			UpdatedBy: req.GetOperateUid(),
		})
	}
	err = mysql.SpaceGroupMemberBatchCreate(invoker.Db.WithContext(ctx), list)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AddSpaceGroupMember Fail, err: " + err.Error())
	}

	spaceList, err := mysql.SpaceListBySpaceGroupGuid(invoker.Db.WithContext(ctx), req.GetSpaceGroupGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AddSpaceGroupMember Fail, err: " + err.Error())
	}

	for _, spaceInfo := range spaceList {
		list := make(mysql.SpaceMembers, 0)
		for _, uid := range req.GetAddUids() {
			var nickname string
			if uMap.GetUserMap() != nil {
				nickname = uMap.GetUserMap()[uid].GetNickname()
			}
			list = append(list, mysql.SpaceMember{
				Ctime:     nowTime,
				Utime:     nowTime,
				Uid:       uid,
				Nickname:  nickname,
				Guid:      spaceInfo.Guid,
				CreatedBy: req.GetOperateUid(),
				UpdatedBy: req.GetOperateUid(),
			})
		}

		err = service.Space.CreateBatchMember(ctx, invoker.Db.WithContext(ctx), list)
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("AddSpaceGroupMember Fail, err: " + err.Error())
		}
		// switch spaceInfo.Type {
		// case commonv1.SPC_TYPE_IM:
		//	info, err := invoker.GrpcIm.FindOneGroup(ctx, &imv1.FindOneGroupReq{
		//		CmtGuid:   req.GetCmtGuid(),
		//		SpaceGuid: spaceInfo.Guid,
		//	})
		//	if err != nil {
		//		return nil, errcodev1.ErrInternal().WithMessage("AddSpaceMember FindOneGroup fail, err: " + err.Error())
		//	}
		//	_, err = invoker.GrpcIm.InviteUserToGroup(ctx, &imv1.InviteUserToGroupReq{
		//		GroupId: info.GetInfo().GetGroupId(),
		//		Uids:    req.GetAddUids(),
		//	})
		//	if err != nil {
		//		return nil, errcodev1.ErrInternal().WithMessage("AddSpaceMember InviteUserToGroup fail, err: " + err.Error())
		//	}
		// }
	}

	return &spacev1.AddSpaceGroupMemberRes{}, nil
}

func (GrpcServer) DeleteSpaceGroupMember(ctx context.Context, req *spacev1.DeleteSpaceGroupMemberReq) (*spacev1.DeleteSpaceGroupMemberRes, error) {
	err := mysql.SpaceGroupMemberBatchDelete(invoker.Db.WithContext(ctx), req.GetSpaceGroupGuid(), req.GetDeleteUids())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("DeleteSpaceGroupMember Fail, err: " + err.Error())
	}

	spaceList, err := mysql.SpaceListBySpaceGroupGuid(invoker.Db.WithContext(ctx), req.GetSpaceGroupGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AddSpaceGroupMember Fail, err: " + err.Error())
	}
	for _, spaceInfo := range spaceList {
		err := service.Space.DeleteBatchMember(ctx, invoker.Db.WithContext(ctx), spaceInfo.Guid, req.GetDeleteUids(), commonv1.TRACK_TOTAL_EVENT_SPACE_DELETE_MEMBER)
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("DeleteSpaceMember Fail, err: " + err.Error())
		}
	}
	return &spacev1.DeleteSpaceGroupMemberRes{}, nil
}

func (GrpcServer) GenSpaceAccessOrder(ctx context.Context, req *spacev1.GenSpaceAccessOrderReq) (*spacev1.GenSpaceAccessOrderRes, error) {
	// db := invoker.Db.WithContext(ctx)
	// var space mysql.Space
	// if err := db.Select("id,name,charge_type,origin_price,price").Where("cmt_guid = ? AND guid = ?", req.CmtGuid, req.SpaceGuid).Find(&space).Error; err != nil {
	// 	return nil, errcodev1.ErrDbError().WithMessage("GetSpaceSetInfo fail, " + err.Error())
	// }
	// if space.ChargeType == commonv1.SPC_CT_FREE {
	// 	return nil, errcodev1.ErrDbError().WithMessage("未开启付费")
	// }
	//
	// // 获取用户信息
	// userInfo, err := invoker.GrpcUser.Info(ctx, &userv1.InfoReq{Uid: req.Uid})
	// if err != nil {
	// 	return nil, errcodev1.ErrDbError().WithMessage("member ship info fail3, err: " + err.Error())
	// }
	//
	// // 生成订单
	// ogs := make([]*orderv1.OrderGood, 0, 1)
	// title := space.Name + "的资格购买"
	// ogs = append(ogs, &orderv1.OrderGood{
	// 	Title:   title,
	// 	Price:   space.Price,
	// 	GoodId:  space.Id,
	// 	CmtGuid: req.CmtGuid,
	// })
	// extStr, err := shop.SpaceMembershipExt{SpaceMemberId: req.Uid, SpaceGuid: req.SpaceGuid}.Encode()
	// if err != nil {
	// 	return nil, errcodev1.ErrDbError().WithMessage("member ship info fail4, err: " + err.Error())
	// }
	// // 生成拆分订单
	// o, err := invoker.GrpcOrder.CreateOrder(ctx, &orderv1.CreateOrderReq{
	// 	OrderInfo: &orderv1.CreateOrderInfo{
	// 		CmtGuid:       req.CmtGuid,
	// 		BuyerId:       req.Uid,
	// 		BuyerName:     userInfo.User.Nickname,
	// 		BuyerPhone:    "",
	// 		BuyerEmail:    userInfo.User.Email,
	// 		BuyerAvatar:   userInfo.User.Avatar,
	// 		Remark:        "buy space member",
	// 		TotalAmount:   space.Price,
	// 		ChargeMethod:  req.GetChargeMethod(),
	// 		OrderGoodList: ogs,
	// 		OrderCase:     commonv1.ORDER_CASE_NORMAL,
	// 		OrderType:     commonv1.ORDER_TP_SPACE_MEMBER,
	// 		ClientIp:      "",
	// 		Title:         title,
	// 		To:            commonv1.ORDER_TO_USER,
	// 		TargetId:      cast.ToString(space.Id),
	// 		Ext:           extStr,
	// 	},
	// })
	// if err != nil {
	// 	return nil, errcodev1.ErrInternal().WithMessage("create order fail, err:" + err.Error())
	// }
	//
	// return &spacev1.GenSpaceAccessOrderRes{
	// 	OriginalAmount: space.Price,
	// 	DiscountAmount: space.OriginPrice - space.Price,
	// 	TradeAmount:    space.Price,
	// 	OrderSn:        o.OrderInfo.Sn,
	// }, nil
	return &spacev1.GenSpaceAccessOrderRes{}, nil
}

func (GrpcServer) QuitSpaceMember(ctx context.Context, req *spacev1.QuitSpaceMemberReq) (*spacev1.QuitSpaceMemberRes, error) {
	err := service.Space.DeleteBatchMember(ctx, invoker.Db.WithContext(ctx), req.GetSpaceGuid(), []int64{req.GetUid()}, commonv1.TRACK_TOTAL_EVENT_SPACE_QUIT_MEMBER)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("DeleteSpaceMember Fail, err: " + err.Error())
	}
	return &spacev1.QuitSpaceMemberRes{}, nil
}
