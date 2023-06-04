package space

import (
	"context"
	"time"

	communityv1 "ecodepost/pb/community/v1"
	notifyv1 "ecodepost/pb/notify/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/utils"
	"github.com/spf13/cast"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	spacev1 "ecodepost/pb/space/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/gotomicro/ego/core/elog"
)

// AuditApplySpaceMember 申请加入空间，只有是private空间，需要这个操作
// todo 一天空间最多申请三次
func (GrpcServer) AuditApplySpaceMember(ctx context.Context, req *spacev1.AuditApplySpaceMemberReq) (*spacev1.AuditApplySpaceMemberRes, error) {
	// 先判断空间类型，如果是internal空间，直接加入，如果是private空间，审核加入
	// switch req.GetAuditType() {
	// case commonv1.AUDIT_TYPE_SPACE_GROUP:
	//	return nil, errcodev1.ErrInternal().WithMessage("not support")
	//	//spaceGroupInfo, _ := mysql.GetSpaceGroupInfo(invoker.Db.WithContext(ctx), "name,visibility", req.GetCmtGuid(), req.GetTargetGuid())
	//	//infoVisibility = spaceGroupInfo.Visibility
	// case commonv1.AUDIT_TYPE_SPACE:
	//	spaceInfo, _ := mysql.GetSpaceInfo(invoker.Db.WithContext(ctx), "name,visibility", req.GetCmtGuid(), req.GetTargetGuid())
	//	infoVisibility = spaceInfo.Visibility
	//
	// }

	spaceInfo, err := mysql.GetSpaceInfo(invoker.Db.WithContext(ctx), "guid,name,visibility,access", req.GetTargetGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("AuditApplySpaceMember fail0, err: " + err.Error())
	}
	if spaceInfo.Guid == "" {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("space info empty")
	}
	infoVisibility := spaceInfo.Visibility

	// 查看下是否
	cmtInfo, err := invoker.GrpcCommunity.Info(ctx, &communityv1.InfoReq{
		Uid: req.GetOperateUid(),
	})
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("AuditApplySpaceMember fail1, err: " + err.Error())
	}

	nowTime := time.Now().Unix()

	userInfo, err := invoker.GrpcUser.Info(ctx, &userv1.InfoReq{
		Uid: req.GetOperateUid(),
	})
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("AuditApplySpaceMember fail2, err: " + err.Error())
	}

	// 如果用户不是社区会员，同时不是开放社区，那么报错
	if cmtInfo.Community.Access != commonv1.CMT_ACS_OPEN {
		return nil, errcodev1.ErrInternal().WithMessage("AuditApplySpaceMember not support this type")
	}

	switch infoVisibility {
	case commonv1.CMN_VISBL_INTERNAL:
		switch spaceInfo.Access {
		// 直接加入
		case commonv1.SPC_ACS_OPEN:
			err = service.Space.CreateMember(ctx, invoker.Db.WithContext(ctx), &mysql.SpaceMember{
				Ctime:     nowTime,
				Utime:     nowTime,
				Uid:       req.GetOperateUid(),
				Nickname:  userInfo.User.GetNickname(),
				Guid:      req.GetTargetGuid(),
				CreatedBy: 9999, // 系统加入
			})
			if err != nil {
				return nil, errcodev1.ErrDbError().WithMessage("AuditApplySpaceMember fail, err: " + err.Error())
			}
			return &spacev1.AuditApplySpaceMemberRes{BizCode: 1}, nil
		// 审核加入
		case commonv1.SPC_ACS_USER_APPLY:
			// auditInfo, err := mysql.AuditIndexInfoByUidAndGuid(invoker.Db.WithContext(ctx), req.GetOperateUid(), req.GetCmtGuid(), req.GetTargetGuid(), commonv1.AUDIT_TYPE_SPACE)
			// if err != nil {
			//	return nil, errcodev1.ErrDbError().WithMessage("AuditApplySpaceMember fail, err: " + err.Error())
			// }
			// if auditInfo.Id == 0 {
			//	auditIndex := &mysql.AuditIndex{
			//		Uid:       req.GetOperateUid(),
			//		CmtGuid:   req.GetCmtGuid(),
			//		AuditGuid: req.GetTargetGuid(),
			//		AuditType: req.GetAuditType(),
			//		Status:    commonv1.AUDIT_STATUS_APPLY,
			//	}
			// } else {
			//
			// }

			tx := invoker.Db.WithContext(ctx).Begin()
			auditIndex := &mysql.AuditIndex{
				Uid:       req.GetOperateUid(),
				AuditGuid: req.GetTargetGuid(),
				AuditType: req.GetAuditType(),
			}
			if err = mysql.AuditLogCreate(tx, auditIndex, req.GetReason()); err != nil {
				tx.Rollback()
				return nil, errcodev1.ErrDbError().WithMessage("AuditLogCreate fail, err: " + err.Error())
			}
			tx.Commit()

			// 查询到社区的超级管理员
			allManagerMembers, err := mysql.GetPmsSuperAdminMembers(invoker.Db.WithContext(ctx))
			if err != nil {
				return nil, errcodev1.ErrDbError().WithMessage("AuditLogCreate fail2, err: " + err.Error())
			}
			// todo 查询到空间有该权限的管理员，后面在写
			metaData := make(map[string]any, 0)
			notificationType := commonv1.NTF_TYPE_INVALID
			switch req.GetAuditType() {
			// case commonv1.AUDIT_TYPE_SPACE_GROUP:
			//	notificationType = commonv1.NTF_TYPE_APPLY_SPACE_GROUP
			//	notificationType = commonv1.NTF_TYPE_APPLY_SPACE
			//	metaData["uid"] = req.GetOperateUid()
			//	metaData["nickname"] = userInfo.GetNickname()
			//	metaData["avatar"] = userInfo.GetAvatar()
			//	metaData["spaceGroupGuid"] = req.GetTargetGuid()
			//	metaData["spaceGroupName"] = infoName
			case commonv1.AUDIT_TYPE_SPACE:
				notificationType = commonv1.NTF_TYPE_APPLY_SPACE
				metaData["uid"] = req.GetOperateUid()
				metaData["nickname"] = userInfo.User.GetNickname()
				metaData["avatar"] = userInfo.User.GetAvatar()
				metaData["spaceGuid"] = req.GetTargetGuid()
				metaData["spaceName"] = spaceInfo.Name
			}
			uids := allManagerMembers.ToUids()
			msgs := make([]*notifyv1.Msg, 0, len(uids))
			for _, v := range uids {
				msgs = append(msgs, &notifyv1.Msg{Receiver: cast.ToString(v)})
			}
			_, err = invoker.GrpcNotify.SendMsg(ctx, &notifyv1.SendMsgReq{
				TplId: 3, // 默认站内信模板id
				Msgs:  msgs,
				VarLetter: &notifyv1.Letter{
					Type:     notificationType,
					TargetId: cast.ToString(auditIndex.Id),
					Meta:     utils.NewMeta(metaData),
				},
			})
			if err != nil {
				return nil, errcodev1.ErrDbError().WithMessage("AuditLogCreate notify fail, err: " + err.Error())
			}
			// 不允许通过这种方式付费加入
		// case commonv1.SPC_ACS_USER_PAY:
		// 自己不允许加入
		// case commonv1.SPC_ACS_DENY_ALL:
		default:
			return nil, errcodev1.ErrInvalidArgument().WithMessage("not support type: ")
		}
	// 其他例如私密类型直接报错
	default:
		return nil, errcodev1.ErrInvalidArgument().WithMessage("visibility type fail, not support type, visibility: ")
	}

	return &spacev1.AuditApplySpaceMemberRes{}, nil
}

