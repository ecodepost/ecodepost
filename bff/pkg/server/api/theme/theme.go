package theme

import (
	"encoding/json"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"

	communityv1 "ecodepost/pb/community/v1"
)

type GetRes struct {
	IsCustom          bool        `json:"isCustom"`          // 是否自定义
	ThemeName         string      `json:"themeName"`         // 主题名称
	CustomColor       CustomColor `json:"customColor"`       // 自定义颜色
	DefaultAppearance string      `json:"defaultAppearance"` // 白天模式，暗黑模式
}

type CustomColor struct {
	ThemeColorPrimary    string `json:"themeColorPrimary"`
	ThemeColorStatus     string `json:"themeColorStatus"`
	ThemeColorButtonText string `json:"themeColorButtonText"`
	ThemeColorBackground string `json:"themeColorBackground"`
}

// Get 获取自定义主题颜色信息
func Get(c *bffcore.Context) {
	themeRes, err := invoker.GrpcCommunity.GetTheme(c.Ctx(), &communityv1.GetThemeReq{})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	customColor := CustomColor{}
	if themeRes.CustomColor != "" {
		err = json.Unmarshal([]byte(themeRes.CustomColor), &customColor)
		if err != nil {
			c.JSONE(1, "解析失败", err.Error())
			return
		}
	}
	c.JSONOK(GetRes{
		IsCustom:          themeRes.IsCustom,
		ThemeName:         themeRes.ThemeName,
		CustomColor:       customColor,
		DefaultAppearance: themeRes.DefaultAppearance,
	})
}

type PutRes struct {
	IsCustom          bool        `json:"isCustom"`          // 是否自定义
	ThemeName         string      `json:"themeName"`         // 主题名称
	CustomColor       CustomColor `json:"customColor"`       // 自定义颜色
	DefaultAppearance string      `json:"defaultAppearance"` // 白天模式，暗黑模式
}

// Put 更新自定义主题颜色信息
func Put(c *bffcore.Context) {
	req := PutRes{}
	err := c.Bind(&req)
	if err != nil {
		c.JSONE(1, err.Error(), err)
		return
	}

	customColor, err := json.Marshal(req.CustomColor)
	if err != nil {
		c.JSONE(1, "json编码失败", nil)
		return
	}
	_, err = invoker.GrpcCommunity.SetTheme(c.Ctx(), &communityv1.SetThemeReq{
		IsCustom:          req.IsCustom,
		ThemeName:         req.ThemeName,
		CustomColor:       string(customColor),
		DefaultAppearance: req.DefaultAppearance,
	})
	if err != nil {
		c.EgoJsonI18N(err)
		return
	}
	c.JSONOK()
}
