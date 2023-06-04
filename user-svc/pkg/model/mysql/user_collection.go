package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	statv1 "ecodepost/pb/stat/v1"
	"gorm.io/gorm"
	sdelete "gorm.io/plugin/soft_delete"
)

type UserCollection struct {
	ID    int64             `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Ctime int64             `gorm:"bigint;autoCreateTime;comment:创建时间" json:"ctime"`
	Dtime sdelete.DeletedAt `gorm:"bigint;comment:删除时间" json:"dtime"`
	// 在这三列上创建唯一索引，确保每个用户和特定资源只有一个关系
	Uid     int64            `gorm:"not null;unique_index:bid_bt_uid" json:"uid"` // 用户ID
	BizGuid string           `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT ''; unique_index:bid_bt_uid; comment:业务ID"`
	BizType commonv1.CMN_BIZ `gorm:"not null;unique_index:bid_bt_uid" json:"bizType"` // 类型id
	GroupId int64            `grom:"not null;default:0;comment:收藏分组"`
}

func (UserCollection) TableName() string {
	return "user_collection"
}

func (c UserCollection) ToPb() *statv1.CollectionInfo {
	return &statv1.CollectionInfo{
		Id:      c.ID,
		Uid:     c.Uid,
		BizGuid: c.BizGuid,
		BizType: c.BizType,
		// CollectionGroupIds: nil,
	}
}

type UserCollections []*UserCollection

func (list UserCollections) ToPb() []*statv1.CollectionInfo {
	output := make([]*statv1.CollectionInfo, 0)
	for _, value := range list {
		output = append(output, value.ToPb())
	}
	return output
}

// UserCollectionCreate 创建一条记录
func UserCollectionCreate(db *gorm.DB, data *UserCollection) (err error) {
	if err = db.Create(data).Error; err != nil {
		err = fmt.Errorf("UserCollectionCreate fail, err: %w", err)
		return
	}

	err = db.Model(UserCollectionGroup{}).Where("group_id = ?", data.GroupId).Update("cnt", gorm.Expr("cnt + ?", 1)).Error
	if err != nil {
		err = fmt.Errorf("UserCollectionCreate cnt fail, err: %w", err)
		return
	}
	return
}

// UserCollectionDelete 创建一条记录
func UserCollectionDelete(db *gorm.DB, id int64, uid int64) (err error) {
	var data UserCollection
	err = db.Where("id = ?", id).Find(&data).Error
	if err != nil {
		err = fmt.Errorf("UserCollectionDelete fail, err: %w", err)
		return
	}
	if id == 0 {
		err = fmt.Errorf("not exist")
		return
	}
	if data.Uid != uid {
		err = fmt.Errorf("no auth")
		return
	}
	if err = db.Where("id = ?", id).Delete(&UserCollection{}).Error; err != nil {
		err = fmt.Errorf("UserCollectionDelete fail2, err: %w", err)
		return
	}
	err = db.Model(UserCollectionGroup{}).Where("group_id = ?", data.GroupId).Update("cnt", gorm.Expr("cnt - ?", 1)).Error
	if err != nil {
		err = fmt.Errorf("UserCollectionCreate cnt fail, err: %w", err)
		return
	}
	return
}

// GetCollectionInfo 创建一条记录
func GetCollectionInfo(db *gorm.DB, uid int64, bizGuid string, bizType commonv1.CMN_BIZ, groupId int64) (list UserCollection, err error) {
	err = db.Select("id").Where("uid = ? and  biz_guid = ? and biz_type = ? and group_id = ?", uid, bizGuid, bizType.Number(), groupId).Find(&list).Error
	return
}

// CollectionCreateToGroup 创建一条记录到某个收藏夹
func CollectionCreateToGroup(db *gorm.DB, list UserCollections) (err error) {
	err = db.CreateInBatches(list, len(list)).Error
	if err != nil {
		err = fmt.Errorf("CollectionCreateToGroup fail, err: %w", err)
		return
	}

	for _, value := range list {
		err = db.Model(UserCollectionGroup{}).Where("id = ?", value.GroupId).Update("cnt", gorm.Expr("cnt + ?", 1)).Error
		if err != nil {
			err = fmt.Errorf("UserCollectionCreate cnt fail, err: %w", err)
			return
		}
	}
	return
}

