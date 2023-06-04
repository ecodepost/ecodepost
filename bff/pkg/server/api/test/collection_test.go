package api

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"ecodepost/bff/pkg/server/api/my"
	"ecodepost/bff/pkg/server/bffcore"

	commonv1 "ecodepost/pb/common/v1"

	"github.com/gotomicro/unittest/gintest"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestCRUDCollection(t *testing.T) {
	var cgid int64
	// 创建分组
	tt.POST(bffcore.Handle(my.CollectionGroupCreate), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithJsonBody(my.CollectionGroupCreateReq{
			Title: "title-1",
			Desc:  "desc-1",
		}))
		var res bffcore.Res
		err := json.Unmarshal(info, &res)
		assert.NoError(t, err)
		cgid = cast.ToInt64(res.Data)
		log.Printf("cgid--------------->"+"%+v\n", cgid)
		return nil
	})
	// 修改分组
	tt.PUT(bffcore.Handle(my.CollectionGroupUpdate), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri("/api/my/collection-groups/"+cast.ToString(cgid)), gintest.WithJsonBody(my.CollectionGroupUpdateReq{
			Title: newString("title-1-1"),
			Desc:  newString("desc-1-1"),
		}))
		fmt.Printf("Update info--------------->"+"%+v\n", string(info))
		return nil
	}, gintest.WithRoutePath("/api/my/collection-groups/:id"))
	// 查询分组列表
	tt.GET(bffcore.Handle(my.CollectionGroupList), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri("/api/my/collection-groups"))
		fmt.Printf("CollectionGroupList info--------------->"+"%+v\n", string(info))
		return nil
	}, gintest.WithRoutePath("/api/my/collection-groups"))

	// 收藏某个目标
	var cid int64
	tt.POST(bffcore.Handle(my.CollectionCreate), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithJsonBody(my.CollectionCreateReq{
			CollectionGroupIds: []int64{cgid},
			Guid:               "gv6EJrHyo4",
			Type:               commonv1.CMN_BIZ_ARTICLE,
		}))
		var res bffcore.Res
		err := json.Unmarshal(info, &res)
		assert.NoError(t, err)
		cid = int64(cast.ToIntSlice(res.Data)[0])
		log.Printf("cid--------------->"+"%+v\n", cid)
		return nil
	})
	// 取消收藏某个目标
	// tt.PUT(bffcore.Handle(my.CollectionDelete), func(m *gintest.Mock) error {
	// 	info := m.Exec(gintest.WithUri("/my/collection-groups/-/collections"), gintest.WithJsonBody(my.CollectionDeleteReq{
	// 		CollectionGroupIds: []int64{cast.ToInt64(cgid)},
	// 		Guid:               "article-xxx",
	// 		Type:               commonv1.CMN_BIZ_ARTICLE,
	// 	}))
	// 	fmt.Printf("CollectionDelete info--------------->"+"%+v\n", string(info))
	// 	return nil
	// }, gintest.WithRoutePath("/my/collection-groups/-/collections"))
	// 查看收藏列表
	tt.PUT(bffcore.Handle(my.CollectionList), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri("/my/collection-groups/"+cast.ToString(cgid)+"/collections"), gintest.WithJsonBody(my.CollectionListReq{
			Pagination: bffcore.Pagination{CurrentPage: 1, PageSize: 10},
		}))
		fmt.Printf("CollectionList info--------------->"+"%+v\n", string(info))
		return nil
	}, gintest.WithRoutePath("/my/collection-groups/:cgid/collections"))

	tt.Run()
}
