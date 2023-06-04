package bffcore

import (
	ssov1 "ecodepost/pb/sso/v1"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserInfoKey   = "ctx-user-info"
	ContextLanguage      = "ctx-language"
	oauthParentTokenName = "saas_ticket"
)

type XEcodePostUserInfo struct{}

// IsAuthenticated 判断用户是否登录
func (c *Context) IsAuthenticated() bool {
	if user := contextUser(c.Context); user != nil && user.Uid > 0 {
		return true
	}
	return false
}

// User 返回当前用户全部数据
func (c *Context) User() *ssov1.User {
	return contextUser(c.Context)
}

// Uid 返回当前用户uid
func (c *Context) Uid() int64 {
	return userUid(c.Context)
}

// userUid 返回当前用户uid，入参使用gin.Context
func userUid(c *gin.Context) int64 {
	return contextUser(c).Uid
}

// contextUser 从context取用户，入参使用gin.Context
func contextUser(c *gin.Context) *ssov1.User {
	resp := &ssov1.User{}
	respI, flag := c.Get(ContextUserInfoKey)
	if flag {
		resp = respI.(*ssov1.User)
	}
	return resp
}

func (c *Context) GetParentToken() (string, error) {
	return c.Cookie(oauthParentTokenName)
}

func (c *Context) Login(token *ssov1.Token) {
	c.SetCookie(oauthParentTokenName, token.Token, int(token.ExpiresIn), "/", token.Domain, false, true)
}
