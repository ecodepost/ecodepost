package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	notifyv1 "ecodepost/pb/notify/v1"
	"ecodepost/resource-svc/pkg/utils/x"
	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/elog"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type GrpcServer struct{}

var _ notifyv1.NotifyServer = &GrpcServer{}

type S2S = map[string]string

// SendMsg 发送消息
func (*GrpcServer) SendMsg(ctx context.Context, req *notifyv1.SendMsgReq) (*notifyv1.SendMsgRes, error) {
	var (
		response = &notifyv1.SendMsgRes{MsgResults: make([]*notifyv1.MsgResult, 0, len(req.Msgs))}
		err      error
	)

	// 查询消息模板渠道数据
	channelInfo, err := getChannel(ctx, req)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("getChannel fail," + err.Error())
	}

	for _, msg := range req.Msgs {
		ret := &notifyv1.MsgResult{Code: 0, ExtraId: msg.ExtraId}
		response.MsgResults = append(response.MsgResults, ret)
		// 构造MsgReceiverInfo
		if _, err = CreateNotify(msg, req, channelInfo); err != nil {
			SetErr(ret, err)
			continue
		}
	}
	return response, nil
}

func SetErr(ret *notifyv1.MsgResult, err error) {
	elog.Error("db insert error", elog.FieldErr(err))
	ret.Reason = err.Error()
	ret.Code = 500
}

// CreateNotify 创建通知
func CreateNotify(msg *notifyv1.Msg, req *notifyv1.SendMsgReq, channelId commonv1.NOTIFY_CHANNEL) (*mysql.Notify, error) {
	var (
		jsonByte []byte
		err      error
	)
	switch channelId {
	case commonv1.NOTIFY_CHANNEL_EMAIL_COMMON:
		jsonByte, err = json.Marshal(req.VarEmail)
		if err != nil {
			return nil, fmt.Errorf("email marshal fail, err: %w", err)
		}
	case commonv1.NOTIFY_CHANNEL_SMS_ALI:
		jsonByte, err = json.Marshal(req.VarSms)
		if err != nil {
			return nil, fmt.Errorf("email marshal fail, err: %w", err)
		}
	case commonv1.NOTIFY_CHANNEL_LETTER:
		jsonByte, err = json.Marshal(req.VarLetter)
		if err != nil {
			return nil, fmt.Errorf("email marshal fail, err: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported channelType, %d", channelId)
	}

	data := &mysql.Notify{
		MsgTmplId:    int(req.TplId),
		Channel:      channelId,
		Status:       commonv1.NOTIFY_STATUS_INIT,
		ExtraId:      msg.ExtraId,
		ExtraContent: msg.ExtraContent,
		TplData:      msg.TplData,
		Sender:       msg.Sender,
		Vars:         datatypes.JSON(jsonByte),
	}
	if msg.Vars != nil {
		bytes, err := json.Marshal(msg.Vars)
		if err != nil {
			elog.Error("json Marshal vars error", elog.FieldErr(err))
		} else {
			data.Vars = bytes
		}
	}
	switch channelId {
	case commonv1.NOTIFY_CHANNEL_EMAIL_COMMON:
		data.Email = msg.Receiver
	case commonv1.NOTIFY_CHANNEL_SMS_ALI:
		data.Phone = msg.Receiver
	case commonv1.NOTIFY_CHANNEL_LETTER:
		data.Uid = cast.ToInt64(msg.Receiver)
	default:
		return nil, fmt.Errorf("unsupported channelType, %d", channelId)
	}
	return data, mysql.NotifyCreate(invoker.Db, data)
}

func getChannel(ctx context.Context, msgs *notifyv1.SendMsgReq) (commonv1.NOTIFY_CHANNEL, error) {
	tplCh := mysql.NotifyTplChannel{}
	err := invoker.Db.WithContext(ctx).
		Where("`tpl_id` = ? and ch_status = ?", msgs.TplId, uint8(commonv1.NOTIFY_CS_SUCCESS)).
		Find(&tplCh).Error
	if err != nil {
		return commonv1.NOTIFY_CHANNEL_INVALID, fmt.Errorf("find channel fail, %w", err)
	}

	return commonv1.NOTIFY_CHANNEL(tplCh.ChannelId), nil
}

// UpdateReadStatus 更改通知状态
func (*GrpcServer) UpdateReadStatus(ctx context.Context, req *notifyv1.UpdateReadStatusReq) (*notifyv1.UpdateReadStatusRes, error) {
	if req.GetUid() <= 0 {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("uid cant empty")
	}
	if req.GetNotificationId() <= 0 {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("notification id cant empty")
	}
	conds := egorm.Conds{"uid": req.Uid, "notification_id": req.NotificationId}
	oriUserNotification, err := mysql.NotifyInfoX(invoker.Db.WithContext(ctx), conds)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("ChangeUserNotificationStatus: UserNotificationInfoX fail, err: " + err.Error())
	}

	// userNotification 不存在时直接返回
	if oriUserNotification.ID <= 0 {
		elog.Warn("ChangeUserNotificationStatus nid not exist", zap.Int64("nid", req.NotificationId))
		return &notifyv1.UpdateReadStatusRes{}, nil
	}
	if oriUserNotification.Status == commonv1.NTF_STATUS_READED {
		return &notifyv1.UpdateReadStatusRes{}, nil
	}
	if err := mysql.NotifyUpdate(invoker.Db.WithContext(ctx), oriUserNotification.ID, map[string]any{
		"status": int32(commonv1.NTF_STATUS_READED),
	}); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("ChangeUserNotificationStatus: UserNotificationUpdate fail, err: " + err.Error())
	}
	return &notifyv1.UpdateReadStatusRes{}, nil
}

