package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"ecodepost/bff/pkg/server/api/article"
	"ecodepost/bff/pkg/server/api/home"
	"ecodepost/bff/pkg/server/bffcore"

	"github.com/gotomicro/ego/core/eapp"

	"github.com/gotomicro/unittest/gintest"
)

func TestCommunityHome(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(home.Page), func(m *gintest.Mock) error {
		info := m.Exec()
		fmt.Printf("info--------------->"+"%+v\n", string(info))
		return nil
	})
	tt.Run()
}

// content: "<p>a hahasdfasdf</p>"
// guid: "GXRDbqJKvQ"
// headImage: "https://cdn.gocn.vip/ofimage/art_head_img/G7joVBokn3/20220520/4272e36f6d7841dda9d369f40fb1c782.jpeg"
// name: "我jude"
func TestFilePublish(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(article.UpdateArticle), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri("/api/articles/"+"GXRDbqJKvQ"), gintest.WithJsonBody(article.UpdateArticleRequest{
			HeadImage: "https://cdn.gocn.vip/ofimage/art_head_img/G7joVBokn3/20220520/4272e36f6d7841dda9d369f40fb1c782.jpeg",
			Name:      "test",
			Content:   json.RawMessage("asdfasdfasdf"),
		}))
		fmt.Printf("info--------------->"+"%+v\n", string(info))
		return nil
	}, gintest.WithRoutePath("/api/articles/:guid"))
	tt.Run()
}
