package ssoapi

import (
	"fmt"

	"ecodepost/bff/pkg/dto"
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/service/ssoservice"
	userv1 "ecodepost/pb/user/v1"
	"github.com/gotomicro/ego/core/elog"
)

func Register(c *bffcore.Context) {
	fmt.Println(333)

	var req dto.OauthDirectRequest
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "参数错误", err)
		return
	}

	showTips, err := ssoservice.User.RequestValid(c.Request.Context(), dto.OauthRequest{
		Account:  req.Account,
		Password: req.Password,
		Code:     req.Account,
	}, ssoservice.RegisterLoginType)
	if err != nil {
		c.JSONE(1, showTips, err)
		return
	}

	verifyShowTips, err := ssoservice.Code.VerifyCode(c.Request.Context(), req.Account, req.Code, ssoservice.RegisterCodeType)
	if err != nil {
		c.JSONE(1, verifyShowTips, err)
		return
	}

	newUser, err := invoker.GrpcUser.Create(c.Request.Context(), &userv1.CreateReq{
		Password:   &req.Password,
		Phone:      req.Account,
		RegisterIp: c.ClientIP(),
	})
	elog.Info("newUser info", elog.Any("newUser", newUser))
	if err != nil {
		c.JSONE(1, "注册失败", err)
		return
	}

	// 直接登录模式
	responseTypeLoginDirect(c, req, ssoservice.RegisterLoginType)
}
