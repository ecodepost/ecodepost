package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	"github.com/ego-component/egorm"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// NotifyLetter 站内信
type NotifyLetter struct {
	BaseModel
	NotificationId int64               `gorm:"not null; index:idx_notification_id; comment:消息通知ID"`
	Type           commonv1.NTF_TYPE   `gorm:"not null; default:0; comment:消息类型"`                      // 冗余字段
	TargetId       string              `gorm:"not null; default:0; index:idx_target_id; comment:目标ID"` // 冗余字段
	Uid            int64               `gorm:"not null; index:idx_uid; comment:用户ID"`
	Status         commonv1.NTF_STATUS `gorm:"not null; comment:通知状态"`
	Link           string              `gorm:"type:varchar(255); comment:消息链接"`
	Meta           string              `gorm:"not null; type:text; comment:附属数据"`
	SenderId       string              `gorm:"-"` // 暂不使用
}

// TableName 数据表名
func (NotifyLetter) TableName() string {
	return "notify_letter"
}

type Letters []*NotifyLetter

func (us Letters) MapNotificationIds(fn func(u *NotifyLetter) int64) []int64 {
	res := make([]int64, 0)
	for _, u := range us {
		res = append(res, fn(u))
	}
	return lo.Uniq(res)
}

func (us Letters) Filter(fn func(u *NotifyLetter) bool) Letters {
	res := make(Letters, 0)
	for _, u := range us {
		if !fn(u) {
			continue
		}
		res = append(res, u)
	}
	return res
}

// LetterCreate 创建一条记录
func LetterCreate(db *gorm.DB, data *NotifyLetter) (err error) {
	if err = db.Create(data).Error; err != nil {
		err = fmt.Errorf("LetterCreate err:%w", err)
		return
	}
	return
}

// NotifyUpdate 根据主键更新一条记录
func NotifyUpdate(db *gorm.DB, paramId int64, ups map[string]any) (err error) {
	var sql = "`id`=?"
	var binds = []interface{}{paramId}
	if err = db.Model(NotifyLetter{}).Where(sql, binds...).Updates(ups).Error; err != nil {
		err = fmt.Errorf("NotifyUpdate err:%w", err)
		return
	}
	return
}

// NotifyUpdateX Update的扩展方法，根据Cond更新一条或多条记录
func NotifyUpdateX(db *gorm.DB, conds egorm.Conds, ups map[string]any) (err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = db.Model(NotifyLetter{}).Where(sql, binds...).Updates(ups).Error; err != nil {
		err = fmt.Errorf("NotifyUpdateX err:%w", err)
		return
	}
	return
}

// NotifyInfoX ...
func NotifyInfoX(db *gorm.DB, conds egorm.Conds) (resp NotifyLetter, err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = db.Model(NotifyLetter{}).Where(sql, binds...).Find(&resp).Error; err != nil {
		err = fmt.Errorf("NotifyInfoX err:%w", err)
		return
	}
	return
}

// NotifyCount 根据查询条件查询结果数量
func NotifyCount(db *gorm.DB, conds egorm.Conds) (count int64, err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = db.Model(NotifyLetter{}).Where(sql, binds...).Count(&count).Error; err != nil {
		err = fmt.Errorf("NotifyCount err:%w", err)
		return
	}
	return
}

// NotifyListPage 根据分页条件查询list
func NotifyListPage(db *gorm.DB, conds egorm.Conds, reqList *commonv1.Pagination) (respList Letters, err error) {
	respList = make(Letters, 0)
	if reqList.PageSize == 0 {
		reqList.PageSize = 20
	}
	if reqList.CurrentPage <= 0 {
		reqList.CurrentPage = 1
	}
	sql, binds := egorm.BuildQuery(conds)
	db = db.Model(NotifyLetter{}).Where(sql, binds...)
	db.Count(&reqList.Total)
	if reqList.Sort == "" {
		reqList.Sort = "utime desc"
	}
	err = db.Order(reqList.Sort).Offset(int(reqList.CurrentPage-1) * int(reqList.PageSize)).Limit(int(reqList.PageSize)).Find(&respList).Error
	if err != nil {
		err = fmt.Errorf("NotifyListPage err:%w", err)
	}

	return
}
