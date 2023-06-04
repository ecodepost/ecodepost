package logger

import (
	"context"
	"time"

	"ecodepost/user-svc/pkg/invoker"
	"ecodepost/user-svc/pkg/model/mysql"
	"ecodepost/user-svc/pkg/service"
	"gorm.io/datatypes"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	loggerv1 "ecodepost/pb/logger/v1"
	"github.com/ego-component/egorm"
)

type GrpcServer struct{}

var _ loggerv1.LoggerServer = &GrpcServer{}

func (GrpcServer) Create(ctx context.Context, req *loggerv1.CreateReq) (*loggerv1.CreateRes, error) {
	err := mysql.LoggerCreate(invoker.Db.WithContext(ctx), &mysql.Logger{
		Event:          req.GetEvent(),
		Group:          req.GetGroup(),
		TargetUid:      req.GetTargetUid(),
		OperateUid:     req.GetOperateUid(),
		SpaceGuid:      req.GetSpaceGuid(),
		SpaceGroupGuid: req.GetSpaceGroupGuid(),
		FileGuid:       req.GetFileGuid(),
		Metadata:       datatypes.JSON(req.GetMetadata()),
		Ctime:          time.Now().Unix(),
	})
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("logger create fail, err: " + err.Error())
	}
	return &loggerv1.CreateRes{}, nil
}

// BatchCreate 批量增加日志事件
func (GrpcServer) BatchCreate(ctx context.Context, req *loggerv1.BatchCreateReq) (*loggerv1.BatchCreateRes, error) {
	createBatchData := make([]mysql.Logger, 0)
	nowTime := time.Now().Unix()
	for _, value := range req.GetList() {
		createBatchData = append(createBatchData, mysql.Logger{
			Event:          value.GetEvent(),
			Group:          value.GetGroup(),
			TargetUid:      value.GetTargetUid(),
			OperateUid:     value.GetOperateUid(),
			SpaceGuid:      value.GetSpaceGuid(),
			SpaceGroupGuid: value.GetSpaceGroupGuid(),
			FileGuid:       value.GetFileGuid(),
			Metadata:       datatypes.JSON(value.GetMetadata()),
			Ctime:          nowTime,
		})
	}
	err := mysql.LoggerCreateInBatches(invoker.Db.WithContext(ctx), createBatchData)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("logger create fail, err: " + err.Error())
	}
	return &loggerv1.BatchCreateRes{}, nil
}

// BatchCreateByTargetUids 批量增加日志ByTargetUids
func (GrpcServer) BatchCreateByTargetUids(ctx context.Context, req *loggerv1.BatchCreateByTargetUidsReq) (*loggerv1.BatchCreateByTargetUidsRes, error) {
	createBatchData := make([]mysql.Logger, 0)
	nowTime := time.Now().Unix()
	for _, targetUid := range req.GetTargetUids() {
		createBatchData = append(createBatchData, mysql.Logger{
			Event:          req.GetEvent(),
			Group:          req.GetGroup(),
			TargetUid:      targetUid,
			OperateUid:     req.GetOperateUid(),
			SpaceGuid:      req.GetSpaceGuid(),
			SpaceGroupGuid: req.GetSpaceGroupGuid(),
			FileGuid:       req.GetFileGuid(),
			Metadata:       datatypes.JSON(req.GetMetadata()),
			Ctime:          nowTime,
		})
	}
	err := mysql.LoggerCreateInBatches(invoker.Db.WithContext(ctx), createBatchData)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("logger create fail, err: " + err.Error())
	}
	return &loggerv1.BatchCreateByTargetUidsRes{}, nil
}

func (GrpcServer) ListPage(ctx context.Context, req *loggerv1.ListPageReq) (*loggerv1.ListPageRes, error) {
	conds := egorm.Conds{}
	if req.GetSearchEvent() != commonv1.LOG_EVENT_INVALID {
		conds["event"] = int32(req.GetSearchEvent())
	}
	if req.GetSearchGroup() != commonv1.LOG_GROUP_INVALID {
		conds["group"] = int32(req.GetSearchGroup())
	}
	if req.GetSearchOperateUid() != 0 {
		conds["operate_uid"] = req.GetSearchOperateUid()
	}
	if req.GetSearchTargetUid() != 0 {
		conds["target_uid"] = req.GetSearchTargetUid()
	}
	list, err := mysql.LoggerListPage(invoker.Db.WithContext(ctx), conds, req.GetPagination())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("fail1, err: " + err.Error())
	}
	uMap, err := mysql.UserMap(invoker.Db.WithContext(ctx), list.ToAllUids())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Map Fail, err: " + err.Error())
	}
	for _, value := range list {
		eventInfo, err := service.Logger.GetEvent(value.Event)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("get event fail, err:" + err.Error())
		}
		groupInfo, err := service.Logger.GetGroup(value.Group)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("get group fail, err:" + err.Error())
		}

		value.EventName = eventInfo.CNName
		value.GroupName = groupInfo.CNName
		value.Message = eventInfo.CNName
		operateUserInfo, flag := uMap[value.OperateUid]
		if flag {
			value.OperateNickname = operateUserInfo.GetNickname()
			value.OperateAvatar = operateUserInfo.GetAvatar()
		}
		if value.TargetUid > 0 {
			targetUserInfo, flag2 := uMap[value.TargetUid]
			if flag2 {
				value.TargetNickname = targetUserInfo.GetNickname()
				value.TargetAvatar = targetUserInfo.GetAvatar()
			}
		}
	}
	return &loggerv1.ListPageRes{
		List:       list.ToPb(),
		Pagination: req.GetPagination(),
	}, nil
}

// ListEventAndGroup 显示日志列表
func (GrpcServer) ListEventAndGroup(ctx context.Context, req *loggerv1.ListEventAndGroupReq) (*loggerv1.ListEventAndGroupRes, error) {
	return &loggerv1.ListEventAndGroupRes{
		EventList: service.LoggerEventCNPbList,
		GroupList: service.LoggerGroupCNPbList,
	}, nil
}
