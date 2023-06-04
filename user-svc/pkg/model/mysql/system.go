package mysql

import (
	"encoding/json"
	"fmt"

	commonv1 "ecodepost/pb/common/v1"
	communityv1 "ecodepost/pb/community/v1"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	systemKeyHomeOption = "HOME_OPTION"
	systemKeyBase       = "BASE"
	systemKeyTheme      = "THEME"
)

type System struct {
	ID          int64  `gorm:"not null;primary_key;AUTO_INCREMENT" json:"id"`
	Key         string `gorm:"not null;"`
	ValueType   int64  `gorm:"not null"` // 1 json, 2 int64 3 string
	ValueJson   string `gorm:"not null"`
	ValueInt64  int64  `gorm:"not null"`
	ValueString string `gorm:"not null"`
}

func (System) TableName() string {
	return "system"
}

type SystemBase struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Logo           string           `json:"logo"`
	Access         commonv1.CMT_ACS `json:"access"`
	GongxinbuBeian string           `json:"gongxinbuBeian"`
	GonganBeian    string           `json:"gonganBeian"`
}

type SystemHomeOption struct {
	IsSetHome                bool                   `json:"isSetHome"`
	IsSetBanner              bool                   `json:"isSetBanner"`           // 是否启用banner
	ArticleSortByLogin       commonv1.CMN_FILE_SORT `json:"articleSortByLogin"`    // 登录用户推荐内容排序规则
	ArticleSortByNotLogin    commonv1.CMN_FILE_SORT `json:"articleSortByNotLogin"` // 未登录用户推荐内容排序规则
	ArticleHotShowSum        int32                  `json:"articleHotShowSum"`     // 展示热门帖子的数量
	ArticleHotShowWithLatest int32                  `json:""`                      // 展示近期多少天内创建的帖子
	ActivityLatestShowSum    int32                  `json:""`                      // 展示近期的活动数量
	BannerImg                string                 `json:""`                      // 启用首页banner
	BannerTitle              string                 `json:""`                      // banner文案
	BannerLink               string                 `json:""`                      // banner的跳转链接
	DefaultPageByNewUser     string                 `json:""`                      // 新用户注册默认访问页面
	DefaultPageByNotLogin    string                 `json:""`                      // 未登录用户默认访问页面
	DefaultPageByLogin       string                 `json:""`                      // 登录用户默认访问页面
}

func (v SystemHomeOption) ToPb() *communityv1.GetHomeOptionRes {
	return &communityv1.GetHomeOptionRes{
		IsSetHome:                v.IsSetHome,
		IsSetBanner:              v.IsSetBanner,
		ArticleSortByLogin:       v.ArticleSortByLogin,
		ArticleSortByNotLogin:    v.ArticleSortByNotLogin,
		ArticleHotShowSum:        v.ArticleHotShowSum,
		ArticleHotShowWithLatest: v.ArticleHotShowWithLatest,
		BannerImg:                v.BannerImg,
		BannerTitle:              v.BannerTitle,
		BannerLink:               v.BannerLink,
		DefaultPageByNewUser:     v.DefaultPageByNewUser,
		DefaultPageByNotLogin:    v.DefaultPageByNotLogin,
		DefaultPageByLogin:       v.DefaultPageByLogin,
	}
}

type SystemTheme struct {
	ThemeName         string      `gorm:"not null;default:''" json:"themeName"`
	IsCustom          bool        `gorm:"not null;default:false;comment:是否自定义" json:"isCustom"`
	CustomColor       CustomColor `gorm:"type:json;not null;comment:自定义颜色" json:"customColor"`
	DefaultAppearance string      `json:"default_appearance"`
}

func InitSystem(db *gorm.DB) (err error) {
	db.Create(&System{
		Key:       systemKeyBase,
		ValueType: 1,
		ValueJson: "{}",
	})
	db.Create(&System{
		Key:       systemKeyTheme,
		ValueType: 1,
		ValueJson: "{}",
	})
	db.Create(&System{
		Key:       systemKeyHomeOption,
		ValueType: 1,
		ValueJson: "{}",
	})
	return nil
}

func getSystemValueJson(db *gorm.DB, key string) (info System, err error) {
	err = db.Where("`key` =  ? and value_type = ? ", key, 1).Find(&info).Error
	if err != nil {
		err = fmt.Errorf("getSystemValueJson fail, err: %w", err)
		return
	}
	return
}

func putSystemValueJson(db *gorm.DB, key string, value string) (err error) {
	err = db.Model(System{}).Where("`key` =  ? and value_type = ? ", key, 1).Update("value_json", value).Error
	if err != nil {
		err = fmt.Errorf("getSystemValueJson fail, err: %w", err)
		return
	}
	return
}

