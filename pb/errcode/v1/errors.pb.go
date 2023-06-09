// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: errcode/v1/errors.proto

package errcodev1

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

// @plugins=protoc-gen-go-errors
// 错误定义集合
type Error int32

const (
	// @code=UNKNOWN
	// @i18n.cn="未知类型"
	// @i18n.en="UNKNOWN error"
	Error_ERR_UNKNOWN Error = 0
	// @code=INTERNAL
	// @i18n.cn="内部错误"
	// @i18n.en="INTERNAL error"
	Error_ERR_INTERNAL Error = 1
	// @code=INVALID_ARGUMENT
	// @i18n.cn="参数错误"
	// @i18n.en="INVALID_ARGUMENT error"
	Error_ERR_INVALID_ARGUMENT Error = 2
	// @code=NOT_FOUND
	// @i18n.cn="找不到资源"
	// @i18n.en="NOT_FOUND error"
	Error_ERR_NOT_FOUND Error = 3
	// @code=ALREADY_EXISTS
	// @i18n.cn="已经存在"
	// @i18n.en="ALREADY_EXISTS error"
	Error_ERR_ALREADY_EXIST Error = 4
	// @code=PERMISSION_DENIED
	// @i18n.cn="没有权限"
	// @i18n.en="PERMISSION_DENIED error"
	Error_ERR_ALREADY_PERMISSION_DENIED Error = 5
	// @code=ABORTED
	// @i18n.cn="放弃操作"
	// @i18n.en="ABORTED error"
	Error_ERR_ABORTED Error = 6
	// @code=OUT_OF_RANGE
	// @i18n.cn="超过限制"
	// @i18n.en="OUT_OF_RANGE error"
	Error_ERR_OUT_OF_RANGE Error = 7
	// @code=DATA_LOSS
	// @i18n.cn="超过限制"
	// @i18n.en="DATA_LOSS error"
	Error_ERR_DATA_LOSS Error = 8
	// @code=INTERNAL
	// @i18n.cn="DB错误"
	// @i18n.en="Intern DB error"
	Error_ERR_DB_ERROR Error = 9
	// @code=INTERNAL
	// @i18n.cn="Redis错误"
	// @i18n.en="Intern Redis error"
	Error_ERR_REDIS_ERROR Error = 10
	// @code=INTERNAL
	// @i18n.cn="编码 JSON 错误"
	// @i18n.en="Marshal JSON error"
	Error_ERR_JSON_MARSHAL_ERROR Error = 11
	// @code=INTERNAL
	// @i18n.cn="编码 PROTO 错误"
	// @i18n.en="Marshal PROTO error"
	Error_ERR_PROTO_MARSHAL_ERROR Error = 12
	// @code=INVALID_ARGUMENT
	// @i18n.cn="UID不能为空"
	// @i18n.en="UID can't be empty"
	Error_ERR_UID_EMPTY Error = 101
	// @code=INVALID_ARGUMENT
	// @i18n.cn="社区不能为空"
	// @i18n.en="Community can't be empty"
	//  ERR_COMMUNITY_GUID_EMPTY = 102;
	// @code=INVALID_ARGUMENT
	// @i18n.cn="名字不能为空"
	// @i18n.en="Name can't be empty"
	Error_ERR_NAME_EMPTY Error = 103
	// @code=INVALID_ARGUMENT
	// @i18n.cn="LOGO不能为空"
	// @i18n.en="Logo can't be empty"
	Error_ERR_LOGO_EMPTY Error = 104
	// @code=INVALID_ARGUMENT
	// @i18n.cn="自定义域名不能为空"
	// @i18n.en="Custom domain can't be empty"
	Error_ERR_CUSTOM_DOMAIN_EMPTY Error = 105
	// @code=ALREADY_EXISTS
	// @i18n.cn="自定义域名已经存在"
	// @i18n.en="Custom domain already exist"
	Error_ERR_CUSTOM_DOMAIN_ALREADY_EXIST Error = 106
	// @code=INVALID_ARGUMENT
	// @i18n.cn="文件GUID不能为空"
	// @i18n.en="File guid can't be empty"
	Error_ERR_FILE_GUID_EMPTY Error = 107
	// @code=INVALID_ARGUMENT
	// @i18n.cn="空间分组不能为空"
	// @i18n.en="space group can't be empty"
	Error_ERR_SPACE_GROUP_EMPTY Error = 108
	// @code=INVALID_ARGUMENT
	// @i18n.cn="空间不能为空"
	// @i18n.en="space can't be empty"
	Error_ERR_SPACE_EMPTY Error = 109
	// @code=PERMISSION_DENIED
	// @i18n.cn="社区版本到最大限制"
	// @i18n.en="community edition reach limited"
	Error_ERR_EDITION_REACH_LIMITED Error = 110
	// @code=PERMISSION_DENIED
	// @i18n.cn="你没有权限"
	// @i18n.en="you don't have permission"
	Error_ERR_PMS_PERMISSION_DENIED Error = 111
	// @code=INVALID_ARGUMENT
	// @i18n.cn="操作用户不能为空"
	// @i18n.en="operate uid can't be empty"
	Error_ERR_OPERATE_UID_EMPTY Error = 112
	// @code=ALREADY_EXISTS
	// @i18n.cn="不能重复申请"
	// @i18n.en="repeated apply forbidden"
	Error_ERR_REPEATED_APPLY Error = 113
	// @code=INVALID_ARGUMENT
	// @i18n.cn="guid不匹配"
	// @i18n.en="guid not match"
	Error_ERR_COMMUNITY_GUID_NOT_MATCH Error = 114
	// @code=INVALID_ARGUMENT
	// @i18n.cn="guid不存在"
	// @i18n.en="guid not exist"
	Error_ERR_COMMUNITY_GUID_NOT_EXIST Error = 115
	// @code=INVALID_ARGUMENT
	// @i18n.cn="实体Guid不能为空"
	// @i18n.en="Guid can't be empty"
	Error_ERR_GUID_EMPTY Error = 116
	// @code=INVALID_ARGUMENT
	// @i18n.cn="社区可见度非法"
	// @i18n.en="cmt not match"
	Error_ERR_COMMUNITY_VISIBILITY_INVALID Error = 118
	// 文件名称长度不合法
	// @code=INVALID_ARGUMENT
	// @i18n.cn="file name should not less than 0 or bigger than 100"
	// @i18n.en="get parent token fail"
	Error_ERR_FILE_NAME_LENGTH Error = 119
	// @code=INVALID_ARGUMENT
	// @i18n.cn="业务Guid不能为空"
	// @i18n.en="Business guid can't be empty"
	Error_ERR_BIZ_GUID_EMPTY Error = 120
	// @code=INVALID_ARGUMENT
	// @i18n.cn="文本内容不能为空"
	// @i18n.en="Content can't be empty"
	Error_ERR_FILE_CONTENT_EMPTY Error = 121
	// @code=UNKNOWN
	// @i18n.cn="授权未知错误"
	// @i18n.en="auth invalid"
	Error_AUTH_ERR_INVALID Error = 200
	// 获取sub token的浏览器系统错误
	// @code=INTERNAL
	// @i18n.cn="浏览器token cookie系统错误"
	// @i18n.en="token cookie invalid"
	Error_AUTH_ERR_BROWSER_COOKIE_SYSTEM_ERROR Error = 201
	// 获取parent token的浏览器系统错误
	// @code=INTERNAL
	// @i18n.cn="浏览器parent token cookie系统错误"
	// @i18n.en="parent token cookie invalid"
	Error_AUTH_ERR_BROWSER_PARENT_COOKIE_SYSTEM_ERROR Error = 202
	// parent token不存在
	// @code=INVALID_ARGUMENT
	// @i18n.cn="浏览器parent token不存在"
	// @i18n.en="parent token cookie empty"
	Error_AUTH_ERR_BROWSER_PARENT_COOKIE_EMPTY Error = 203
	// grpc获取parent token失败
	// @code=INTERNAL
	// @i18n.cn="获取access失败"
	// @i18n.en="get parent token fail"
	Error_AUTH_ERR_GET_ACCESS_BY_PARENT_COOKIE_ERROR Error = 204
	// refresh token失败
	// @code=INTERNAL
	// @i18n.cn="刷新token失败"
	// @i18n.en="get parent token fail"
	Error_AUTH_ERR_REFRESH_TOKEN_ERROR Error = 205
	// 获取用户信息失败
	// @code=INTERNAL
	// @i18n.cn="获取用户信息失败"
	// @i18n.en="get parent token fail"
	Error_AUTH_ERR_GET_USER_INFO_BY_TOKEN_ERROR Error = 206
	// @code=UNKNOWN
	// @i18n.cn="个人中心未知错误"
	// @i18n.en="profile invalid"
	Error_ERR_PROFILE_INVALID Error = 300
	// @code=INVALID_ARGUMENT
	// @i18n.cn="手机号已存在"
	// @i18n.en="phone exist"
	Error_ERR_PROFILE_PHONE_EXIST Error = 301
	// @code=INVALID_ARGUMENT
	// @i18n.cn="邮箱已存在"
	// @i18n.en="email exist"
	Error_ERR_PROFILE_EMAIL_EXIST Error = 302
	// 用户已加入社区
	// @code=ALREADY_EXISTS
	// @i18n.cn="用户已加入社区"
	// @i18n.en="get parent token fail"
	Error_ERR_USER_HAS_JOINED_CMT Error = 400
)

