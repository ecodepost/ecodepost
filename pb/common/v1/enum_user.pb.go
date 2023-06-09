// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: common/v1/enum_user.proto

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

// STATUS:Status 用户状态
type USER_STATUS int32

const (
	// 无效枚举值
	USER_STATUS_INVALID USER_STATUS = 0
	// 激活状态
	USER_STATUS_ACTIVE USER_STATUS = 1
	// 禁用状态
	USER_STATUS_BAN USER_STATUS = 2
)

// Enum value maps for USER_STATUS.
var (
	USER_STATUS_name = map[int32]string{
		0: "STATUS_INVALID",
		1: "STATUS_ACTIVE",
		2: "STATUS_BAN",
	}
	USER_STATUS_value = map[string]int32{
		"STATUS_INVALID": 0,
		"STATUS_ACTIVE":  1,
		"STATUS_BAN":     2,
	}
)

func (x USER_STATUS) Enum() *USER_STATUS {
	p := new(USER_STATUS)
	*p = x
	return p
}

func (x USER_STATUS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (USER_STATUS) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_user_proto_enumTypes[0].Descriptor()
}

func (USER_STATUS) Type() protoreflect.EnumType {
	return &file_common_v1_enum_user_proto_enumTypes[0]
}

func (x USER_STATUS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use USER_STATUS.Descriptor instead.
func (USER_STATUS) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_user_proto_rawDescGZIP(), []int{0, 0}
}

// AS:AuthStatus 用户身份认证状态
type USER_AS int32

const (
	// 无效枚举值
	USER_AS_INVALID USER_AS = 0
	// 未认证
	USER_AS_NOT USER_AS = 1
	// 初级认证
	USER_AS_PRIMARY USER_AS = 2
	// 高级认证
	USER_AS_HIGH USER_AS = 3
)

// Enum value maps for USER_AS.
var (
	USER_AS_name = map[int32]string{
		0: "AS_INVALID",
		1: "AS_NOT",
		2: "AS_PRIMARY",
		3: "AS_HIGH",
	}
	USER_AS_value = map[string]int32{
		"AS_INVALID": 0,
		"AS_NOT":     1,
		"AS_PRIMARY": 2,
		"AS_HIGH":    3,
	}
)

func (x USER_AS) Enum() *USER_AS {
	p := new(USER_AS)
	*p = x
	return p
}

func (x USER_AS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (USER_AS) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_user_proto_enumTypes[1].Descriptor()
}

func (USER_AS) Type() protoreflect.EnumType {
	return &file_common_v1_enum_user_proto_enumTypes[1]
}

func (x USER_AS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use USER_AS.Descriptor instead.
func (USER_AS) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_user_proto_rawDescGZIP(), []int{0, 1}
}

// EBS:EmailBindStatus 邮箱绑定状态
type USER_EBS int32

const (
	// 无效枚举值
	USER_EBS_INVALID USER_EBS = 0
	// 待确认
	USER_EBS_TO_BE_CONFIRMED USER_EBS = 1
	// 已确认
	USER_EBS_CONFIRMED USER_EBS = 2
)

// Enum value maps for USER_EBS.
var (
	USER_EBS_name = map[int32]string{
		0: "EBS_INVALID",
		1: "EBS_TO_BE_CONFIRMED",
		2: "EBS_CONFIRMED",
	}
	USER_EBS_value = map[string]int32{
		"EBS_INVALID":         0,
		"EBS_TO_BE_CONFIRMED": 1,
		"EBS_CONFIRMED":       2,
	}
)

func (x USER_EBS) Enum() *USER_EBS {
	p := new(USER_EBS)
	*p = x
	return p
}

func (x USER_EBS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (USER_EBS) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_user_proto_enumTypes[2].Descriptor()
}

func (USER_EBS) Type() protoreflect.EnumType {
	return &file_common_v1_enum_user_proto_enumTypes[2]
}

func (x USER_EBS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use USER_EBS.Descriptor instead.
func (USER_EBS) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_user_proto_rawDescGZIP(), []int{0, 2}
}

// ETT:EmailTokenType 邮件验证码类型
type USER_ETT int32

