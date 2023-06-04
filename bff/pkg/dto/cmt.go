package dto

import (
	"encoding/json"

	commonv1 "ecodepost/pb/common/v1"
)

// CmtDetail 社区的所有信息
type CmtDetail struct {
	CmtInfo    CmtInfo       `json:"cmtInfo"` // 社区信息
	CmtTheme   CmtTheme      `json:"cmtTheme"`
	UserInfo   CmtUserInfo   `json:"userInfo"`   // 社区的用户信息
	Permission CmtPermission `json:"permission"` // 用户所在社区的权限信息
}

type CmtInfo struct {
	Name           string             `json:"name"`                  // 团队名称
	Description    string             `json:"description,omitempty"` // 团队描述
	Logo           string             `json:"logo,omitempty"`        // LOGO
	Visibility     commonv1.CMN_VISBL `json:"visibility,omitempty"`
	Ctime          int64              `json:"ctime,omitempty"`
	IsSetHome      bool               `json:"isSetHome,omitempty"`
	FirstVisitUrl  string             `json:"firstVisitUrl"` // 第一次访问的URL
	Access         commonv1.CMT_ACS   `json:"access"`
	GongxinbuBeian string             `json:"gongxinbuBeian"`
}

type CmtTheme struct {
	IsCustom          bool            `json:"isCustom,omitempty"`
	ThemeName         string          `json:"themeName,omitempty"`
	CustomColor       json.RawMessage `json:"customColor,omitempty"`
	DefaultAppearance string          `json:"defaultAppearance,omitempty"`
}

type CmtUserInfo struct {
	// 是否登录
	IsLogin bool `json:"isLogin"`
	// 是否在这个社区
	IsExist bool `json:"isExist"`
	// 是否为会员，如果为false，那么可能是他没购买，也可能是之前购买了，但过期了
	IsMemberShip bool `json:"isMemberShip"`
	// 是否为第一次购买，第一次请求会展示，如果购买后，会更新这个字段变为false
	IsFirstFollow bool `json:"isFirstFollow"`
	// 过期信息
	ExpireMsg string `json:"expireMsg"`
}

type AppInfo struct {
	AppId commonv1.CMN_APP `json:"appId"`
	Name  string           `json:"name"`
}

type CmtPermission struct {
	IsAllowManageCommunity  bool      `json:"isAllowManageCommunity"`
	IsAllowCreateSpaceGroup bool      `json:"isAllowCreateSpaceGroup"`
	IsAllowCreateSpace      bool      `json:"isAllowCreateSpace"`
	IsAllowUpgradeEdition   bool      `json:"isAllowUpgradeEdition"`
	AppList                 []AppInfo `json:"appList"`
}
