package service

import (
	"context"
	"fmt"
	"time"

	"ecodepost/resource-svc/pkg/constx"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/dao"
	"ecodepost/resource-svc/pkg/model/mysql"

	commentv1 "ecodepost/pb/comment/v1"
	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"

	"github.com/ego-component/egorm"
	"gorm.io/gorm"
)

func CommentSubjectCreate(db *gorm.DB, data *mysql.CommentSubject) (subjectId int64, err error) {
	var info mysql.CommentSubject
	err = db.Select("id").Where("biz_guid = ? and biz_type = ?", data.BizGuid, int32(data.BizType)).Find(&info).Error
	if err != nil {
		return 0, fmt.Errorf("comment subject create failed1, err: %w", err)
	}

	if info.Id != 0 {
		return info.Id, nil
	}

	if err = db.Create(data).Error; err != nil {
		return 0, fmt.Errorf("comment subject create failed2, err: %w", err)
	}
	return data.Id, nil
}

func (*comment) Create(ctx context.Context, req *commentv1.CreateReq, commentSubjectId int64) (*commentv1.CreateRes, error) {
	// 更新总数
	err := dao.CommentSubject.UpdateExprById(invoker.Db.WithContext(ctx), commentSubjectId, "cnt_comment", 1)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Create error1").WithMetadata(map[string]string{"err": err.Error()})
	}
	// 更新根节点总数
	err = dao.CommentSubject.UpdateExprById(invoker.Db.WithContext(ctx), commentSubjectId, "cnt_root_comment", 1)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Create error2").WithMetadata(map[string]string{"err": err.Error()})
	}

	guid, err := Resource.GenerateGuid(ctx, commonv1.CMN_GUID_COMMENT, req.GetUid())
	if err != nil {
		return nil, fmt.Errorf("comment subject create generate guid failed1, err: %w", err)
	}

	nowTime := time.Now().Unix()
	ctime := nowTime
	utime := nowTime
	if req.GetCtime() > 0 {
		ctime = req.GetCtime()
	}
	if req.GetUtime() > 0 {
		utime = req.GetUtime()
	}

	cIndex := &mysql.CommentIndex{
		Guid:       guid,
		Uid:        req.Uid,
		SubjectId:  commentSubjectId,
		BizType:    req.GetBizType(),
		BizGuid:    req.GetBizGuid(),
		ActionType: req.GetActionType(),
		ActionGuid: req.GetActionGuid(),
		Status:     constx.AuditSuccess,
		Ctime:      ctime,
		Utime:      utime,
		Ip:         req.Ip,
	}

	err = invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 存index表
		if err = storeComment(ctx, tx, req.Content, cIndex); err != nil {
			return errcodev1.ErrDbError().WithMessage("Create error3").WithMetadata(map[string]string{"err": err.Error()})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &commentv1.CreateRes{CommentGuid: guid}, nil
}

func ReplyCommentByGuid(ctx context.Context, req *commentv1.CreateReq) (*commentv1.CreateRes, error) {
	guid, err := Resource.GenerateGuid(ctx, commonv1.CMN_GUID_COMMENT, req.Uid)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateCommentSubject Err1").WithMetadata(map[string]string{
			"err": err.Error(),
		})
	}

	err = invoker.CommentDb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 悲观锁，防止楼层混乱
		resp, err := dao.CommentIndex.InfoXLock(
			tx,
			egorm.Conds{"guid": req.GetCommentGuid()},
		)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("ReplyComment Err1").WithMetadata(map[string]string{
				"err": err.Error(),
			})
		}
		var rootCommentId string
		// 判断是否为根 comment id。说明是顶级id
		if resp.ReplyToGuid == "" {
			rootCommentId = resp.Guid
			// 不是顶级id，需要在找上一级
		} else {
			rootCommentId = resp.ReplyToRootGuid
		}

		// 更新该主题的评论总数
		err = dao.CommentSubject.UpdateExprById(tx, resp.SubjectId, "cnt_comment", 1)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("ReplyComment Err2").WithMetadata(map[string]string{
				"err": err.Error(),
			})
		}
		// 更新子comment总数
		err = dao.CommentIndex.UpdateExpr(tx, rootCommentId, "cnt_child_comment", 1)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("ReplyComment Err3").WithMetadata(map[string]string{
				"err": err.Error(),
			})
		}

		cIndex := &mysql.CommentIndex{
			Guid:            guid,
			SubjectId:       resp.SubjectId,
			Uid:             req.Uid,
			ReplyToUid:      resp.Uid,
			ReplyToGuid:     req.GetCommentGuid(),
			ReplyToRootGuid: rootCommentId,
			BizGuid:         req.GetBizGuid(),
			BizType:         req.BizType,
			Status:          constx.AuditSuccess,
			Ctime:           time.Now().Unix(),
			Utime:           time.Now().Unix(),
			ActionType:      req.GetActionType(),
			ActionGuid:      req.GetActionGuid(),
		}
		err = storeComment(ctx, tx, req.Content, cIndex)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("ReplyComment Err3").WithMetadata(map[string]string{
				"err": err.Error(),
			})
		}
		return nil
	})

	//	go func() {
	//		token := uuid.NewV4().String()
	//		invoker.CommentRedis.SetEX(ctx, "/comment/audit/"+token, req.GetCommentGuid(), 24*time.Hour)
	//		Dingtalk.SendAlarmToken(ctx, "回复评论", `
	//		* 用户uid:`+cast.ToString(req.Uid)+`
	//		* 主题id:`+cast.ToString(req.GetBizGuid())+`
	//		* 回复uid:`+cast.ToString(resp.Uid)+`
	//		* 回复评论id:`+cast.ToString(req.GetCommentGuid())+`
	//		* 内容:`+cast.ToString(req.Content)+`
	// `, token)
	//	}()
	return &commentv1.CreateRes{
		CommentGuid: guid,
	}, err
}
