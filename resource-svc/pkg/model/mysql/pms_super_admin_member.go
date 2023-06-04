package mysql

import (
	"fmt"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	pmsv1 "ecodepost/pb/pms/v1"
	userv1 "ecodepost/pb/user/v1"
	"gorm.io/gorm"
)

// PmsSuperAdminMember 用户组的权限分配
type PmsSuperAdminMember struct {
	Id         int64 `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	Uid        int64 `gorm:"not null; default:0;unique_index:idx_role_uid; comment:使用者UID"`
	CreatedUid int64 `gorm:"not null; default:0; comment:创建者"`
	Ctime      int64 `gorm:"not null; default:0; comment:创建时间"`
}

func (PmsSuperAdminMember) TableName() string {
	return "pms_super_admin_member"
}

type PmsSuperAdminMembers []*PmsSuperAdminMember

func (list PmsSuperAdminMembers) ToUids() []int64 {
	output := make([]int64, 0)
	for _, value := range list {
		output = append(output, value.Uid)
	}
	return output
}

func (list PmsSuperAdminMembers) ToPb(userMap map[int64]*userv1.UserInfo) []*pmsv1.MemberInfo {
	output := make([]*pmsv1.MemberInfo, 0)
	for _, value := range list {
		output = append(output, value.ToPb(userMap))
	}
	return output
}
func (member PmsSuperAdminMember) ToPb(userMap map[int64]*userv1.UserInfo) *pmsv1.MemberInfo {
	activeTime := userMap[member.Uid].GetActiveTime()
	aDateStr := "未活跃过"
	if activeTime != 0 {
		activeDate := time.Unix(activeTime, 0).Format("2006-01-02")
		today := time.Now().Format("2006-01-02")
		yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		if activeDate == today {
			aDateStr = "今天活跃过"
		} else if activeDate == yesterday {
			aDateStr = "昨天活跃过"
		} else {
			aDateStr = activeDate + "活跃过"
		}
	}

	return &pmsv1.MemberInfo{
		Uid:            member.Uid,
		Name:           userMap[member.Uid].GetName(),
		Nickname:       userMap[member.Uid].GetNickname(),
		Avatar:         userMap[member.Uid].GetAvatar(),
		Ctime:          member.Ctime,
		PmsManagerType: commonv1.PMS_MANAGER_SUPER_ADMIN,
		Position:       userMap[member.Uid].GetPosition(),
		ActiveTime:     aDateStr,
	}
}

// GetPmsSuperAdminMembers 超级管理员列表
func GetPmsSuperAdminMembers(db *gorm.DB) (superAdminMembers PmsSuperAdminMembers, err error) {
	superAdminMembers = make(PmsSuperAdminMembers, 0)
	err = db.Find(&superAdminMembers).Error
	if err != nil {
		err = fmt.Errorf("GetPmsSuperAdminMembers fail, err: %w", err)
		return
	}
	return
}

// GetSuperAdminMemberId 创建一条记录
func GetSuperAdminMemberId(db *gorm.DB, uid int64) (info PmsSuperAdminMember, err error) {
	err = db.Select("id").Where(" uid = ?", uid).Find(&info).Error
	return
}

// CreateSuperAdminMemberInBatches 设置了唯一索引，所以有重复创建会直接报错
func CreateSuperAdminMemberInBatches(db *gorm.DB, members []PmsSuperAdminMember) (err error) {
	err = db.CreateInBatches(members, len(members)).Error
	if err != nil {
		err = fmt.Errorf("CreateSuperAdminMemberInBatches fail, err: %w", err)
		return
	}
	return
}

// DeleteSuperAdminMember 设置了唯一索引，所以有重复创建会直接报错
func DeleteSuperAdminMember(db *gorm.DB, uid int64) (err error) {
	err = db.Where("uid = ?", uid).Delete(&PmsSuperAdminMember{}).Error
	if err != nil {
		err = fmt.Errorf("DeleteSuperAdminMember fail, err: %w", err)
		return
	}
	return
}

// SuperAdminMemberCnt 超级管理员个数
func SuperAdminMemberCnt(db *gorm.DB) (cnt int64, err error) {
	if err = db.Model(PmsSuperAdminMember{}).Count(&cnt).Error; err != nil {
		err = fmt.Errorf("CommunityMemberCnt fail, err: %w", err)
		return
	}
	return
}
