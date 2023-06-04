package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	commonv1 "ecodepost/pb/common/v1"

	"github.com/ego-component/egorm"
	"gorm.io/gorm"
	sdelete "gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	ID    int64             `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Ctime int64             `gorm:"bigint(20);autoCreateTime;comment:创建时间" json:"ctime"`
	Utime int64             `gorm:"bigint(20);autoUpdateTime;comment:更新时间" json:"utime"`
	Dtime sdelete.DeletedAt `gorm:"bigint(20);comment:删除时间" json:"dtime"`
}

type Ups = map[string]any

type I64s []int64

func (t I64s) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func (t *I64s) Scan(input interface{}) error {
	if len(input.([]byte)) == 0 {
		return json.Unmarshal([]byte("[]"), t)
	}
	if err := json.Unmarshal(input.([]byte), t); err != nil {
		return json.Unmarshal([]byte("[]"), t)
	}
	return json.Unmarshal(input.([]byte), t)
}

func Count(tx *gorm.DB, info interface{}, sql string, args ...interface{}) (cnt int64) {
	tx.Model(info).Where(sql, args...).Count(&cnt)
	return
}

func List(db *gorm.DB, table interface{}, conds egorm.Conds, list interface{}) (err error) {
	sql, binds := egorm.BuildQuery(conds)
	query := db.Model(table).Where(sql, binds...)
	err = query.Find(list).Error
	return
}

func ListPage(db *gorm.DB, table interface{}, conds egorm.Conds, reqList *commonv1.Pagination, respList interface{}) (err error) {
	if reqList.PageSize == 0 || reqList.PageSize > 200 {
		reqList.PageSize = 20
	}
	if reqList.CurrentPage == 0 {
		reqList.CurrentPage = 1
	}
	sql, binds := egorm.BuildQuery(conds)

	query := db.Model(table).Where(sql, binds...)
	query.Count(&reqList.Total)
	if reqList.Sort != "" {
		query = query.Order(reqList.Sort)
	}
	err = query.Offset(int((reqList.CurrentPage - 1) * reqList.PageSize)).Limit(int(reqList.PageSize)).Find(respList).Error
	return
}

func Find(tx *gorm.DB, field string, info interface{}, sql string, args ...interface{}) error {
	return tx.Model(info).Select(field).Where(sql, args...).Find(info).Error
}

func Create(tx *gorm.DB, create interface{}) error {
	return tx.Create(create).Error
}

func Update(tx *gorm.DB, model interface{}, id int64, ups map[string]interface{}) error {
	return tx.Model(model).Where("id = ?", id).Updates(ups).Error
}

func UpdateByUid(tx *gorm.DB, model interface{}, id int64, uid int64, ups map[string]interface{}) error {
	return tx.Model(model).Where("id = ? and uid = ?", id, uid).Updates(ups).Error
}

func UpdateByUidAndGuid(tx *gorm.DB, model interface{}, guid string, uid int64, ups map[string]interface{}) error {
	return tx.Model(model).Where("guid = ? and uid = ?", guid, uid).Updates(ups).Error
}

func Delete(tx *gorm.DB, model interface{}, id int64) error {
	return tx.Model(model).Where("id = ?", id).Updates(map[string]interface{}{
		"dtime": time.Now().Unix(),
	}).Error
}

func DeleteByUid(tx *gorm.DB, model interface{}, id int64, uid int64) error {
	return tx.Model(model).Where("id = ?", id).Updates(map[string]interface{}{
		"dtime": time.Now().Unix(),
	}).Error
}

func DeleteByUidAndGuid(tx *gorm.DB, model interface{}, guid string, uid int64) error {
	return tx.Model(model).Where("guid = ?", guid).Updates(map[string]interface{}{
		"dtime": time.Now().Unix(),
	}).Error
}
