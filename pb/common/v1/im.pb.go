// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: common/v1/im.proto

package commonv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// IM群组信息
type ImGroupInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 群组名称
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// 群组头像
	Avatar string `protobuf:"bytes,2,opt,name=avatar,proto3" json:"avatar,omitempty"`
	// 描述
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// 类型
	Type int32 `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"`
	// 需要加入社区的uid列表
	Uids []int64 `protobuf:"varint,5,rep,packed,name=uids,proto3" json:"uids,omitempty"`
	// 空间guid
	SpaceGuid string `protobuf:"bytes,7,opt,name=spaceGuid,proto3" json:"spaceGuid,omitempty"`
	// 群主uid(我方)
	OwnerUid int64 `protobuf:"varint,8,opt,name=ownerUid,proto3" json:"ownerUid,omitempty"`
	// 群组id(蓝莺)
	GroupId int64 `protobuf:"varint,9,opt,name=groupId,proto3" json:"groupId,omitempty"`
	// 创建时间
	Ctime int64 `protobuf:"varint,10,opt,name=ctime,proto3" json:"ctime,omitempty"`
	// 更新时间
	Utime int64 `protobuf:"varint,11,opt,name=utime,proto3" json:"utime,omitempty"`
}

func (x *ImGroupInfo) Reset() {
	*x = ImGroupInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_v1_im_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImGroupInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImGroupInfo) ProtoMessage() {}

func (x *ImGroupInfo) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1_im_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImGroupInfo.ProtoReflect.Descriptor instead.
func (*ImGroupInfo) Descriptor() ([]byte, []int) {
	return file_common_v1_im_proto_rawDescGZIP(), []int{0}
}

func (x *ImGroupInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ImGroupInfo) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *ImGroupInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ImGroupInfo) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *ImGroupInfo) GetUids() []int64 {
	if x != nil {
		return x.Uids
	}
	return nil
}

func (x *ImGroupInfo) GetSpaceGuid() string {
	if x != nil {
		return x.SpaceGuid
	}
	return ""
}

func (x *ImGroupInfo) GetOwnerUid() int64 {
	if x != nil {
		return x.OwnerUid
	}
	return 0
}

func (x *ImGroupInfo) GetGroupId() int64 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *ImGroupInfo) GetCtime() int64 {
	if x != nil {
		return x.Ctime
	}
	return 0
}

func (x *ImGroupInfo) GetUtime() int64 {
	if x != nil {
		return x.Utime
	}
	return 0
}

var File_common_v1_im_proto protoreflect.FileDescriptor

var file_common_v1_im_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x22,
	0x83, 0x02, 0x0a, 0x0b, 0x49, 0x6d, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x69, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52,
	0x04, 0x75, 0x69, 0x64, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x70, 0x61, 0x63, 0x65, 0x47, 0x75,
	0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x70, 0x61, 0x63, 0x65, 0x47,
	0x75, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x55, 0x69, 0x64, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x55, 0x69, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x75, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x75, 0x74, 0x69, 0x6d, 0x65, 0x42, 0x7e, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x07, 0x49, 0x6d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x1f, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x74, 0x2f, 0x70, 0x62, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x15, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_v1_im_proto_rawDescOnce sync.Once
	file_common_v1_im_proto_rawDescData = file_common_v1_im_proto_rawDesc
)

func file_common_v1_im_proto_rawDescGZIP() []byte {
	file_common_v1_im_proto_rawDescOnce.Do(func() {
		file_common_v1_im_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_v1_im_proto_rawDescData)
	})
	return file_common_v1_im_proto_rawDescData
}

var file_common_v1_im_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_v1_im_proto_goTypes = []interface{}{
	(*ImGroupInfo)(nil), // 0: common.v1.ImGroupInfo
}
var file_common_v1_im_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_v1_im_proto_init() }
func file_common_v1_im_proto_init() {
	if File_common_v1_im_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_v1_im_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImGroupInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_v1_im_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_v1_im_proto_goTypes,
		DependencyIndexes: file_common_v1_im_proto_depIdxs,
		MessageInfos:      file_common_v1_im_proto_msgTypes,
	}.Build()
	File_common_v1_im_proto = out.File
	file_common_v1_im_proto_rawDesc = nil
	file_common_v1_im_proto_goTypes = nil
	file_common_v1_im_proto_depIdxs = nil
}
