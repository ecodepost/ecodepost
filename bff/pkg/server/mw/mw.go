package mw

import (
	bbfcore "ecodepost/bff/pkg/server/bffcore"

	"github.com/gin-gonic/gin"
)

const CookieI18n = "lan"

func I18nCookie() gin.HandlerFunc {
	return bbfcore.Handle(func(c *bbfcore.Context) {
		lan, err := c.Cookie(CookieI18n)
		// 没有这个cookie，那么设置为0
		if err != nil {
			c.Set(bbfcore.ContextLanguage, "cn")
			c.Next()
			return
		}
		if lan == "" {
			c.Set(bbfcore.ContextLanguage, "cn")
			c.Next()
			return
		}
		if lan == "cn" || lan == "en" {
			c.Set(bbfcore.ContextLanguage, lan)
		} else {
			c.Set(bbfcore.ContextLanguage, "cn")
		}
		c.Next()
	})
}
