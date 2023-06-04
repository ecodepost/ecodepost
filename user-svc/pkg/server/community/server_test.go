//go:build microservice

package community

import (
	"context"
	"fmt"
	"testing"

	communityv1 "ecodepost/pb/community/v1"
	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/eerrors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProvideReferralRegisterCode(t *testing.T) {
	eapp.SetEgoDebug("true")
	cli := communityv1.NewCommunityClient(cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener())).ClientConn)
	ctx := context.Background()
	res, err := cli.Apply(ctx, &communityv1.ApplyReq{})
	spbStatus := Convert(err)
	egoErr := eerrors.FromError(err)
	fmt.Println(egoErr.GetMetadata()["showMsg"])
	fmt.Println(spbStatus.Details())
	assert.NoError(t, err)
	t.Logf("res: %+v", res)
}

// Convert 内部转换，为了让err=nil的时候，监控数据里有OK信息
func Convert(err error) *status.Status {
	if err == nil {
		return status.New(codes.OK, "OK")
	}

	if se, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return se.GRPCStatus()
	}

	switch err {
	case context.DeadlineExceeded:
		return status.New(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.New(codes.Canceled, err.Error())
	}

	return status.New(codes.Unknown, err.Error())
}