// GetUnViewCount 查询用户未查看的通知数量
func (*GrpcServer) GetUnViewCount(ctx context.Context, req *notifyv1.GetUnViewCountReq) (*notifyv1.GetUnViewCountRes, error) {
	if req.Uid <= 0 {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("GetUnViewCount: uid is empty")
	}
	conds := egorm.Conds{"uid": req.Uid}
	viewRecord, err := mysql.NotifyUserTsInfoX(invoker.Db.WithContext(ctx), conds)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetUnViewCount: UserTimestampInfoX fail, err: " + err.Error())
	}

	conds["ctime"] = egorm.Cond{Op: ">", Val: viewRecord.ViewTimestamp}
	count, err := mysql.NotifyCount(invoker.Db.WithContext(ctx), conds)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("GetUnViewCount: UserNotificationCount failed").WithMetadata(S2S{"err": err.Error()})
	}
	return &notifyv1.GetUnViewCountRes{Count: count}, nil
}

// ListUserNotification 查询用户下的通知列表
func (*GrpcServer) ListUserNotification(ctx context.Context, req *notifyv1.ListUserNotificationReq) (*notifyv1.ListUserNotificationRes, error) {
	if req.Uid <= 0 {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("uid cant empty")
	}
	var conds = egorm.Conds{}
	var uid = req.GetUid()
	conds["uid"] = uid

	if len(req.Types) != 0 {
		conds["type"] = x.Es2I32s(req.Types)
	}

	list, err := mysql.NotifyListPage(invoker.Db.WithContext(ctx), conds, req.GetPagination())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("ListUserNotification: UserNotificationListPage fail, err: " + err.Error())
	}
	listUserNotificationItems := make([]*notifyv1.ListUserNotificationItem, 0)
	for _, userNotification := range list {
		listUserNotificationItems = append(listUserNotificationItems, &notifyv1.ListUserNotificationItem{
			Type:             userNotification.Type,
			TargetId:         userNotification.TargetId,
			Link:             userNotification.Link,
			Meta:             []byte(userNotification.Meta),
			Status:           userNotification.Status,
			NotificationId:   userNotification.NotificationId,
			NotificationTime: userNotification.Ctime,
			SenderId:         userNotification.SenderId,
			Id:               userNotification.ID,
		})
	}
	res := &notifyv1.ListUserNotificationRes{
		Pagination: req.GetPagination(),
		List:       listUserNotificationItems,
	}
	// go func() {
	_ = UpdateUserTsByView(ctx, req.GetUid())
	// }()
	return res, nil
}

