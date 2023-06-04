package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"

	"gorm.io/gorm"
)

// PmsRoleSpace 角色
type PmsRoleSpace struct {
	Id       int64             `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	RoleId   int64             `gorm:"not null;type:varchar(160)"`
	Guid     string            `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index:guid"`
	GuidType commonv1.CMN_GUID `gorm:"not null; default:0; comment:Guid Type"`
	Ctime    int64             `gorm:"not null; default:0; comment:创建时间"`
}

type PmsRoleSpaces []PmsRoleSpace

func (PmsRoleSpace) TableName() string {
	return "pms_role_space"
}

func (list PmsRoleSpaces) ToGuids() []string {
	output := make([]string, 0)
	for _, value := range list {
		output = append(output, value.Guid)
	}
	return output
}

// PmsRoleSpaceList 根据主键查询多条记录
func PmsRoleSpaceList(db *gorm.DB, roleId int64, guidType commonv1.CMN_GUID) (list PmsRoleSpaces, err error) {
	if err = db.Select("guid").Where("role_id = ? and guid_type = ? ", roleId, guidType.Number()).Find(&list).Error; err != nil {
		err = fmt.Errorf("PmsRoleSpaceList failed,err: %w", err)
		return
	}
	return
}
