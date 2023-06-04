package article

import (
	"encoding/json"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/util/filecontent"

	articlev1 "ecodepost/pb/article/v1"

	commonv1 "ecodepost/pb/common/v1"
)

type CreateArticleRequest struct {
	Name      string          `json:"name" binding:"required" label:"名称"`        // 名称
	SpaceGuid string          `json:"spaceGuid" binding:"required" label:"空间ID"` // 空间GUID
	Content   json.RawMessage `json:"content" binding:"required" label:"内容"`     // 内容，json数据结构
	HeadImage string          `json:"headImage"  label:"头图"`                     // 头图
}

// CreateArticle 创建文档
// @Tags Article
func CreateArticle(c *bffcore.Context) {
	var req CreateArticleRequest
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
	res, err := invoker.GrpcArticle.CreateDocument(c.Ctx(), &articlev1.CreateDocumentReq{
		Uid:       c.Uid(),
		Name:      req.Name,
		SpaceGuid: req.SpaceGuid,
		Content:   string(strByte),
		HeadImage: req.HeadImage,
		Format:    commonv1.FILE_FORMAT_DOCUMENT_SLATE,
		Ip:        c.ClientIP(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res.File)
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
		c.JSONE(1, err.Error(), err)
		return
	}
	c.InjectGuid(guid)

	pContent, err := filecontent.GetPointerContent(req.Content)
	if err != nil {
		c.JSONE(1, "内容参数错误", err)
		return
	}
	_, err = invoker.GrpcArticle.UpdateDocument(c.Ctx(), &articlev1.UpdateDocumentReq{
		Uid:        c.Uid(),
		Guid:       guid,
		Name:       req.Name,
		Content:    pContent,
		HeadImage:  req.HeadImage,
		FileFormat: commonv1.FILE_FORMAT_DOCUMENT_SLATE,
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

	_, err := invoker.GrpcArticle.DeleteDocument(c.Ctx(), &articlev1.DeleteDocumentReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

func SpaceTop(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcArticle.SetDocumentSpaceTop(c.Ctx(), &articlev1.SetDocumentSpaceTopReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

func CancelSpaceTop(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcArticle.CancelDocumentSpaceTop(c.Ctx(), &articlev1.CancelDocumentSpaceTopReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

func Recommend(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcArticle.SetDocumentRecommend(c.Ctx(), &articlev1.SetDocumentRecommendReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

func CancelRecommend(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcArticle.CancelDocumentRecommend(c.Ctx(), &articlev1.CancelDocumentRecommendReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

func OpenComment(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcArticle.OpenDocumentComment(c.Ctx(), &articlev1.OpenDocumentCommentReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

func CloseComment(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)
	_, err := invoker.GrpcArticle.CloseDocumentComment(c.Ctx(), &articlev1.CloseDocumentCommentReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}
