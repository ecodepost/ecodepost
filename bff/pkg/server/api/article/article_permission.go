package article

import (
	"ecodepost/bff/pkg/server/bffcore"

	"github.com/gotomicro/ego/core/econf"
)

type ListCoversRes struct {
	List []string `json:"list"` // 封面数组
}

func ListCovers(c *bffcore.Context) {
	c.JSONOK(ListCoversRes{
		List: econf.GetStringSlice("recommend.articleCovers"),
	})
}