// AuditListSpaceMember 空间的审核列表，管理者使用
func (GrpcServer) AuditListSpaceMember(ctx context.Context, req *spacev1.AuditListSpaceMemberReq) (*spacev1.AuditListSpaceMemberRes, error) {
	list, err := service.Audit.ListPage(invoker.Db.WithContext(ctx), req.GetTargetGuid(), req.GetAuditType(), req.GetPagination())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AuditListSpaceMember fail, err: " + err.Error())
	}

	resp, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{
		UidList: list.ToUids(),
	})
	// 不return，用户服务不影响，展示，记录log
	if err != nil {
		elog.Error("list document user error", elog.FieldErr(err))
	}
	return &spacev1.AuditListSpaceMemberRes{
		List:       list.ToPb(resp.GetUserMap()),
		Pagination: req.GetPagination(),
	}, nil

}

// AuditPassSpaceMember 通过用户，管理者使用
func (GrpcServer) AuditPassSpaceMember(ctx context.Context, req *spacev1.AuditPassSpaceMemberReq) (*spacev1.AuditPassSpaceMemberRes, error) {
	_, err := service.Audit.Pass(ctx, req.GetId(), req.GetOperateUid(), req.GetOpReason())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AuditPassSpaceMember fail, err: " + err.Error())
	}
	return &spacev1.AuditPassSpaceMemberRes{}, nil
}

// AuditRejectSpaceMember 拒绝原因，并且可以，禁止在申请，管理者使用
func (GrpcServer) AuditRejectSpaceMember(ctx context.Context, req *spacev1.AuditRejectSpaceMemberReq) (*spacev1.AuditRejectSpaceMemberRes, error) {
	err := service.Audit.Reject(ctx, req.GetId(), req.GetOperateUid(), req.GetOpReason())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AuditRejectSpaceMember fail, err: " + err.Error())
	}
	return &spacev1.AuditRejectSpaceMemberRes{}, nil
}

// AuditMapByIds 根据IDS，获取审核的map信息
func (GrpcServer) AuditMapByIds(ctx context.Context, req *spacev1.AuditMapByIdsReq) (*spacev1.AuditMapByIdsRes, error) {
	list, err := service.Audit.Map(invoker.Db.WithContext(ctx), req.GetAuditIds())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("AuditMapByIds fail, err: " + err.Error())
	}
	return &spacev1.AuditMapByIdsRes{AuditMap: list}, nil
}
