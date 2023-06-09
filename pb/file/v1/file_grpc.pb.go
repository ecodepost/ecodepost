// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package filev1

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

// FileClient is the client API for File service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileClient interface {
	// Emoji List
	EmojiList(ctx context.Context, in *EmojiListReq, opts ...grpc.CallOption) (*EmojiListRes, error)
	// 根据文章guid和用户信息，返回他的emoji
	MyEmojiList(ctx context.Context, in *MyEmojiListReq, opts ...grpc.CallOption) (*MyEmojiListRes, error)
	// 根据文章guid和用户信息，返回他的emoji
	MyEmojiListByFileGuids(ctx context.Context, in *MyEmojiListByFileGuidsReq, opts ...grpc.CallOption) (*MyEmojiListByFileGuidsRes, error)
	// Emoji Create
	CreateEmoji(ctx context.Context, in *CreateEmojiReq, opts ...grpc.CallOption) (*CreateEmojiRes, error)
	// Emoji Delete
	DeleteEmoji(ctx context.Context, in *DeleteEmojiReq, opts ...grpc.CallOption) (*DeleteEmojiRes, error)
	// 收藏某个目标到某几个收藏夹
	CollectionCreate(ctx context.Context, in *CollectionCreateReq, opts ...grpc.CallOption) (*CollectionCreateRes, error)
	// 从几个收藏夹取消收藏某个目标
	CollectionDelete(ctx context.Context, in *CollectionDeleteReq, opts ...grpc.CallOption) (*CollectionDeleteRes, error)
	// 直接更新File size
	UpdateFileSize(ctx context.Context, in *UpdateFileSizeReq, opts ...grpc.CallOption) (*UpdateFileSizeRes, error)
	// 获取文档内容的链接，前端读取文件内容渲染，有缓存，走CDN流量鉴权到OSS
	GetShowInfo(ctx context.Context, in *GetShowInfoReq, opts ...grpc.CallOption) (*GetShowInfoRes, error)
	// 单个file的permission
	Permission(ctx context.Context, in *PermissionReq, opts ...grpc.CallOption) (*PermissionRes, error)
	// 列表的permission
	PermissionList(ctx context.Context, in *PermissionListReq, opts ...grpc.CallOption) (*PermissionListRes, error)
	// List文章
	ListPage(ctx context.Context, in *ListPageReq, opts ...grpc.CallOption) (*ListPageRes, error)
	// List子文章
	ListPageByParent(ctx context.Context, in *ListPageByParentReq, opts ...grpc.CallOption) (*ListPageByParentRes, error)
	// 文章置顶列表
	SpaceTopList(ctx context.Context, in *SpaceTopListReq, opts ...grpc.CallOption) (*SpaceTopListRes, error)
	// 文章推荐列表
	RecommendList(ctx context.Context, in *RecommendListReq, opts ...grpc.CallOption) (*RecommendListRes, error)
	// 打开评论
	OpenComment(ctx context.Context, in *OpenCommentReq, opts ...grpc.CallOption) (*OpenCommentRes, error)
	// 关闭评论
	CloseComment(ctx context.Context, in *CloseCommentReq, opts ...grpc.CallOption) (*CloseCommentRes, error)
}

type fileClient struct {
	cc grpc.ClientConnInterface
}

func NewFileClient(cc grpc.ClientConnInterface) FileClient {
	return &fileClient{cc}
}

