package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"

	"gorm.io/gorm"
)

type Resource struct {
	Id        int64             `gorm:"not null;primary_key;AUTO_INCREMENT = 1;comment:主键id" json:"id"`
	Guid      string            `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:guid信息"`
	GuidType  commonv1.CMN_GUID `gorm:"not null;default:0;comment:guid类型"`
	Ctime     int64             `gorm:"not null;default:0;comment:创建时间"`
	CreatedBy int64             `gorm:"not null;default:0;comment:创建人"`
}

func (Resource) TableName() string {
	return "resource"
}

// ResourceCreate 创建一条记录
func ResourceCreate(db *gorm.DB, data *Resource) (err error) {
	if err = db.Create(data).Error; err != nil {
		return fmt.Errorf("ResourceCreate failed,err: %w", err)
	}
	return
}
