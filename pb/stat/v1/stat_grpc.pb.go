// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package statv1

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

// StatClient is the client API for Stat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatClient interface {
	// 创建收藏夹
	CollectionGroupCreate(ctx context.Context, in *CollectionGroupCreateReq, opts ...grpc.CallOption) (*CollectionGroupCreateRes, error)
	// 获取收藏夹
	CollectionGroupList(ctx context.Context, in *CollectionGroupListReq, opts ...grpc.CallOption) (*CollectionGroupListRes, error)
	// 更新收藏夹
	CollectionGroupUpdate(ctx context.Context, in *CollectionGroupUpdateReq, opts ...grpc.CallOption) (*CollectionGroupUpdateRes, error)
	// 删除收藏夹
	CollectionGroupDelete(ctx context.Context, in *CollectionGroupDeleteReq, opts ...grpc.CallOption) (*CollectionGroupDeleteRes, error)
	// 收藏某个目标到某几个收藏夹
	CollectionCreate(ctx context.Context, in *CollectionCreateReq, opts ...grpc.CallOption) (*CollectionCreateRes, error)
	// 从几个收藏夹取消收藏某个目标
	CollectionDelete(ctx context.Context, in *CollectionDeleteReq, opts ...grpc.CallOption) (*CollectionDeleteRes, error)
	// 查看某个收藏夹收藏列表
	CollectionList(ctx context.Context, in *CollectionListReq, opts ...grpc.CallOption) (*CollectionListRes, error)
	// 根据文件GUIDS查看收藏列表
	MyCollectionListByFileGuids(ctx context.Context, in *MyCollectionListByFileGuidsReq, opts ...grpc.CallOption) (*MyCollectionListByFileGuidsRes, error)
	// 是否收藏
	IsCollection(ctx context.Context, in *IsCollectionReq, opts ...grpc.CallOption) (*IsCollectionRes, error)
}

type statClient struct {
	cc grpc.ClientConnInterface
}

func NewStatClient(cc grpc.ClientConnInterface) StatClient {
	return &statClient{cc}
}

