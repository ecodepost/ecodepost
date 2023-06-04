package bffsso

import (
	ssov1 "ecodepost/pb/sso/v1"

	"github.com/gin-gonic/gin"
)

// Config oauth2配置
type Config struct {
	SsoAddr                       string // 单点登录地址
	SsoOnFail                     string // 失败后处理模式，失败后的处理方式，panic | error，默认panic
	SsoEnableAccessInterceptor    bool
	SsoEnableAccessInterceptorReq bool
	SsoEnableAccessInterceptorRes bool
	TokenDomain                   string // 域名，目前sub token，parent token，同域名
	TokenSecure                   bool   // 写入token，开启HTTPS，默认开启
	AuthURL                       string
	RegisterURL                   string
	ClientID                      string
	ClientSecret                  string
	Scopes                        []string // Scope specifies optional requested permissions.
	RedirectURLFormat             string
	NeedRefreshTokenDuration      int64

	grpcClient ssov1.SsoClient
	router     *gin.RouterGroup
}

// DefaultConfig 定义了moauth默认配置
func DefaultConfig() *Config {
	return &Config{
		SsoOnFail:   "panic",
		TokenSecure: true, // 默认开启HTTPS
	}
}
