package space

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	commonv1 "ecodepost/pb/common/v1"
	spacev1 "ecodepost/pb/space/v1"

	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestSpacePermission(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.GetSpacePermissionByUid(ctx, &spacev1.GetSpacePermissionByUidReq{
		OperateUid: 125,
		TargetGuid: "6BaQHP",
		GuidType:   commonv1.CMN_GUID_SPACE,
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestSpaceInfo(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.SpaceInfo(ctx, &commonv1.SpaceInfoReq{
	// 	CmtGuid:    "8YNKvqxoa5",
	// 	OperateUid: 125,
	// 	SpaceGuid:  "6BaQHP",
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
}

func TestCreateSpaceGroup(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.CreateSpaceGroup(ctx, &spacev1.CreateSpaceGroupReq{
		OperateUid: 1,
		Name:       "hello",
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestUpdateSpaceGroup(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.UpdateSpaceGroup(ctx, &spacev1.UpdateSpaceGroupReq{
		OperateUid:     1,
		Name:           "hello2",
		SpaceGroupGuid: "QRKxlTEKPZ",
	})
	assert.NoError(t, err)
	t.Logf("res: %+v", res)
}

func TestDeleteSpaceGroup(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.DeleteSpaceGroup(ctx, &spacev1.DeleteSpaceGroupReq{
		OperateUid:     1,
		SpaceGroupGuid: "QRKxlTEKPZ",
	})
	assert.NoError(t, err)
	t.Logf("res: %+v", res)
}

func TestCreateSpace(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.CreateSpace(ctx, &spacev1.CreateSpaceReq{
		OperateUid:     1,
		SpaceGroupGuid: "QRKxlTEKPZ",
		Name:           "我的空间2",
	})
	assert.NoError(t, err)
	t.Logf("res: %+v", res)
}

func TestUpdateSpace(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.UpdateSpace(ctx, &spacev1.UpdateSpaceReq{
	// 	OperateUid:     1,
	// 	CmtGuid:        "xxx",
	// 	SpaceGuid:      "zMolXTJol3",
	// 	Name:           "我觉得空间很不错",
	// 	IconType:       commonv1.FILE_IT_EMOJI,
	// 	SpaceGroupGuid: "QRKxlTEKPZ",
	// })
	// assert.NoError(t, err)
	// t.Logf("res: %+v", res)
}

func TestDeleteSpace(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.DeleteSpace(ctx, &spacev1.DeleteSpaceReq{
		OperateUid: 1,
		SpaceGuid:  "QRKxlTEKPZ",
	})
	assert.NoError(t, err)
	t.Logf("res: %+v", res)
}

func TestChangeSpaceSort(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := spacev1.NewSpaceClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.ChangeSpaceSort(ctx, &spacev1.ChangeSpaceSortReq{
		OperateUid: 1,
		SpaceGuid:  "zMolXTJol3",
	})
	assert.NoError(t, err)
	t.Logf("res: %+v", res)
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
