package pms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	pmsv1 "ecodepost/pb/pms/v1"

	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestCreateRole(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := pmsv1.NewPmsClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.CreateRole(ctx, &pmsv1.CreateRoleReq{
		OperateUid: 1,
		Name:       "第一个自定义role",
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestGetRolePermission(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := pmsv1.NewPmsClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()

	res, err := cli.GetRolePermission(ctx, &pmsv1.GetRolePermissionReq{
		OperateUid: 1,
		RoleId:     1,
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestGetInitActionOptionPermission(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := pmsv1.NewPmsClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.GetInitActionOptionPermission(ctx, &pmsv1.GetInitActionOptionPermissionReq{
		OperateUid: 1,
		Guid:       "KO3Msd",
		Type:       commonv1.CMN_GUID_SPACE,
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestPutRolePermission(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := pmsv1.NewPmsClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()

	items := make([]*commonv1.PmsItem, 0)
	items = append(items, &commonv1.PmsItem{
		ActionName: "SPACE_CREATE",
		Flag:       1,
	})

	res, err := cli.PutRolePermission(ctx, &pmsv1.PutRolePermissionReq{
		OperateUid: 1,
		RoleId:     19,
		List:       nil,
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestPutRoleSpacePermission(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := pmsv1.NewPmsClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()

	spaceList := make([]*pmsv1.SpacePmsItem, 0)
	items := make([]*commonv1.PmsItem, 0)
	items = append(items, &commonv1.PmsItem{
		ActionName: "SPACE_CREATE_ARTICLE",
		Flag:       1,
	})
	items = append(items, &commonv1.PmsItem{
		ActionName: "SPACE_SET",
		Flag:       1,
	})
	spaceList = append(spaceList, &pmsv1.SpacePmsItem{
		Guid: "KO3Msd",
		List: items,
	})
	res, err := cli.PutRoleSpacePermission(ctx, &pmsv1.PutRoleSpacePermissionReq{
		OperateUid: 1,
		RoleId:     19,
		List:       nil,
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func prettyJsonPrint(protoMsg proto.Message) {
	jsonStr, _ := json.Marshal(protoMsg)
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, jsonStr, "", "\t")
	if error != nil {
		log.Println("prettyJsonPrint error: ", error)
		log.Println("origin str: " + string(jsonStr))
		return
	}

	fmt.Println(string(prettyJSON.Bytes()))
}
