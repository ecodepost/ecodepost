package api

import (
	"testing"

	"ecodepost/bff/pkg/server/api/pms"
	"ecodepost/bff/pkg/server/bffcore"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/unittest/gintest"
)

func TestPmsCreateRole(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(pms.CreateRole), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithJsonBody(pms.CreateRoleRequest{
			Name: "第一个角色",
		}))
		prettyJsonPrint(info)
		return nil
	})
	tt.Run()
}

func TestPmsRoleList(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.GET(bffcore.Handle(pms.RoleList), func(m *gintest.Mock) error {
		info := m.Exec()
		prettyJsonPrint(info)
		return nil
	})
	tt.Run()
}
