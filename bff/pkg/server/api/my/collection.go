package my

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/api/common"
	"ecodepost/bff/pkg/server/bffcore"
	filev1 "ecodepost/pb/file/v1"

	commonv1 "ecodepost/pb/common/v1"
	statv1 "ecodepost/pb/stat/v1"

	"github.com/spf13/cast"
)

func CollectionGroupList(c *bffcore.Context) {
	list, err := invoker.GrpcStat.CollectionGroupList(c.Ctx(), &statv1.CollectionGroupListReq{
		Uid: c.Uid(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(list)
}

type CollectionGroupUpdateReq struct {
	Title *string `json:"title" binding:"required" label:"标题"`
	Desc  *string `json:"desc" label:"描述"`
}

func CollectionGroupUpdate(c *bffcore.Context) {
	cgid := cast.ToInt64(c.Param("cgid"))
	if cgid == 0 {
		c.JSONE(1, "CollectionGroupId can't be 0", nil)
		return
	}
	var req CollectionGroupUpdateReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	_, err := invoker.GrpcStat.CollectionGroupUpdate(c.Ctx(), &statv1.CollectionGroupUpdateReq{
		Id:    cgid,
		Uid:   c.Uid(),
		Title: req.Title,
		Desc:  req.Desc,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

// CollectionGroupDelete 会移除收藏分组及关联收藏目标记录，前端需给提示
func CollectionGroupDelete(c *bffcore.Context) {
	cgid := cast.ToInt64(c.Param("cgid"))
	if cgid == 0 {
		c.JSONE(1, "CollectionGroupId can't be 0", nil)
		return
	}
	_, err := invoker.GrpcStat.CollectionGroupDelete(c.Ctx(), &statv1.CollectionGroupDeleteReq{
		Uid: c.Uid(),
		Id:  cgid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type CollectionGroupCreateReq struct {
	Title string `json:"title" binding:"required" label:"标题"`
	Desc  string `json:"desc" label:"描述"`
}

func CollectionGroupCreate(c *bffcore.Context) {
	var req CollectionGroupCreateReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	_, err := invoker.GrpcStat.CollectionGroupCreate(c.Ctx(), &statv1.CollectionGroupCreateReq{
		Uid:   c.Uid(),
		Title: req.Title,
		Desc:  req.Desc,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type CollectionListReq struct {
	bffcore.Pagination
}

type CollectionInfo struct {
	// id
	Id int64 `json:"id"`
	// 收藏人Uid
	Uid int64 `json:"uid"`
	// 需要添加的收藏夹ID列表
	CollectionGroupIds []int64 `json:"collectionGroupIds"`
	// 业务Guid
	BizGuid string `json:"bizGuid"`
	// 业务类型
	BizType commonv1.CMN_BIZ `json:"bizType"`
	// 业务实例
	BizItem any `json:"bizItem"`
}

func CollectionList(c *bffcore.Context) {
	cgid := cast.ToInt64(c.Param("cgid"))
	if cgid == 0 {
		c.JSONE(1, "CollectionGroupId can't be 0", nil)
		return
	}
	var req CollectionListReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}

	res, err := invoker.GrpcStat.CollectionList(c.Ctx(), &statv1.CollectionListReq{
		Uid:               c.Uid(),
		CollectionGroupId: cgid,
		Pagination:        req.Pagination.ToPb(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	items := make([]common.BizItem, 0, len(res.List))
	for _, v := range res.List {
		items = append(items, common.BizItem{Type: v.BizType, Guid: v.BizGuid})
	}
	ret, e := common.ListBizItem(c.Ctx(), c.Uid(), items)
	if e != nil {
		c.EgoJsonI18N(err)
		return
	}

	list := make([]CollectionInfo, 0, len(res.List))
	for _, v := range res.List {
		// 如果存在BizType
		if bizItems, ok1 := ret[v.BizType]; ok1 {
			// 如果存在guid
			if v2, ok2 := bizItems[v.BizGuid]; ok2 {
				list = append(list, CollectionInfo{
					Id:                 v.Id,
					Uid:                v.Uid,
					CollectionGroupIds: v.CollectionGroupIds,
					BizGuid:            v.BizGuid,
					BizType:            v.BizType,
					BizItem:            v2,
				})
			}
		}
	}

	c.JSONListPage(list, res.GetPagination())
}

type CollectionCreateReq struct {
	CollectionGroupIds []int64          `json:"collectionGroupIds"  binding:"required" label:"分组Id"`
	Guid               string           `json:"guid"  binding:"required" label:"业务Guid"`
	Type               commonv1.CMN_BIZ `json:"type"  binding:"required" label:"业务类型"`
}

func CollectionCreate(c *bffcore.Context) {
	var req CollectionCreateReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	c.InjectGuid(req.Guid)
	_, err := invoker.GrpcFile.CollectionCreate(c.Ctx(), &filev1.CollectionCreateReq{
		Uid:                c.Uid(),
		CollectionGroupIds: req.CollectionGroupIds,
		BizGuid:            req.Guid,
		BizType:            req.Type,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type CollectionDeleteReq struct {
	CollectionGroupIds []int64          `json:"collectionGroupIds"  binding:"required" label:"分组Id"`
	Guid               string           `json:"guid"  binding:"required" label:"唯一ID"`
	Type               commonv1.CMN_BIZ `json:"type"  binding:"required" label:"类型"`
}

func CollectionDelete(c *bffcore.Context) {
	var req CollectionDeleteReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	c.InjectGuid(req.Guid)
	_, err := invoker.GrpcFile.CollectionDelete(c.Ctx(), &filev1.CollectionDeleteReq{
		Uid:                c.Uid(),
		CollectionGroupIds: req.CollectionGroupIds,
		BizGuid:            req.Guid,
		BizType:            req.Type,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}

	c.JSONOK()
}
