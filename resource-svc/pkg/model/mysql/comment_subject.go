package mysql

import (
	commonv1 "ecodepost/pb/common/v1"
)

// CommentSubject 主题表
type CommentSubject struct {
	Id             int64            `gorm:"not null;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	BizGuid        string           `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:comment index guid"`
	BizType        commonv1.CMN_BIZ `gorm:"not null;index:idx_column,unique; comment:对象类型" json:"bizType"`
	CntComment     int32            `gorm:"not null;default:0;" json:"cntComment"`     // 评论总数
	CntRootComment int32            `gorm:"not null;default:0;" json:"cntRootComment"` // 根评论总数，一级评论总数
	Status         int8             `gorm:"not null;index:idx_column,unique;comment:0初始化、1正常，2删除等" json:"status"`
	Ctime          int64            `gorm:"comment:创建时间" json:"ctime"`
	Utime          int64            `gorm:"comment:更新时间" json:"utime"` // 更新时间
	Dtime          int64            `gorm:"comment:删除时间" json:"dtime"`
}

func (CommentSubject) TableName() string {
	return "comment_subject"
}
