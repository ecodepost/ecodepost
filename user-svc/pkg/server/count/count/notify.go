package count

import (
	commonv1 "ecodepost/pb/common/v1"
	countv1 "ecodepost/pb/count/v1"
	"ecodepost/user-svc/pkg/code"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// HandleIncr 只处理ADD/SUB操作
func (eng *Engine) HandleIncr(ctx context.Context, in *countv1.SetReq) (data *countv1.SetRes, err error) {
	defer eng.handleHook(ctx, in)

	// 更新redis
	data, err = eng.countCache.Incr(ctx, in)
	if err != nil {
		elog.Error("redis Incr fail", elog.FieldErr(err))
		return nil, code.IncrError.WithMetadata(map[string]string{"err": err.Error()})
	}

	// 更新数据库
	if in.Act == commonv1.CNT_ACT_VIEW || in.Act == commonv1.CNT_ACT_CNT {
		err = eng.RepeatedUpdate(ctx, in.Tid, in.Fid, in.Biz, in.Act, data.Num, data.RealNum, in.Ct, in.Did, in.Ip)
		if err != nil {
			elog.Error("RepeatedUpdate fail", elog.FieldErr(err))
			return nil, code.IncrError.WithMetadata(map[string]string{"err": err.Error()})
		}
	} else {
		err = eng.IncrOnce(invoker.Db.WithContext(ctx), in.Tid, in.Fid, in.Biz, in.Act, data.Num, data.RealNum, in.Ct, in.Did, in.Ip)
		if err != nil {
			elog.Error("IncrOnce fail", elog.FieldErr(err))
			return
		}
	}

	return
}

func (eng *Engine) HandleUpdate(ctx context.Context, in *countv1.SetReq) (data *countv1.SetRes, err error) {
	defer eng.handleHook(ctx, in)

	// 更新redis
	data, err = eng.countCache.Update(ctx, in)
	if err != nil {
		return nil, code.UpdateError.WithMetadata(map[string]string{"err": err.Error()})
	}

	// 写数据库 TODO 是否需要限定ACT类型？
	err = eng.RepeatedUpdate(ctx, in.Tid, in.Fid, in.Biz, in.Act, data.Num, data.RealNum, in.Ct, in.Did, in.Ip)
	if err != nil {
		elog.Error("Update.RepeatedUpdate fail", elog.FieldErr(err), elog.String("tid", in.Tid), elog.Int64("num", data.Num), elog.Int64("real_num", data.RealNum), elog.Any("in", in))
		return
	}
	elog.Info("Update.success", elog.String("tid", in.Tid), elog.Int64("num", data.Num), elog.Int64("real_num", data.RealNum), zap.Any("in", in))

	return
}

func (eng *Engine) HandleReset(ctx context.Context, in *countv1.SetReq) (data *countv1.SetRes, err error) {
	defer eng.handleHook(ctx, in)

	// 更新redis缓存
	data, err = eng.countCache.Reset(ctx, in)
	if err != nil {
		return nil, code.ResetError.WithMetadata(map[string]string{"err": err.Error()})
	}

	// TODO 更新数据库
	return
}

func (eng *Engine) handleHook(ctx context.Context, in *countv1.SetReq) {
	// 日志记录计数流水
	elog.Info("Engine.DealCountEvent.success", elog.String("tid", in.Tid), elog.String("fid", in.Fid),
		elog.String("biz", in.Biz.String()), elog.String("act", in.Act.String()), elog.String("acti", in.Acti.String()),
		elog.String("ct", in.Ct), elog.String("did", in.Did), elog.String("ip", in.Ip), elog.Int32("val", in.Val),
	)
	// 通知处理
	// go eng.NotifyMQ(in)

	// 记录日志行为
	go eng.AddBehaviorLog(ctx, &mysql.CountBehaviorLog{
		Fid:  in.Fid,
		Biz:  in.Biz,
		Act:  in.Act,
		Acti: in.Acti,
		Tid:  in.Tid,
	})
}

// DealCountEvent 处理计数消息事件
func (eng *Engine) DealCountEvent(ctx context.Context, in *countv1.SetReq) (res *countv1.SetRes, err error) {
	// 默认为+1操作
	if in.Val == 0 {
		in.Val = 1
	}

	switch in.Acti {
	case commonv1.CNT_ACTI_ADD, commonv1.CNT_ACTI_SUB:
		return eng.HandleIncr(ctx, in)
	case commonv1.CNT_ACTI_UPDATE:
		return eng.HandleUpdate(ctx, in)
	case commonv1.CNT_ACTI_RESET:
		return eng.HandleReset(ctx, in)
	default:
		return nil, code.OpTypeNotSupport
	}
}
