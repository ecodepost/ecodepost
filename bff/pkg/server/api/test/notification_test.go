package api

import (
	"fmt"
	"testing"
	"time"

	"ecodepost/bff/pkg/server/api/my"
	"ecodepost/bff/pkg/server/bffcore"

	commonv1 "ecodepost/pb/common/v1"

	"github.com/gotomicro/ego/core/eapp"

	"github.com/gotomicro/unittest/gintest"
)

func TestCreateNotification(t *testing.T) {
	eapp.SetEgoDebug("true")
	tt.POST(bffcore.Handle(my.NotificationCreate), func(m *gintest.Mock) error {
		info := m.Exec(gintest.WithJsonBody(my.NotificationCreateReq{
			Type:     commonv1.NTF_TYPE_SYSTEM,
			TargetId: "xxx",
			Uids:     []int64{147},
			Link:     "xxx",
			Ctime:    time.Now().Unix(),
		}))
		fmt.Printf("info--------------->"+"%+v\n", string(info))
		return nil
	})
	tt.Run()
}