const (
	// 未知类型
	USER_ETT_INVALID USER_ETT = 0
	// 绑定验证
	USER_ETT_BIND USER_ETT = 1
	// 忘记密码
	USER_ETT_FORGET_PASSWORD USER_ETT = 2
	// 修改密码
	USER_ETT_CHANGE_PASSWORD USER_ETT = 3
	// 修改邮箱
	USER_ETT_CHANGE_EMAIL USER_ETT = 4
)

// Enum value maps for USER_ETT.
var (
	USER_ETT_name = map[int32]string{
		0: "ETT_INVALID",
		1: "ETT_BIND",
		2: "ETT_FORGET_PASSWORD",
		3: "ETT_CHANGE_PASSWORD",
		4: "ETT_CHANGE_EMAIL",
	}
	USER_ETT_value = map[string]int32{
		"ETT_INVALID":         0,
		"ETT_BIND":            1,
		"ETT_FORGET_PASSWORD": 2,
		"ETT_CHANGE_PASSWORD": 3,
		"ETT_CHANGE_EMAIL":    4,
	}
)

func (x USER_ETT) Enum() *USER_ETT {
	p := new(USER_ETT)
	*p = x
	return p
}

func (x USER_ETT) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (USER_ETT) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_user_proto_enumTypes[3].Descriptor()
}

func (USER_ETT) Type() protoreflect.EnumType {
	return &file_common_v1_enum_user_proto_enumTypes[3]
}

func (x USER_ETT) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use USER_ETT.Descriptor instead.
func (USER_ETT) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_user_proto_rawDescGZIP(), []int{0, 3}
}

// U2U_STAT_TYPE:U2UStatType 用户和其他用户关联状态类型
type USER_U2U_STAT_TYPE int32

const (
	// 无效枚举值
	USER_U2U_STAT_TYPE_INVALID USER_U2U_STAT_TYPE = 0
	// 是否关注
	USER_U2U_STAT_TYPE_FOLLOW USER_U2U_STAT_TYPE = 1
	// 是否拉黑
	USER_U2U_STAT_TYPE_BLOCK USER_U2U_STAT_TYPE = 2
)

// Enum value maps for USER_U2U_STAT_TYPE.
var (
	USER_U2U_STAT_TYPE_name = map[int32]string{
		0: "U2U_STAT_TYPE_INVALID",
		1: "U2U_STAT_TYPE_FOLLOW",
		2: "U2U_STAT_TYPE_BLOCK",
	}
	USER_U2U_STAT_TYPE_value = map[string]int32{
		"U2U_STAT_TYPE_INVALID": 0,
		"U2U_STAT_TYPE_FOLLOW":  1,
		"U2U_STAT_TYPE_BLOCK":   2,
	}
)

func (x USER_U2U_STAT_TYPE) Enum() *USER_U2U_STAT_TYPE {
	p := new(USER_U2U_STAT_TYPE)
	*p = x
	return p
}

