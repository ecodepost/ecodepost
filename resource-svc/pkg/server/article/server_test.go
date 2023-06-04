package article

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	// "git.yitum.com/gopkg/quill"
	articlev1 "ecodepost/pb/article/v1"

	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func Test(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := articlev1.NewArticleClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uid int64 = 1
	res, err := cli.HomeArticleHotList(ctx, &articlev1.HomeArticleHotListReq{
		Uid:        &uid,
		Limit:      5,
		LatestTime: 86400 * 30,
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

//
// func TestCreateDocument(t *testing.T) {
//	eapp.SetEgoDebug("true")
//	cli := articlev1.NewArticleClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	content := `{"ops":[{"insert":"hello world"}]}`
//	res, err := cli.CreateDocument(ctx, &articlev1.CreateDocumentReq{
//		Uid:       1,
//		CmtGuid:   "bPMDkrxK5x",
//		Name:      "hello world",
//		SpaceGuid: "RLXJcB",
//		Format:    commonv1.FILE_FORMAT_DOCUMENT_SLATE,
//		Content:   content,
//	})
//	str, _ := quill.Render([]byte(content))
//	fmt.Printf("str--------------->"+"%+v\n", string(str))
//	assert.NoError(t, err)
//	prettyJsonPrint(res)
// }

func TestDocumentRecommendList(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := articlev1.NewArticleClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.DocumentRecommendList(ctx, &articlev1.DocumentRecommendListReq{
		Uid:       2000009,
		SpaceGuid: "XLAYiB",
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)

}

func TestUpdateDocument(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := articlev1.NewArticleClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.UpdateDocument(ctx, &articlev1.UpdateDocumentReq{
		Uid:     1,
		Guid:    "RZzoPnDG5e",
		Name:    "hello world",
		Content: lo.ToPtr("我觉得还不错22222"),
	})
	assert.NoError(t, err)
	prettyJsonPrint(res)
}

func TestListDocument(t *testing.T) {
	// eapp.SetEgoDebug("true")
	// cli := articlev1.NewArticleClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// res, err := cli.ListDocument(ctx, &articlev1.ListDocumentReq{
	// 	Uid:        1,
	// 	CmtGuid:    "bPMDkrxK5x",
	// 	SpaceGuid:  "RLXJcB",
	// 	Pagination: &commonv1.Pagination{},
	// })
	// assert.NoError(t, err)
	// prettyJsonPrint(res)
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
