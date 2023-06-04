package column

import (
	"encoding/json"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/util/filecontent"

	columnv1 "ecodepost/pb/column/v1"
)

type ListFilesReq struct {
	SpaceGuid string `form:"spaceGuid" binding:"required" label:"空间ID"`
}

func ListFiles(c *bffcore.Context) {
	var req ListFilesReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectSpc(req.SpaceGuid)

	res, err := invoker.GrpcColumn.ListFile(c.Ctx(), &columnv1.ListFileReq{
		Uid:       c.Uid(),
		SpaceGuid: req.SpaceGuid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	// todo 因为grpc转protobuf omity原因，没有数据会是null，判断一下，返回一个数组
	if len(res.Files) == 0 {
		c.JSONOK([]struct{}{})
		return
	}
	c.JSONOK(res)
}

type SidebarChangeSortRequest struct {
	FileGuid       string  `json:"fileGuid" binding:"required" label:"文件guid"`
	TargetFileGuid *string `json:"targetFileGuid"`
	DropPosition   *string `json:"dropPosition"`
	ParentFileGuid *string `json:"parentFileGuid"`
}

// SidebarChangeSort 修改树型目录顺序
// @Tags Article
// @Success 200 {object} bffcore.Res{data=filev1.ChangeSortRes}
func SidebarChangeSort(c *bffcore.Context) {
	var req SidebarChangeSortRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectGuid(req.FileGuid)

	resp, err := invoker.GrpcColumn.ChangeSort(c.Ctx(), &columnv1.ChangeSortReq{
		Uid:            c.Uid(),
		FileGuid:       req.FileGuid,
		DropPosition:   req.DropPosition,
		TargetFileGuid: req.TargetFileGuid,
		ParentFileGuid: req.ParentFileGuid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type ListPermissionRequest struct {
	SpaceGuid string `form:"spaceGuid"`
}

// ListPermission 权限
// @Tags Article
// @Success 200 {object} bffcore.Res{data=filev1.GetDocumentTreeRes}
func ListPermission(c *bffcore.Context) {
	var req ListPermissionRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误:"+err.Error(), err)
		return
	}
	c.InjectSpc(req.SpaceGuid)

	resp, err := invoker.GrpcColumn.ListPermission(c.Ctx(), &columnv1.ListPermissionReq{
		Uid:       c.Uid(),
		SpaceGuid: req.SpaceGuid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(resp)
}

type UpdateArticleRequest struct {
	HeadImage *string         `json:"headImage"  label:"头图"`
	Name      string          `json:"name" binding:"required" label:"名称"`
	Content   json.RawMessage `json:"content" label:"内容"`
}

// UpdateArticle 更新或发布文档
// @Tags Article
// @Success 200 {object} bffcore.Res{}
func UpdateArticle(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	var req UpdateArticleRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	pContent, err := filecontent.GetPointerContent(req.Content)
	if err != nil {
		c.JSONE(1, "内容参数错误", err)
		return
	}
	c.InjectGuid(guid)

	_, err = invoker.GrpcColumn.Update(c.Ctx(), &columnv1.UpdateReq{
		Uid:       c.Uid(),
		FileGuid:  guid,
		Name:      req.Name,
		Content:   pContent,
		HeadImage: req.HeadImage,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

func DeleteArticle(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)
	_, err := invoker.GrpcColumn.Delete(c.Ctx(), &columnv1.DeleteReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}