func (c *fileClient) EmojiList(ctx context.Context, in *EmojiListReq, opts ...grpc.CallOption) (*EmojiListRes, error) {
	out := new(EmojiListRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/EmojiList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) MyEmojiList(ctx context.Context, in *MyEmojiListReq, opts ...grpc.CallOption) (*MyEmojiListRes, error) {
	out := new(MyEmojiListRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/MyEmojiList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) MyEmojiListByFileGuids(ctx context.Context, in *MyEmojiListByFileGuidsReq, opts ...grpc.CallOption) (*MyEmojiListByFileGuidsRes, error) {
	out := new(MyEmojiListByFileGuidsRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/MyEmojiListByFileGuids", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) CreateEmoji(ctx context.Context, in *CreateEmojiReq, opts ...grpc.CallOption) (*CreateEmojiRes, error) {
	out := new(CreateEmojiRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/CreateEmoji", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) DeleteEmoji(ctx context.Context, in *DeleteEmojiReq, opts ...grpc.CallOption) (*DeleteEmojiRes, error) {
	out := new(DeleteEmojiRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/DeleteEmoji", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) CollectionCreate(ctx context.Context, in *CollectionCreateReq, opts ...grpc.CallOption) (*CollectionCreateRes, error) {
	out := new(CollectionCreateRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/CollectionCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) CollectionDelete(ctx context.Context, in *CollectionDeleteReq, opts ...grpc.CallOption) (*CollectionDeleteRes, error) {
	out := new(CollectionDeleteRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/CollectionDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) UpdateFileSize(ctx context.Context, in *UpdateFileSizeReq, opts ...grpc.CallOption) (*UpdateFileSizeRes, error) {
	out := new(UpdateFileSizeRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/UpdateFileSize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) GetShowInfo(ctx context.Context, in *GetShowInfoReq, opts ...grpc.CallOption) (*GetShowInfoRes, error) {
	out := new(GetShowInfoRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/GetShowInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) Permission(ctx context.Context, in *PermissionReq, opts ...grpc.CallOption) (*PermissionRes, error) {
	out := new(PermissionRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/Permission", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) PermissionList(ctx context.Context, in *PermissionListReq, opts ...grpc.CallOption) (*PermissionListRes, error) {
	out := new(PermissionListRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/PermissionList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) ListPage(ctx context.Context, in *ListPageReq, opts ...grpc.CallOption) (*ListPageRes, error) {
	out := new(ListPageRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/ListPage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) ListPageByParent(ctx context.Context, in *ListPageByParentReq, opts ...grpc.CallOption) (*ListPageByParentRes, error) {
	out := new(ListPageByParentRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/ListPageByParent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) SpaceTopList(ctx context.Context, in *SpaceTopListReq, opts ...grpc.CallOption) (*SpaceTopListRes, error) {
	out := new(SpaceTopListRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/SpaceTopList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) RecommendList(ctx context.Context, in *RecommendListReq, opts ...grpc.CallOption) (*RecommendListRes, error) {
	out := new(RecommendListRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/RecommendList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) OpenComment(ctx context.Context, in *OpenCommentReq, opts ...grpc.CallOption) (*OpenCommentRes, error) {
	out := new(OpenCommentRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/OpenComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) CloseComment(ctx context.Context, in *CloseCommentReq, opts ...grpc.CallOption) (*CloseCommentRes, error) {
	out := new(CloseCommentRes)
	err := c.cc.Invoke(ctx, "/file.v1.File/CloseComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileServer is the server API for File service.
// All implementations should embed UnimplementedFileServer
// for forward compatibility
type FileServer interface {
	// Emoji List
	EmojiList(context.Context, *EmojiListReq) (*EmojiListRes, error)
	// 根据文章guid和用户信息，返回他的emoji
	MyEmojiList(context.Context, *MyEmojiListReq) (*MyEmojiListRes, error)
	// 根据文章guid和用户信息，返回他的emoji
	MyEmojiListByFileGuids(context.Context, *MyEmojiListByFileGuidsReq) (*MyEmojiListByFileGuidsRes, error)
	// Emoji Create
	CreateEmoji(context.Context, *CreateEmojiReq) (*CreateEmojiRes, error)
	// Emoji Delete
	DeleteEmoji(context.Context, *DeleteEmojiReq) (*DeleteEmojiRes, error)
	// 收藏某个目标到某几个收藏夹
	CollectionCreate(context.Context, *CollectionCreateReq) (*CollectionCreateRes, error)
	// 从几个收藏夹取消收藏某个目标
	CollectionDelete(context.Context, *CollectionDeleteReq) (*CollectionDeleteRes, error)
	// 直接更新File size
	UpdateFileSize(context.Context, *UpdateFileSizeReq) (*UpdateFileSizeRes, error)
	// 获取文档内容的链接，前端读取文件内容渲染，有缓存，走CDN流量鉴权到OSS
	GetShowInfo(context.Context, *GetShowInfoReq) (*GetShowInfoRes, error)
	// 单个file的permission
	Permission(context.Context, *PermissionReq) (*PermissionRes, error)
	// 列表的permission
	PermissionList(context.Context, *PermissionListReq) (*PermissionListRes, error)
	// List文章
	ListPage(context.Context, *ListPageReq) (*ListPageRes, error)
	// List子文章
	ListPageByParent(context.Context, *ListPageByParentReq) (*ListPageByParentRes, error)
	// 文章置顶列表
	SpaceTopList(context.Context, *SpaceTopListReq) (*SpaceTopListRes, error)
	// 文章推荐列表
	RecommendList(context.Context, *RecommendListReq) (*RecommendListRes, error)
	// 打开评论
	OpenComment(context.Context, *OpenCommentReq) (*OpenCommentRes, error)
	// 关闭评论
	CloseComment(context.Context, *CloseCommentReq) (*CloseCommentRes, error)
}

// UnimplementedFileServer should be embedded to have forward compatible implementations.
type UnimplementedFileServer struct {
}

func (UnimplementedFileServer) EmojiList(context.Context, *EmojiListReq) (*EmojiListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmojiList not implemented")
}
func (UnimplementedFileServer) MyEmojiList(context.Context, *MyEmojiListReq) (*MyEmojiListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MyEmojiList not implemented")
}
func (UnimplementedFileServer) MyEmojiListByFileGuids(context.Context, *MyEmojiListByFileGuidsReq) (*MyEmojiListByFileGuidsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MyEmojiListByFileGuids not implemented")
}
func (UnimplementedFileServer) CreateEmoji(context.Context, *CreateEmojiReq) (*CreateEmojiRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEmoji not implemented")
}
func (UnimplementedFileServer) DeleteEmoji(context.Context, *DeleteEmojiReq) (*DeleteEmojiRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEmoji not implemented")
}
func (UnimplementedFileServer) CollectionCreate(context.Context, *CollectionCreateReq) (*CollectionCreateRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionCreate not implemented")
}
func (UnimplementedFileServer) CollectionDelete(context.Context, *CollectionDeleteReq) (*CollectionDeleteRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionDelete not implemented")
}
func (UnimplementedFileServer) UpdateFileSize(context.Context, *UpdateFileSizeReq) (*UpdateFileSizeRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFileSize not implemented")
}
func (UnimplementedFileServer) GetShowInfo(context.Context, *GetShowInfoReq) (*GetShowInfoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShowInfo not implemented")
}
func (UnimplementedFileServer) Permission(context.Context, *PermissionReq) (*PermissionRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Permission not implemented")
}
func (UnimplementedFileServer) PermissionList(context.Context, *PermissionListReq) (*PermissionListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PermissionList not implemented")
}
func (UnimplementedFileServer) ListPage(context.Context, *ListPageReq) (*ListPageRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPage not implemented")
}
func (UnimplementedFileServer) ListPageByParent(context.Context, *ListPageByParentReq) (*ListPageByParentRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPageByParent not implemented")
}
func (UnimplementedFileServer) SpaceTopList(context.Context, *SpaceTopListReq) (*SpaceTopListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SpaceTopList not implemented")
}
func (UnimplementedFileServer) RecommendList(context.Context, *RecommendListReq) (*RecommendListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecommendList not implemented")
}
func (UnimplementedFileServer) OpenComment(context.Context, *OpenCommentReq) (*OpenCommentRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenComment not implemented")
}
func (UnimplementedFileServer) CloseComment(context.Context, *CloseCommentReq) (*CloseCommentRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseComment not implemented")
}

// UnsafeFileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileServer will
// result in compilation errors.
type UnsafeFileServer interface {
	mustEmbedUnimplementedFileServer()
}

func RegisterFileServer(s grpc.ServiceRegistrar, srv FileServer) {
	s.RegisterService(&File_ServiceDesc, srv)
}

func _File_EmojiList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmojiListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).EmojiList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/EmojiList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).EmojiList(ctx, req.(*EmojiListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_MyEmojiList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MyEmojiListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).MyEmojiList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/MyEmojiList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).MyEmojiList(ctx, req.(*MyEmojiListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_MyEmojiListByFileGuids_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MyEmojiListByFileGuidsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).MyEmojiListByFileGuids(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/MyEmojiListByFileGuids",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).MyEmojiListByFileGuids(ctx, req.(*MyEmojiListByFileGuidsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_CreateEmoji_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEmojiReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).CreateEmoji(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/CreateEmoji",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).CreateEmoji(ctx, req.(*CreateEmojiReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_DeleteEmoji_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEmojiReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).DeleteEmoji(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/DeleteEmoji",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).DeleteEmoji(ctx, req.(*DeleteEmojiReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_CollectionCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).CollectionCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/CollectionCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).CollectionCreate(ctx, req.(*CollectionCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_CollectionDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionDeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).CollectionDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/CollectionDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).CollectionDelete(ctx, req.(*CollectionDeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_UpdateFileSize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFileSizeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).UpdateFileSize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/UpdateFileSize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).UpdateFileSize(ctx, req.(*UpdateFileSizeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_GetShowInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShowInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).GetShowInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/GetShowInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).GetShowInfo(ctx, req.(*GetShowInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_Permission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).Permission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/Permission",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).Permission(ctx, req.(*PermissionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_PermissionList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).PermissionList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/PermissionList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).PermissionList(ctx, req.(*PermissionListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_ListPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPageReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).ListPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/ListPage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).ListPage(ctx, req.(*ListPageReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_ListPageByParent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPageByParentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).ListPageByParent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/ListPageByParent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).ListPageByParent(ctx, req.(*ListPageByParentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_SpaceTopList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpaceTopListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).SpaceTopList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/SpaceTopList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).SpaceTopList(ctx, req.(*SpaceTopListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_RecommendList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecommendListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).RecommendList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/RecommendList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).RecommendList(ctx, req.(*RecommendListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_OpenComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenCommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).OpenComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/OpenComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).OpenComment(ctx, req.(*OpenCommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_CloseComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseCommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).CloseComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.v1.File/CloseComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).CloseComment(ctx, req.(*CloseCommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

// File_ServiceDesc is the grpc.ServiceDesc for File service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var File_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "file.v1.File",
	HandlerType: (*FileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EmojiList",
			Handler:    _File_EmojiList_Handler,
		},
		{
			MethodName: "MyEmojiList",
			Handler:    _File_MyEmojiList_Handler,
		},
		{
			MethodName: "MyEmojiListByFileGuids",
			Handler:    _File_MyEmojiListByFileGuids_Handler,
		},
		{
			MethodName: "CreateEmoji",
			Handler:    _File_CreateEmoji_Handler,
		},
		{
			MethodName: "DeleteEmoji",
			Handler:    _File_DeleteEmoji_Handler,
		},
		{
			MethodName: "CollectionCreate",
			Handler:    _File_CollectionCreate_Handler,
		},
		{
			MethodName: "CollectionDelete",
			Handler:    _File_CollectionDelete_Handler,
		},
		{
			MethodName: "UpdateFileSize",
			Handler:    _File_UpdateFileSize_Handler,
		},
		{
			MethodName: "GetShowInfo",
			Handler:    _File_GetShowInfo_Handler,
		},
		{
			MethodName: "Permission",
			Handler:    _File_Permission_Handler,
		},
		{
			MethodName: "PermissionList",
			Handler:    _File_PermissionList_Handler,
		},
		{
			MethodName: "ListPage",
			Handler:    _File_ListPage_Handler,
		},
		{
			MethodName: "ListPageByParent",
			Handler:    _File_ListPageByParent_Handler,
		},
		{
			MethodName: "SpaceTopList",
			Handler:    _File_SpaceTopList_Handler,
		},
		{
			MethodName: "RecommendList",
			Handler:    _File_RecommendList_Handler,
		},
		{
			MethodName: "OpenComment",
			Handler:    _File_OpenComment_Handler,
		},
		{
			MethodName: "CloseComment",
			Handler:    _File_CloseComment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "file/v1/file.proto",
}