// CollectionDeleteCount 从多个收藏夹移除一条收藏记录
func CollectionDeleteCount(db *gorm.DB, uid int64, bizGuid string, bizType commonv1.CMN_BIZ, groupIds []int64) (cnt int64, err error) {
	err = db.Model(UserCollection{}).Where("uid = ? and  biz_guid = ? and biz_type = ? and group_id in (?)", uid, bizGuid, bizType.Number(), groupIds).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("CollectionDeleteCount fail, err: %w", err)
		return
	}
	return
}

// CollectionCount 从多个收藏夹移除一条收藏记录
func CollectionCount(db *gorm.DB, uid int64, bizGuid string, bizType commonv1.CMN_BIZ) (cnt int64, err error) {
	err = db.Model(UserCollection{}).Where("uid = ? and  biz_guid = ? and biz_type = ? ", uid, bizGuid, bizType.Number()).Count(&cnt).Error
	if err != nil {
		err = fmt.Errorf("CollectionDeleteCount fail, err: %w", err)
		return
	}
	return
}

// CollectionDeleteFromGroup 从多个收藏夹移除一条收藏记录
func CollectionDeleteFromGroup(db *gorm.DB, uid int64, bizGuid string, bizType commonv1.CMN_BIZ, groupIds []int64) (err error) {
	err = db.Where("uid = ? and  biz_guid = ? and biz_type = ? and group_id in (?)", uid, bizGuid, bizType.Number(), groupIds).Delete(&UserCollection{}).Error
	if err != nil {
		err = fmt.Errorf("SpaceGroupMemberBatchDelete fail, err: %w", err)
		return
	}
	for _, groupId := range groupIds {
		err = db.Model(UserCollectionGroup{}).Where("id = ?", groupId).Update("cnt", gorm.Expr("cnt - ?", 1)).Error
		if err != nil {
			err = fmt.Errorf("SpaceGroupMemberBatchDelete cnt fail, err: %w", err)
			return
		}
	}

	return
}

// CollectionGroupDelete 删除一个收藏夹和其中收藏品
func CollectionGroupDelete(db *gorm.DB, uid int64, groupId int64) (err error) {
	if uid == 0 || groupId == 0 {
		return fmt.Errorf("uid or cgid can't be 0")
	}
	// 删除收藏关联表记录
	if err := db.Where("group_id = ? and uid = ?", groupId, uid).Delete(&UserCollection{}).Error; err != nil {
		return fmt.Errorf("find UserCollection fail, %w", err)
	}
	// 删除收藏分组记录
	if err := db.Where("id = ? and uid = ?", groupId, uid).Delete(&UserCollectionGroup{}).Error; err != nil {
		return fmt.Errorf("find UserCollectionGroup fail, %w", err)
	}

	return nil
}

// CollectionListByGroup 查询指定收藏夹下所有搜藏列表
func CollectionListByGroup(db *gorm.DB, uid int64, groupId int64, page *commonv1.Pagination) (res UserCollections, err error) {
	if uid == 0 || groupId == 0 {
		return nil, fmt.Errorf("uid or cgid can't be 0")
	}
	if page.PageSize == 0 || page.PageSize > 200 {
		page.PageSize = 200
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	query := db.Model(&UserCollection{}).Where("group_id = ? and uid = ?", groupId, uid)
	query.Count(&page.Total)
	err = query.Order("id desc ").Offset(int((page.CurrentPage - 1) * page.PageSize)).Limit(int(page.PageSize)).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CollectionListByFileGuids 查询指定收藏夹下所有搜藏列表
func CollectionListByFileGuids(db *gorm.DB, uid int64, fileGuids []string) (res UserCollections, err error) {
	if uid == 0 {
		return nil, fmt.Errorf("uid can't be 0")
	}

	err = db.Where("uid = ? and biz_guid in (?)", uid, fileGuids).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
