package service

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/constx"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/dao"
	"ecodepost/resource-svc/pkg/model/mysql"

	commentv1 "ecodepost/pb/comment/v1"
	commonv1 "ecodepost/pb/common/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/ego-component/egorm"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

var CommentInfoNotFound = fmt.Errorf("CommentInfoNotFound")

type indexService struct{}

func (indexService) GetContentKeyByGuid(commentGuid string) string {
	return constx.CommentContentCache + commentGuid
}

func (i *indexService) DelIndexContentCacheGuid(ctx context.Context, commentGuid string) {
	invoker.Redis.Del(ctx, i.GetContentKeyByGuid(commentGuid))
}

func (i *indexService) GetIndexListByMySQLByGuid(ctx context.Context, subjectId int64, page *commonv1.Pagination) (indexGuids []string, respList map[string]*commentv1.CommentDetail, err error) {
	indexGuids = dao.CommentListPage(ctx, egorm.Conds{
		"subject_id":    subjectId,
		"reply_to_guid": "",
		"dtime":         0,
		"status":        constx.AuditSuccess,
	}, page)

	if len(indexGuids) == 0 {
		return
	}

	mInfo, err := i.GetCommentInfoByCacheByGuid(ctx, indexGuids)
	if err != nil {
		return
	}

	return indexGuids, mInfo, nil
}

func (i *indexService) GetChildIndexListByMySQLByGuid(ctx context.Context, commentGuid string, page *commonv1.Pagination) (indexGuids []string, respList map[string]*commentv1.CommentDetail, err error) {
	indexGuids = dao.CommentListPage(ctx, egorm.Conds{
		"reply_to_root_guid": commentGuid,
		"dtime":              0,
		"status":             constx.AuditSuccess,
	}, page)

	if len(indexGuids) == 0 {
		return
	}

	mInfo, err := i.GetCommentInfoByCacheByGuid(ctx, indexGuids)
	if err != nil {
		return
	}

	return indexGuids, mInfo, nil
}

func (i *indexService) GetIndexListByMysqlUserByGuid(ctx context.Context, uid int64, page *commonv1.Pagination) (indexGuids []string, respList map[string]*commentv1.CommentDetail, err error) {
	indexGuids = dao.CommentListPage(ctx, egorm.Conds{
		"uid":    uid,
		"dtime":  0,
		"status": constx.AuditSuccess,
	}, page)

	if len(indexGuids) == 0 {
		return
	}

	mInfo, err := i.GetCommentInfoByCacheByGuid(ctx, indexGuids)
	if err != nil {
		return
	}

	return indexGuids, mInfo, nil
}

func (i *indexService) GetSubIndexListByReplyToRootIdByGuid(ctx context.Context, db *gorm.DB, commentGuid string) (pbList []*commentv1.CommentDetail, total int64, hasMore bool, err error) {
	indexGuids, total, hasMore, err := dao.CommentIndexListByReplyToRootGuid(db, commentGuid)
	if err != nil {
		return nil, 0, false, err
	}
	mInfo, err := i.GetCommentInfoByCacheByGuid(ctx, indexGuids)
	if err != nil {
		return nil, 0, false, err
	}
	for _, indexGuid := range indexGuids {
		pbList = append(pbList, mInfo[indexGuid])
	}
	return pbList, total, hasMore, nil
}

