package user

import (
	"time"

	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"google.golang.org/grpc"
)

func newCC() grpc.ClientConnInterface {
	return cegrpc.DefaultContainer().Build(
		cegrpc.WithBufnetServerListener(svc.Listener()),
		cegrpc.WithDialTimeout(5*time.Second),
		cegrpc.WithReadTimeout(5*time.Second),
	).ClientConn
}