func GetSystemBase(db *gorm.DB) (res SystemBase, err error) {
	info, err := getSystemValueJson(db, systemKeyBase)
	if err != nil {
		err = fmt.Errorf("GetSystemBase fail, err: %w", err)
		return
	}
	err = json.Unmarshal([]byte(info.ValueJson), &res)
	if err != nil {
		err = fmt.Errorf("GetSystemBase fail2, err: %w", err)
		return
	}
	return
}

func PutSystemBase(db *gorm.DB, req *communityv1.UpdateReq) (err error) {
	var res SystemBase
	info, err := getSystemValueJson(db.Clauses(clause.Locking{Strength: "UPDATE"}), systemKeyBase)
	if err != nil {
		err = fmt.Errorf("PutSystemBase fail, err: %w", err)
		return
	}
	// todo access，后续可以设置
	res.Access = commonv1.CMT_ACS_OPEN
	err = json.Unmarshal([]byte(info.ValueJson), &res)
	if err != nil {
		err = fmt.Errorf("PutSystemBase fail2, err: %w", err)
		return
	}
	if req.Name != nil {
		res.Name = req.GetName()
	}
	if req.Description != nil {
		res.Description = req.GetDescription()
	}
	if req.Logo != nil {
		res.Logo = req.GetLogo()
	}
	if req.GongxinbuBeian != nil {
		res.GongxinbuBeian = req.GetGongxinbuBeian()
	}
	if req.GonganbuBeian != nil {
		//res.GonganBeian = req.GetGonganBeian()
	}

	jsonByte, err := json.Marshal(res)
	if err != nil {
		err = fmt.Errorf("PutSystemBase json marshal fail, err: %w", err)
		return
	}
	err = putSystemValueJson(db, systemKeyBase, string(jsonByte))
	if err != nil {
		err = fmt.Errorf("PutSystemBase putSystemValueJson fail, err: %w", err)
		return
	}
	return
}

func GetSystemHomeOption(db *gorm.DB) (res SystemHomeOption, err error) {
	info, err := getSystemValueJson(db, systemKeyHomeOption)
	if err != nil {
		err = fmt.Errorf("GetSystemHomeOption fail, err: %w", err)
		return
	}
	err = json.Unmarshal([]byte(info.ValueJson), &res)
	if err != nil {
		err = fmt.Errorf("GetSystemHomeOption fail2, err: %w", err)
		return
	}
	return
}

func PutSystemHomeOption(db *gorm.DB, req *communityv1.PutHomeOptionReq) (err error) {
	var res SystemHomeOption
	info, err := getSystemValueJson(db.Clauses(clause.Locking{Strength: "UPDATE"}), systemKeyHomeOption)
	if err != nil {
		err = fmt.Errorf("PutSystemHomeOption fail, err: %w", err)
		return
	}
	err = json.Unmarshal([]byte(info.ValueJson), &res)
	if err != nil {
		err = fmt.Errorf("PutSystemHomeOption fail2, err: %w", err)
		return
	}

	if req.IsSetHome != nil {
		res.IsSetHome = req.GetIsSetHome()
		res.ArticleSortByLogin = req.GetArticleSortByLogin()
		res.ArticleSortByNotLogin = req.GetArticleSortByNotLogin()
		if req.ArticleHotShowSum != nil {
			res.ArticleHotShowSum = req.GetArticleHotShowSum()
		}
		if req.ArticleHotShowWithLatest != nil {
			res.ArticleHotShowWithLatest = req.GetArticleHotShowWithLatest()
		}
	}
	if req.IsSetBanner != nil {
		res.IsSetBanner = req.GetIsSetBanner()
		res.BannerLink = req.GetBannerLink()
	}
	if req.BannerImg != nil {
		res.BannerImg = req.GetBannerImg()
	}
	if req.BannerTitle != nil {
		res.BannerTitle = req.GetBannerTitle()
	}
	if req.BannerLink != nil {
		res.BannerLink = req.GetBannerLink()
	}
	if req.DefaultPageByNewUser != nil {
		res.DefaultPageByNewUser = req.GetDefaultPageByNewUser()
	}
	if req.DefaultPageByLogin != nil {
		res.DefaultPageByLogin = req.GetDefaultPageByLogin()
	}
	if req.DefaultPageByNotLogin != nil {
		res.DefaultPageByNotLogin = req.GetDefaultPageByNotLogin()
	}
	jsonByte, err := json.Marshal(res)
	if err != nil {
		err = fmt.Errorf("PutSystemHomeOption json marshal fail, err: %w", err)
		return
	}
	err = putSystemValueJson(db, systemKeyHomeOption, string(jsonByte))
	if err != nil {
		err = fmt.Errorf("PutSystemHomeOption putSystemValueJson fail, err: %w", err)
		return
	}
	return
}
