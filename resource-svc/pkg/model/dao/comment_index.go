package dao

import (
	"context"
	"fmt"
	"time"

	"ecodepost/resource-svc/pkg/constx"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"

	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type commentIndex struct {
	logger *elog.Component
	db     *gorm.DB
}

func InitCommentIndex(logger *elog.Component, db *gorm.DB) *commentIndex {
	return &commentIndex{
		logger: logger,
		db:     db,
	}
}

// Create 新增一条记
func (g *commentIndex) Create(db *gorm.DB, data *mysql.CommentIndex) (err error) {
	if err = db.Create(data).Error; err != nil {
		g.logger.Error("create commentIndex create error", zap.Error(err))
		return
	}
	return nil
}

func (g *commentIndex) UpdateExpr(db *gorm.DB, guid string, column string, delta int) (err error) {
	if err = db.Table("comment_index").Where("guid = ? ", guid).
		Updates(map[string]interface{}{column: gorm.Expr(column+"+ ?", delta), "utime": time.Now().Unix()}).Error; err != nil {
		g.logger.Error("comment_index update error", zap.Error(err))
		return
	}
	return
}

func CommentIndexDelete(ctx context.Context, uid int64, commentGuid string, bizGuid string, bizType commonv1.CMN_BIZ, deleteType commonv1.FILE_CMET_DEL) (commentIndexInfo mysql.CommentIndex, err error) {
	// id， uid，是用户侧删除数据
	// bizid  biztype  actiontype是后台删除数据
	switch deleteType {
	case commonv1.FILE_CMET_DEL_USER:
		sql, binds := egorm.BuildQuery(egorm.Conds{
			"guid": commentGuid,
			"uid":  uid,
		})
		err = invoker.Db.WithContext(ctx).Where(sql, binds...).Find(&commentIndexInfo).Error
		if err != nil {
			return
		}
	case commonv1.FILE_CMET_DEL_EXCELLENT:
		sql, binds := egorm.BuildQuery(egorm.Conds{
			"biz_guid":    bizGuid,
			"biz_type":    bizType,
			"action_type": commonv1.FILE_ACT_EXCELLENT,
		})
		err = invoker.Db.WithContext(ctx).Where(sql, binds...).Find(&commentIndexInfo).Error
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("no delete type")
		return
	}

	if commentIndexInfo.Id == 0 {
		return mysql.CommentIndex{}, fmt.Errorf("CommentIndexDelete not found id, %d", commentIndexInfo.Id)
	}

	// 故意不要事务的，否则提交到mysql，刚好查询的时候，数据不会变
	if err = invoker.Db.WithContext(ctx).Model(mysql.CommentIndex{}).Where("id = ?", commentIndexInfo.Id).Updates(map[string]interface{}{
		"dtime": time.Now().Unix(),
	}).Error; err != nil {
		return mysql.CommentIndex{}, fmt.Errorf("CommentIndexDeleteX failed,err :%w", err)
	}

	err = CommentSubject.UpdateExprById(invoker.Db.WithContext(ctx), commentIndexInfo.SubjectId, "cnt_comment", -1)
	if err != nil {
		return mysql.CommentIndex{}, fmt.Errorf("CommentIndexDeleteX failed2,err :%w", err)
	}
	// 根节点
	if commentIndexInfo.ReplyToGuid == "" {
		err = CommentSubject.UpdateExprById(invoker.Db.WithContext(ctx), commentIndexInfo.SubjectId, "cnt_root_comment", -1)
		if err != nil {
			return mysql.CommentIndex{}, fmt.Errorf("CommentIndexDeleteX failed2,err :%w", err)
		}
	}
	return
}

// InfoX Info的扩展方法，根据Cond查询单条记录
func (g *commentIndex) InfoXLock(db *gorm.DB, conds egorm.Conds) (resp mysql.CommentIndex, err error) {
	sql, binds := egorm.BuildQuery(conds)

	if err = db.Clauses(clause.Locking{Strength: "UPDATE"}).Table("comment_index").
		Where(sql, binds...).First(&resp).Error; err != nil {
		g.logger.Error("comment_index info error", zap.Error(err))
		return
	}
	return
}

func CommentIndexListByGuids(db *gorm.DB, guids []string) (resp []mysql.CommentIndexStoreRedis, err error) {
	err = db.Select("guid,biz_guid,biz_type,uid,subject_id,reply_to_uid,ctime,action_type,action_guid,reply_to_guid,reply_to_root_guid,cnt_child_comment,ip").Where("guid in (?)", guids).Find(&resp).Error
	return
}

func CommentIndexListByReplyToRootGuid(db *gorm.DB, commentGuid string) (output []string, total int64, hasMore bool, err error) {
	resp := make([]mysql.CommentIndexPage, 0)
	sql, binds := egorm.BuildQuery(egorm.Conds{
		"reply_to_root_guid": commentGuid,
		"dtime":              0,
		"status":             constx.AuditSuccess,
	})
	dbSQL := db.Table("comment_index").Where(sql, binds...)
	dbSQL.Count(&total)
	if err = dbSQL.Order("id desc").Limit(2).Find(&resp).Error; err != nil {
		invoker.Logger.Error("commentIndex info error", zap.Error(err))
		return
	}
	if total > int64(len(resp)) {
		hasMore = true
	}

	output = make([]string, 0)
	for _, guidInfo := range resp {
		output = append(output, guidInfo.Guid)
	}
	return
}

// CommentListPage 根据分页条件查询list
func CommentListPage(ctx context.Context, conds egorm.Conds, reqList *commonv1.Pagination) (output []string) {
	if reqList.PageSize == 0 || reqList.PageSize > 200 {
		reqList.PageSize = 20
	}
	if reqList.CurrentPage == 0 {
		reqList.CurrentPage = 1
	}
	sql, binds := egorm.BuildQuery(conds)

	db := invoker.Db.WithContext(ctx).Table("comment_index").Where(sql, binds...)
	respList := make([]mysql.CommentIndexPage, 0)
	db.Count(&reqList.Total)
	db.Select("guid").Order("id desc").Offset(int((reqList.CurrentPage - 1) * reqList.PageSize)).Limit(int(reqList.PageSize)).Find(&respList)
	output = make([]string, 0)
	for _, guidInfo := range respList {
		output = append(output, guidInfo.Guid)
	}
	return
}
