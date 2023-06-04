// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package uploadv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UploadClient is the client API for Upload service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UploadClient interface {
	// 获取一次性上传文件的token
	GetToken(ctx context.Context, in *GetTokenReq, opts ...grpc.CallOption) (*GetTokenRes, error)
	// 获取上传文件的Path
	GetPath(ctx context.Context, in *GetPathReq, opts ...grpc.CallOption) (*GetPathRes, error)
	// 上传本地文件
	UploadLocalFile(ctx context.Context, in *UploadLocalFileReq, opts ...grpc.CallOption) (*UploadLocalFileRes, error)
	// 展示图片
	ShowImage(ctx context.Context, in *ShowImageReq, opts ...grpc.CallOption) (*ShowImageRes, error)
}

type uploadClient struct {
	cc grpc.ClientConnInterface
}

func NewUploadClient(cc grpc.ClientConnInterface) UploadClient {
	return &uploadClient{cc}
}

func (c *uploadClient) GetToken(ctx context.Context, in *GetTokenReq, opts ...grpc.CallOption) (*GetTokenRes, error) {
	out := new(GetTokenRes)
	err := c.cc.Invoke(ctx, "/upload.v1.Upload/GetToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadClient) GetPath(ctx context.Context, in *GetPathReq, opts ...grpc.CallOption) (*GetPathRes, error) {
	out := new(GetPathRes)
	err := c.cc.Invoke(ctx, "/upload.v1.Upload/GetPath", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadClient) UploadLocalFile(ctx context.Context, in *UploadLocalFileReq, opts ...grpc.CallOption) (*UploadLocalFileRes, error) {
	out := new(UploadLocalFileRes)
	err := c.cc.Invoke(ctx, "/upload.v1.Upload/UploadLocalFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadClient) ShowImage(ctx context.Context, in *ShowImageReq, opts ...grpc.CallOption) (*ShowImageRes, error) {
	out := new(ShowImageRes)
	err := c.cc.Invoke(ctx, "/upload.v1.Upload/ShowImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UploadServer is the server API for Upload service.
// All implementations should embed UnimplementedUploadServer
// for forward compatibility
type UploadServer interface {
	// 获取一次性上传文件的token
	GetToken(context.Context, *GetTokenReq) (*GetTokenRes, error)
	// 获取上传文件的Path
	GetPath(context.Context, *GetPathReq) (*GetPathRes, error)
	// 上传本地文件
	UploadLocalFile(context.Context, *UploadLocalFileReq) (*UploadLocalFileRes, error)
	// 展示图片
	ShowImage(context.Context, *ShowImageReq) (*ShowImageRes, error)
}

// UnimplementedUploadServer should be embedded to have forward compatible implementations.
type UnimplementedUploadServer struct {
}

func (UnimplementedUploadServer) GetToken(context.Context, *GetTokenReq) (*GetTokenRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetToken not implemented")
}
func (UnimplementedUploadServer) GetPath(context.Context, *GetPathReq) (*GetPathRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPath not implemented")
}
func (UnimplementedUploadServer) UploadLocalFile(context.Context, *UploadLocalFileReq) (*UploadLocalFileRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadLocalFile not implemented")
}
func (UnimplementedUploadServer) ShowImage(context.Context, *ShowImageReq) (*ShowImageRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShowImage not implemented")
}

// UnsafeUploadServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UploadServer will
// result in compilation errors.
type UnsafeUploadServer interface {
	mustEmbedUnimplementedUploadServer()
}

func RegisterUploadServer(s grpc.ServiceRegistrar, srv UploadServer) {
	s.RegisterService(&Upload_ServiceDesc, srv)
}

func _Upload_GetToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadServer).GetToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/upload.v1.Upload/GetToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadServer).GetToken(ctx, req.(*GetTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Upload_GetPath_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPathReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadServer).GetPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/upload.v1.Upload/GetPath",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadServer).GetPath(ctx, req.(*GetPathReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Upload_UploadLocalFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadLocalFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadServer).UploadLocalFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/upload.v1.Upload/UploadLocalFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadServer).UploadLocalFile(ctx, req.(*UploadLocalFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Upload_ShowImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShowImageReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadServer).ShowImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/upload.v1.Upload/ShowImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadServer).ShowImage(ctx, req.(*ShowImageReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Upload_ServiceDesc is the grpc.ServiceDesc for Upload service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Upload_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "upload.v1.Upload",
	HandlerType: (*UploadServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetToken",
			Handler:    _Upload_GetToken_Handler,
		},
		{
			MethodName: "GetPath",
			Handler:    _Upload_GetPath_Handler,
		},
		{
			MethodName: "UploadLocalFile",
			Handler:    _Upload_UploadLocalFile_Handler,
		},
		{
			MethodName: "ShowImage",
			Handler:    _Upload_ShowImage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "upload/v1/upload.proto",
}
