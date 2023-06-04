package file

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	commonv1 "ecodepost/pb/common/v1"
	filev1 "ecodepost/pb/file/v1"
)

type ListReq struct {
	SpaceGuid   string                 `form:"spaceGuid" binding:"required" label:"空间ID"`
	CurrentPage int32                  `form:"currentPage"` // 当前页数
	Sort        commonv1.CMN_FILE_SORT `form:"sort"`        // 排序值
}

// ListPage 文件分页列表
func ListPage(c *bffcore.Context) {
	var req ListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	resp, err := invoker.GrpcFile.ListPage(c.Ctx(), &filev1.ListPageReq{
		Uid:       c.Uid(),
		SpaceGuid: req.SpaceGuid,
		Pagination: &commonv1.Pagination{
			CurrentPage: req.CurrentPage,
		},
		Sort: req.Sort,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(resp.List, resp.GetPagination())
}

type SubListPageReq struct {
	CurrentPage int32                  `form:"currentPage"` // 当前页数
	Sort        commonv1.CMN_FILE_SORT `form:"sort"`        // 排序值
}

func SubListPage(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)
	var req SubListPageReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	resp, err := invoker.GrpcFile.ListPageByParent(c.Ctx(), &filev1.ListPageByParentReq{
		Uid:        c.Uid(),
		ParentGuid: guid,
		Pagination: &commonv1.Pagination{
			CurrentPage: req.CurrentPage,
		},
		Sort: req.Sort,
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}
	c.JSONListPage(resp.List, resp.GetPagination())
}

type RecommendsReq struct {
	SpaceGuid string `form:"spaceGuid" binding:"required" label:"空间ID"`
}

// Recommends 推荐列表 (0522 新增用户昵称，头像，摘要)
// @Tags Article
// @Success 200 {object} bffcore.Res{data=bffcore.ResPageData{list=[]filev1.ArticleShow}}
func Recommends(c *bffcore.Context) {
	var req RecommendsReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectSpc(req.SpaceGuid)
	resp, err := invoker.GrpcFile.RecommendList(c.Ctx(), &filev1.RecommendListReq{
		Uid:       c.Uid(),
		SpaceGuid: req.SpaceGuid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(resp.List, resp.GetPagination())
}

type SpaceTopsReq struct {
	SpaceGuid string `form:"spaceGuid" binding:"required" label:"空间ID"`
}

// SpaceTops 置顶列表 (0522 新增用户昵称，头像，摘要)
func SpaceTops(c *bffcore.Context) {
	var req SpaceTopsReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectSpc(req.SpaceGuid)
	resp, err := invoker.GrpcFile.SpaceTopList(c.Ctx(), &filev1.SpaceTopListReq{
		Uid:       c.Uid(),
		SpaceGuid: req.SpaceGuid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(resp.List, resp.GetPagination())
}
