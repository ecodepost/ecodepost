package dao

import (
	"ecodepost/resource-svc/pkg/model/mysql"

	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type commentContent struct {
	logger *elog.Component
	db     *gorm.DB
}

func InitCommentContent(logger *elog.Component, db *gorm.DB) *commentContent {
	return &commentContent{
		logger: logger,
		db:     db,
	}
}

func (g *commentContent) Create(db *gorm.DB, data *mysql.CommentContent) (err error) {
	if err = db.Create(data).Error; err != nil {
		g.logger.Error("create commentContent create error", zap.Error(err))
		return
	}
	return nil
}

func (g *commentContent) Detail(db *gorm.DB, commentGuid string) (resp mysql.CommentContent, err error) {
	if err = db.Where("comment_guid = ?", commentGuid).First(&resp).Error; err != nil {
		g.logger.Error("comment_content info error", zap.Error(err))
		return
	}
	return
}
