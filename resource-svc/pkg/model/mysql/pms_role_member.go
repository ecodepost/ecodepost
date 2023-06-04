package mysql

import (
	"fmt"

	pmsv1 "ecodepost/pb/pms/v1"
	userv1 "ecodepost/pb/user/v1"

	"gorm.io/gorm"
)

// PmsRoleMember 用户组的权限分配
type PmsRoleMember struct {
	Id        int64  `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	RoleId    int64  `gorm:"not null; default:0;unique_index:idx_role_uid; comment:角色ID"`
	Uid       int64  `gorm:"not null; default:0;unique_index:idx_role_uid; comment:使用者UID"`
	CmtGuid   string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index:cmt_guid"`
	CreatedBy int64  `gorm:"not null; default:0; comment:创建者"`
	Ctime     int64  `gorm:"not null; default:0; comment:创建时间"`
}

func (PmsRoleMember) TableName() string {
	return "pms_role_member"
}

type PmsRoleMembers []PmsRoleMember

func (list PmsRoleMembers) ToPb(userMap map[int64]*userv1.UserInfo) []*pmsv1.MemberInfo {
	output := make([]*pmsv1.MemberInfo, 0)
	for _, value := range list {
		output = append(output, value.ToPb(userMap))
	}
	return output
}

func (member PmsRoleMember) ToPb(userMap map[int64]*userv1.UserInfo) *pmsv1.MemberInfo {
	return &pmsv1.MemberInfo{
		Uid:      member.Uid,
		Nickname: userMap[member.Uid].GetNickname(),
		Avatar:   userMap[member.Uid].GetAvatar(),
		Ctime:    member.Ctime,
	}
}

func (list PmsRoleMembers) ToRoleIds() []int64 {
	output := make([]int64, 0)
	for _, value := range list {
		output = append(output, value.RoleId)
	}
	return output
}

func (list PmsRoleMembers) ToUids() []int64 {
	output := make([]int64, 0)
	for _, value := range list {
		output = append(output, value.Uid)
	}
	return output
}

// GetRoleMemberId 创建一条记录
func GetRoleMemberId(db *gorm.DB, uid int64, roleId int64) (info PmsRoleMember, err error) {
	err = db.Select("id").Where("  uid = ? and role_id = ?", uid, roleId).Find(&info).Error
	return
}

// RoleIdsByUid 查找一个用户所有role ids
func RoleIdsByUid(db *gorm.DB, uid int64) (roleIds []int64, err error) {
	roleMembers := make(PmsRoleMembers, 0)
	err = db.Select("role_id").Where("uid = ?", uid).Find(&roleMembers).Error
	if err != nil {
		err = fmt.Errorf("roleMembers find fail, err: %w", err)
		return
	}
	return roleMembers.ToRoleIds(), nil
}

// DeleteManagerMembers 查询管理员的成员：通过社区ID，角色类型来查询
// 1 先找到role id
// 2 再找到members
func DeleteManagerMembers(db *gorm.DB, req *pmsv1.DeleteManagerMemberReq) (err error) {
	err = DeleteSuperAdminMember(db, req.GetUid())
	if err != nil {
		err = fmt.Errorf("DeleteSuperAdminMember fail3, err: %w", err)
		return
	}
	return
}

func GetRoleMembers(db *gorm.DB, roleId int64) (roleMembers PmsRoleMembers, err error) {
	roleMembers = make(PmsRoleMembers, 0)
	err = db.Where("role_id = ?", roleId).Find(&roleMembers).Error
	if err != nil {
		err = fmt.Errorf("GetRoleMembers fail, err: %w", err)
		return
	}
	return
}

// CreateRoleMemberInBatches 设置了唯一索引，所以有重复创建会直接报错
func CreateRoleMemberInBatches(db *gorm.DB, members []PmsRoleMember) (err error) {
	err = db.CreateInBatches(members, len(members)).Error
	if err != nil {
		err = fmt.Errorf("CreateRoleMemberInBatches fail, err: %w", err)
		return
	}
	return
}

// DeleteRoleMember 设置了唯一索引，所以有重复创建会直接报错
func DeleteRoleMember(db *gorm.DB, roleId, uid int64) (err error) {
	err = db.Where("uid = ? and role_id = ?", uid, roleId).Delete(&PmsRoleMember{}).Error
	if err != nil {
		err = fmt.Errorf("DeleteRoleMember fail, err: %w", err)
		return
	}
	return
}
func IsExistRoleMember(db *gorm.DB, roleId int64, uid int64) (flag bool, err error) {
	info := PmsRoleMember{}
	err = db.Select("id").Where("role_id = ? and uid = ?", roleId, uid).Find(&info).Error
	if err != nil {
		err = fmt.Errorf("IsExistRole fail, err: %w", err)
		return
	}
	if info.Id > 0 {
		flag = true
		return
	}
	return
}

// RoleMemberCnt 超级管理员个数
func RoleMemberCnt(db *gorm.DB) (cnt int64, err error) {
	if err = db.Model(PmsRoleMember{}).Count(&cnt).Error; err != nil {
		err = fmt.Errorf("RoleMemberCnt fail, err: %w", err)
		return
	}
	return
}
