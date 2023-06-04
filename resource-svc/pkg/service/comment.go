package service

import (
	"context"
	"errors"

	"ecodepost/resource-svc/pkg/constx"
	"ecodepost/resource-svc/pkg/model/dao"
	"ecodepost/resource-svc/pkg/model/mysql"

	commentv1 "ecodepost/pb/comment/v1"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type comment struct {
}

func storeComment(ctx context.Context, tx *gorm.DB, content string, cIndex *mysql.CommentIndex) error {
	// 存index表
	err := dao.CommentIndex.Create(tx, cIndex)
	if err != nil {
		return err
	}

	// 存context表
	err = dao.CommentContent.Create(tx, &mysql.CommentContent{
		CommentGuid: cIndex.Guid,
		Content:     content,
		Ctime:       cIndex.Ctime,
		Utime:       cIndex.Utime,
		Status:      constx.AuditSuccess, // 不审核了，这里，反正都是按照index来检索的
	})
	if err != nil {
		return err
	}
	Index.DelIndexContentCacheGuid(ctx, cIndex.ReplyToRootGuid)
	Index.DelIndexContentCacheGuid(ctx, cIndex.ReplyToGuid)
	return nil
}

func SetMultiCacheContextByGuid(ctx context.Context, req []mysql.CommentIndexStoreRedis) (replyList []*commentv1.CommentDetail, err error) {
	g, ctx := errgroup.WithContext(ctx)
	// var cacheKeyList []string = make([]string, len(req))
	replyList = make([]*commentv1.CommentDetail, 0)
	for _, value := range req {
		// 一定要注意这个copy
		v := value
		g.Go(func() error {
			res, e := Index.GetContentDetailByCache(ctx, v.Guid)
			if e != nil && !errors.Is(e, CommentInfoNotFound) {
				return e
			}

			if errors.Is(e, CommentInfoNotFound) {
				res, e = Index.SetContentDetailCache(ctx, v.Guid, &v)
				if e != nil {
					return e
				}
			}
			replyList = append(replyList, res)
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		return
	}
	return
}
