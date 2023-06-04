package api

import (
	"fmt"
	"testing"

	"ecodepost/bff/pkg/server/api/article"
	"ecodepost/bff/pkg/server/bffcore"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/unittest/gintest"
)

func TestArticleCreate(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(article.CreateArticle), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithJsonBody(article.CreateArticleRequest{
			Name:      "test",
			SpaceGuid: "gYGzhg",
			Content:   []byte("test"),
		}))
		prettyJsonPrint(info)
		return nil
	})
	tt.Run()
}

// GeZo9weD5w
func TestArticleRecommend(t *testing.T) {
	tt.POST(bffcore.Handle(article.Recommend), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri(fmt.Sprintf("/api/articles/%s/recommend", "GeZo9weD5w")))
		prettyJsonPrint(info)
		return nil
	}, gintest.WithRoutePath("/api/articles/:guid/recommend"))
	tt.Run()
}
