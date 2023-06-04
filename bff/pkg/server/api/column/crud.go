package column

import (
	"encoding/json"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	columnv1 "ecodepost/pb/column/v1"
)

type CreateRequest struct {
	Name       string          `json:"name" binding:"required" label:"名称"`
	SpaceGuid  string          `json:"spaceGuid" binding:"required" label:"空间ID"`
	ParentGuid string          `json:"parentGuid"`
	Content    json.RawMessage `json:"content" binding:"required" label:"内容"`
	HeadImage  string          `json:"headImage"  label:"头图"`
}

// Create 创建文档
// @Tags Article
// @Success 200 {object} bffcore.Res{data=commonv1.FileInfo}
func Create(c *bffcore.Context) {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	strByte, err := req.Content.MarshalJSON()
	if err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	c.InjectSpc(req.SpaceGuid)
	res, err := invoker.GrpcColumn.Create(c.Ctx(), &columnv1.CreateReq{
		Uid:        c.Uid(),
		Name:       req.Name,
		SpaceGuid:  req.SpaceGuid,
		ParentGuid: req.ParentGuid,
		Content:    string(strByte),
		HeadImage:  req.HeadImage,
		Ip:         c.ClientIP(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res.File)
}
