package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"ecodepost/resource-svc/pkg/service"

	commonv1 "ecodepost/pb/common/v1"
	filev1 "ecodepost/pb/file/v1"

	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestUpdateFileSize(t *testing.T) {
	eapp.SetEgoDebug("true")
	_, err := service.File.CreateFile(context.Background(), service.CreateOrCopyFileReq{
		Name:       "ttttt",
		Uid:        151,
		Content:    "123123",
		SpaceGuid:  "RLXJcB",
		CreateTime: time.Now().Unix(),
	}, commonv1.FILE_TYPE_DOCUMENT)
	if err != nil {
		fmt.Printf("err--------------->"+"%+v\n", err)
		return
	}
	cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.UpdateFileSize(ctx, &filev1.UpdateFileSizeReq{
		SpaceGuid: "RLXJcB",
		Guid:      "MKQJNyFgDA",
		Size:      9,
		Uid:       151,
	})

	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestEmojiCreate(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.CreateEmoji(ctx, &filev1.CreateEmojiReq{
		Uid:  1,
		Guid: "GeZo9weD5w",
		V:    1,
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestEmojiDecrease(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.DeleteEmoji(ctx, &filev1.DeleteEmojiReq{
		Uid:  1,
		Guid: "GeZo9weD5w",
		V:    1,
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestMyEmojiList(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.MyEmojiList(ctx, &filev1.MyEmojiListReq{
		Uid:   2,
		Guids: []string{"LbyKxlXKxX"},
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestFileListPage(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.ListPage(ctx, &filev1.ListPageReq{
		Uid:        2,
		SpaceGuid:  "11J9U2",
		Pagination: &commonv1.Pagination{},
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