func (i *indexService) GetCommentInfoByCacheByGuid(ctx context.Context, indexGuids []string) (respList map[string]*commentv1.CommentDetail, err error) {
	cacheKeyList := make([]string, 0)
	cacheKeyMapToId := make(map[string]string)

	respList = make(map[string]*commentv1.CommentDetail, 0)
	for _, indexGuid := range indexGuids {
		cacheKeyList = append(cacheKeyList, i.GetContentKeyByGuid(indexGuid))
		cacheKeyMapToId[i.GetContentKeyByGuid(indexGuid)] = indexGuid
	}
	commentList, err := invoker.Redis.MGet(ctx, cacheKeyList)
	if err != nil {
		return nil, err
	}

	needSqlQueryIds := make([]string, 0)
	for index, v := range commentList {
		// 说明没有缓存
		if v == nil {
			cacheKey := cacheKeyList[index]
			needSqlQueryIds = append(needSqlQueryIds, cacheKeyMapToId[cacheKey])
			continue
		}
		info := &commentv1.CommentDetail{}
		// redis获取的是string
		if err = proto.Unmarshal([]byte(v.(string)), info); err != nil {
			invoker.Logger.Info("unmarshal err")
			continue
		}
		respList[info.CommentGuid] = info
	}

	// 回源
	if len(needSqlQueryIds) > 0 {
		indexList, err := dao.CommentIndexListByGuids(invoker.Db.WithContext(ctx), needSqlQueryIds)
		if err != nil {
			return nil, err
		}
		infoList, err := SetMultiCacheContextByGuid(ctx, indexList)
		if err != nil {
			return nil, err
		}
		for _, v := range infoList {
			respList[v.GetCommentGuid()] = v
		}
	}

	for _, indexGuid := range indexGuids {
		go func(v string) {
			invoker.Redis.Expire(ctx, i.GetContentKeyByGuid(v), constx.CommentExpire)
		}(indexGuid)
	}

	// 将用户昵称，头像，放入评论中
	uids := make([]int64, 0)
	for _, value := range respList {
		uids = append(uids, value.Uid)
		if value.ReplyToUid != 0 {
			uids = append(uids, value.ReplyToUid)
		}
	}

	userMapResp, err := invoker.GrpcUser.Map(ctx, &userv1.MapReq{
		UidList: uids,
	})
	if err != nil {
		err = fmt.Errorf("user srv get info fail, err: %w", err)
		return
	}

	for key, value := range respList {
		value.UserNickname = userMapResp.GetUserMap()[value.Uid].GetNickname()
		value.UserAvatar = userMapResp.GetUserMap()[value.Uid].GetAvatar()
		value.ReplyNickname = userMapResp.GetUserMap()[value.ReplyToUid].GetNickname()
		value.ReplyAvatar = userMapResp.GetUserMap()[value.ReplyToUid].GetAvatar()
		respList[key] = value
	}
	return
}

func (i *indexService) GetContentDetailByCache(ctx context.Context, commentGuid string) (resp *commentv1.CommentDetail, err error) {
	res, err := invoker.Redis.Expire(ctx, i.GetContentKeyByGuid(commentGuid), constx.CommentExpire)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, CommentInfoNotFound
	}
	content, err := invoker.Redis.Get(ctx, i.GetContentKeyByGuid(commentGuid))
	if err != nil {
		return nil, err
	}
	resp = &commentv1.CommentDetail{}
	if err = proto.Unmarshal([]byte(content), resp); err != nil {
		return nil, err
	}
	return
}

func (i *indexService) SetContentDetailCache(ctx context.Context, commentGuid string, commentIndexRedis *mysql.CommentIndexStoreRedis) (resp *commentv1.CommentDetail, err error) {
	userReply, err := invoker.GrpcUser.List(ctx, &userv1.ListReq{
		UidList: []int64{commentIndexRedis.Uid, commentIndexRedis.ReplyToUid},
		// ReplyType: 1,
	})
	if err != nil {
		return nil, err
	}

	cContent, err := dao.CommentContent.Detail(invoker.Db.WithContext(ctx), commentGuid)
	if err != nil {
		return nil, err
	}
	commentInfo := commentIndexRedis.ToPBDetail(userReply, cContent.Content)
	go func() {
		i.setContentNxByGuid(ctx, commentGuid, commentInfo)
	}()
	return commentInfo, nil
}

func (i *indexService) setContentNxByGuid(ctx context.Context, commentGuid string, msg *commentv1.CommentDetail) (err error) {
	var msgBytes []byte
	if msgBytes, err = proto.Marshal(msg); err != nil {
		return
	}
	invoker.Redis.SetEX(ctx, i.GetContentKeyByGuid(commentGuid), msgBytes, constx.CommentContentExpire)
	return nil
}
