package bffsso

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"ecodepost/bff/pkg/consts"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/util/utiltoken"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	ssov1 "ecodepost/pb/sso/v1"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/transport"
	"golang.org/x/sync/singleflight"
)

const (
	oauthTokenName       = "sid"
	oauthParentTokenName = "saas_ticket"
)

type Component struct {
	config       *Config
	logger       *elog.Component
	name         string
	grpcClient   ssov1.SsoClient
	SingleFlight singleflight.Group
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	grpcConn := egrpc.DefaultContainer().Build(
		egrpc.WithReadTimeout(3*time.Second),
		egrpc.WithAddr(config.SsoAddr),
		egrpc.WithOnFail(config.SsoOnFail),
		egrpc.WithEnableAccessInterceptor(config.SsoEnableAccessInterceptor),
		egrpc.WithEnableAccessInterceptorReq(config.SsoEnableAccessInterceptorReq),
		egrpc.WithEnableAccessInterceptorRes(config.SsoEnableAccessInterceptorRes),
	)

	return &Component{
		name:         name,
		logger:       logger,
		config:       config,
		SingleFlight: singleflight.Group{},
		grpcClient:   ssov1.NewSsoClient(grpcConn.ClientConn),
	}
}

// OauthState base 64 编码 referer和state信息
type OauthState struct {
	State   string `json:"state,omitempty"`   // State
	Referer string `json:"referer,omitempty"` // Referer
	Ref     string `json:"ref,omitempty"`     // 邀请码
}

func (c *Component) Config() *Config {
	return c.config
}

func (c *Component) MustCheckToken() gin.HandlerFunc {
	return bffcore.Handle(func(ctx *bffcore.Context) {
		token, err := c.getToken(ctx.Context)
		if err != nil {
			elog.Error("MustCheckToken error", elog.FieldErr(err))
			ctx.EgoJsonI18N(err)
			ctx.Abort()
			return
		}
		userByToken, err := c.grpcClient.GetUserByToken(ctx, &ssov1.GetUserByTokenReq{
			Token: token,
		})
		if err != nil {
			ctx.EgoJsonI18N(errcodev1.AuthErrGetUserInfoByTokenError())
			ctx.Abort()
			return
		}
		user := &ssov1.User{
			Uid:      userByToken.Uid,
			Nickname: userByToken.Nickname,
			Username: userByToken.Username,
			Avatar:   userByToken.Avatar,
			Email:    userByToken.Email,
			Name:     userByToken.Name,
		}
		parentContext := transport.WithValue(ctx.Request.Context(), consts.XLinkSpaceUid, user.Uid)
		ctx.Request = ctx.Request.WithContext(parentContext)
		ctx.Set(bffcore.ContextUserInfoKey, user)
		ctx.Next()
	})
}

func (c *Component) refreshTokenFn(token string, ctx *gin.Context) (result interface{}, err error) {
	basicAuthEncode := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.ClientID, c.config.ClientSecret)))
	refreshToken, err := c.grpcClient.RefreshToken(ctx.Request.Context(), &ssov1.RefreshTokenReq{
		Code:          token,
		Authorization: "Basic " + basicAuthEncode,
		ClientIP:      ctx.ClientIP(),
		ClientUA:      ctx.Request.Header.Get("User-Agent"),
	})
	if err != nil {
		return
	}
	return refreshToken, nil
}

