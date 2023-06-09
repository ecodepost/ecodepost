// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: common/v1/enum_community.proto

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

type CMT_MS int32

const (
	// 无效枚举值
	CMT_MS_INVALID CMT_MS = 0
	// 已邀请
	CMT_MS_INVITED CMT_MS = 1
	// 已接受
	CMT_MS_ACCEPTED CMT_MS = 2
	// 已拒绝
	CMT_MS_REJECTED CMT_MS = 3
	// 申请加入，审核中
	CMT_MS_APPLY_JOIN_AUDITING CMT_MS = 4
	// 申请加入，审核中
	CMT_MS_APPLY_JOIN_AUDIT_REJECT CMT_MS = 5
)

// Enum value maps for CMT_MS.
var (
	CMT_MS_name = map[int32]string{
		0: "MS_INVALID",
		1: "MS_INVITED",
		2: "MS_ACCEPTED",
		3: "MS_REJECTED",
		4: "MS_APPLY_JOIN_AUDITING",
		5: "MS_APPLY_JOIN_AUDIT_REJECT",
	}
	CMT_MS_value = map[string]int32{
		"MS_INVALID":                 0,
		"MS_INVITED":                 1,
		"MS_ACCEPTED":                2,
		"MS_REJECTED":                3,
		"MS_APPLY_JOIN_AUDITING":     4,
		"MS_APPLY_JOIN_AUDIT_REJECT": 5,
	}
)

func (x CMT_MS) Enum() *CMT_MS {
	p := new(CMT_MS)
	*p = x
	return p
}

func (x CMT_MS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CMT_MS) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_community_proto_enumTypes[0].Descriptor()
}

func (CMT_MS) Type() protoreflect.EnumType {
	return &file_common_v1_enum_community_proto_enumTypes[0]
}

func (x CMT_MS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CMT_MS.Descriptor instead.
func (CMT_MS) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_community_proto_rawDescGZIP(), []int{0, 0}
}

// ROLE:Role 社区角色
type CMT_ROLE int32

const (
	// 未知状态
	CMT_ROLE_INVALID CMT_ROLE = 0
	// Owner
	CMT_ROLE_OWNER CMT_ROLE = 1
	// 管理者
	CMT_ROLE_ADMIN CMT_ROLE = 2
	// 成员
	CMT_ROLE_MEMBER CMT_ROLE = 3
	// Guest
	CMT_ROLE_GUEST CMT_ROLE = 4
)

// Enum value maps for CMT_ROLE.
var (
	CMT_ROLE_name = map[int32]string{
		0: "ROLE_INVALID",
		1: "ROLE_OWNER",
		2: "ROLE_ADMIN",
		3: "ROLE_MEMBER",
		4: "ROLE_GUEST",
	}
	CMT_ROLE_value = map[string]int32{
		"ROLE_INVALID": 0,
		"ROLE_OWNER":   1,
		"ROLE_ADMIN":   2,
		"ROLE_MEMBER":  3,
		"ROLE_GUEST":   4,
	}
)

func (x CMT_ROLE) Enum() *CMT_ROLE {
	p := new(CMT_ROLE)
	*p = x
	return p
}

func (x CMT_ROLE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CMT_ROLE) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_community_proto_enumTypes[1].Descriptor()
}

func (CMT_ROLE) Type() protoreflect.EnumType {
	return &file_common_v1_enum_community_proto_enumTypes[1]
}

func (x CMT_ROLE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CMT_ROLE.Descriptor instead.
func (CMT_ROLE) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_community_proto_rawDescGZIP(), []int{0, 1}
}

// IS:InvitationStatus 邀请状态
type CMT_IS int32

const (
	// 未知状态
	CMT_IS_INVALID CMT_IS = 0
	// 邀请中
	CMT_IS_INVITING CMT_IS = 1
	// 已邀请成功
	CMT_IS_ACCEPTED CMT_IS = 2
	// 邀请被拒绝
	CMT_IS_REJECTED CMT_IS = 3
	// 邀请被取消
	CMT_IS_CANCELED CMT_IS = 4
	// 邀请过期
	CMT_IS_EXPIRED CMT_IS = 5
)

