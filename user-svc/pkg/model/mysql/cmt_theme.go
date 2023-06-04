package mysql

// 社区自定义话颜色
/*
"theme_settings": {
			"is_custom": true,
			"current_theme": "custom",
			"custom_colors": {
				"active_item_color": "#EF9072",
				"sidebar_text_color": "#7941c0",
				"mention_badge_color": "#4BABAC",
				"sidebar_hover_color": "#FFEADE",
				"active_item_text_color": "#FFFFFF",
				"online_indicator_color": "#17A7BB",
				"sidebar_background_color": "#FEF3ED"
			},
			"default_appearance": "light"
		},
*/

/*
--themeColorPrimary: #05a48a;
    --themeColorStatus: #05a48a;
    --themeColorButtonText: #28354a;
    --themeColorBackground: #f7f9fb;
*/

var DefaultColor = CustomColor{
	ThemeColorPrimary:    "#05a48a",
	ThemeColorStatus:     "#05a48a",
	ThemeColorButtonText: "#28354a",
	ThemeColorBackground: "#f7f9fb",
}

type CustomColor struct {
	ThemeColorPrimary    string `json:"themeColorPrimary"`
	ThemeColorStatus     string `json:"themeColorStatus"`
	ThemeColorButtonText string `json:"themeColorButtonText"`
	ThemeColorBackground string `json:"themeColorBackground"`
	//ActiveItemColor        string `json:"active_item_color"`
	//SidebarTextColor       string `json:"sidebar_text_color"`
	//MentionBadgeColor      string `json:"mention_badge_color"`
	//SidebarHoverColor      string `json:"sidebar_hover_color"`
	//ActiveItemTextColor    string `json:"active_item_text_color"`
	//OnlineIndicatorColor   string `json:"online_indicator_color"`
	//SidebarBackgroundColor string `json:"sidebar_background_color"`
}
