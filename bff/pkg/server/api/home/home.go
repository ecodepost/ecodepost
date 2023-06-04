package home

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	// activityv1 "ecodepost/pb/activity/v1"
	articlev1 "ecodepost/pb/article/v1"
	commonv1 "ecodepost/pb/common/v1"
	communityv1 "ecodepost/pb/community/v1"
	spacev1 "ecodepost/pb/space/v1"

	"github.com/samber/lo"
)

type PageRes struct {
	// 是否启用banner
	IsSetBanner bool `json:"isSetBanner"`
	// 启用首页banner
	BannerImg string `json:"bannerImg"`
	// banner文案
	BannerTitle string `json:"bannerTitle"`
	// banner的跳转链接
	BannerLink string `json:"bannerLink"`
	// 默认访问页面
	DefaultPage string `json:"defaultPage"`
	// ArticlePage
	ArticlePageList ArticlePageList `json:"articlePageList"`
	// ArticleHot
	ArticleHot ArticleHot `json:"articleHot"`
}

type ArticlePageList struct {
	// 列表
	List []*commonv1.FileShow `json:"list"`
	// 分页
	Pagination *commonv1.Pagination `json:"pagination"`
}

type ArticleHot struct {
	// 列表
	List []*commonv1.FileShow `json:"list"`
}

func Page(c *bffcore.Context) {
	homeOption, err := invoker.GrpcCommunity.GetHomeOption(c.Ctx(), &communityv1.GetHomeOptionReq{
		Uid: c.Uid(),
	})
	if err != nil {
		c.JSONE(1, "error", err)
		return
	}

	resp := PageRes{
		IsSetBanner: homeOption.GetIsSetBanner(),
		BannerImg:   homeOption.GetBannerImg(),
		BannerTitle: homeOption.GetBannerTitle(),
		BannerLink:  homeOption.GetBannerLink(),
		DefaultPage: "",
	}
	// todo 新注册用户后面在处理
	if c.Uid() == 0 {
		resp.DefaultPage = homeOption.DefaultPageByNotLogin
	} else {
		resp.DefaultPage = homeOption.DefaultPageByLogin
	}

	var sort commonv1.CMN_FILE_SORT
	var pointerUid *int64
	// 如果未登录
	if c.Uid() == 0 {
		sort = homeOption.ArticleSortByNotLogin
	} else {
		sort = homeOption.ArticleSortByLogin
		uid := c.Uid()
		pointerUid = &uid
	}

	articlePage, err := invoker.GrpcArticle.HomeArticlePageList(c.Ctx(), &articlev1.HomeArticlePageListReq{
		Uid: pointerUid,
		Pagination: &commonv1.Pagination{
			CurrentPage: 0,
		},
		Sort: sort,
	})
	if err != nil {
		c.JSONE(1, "error", err)
		return
	}
	resp.ArticlePageList = ArticlePageList{
		List:       articlePage.GetList(),
		Pagination: articlePage.GetPagination(),
	}

	articleHot, err := invoker.GrpcArticle.HomeArticleHotList(c.Ctx(), &articlev1.HomeArticleHotListReq{
		Uid:        pointerUid,
		Limit:      homeOption.ArticleHotShowSum,
		LatestTime: homeOption.ArticleHotShowWithLatest,
	})
	// 列表
	resp.ArticleHot = ArticleHot{
		List: articleHot.GetList(),
	}

	c.JSONOK(resp)
}

type ListReq struct {
	CurrentPage int32 `form:"currentPage"` // 当前页数
}

func Files(c *bffcore.Context) {
	var req ListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	homeOption, err := invoker.GrpcCommunity.GetHomeOption(c.Ctx(), &communityv1.GetHomeOptionReq{
		Uid: c.Uid(),
	})
	if err != nil {
		c.JSONE(1, "error", err)
		return
	}

	var sort commonv1.CMN_FILE_SORT
	var pointerUid *int64
	// 如果未登录
	if c.Uid() == 0 {
		sort = homeOption.ArticleSortByNotLogin
	} else {
		sort = homeOption.ArticleSortByLogin
		uid := c.Uid()
		pointerUid = &uid
	}

	articlePage, err := invoker.GrpcArticle.HomeArticlePageList(c.Ctx(), &articlev1.HomeArticlePageListReq{
		Uid: pointerUid,
		Pagination: &commonv1.Pagination{
			CurrentPage: req.CurrentPage,
		},
		Sort: sort,
	})
	if err != nil {
		c.JSONE(1, "error", err)
		return
	}

	c.JSONListPage(articlePage.List, articlePage.GetPagination())
}

type GetRes struct {
	// 是否启用首页
	IsSetHome bool `json:"isSetHome,omitempty"`
	// 是否启用banner
	IsSetBanner bool `json:"isSetBanner,omitempty"`
	// 登录用户推荐内容排序规则
	ArticleSortByLogin commonv1.CMN_FILE_SORT `json:"articleSortByLogin,omitempty"`
	// 未登录用户推荐内容排序规则
	ArticleSortByNotLogin commonv1.CMN_FILE_SORT `json:"articleSortByNotLogin,omitempty"`
	// 展示热门帖子的数量
	ArticleHotShowSum int32 `json:"articleHotShowSum,omitempty"`
	// 展示近期多少天内创建的帖子
	ArticleHotShowWithLatest int32 `json:"articleHotShowWithLatest,omitempty"`
	// 展示近期的活动数量
	ActivityLatestShowSum int32 `json:"activityLatestShowSum,omitempty"`
	// 启用首页banner
	BannerImg string `json:"bannerImg,omitempty"`
	// banner文案
	BannerTitle string `json:"bannerTitle,omitempty"`
	// banner的跳转链接
	BannerLink string `json:"bannerLink,omitempty"`
	// 新用户注册默认访问页面
	DefaultPageByNewUser string `json:"defaultPageByNewUser,omitempty"`
	// 未登录用户默认访问页面
	DefaultPageByNotLogin string `json:"defaultPageByNotLogin,omitempty"`
	// 登录用户默认访问页面
	DefaultPageByLogin string `json:"defaultPageByLogin,omitempty"`
	// 页面选项
	PageOptions []PageOption `json:"pageOptions"`
}

