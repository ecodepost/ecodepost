package mysql

import (
	commonv1 "ecodepost/pb/common/v1"
)

// CountConfig 计数配置表
type CountConfig struct {
	Id         int64             `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Biz        commonv1.CMN_BIZ  `gorm:"column:biz"`  // 业务类型
	Act        commonv1.CNT_ACT  `gorm:"column:act"`  // 动作名称
	ActCommand commonv1.CNT_ACTI `gorm:"column:acti"` // 动作指令
}
