package mysql

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// FileRecommend 文章推荐
type FileRecommend struct {
	Id        int64  `json:"id" gorm:"not null;primary_key;auto_increment"`
	Guid      string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"`
	SpaceGuid string `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:空间ID"`
	CreatedBy int64  `gorm:"not null; default:0; comment:创建人"`
	Sort      int64  `gorm:"not null; default:0; comment:排序"`
	Ctime     int64  `gorm:"not null; default:0; comment:创建时间"`
}

func (FileRecommend) TableName() string {
	return "file_recommend"
}

// FileRecommendCreate 创建space top
func FileRecommendCreate(db *gorm.DB, fileGuid, spaceGuid string, uid int64) (err error) {
	err = db.Create(&FileRecommend{
		Guid:      fileGuid,
		SpaceGuid: spaceGuid,
		Sort:      time.Now().UnixMilli(),
		CreatedBy: uid,
		Ctime:     time.Now().Unix(),
	}).Error
	if err != nil {
		return fmt.Errorf("FileRecommendCreate fail, err: %w", err)
	}
	return nil
}

// FileRecommendId 获取id号
func FileRecommendId(db *gorm.DB, fileGuid string) (id int64, err error) {
	var info FileRecommend
	err = db.Select("id").Where("guid = ?", fileGuid).Find(&info).Error
	if err != nil {
		return 0, fmt.Errorf("FileRecommendId fail, err: %w", err)
	}
	id = info.Id
	return
}

// FileRecommendDelete 获取id号
func FileRecommendDelete(db *gorm.DB, fileGuid string) (err error) {
	err = db.Where("guid = ?", fileGuid).Delete(&FileRecommend{}).Error
	if err != nil {
		return fmt.Errorf("FileRecommendDelete fail, err: %w", err)
	}
	return
}

// FileRecommendGuids 获取list
func FileRecommendGuids(db *gorm.DB, spaceGuid string) (fileGuids []string, err error) {
	var list []FileRecommend
	err = db.Select("guid").Where("space_guid = ?", spaceGuid).Order("sort asc").Limit(5).Find(&list).Error
	if err != nil {
		return
	}
	fileGuids = make([]string, 0)
	for _, value := range list {
		fileGuids = append(fileGuids, value.Guid)
	}
	return
}
