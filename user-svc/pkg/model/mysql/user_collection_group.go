package mysql

import (
	"fmt"

	statv1 "ecodepost/pb/stat/v1"
	"github.com/ego-component/egorm"
	"gorm.io/gorm"
)

const (
	tableUserCollectionGroup = "user_collection_group"
)

type UserCollectionGroup struct {
	BaseModel
	Uid   int64  `gorm:"not null;index:idx_uid"`         // 用户ID
	Title string `gorm:"not null;"`                      // 标题
	Desc  string `gorm:"not null;"`                      // 说明
	Cnt   int64  `gorm:"not null;default:0;comment:收藏数"` // 收藏数
	// VisibleType commonv1.CollectionVisibleType `gorm:"not null;default:0;comment:可见类型"`             // 可见类型
}

func (UserCollectionGroup) TableName() string {
	return tableUserCollectionGroup
}

type UserCollectionGroups []*UserCollectionGroup

func (list UserCollectionGroups) ToPb() []*statv1.CollectionGroupInfo {
	output := make([]*statv1.CollectionGroupInfo, 0)
	for _, value := range list {
		output = append(output, value.ToPb())
	}
	return output
}

func (info *UserCollectionGroup) ToPb() *statv1.CollectionGroupInfo {
	return &statv1.CollectionGroupInfo{
		Id:        info.ID,
		Title:     info.Title,
		IsCollect: false,
		Cnt:       info.Cnt,
	}
}

// UserCollectionGroupCreate 创建一条记录
func UserCollectionGroupCreate(db *gorm.DB, data *UserCollectionGroup) (err error) {
	if err = db.Create(data).Error; err != nil {
		err = fmt.Errorf("UserCollectionGroupCreate fail, err: %w", err)
		return
	}
	return
}

// UserCollectionGroupList 创建一条记录
func UserCollectionGroupList(tx *gorm.DB, uid int64) (res UserCollectionGroups, err error) {
	resp := tx.Where("uid = ?", uid).Find(&res)
	return res, resp.Error
}

// UserCollectionGroupUpdateX 修改一条记录
func UserCollectionGroupUpdateX(db *gorm.DB, conds egorm.Conds, ups map[string]any) (err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = db.Table(tableUserCollectionGroup).Where(sql, binds...).Updates(ups).Error; err != nil {
		err = fmt.Errorf("UserCollectionGroupUpdate fail, err: %w", err)
		return
	}
	return
}
