package mysql

import (
	commonv1 "ecodepost/pb/common/v1"
	userv1 "ecodepost/pb/user/v1"
)

func toMemberRole(isSuperAdmin bool, memberUid int64, memberCtime int64, userMap map[int64]*userv1.UserInfo, createdUid int64) *commonv1.MemberRole {
	pmsManagerType := commonv1.PMS_MANAGER_INVALID
	// todo 目前只有超级管理员
	if isSuperAdmin {
		pmsManagerType = commonv1.PMS_MANAGER_SUPER_ADMIN
	}

	if memberUid == createdUid {
		pmsManagerType = commonv1.PMS_MANAGER_CREATE
	}
	return &commonv1.MemberRole{
		Uid:            memberUid,
		Nickname:       userMap[memberUid].GetNickname(),
		Avatar:         userMap[memberUid].GetAvatar(),
		Ctime:          memberCtime,
		PmsManagerType: pmsManagerType,
	}
}
