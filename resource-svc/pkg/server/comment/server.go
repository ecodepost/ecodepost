package comment

import (
	"context"
	"errors"
	"fmt"
	"html"
	"log"
	"time"

	commentv1 "ecodepost/pb/comment/v1"
	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	notifyv1 "ecodepost/pb/notify/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/dao"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"
	"ecodepost/resource-svc/pkg/utils"
	"github.com/spf13/cast"

	// trackv1 "ecodepost/pb/track/v1"
	userv1 "ecodepost/pb/user/v1"

	"github.com/gotomicro/ego/core/elog"
	"gorm.io/gorm"
)

// GrpcServer ...
type GrpcServer struct{}

var _ commentv1.CommentServer = (*GrpcServer)(nil)

func (GrpcServer) Create(ctx context.Context, req *commentv1.CreateReq) (*commentv1.CreateRes, error) {
	commentAuthor, err := invoker.GrpcUser.Info(ctx, &userv1.InfoReq{Uid: req.Uid})
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("GrpcUser Info fail, err:" + err.Error())
	}

	switch req.GetBizType() {
	case commonv1.CMN_BIZ_ARTICLE:
		if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_CREATE_COMMENT, req.GetUid(), req.GetBizGuid()); err != nil {
			return nil, err
		}
	}

	// 如果传入了时间，那么就说明是同步数据
	nowTime := time.Now().Unix()
	if req.GetCtime() > 0 {
		nowTime = req.GetCtime()
	}

	req.Content = html.EscapeString(req.Content)
	data := &mysql.CommentSubject{
		BizGuid: req.GetBizGuid(),
		BizType: req.GetBizType(),
		Status:  1,
		Ctime:   nowTime,
		Utime:   nowTime,
	}
	commentSubjectId, err := service.CommentSubjectCreate(invoker.Db.WithContext(ctx), data)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("CreateCommentSubject Err1").WithMetadata(map[string]string{"err": err.Error()})
	}

	spaceGuid, err := mysql.FileSpaceGuidByGuid(invoker.Db.WithContext(ctx), req.GetBizGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("DeleteDocument FileSpaceGuidByGuid fail, err: " + err.Error())
	}
	fmt.Printf("spaceGuid--------------->"+"%+v\n", spaceGuid)

	// 增加评论
	if req.GetCommentGuid() == "" {
		resp, err := service.Comment.Create(ctx, req, commentSubjectId)
		switch req.GetBizType() {
		case commonv1.CMN_BIZ_ARTICLE, commonv1.CMN_BIZ_QUESTION, commonv1.CMN_BIZ_ANSWER:
			f, e := mysql.FileByGuid(invoker.Db, req.GetBizGuid())
			if e != nil {
				return nil, errcodev1.ErrDbError().WithMessage("FileByGuid fail").WithMetadata(map[string]string{"err": e.Error()})
			}
			_, e = invoker.GrpcNotify.SendMsg(ctx, &notifyv1.SendMsgReq{
				TplId: 3, // 默认站内信模板id
				Msgs:  []*notifyv1.Msg{{Receiver: cast.ToString(f.CreatedBy)}},
				VarLetter: &notifyv1.Letter{
					Type:     commonv1.NTF_TYPE_NEW_COMMENT,
					TargetId: req.BizGuid,
					Meta:     utils.NewMeta(utils.Meta{"commentAuthor": commentAuthor, "fileTitle": f.Name}),
				},
			})
			if e != nil {
				elog.Warn("CreateLetter fail", elog.FieldErr(err))
			}
		}

		switch req.GetBizType() {
		case commonv1.CMN_BIZ_ARTICLE, commonv1.CMN_BIZ_QUESTION, commonv1.CMN_BIZ_ANSWER:
			err = mysql.FileUpdateExpr(invoker.Db.WithContext(ctx), req.GetBizGuid(), "cnt_comment", 1)
			if err != nil {
				return nil, fmt.Errorf("create comment track total fail2, err: %w", err)
			}
		}
		return resp, err
	}

	// 如果是回复评论，并且是评论模式，那么可以这么操作
	if req.GetCommentGuid() != "" && req.ActionType == commonv1.FILE_ACT_COMMENT {
		cr, err := service.ReplyCommentByGuid(ctx, req)
		if err != nil {
			return nil, err
		}

		switch req.GetBizType() {
		case commonv1.CMN_BIZ_ARTICLE, commonv1.CMN_BIZ_QUESTION, commonv1.CMN_BIZ_ANSWER:
			err = mysql.FileUpdateExpr(invoker.Db.WithContext(ctx), req.GetBizGuid(), "cnt_comment", 1)
			if err != nil {
				return nil, fmt.Errorf("create comment track total fail2, err: %w", err)
			}
		}
		return cr, nil
	}

	return nil, errcodev1.ErrInvalidArgument().WithMessage("不支持该方式")
}