// Enum value maps for Error.
var (
	Error_name = map[int32]string{
		0:   "ERR_UNKNOWN",
		1:   "ERR_INTERNAL",
		2:   "ERR_INVALID_ARGUMENT",
		3:   "ERR_NOT_FOUND",
		4:   "ERR_ALREADY_EXIST",
		5:   "ERR_ALREADY_PERMISSION_DENIED",
		6:   "ERR_ABORTED",
		7:   "ERR_OUT_OF_RANGE",
		8:   "ERR_DATA_LOSS",
		9:   "ERR_DB_ERROR",
		10:  "ERR_REDIS_ERROR",
		11:  "ERR_JSON_MARSHAL_ERROR",
		12:  "ERR_PROTO_MARSHAL_ERROR",
		101: "ERR_UID_EMPTY",
		103: "ERR_NAME_EMPTY",
		104: "ERR_LOGO_EMPTY",
		105: "ERR_CUSTOM_DOMAIN_EMPTY",
		106: "ERR_CUSTOM_DOMAIN_ALREADY_EXIST",
		107: "ERR_FILE_GUID_EMPTY",
		108: "ERR_SPACE_GROUP_EMPTY",
		109: "ERR_SPACE_EMPTY",
		110: "ERR_EDITION_REACH_LIMITED",
		111: "ERR_PMS_PERMISSION_DENIED",
		112: "ERR_OPERATE_UID_EMPTY",
		113: "ERR_REPEATED_APPLY",
		114: "ERR_COMMUNITY_GUID_NOT_MATCH",
		115: "ERR_COMMUNITY_GUID_NOT_EXIST",
		116: "ERR_GUID_EMPTY",
		118: "ERR_COMMUNITY_VISIBILITY_INVALID",
		119: "ERR_FILE_NAME_LENGTH",
		120: "ERR_BIZ_GUID_EMPTY",
		121: "ERR_FILE_CONTENT_EMPTY",
		200: "AUTH_ERR_INVALID",
		201: "AUTH_ERR_BROWSER_COOKIE_SYSTEM_ERROR",
		202: "AUTH_ERR_BROWSER_PARENT_COOKIE_SYSTEM_ERROR",
		203: "AUTH_ERR_BROWSER_PARENT_COOKIE_EMPTY",
		204: "AUTH_ERR_GET_ACCESS_BY_PARENT_COOKIE_ERROR",
		205: "AUTH_ERR_REFRESH_TOKEN_ERROR",
		206: "AUTH_ERR_GET_USER_INFO_BY_TOKEN_ERROR",
		300: "ERR_PROFILE_INVALID",
		301: "ERR_PROFILE_PHONE_EXIST",
		302: "ERR_PROFILE_EMAIL_EXIST",
		400: "ERR_USER_HAS_JOINED_CMT",
	}
	Error_value = map[string]int32{
		"ERR_UNKNOWN":                                 0,
		"ERR_INTERNAL":                                1,
		"ERR_INVALID_ARGUMENT":                        2,
		"ERR_NOT_FOUND":                               3,
		"ERR_ALREADY_EXIST":                           4,
		"ERR_ALREADY_PERMISSION_DENIED":               5,
		"ERR_ABORTED":                                 6,
		"ERR_OUT_OF_RANGE":                            7,
		"ERR_DATA_LOSS":                               8,
		"ERR_DB_ERROR":                                9,
		"ERR_REDIS_ERROR":                             10,
		"ERR_JSON_MARSHAL_ERROR":                      11,
		"ERR_PROTO_MARSHAL_ERROR":                     12,
		"ERR_UID_EMPTY":                               101,
		"ERR_NAME_EMPTY":                              103,
		"ERR_LOGO_EMPTY":                              104,
		"ERR_CUSTOM_DOMAIN_EMPTY":                     105,
		"ERR_CUSTOM_DOMAIN_ALREADY_EXIST":             106,
		"ERR_FILE_GUID_EMPTY":                         107,
		"ERR_SPACE_GROUP_EMPTY":                       108,
		"ERR_SPACE_EMPTY":                             109,
		"ERR_EDITION_REACH_LIMITED":                   110,
		"ERR_PMS_PERMISSION_DENIED":                   111,
		"ERR_OPERATE_UID_EMPTY":                       112,
		"ERR_REPEATED_APPLY":                          113,
		"ERR_COMMUNITY_GUID_NOT_MATCH":                114,
		"ERR_COMMUNITY_GUID_NOT_EXIST":                115,
		"ERR_GUID_EMPTY":                              116,
		"ERR_COMMUNITY_VISIBILITY_INVALID":            118,
		"ERR_FILE_NAME_LENGTH":                        119,
		"ERR_BIZ_GUID_EMPTY":                          120,
		"ERR_FILE_CONTENT_EMPTY":                      121,
		"AUTH_ERR_INVALID":                            200,
		"AUTH_ERR_BROWSER_COOKIE_SYSTEM_ERROR":        201,
		"AUTH_ERR_BROWSER_PARENT_COOKIE_SYSTEM_ERROR": 202,
		"AUTH_ERR_BROWSER_PARENT_COOKIE_EMPTY":        203,
		"AUTH_ERR_GET_ACCESS_BY_PARENT_COOKIE_ERROR":  204,
		"AUTH_ERR_REFRESH_TOKEN_ERROR":                205,
		"AUTH_ERR_GET_USER_INFO_BY_TOKEN_ERROR":       206,
		"ERR_PROFILE_INVALID":                         300,
		"ERR_PROFILE_PHONE_EXIST":                     301,
		"ERR_PROFILE_EMAIL_EXIST":                     302,
		"ERR_USER_HAS_JOINED_CMT":                     400,
	}
)

