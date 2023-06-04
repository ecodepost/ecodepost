package column

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	filev1 "ecodepost/pb/file/v1"

	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestCreateDocument(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.CreateDocument(ctx, &filev1.CreateDocumentReq{
	// 	Uid:        1,
	// 	CmtGuid:    "xxx",
	// 	Name:       "hello world2222",
	// 	SpaceGuid:  "zMolXTJol3",
	// 	ParentGuid: "NlLKz70DQk",
	// 	Format:     0,
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
}

func TestGetDocumentTree(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.GetDocumentTree(ctx, &filev1.GetDocumentTreeReq{
	// 	Uid:       1,
	// 	CmtGuid:   "xxx",
	// 	SpaceGuid: "zMolXTJol3",
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
}

func TestUpdateDocument(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.UpdateDocument(ctx, &filev1.UpdateDocumentReq{
	// 	Uid:     1,
	// 	CmtGuid: "xxx",
	// 	Guid:    "RZzoPnDG5e",
	//
	// 	Name:    "hello world",
	// 	Content: "我觉得还不错22222",
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
}

// func TestGetDocumentDraftContentByCreator(t *testing.T) {
//	eapp.SetEgoDebug("true")
//	cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
//	ctx := context.Background()
//	res, err := cli.GetDocumentDraftContentByCreator(ctx, &filev1.GetDocumentDraftContentByCreator{
//		Uid:     1,
//		CmtGuid: "xxx",
//		suid:    "RZzoPnDG5e",
//	})
//	assert.NoError(t, err)
//	prettyJsonPrint(res)
// }
//
// func TestDocumentPublish(t *testing.T) {
//	eapp.SetEgoDebug("true")
//	cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
//	ctx := context.Background()
//	res, err := cli.PublishDocument(ctx, &filev1.PublishDocumentReq{
//		Uid:     1,
//		CmtGuid: "xxx",
//		suid:    "RZzoPnDG5e",
//	})
//	assert.NoError(t, err)
//	prettyJsonPrint(res)
// }

func TestGetDocumentContentByCreator(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.GetDocumentContentByCreator(ctx, &filev1.GetDocumentContentByCreatorReq{
	// 	Uid:     1,
	// 	CmtGuid: "xxx",
	// 	Guid:    "RZzoPnDG5e",
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
}

func TestGetDocumentContent(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.GetDocumentContent(ctx, &filev1.GetDocumentContentReq{
	// 	Uid:     1,
	// 	CmtGuid: "xxx",
	// 	Guid:    "RZzoPnDG5e",
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
}

func TestDocumentInfo(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// res, err := cli.GetDocument(ctx, &filev1.GetDocumentReq{
	// 	Uid:     2,
	// 	CmtGuid: "2QA6BWAoJx",
	// 	Guid:    "XED7XRH86a",
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
}

func TestListDocument(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// res, err := cli.ListDocument(ctx, &filev1.ListDocumentReq{
	// 	Uid:        1,
	// 	CmtGuid:    "aAWKgMl6qn",
	// 	SpaceGuid:  "oJWXhG",
	// 	Pagination: &commonv1.Pagination{},
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
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

func TestGetDocumentNewTree(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.GetDocumentTree(ctx, &filev1.GetDocumentTreeReq{
	// 	Uid:       2,
	// 	CmtGuid:   "PdkD056Zw2",
	// 	SpaceGuid: "KO3VSJ",
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
}

func TestGetDocumentSort(t *testing.T) {
	eapp.SetEgoDebug("true")
	// cli := filev1.NewFileClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx := context.Background()
	// res, err := cli.ChangeSortByTargetGuid(ctx, &filev1.ChangeSortReq{
	// 	Uid:           2,
	// 	CmtGuid:       "PdkD056Zw2",
	// 	FileGuid:      "ZBN6q7kKjA",
	// 	AfterFileGuid: "GYlDyjW6XR",
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
	//
	// res2, err2 := cli.GetDocumentTree(ctx, &filev1.GetDocumentTreeReq{
	// 	Uid:       2,
	// 	CmtGuid:   "PdkD056Zw2",
	// 	SpaceGuid: "KO3VSJ",
	// })
	// assert.NoError(t, err2)
	// prettyJsonPrint(res2)
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