type PageOption struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

// Get 获取首页配置信息
func Get(c *bffcore.Context) {
	//_, err := invoker.GrpcCommunity.Home(c.Ctx(), &communityv1.HomeReq{})
	//if err != nil {
	//	c.EgoJsonI18N(err)
	//	return
	//}
	homeOption, err := invoker.GrpcCommunity.GetHomeOption(c.Ctx(), &communityv1.GetHomeOptionReq{
		Uid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	// 获取开放的空间列表
	spaceResp, err := invoker.GrpcSpace.ListPublicSpace(c.Ctx(), &spacev1.ListPublicSpaceReq{
		Uid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	pageOptions := []PageOption{
		{
			Title: "默认首页",
			Value: "sys_home",
		},
	}

	// 获取开放的空间列表
	spaces := lo.Map(spaceResp.SimpleSpaceInfo, func(i *spacev1.SimpleSpaceInfo, v int) PageOption {
		return PageOption{
			Title: i.Name,
			Value: "space_" + i.SpaceGuid,
		}
	})

	pageOptions = append(pageOptions, spaces...)

	resp := GetRes{
		IsSetHome:                homeOption.IsSetHome,
		IsSetBanner:              homeOption.GetIsSetBanner(),
		ArticleSortByLogin:       homeOption.GetArticleSortByLogin(),
		ArticleSortByNotLogin:    homeOption.ArticleSortByNotLogin,
		ArticleHotShowSum:        homeOption.ArticleHotShowSum,
		ArticleHotShowWithLatest: homeOption.ArticleHotShowWithLatest,
		BannerImg:                homeOption.GetBannerImg(),
		BannerTitle:              homeOption.GetBannerTitle(),
		BannerLink:               homeOption.GetBannerLink(),
		DefaultPageByNewUser:     homeOption.DefaultPageByNewUser,
		DefaultPageByNotLogin:    homeOption.DefaultPageByNotLogin,
		DefaultPageByLogin:       homeOption.DefaultPageByLogin,
		PageOptions:              pageOptions,
	}
	c.JSONOK(resp)
}

type PutReq struct {
	// 是否启用首页
	IsSetHome *bool `json:"isSetHome,omitempty"`
	// 是否启用banner
	IsSetBanner *bool `json:"isSetBanner,omitempty"`
	// 登录用户推荐内容排序规则
	ArticleSortByLogin commonv1.CMN_FILE_SORT `json:"articleSortByLogin,omitempty"`
	// 未登录用户推荐内容排序规则
	ArticleSortByNotLogin commonv1.CMN_FILE_SORT `json:"articleSortByNotLogin,omitempty"`
	// 展示热门帖子的数量
	ArticleHotShowSum *int32 `json:"articleHotShowSum,omitempty"`
	// 展示近期多少天内创建的帖子
	ArticleHotShowWithLatest *int32 `json:"articleHotShowWithLatest,omitempty"`
	// 展示近期的活动数量
	ActivityLatestShowSum *int32 `json:"activityLatestShowSum,omitempty"`
	// 启用首页banner
	BannerImg *string `json:"bannerImg,omitempty"`
	// banner文案
	BannerTitle *string `json:"bannerTitle,omitempty"`
	// banner的跳转链接
	BannerLink *string `json:"bannerLink,omitempty"`
	// 新用户注册默认访问页面
	DefaultPageByNewUser *string `json:"defaultPageByNewUser,omitempty"`
	// 未登录用户默认访问页面
	DefaultPageByNotLogin *string `json:"defaultPageByNotLogin,omitempty"`
	// 登录用户默认访问页面
	DefaultPageByLogin *string `json:"defaultPageByLogin,omitempty"`
}

func Put(c *bffcore.Context) {
	req := PutReq{}
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "错误", err)
		return
	}

	_, err = invoker.GrpcCommunity.PutHomeOption(c.Ctx(), &communityv1.PutHomeOptionReq{
		Uid:                      c.Uid(),
		IsSetHome:                req.IsSetHome,
		IsSetBanner:              req.IsSetBanner,
		ArticleSortByLogin:       req.ArticleSortByLogin,
		ArticleSortByNotLogin:    req.ArticleSortByNotLogin,
		ArticleHotShowSum:        req.ArticleHotShowSum,
		ArticleHotShowWithLatest: req.ArticleHotShowWithLatest,
		BannerImg:                req.BannerImg,
		BannerTitle:              req.BannerTitle,
		BannerLink:               req.BannerLink,
		DefaultPageByNewUser:     req.DefaultPageByNewUser,
		DefaultPageByNotLogin:    req.DefaultPageByNotLogin,
		DefaultPageByLogin:       req.DefaultPageByLogin,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}
