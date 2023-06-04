package mysql

import (
	"fmt"
	"time"

	commonv1 "ecodepost/pb/common/v1"

	"gorm.io/gorm"
)

// Edge 节点关系
// Edge用于存储file之间的关系
type Edge struct {
	Id         int64              `json:"id" gorm:"not null;primary_key;auto_increment"`
	SpaceGuid  string             `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:空间GUID"`
	Ctime      int64              `gorm:"not null; comment:创建时间;" json:"ctime"`
	Utime      int64              `gorm:"not null; comment:更新时间;" json:"utime"`
	FileType   commonv1.FILE_TYPE `gorm:"not null; default:0; comment:节点File类型"`
	ParentType commonv1.FILE_TYPE `gorm:"not null; default:0; comment:父节点File类型"`
	FileGuid   string             `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:节点"`
	ParentGuid string             `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:上级节点"`
	Sort       int64              `gorm:"sort; not null; default:0; comment:排序号"`
}

// TableName 数据表名
func (Edge) TableName() string {
	return "edge"
}

type Edges []*Edge

func (list Edges) FindByGuid(guid string) *Edge {
	return list.Find(func(e *Edge) bool {
		return e.FileGuid == guid
	})
}

func (list Edges) Find(fn func(e *Edge) bool) *Edge {
	for _, spaceInfo := range list {
		if fn(spaceInfo) {
			return spaceInfo
		}
	}
	return nil
}

// EdgeCreate 创建一条记录
func EdgeCreate(db *gorm.DB, data *Edge) (err error) {
	if err = db.Create(data).Error; err != nil {
		return fmt.Errorf("EdgeCreate failed,err: %w", err)
	}
	return
}

// EdgeInfoByFileGuidAndFileType Info的扩展方法，根据Cond查询单条记录
func EdgeInfoByFileGuidAndFileType(db *gorm.DB, cmtGuid string, fileGuid string, fileType commonv1.FILE_TYPE) (resp Edge, err error) {
	if err = db.Model(Edge{}).Where("cmt_guid =? and file_guid = ? and file_type = ?", cmtGuid, fileGuid, fileType.Number()).Find(&resp).Error; err != nil {
		return Edge{}, fmt.Errorf("EdgeInfoX failed,err: %w", err)
	}
	return
}

// GetNodeParentGuid 创建一条记录
func GetNodeParentGuid(db *gorm.DB, fileGuid string) (parentGuid string, err error) {
	var info Edge
	err = db.Select("parent_guid").Where("file_guid = ?", fileGuid).Find(&info).Error
	if err != nil {
		return "", fmt.Errorf("GetNodeParentGuid fail, err: %w", err)
	}
	parentGuid = info.ParentGuid
	return
}

// GetNodeSortInfoByIn 创建一条记录
func GetNodeSortInfoByIn(db *gorm.DB, guids []string) (list Edges, err error) {
	err = db.Select("file_guid,file_type,parent_type,sort").Where("file_guid in (?) ", guids).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("GetNodeSortInfoByIn fail, err: %w", err)
	}
	return
}

// UpdateEdgeParentGuid 修改parent_guid，事务由外部保证，必须保证oldParentGuid和newParentGuid一致
func UpdateEdgeParentGuid(db *gorm.DB, fileGuid string, newParentGuid string) (err error) {
	// 直接将oldParentGuid更新为newParentGuid，就不做软删除处理了
	if e := db.Model(Edge{}).Where(`file_guid = ?`, fileGuid).Updates(map[string]any{
		"parent_guid": newParentGuid,
		"utime":       time.Now().Unix(),
	}).Error; e != nil {
		return fmt.Errorf("UpdateEdgeParentGuid fail, %w", err)
	}
	return
}

// NodeDelete 文件放到回收站
func NodeDelete(db *gorm.DB, uid int64, guid string) (err error) {
	var sonEdgeCount int64
	if err = db.Model(Edge{}).Where("parent_guid = ?", guid).Count(&sonEdgeCount).Error; err != nil {
		return fmt.Errorf("NodeDelete son edge count fail, err: %w", err)
	}

	// 说明有子节点，那么先要把parent的parent数据给到他们
	if sonEdgeCount > 0 {
		var parentGuid string
		parentGuid, err = GetNodeParentGuid(db, guid)
		if err != nil {
			return fmt.Errorf("NodeDelete GetNodeParentGuid fail, err: %w", err)
		}
		err = db.Model(Edge{}).Where("parent_guid = ?", guid).Updates(map[string]any{"parent_guid": parentGuid}).Error
		if err != nil {
			return fmt.Errorf("NodeDelete update parent guid fail, err: %w", err)
		}
	}

	if err = db.Model(Edge{}).Where("file_guid = ?", guid).Delete(&Edge{}).Error; err != nil {
		return fmt.Errorf("NodeDelete fail, err: %w", err)
	}
	return
}
