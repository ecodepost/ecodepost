package pmspolicy

import (
	"context"

	"google.golang.org/grpc"
)

const (
	XAuthToken = "x-auth-token"
	ClientID   = "x-client-id"
)

var ignoreMethods = []string{
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
	"/grpc.health.v1.Health/Check",
}

func checkMethodIsIgnored(method string) bool {
	for _, v := range ignoreMethods {
		if v == method {
			return true
		}
	}
	return false
}

// UnaryServerInterceptor 服务端Unary拦截器
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		// 如果是reflection请求则不鉴权
		if checkMethodIsIgnored(info.FullMethod) {
			return handler(ctx, req)
		}

		return handler(ctx, req)
	}
}
