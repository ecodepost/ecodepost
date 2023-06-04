package space

import (
	"context"
	"errors"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service/pmspolicy"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	spacev1 "ecodepost/pb/space/v1"
)

// GetSpacePermissionByUid 获取某个用户的空间/分组权限
// 显示、允许任何人加入 commonv1.CMN_VISBL_LV_INTERNAL， commonv1.SPC_ACS_OPEN
func (GrpcServer) GetSpacePermissionByUid(ctx context.Context, req *spacev1.GetSpacePermissionByUidReq) (resp *spacev1.GetSpacePermissionByUidRes, err error) {
	resp = &spacev1.GetSpacePermissionByUidRes{}
	var spaceInfo mysql.Space
	var spaceGroupInfo mysql.SpaceGroup

	// 看空间是不是可以查看的，那么需要先看他是不是private，secret
	switch req.GetGuidType() {
	case commonv1.CMN_GUID_SPACE:
		spaceInfo, err = mysql.GetSpaceInfo(invoker.Db.WithContext(ctx), "visibility,is_allow_read_member_list,`type`,`access`", req.GetTargetGuid())
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("spaceInfo fail1, err: " + err.Error())
		}

		switch spaceInfo.Visibility {
		case commonv1.CMN_VISBL_INTERNAL:
			switch spaceInfo.Access {
			case commonv1.SPC_ACS_OPEN:
				resp.IsAllowView = true
			case commonv1.SPC_ACS_DENY_ALL:
			case commonv1.SPC_ACS_USER_APPLY:
			case commonv1.SPC_ACS_USER_PAY:
			default:
				return nil, errcodev1.ErrInvalidArgument().WithMessage("spaceInfo.Access invalid")
			}
		case commonv1.CMN_VISBL_SECRET:
			switch spaceInfo.Access {
			case commonv1.SPC_ACS_OPEN:
			case commonv1.SPC_ACS_DENY_ALL:
			case commonv1.SPC_ACS_USER_APPLY:
			case commonv1.SPC_ACS_USER_PAY:
			default:
				return nil, errcodev1.ErrInvalidArgument().WithMessage("spaceInfo.Access invalid2")
			}
		}

		// 如果没有登录，那么看空间的属性，看是否可读
		if req.GetOperateUid() == 0 {
			return resp, nil
		}
	case commonv1.CMN_GUID_SPACE_GROUP:
		spaceGroupInfo, err = mysql.GetSpaceGroupInfo(invoker.Db.WithContext(ctx), "visibility", req.GetTargetGuid())
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("spaceGroupInfo fail2, err: " + err.Error())
		}
		// 如果没有登录，那么看空间的属性，看是否可读
		if req.GetOperateUid() == 0 {
			switch spaceGroupInfo.Visibility {
			// 不是空间成员，也是可以查看的
			case commonv1.CMN_VISBL_INTERNAL:
				resp.IsAllowView = true
			}
			return resp, nil
		}
	}

	// 如果是超级管理员
	adminFlag, err := pmspolicy.IsCommunitySuperAdmin(ctx, req.GetOperateUid())
	if err != nil {
		return resp, errcodev1.ErrInternal().WithMessage("fail3, err: " + err.Error())
	}

	// 看空间是不是可以查看的，那么需要先看他是不是private，secret
	switch req.GetGuidType() {
	case commonv1.CMN_GUID_SPACE:
		memberInfo, err := mysql.SpaceMemberInfo(invoker.Db.WithContext(ctx), "id", req.GetTargetGuid(), req.GetOperateUid())
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("fail3, err: " + err.Error())
		}
		if memberInfo.Id > 0 {
			resp.IsMember = true
		}
		resp.IsAllowReadMemberList = spaceInfo.IsAllowReadMemberList

		// 公用权限
		isAllowCreateFile, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_CREATE, req.GetOperateUid(), req.GetTargetGuid())
		if err != nil && !errors.Is(err, pmspolicy.NotSpaceMemberError) {
			return nil, errcodev1.ErrInternal().WithMessage("policy check PMS_FILE_CREATE, err:" + err.Error())
		}
		resp.IsAllowCreateFile = isAllowCreateFile

		isAllowManage, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_MANAGE, req.GetOperateUid(), req.GetTargetGuid())
		if err != nil && !errors.Is(err, pmspolicy.NotSpaceMemberError) {
			return nil, errcodev1.ErrInternal().WithMessage("policy check PMS_FILE_MANAGE, err:" + err.Error())
		}
		resp.IsAllowManage = isAllowManage

		if adminFlag {
			resp.IsAllowView = true
			resp.IsAllowWrite = true
			resp.IsAllowReadMemberList = true
			return resp, nil
		}

		switch spaceInfo.Visibility {
		// 不是空间成员，也是可以查看的
		case commonv1.CMN_VISBL_INTERNAL:
			if memberInfo.Id > 0 {
				resp.IsAllowView = true
				resp.IsAllowWrite = true
			}
			switch spaceInfo.Access {
			case commonv1.SPC_ACS_OPEN:
			case commonv1.SPC_ACS_DENY_ALL:
			case commonv1.SPC_ACS_USER_APPLY:
				// 当还不是成员的时候，查询审核状态
				if memberInfo.Id == 0 {
					// 查看他的审核状态
					auditInfo, err := mysql.AuditIndexInfoByUidAndGuid(invoker.Db.WithContext(ctx), req.GetOperateUid(), req.GetTargetGuid(), commonv1.AUDIT_TYPE_SPACE)
					if err != nil {
						return nil, errcodev1.ErrDbError().WithMessage("AuditApplySpaceMember fail, err: " + err.Error())
					}
					if auditInfo.Id > 0 {
						resp.AuditStatus = auditInfo.Status
					}
				}
			case commonv1.SPC_ACS_USER_PAY:
			}

		// 必须是空间成员才能查看
		// case commonv1.CMN_VISBL_PRIVATE:
		//	// 说明是为该空间成员，那么是可以查看的。
		//	if memberInfo.Id > 0 {
		//		resp.IsAllowView = true
		//		resp.IsAllowWrite = true
		//		// 当还不是成员的时候，查询审核状态
		//	} else {
		//		// 查看他的审核状态
		//		auditInfo, err := mysql.AuditIndexInfoByUidAndGuid(invoker.Db.WithContext(ctx), req.GetOperateUid(), req.GetCmtGuid(), req.GetTargetGuid(), commonv1.AUDIT_TYPE_SPACE)
		//		if err != nil {
		//			return nil, errcodev1.ErrDbError().WithMessage("AuditApplySpaceMember fail, err: " + err.Error())
		//		}
		//		if auditInfo.Id > 0 {
		//			resp.AuditStatus = auditInfo.Status
		//		}
		//	}

		case commonv1.CMN_VISBL_SECRET:
			// 说明是为该空间成员，那么是可以查看的。
			if memberInfo.Id > 0 {
				resp.IsAllowView = true
				resp.IsAllowWrite = true
			}
		}
	case commonv1.CMN_GUID_SPACE_GROUP:
		memberInfo, err := mysql.SpaceGroupMemberInfo(invoker.Db.WithContext(ctx), "id", req.GetTargetGuid(), req.GetOperateUid())
		if err != nil {
			return nil, errcodev1.ErrDbError().WithMessage("fail3, err: " + err.Error())
		}
		if memberInfo.Id > 0 {
			resp.IsMember = true
		}
		if adminFlag {
			resp.IsAllowView = true
			resp.IsAllowWrite = true
			resp.IsAllowReadMemberList = true
			return resp, nil
		}

		resp.IsAllowReadMemberList = spaceGroupInfo.IsAllowReadMemberList

		switch spaceGroupInfo.Visibility {
		// 不是空间成员，也是可以查看的
		case commonv1.CMN_VISBL_INTERNAL:
			resp.IsAllowView = true
			if memberInfo.Id > 0 {
				resp.IsAllowWrite = true
			}
		// 必须是空间成员才能查看
		// case commonv1.CMN_VISBL_PRIVATE:
		//	// 说明是为该空间成员，那么是可以查看的。
		//	if memberInfo.Id > 0 {
		//		resp.IsAllowView = true
		//		resp.IsAllowWrite = true
		//		// 当还不是成员的时候，查询审核状态
		//	} else {
		//		// 查看他的审核状态
		//		auditInfo, err := mysql.AuditIndexInfoByUidAndGuid(invoker.Db.WithContext(ctx), req.GetOperateUid(), req.GetCmtGuid(), req.GetTargetGuid(), commonv1.AUDIT_TYPE_SPACE_GROUP)
		//		if err != nil {
		//			return nil, errcodev1.ErrDbError().WithMessage("AuditApplySpaceMember fail, err: " + err.Error())
		//		}
		//		if auditInfo.Id > 0 {
		//			resp.AuditStatus = auditInfo.Status
		//		}
		//	}
		case commonv1.CMN_VISBL_SECRET:
			// 说明是为该空间成员，那么是可以查看的。
			if memberInfo.Id > 0 {
				resp.IsAllowView = true
				resp.IsAllowWrite = true
			}
		}
	}
	return resp, nil

}
