package profile

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	articlev1 "ecodepost/pb/article/v1"
	commonv1 "ecodepost/pb/common/v1"
	countv1 "ecodepost/pb/count/v1"
	questionv1 "ecodepost/pb/question/v1"
	userv1 "ecodepost/pb/user/v1"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type UserTotalRes struct {
	Uid          int64  `json:"uid"`
	Nickname     string `json:"nickname"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	RegisterTime int64  `json:"registerTime"`
	FollowingCnt int64  `json:"followingCnt"`
	FollowerCnt  int64  `json:"followerCnt"`
	HasFollowed  bool   `json:"hasFollowed"`
}

// UserTotal 用户统计数据
func UserTotal(c *bffcore.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSONE(1, "user name can't be empty", nil)
		return
	}

	u, err := invoker.GrpcUser.ProfileInfo(c.Ctx(), &userv1.ProfileInfoReq{Name: name})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// TODO 提供更原子化更专一的rpc接口
	myFollowersRes, err := invoker.GrpcCount.GetFidsTdetailByTid(c.Ctx(), &countv1.GetFidsTdetailByTidReq{
		Fid: cast.ToString(c.User().Uid),
		Tid: cast.ToString(u.Uid),
		Biz: commonv1.CMN_BIZ_USER,
		Act: commonv1.CNT_ACT_FOLLOW,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	myFollowingRes, err := invoker.GrpcCount.DBGetTdetailsByFid(c.Ctx(), &countv1.DBGetTdetailsByFidReq{
		Fid:        cast.ToString(u.Uid),
		Biz:        commonv1.CMN_BIZ_USER,
		Act:        commonv1.CNT_ACT_FOLLOW,
		Pagination: &commonv1.Pagination{CurrentPage: 1, PageSize: 1},
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(UserTotalRes{
		Uid:          u.Uid,
		Nickname:     u.Nickname,
		Name:         u.Name,
		Avatar:       u.Avatar,
		RegisterTime: u.RegisterTime,
		FollowingCnt: myFollowingRes.GetPagination().GetTotal(),
		FollowerCnt:  myFollowersRes.Num,
		HasFollowed:  myFollowersRes.Status == 1,
	})
}

type ArticlesListReq struct {
	bffcore.Pagination
}

func ArticlesList(c *bffcore.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSONE(1, "user name can't be empty", nil)
		return
	}

	var req ArticlesListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误: "+err.Error(), err)
		return
	}
	u, err := invoker.GrpcUser.Info(c.Ctx(), &userv1.InfoReq{Name: name})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	res, err := invoker.GrpcArticle.PublicListByUserCreated(c.Ctx(), &articlev1.PublicListByUserCreatedReq{
		Uid:        c.Uid(),
		CreatedUid: u.GetUser().Uid,
		Pagination: req.Pagination.ToPb(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(res.List, res.GetPagination())
}

type QAListReq struct {
	bffcore.Pagination
}

func QAList(c *bffcore.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSONE(1, "user name can't be empty", nil)
		return
	}

	var req QAListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误: "+err.Error(), err)
		return
	}
	u, err := invoker.GrpcUser.Info(c.Ctx(), &userv1.InfoReq{Name: name})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	res, err := invoker.GrpcQuestion.PublicListByUserCreated(c.Ctx(), &questionv1.PublicListByUserCreatedReq{
		Uid:        c.Uid(),
		CreatedUid: u.GetUser().Uid,
		Pagination: req.Pagination.ToPb(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(res.List, res.GetPagination())
}

type FollowersListReq struct {
	bffcore.Pagination
}

type FollowerItem struct {
	Nickname     string `json:"nickname"`     // 昵称
	Email        string `json:"email"`        // 邮箱
	Avatar       string `json:"avatar"`       // 头像
	Uid          int64  `json:"uid"`          // 用户UID
	Identify     int32  `json:"identify"`     // identity
	Name         string `json:"name"`         // 用户名称
	FollowersCnt int64  `json:"followersCnt"` // 被多少人关注
	FollowingCnt int64  `json:"followingCnt"` // 关注的人总数
	HasFollowed  bool   `json:"hasFollowed"`  // 当前用户是否关注
}

// FollowersList 查找指定用户的followers列表
// 对每个follower还继续往下查新其自身的followersCnt和followingCnt
func FollowersList(c *bffcore.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSONE(1, "user name can't be empty", nil)
		return
	}
	u, err := invoker.GrpcUser.Info(c.Ctx(), &userv1.InfoReq{Name: name})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	var req FollowersListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误: "+err.Error(), err)
		return
	}

	// 1.查找指定用户的Follower列表
	res, err := invoker.GrpcCount.DBGetFidsTdetailByTid(c.Ctx(), &countv1.DBGetFidsTdetailByTidReq{
		Tid:        cast.ToString(u.GetUser().Uid), // 指定name用户
		Biz:        commonv1.CMN_BIZ_USER,
		Act:        commonv1.CNT_ACT_FOLLOW,
		Pagination: req.Pagination.ToPb(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	// 如果没有followers，则直接返回结果
	followerItems := make([]FollowerItem, 0, len(res.Fids))
	if len(res.Fids) == 0 {
		c.JSONListPage(followerItems, res.Pagination)
		return
	}

	// 如果有followers，则继续查询所有followers详情
	// 2.查询指定用户的每个follower的FollowersCnt
	followersDetailRes, err := invoker.GrpcCount.GetTdetailsByTids(c.Ctx(), &countv1.GetTdetailsByTidsReq{
		Fid:  cast.ToString(c.Uid()),
		Tids: res.Fids,
		Biz:  commonv1.CMN_BIZ_USER,
		Act:  commonv1.CNT_ACT_FOLLOW,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// 3.查询指定用户的每个follower的FollowersCnt
	followingCntRes, err := invoker.GrpcCount.GetTnumByFids(c.Ctx(), &countv1.GetTnumByFidsReq{
		Fids: res.Fids,
		Biz:  commonv1.CMN_BIZ_USER,
		Act:  commonv1.CNT_ACT_FOLLOW,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	// 通过fuids查询其UserInfo
	fuids := lo.Map(res.Fids, func(v string, i int) int64 { return cast.ToInt64(v) })
	res2, err := invoker.GrpcUser.List(c.Ctx(), &userv1.ListReq{UidList: fuids})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	// 拼接item
	for _, v := range res2.UserList {
		followerDetail := followersDetailRes.Map[cast.ToString(v.Uid)]
		followerItems = append(followerItems, FollowerItem{
			Nickname:     v.Nickname,
			Email:        v.Email,
			Avatar:       v.Avatar,
			Uid:          v.Uid,
			Identify:     v.Identify,
			Name:         v.Name,
			FollowersCnt: followerDetail.Num,
			FollowingCnt: followingCntRes.Map[cast.ToString(v.Uid)],
			HasFollowed:  followerDetail.Status == 1,
		})
	}
	c.JSONListPage(followerItems, res.Pagination)
}

type FollowingListReq struct {
	bffcore.Pagination
}

type FollowingItem struct {
	Nickname     string `json:"nickname"`     // 昵称
	Email        string `json:"email"`        // 邮箱
	Avatar       string `json:"avatar"`       // 头像
	Uid          int64  `json:"uid"`          // 用户UID
	Identify     int32  `json:"identify"`     // identity
	Name         string `json:"name"`         // 用户名称
	FollowersCnt int64  `json:"followersCnt"` // 被多少人关注
	FollowingCnt int64  `json:"followingCnt"` // 关注的人总数
	HasFollowed  bool   `json:"hasFollowed"`  // 当前用户是否关注
}

// FollowingList 查找指定用户的following列表
// 对每个following还继续往下查新其自身的followersCnt和followingCnt
func FollowingList(c *bffcore.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSONE(1, "user name can't be empty", nil)
		return
	}
	u, err := invoker.GrpcUser.Info(c.Ctx(), &userv1.InfoReq{Name: name})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	var req FollowingListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误: "+err.Error(), err)
		return
	}

	// 1. 查询指定用户的Following列表
	res, err := invoker.GrpcCount.DBGetTdetailsByFid(c.Ctx(), &countv1.DBGetTdetailsByFidReq{
		Fid:        cast.ToString(u.GetUser().Uid),
		Biz:        commonv1.CMN_BIZ_USER,
		Act:        commonv1.CNT_ACT_FOLLOW,
		Pagination: req.Pagination.ToPb(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	var tuids = make([]int64, 0, len(res.List))
	var tuidsStr = make([]string, 0, len(res.List))
	for _, v := range res.List {
		// 只有Num>=1才表示目前关注中
		if v.Num >= 1 {
			tuids = append(tuids, cast.ToInt64(v.Tid))
			tuidsStr = append(tuidsStr, v.Tid)
		}
	}
	// 如果结果为空，则直接返回
	followingItems := make([]FollowingItem, 0, len(res.List))
	if len(tuidsStr) == 0 {
		c.JSONListPage(followingItems, res.Pagination)
		return
	}

	// 2.对指定用户的所有Following，查询其FollowersCnt
	followersCntRes, err := invoker.GrpcCount.GetTdetailsByTids(c.Ctx(), &countv1.GetTdetailsByTidsReq{
		Fid:  cast.ToString(c.Uid()),
		Tids: tuidsStr,
		Biz:  commonv1.CMN_BIZ_USER,
		Act:  commonv1.CNT_ACT_FOLLOW,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// 3.对指定用户所有Following，查询其FollowingCnt
	followingCntRes, err := invoker.GrpcCount.GetTnumByFids(c.Ctx(), &countv1.GetTnumByFidsReq{
		Fids: tuidsStr,
		Biz:  commonv1.CMN_BIZ_USER,
		Act:  commonv1.CNT_ACT_FOLLOW,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// 通过tuids查询其UserInfo
	ures, err := invoker.GrpcUser.List(c.Ctx(), &userv1.ListReq{UidList: tuids})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	for _, v := range ures.UserList {
		followerDetail := followersCntRes.Map[cast.ToString(v.Uid)]
		followingItems = append(followingItems, FollowingItem{
			Nickname:     v.Nickname,
			Email:        v.Email,
			Avatar:       v.Avatar,
			Uid:          v.Uid,
			Identify:     v.Identify,
			Name:         v.Name,
			FollowersCnt: followerDetail.Num,
			FollowingCnt: followingCntRes.Map[cast.ToString(v.Uid)],
			HasFollowed:  followerDetail.Status == 1,
		})
	}
	c.JSONListPage(followingItems, res.Pagination)
}