// Enum value maps for CMT_IS.
var (
	CMT_IS_name = map[int32]string{
		0: "IS_INVALID",
		1: "IS_INVITING",
		2: "IS_ACCEPTED",
		3: "IS_REJECTED",
		4: "IS_CANCELED",
		5: "IS_EXPIRED",
	}
	CMT_IS_value = map[string]int32{
		"IS_INVALID":  0,
		"IS_INVITING": 1,
		"IS_ACCEPTED": 2,
		"IS_REJECTED": 3,
		"IS_CANCELED": 4,
		"IS_EXPIRED":  5,
	}
)

func (x CMT_IS) Enum() *CMT_IS {
	p := new(CMT_IS)
	*p = x
	return p
}

func (x CMT_IS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CMT_IS) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_community_proto_enumTypes[2].Descriptor()
}

func (CMT_IS) Type() protoreflect.EnumType {
	return &file_common_v1_enum_community_proto_enumTypes[2]
}

func (x CMT_IS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CMT_IS.Descriptor instead.
func (CMT_IS) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_community_proto_rawDescGZIP(), []int{0, 2}
}

// AS:AuditStatus 团队审核状态
type CMT_AS int32

const (
	// 未知状态
	CMT_AS_INVALID CMT_AS = 0
	// 已申请
	CMT_AS_APPLIED CMT_AS = 1
	// 已通过
	CMT_AS_PASSED CMT_AS = 2
	// 已拒绝
	CMT_AS_REJECTED CMT_AS = 3
)

// Enum value maps for CMT_AS.
var (
	CMT_AS_name = map[int32]string{
		0: "AS_INVALID",
		1: "AS_APPLIED",
		2: "AS_PASSED",
		3: "AS_REJECTED",
	}
	CMT_AS_value = map[string]int32{
		"AS_INVALID":  0,
		"AS_APPLIED":  1,
		"AS_PASSED":   2,
		"AS_REJECTED": 3,
	}
)

func (x CMT_AS) Enum() *CMT_AS {
	p := new(CMT_AS)
	*p = x
	return p
}

func (x CMT_AS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CMT_AS) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_community_proto_enumTypes[3].Descriptor()
}

func (CMT_AS) Type() protoreflect.EnumType {
	return &file_common_v1_enum_community_proto_enumTypes[3]
}

func (x CMT_AS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CMT_AS.Descriptor instead.
func (CMT_AS) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_community_proto_rawDescGZIP(), []int{0, 3}
}

// DU:DurationUnit 时长单位
type CMT_DU int32

const (
	// 无效枚举值
	CMT_DU_INVALID CMT_DU = 0
	// 计算时长为：秒
	CMT_DU_SECOND CMT_DU = 1
	// 计算时长为：天
	CMT_DU_DAY CMT_DU = 2
	// 计算时长为：月
	CMT_DU_MONTH CMT_DU = 3
	// 计算时长为：年
	CMT_DU_YEAR CMT_DU = 4
)

// Enum value maps for CMT_DU.
var (
	CMT_DU_name = map[int32]string{
		0: "DU_INVALID",
		1: "DU_SECOND",
		2: "DU_DAY",
		3: "DU_MONTH",
		4: "DU_YEAR",
	}
	CMT_DU_value = map[string]int32{
		"DU_INVALID": 0,
		"DU_SECOND":  1,
		"DU_DAY":     2,
		"DU_MONTH":   3,
		"DU_YEAR":    4,
	}
)

func (x CMT_DU) Enum() *CMT_DU {
	p := new(CMT_DU)
	*p = x
	return p
}

func (x CMT_DU) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CMT_DU) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_community_proto_enumTypes[4].Descriptor()
}

func (CMT_DU) Type() protoreflect.EnumType {
	return &file_common_v1_enum_community_proto_enumTypes[4]
}

func (x CMT_DU) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CMT_DU.Descriptor instead.
func (CMT_DU) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_community_proto_rawDescGZIP(), []int{0, 4}
}

// Access:可访问配置
type CMT_ACS int32