func (c *Component) getToken(ctx *gin.Context) (token string, err error) {
	var pToken string
	cookie, err := ctx.Request.Cookie(oauthTokenName)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			err = errcodev1.AuthErrBrowserCookieSystemError()
			return
		}
		// 当没有短token短时候，直接发放
		pToken, err = ctx.Cookie(oauthParentTokenName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				err = errcodev1.AuthErrBrowserParentCookieEmpty()
				return
			}
			err = fmt.Errorf("ptoken originErr: %v err: %w", err, errcodev1.AuthErrBrowserParentCookieSystemError())
			return
		}
		var access *ssov1.GetAccessByPTokenRes
		access, err = c.grpcClient.GetAccessByPToken(ctx.Request.Context(), &ssov1.GetAccessByPTokenReq{
			ClientId:     c.config.ClientID,
			ClientSecret: c.config.ClientSecret,
			PToken:       pToken,
			Platform:     commonv1.SSO_PLATFORM_WEB,
			ClientIp:     ctx.ClientIP(),
			ClientUa:     ctx.Request.Header.Get("User-Agent"),
		})
		if err != nil {
			ctx.SetCookie(oauthParentTokenName, pToken, -1, "/", utiltoken.GetDomain(ctx, c.config.TokenDomain), c.config.TokenSecure, true)
			err = fmt.Errorf("ptoken originErr: %v err: %w", err, errcodev1.AuthErrGetAccessByParentCookieError())
			return
		}

		token = access.Token
		ctx.SetCookie(oauthTokenName, access.Token, int(access.ExpiresIn), "/", utiltoken.GetDomain(ctx, c.config.TokenDomain), c.config.TokenSecure, true)
	} else {
		// 有短token，作refresh操作
		token, _ = url.QueryUnescape(cookie.Value)
		// 如果小于7天，需要换token
		if cookie.Expires.Unix() < c.config.NeedRefreshTokenDuration {
			key := fmt.Sprintf("refreshToken-%s", token)
			result, e, _ := c.SingleFlight.Do(key, func() (interface{}, error) {
				return c.refreshTokenFn(token, ctx)
			})
			if e != nil {
				err = fmt.Errorf("ptoken originErr: %v err: %w", e, errcodev1.AuthErrRefreshTokenError())
				return
			}
			refreshTokenResult := result.(*ssov1.RefreshTokenRes)
			if refreshTokenResult.Token != token {
				token = refreshTokenResult.Token
				ctx.SetCookie(oauthTokenName, refreshTokenResult.Token, int(refreshTokenResult.ExpiresIn), "/", utiltoken.GetDomain(ctx, c.config.TokenDomain), c.config.TokenSecure, true)
			}
		}
	}
	return
}

// MayBeLogin 尝试看是否登录，如果登录就设置用户，如果没登录，那么就直接不设置变量
func (c *Component) MayBeLogin() gin.HandlerFunc {
	return bffcore.Handle(func(ctx *bffcore.Context) {
		token, err := c.getToken(ctx.Context)
		if err != nil {
			elog.Error("MayBeLogin error", elog.FieldErr(err))
			ctx.Next()
			return
		}

		userByToken, err := c.grpcClient.GetUserByToken(ctx.Request.Context(), &ssov1.GetUserByTokenReq{Token: token})
		if err != nil {
			ctx.Next()
			return
		}
		user := &ssov1.User{
			Uid:      userByToken.Uid,
			Nickname: userByToken.Nickname,
			Username: userByToken.Username,
			Avatar:   userByToken.Avatar,
			Email:    userByToken.Email,
			Name:     userByToken.Name,
		}

		parentContext := transport.WithValue(ctx.Request.Context(), consts.XLinkSpaceUid, user.Uid)
		ctx.Request = ctx.Request.WithContext(parentContext)
		ctx.Set(bffcore.ContextUserInfoKey, user)
		ctx.Next()
	})
}

// OauthLogout 退出登录Handler
func (c *Component) OauthLogout(ctx *bffcore.Context) {
	token, err := ctx.Cookie(oauthTokenName)
	if err != nil {
		ctx.JSONE(1, "获取登录态token失败: "+err.Error(), err)
		return
	}

	_, err = c.grpcClient.RemoveToken(ctx, &ssov1.RemoveTokenReq{
		Token: token,
	})
	if err != nil {
		ctx.JSONE(bffcore.CodeErr, "获取remove access token信息失败", err)
		return
	}
	ctx.JSONOK()
	return
}