func (x USER_U2U_STAT_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (USER_U2U_STAT_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_user_proto_enumTypes[4].Descriptor()
}

func (USER_U2U_STAT_TYPE) Type() protoreflect.EnumType {
	return &file_common_v1_enum_user_proto_enumTypes[4]
}

func (x USER_U2U_STAT_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use USER_U2U_STAT_TYPE.Descriptor instead.
func (USER_U2U_STAT_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_user_proto_rawDescGZIP(), []int{0, 4}
}

// USER:User 用户相关枚举值
type USER struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *USER) Reset() {
	*x = USER{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_v1_enum_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *USER) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*USER) ProtoMessage() {}

func (x *USER) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1_enum_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use USER.ProtoReflect.Descriptor instead.
func (*USER) Descriptor() ([]byte, []int) {
	return file_common_v1_enum_user_proto_rawDescGZIP(), []int{0}
}

var File_common_v1_enum_user_proto protoreflect.FileDescriptor

var file_common_v1_enum_user_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x75, 0x6d,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x22, 0x97, 0x03, 0x0a, 0x04, 0x55, 0x53, 0x45, 0x52, 0x22,
	0x3f, 0x0a, 0x06, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x11, 0x0a,
	0x0d, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01,
	0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x42, 0x41, 0x4e, 0x10, 0x02,
	0x22, 0x3d, 0x0a, 0x02, 0x41, 0x53, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x53, 0x5f, 0x49, 0x4e, 0x56,
	0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x53, 0x5f, 0x4e, 0x4f, 0x54,
	0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x53, 0x5f, 0x50, 0x52, 0x49, 0x4d, 0x41, 0x52, 0x59,
	0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x53, 0x5f, 0x48, 0x49, 0x47, 0x48, 0x10, 0x03, 0x22,
	0x42, 0x0a, 0x03, 0x45, 0x42, 0x53, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x42, 0x53, 0x5f, 0x49, 0x4e,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x42, 0x53, 0x5f, 0x54,
	0x4f, 0x5f, 0x42, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x52, 0x4d, 0x45, 0x44, 0x10, 0x01,
	0x12, 0x11, 0x0a, 0x0d, 0x45, 0x42, 0x53, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x52, 0x4d, 0x45,
	0x44, 0x10, 0x02, 0x22, 0x6c, 0x0a, 0x03, 0x45, 0x54, 0x54, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x54,
	0x54, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x45,
	0x54, 0x54, 0x5f, 0x42, 0x49, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x54, 0x54,
	0x5f, 0x46, 0x4f, 0x52, 0x47, 0x45, 0x54, 0x5f, 0x50, 0x41, 0x53, 0x53, 0x57, 0x4f, 0x52, 0x44,
	0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x54, 0x54, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45,
	0x5f, 0x50, 0x41, 0x53, 0x53, 0x57, 0x4f, 0x52, 0x44, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10, 0x45,
	0x54, 0x54, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x5f, 0x45, 0x4d, 0x41, 0x49, 0x4c, 0x10,
	0x04, 0x22, 0x5d, 0x0a, 0x0d, 0x55, 0x32, 0x55, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x12, 0x19, 0x0a, 0x15, 0x55, 0x32, 0x55, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a,
	0x14, 0x55, 0x32, 0x55, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x46,
	0x4f, 0x4c, 0x4c, 0x4f, 0x57, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x55, 0x32, 0x55, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x10, 0x02,
	0x42, 0x84, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x76, 0x31, 0x42, 0x0d, 0x45, 0x6e, 0x75, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x1f, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x74, 0x2f, 0x70,
	0x62, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x43, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c,
	0x56, 0x31, 0xe2, 0x02, 0x15, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x43, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_v1_enum_user_proto_rawDescOnce sync.Once
	file_common_v1_enum_user_proto_rawDescData = file_common_v1_enum_user_proto_rawDesc
)

func file_common_v1_enum_user_proto_rawDescGZIP() []byte {
	file_common_v1_enum_user_proto_rawDescOnce.Do(func() {
		file_common_v1_enum_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_v1_enum_user_proto_rawDescData)
	})
	return file_common_v1_enum_user_proto_rawDescData
}

var file_common_v1_enum_user_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_common_v1_enum_user_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_v1_enum_user_proto_goTypes = []interface{}{
	(USER_STATUS)(0),        // 0: common.v1.USER.STATUS
	(USER_AS)(0),            // 1: common.v1.USER.AS
	(USER_EBS)(0),           // 2: common.v1.USER.EBS
	(USER_ETT)(0),           // 3: common.v1.USER.ETT
	(USER_U2U_STAT_TYPE)(0), // 4: common.v1.USER.U2U_STAT_TYPE
	(*USER)(nil),            // 5: common.v1.USER
}
var file_common_v1_enum_user_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_v1_enum_user_proto_init() }
func file_common_v1_enum_user_proto_init() {
	if File_common_v1_enum_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_v1_enum_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*USER); i {
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
			RawDescriptor: file_common_v1_enum_user_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_v1_enum_user_proto_goTypes,
		DependencyIndexes: file_common_v1_enum_user_proto_depIdxs,
		EnumInfos:         file_common_v1_enum_user_proto_enumTypes,
		MessageInfos:      file_common_v1_enum_user_proto_msgTypes,
	}.Build()
	File_common_v1_enum_user_proto = out.File
	file_common_v1_enum_user_proto_rawDesc = nil
	file_common_v1_enum_user_proto_goTypes = nil
	file_common_v1_enum_user_proto_depIdxs = nil
}
