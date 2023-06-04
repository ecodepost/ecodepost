package api

import (
	"fmt"
	"testing"

	"ecodepost/bff/pkg/server/api/space"
	"ecodepost/bff/pkg/server/bffcore"

	commonv1 "ecodepost/pb/common/v1"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/unittest/gintest"
)

func TestCreateGroup(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(space.CreateGroup), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithJsonBody(space.CreateOrUpdateGroupRequest{
			Name:                 "xxxx",
			Icon:                 "",
			Visibility:           commonv1.CMN_VISBL_INTERNAL,
			IsOpenReadMemberList: true,
		}))
		fmt.Printf("info--------------->"+"%+v\n", string(info))
		return nil
	})
	tt.Run()
}

func TestChangeGroupSort(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(space.ChangeGroupSort), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithJsonBody(space.ChangeGroupSortRequest{
			AfterSpaceGroupGuid: "xx",
			SpaceGroupGuid:      "yy",
		}))
		fmt.Printf("info--------------->"+"%+v\n", string(info))
		return nil
	})
	tt.Run()
}

func TestCreateSpace(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(space.Create), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithJsonBody(space.CreateRequest{
			SpaceGroupGuid: "Kd3Ubo",
			Name:           "新的space",
			Icon:           "",
			SpaceType:      commonv1.CMN_APP_ARTICLE,
			SpaceLayout:    commonv1.SPC_LAYOUT_ARTICLE_FEED,
			Visibility:     commonv1.CMN_VISBL_INTERNAL,
		}))
		fmt.Printf("info--------------->"+"%+v\n", string(info))
		return nil
	})
	tt.Run()
}

func TestAddSpaceMember(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(space.CreateMember), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithUri(fmt.Sprintf("/api/spaces/%s/members", "KdpJiz")), gintest.WithJsonBody(space.CreateMemberRequest{
			AddUids: []int64{77},
		}))
		prettyJsonPrint(info)
		return nil
	}, gintest.WithRoutePath("/api/spaces/:guid/members"))
	tt.Run()
}

func TestUpdateSpacePermission(t *testing.T) {
	eapp.SetEgoDebug("true")
	// tt.POST(bffcore.Handle(space.UpdatePermission), func(m *gintest.Mock) error {
	// 	info := m.Exec(gintest.WithJsonBody(space.UpdatePermissionRequest{
	// 		IsOpenMemberPost:               false,
	// 		IsOpenMemberComment:            false,
	// 		IsOpenSpaceTopReadMore:         false,
	// 		IsOpenSpaceGroupReadMemberList: false,
	// 		IsOpenSpaceReadMemberList:      false,
	// 	}))
	// 	fmt.Printf("info--------------->"+"%+v\n", string(info))
	// 	return nil
	// })
	// tt.Run()
}

func TestSpaceTree(t *testing.T) {
	// eapp.SetEgoDebug("true")
	// tt.POST(bffcore.Handle(space.Tree), func(m *gintest.Mock) error {
	// 	info := m.Exec()
	// 	prettyJsonPrint(info)
	// 	return nil
	// })
	// tt.Run()
}