const (
	// 未知类型
	CMT_ACS_INVALID CMT_ACS = 0
	// 可公开访问，（管理员仍可手动添加），任何能进入到社区的用户都能wei访问此社区
	CMT_ACS_OPEN CMT_ACS = 1
	// 外部用户禁止进入社区（管理员仍可手动添加）
	CMT_ACS_DENY_ALL CMT_ACS = 2
	// 私密访问，外部用户主动申请（管理员仍可手动添加），才能加入此社区
	CMT_ACS_USER_APPLY CMT_ACS = 3
	// 私密访问，外部用户购买资格（管理员仍可手动添加），才能加入此社区
	CMT_ACS_USER_PAY CMT_ACS = 4
)

// Enum value maps for CMT_ACS.
var (
	CMT_ACS_name = map[int32]string{
		0: "ACS_INVALID",
		1: "ACS_OPEN",
		2: "ACS_DENY_ALL",
		3: "ACS_USER_APPLY",
		4: "ACS_USER_PAY",
	}
	CMT_ACS_value = map[string]int32{
		"ACS_INVALID":    0,
		"ACS_OPEN":       1,
		"ACS_DENY_ALL":   2,
		"ACS_USER_APPLY": 3,
		"ACS_USER_PAY":   4,
	}
)

func (x CMT_ACS) Enum() *CMT_ACS {
	p := new(CMT_ACS)
	*p = x
	return p
}

func (x CMT_ACS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CMT_ACS) Descriptor() protoreflect.EnumDescriptor {
	return file_common_v1_enum_community_proto_enumTypes[5].Descriptor()
}

func (CMT_ACS) Type() protoreflect.EnumType {
	return &file_common_v1_enum_community_proto_enumTypes[5]
}

func (x CMT_ACS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CMT_ACS.Descriptor instead.
func (CMT_ACS) EnumDescriptor() ([]byte, []int) {
	return file_common_v1_enum_community_proto_rawDescGZIP(), []int{0, 5}
}

// CMT:Community 社区相关枚举值
type CMT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CMT) Reset() {
	*x = CMT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_v1_enum_community_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CMT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CMT) ProtoMessage() {}

func (x *CMT) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1_enum_community_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CMT.ProtoReflect.Descriptor instead.
func (*CMT) Descriptor() ([]byte, []int) {
	return file_common_v1_enum_community_proto_rawDescGZIP(), []int{0}
}

var File_common_v1_enum_community_proto protoreflect.FileDescriptor

