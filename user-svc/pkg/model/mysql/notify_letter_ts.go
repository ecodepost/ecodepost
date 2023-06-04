package mysql

import (
	"fmt"

	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NotifyLetterUserTs 用户消息时间戳
type NotifyLetterUserTs struct {
	BaseModel
	Uid           int64 `gorm:"not null; uniqueIndex:idx_user_team; comment:用户ID"`
	ViewTimestamp int64 `gorm:"type:bigint(20) NOT NULL;comment:查看时间戳;"`
	ReadTimestamp int64 `gorm:"type:bigint(20) NOT NULL;comment:已读时间戳;"`
}

// TableName 数据表名
func (NotifyLetterUserTs) TableName() string {
	return "notify_letter_user_ts"
}

// NotifyUserTsCreate 创建一条记录
func NotifyUserTsCreate(db *gorm.DB, data *NotifyLetterUserTs) (err error) {
	if err = db.Create(data).Error; err != nil {
		err = fmt.Errorf("NotifyUserTsCreate err:%w", err)
		return
	}
	return
}

// NotifyUserTsUpdateRead 根据主键更新一条记录
func NotifyUserTsUpdateRead(db *gorm.DB, id int64, ts int64) (err error) {
	if err = db.Model(NotifyLetterUserTs{}).Where("id = ?", id).Update("read_timestamp", ts).Error; err != nil {
		err = fmt.Errorf("NotifyUserTsUpdateRead err:%w", err)
		return
	}
	return
}

// NotifyUserTsUpdateView 根据主键更新一条记录
func NotifyUserTsUpdateView(db *gorm.DB, id int64, ts int64) (err error) {
	if err = db.Model(NotifyLetterUserTs{}).Where("id = ?", id).Update("view_timestamp", ts).Error; err != nil {
		err = fmt.Errorf("NotifyUserTsUpdateView err:%w", err)
		return
	}
	return
}

// NotifyUserTsInfoX ...
func NotifyUserTsInfoX(db *gorm.DB, conds egorm.Conds) (resp NotifyLetterUserTs, err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = db.Model(NotifyLetterUserTs{}).Where(sql, binds...).Find(&resp).Error; err != nil {
		err = fmt.Errorf("NotifyUserTsInfoX err:%w", err)
		return
	}
	return
}

func ListNotifyUserTssByUID(tx *gorm.DB, userIDs []int64) (res map[int64]*NotifyLetterUserTs, err error) { // uid => timestamp
	var timestamps []NotifyLetterUserTs
	err = tx.Where("uid in (?)", userIDs).Find(&timestamps).Error
	if err != nil {
		elog.Error("ListNotifyUserTssByUID: find user timestamps failed", zap.Error(err))
		return
	}

	res = make(map[int64]*NotifyLetterUserTs)
	for _, timestamp := range timestamps {
		res[timestamp.Uid] = &timestamp
	}

	return
}
