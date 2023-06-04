package user

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	userv1 "ecodepost/pb/user/v1"
)

func My(c *bffcore.Context) {
	res, err := invoker.GrpcUser.ProfileInfo(c.Ctx(), &userv1.ProfileInfoReq{Uid: c.Uid()})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res)
	return
}

// Info 查询单个用户
func Info(c *bffcore.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSONE(1, "user name can't be empty", nil)
		return
	}
	res, err := invoker.GrpcUser.Info(c.Ctx(), &userv1.InfoReq{Name: name})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK(res)
}

type UpdateAvatarRequest struct {
	Url string `json:"url"`
}

func UpdateAvatar(c *bffcore.Context) {
	req := UpdateAvatarRequest{}
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}

	_, err = invoker.GrpcUser.Update(c.Ctx(), &userv1.UpdateReq{
		Uid:    c.Uid(),
		Avatar: &req.Url,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type UpdateNicknameRequest struct {
	Nickname string `json:"nickname"`
}

func UpdateNickname(c *bffcore.Context) {
	var req UpdateNicknameRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	_, err := invoker.GrpcUser.Update(c.Ctx(), &userv1.UpdateReq{
		Uid:      c.Uid(),
		Nickname: &req.Nickname,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type UpdateAttrReq struct {
	Avatar   *string `json:"avatar"`
	Sex      *int32  `json:"sex"`
	Birthday *int64  `json:"birthday"`
	Intro    *string `json:"intro"`
}

func UpdateAttr(c *bffcore.Context) {
	var req UpdateAttrReq
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}
	_, err := invoker.GrpcUser.Update(c.Ctx(), &userv1.UpdateReq{
		Uid:      c.Uid(),
		Sex:      req.Sex,
		Birthday: req.Birthday,
		Intro:    req.Intro,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type UpdatePhoneRequest struct {
	Phone string `json:"phone"`
}

func UpdatePhone(c *bffcore.Context) {
	var req UpdatePhoneRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}

	_, err := invoker.GrpcUser.UpdatePhone(c.Ctx(), &userv1.UpdatePhoneReq{
		Uid:   c.Uid(),
		Phone: req.Phone,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

func UpdateEmail(c *bffcore.Context) {
	var req UpdateEmailRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}

	_, err := invoker.GrpcUser.UpdateEmail(c.Ctx(), &userv1.UpdateEmailReq{
		Uid:   c.Uid(),
		Email: req.Email,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}
