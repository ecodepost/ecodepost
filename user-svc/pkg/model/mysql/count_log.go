package mysql

import (
	commonv1 "ecodepost/pb/common/v1"
	"github.com/ego-component/egorm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CountLog 日志计数信息
type CountLog struct {
	Id      int64            `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Biz     commonv1.CMN_BIZ `gorm:"not null;type:int(10);uniqueIndex:idx_biz_act_tid_fid"`      // 业务类型
	Act     commonv1.CNT_ACT `gorm:"not null;type:int(10);uniqueIndex:idx_biz_act_tid_fid"`      // 动作类型
	Tid     string           `gorm:"not null;type:varchar(255);uniqueIndex:idx_biz_act_tid_fid"` // 目标ID
	Fid     string           `gorm:"not null;type:varchar(255);uniqueIndex:idx_biz_act_tid_fid"` // 来源ID
	Num     int64            // 真实计数
	RealNum int64            // 真实计数
	BaseNum int64            // 真实计数
	Utime   int64            // 更新时间戳
	Ct      string           // 客户端
	Did     string           // 设备id
	Ip      string           // IP
}

func (t CountLog) TableName() string {
	return "count_log"
}

// UpdateRow 更新或新增一行
func (t CountLog) UpdateRow(db *gorm.DB, ts CountLog) error {
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tid"}, {Name: "event"}, {Name: "fid"}},
		DoUpdates: clause.AssignmentColumns([]string{"fid", "num", "real_num", "base_num", "utime"}),
	}).Create(&ts)

	return nil
}

// Update 更新一行或多行
func (t CountLog) Update(db *gorm.DB, conds egorm.Conds, ups map[string]any) (err error) {
	sql, binds := egorm.BuildQuery(conds)
	err = db.Model(CountLog{}).Where(sql, binds).Updates(ups).Error
	return nil
}

// Create 新增一行
func (t CountLog) Create(db *gorm.DB, s CountLog) (err error) {
	err = db.Table(t.TableName()).Create(&s).Error
	return nil
}

// Find 获取一行
func (t CountLog) Find(db *gorm.DB, conds egorm.Conds) (res *CountLog, err error) {
	sql, binds := egorm.BuildQuery(conds)
	err = db.Table(t.TableName()).Where(sql, binds...).Find(&res).Error
	return
}

// List 获取多行
func (t CountLog) List(db *gorm.DB, conds egorm.Conds, p *commonv1.Pagination) (res []*CountLog, err error) {
	sql, binds := egorm.BuildQuery(conds)
	query := db.Table(t.TableName()).Where(sql, binds...)
	if p == nil {
		err = query.Find(&res).Error
	} else {
		if p.PageSize == 0 || p.PageSize > 200 {
			p.PageSize = 20
		}
		if p.CurrentPage == 0 {
			p.CurrentPage = 1
		}
		query.Count(&p.Total)
		if p.Sort != "" {
			query = query.Order(p.Sort)
		}
		err = query.Offset(int((p.CurrentPage - 1) * p.PageSize)).Limit(int(p.PageSize)).Find(&res).Error
	}
	return
}
