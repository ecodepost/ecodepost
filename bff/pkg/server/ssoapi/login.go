package ssoapi

import (
	"context"

	"ecodepost/bff/pkg/dto"
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/service/ssoservice"
	ssov1 "ecodepost/pb/sso/v1"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"

	"go.uber.org/zap"
)

func LoginDirect(c *bffcore.Context) {
	// 如果已经登录
	reqView := dto.OauthDirectRequest{}
	err := c.Bind(&reqView)
	if err != nil {
		elog.Error("json marshal request body error", zap.Error(err))
		c.JSONE(1, "参数错误", err)
		return
	}

	showTips, err := ssoservice.User.RequestValid(c.Request.Context(), dto.OauthRequest{
		Account:  reqView.Account,
		Password: reqView.Password,
	}, ssoservice.PasswordLoginType)
	if err != nil {
		c.JSONE(1, showTips, err)
		return
	}
	// 直接登录模式
	responseTypeLogin(c, reqView, ssoservice.PasswordLoginType, false)
}

func responseTypeLogin(c *bffcore.Context, reqView dto.OauthDirectRequest, loginType ssoservice.LoginType, isDirect bool) {
	userInfo, showTips, err := ssoservice.User.Verify(c.Request.Context(), dto.OauthRequest{
		Account:  reqView.Account,
		Password: reqView.Password,
		Code:     reqView.Code,
	}, loginType)
	if err != nil {
		c.JSONE(1, showTips, err)
		return
	}
	responseSsoLogin(c, userInfo.User.Uid, isDirect)
}

func responseSsoLogin(c *bffcore.Context, uid int64, isDirect bool) {
	parentToken, _ := c.GetParentToken()
	loginRes, err := invoker.GrpcSso.Login(c.Request.Context(), &ssov1.LoginReq{
		ClientId:     econf.GetString("bff.sso.clientId"),
		ClientSecret: econf.GetString("bff.sso.clientSecret"),
		RedirectUri:  econf.GetString("bff.sso.redirectUri"),
		ParentToken:  parentToken,
		Uid:          uid,
		ClientIp:     c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
	})
	if err != nil {
		c.JSONE(1, "登录失败", err)
		return
	}

	c.Login(loginRes.Parent)
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "x-saas-uid", uid))
	c.SetCookie(oauthTokenName, loginRes.Sub.Token, int(loginRes.Sub.ExpiresIn), "/", loginRes.Sub.Domain, false, true)

	if isDirect {
		c.Redirect(302, loginRes.Sub.RedirectUri)
		return
	}
	c.JSONOK()
}
