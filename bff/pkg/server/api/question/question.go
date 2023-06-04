package question

import (
	"encoding/json"
	"strings"

	"ecodepost/bff/pkg/dto"
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/api/common"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/util/filecontent"
	notifyv1 "ecodepost/pb/notify/v1"
	userv1 "ecodepost/pb/user/v1"
	"github.com/spf13/cast"

	commonv1 "ecodepost/pb/common/v1"
	filev1 "ecodepost/pb/file/v1"
	questionv1 "ecodepost/pb/question/v1"
)

type CreateRequest struct {
	Name      string          `json:"name" binding:"required" label:"名称"`
	SpaceGuid string          `json:"spaceGuid" binding:"required" label:"空间ID"`
	Content   json.RawMessage `json:"content" binding:"required" label:"内容"`
}

func Create(c *bffcore.Context) {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	c.InjectSpc(req.SpaceGuid)

	strByte, err := req.Content.MarshalJSON()
	if err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}

	res, err := invoker.GrpcQuestion.Create(c.Ctx(), &questionv1.CreateReq{
		Uid:       c.Uid(),
		Name:      req.Name,
		SpaceGuid: req.SpaceGuid,
		Content:   string(strByte),
		Ip:        c.ClientIP(),
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}
	c.JSONOK(res.File)
}

func CreateAnswer(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}
	c.InjectGuid(guid).InjectSpc(req.SpaceGuid)

	// 查询问题
	q, err := invoker.GrpcFile.GetShowInfo(c.Ctx(), &filev1.GetShowInfoReq{
		Guid: guid,
		Uid:  c.Uid(),
	})
	if err != nil {
		c.JSONE(1, "query question fail", err)
		return
	}

	strByte, err := req.Content.MarshalJSON()
	if err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}

	// 创建回答
	res, err := invoker.GrpcQuestion.Create(c.Ctx(), &questionv1.CreateReq{
		Uid:        c.Uid(),
		Name:       req.Name,
		SpaceGuid:  req.SpaceGuid,
		ParentGuid: guid,
		Content:    string(strByte),
		Ip:         c.ClientIP(),
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}

	// 发送站内通知
	u := c.User()
	md, _ := json.Marshal(dto.MetaAnswer{
		QuestionGuid:  q.File.Guid,
		QuestionTitle: q.File.Name,
		AnswerAuthor:  userv1.UserInfo{Uid: u.Uid, Avatar: u.Avatar, Nickname: u.Nickname},
	})
	_, err = invoker.GrpcNotify.SendMsg(c.Ctx(), &notifyv1.SendMsgReq{
		TplId: 3,
		Msgs:  []*notifyv1.Msg{{Receiver: cast.ToString(q.File.Uid)}},
		VarLetter: &notifyv1.Letter{
			Type:     commonv1.NTF_TYPE_NEW_ANSWER,
			TargetId: guid,
			Meta:     md,
		},
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}
	c.JSONOK(res.File)
}

type UpdateArticleRequest struct {
	Name    string          `json:"name" binding:"required" label:"名称"`
	Content json.RawMessage `json:"content" label:"内容"`
}

func Update(c *bffcore.Context) {
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
	c.InjectGuid(guid)
	pContent, err := filecontent.GetPointerContent(req.Content)
	if err != nil {
		c.JSONE(1, "内容参数错误", err)
		return
	}

	_, err = invoker.GrpcQuestion.Update(c.Ctx(), &questionv1.UpdateReq{
		Uid:     c.Uid(),
		Guid:    guid,
		Name:    req.Name,
		Content: pContent,
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}
	c.JSONOK()
}

func UpdateAnswer(c *bffcore.Context) {
	guid := c.Param("answerGuid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	var req UpdateArticleRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	c.InjectGuid(guid)

	pContent, err := filecontent.GetPointerContent(req.Content)
	if err != nil {
		c.JSONE(1, "内容参数错误", err)
		return
	}

	_, err = invoker.GrpcQuestion.Update(c.Ctx(), &questionv1.UpdateReq{
		Uid:     c.Uid(),
		Guid:    guid,
		Content: pContent,
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}
	c.JSONOK()
}

func LikeQuestion(c *bffcore.Context) {
	guid := strings.TrimSpace(c.Param("guid"))
	if guid == "" {
		c.JSONE(1, "guid can't be empty", nil)
		return
	}
	common.UserActOnBiz(c, commonv1.CMN_BIZ_QUESTION, guid, commonv1.CNT_ACT_LIKE, commonv1.CNT_ACTI_ADD)
}

func UndoLikeQuestion(c *bffcore.Context) {
	guid := strings.TrimSpace(c.Param("guid"))
	if guid == "" {
		c.JSONE(1, "guid can't be empty", nil)
		return
	}
	common.UserActOnBiz(c, commonv1.CMN_BIZ_QUESTION, guid, commonv1.CNT_ACT_LIKE, commonv1.CNT_ACTI_SUB)
}

func LikeAnswer(c *bffcore.Context) {
	guid := strings.TrimSpace(c.Param("answerGuid"))
	if guid == "" {
		c.JSONE(1, "guid can't be empty", nil)
		return
	}
	common.UserActOnBiz(c, commonv1.CMN_BIZ_ANSWER, guid, commonv1.CNT_ACT_LIKE, commonv1.CNT_ACTI_ADD)
}

func UndoLikeAnswer(c *bffcore.Context) {
	guid := strings.TrimSpace(c.Param("answerGuid"))
	if guid == "" {
		c.JSONE(1, "guid can't be empty", nil)
		return
	}
	common.UserActOnBiz(c, commonv1.CMN_BIZ_ANSWER, guid, commonv1.CNT_ACT_LIKE, commonv1.CNT_ACTI_SUB)
}

func LikeComment(c *bffcore.Context) {
	guid := strings.TrimSpace(c.Param("commentGuid"))
	if guid == "" {
		c.JSONE(1, "guid can't be empty", nil)
		return
	}
	common.UserActOnBiz(c, commonv1.CMN_BIZ_COMMENT, guid, commonv1.CNT_ACT_LIKE, commonv1.CNT_ACTI_ADD)
}

func UndoLikeComment(c *bffcore.Context) {
	guid := strings.TrimSpace(c.Param("commentGuid"))
	if guid == "" {
		c.JSONE(1, "guid can't be empty", nil)
		return
	}
	common.UserActOnBiz(c, commonv1.CMN_BIZ_COMMENT, guid, commonv1.CNT_ACT_LIKE, commonv1.CNT_ACTI_SUB)
}

func Info(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)
	resp, err := invoker.GrpcQuestion.MyInfo(c.Ctx(), &questionv1.MyInfoReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}
	c.JSONOK(resp)
}

func Delete(c *bffcore.Context) {
	guid := c.Param("guid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcQuestion.Delete(c.Ctx(), &questionv1.DeleteReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}
	c.JSONOK()
}

func DeleteAnswer(c *bffcore.Context) {
	guid := c.Param("answerGuid")
	if guid == "" {
		c.JSONE(1, "文档ID不能为空", nil)
		return
	}
	c.InjectGuid(guid)

	_, err := invoker.GrpcQuestion.Delete(c.Ctx(), &questionv1.DeleteReq{
		Uid:  c.Uid(),
		Guid: guid,
	})
	if err != nil {
		c.JSONE(1, "系统错误", err)
		return
	}
	c.JSONOK()
}
