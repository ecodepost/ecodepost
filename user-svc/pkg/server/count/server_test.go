package count

import (
	"context"
	"log"
	"testing"

	commonv1 "ecodepost/pb/common/v1"
	countv1 "ecodepost/pb/count/v1"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	eapp.SetEgoDebug("true")
	// 清除已有数据
	// err := invoker.Db.Model(mysql.U2uStat{}).Where("1=1").Delete(&mysql.U2uStat{}).Error
	// assert.NoError(t, err)
	// err = invoker.Db.Model(mysql.U2uStat{}).Where("1=1").Delete(&mysql.U2uStat{}).Error
	// assert.NoError(t, err)
	var err error
	const fromUid1 = "100"
	const fromUid2 = "200"
	const targetUid = "999"
	const acti = commonv1.CNT_ACTI_ADD
	const act = commonv1.CNT_ACT_LIKE
	const biz = commonv1.CMN_BIZ_USER

	// 发起测试请求
	cli := countv1.NewCountClient(newCC())
	ctx := context.Background()

	// 用户fromUid1点赞了用户targetUid
	res1, err := cli.Set(ctx, &countv1.SetReq{
		Fid:  fromUid1,
		Tid:  targetUid,
		Acti: acti,
		Ip:   "1.1.1.1",
		Biz:  biz,
		Ct:   "1",
		Did:  "1",
		Act:  act,
	})
	assert.NoError(t, err)
	t.Logf("res1------------%+v", res1)

	res8, err := cli.GetTnumByFids(ctx, &countv1.GetTnumByFidsReq{
		Fids: []string{fromUid1, fromUid2},
		Biz:  biz,
		Act:  act,
	})
	assert.NoError(t, err)
	log.Println(`res8--------------->`, res8)

	// 用户fromUid2点赞了用户targetUid
	res2, err := cli.Set(ctx, &countv1.SetReq{
		Fid:  fromUid2,
		Tid:  targetUid,
		Acti: acti,
		Ip:   "1.1.1.1",
		Biz:  biz,
		Ct:   "1",
		Did:  "1",
		Act:  act,
	})
	assert.NoError(t, err)
	log.Printf("res2--------------->"+"%+v\n", res2)

	// 查询targetUid被点赞详情
	// res3, err := cli.GetTdetail(ctx, &countv1.GetTdetailReq{
	// 	Tid: targetUid,
	// 	Biz: biz,
	// 	Act: act,
	// 	Bid: "of",
	// })
	// assert.NoError(t, err)
	// log.Printf("res3--------------->"+"%+v\n", res3)

	// 查询targetUid被点赞详情，及fromUid1对targetUid状态
	res4, err := cli.GetTdetailsByTids(ctx, &countv1.GetTdetailsByTidsReq{
		Fid:  fromUid1,
		Tids: []string{targetUid},
		Biz:  biz,
		Act:  act,
	})
	assert.NoError(t, err)
	log.Printf("res4--------------->"+"%+v\n", res4)

	// 查询targetUid被点赞详情
	res5, err := cli.GetTdetailsByTids(ctx, &countv1.GetTdetailsByTidsReq{
		Tids: []string{targetUid},
		Biz:  biz,
		Act:  act,
	})
	assert.NoError(t, err)
	log.Println(`res5--------------->`, res5)

	// 查询targetUid被点赞详情
	res6, err := cli.GetTdetailsByFid(ctx, &countv1.GetTdetailsByFidReq{
		Fid:    fromUid1,
		Biz:    biz,
		Act:    act,
		Offset: 0,
		Limit:  10,
	})
	assert.NoError(t, err)
	log.Println(`res6--------------->`, res6)

	// 查询targetUid被点赞详情
	res7, err := cli.GetTnumByBaksAndTids(ctx, &countv1.GetTnumByBaksAndTidsReq{
		Tids: []string{targetUid},
		Baks: []*countv1.BAK{
			{Biz: biz, Act: act},
		},
	})
	assert.NoError(t, err)
	log.Println(`res7--------------->`, res7)

	return
}
