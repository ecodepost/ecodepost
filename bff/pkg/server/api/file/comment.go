package file

import (
	"errors"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	commentv1 "ecodepost/pb/comment/v1"
	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"

	"github.com/gotomicro/ego/core/eerrors"
)

type TopicCommentListRequest struct {
	CurrentPage int32 `form:"currentPage" label:"页码"`
}

func CommentList(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	var req TopicCommentListRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}

	resp, err := invoker.GrpcComment.List(c.Ctx(), &commentv1.ListReq{
		BizGuid: guid,
		Pagination: &commonv1.Pagination{
			CurrentPage: req.CurrentPage,
		},
	})
	if err != nil {
		egoErr := eerrors.FromError(err)
		if errors.Is(egoErr, errcodev1.ErrNotFound()) {
			c.JSONListPage([]struct{}{}, nil)
			return
		}
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(resp.List, resp.GetPagination())
}

type TopicChildCommentListRequest struct {
	CurrentPage int32 `form:"currentPage" label:"页码"`
}

func ChildCommentList(c *bffcore.Context) {
	commentGuid := c.Param("commentGuid")
	if commentGuid == "" {
		c.JSONE(1, "commentID不能为空", nil)
		return
	}
	c.InjectGuid(commentGuid)

	var req TopicChildCommentListRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	resp, err := invoker.GrpcComment.ChildList(c.Ctx(), &commentv1.ChildListReq{
		CommentGuid: commentGuid,
		Pagination: &commonv1.Pagination{
			CurrentPage: req.CurrentPage,
		},
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONListPage(resp.List, resp.GetPagination())
}

type CreateCommentReq struct {
	Guid        string `json:"guid" binding:"required"`
	CommentGuid string `json:"commentGuid"`
	Content     string `json:"content" binding:"required"`
}

// CreateComment 创建评论
func CreateComment(c *bffcore.Context) {
	var req CreateCommentReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err.Error())
		return
	}
	c.InjectGuid(req.Guid)

	_, err := invoker.GrpcComment.Create(c.Ctx(), &commentv1.CreateReq{
		BizGuid:     req.Guid,
		BizType:     commonv1.CMN_BIZ_ARTICLE,
		Uid:         c.Uid(),
		CommentGuid: req.CommentGuid,
		Content:     req.Content,
		ActionType:  commonv1.FILE_ACT_COMMENT,
		Ip:          c.ClientIP(),
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

// DeleteComment 删除评论
func DeleteComment(c *bffcore.Context) {
	commentGuid := c.Param("commentGuid")
	if commentGuid == "" {
		c.JSONE(1, "commentGuid不能为空", nil)
		return
	}
	c.InjectGuid(commentGuid)
	_, err := invoker.GrpcComment.Delete(c.Ctx(), &commentv1.DeleteReq{
		Uid:         c.Uid(),
		CommentGuid: commentGuid,
		DeleteType:  commonv1.FILE_CMET_DEL_USER,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}
