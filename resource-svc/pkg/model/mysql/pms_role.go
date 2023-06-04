package mysql

import (
	"fmt"
	"time"

	pmsv1 "ecodepost/pb/pms/v1"

	"gorm.io/gorm"
)

// PmsRole 角色
type PmsRole struct {
	Id    int64  `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	Name  string `gorm:"not null;type:varchar(160)"`
	Ctime int64  `gorm:"not null; default:0; comment:创建时间"`
	Utime int64  `gorm:"not null; default:0; comment:更新时间"`
}

func (PmsRole) TableName() string {
	return "pms_role"
}

type PmsRoles []*PmsRole

func (list PmsRoles) ToPb() []*pmsv1.RoleInfo {
	output := make([]*pmsv1.RoleInfo, 0)
	for _, value := range list {
		output = append(output, value.ToPb())
	}
	return output
}

func (info *PmsRole) ToPb() *pmsv1.RoleInfo {
	return &pmsv1.RoleInfo{
		Id:   info.Id,
		Name: info.Name,
	}
}

func (list PmsRoles) ToRoleIds() []int64 {
	output := make([]int64, 0)
	for _, value := range list {
		output = append(output, value.Id)
	}
	return output
}

// GetPmsManagerRoleCnt 获取role信息，如果不存在自动创建
func GetPmsManagerRoleCnt(db *gorm.DB, cmtGuid string) (cnt int64, err error) {
	// 根据社区id和role type，找到对应的role id
	err = db.Model(PmsSuperAdminMember{}).Where("cmt_guid = ?", cmtGuid).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("GetPmsManagerRoleCnt fail, err: %w", err)
		return
	}
	return
}

// GetPmsRoles 获取role list
func GetPmsRoles(db *gorm.DB) (roleInfo PmsRoles, err error) {
	// 根据社区id和role type，找到对应的role id
	err = db.Where(" dtime = 0 ").Find(&roleInfo).Error
	if err != nil {
		err = fmt.Errorf("GetPmsManagerRole fail, err: %w", err)
		return
	}
	return
}

// CreatePmsRole 创建一条记录
func CreatePmsRole(db *gorm.DB, data *PmsRole) (err error) {
	if err = db.Create(data).Error; err != nil {
		return fmt.Errorf("CreatePmsRole fail, err: %w", err)
	}
	return
}

// UpdatePmsRole 创建一条记录
func UpdatePmsRole(db *gorm.DB, id int64, name string) (err error) {
	if err = db.Model(PmsRole{}).Where("id = ? ", id).Updates(map[string]any{
		"name":  name,
		"utime": time.Now().Unix(),
	}).Error; err != nil {
		return fmt.Errorf("UpdatePmsRole fail, err: %w", err)
	}
	return
}

// GetPmsRole 创建一条记录
func GetPmsRole(db *gorm.DB, id int64) (info PmsRole, err error) {
	if err = db.Select("id,name").Where("id = ? ", id).Find(&info).Error; err != nil {
		err = fmt.Errorf("GetPmsRole fail, err: %w", err)
		return
	}
	return
}

// DeletePmsRole 删除数据
func DeletePmsRole(db *gorm.DB, id int64) (err error) {
	if err = db.Where("id = ? ", id).Delete(&PmsRole{}).Error; err != nil {
		return fmt.Errorf("DeletePmsRole fail, err: %w", err)
	}
	err = db.Where("role_id = ?", id).Delete(&PmsRoleMember{}).Error
	if err != nil {
		return fmt.Errorf("DeletePmsRole fail2, err: %w", err)
	}
	err = db.Where("role_id = ?", id).Delete(&PmsRoleSpace{}).Error
	if err != nil {
		return fmt.Errorf("DeletePmsRole fail3, err: %w", err)
	}
	return
}