var file_common_v1_enum_community_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x75, 0x6d,
	0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x22, 0xbf, 0x04, 0x0a, 0x03,
	0x43, 0x4d, 0x54, 0x22, 0x82, 0x01, 0x0a, 0x02, 0x4d, 0x53, 0x12, 0x0e, 0x0a, 0x0a, 0x4d, 0x53,
	0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x4d, 0x53,
	0x5f, 0x49, 0x4e, 0x56, 0x49, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x4d, 0x53,
	0x5f, 0x41, 0x43, 0x43, 0x45, 0x50, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x4d,
	0x53, 0x5f, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x03, 0x12, 0x1a, 0x0a, 0x16,
	0x4d, 0x53, 0x5f, 0x41, 0x50, 0x50, 0x4c, 0x59, 0x5f, 0x4a, 0x4f, 0x49, 0x4e, 0x5f, 0x41, 0x55,
	0x44, 0x49, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x1e, 0x0a, 0x1a, 0x4d, 0x53, 0x5f, 0x41,
	0x50, 0x50, 0x4c, 0x59, 0x5f, 0x4a, 0x4f, 0x49, 0x4e, 0x5f, 0x41, 0x55, 0x44, 0x49, 0x54, 0x5f,
	0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x10, 0x05, 0x22, 0x59, 0x0a, 0x04, 0x52, 0x4f, 0x4c, 0x45,
	0x12, 0x10, 0x0a, 0x0c, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44,
	0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x4f, 0x57, 0x4e, 0x45, 0x52,
	0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e,
	0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x4d, 0x45, 0x4d, 0x42, 0x45,
	0x52, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x47, 0x55, 0x45, 0x53,
	0x54, 0x10, 0x04, 0x22, 0x68, 0x0a, 0x02, 0x49, 0x53, 0x12, 0x0e, 0x0a, 0x0a, 0x49, 0x53, 0x5f,
	0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x53, 0x5f,
	0x49, 0x4e, 0x56, 0x49, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x53,
	0x5f, 0x41, 0x43, 0x43, 0x45, 0x50, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x49,
	0x53, 0x5f, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x03, 0x12, 0x0f, 0x0a, 0x0b,
	0x49, 0x53, 0x5f, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x12, 0x0e, 0x0a,
	0x0a, 0x49, 0x53, 0x5f, 0x45, 0x58, 0x50, 0x49, 0x52, 0x45, 0x44, 0x10, 0x05, 0x22, 0x44, 0x0a,
	0x02, 0x41, 0x53, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x53, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49,
	0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x53, 0x5f, 0x41, 0x50, 0x50, 0x4c, 0x49, 0x45,
	0x44, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x53, 0x5f, 0x50, 0x41, 0x53, 0x53, 0x45, 0x44,
	0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x41, 0x53, 0x5f, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x45,
	0x44, 0x10, 0x03, 0x22, 0x4a, 0x0a, 0x02, 0x44, 0x55, 0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x55, 0x5f,
	0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x55, 0x5f,
	0x53, 0x45, 0x43, 0x4f, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x55, 0x5f, 0x44,
	0x41, 0x59, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x55, 0x5f, 0x4d, 0x4f, 0x4e, 0x54, 0x48,
	0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x55, 0x5f, 0x59, 0x45, 0x41, 0x52, 0x10, 0x04, 0x22,
	0x5c, 0x0a, 0x03, 0x41, 0x43, 0x53, 0x12, 0x0f, 0x0a, 0x0b, 0x41, 0x43, 0x53, 0x5f, 0x49, 0x4e,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x43, 0x53, 0x5f, 0x4f,
	0x50, 0x45, 0x4e, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x41, 0x43, 0x53, 0x5f, 0x44, 0x45, 0x4e,
	0x59, 0x5f, 0x41, 0x4c, 0x4c, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x41, 0x43, 0x53, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x5f, 0x41, 0x50, 0x50, 0x4c, 0x59, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x41,
	0x43, 0x53, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x50, 0x41, 0x59, 0x10, 0x04, 0x42, 0x89, 0x01,
	0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42,
	0x12, 0x45, 0x6e, 0x75, 0x6d, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x1f, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x74,
	0x2f, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x43,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x43,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_common_v1_enum_community_proto_rawDescOnce sync.Once
	file_common_v1_enum_community_proto_rawDescData = file_common_v1_enum_community_proto_rawDesc
)

func file_common_v1_enum_community_proto_rawDescGZIP() []byte {
	file_common_v1_enum_community_proto_rawDescOnce.Do(func() {
		file_common_v1_enum_community_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_v1_enum_community_proto_rawDescData)
	})
	return file_common_v1_enum_community_proto_rawDescData
}

var file_common_v1_enum_community_proto_enumTypes = make([]protoimpl.EnumInfo, 6)
var file_common_v1_enum_community_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_v1_enum_community_proto_goTypes = []interface{}{
	(CMT_MS)(0),   // 0: common.v1.CMT.MS
	(CMT_ROLE)(0), // 1: common.v1.CMT.ROLE
	(CMT_IS)(0),   // 2: common.v1.CMT.IS
	(CMT_AS)(0),   // 3: common.v1.CMT.AS
	(CMT_DU)(0),   // 4: common.v1.CMT.DU
	(CMT_ACS)(0),  // 5: common.v1.CMT.ACS
	(*CMT)(nil),   // 6: common.v1.CMT
}
var file_common_v1_enum_community_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_v1_enum_community_proto_init() }
func file_common_v1_enum_community_proto_init() {
	if File_common_v1_enum_community_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_v1_enum_community_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CMT); i {
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
			RawDescriptor: file_common_v1_enum_community_proto_rawDesc,
			NumEnums:      6,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_v1_enum_community_proto_goTypes,
		DependencyIndexes: file_common_v1_enum_community_proto_depIdxs,
		EnumInfos:         file_common_v1_enum_community_proto_enumTypes,
		MessageInfos:      file_common_v1_enum_community_proto_msgTypes,
	}.Build()
	File_common_v1_enum_community_proto = out.File
	file_common_v1_enum_community_proto_rawDesc = nil
	file_common_v1_enum_community_proto_goTypes = nil
	file_common_v1_enum_community_proto_depIdxs = nil
}
