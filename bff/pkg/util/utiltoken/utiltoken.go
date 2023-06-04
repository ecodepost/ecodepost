package utiltoken

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetDomain 如果不是自定义域名，那么就使用c.Request.Host
// 如果是子域名，就使用token domain
func GetDomain(c *gin.Context, tokenDomain string) string {
	if strings.Contains(c.Request.Host, tokenDomain) {
		return tokenDomain
	}
	return c.Request.Host
}
