package dao

import (
	"context"
	"fmt"
	"time"

	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"

	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type commentSubject struct {
	logger *elog.Component
	db     *gorm.DB
}

func InitCommentSubject(logger *elog.Component, db *gorm.DB) *commentSubject {
	return &commentSubject{
		logger: logger,
		db:     db,
	}
}

func (g *commentSubject) UpdateExprById(db *gorm.DB, id int64, column string, delta int) (err error) {
	if err = db.Table("comment_subject").Where("id = ? ", id).
		Updates(map[string]interface{}{column: gorm.Expr(column+"+ ?", delta), "utime": time.Now().Unix()}).Error; err != nil {
		err = fmt.Errorf("commentSubject update expr fail, err: %w", err)
		return
	}
	return
}

func CommentSubjectInfoByBizGuidAndBizType(db *gorm.DB, bizGuid string, bizType commonv1.CMN_BIZ) (resp mysql.CommentSubject, err error) {
	sql, binds := egorm.BuildQuery(egorm.Conds{
		"biz_guid": bizGuid,
		"biz_type": int32(bizType),
	})

	if err = db.Table("comment_subject").Select("guid").Where(sql, binds...).Find(&resp).Error; err != nil {
		err = fmt.Errorf("CommentSubjectInfoByBizIdAndBizType failed, err: %w", err)
		return
	}
	return
}

// InfoByBizInfo Info的扩展方法，根据Cond查询单条记录
func (g *commentSubject) InfoByBizInfo(db *gorm.DB, bizGuid string, bizType commonv1.CMN_BIZ) (resp mysql.CommentSubject, err error) {
	// if err = db.Select("id,cnt_comment").Where("biz_guid = ? and biz_type = ?", bizGuid, int32(bizType)).Find(&resp).Error; err != nil {
	if err = db.Select("id,cnt_comment").Where("biz_guid = ?", bizGuid).Find(&resp).Error; err != nil {
		g.logger.Error("commentSubject info error", zap.Error(err))
		return
	}
	return
}

// InfoX Info的扩展方法，根据Cond查询单条记录
func (g *commentSubject) InfoX(c context.Context, conds egorm.Conds) (resp mysql.CommentSubject, err error) {
	sql, binds := egorm.BuildQuery(conds)
	if err = g.db.WithContext(c).Table("comment_subject").Where(sql, binds...).Find(&resp).Error; err != nil {
		g.logger.Error("commentSubject info error", zap.Error(err))
		return
	}
	return
}

// InfoX Info的扩展方法，根据Cond查询单条记录
func (g *commentSubject) InfoXLock(db *gorm.DB, conds egorm.Conds) (resp mysql.CommentSubject, err error) {
	sql, binds := egorm.BuildQuery(conds)

	if err = db.Clauses(clause.Locking{Strength: "UPDATE"}).Table("comment_subject").
		Where(sql, binds...).First(&resp).Error; err != nil {
		err = fmt.Errorf("commentSubject info error: %w", err)
		return
	}
	return
}

// GetCommentSubjectId Info的扩展方法，根据Cond查询单条记录
func GetCommentSubjectId(db *gorm.DB, bizGuid string, bizType commonv1.CMN_BIZ) (subjectId int64, err error) {
	var resp mysql.CommentSubject
	if err = db.Table("id").Select("guid").Where("biz_guid = ? and biz_type = ? and dtime = 0", bizGuid, bizType.Number()).Find(&resp).Error; err != nil {
		err = fmt.Errorf("GetCommentSubjectId fail, err: %w", err)
		return
	}
	subjectId = resp.Id
	return
}

// List 查询list，extra[0]为sorts
func (g *commentSubject) List(c context.Context, conds egorm.Conds, extra ...string) (resp []mysql.CommentSubject, err error) {
	sql, binds := egorm.BuildQuery(conds)

	sorts := ""
	if len(extra) >= 1 {
		sorts = extra[0]
	}
	if err = g.db.Table("comment_subject").Where(sql, binds...).Order(sorts).Find(&resp).Error; err != nil {
		g.logger.Error("commentSubject info error", zap.Error(err))
		return
	}
	return
}

// ListMap 查询map，map遍历的时候是无序的，所以指定sorts参数没有意义
func (g *commentSubject) ListMap(c context.Context, conds egorm.Conds) (resp map[int64]mysql.CommentSubject, err error) {
	sql, binds := egorm.BuildQuery(conds)

	mysqlSlice := make([]mysql.CommentSubject, 0)
	resp = make(map[int64]mysql.CommentSubject, 0)
	if err = g.db.Table("comment_subject").Where(sql, binds...).Find(&mysqlSlice).Error; err != nil {
		g.logger.Error("commentSubject info error", zap.Error(err))
		return
	}
	for _, value := range mysqlSlice {
		resp[value.Id] = value
	}
	return
}
