package mysql

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// FileSpaceTop 空间置顶
type FileSpaceTop struct {
	Id        int64  `json:"id" gorm:"not null;primary_key;auto_increment"`
	Guid      string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"`
	SpaceGuid string `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:空间ID"`
	CreatedBy int64  `gorm:"not null; default:0; comment:创建人"`
	Sort      int64  `gorm:"not null; default:0; comment:排序"`
	Ctime     int64  `gorm:"not null; default:0; comment:创建时间"`
}

func (FileSpaceTop) TableName() string {
	return "file_space_top"
}

// FileSpaceTopCreate 创建space top
func FileSpaceTopCreate(db *gorm.DB, fileGuid, spaceGuid string, uid int64) (err error) {
	err = db.Create(&FileSpaceTop{
		Guid:      fileGuid,
		SpaceGuid: spaceGuid,
		CreatedBy: uid,
		Sort:      time.Now().UnixMilli(),
		Ctime:     time.Now().Unix(),
	}).Error
	if err != nil {
		return fmt.Errorf("FileSpaceTopCreate fail, err: %w", err)
	}
	return nil
}

func FileSpaceTopCnt(db *gorm.DB, spaceGuid string) (cnt int64, err error) {
	err = db.Model(FileSpaceTop{}).Where("space_guid = ?", spaceGuid).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("FileSpaceTopCnt fail, err: %w", err)
		return
	}
	return
}

// FileSpaceTopId 获取id号
func FileSpaceTopId(db *gorm.DB, fileGuid, spaceGuid string) (id int64, err error) {
	var info FileSpaceTop
	err = db.Select("id").Where("guid = ? and  space_guid = ?", fileGuid, spaceGuid).Find(&info).Error
	if err != nil {
		return 0, fmt.Errorf("FileSpaceTopId fail, err: %w", err)
	}
	id = info.Id
	return
}

// FileSpaceTopDelete 获取id号
func FileSpaceTopDelete(db *gorm.DB, fileGuid string) (err error) {
	err = db.Where("guid = ?", fileGuid).Delete(&FileSpaceTop{}).Error
	if err != nil {
		return fmt.Errorf("FileSpaceTopDelete fail, err: %w", err)
	}
	return
}

// FileSpaceTopGuids 获取list
func FileSpaceTopGuids(db *gorm.DB, spaceGuid string) (fileGuids []string, err error) {
	var list []FileSpaceTop
	err = db.Select("guid").Where(" space_guid = ?", spaceGuid).Find(&list).Error
	if err != nil {
		return
	}
	fileGuids = make([]string, 0)
	for _, value := range list {
		fileGuids = append(fileGuids, value.Guid)
	}
	return
}
