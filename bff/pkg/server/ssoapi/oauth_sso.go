package ssoapi

import (
	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	ssov1 "ecodepost/pb/sso/v1"
	"github.com/gin-gonic/gin"
)

const oauthTokenName = "sid"

// AuthExist 如果auth，已经存在，那么就不需要执行下面方法
func AuthExist() gin.HandlerFunc {
	return bffcore.Handle(func(ctx *bffcore.Context) {
		token, err := ctx.GetParentToken()
		if err != nil {
			ctx.Next()
			return
		}

		_, err = invoker.GrpcSso.GetUserByParentToken(ctx.Request.Context(), &ssov1.GetUserByParentTokenReq{
			ParentToken: token,
		})
		if err != nil {
			ctx.Next()
			return
		}
		ctx.JSONE(200, "已经登录", nil)
		ctx.Abort()
	})
}
