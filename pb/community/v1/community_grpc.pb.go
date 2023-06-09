// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package communityv1

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

// CommunityClient is the client API for Community service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunityClient interface {
	// 社区首页信息
	Home(ctx context.Context, in *HomeReq, opts ...grpc.CallOption) (*HomeRes, error)
	// 获取社区主题
	GetTheme(ctx context.Context, in *GetThemeReq, opts ...grpc.CallOption) (*GetThemeRes, error)
	// 设置社区主题
	SetTheme(ctx context.Context, in *SetThemeReq, opts ...grpc.CallOption) (*SetThemeRes, error)
	// 社区首页可选项信息
	GetHomeOption(ctx context.Context, in *GetHomeOptionReq, opts ...grpc.CallOption) (*GetHomeOptionRes, error)
	// 更新首页可选项信息
	PutHomeOption(ctx context.Context, in *PutHomeOptionReq, opts ...grpc.CallOption) (*PutHomeOptionRes, error)
	// 获取社区信息
	Info(ctx context.Context, in *InfoReq, opts ...grpc.CallOption) (*InfoRes, error)
	// 修改社区信息
	Update(ctx context.Context, in *UpdateReq, opts ...grpc.CallOption) (*UpdateRes, error)
}

type communityClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunityClient(cc grpc.ClientConnInterface) CommunityClient {
	return &communityClient{cc}
}

func (c *communityClient) Home(ctx context.Context, in *HomeReq, opts ...grpc.CallOption) (*HomeRes, error) {
	out := new(HomeRes)
	err := c.cc.Invoke(ctx, "/community.v1.Community/Home", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communityClient) GetTheme(ctx context.Context, in *GetThemeReq, opts ...grpc.CallOption) (*GetThemeRes, error) {
	out := new(GetThemeRes)
	err := c.cc.Invoke(ctx, "/community.v1.Community/GetTheme", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communityClient) SetTheme(ctx context.Context, in *SetThemeReq, opts ...grpc.CallOption) (*SetThemeRes, error) {
	out := new(SetThemeRes)
	err := c.cc.Invoke(ctx, "/community.v1.Community/SetTheme", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communityClient) GetHomeOption(ctx context.Context, in *GetHomeOptionReq, opts ...grpc.CallOption) (*GetHomeOptionRes, error) {
	out := new(GetHomeOptionRes)
	err := c.cc.Invoke(ctx, "/community.v1.Community/GetHomeOption", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communityClient) PutHomeOption(ctx context.Context, in *PutHomeOptionReq, opts ...grpc.CallOption) (*PutHomeOptionRes, error) {
	out := new(PutHomeOptionRes)
	err := c.cc.Invoke(ctx, "/community.v1.Community/PutHomeOption", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communityClient) Info(ctx context.Context, in *InfoReq, opts ...grpc.CallOption) (*InfoRes, error) {
	out := new(InfoRes)
	err := c.cc.Invoke(ctx, "/community.v1.Community/Info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communityClient) Update(ctx context.Context, in *UpdateReq, opts ...grpc.CallOption) (*UpdateRes, error) {
	out := new(UpdateRes)
	err := c.cc.Invoke(ctx, "/community.v1.Community/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunityServer is the server API for Community service.
// All implementations should embed UnimplementedCommunityServer
// for forward compatibility
type CommunityServer interface {
	// 社区首页信息
	Home(context.Context, *HomeReq) (*HomeRes, error)
	// 获取社区主题
	GetTheme(context.Context, *GetThemeReq) (*GetThemeRes, error)
	// 设置社区主题
	SetTheme(context.Context, *SetThemeReq) (*SetThemeRes, error)
	// 社区首页可选项信息
	GetHomeOption(context.Context, *GetHomeOptionReq) (*GetHomeOptionRes, error)
	// 更新首页可选项信息
	PutHomeOption(context.Context, *PutHomeOptionReq) (*PutHomeOptionRes, error)
	// 获取社区信息
	Info(context.Context, *InfoReq) (*InfoRes, error)
	// 修改社区信息
	Update(context.Context, *UpdateReq) (*UpdateRes, error)
}

// UnimplementedCommunityServer should be embedded to have forward compatible implementations.
type UnimplementedCommunityServer struct {
}

func (UnimplementedCommunityServer) Home(context.Context, *HomeReq) (*HomeRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Home not implemented")
}
func (UnimplementedCommunityServer) GetTheme(context.Context, *GetThemeReq) (*GetThemeRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTheme not implemented")
}
func (UnimplementedCommunityServer) SetTheme(context.Context, *SetThemeReq) (*SetThemeRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetTheme not implemented")
}
func (UnimplementedCommunityServer) GetHomeOption(context.Context, *GetHomeOptionReq) (*GetHomeOptionRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHomeOption not implemented")
}
func (UnimplementedCommunityServer) PutHomeOption(context.Context, *PutHomeOptionReq) (*PutHomeOptionRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutHomeOption not implemented")
}
func (UnimplementedCommunityServer) Info(context.Context, *InfoReq) (*InfoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedCommunityServer) Update(context.Context, *UpdateReq) (*UpdateRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

// UnsafeCommunityServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunityServer will
// result in compilation errors.
type UnsafeCommunityServer interface {
	mustEmbedUnimplementedCommunityServer()
}

func RegisterCommunityServer(s grpc.ServiceRegistrar, srv CommunityServer) {
	s.RegisterService(&Community_ServiceDesc, srv)
}

func _Community_Home_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HomeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServer).Home(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/community.v1.Community/Home",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServer).Home(ctx, req.(*HomeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Community_GetTheme_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetThemeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServer).GetTheme(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/community.v1.Community/GetTheme",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServer).GetTheme(ctx, req.(*GetThemeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Community_SetTheme_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetThemeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServer).SetTheme(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/community.v1.Community/SetTheme",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServer).SetTheme(ctx, req.(*SetThemeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Community_GetHomeOption_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHomeOptionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServer).GetHomeOption(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/community.v1.Community/GetHomeOption",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServer).GetHomeOption(ctx, req.(*GetHomeOptionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Community_PutHomeOption_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutHomeOptionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServer).PutHomeOption(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/community.v1.Community/PutHomeOption",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServer).PutHomeOption(ctx, req.(*PutHomeOptionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Community_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/community.v1.Community/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServer).Info(ctx, req.(*InfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Community_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/community.v1.Community/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServer).Update(ctx, req.(*UpdateReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Community_ServiceDesc is the grpc.ServiceDesc for Community service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Community_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "community.v1.Community",
	HandlerType: (*CommunityServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Home",
			Handler:    _Community_Home_Handler,
		},
		{
			MethodName: "GetTheme",
			Handler:    _Community_GetTheme_Handler,
		},
		{
			MethodName: "SetTheme",
			Handler:    _Community_SetTheme_Handler,
		},
		{
			MethodName: "GetHomeOption",
			Handler:    _Community_GetHomeOption_Handler,
		},
		{
			MethodName: "PutHomeOption",
			Handler:    _Community_PutHomeOption_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _Community_Info_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Community_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "community/v1/community.proto",
}
