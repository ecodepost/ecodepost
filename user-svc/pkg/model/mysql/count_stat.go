package mysql

import (
	"context"

	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/user-svc/pkg/invoker"
	"github.com/ego-component/egorm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NOTICE: 外部类型，INSERT/UPDATE 无影响，SELECT 最好转成基础类型
// 比如 commonv1.CMN_BIZ 和 commonv1.CNT_ACT

// CountStat 目标计数总值
type CountStat struct {
	Id      int64            `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Biz     commonv1.CMN_BIZ `gorm:"not null;type:int(10);uniqueIndex:idx_biz_act_tid"`      // 业务类型
	Act     commonv1.CNT_ACT `gorm:"not null;type:int(10);uniqueIndex:idx_biz_act_tid"`      // 动作类型
	Tid     string           `gorm:"not null;type:varchar(255);uniqueIndex:idx_biz_act_tid"` // 目标ID
	Num     int64            // 前台计数值
	RealNum int64            // 真实计数
	BaseNum int64            // 真实计数
	Utime   int64            // 更新时间戳
}

func (t CountStat) TableName() string {
	return "count_stat"
}

// InsertOrUpdate 更新或新增一行
func (t CountStat) InsertOrUpdate(db *gorm.DB, ts CountStat) (int64, error) {
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "biz"}, {Name: "act"}, {Name: "tid"}},
		DoUpdates: clause.AssignmentColumns([]string{"num", "real_num", "base_num", "utime"}),
	}).Create(&ts)
	return ts.RealNum, nil
}

// Update 更新一行或多行
func (t CountStat) Update(ctx context.Context, conds egorm.Conds, ups Ups) (err error) {
	sql, binds := egorm.BuildQuery(conds)
	err = invoker.Db.WithContext(ctx).Table(t.TableName()).Where(sql, binds).Updates(ups).Error
	return nil
}

// Create 新增一行
func (t CountStat) Create(ctx context.Context, s *CountStat) (err error) {
	err = invoker.Db.WithContext(ctx).Table(t.TableName()).Create(&s).Error
	return nil
}

// Find 获取一行
func (t CountStat) Find(ctx context.Context, conds egorm.Conds, s CountStat) (res *CountStat, err error) {
	sql, binds := egorm.BuildQuery(conds)
	err = invoker.Db.WithContext(ctx).Table(t.TableName()).Where(sql, binds...).Find(&res).Error
	return
}

// List 获取多行
func (t CountStat) List(db *gorm.DB, conds egorm.Conds, p *commonv1.Pagination) (res []*CountStat, err error) {
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

func (t CountStat) ListX(db *gorm.DB, conds egorm.Conds, p *commonv1.Pagination, exSql string, exBinds []any) (res []*CountStat, err error) {
	sql, binds := egorm.BuildQuery(conds)
	sql = sql + " AND " + exSql
	binds = append(binds, exBinds...)
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
