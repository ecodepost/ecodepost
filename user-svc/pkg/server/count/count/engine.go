package count

import (
	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"ecodepost/user-svc/pkg/server/count/cache"
	"ecodepost/user-svc/pkg/util"
	"github.com/ego-component/eredis"
	"github.com/gotomicro/ego/core/econf"
)

type Engine struct {
	countConfig map[string][]commonv1.CNT_ACTI
	countDB     *mysql.CountDB
	countCache  *cache.CountCache
	checkCache  *eredis.Component
	limitCache  *eredis.Component
	db          *mysql.CountTable
	baseMinNum  int64 // baseMinNum  打底最小随机值
	baseMaxNum  int64 // baseMaxNum  打底最大随机值
	multipleNum int64 // multipleNum 打底倍率
	LimitNum    int64 // LimitNum 频率限制次数
}

func New() *Engine {
	eng := &Engine{
		countConfig: make(map[string][]commonv1.CNT_ACTI),
		countDB:     mysql.NewCountDB(invoker.Db),
		countCache:  cache.NewCountCache(invoker.Redis),
		checkCache:  invoker.Redis,
		limitCache:  invoker.Redis,
		db:          mysql.NewCountTable(),
		baseMinNum:  econf.GetInt64("common.base_num.min"),
		baseMaxNum:  econf.GetInt64("common.base_num.max"),
		multipleNum: econf.GetInt64("common.base_num.multiple"),
		LimitNum:    econf.GetInt64("common.rate_limit.content"),
	}
	// SetConfig 设置配置信息
	eng.countConfig = mysql.GetConfig()
	return eng
}

// CheckConfig 检测配置信息
func (eng *Engine) checkConfig(biz commonv1.CMN_BIZ, act commonv1.CNT_ACT, acti commonv1.CNT_ACTI) bool {
	key := util.BizActKey(biz, act)
	// 检查BizKey是否存在
	cfgActis, ok := eng.countConfig[key]
	if !ok {
		return false
	}
	// 检查Acti是否存在
	for _, a := range cfgActis {
		if a == acti {
			return true
		}
	}
	return false
}
