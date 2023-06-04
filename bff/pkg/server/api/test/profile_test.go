package api

import (
	"testing"

	"ecodepost/bff/pkg/server/api/my"
	profile "ecodepost/bff/pkg/server/api/public-profile"
	"ecodepost/bff/pkg/server/bffcore"

	ssov1 "ecodepost/pb/sso/v1"

	"github.com/gotomicro/unittest/gintest"
)

func TestFollowing(t *testing.T) {
	page := bffcore.Pagination{CurrentPage: 1, PageSize: 10}
	// -> 表示关注状态， 测试场景: 用户 A -> B
	uidB1 := "3"
	uidB2 := "4"
	uidB3 := "5"

	// 测试follow用户B1\B2\B3
	tt.GET(bffcore.Handle(my.FollowingCreate), func(m *gintest.Mock) error {
		info1 := m.Exec(gintest.WithUri("/api/my/following/" + uidB1))
		info2 := m.Exec(gintest.WithUri("/api/my/following/" + uidB2))
		info3 := m.Exec(gintest.WithUri("/api/my/following/" + uidB3))
		prettyJsonPrint(info1)
		prettyJsonPrint(info2)
		prettyJsonPrint(info3)
		return nil
	}, gintest.WithRoutePath("/api/my/following/:uid"))
	tt.GET(bffcore.Handle(profile.FollowingList), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri("/api/users/name_TomSawyer2/following"), gintest.WithQuery(profile.FollowersListReq{Pagination: page}))
		prettyJsonPrint(info)
		return nil
	}, gintest.WithRoutePath("/api/users/:name/following"))

	// 测试对用户B1移除follow
	tt.DELETE(bffcore.Handle(my.FollowingDelete), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri("/api/my/following/" + uidB1))
		prettyJsonPrint(info)
		return nil
	}, gintest.WithRoutePath("/api/my/following/:uid"))
	tt.GET(bffcore.Handle(profile.FollowingList), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri("/api/users/name_TomSawyer2/following"), gintest.WithQuery(profile.FollowersListReq{Pagination: page}))
		prettyJsonPrint(info)
		return nil
	}, gintest.WithRoutePath("/api/users/:name/following"))

	// 测试用户B1查看自己的关注者列表
	fn := func(c *bffcore.Context) {
		u := &ssov1.User{Uid: 4, Name: "name_1361****255"}
		c.Set(bffcore.ContextUserInfoKey, u)
		profile.FollowersList(c)
	}
	tt.GET(bffcore.Handle(fn), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri("/api/users/name_1361****255/followers"), gintest.WithQuery(profile.FollowersListReq{Pagination: page}))
		prettyJsonPrint(info)
		return nil
	}, gintest.WithRoutePath("/api/users/:name/followers"))

	tt.Run()
	return
}
