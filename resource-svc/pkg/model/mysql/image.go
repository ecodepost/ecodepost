package mysql

import (
	"gorm.io/gorm"
)

type Image struct {
	Id       int    `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	Ctime    int64  `gorm:"not null" json:"ctime" `          // 创建时间
	Utime    int64  `gorm:"not null" json:"utime" `          // 更新时间
	Dtime    int64  `gorm:"not null" json:"dtime"`           // 删除时间
	Uid      int64  `gorm:"not null;" json:"uid"`            // 创建者
	CmtGuid  string `gorm:"not null"`
	Name     string `gorm:"not null;" json:"name" form:"name" ` //
	FileGuid string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:唯一标识"`
	Type     string `gorm:"not null;" json:"type" form:"type" ` //
	CdnName  string `gorm:"not null;"`
	Path     string `gorm:"not null;" json:"path" form:"path" ` //
}

func (Image) TableName() string {
	return "image"
}

func AddImage(db *gorm.DB, image *Image) error {
	return db.Create(image).Error
}