// ChangeUserAllNotificationStatus 更改用户全部通知消息状态
func (*GrpcServer) ChangeUserAllNotificationStatus(ctx context.Context, req *notifyv1.ChangeUserAllNotificationStatusReq) (*notifyv1.ChangeUserAllNotificationStatusRes, error) {
	if req.Uid <= 0 {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("uid cant empty")
	}
	conds := egorm.Conds{"uid": req.Uid}
	tx := invoker.Db.Begin()
	if err := mysql.NotifyUpdateX(tx.WithContext(ctx), conds, map[string]any{"status": uint8(req.Status)}); err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("ChangeUserAllNotificationStatus: UserNotificationUpdateX fail, err: " + err.Error())
	}
	if req.Status != commonv1.NTF_STATUS_READED {
		tx.Commit()
		return &notifyv1.ChangeUserAllNotificationStatusRes{}, nil
	}

	oriUserTimestamp, err := mysql.NotifyUserTsInfoX(tx.WithContext(ctx), conds)
	if err != nil {
		tx.Rollback()
		return nil, errcodev1.ErrDbError().WithMessage("ChangeUserAllNotificationStatus: UserTimestampInfoX fail, err: " + err.Error())
	}
	// create
	nowTime := time.Now().Unix()
	if oriUserTimestamp.ID <= 0 {
		if err := mysql.NotifyUserTsCreate(tx.WithContext(ctx), &mysql.NotifyLetterUserTs{
			Uid:           req.Uid,
			ReadTimestamp: nowTime,
		}); err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("ChangeUserAllNotificationStatus: UserTimestampCreate fail, err: " + err.Error())
		}
	} else {
		// update
		if err := mysql.NotifyUserTsUpdateRead(tx.WithContext(ctx), oriUserTimestamp.ID, nowTime); err != nil {
			tx.Rollback()
			return nil, errcodev1.ErrDbError().WithMessage("ChangeUserAllNotificationStatus: UserTimestampUpdate failed").WithMetadata(S2S{"err": err.Error()})
		}
	}
	tx.Commit()
	return &notifyv1.ChangeUserAllNotificationStatusRes{}, nil
}

// UpdateUserTsByView 查看后更新时间戳，查看列表的时候更新这个数据
func UpdateUserTsByView(ctx context.Context, uid int64) error {
	if uid <= 0 {
		return errcodev1.ErrInvalidArgument().WithMessage("UpdateUserTsByView: uid is empty")
	}
	conds := egorm.Conds{"uid": uid}
	oriUserTimestamp, err := mysql.NotifyUserTsInfoX(invoker.Db.WithContext(ctx), conds)
	if err != nil {
		return errcodev1.ErrDbError().WithMessage("UpdateUserTsByView: UserTimestampInfoX fail,err: " + err.Error())
	}

	// create
	nowTime := time.Now().Unix()
	if oriUserTimestamp.ID <= 0 {
		if err := mysql.NotifyUserTsCreate(invoker.Db.WithContext(ctx), &mysql.NotifyLetterUserTs{
			Uid:           uid,
			ViewTimestamp: nowTime,
		}); err != nil {
			return errcodev1.ErrDbError().WithMessage("UpdateUserTsByView: UserTimestampInfoX fail2,err: " + err.Error())
		}
		return nil
	}
	// update
	if err := mysql.NotifyUserTsUpdateView(invoker.Db.WithContext(ctx), oriUserTimestamp.ID, nowTime); err != nil {
		return errcodev1.ErrDbError().WithMessage("UpdateUserTsByView: UserTimestampUpdate fail3, err: " + err.Error())
	}
	return nil
}
