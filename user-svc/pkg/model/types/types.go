package types

import (
	commonv1 "ecodepost/pb/common/v1"
)

// PlatformMap 登录平台映射关系
var PlatformMap = map[commonv1.SSO_Platform]string{
	commonv1.SSO_PLATFORM_WEB:     "web",
	commonv1.SSO_PLATFORM_PC:      "pc",
	commonv1.SSO_PLATFORM_ANDROID: "android",
	commonv1.SSO_PLATFORM_IOS:     "ios",
}
