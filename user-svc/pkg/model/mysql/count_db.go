package mysql

import (
	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/user-svc/pkg/util"

	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/econf"
)

var (
	DB *CountDB
)

// NewCountDB 初始化DB
func NewCountDB(db *egorm.Component) *CountDB {
	ct := &CountDB{
		Component:   db,
		baseMinNum:  econf.GetInt64("common.base_num.min"),
		baseMaxNum:  econf.GetInt64("common.base_num.max"),
		multipleNum: econf.GetInt64("common.base_num.multiple"),
	}
	DB = ct
	return ct
}

var (
	countConfigs = make([]CountConfig, 0)
	configMap    = make(map[string][]commonv1.CNT_ACTI)
)

func init() {
	countConfigs = []CountConfig{
		{
			Id:         1,
			Biz:        commonv1.CMN_BIZ_ARTICLE,
			Act:        commonv1.CNT_ACT_LIKE,
			ActCommand: commonv1.CNT_ACTI_ADD,
		},
		{
			Id:         2,
			Biz:        commonv1.CMN_BIZ_ARTICLE,
			Act:        commonv1.CNT_ACT_LIKE,
			ActCommand: commonv1.CNT_ACTI_SUB,
		},
		{
			Id:         3,
			Biz:        commonv1.CMN_BIZ_ARTICLE,
			Act:        commonv1.CNT_ACT_LIKE,
			ActCommand: commonv1.CNT_ACTI_UPDATE,
		},
		{
			Id:         4,
			Biz:        commonv1.CMN_BIZ_ARTICLE,
			Act:        commonv1.CNT_ACT_LIKE,
			ActCommand: commonv1.CNT_ACTI_RESET,
		},
		{
			Id:         5,
			Biz:        commonv1.CMN_BIZ_USER,
			Act:        commonv1.CNT_ACT_FOLLOW,
			ActCommand: commonv1.CNT_ACTI_ADD,
		},
		{
			Id:         6,
			Biz:        commonv1.CMN_BIZ_USER,
			Act:        commonv1.CNT_ACT_FOLLOW,
			ActCommand: commonv1.CNT_ACTI_SUB,
		},
	}

	for _, v := range countConfigs {
		key := util.BizActKey(v.Biz, v.Act)
		_, ok := configMap[key]
		if !ok {
			configMap[key] = make([]commonv1.CNT_ACTI, 0)
		}
		configMap[key] = append(configMap[key], v.ActCommand)
	}
}

// GetConfig 获取计数配置列表
func GetConfig() map[string][]commonv1.CNT_ACTI {
	return DB.getConfig()
}

// CountDB ...
type CountDB struct {
	*egorm.Component
	baseMinNum  int64
	baseMaxNum  int64
	multipleNum int64
}

// getConfig 获取计数配置信息
func (cbd CountDB) getConfig() (res map[string][]commonv1.CNT_ACTI) {
	return configMap
}
