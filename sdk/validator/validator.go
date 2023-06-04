// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package validator

import (
	"context"
	"regexp"

	errcodev1 "ecodepost/pb/errcode/v1"
	"github.com/gotomicro/ego/core/eerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The validate interface starting with protoc-gen-validate v0.6.0.
// See https://github.com/envoyproxy/protoc-gen-validate/pull/455.
type validator interface {
	Validate(all bool) error
}

// The validate interface prior to protoc-gen-validate v0.6.0.
type validatorLegacy interface {
	Validate() error
}

var reg, _ = regexp.Compile(`error: code = InvalidArgument desc = invalid \w+\.(\w+): (.+)`)

const errStrEmpty = `value length must be at least 1 runes`
const errIntZero = `value must be greater than 0`
const errStrLength = `value length must be between 1 and 100 runes, inclusive`

// CastToEeerror 转换error为ego error
// TODO 修改validate插件
func CastToEeerror(err error) eerrors.Error {
	res := reg.FindStringSubmatch(err.Error())
	e := errcodev1.ErrInvalidArgument().WithMessage(err.Error())
	if len(res) != 3 {
		return e
	}
	switch res[1] {
	case "Uid", "OperateUid", "CreatedUid", "UpdatedUid", "ReferrerUid":
		if res[2] == errIntZero {
			return errcodev1.ErrUidEmpty()
		}
	case "SpaceGuid":
		if res[2] == errStrEmpty {
			return errcodev1.ErrSpaceEmpty()
		}
	case "FileGuid":
		if res[2] == errStrEmpty {
			return errcodev1.ErrFileGuidEmpty()
		}
	case "Guid":
		if res[2] == errStrEmpty {
			return errcodev1.ErrGuidEmpty()
		}
	case "Name":
		if res[2] == errStrEmpty {
			return errcodev1.ErrNameEmpty()
		}
	case "SpaceGroupGuid":
		if res[2] == errStrEmpty {
			return errcodev1.ErrSpaceGroupEmpty()
		}
	case "FileName":
		if res[2] == errStrLength {
			return errcodev1.ErrFileNameLength()
		}
	case "BizGuid":
		if res[2] == errStrEmpty {
			return errcodev1.ErrBizGuidEmpty()
		}
	case "Content":
		if res[2] == errStrEmpty {
			return errcodev1.ErrFileContentEmpty()
		}
	}
	return e
}

func validate(req interface{}) error {
	switch v := req.(type) {
	case validatorLegacy:
		if err := v.Validate(); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
		// case validator:
		// 	if err := v.Validate(false); err != nil {
		// 		return status.Error(codes.InvalidArgument, err.Error())
		// 	}
	}
	return nil
}

// UnaryServerInterceptor returns a new unary server interceptor that validates incoming messages.
//
// Invalid messages will be rejected with `InvalidArgument` before reaching any userspace handlers.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate(req); err != nil {
			return nil, CastToEeerror(err)
		}
		return handler(ctx, req)
	}
}

// UnaryClientInterceptor returns a new unary client interceptor that validates outgoing messages.
//
// Invalid messages will be rejected with `InvalidArgument` before sending the request to server.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := validate(req); err != nil {
			return err
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that validates incoming messages.
//
// The stage at which invalid messages will be rejected with `InvalidArgument` varies based on the
// type of the RPC. For `ServerStream` (1:m) requests, it will happen before reaching any userspace
// handlers. For `ClientStream` (n:1) or `BidiStream` (n:m) RPCs, the messages will be rejected on
// calls to `stream.Recv()`.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &recvWrapper{stream}
		return handler(srv, wrapper)
	}
}

type recvWrapper struct {
	grpc.ServerStream
}

func (s *recvWrapper) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	if err := validate(m); err != nil {
		return err
	}

	return nil
}
