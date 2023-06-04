package mysql

import (
	"context"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/user-svc/pkg/invoker"
)

type CountBehaviorLog struct {
	Id    int64             `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Ds    string            `gorm:"column:ds"`
	Fid   string            `gorm:"column:fid"`
	Biz   commonv1.CMN_BIZ  `gorm:"column:biz"`
	Act   commonv1.CNT_ACT  `gorm:"column:act"`
	Acti  commonv1.CNT_ACTI `gorm:"column:acti;type:varchar(64)"` // Act 操作名称
	Tid   string            `gorm:"column:tid"`
	Ctime int64             `gorm:"column:ctime"`
}

func (t CountBehaviorLog) TableName() string {
	return "count_behavior_log"
}

// behaviorDs 获取当月日期
func (t *CountBehaviorLog) behaviorDs() string {
	day := time.Now().Format("20060102")
	return day
}

// StoreBehaviorLog 记录行为日志
func (t *CountBehaviorLog) StoreBehaviorLog(ctx context.Context, info *CountBehaviorLog) (err error) {
	err = invoker.Db.WithContext(ctx).Model(CountBehaviorLog{
		Ds:    t.behaviorDs(),
		Fid:   info.Fid,
		Biz:   info.Biz,
		Act:   info.Act,
		Acti:  info.Acti,
		Tid:   info.Tid,
		Ctime: time.Now().Unix(),
	}).Error
	return err
}
