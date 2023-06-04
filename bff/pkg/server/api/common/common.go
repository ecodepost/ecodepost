package common

import (
	"context"
	"fmt"
	"sync"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	articlev1 "ecodepost/pb/article/v1"
	commonv1 "ecodepost/pb/common/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/spf13/cast"
)

type ListBizItemRes map[commonv1.CMN_BIZ]map[string]any

type BizItem struct {
	Type commonv1.CMN_BIZ `json:"type" form:"type" binding:"required" label:"业务类型"`
	Guid string           `json:"guid" form:"guid" binding:"required" label:"业务Guid"`
}

// ListBizItem 根据业务类型和业务ID列表，查询指定业务详细信息
func ListBizItem(c context.Context, uid int64, items []BizItem) (ListBizItemRes, error) {
	// 参数校验
	cnt := len(items)
	if cnt >= 100 {
		return nil, fmt.Errorf("too many items")
	}
	// 构造map<commonv1.CMN_BIZ, []guid>
	bizM := make(map[commonv1.CMN_BIZ][]string)
	for _, b := range items {
		_, ok := bizM[b.Type]
		if !ok {
			bizM[b.Type] = make([]string, 0)
		}
		bizM[b.Type] = append(bizM[b.Type], b.Guid)
	}
	// 根据bizType数据，起多个goroutine并发查询
	ch := make(chan BizGuidsRes, 10)
	wg := sync.WaitGroup{}
	for bizType, bizGuids := range bizM {
		wg.Add(1)
		go func(bizType commonv1.CMN_BIZ, guids []string) {
			itemMap, err := GetOneBiz(c, uid, bizType, guids)
			ch <- BizGuidsRes{bizType, itemMap, err}
			wg.Done()
		}(bizType, bizGuids)
	}

	wg.Wait()
	close(ch)

	res := make(ListBizItemRes, 0)
	for v := range ch {
		// 只要一个查询失败，则全部失败
		if v.Err != nil {
			return nil, fmt.Errorf("ListBizItem fail, err:%w", v.Err)
		}
		res[v.BizType] = v.Res
	}
	return res, nil
}

type BizGuidsRes struct {
	BizType commonv1.CMN_BIZ // 业务类型
	Res     map[string]any   // 响应结果
	Err     error            // 错误
}

// GetOneBiz 根据一个bizType查询多个guids列表
func GetOneBiz(ctx context.Context, uid int64, bizType commonv1.CMN_BIZ, bizGuids []string) (map[string]any, error) {
	switch bizType {
	case commonv1.CMN_BIZ_ARTICLE, commonv1.CMN_BIZ_ANSWER, commonv1.CMN_BIZ_QUESTION, commonv1.CMN_BIZ_COLUMN:
		res, err := invoker.GrpcArticle.ListDocumentByGuids(ctx, &articlev1.ListDocumentByGuidsReq{
			Uid:        uid,
			Guids:      bizGuids,
			Pagination: &commonv1.Pagination{CurrentPage: 1, PageSize: 10},
		})
		if err != nil {
			return nil, err
		}
		ret := make(map[string]any)
		for _, v := range res.List {
			ret[v.Guid] = v
		}
		return ret, nil
	case commonv1.CMN_BIZ_COMMUNITY:
	case commonv1.CMN_BIZ_SPACE:
	case commonv1.CMN_BIZ_USER:
		uids := make([]int64, 0, len(bizGuids))
		for _, v := range bizGuids {
			uids = append(uids, cast.ToInt64(v))
		}
		res, err := invoker.GrpcUser.List(ctx, &userv1.ListReq{UidList: uids})
		if err != nil {
			return nil, err
		}
		ret := make(map[string]any)
		for _, v := range res.UserList {
			ret[cast.ToString(v.Uid)] = v
		}
		return ret, nil
	default:
	}
	return nil, fmt.Errorf("unsupported bizType, %s", bizType.String())
}

func UserActOnBiz(c *bffcore.Context, biz commonv1.CMN_BIZ, tid string, act commonv1.CNT_ACT, acti commonv1.CNT_ACTI) {
	c.InjectGuid(tid)
	c.JSONOK()
}