func (c GrpcServer) List(ctx context.Context, req *commentv1.ListReq) (resp *commentv1.ListRes, err error) {
	// 获取subject信息
	subjectId, cntComment, err := c.singleSubjectInfoByGuid(invoker.Db.WithContext(ctx), req.GetBizGuid(), req.GetBizType())
	if err != nil {
		return nil, err
	}

	indexList, commentMap, err := service.Index.GetIndexListByMySQLByGuid(ctx, subjectId, req.Pagination)
	resp = &commentv1.ListRes{}
	for _, v := range indexList {
		commentInfo := commentMap[v]
		if commentInfo.CntChildComment > 0 {
			list, _, hasMore, err := service.Index.GetSubIndexListByReplyToRootIdByGuid(ctx, invoker.Db.WithContext(ctx), commentInfo.CommentGuid)
			if err != nil {
				// 因为被删除，所以找不到，那么不能continue
				elog.Error("GetSubIndexListByReplyToRootId error", elog.FieldErr(err))
			}
			if list != nil {
				commentInfo.Children = list
				commentInfo.HasMoreChildComment = hasMore
			}

		}
		resp.List = append(resp.List, commentInfo)
	}

	resp.Pagination = req.Pagination
	resp.CntComment = cntComment
	return
}

func (GrpcServer) ChildList(ctx context.Context, req *commentv1.ChildListReq) (resp *commentv1.ChildListRes, err error) {
	indexList, commentMap, err := service.Index.GetChildIndexListByMySQLByGuid(ctx, req.GetCommentGuid(), req.Pagination)
	if err != nil {
		return nil, err
	}

	resp = &commentv1.ChildListRes{
		Pagination: req.Pagination,
	}
	for _, v := range indexList {
		resp.List = append(resp.List, commentMap[v])
	}
	return
}

func (GrpcServer) Delete(ctx context.Context, req *commentv1.DeleteReq) (*commentv1.DeleteRes, error) {
	switch req.GetBizType() {
	case commonv1.CMN_BIZ_ARTICLE:
		if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_DELETE_COMMENT, req.GetUid(), req.GetBizGuid()); err != nil {
			return nil, err
		}
	}

	spaceGuid, err := mysql.FileSpaceGuidByGuid(invoker.Db.WithContext(ctx), req.GetBizGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("DeleteDocument FileSpaceGuidByGuid fail, err: " + err.Error())
	}
	fmt.Printf("spaceGuid--------------->"+"%+v\n", spaceGuid)

	resp, err := dao.CommentIndexDelete(ctx, req.Uid, req.GetCommentGuid(), req.GetBizGuid(), req.BizType, req.DeleteType)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("DeleteComment err:" + err.Error())
	}
	log.Printf("resp--------------->"+"%+v\n", resp)
	// TODO mqmsg
	// msg := mqmsgv1.CommentUpdate{
	// 	Uid:     req.Uid,
	// 	BizGuid: resp.BizGuid,
	// 	BizType: resp.BizType,
	// 	Delta:   -1,
	// }
	// updateArticleCommentByGuid(ctx, &msg)

	// _, err = invoker.GrpcTrack.Total(ctx, &trackv1.TotalReq{
	// 	Event:     commonv1.TRACK_TOTAL_EVENT_DELETE_COMMENT,
	// 	Tid:       etrace.ExtractTraceID(ctx),
	// 	Uid:       req.Uid,
	// 	FileGuid:  req.GetBizGuid(),
	// 	CmtGuid:   req.GetCmtGuid(),
	// 	SpaceGuid: spaceGuid,
	// 	Ts:        time.Now().UnixMilli(),
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("create comment track total fail, err: %w", err)
	// }

	return &commentv1.DeleteRes{}, nil
}

func (GrpcServer) ListByUser(ctx context.Context, req *commentv1.ListByUserReq) (resp *commentv1.ListByUserRes, err error) {
	indexList, commentMap, err := service.Index.GetIndexListByMysqlUserByGuid(ctx, req.Uid, req.Pagination)
	if err != nil {
		return nil, err
	}

	resp = &commentv1.ListByUserRes{Pagination: req.Pagination}
	for _, v := range indexList {
		resp.List = append(resp.List, commentMap[v])
	}
	return
}

// todo 创建评论的时候，删除索引会有并发问题，所以现在不在使用这个缓存，直接通过mysql获取数据
// 会出现为0的情况
func (GrpcServer) singleSubjectInfoByGuid(db *gorm.DB, bizGuid string, bizType commonv1.CMN_BIZ) (int64, int32, error) {
	subjectId, cntComment, err := service.Subject.GetInfoByMySQLByGuid(db, bizGuid, bizType)
	if err != nil && !errors.Is(err, service.SubjectNotFound) {
		return 0, 0, errcodev1.ErrUnknown().WithMessage("singleSubjectInfo Err2").WithMetadata(map[string]string{"err": err.Error()})
	}
	if errors.Is(err, service.SubjectNotFound) {
		return 0, 0, errcodev1.ErrNotFound()
	}
	return subjectId, cntComment, nil
}