func (c *statClient) CollectionGroupCreate(ctx context.Context, in *CollectionGroupCreateReq, opts ...grpc.CallOption) (*CollectionGroupCreateRes, error) {
	out := new(CollectionGroupCreateRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/CollectionGroupCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statClient) CollectionGroupList(ctx context.Context, in *CollectionGroupListReq, opts ...grpc.CallOption) (*CollectionGroupListRes, error) {
	out := new(CollectionGroupListRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/CollectionGroupList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statClient) CollectionGroupUpdate(ctx context.Context, in *CollectionGroupUpdateReq, opts ...grpc.CallOption) (*CollectionGroupUpdateRes, error) {
	out := new(CollectionGroupUpdateRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/CollectionGroupUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statClient) CollectionGroupDelete(ctx context.Context, in *CollectionGroupDeleteReq, opts ...grpc.CallOption) (*CollectionGroupDeleteRes, error) {
	out := new(CollectionGroupDeleteRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/CollectionGroupDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statClient) CollectionCreate(ctx context.Context, in *CollectionCreateReq, opts ...grpc.CallOption) (*CollectionCreateRes, error) {
	out := new(CollectionCreateRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/CollectionCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statClient) CollectionDelete(ctx context.Context, in *CollectionDeleteReq, opts ...grpc.CallOption) (*CollectionDeleteRes, error) {
	out := new(CollectionDeleteRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/CollectionDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statClient) CollectionList(ctx context.Context, in *CollectionListReq, opts ...grpc.CallOption) (*CollectionListRes, error) {
	out := new(CollectionListRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/CollectionList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statClient) MyCollectionListByFileGuids(ctx context.Context, in *MyCollectionListByFileGuidsReq, opts ...grpc.CallOption) (*MyCollectionListByFileGuidsRes, error) {
	out := new(MyCollectionListByFileGuidsRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/MyCollectionListByFileGuids", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statClient) IsCollection(ctx context.Context, in *IsCollectionReq, opts ...grpc.CallOption) (*IsCollectionRes, error) {
	out := new(IsCollectionRes)
	err := c.cc.Invoke(ctx, "/stat.v1.Stat/IsCollection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatServer is the server API for Stat service.
// All implementations should embed UnimplementedStatServer
// for forward compatibility
type StatServer interface {
	// 创建收藏夹
	CollectionGroupCreate(context.Context, *CollectionGroupCreateReq) (*CollectionGroupCreateRes, error)
	// 获取收藏夹
	CollectionGroupList(context.Context, *CollectionGroupListReq) (*CollectionGroupListRes, error)
	// 更新收藏夹
	CollectionGroupUpdate(context.Context, *CollectionGroupUpdateReq) (*CollectionGroupUpdateRes, error)
	// 删除收藏夹
	CollectionGroupDelete(context.Context, *CollectionGroupDeleteReq) (*CollectionGroupDeleteRes, error)
	// 收藏某个目标到某几个收藏夹
	CollectionCreate(context.Context, *CollectionCreateReq) (*CollectionCreateRes, error)
	// 从几个收藏夹取消收藏某个目标
	CollectionDelete(context.Context, *CollectionDeleteReq) (*CollectionDeleteRes, error)
	// 查看某个收藏夹收藏列表
	CollectionList(context.Context, *CollectionListReq) (*CollectionListRes, error)
	// 根据文件GUIDS查看收藏列表
	MyCollectionListByFileGuids(context.Context, *MyCollectionListByFileGuidsReq) (*MyCollectionListByFileGuidsRes, error)
	// 是否收藏
	IsCollection(context.Context, *IsCollectionReq) (*IsCollectionRes, error)
}

// UnimplementedStatServer should be embedded to have forward compatible implementations.
type UnimplementedStatServer struct {
}

func (UnimplementedStatServer) CollectionGroupCreate(context.Context, *CollectionGroupCreateReq) (*CollectionGroupCreateRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionGroupCreate not implemented")
}
func (UnimplementedStatServer) CollectionGroupList(context.Context, *CollectionGroupListReq) (*CollectionGroupListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionGroupList not implemented")
}
func (UnimplementedStatServer) CollectionGroupUpdate(context.Context, *CollectionGroupUpdateReq) (*CollectionGroupUpdateRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionGroupUpdate not implemented")
}
func (UnimplementedStatServer) CollectionGroupDelete(context.Context, *CollectionGroupDeleteReq) (*CollectionGroupDeleteRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionGroupDelete not implemented")
}
func (UnimplementedStatServer) CollectionCreate(context.Context, *CollectionCreateReq) (*CollectionCreateRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionCreate not implemented")
}
func (UnimplementedStatServer) CollectionDelete(context.Context, *CollectionDeleteReq) (*CollectionDeleteRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionDelete not implemented")
}
func (UnimplementedStatServer) CollectionList(context.Context, *CollectionListReq) (*CollectionListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionList not implemented")
}
func (UnimplementedStatServer) MyCollectionListByFileGuids(context.Context, *MyCollectionListByFileGuidsReq) (*MyCollectionListByFileGuidsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MyCollectionListByFileGuids not implemented")
}
func (UnimplementedStatServer) IsCollection(context.Context, *IsCollectionReq) (*IsCollectionRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsCollection not implemented")
}

// UnsafeStatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatServer will
// result in compilation errors.
type UnsafeStatServer interface {
	mustEmbedUnimplementedStatServer()
}

func RegisterStatServer(s grpc.ServiceRegistrar, srv StatServer) {
	s.RegisterService(&Stat_ServiceDesc, srv)
}

func _Stat_CollectionGroupCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionGroupCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).CollectionGroupCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/CollectionGroupCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).CollectionGroupCreate(ctx, req.(*CollectionGroupCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stat_CollectionGroupList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionGroupListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).CollectionGroupList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/CollectionGroupList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).CollectionGroupList(ctx, req.(*CollectionGroupListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stat_CollectionGroupUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionGroupUpdateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).CollectionGroupUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/CollectionGroupUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).CollectionGroupUpdate(ctx, req.(*CollectionGroupUpdateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stat_CollectionGroupDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionGroupDeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).CollectionGroupDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/CollectionGroupDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).CollectionGroupDelete(ctx, req.(*CollectionGroupDeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stat_CollectionCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).CollectionCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/CollectionCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).CollectionCreate(ctx, req.(*CollectionCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stat_CollectionDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionDeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).CollectionDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/CollectionDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).CollectionDelete(ctx, req.(*CollectionDeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stat_CollectionList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).CollectionList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/CollectionList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).CollectionList(ctx, req.(*CollectionListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stat_MyCollectionListByFileGuids_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MyCollectionListByFileGuidsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).MyCollectionListByFileGuids(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/MyCollectionListByFileGuids",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).MyCollectionListByFileGuids(ctx, req.(*MyCollectionListByFileGuidsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stat_IsCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsCollectionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatServer).IsCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stat.v1.Stat/IsCollection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatServer).IsCollection(ctx, req.(*IsCollectionReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Stat_ServiceDesc is the grpc.ServiceDesc for Stat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stat.v1.Stat",
	HandlerType: (*StatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CollectionGroupCreate",
			Handler:    _Stat_CollectionGroupCreate_Handler,
		},
		{
			MethodName: "CollectionGroupList",
			Handler:    _Stat_CollectionGroupList_Handler,
		},
		{
			MethodName: "CollectionGroupUpdate",
			Handler:    _Stat_CollectionGroupUpdate_Handler,
		},
		{
			MethodName: "CollectionGroupDelete",
			Handler:    _Stat_CollectionGroupDelete_Handler,
		},
		{
			MethodName: "CollectionCreate",
			Handler:    _Stat_CollectionCreate_Handler,
		},
		{
			MethodName: "CollectionDelete",
			Handler:    _Stat_CollectionDelete_Handler,
		},
		{
			MethodName: "CollectionList",
			Handler:    _Stat_CollectionList_Handler,
		},
		{
			MethodName: "MyCollectionListByFileGuids",
			Handler:    _Stat_MyCollectionListByFileGuids_Handler,
		},
		{
			MethodName: "IsCollection",
			Handler:    _Stat_IsCollection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stat/v1/stat.proto",
}