func (x Error) Enum() *Error {
	p := new(Error)
	*p = x
	return p
}

func (x Error) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Error) Descriptor() protoreflect.EnumDescriptor {
	return file_errcode_v1_errors_proto_enumTypes[0].Descriptor()
}

func (Error) Type() protoreflect.EnumType {
	return &file_errcode_v1_errors_proto_enumTypes[0]
}

func (x Error) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Error.Descriptor instead.
func (Error) EnumDescriptor() ([]byte, []int) {
	return file_errcode_v1_errors_proto_rawDescGZIP(), []int{0}
}

var File_errcode_v1_errors_proto protoreflect.FileDescriptor

var file_errcode_v1_errors_proto_rawDesc = []byte{
	0x0a, 0x17, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x65, 0x72, 0x72, 0x63, 0x6f,
	0x64, 0x65, 0x2e, 0x76, 0x31, 0x2a, 0xc1, 0x09, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12,
	0x0f, 0x0a, 0x0b, 0x45, 0x52, 0x52, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x10, 0x0a, 0x0c, 0x45, 0x52, 0x52, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c,
	0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x45, 0x52, 0x52, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49,
	0x44, 0x5f, 0x41, 0x52, 0x47, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d,
	0x45, 0x52, 0x52, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x03, 0x12,
	0x15, 0x0a, 0x11, 0x45, 0x52, 0x52, 0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x45,
	0x58, 0x49, 0x53, 0x54, 0x10, 0x04, 0x12, 0x21, 0x0a, 0x1d, 0x45, 0x52, 0x52, 0x5f, 0x41, 0x4c,
	0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e,
	0x5f, 0x44, 0x45, 0x4e, 0x49, 0x45, 0x44, 0x10, 0x05, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x52, 0x52,
	0x5f, 0x41, 0x42, 0x4f, 0x52, 0x54, 0x45, 0x44, 0x10, 0x06, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x52,
	0x52, 0x5f, 0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x52, 0x41, 0x4e, 0x47, 0x45, 0x10, 0x07,
	0x12, 0x11, 0x0a, 0x0d, 0x45, 0x52, 0x52, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x4c, 0x4f, 0x53,
	0x53, 0x10, 0x08, 0x12, 0x10, 0x0a, 0x0c, 0x45, 0x52, 0x52, 0x5f, 0x44, 0x42, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x10, 0x09, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x52, 0x52, 0x5f, 0x52, 0x45, 0x44,
	0x49, 0x53, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x0a, 0x12, 0x1a, 0x0a, 0x16, 0x45, 0x52,
	0x52, 0x5f, 0x4a, 0x53, 0x4f, 0x4e, 0x5f, 0x4d, 0x41, 0x52, 0x53, 0x48, 0x41, 0x4c, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x0b, 0x12, 0x1b, 0x0a, 0x17, 0x45, 0x52, 0x52, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x5f, 0x4d, 0x41, 0x52, 0x53, 0x48, 0x41, 0x4c, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0x0c, 0x12, 0x11, 0x0a, 0x0d, 0x45, 0x52, 0x52, 0x5f, 0x55, 0x49, 0x44, 0x5f, 0x45,
	0x4d, 0x50, 0x54, 0x59, 0x10, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x45, 0x52, 0x52, 0x5f, 0x4e, 0x41,
	0x4d, 0x45, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x67, 0x12, 0x12, 0x0a, 0x0e, 0x45, 0x52,
	0x52, 0x5f, 0x4c, 0x4f, 0x47, 0x4f, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x68, 0x12, 0x1b,
	0x0a, 0x17, 0x45, 0x52, 0x52, 0x5f, 0x43, 0x55, 0x53, 0x54, 0x4f, 0x4d, 0x5f, 0x44, 0x4f, 0x4d,
	0x41, 0x49, 0x4e, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x69, 0x12, 0x23, 0x0a, 0x1f, 0x45,
	0x52, 0x52, 0x5f, 0x43, 0x55, 0x53, 0x54, 0x4f, 0x4d, 0x5f, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e,
	0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x10, 0x6a,
	0x12, 0x17, 0x0a, 0x13, 0x45, 0x52, 0x52, 0x5f, 0x46, 0x49, 0x4c, 0x45, 0x5f, 0x47, 0x55, 0x49,
	0x44, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x6b, 0x12, 0x19, 0x0a, 0x15, 0x45, 0x52, 0x52,
	0x5f, 0x53, 0x50, 0x41, 0x43, 0x45, 0x5f, 0x47, 0x52, 0x4f, 0x55, 0x50, 0x5f, 0x45, 0x4d, 0x50,
	0x54, 0x59, 0x10, 0x6c, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x52, 0x52, 0x5f, 0x53, 0x50, 0x41, 0x43,
	0x45, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x6d, 0x12, 0x1d, 0x0a, 0x19, 0x45, 0x52, 0x52,
	0x5f, 0x45, 0x44, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x45, 0x41, 0x43, 0x48, 0x5f, 0x4c,
	0x49, 0x4d, 0x49, 0x54, 0x45, 0x44, 0x10, 0x6e, 0x12, 0x1d, 0x0a, 0x19, 0x45, 0x52, 0x52, 0x5f,
	0x50, 0x4d, 0x53, 0x5f, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x44,
	0x45, 0x4e, 0x49, 0x45, 0x44, 0x10, 0x6f, 0x12, 0x19, 0x0a, 0x15, 0x45, 0x52, 0x52, 0x5f, 0x4f,
	0x50, 0x45, 0x52, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x49, 0x44, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59,
	0x10, 0x70, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x52, 0x52, 0x5f, 0x52, 0x45, 0x50, 0x45, 0x41, 0x54,
	0x45, 0x44, 0x5f, 0x41, 0x50, 0x50, 0x4c, 0x59, 0x10, 0x71, 0x12, 0x20, 0x0a, 0x1c, 0x45, 0x52,
	0x52, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x55, 0x4e, 0x49, 0x54, 0x59, 0x5f, 0x47, 0x55, 0x49, 0x44,
	0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x10, 0x72, 0x12, 0x20, 0x0a, 0x1c,
	0x45, 0x52, 0x52, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x55, 0x4e, 0x49, 0x54, 0x59, 0x5f, 0x47, 0x55,
	0x49, 0x44, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x10, 0x73, 0x12, 0x12,
	0x0a, 0x0e, 0x45, 0x52, 0x52, 0x5f, 0x47, 0x55, 0x49, 0x44, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59,
	0x10, 0x74, 0x12, 0x24, 0x0a, 0x20, 0x45, 0x52, 0x52, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x55, 0x4e,
	0x49, 0x54, 0x59, 0x5f, 0x56, 0x49, 0x53, 0x49, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x49,
	0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x76, 0x12, 0x18, 0x0a, 0x14, 0x45, 0x52, 0x52, 0x5f,
	0x46, 0x49, 0x4c, 0x45, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x5f, 0x4c, 0x45, 0x4e, 0x47, 0x54, 0x48,
	0x10, 0x77, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x52, 0x52, 0x5f, 0x42, 0x49, 0x5a, 0x5f, 0x47, 0x55,
	0x49, 0x44, 0x5f, 0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x78, 0x12, 0x1a, 0x0a, 0x16, 0x45, 0x52,
	0x52, 0x5f, 0x46, 0x49, 0x4c, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x45,
	0x4d, 0x50, 0x54, 0x59, 0x10, 0x79, 0x12, 0x15, 0x0a, 0x10, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45,
	0x52, 0x52, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0xc8, 0x01, 0x12, 0x29, 0x0a,
	0x24, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x5f, 0x42, 0x52, 0x4f, 0x57, 0x53, 0x45,
	0x52, 0x5f, 0x43, 0x4f, 0x4f, 0x4b, 0x49, 0x45, 0x5f, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xc9, 0x01, 0x12, 0x30, 0x0a, 0x2b, 0x41, 0x55, 0x54, 0x48,
	0x5f, 0x45, 0x52, 0x52, 0x5f, 0x42, 0x52, 0x4f, 0x57, 0x53, 0x45, 0x52, 0x5f, 0x50, 0x41, 0x52,
	0x45, 0x4e, 0x54, 0x5f, 0x43, 0x4f, 0x4f, 0x4b, 0x49, 0x45, 0x5f, 0x53, 0x59, 0x53, 0x54, 0x45,
	0x4d, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xca, 0x01, 0x12, 0x29, 0x0a, 0x24, 0x41, 0x55,
	0x54, 0x48, 0x5f, 0x45, 0x52, 0x52, 0x5f, 0x42, 0x52, 0x4f, 0x57, 0x53, 0x45, 0x52, 0x5f, 0x50,
	0x41, 0x52, 0x45, 0x4e, 0x54, 0x5f, 0x43, 0x4f, 0x4f, 0x4b, 0x49, 0x45, 0x5f, 0x45, 0x4d, 0x50,
	0x54, 0x59, 0x10, 0xcb, 0x01, 0x12, 0x2f, 0x0a, 0x2a, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45, 0x52,
	0x52, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x42, 0x59, 0x5f,
	0x50, 0x41, 0x52, 0x45, 0x4e, 0x54, 0x5f, 0x43, 0x4f, 0x4f, 0x4b, 0x49, 0x45, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x10, 0xcc, 0x01, 0x12, 0x21, 0x0a, 0x1c, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x45,
	0x52, 0x52, 0x5f, 0x52, 0x45, 0x46, 0x52, 0x45, 0x53, 0x48, 0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e,
	0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xcd, 0x01, 0x12, 0x2a, 0x0a, 0x25, 0x41, 0x55, 0x54,
	0x48, 0x5f, 0x45, 0x52, 0x52, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49,
	0x4e, 0x46, 0x4f, 0x5f, 0x42, 0x59, 0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0xce, 0x01, 0x12, 0x18, 0x0a, 0x13, 0x45, 0x52, 0x52, 0x5f, 0x50, 0x52, 0x4f,
	0x46, 0x49, 0x4c, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0xac, 0x02, 0x12,
	0x1c, 0x0a, 0x17, 0x45, 0x52, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x46, 0x49, 0x4c, 0x45, 0x5f, 0x50,
	0x48, 0x4f, 0x4e, 0x45, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x10, 0xad, 0x02, 0x12, 0x1c, 0x0a,
	0x17, 0x45, 0x52, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x46, 0x49, 0x4c, 0x45, 0x5f, 0x45, 0x4d, 0x41,
	0x49, 0x4c, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x10, 0xae, 0x02, 0x12, 0x1c, 0x0a, 0x17, 0x45,
	0x52, 0x52, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x48, 0x41, 0x53, 0x5f, 0x4a, 0x4f, 0x49, 0x4e,
	0x45, 0x44, 0x5f, 0x43, 0x4d, 0x54, 0x10, 0x90, 0x03, 0x42, 0x89, 0x01, 0x0a, 0x0e, 0x63, 0x6f,
	0x6d, 0x2e, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x21, 0x65, 0x63, 0x6f,
	0x64, 0x65, 0x70, 0x6f, 0x73, 0x74, 0x2f, 0x70, 0x62, 0x2f, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64,
	0x65, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x45, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x0a, 0x45, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x16, 0x45, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x45, 0x72, 0x72, 0x63, 0x6f, 0x64,
	0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errcode_v1_errors_proto_rawDescOnce sync.Once
	file_errcode_v1_errors_proto_rawDescData = file_errcode_v1_errors_proto_rawDesc
)

func file_errcode_v1_errors_proto_rawDescGZIP() []byte {
	file_errcode_v1_errors_proto_rawDescOnce.Do(func() {
		file_errcode_v1_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_errcode_v1_errors_proto_rawDescData)
	})
	return file_errcode_v1_errors_proto_rawDescData
}

var file_errcode_v1_errors_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_errcode_v1_errors_proto_goTypes = []interface{}{
	(Error)(0), // 0: errcode.v1.Error
}
var file_errcode_v1_errors_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_errcode_v1_errors_proto_init() }
func file_errcode_v1_errors_proto_init() {
	if File_errcode_v1_errors_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_errcode_v1_errors_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errcode_v1_errors_proto_goTypes,
		DependencyIndexes: file_errcode_v1_errors_proto_depIdxs,
		EnumInfos:         file_errcode_v1_errors_proto_enumTypes,
	}.Build()
	File_errcode_v1_errors_proto = out.File
	file_errcode_v1_errors_proto_rawDesc = nil
	file_errcode_v1_errors_proto_goTypes = nil
	file_errcode_v1_errors_proto_depIdxs = nil
}
