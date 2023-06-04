package job

import (
	"ecodepost/bff/pkg/server/bffcore"

	ssov1 "ecodepost/pb/sso/v1"

	"github.com/gin-gonic/gin"
)

func DebugToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u := &ssov1.User{
			Uid:      1,
			Nickname: "saas",
		}
		ctx.Set(bffcore.ContextUserInfoKey, u)
		ctx.Next()
	}
}
