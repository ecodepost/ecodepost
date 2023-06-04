package question

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"
	"ecodepost/resource-svc/pkg/utils"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	questionv1 "ecodepost/pb/question/v1"
)

type GrpcServer struct{}

var _ questionv1.QuestionServer = (*GrpcServer)(nil)

func (GrpcServer) Create(ctx context.Context, req *questionv1.CreateReq) (*questionv1.CreateRes, error) {
	// 说明是提问，那么需要检查标题
	if req.GetParentGuid() == "" {
		if err := utils.CheckFileName(req.GetName()); err != nil {
			return nil, errcodev1.ErrInvalidArgument().WithMessage("name is invalid ,err: " + err.Error())
		}
	}
	// 权限检查
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_CREATE, req.GetUid(), req.GetSpaceGuid()); err != nil {
		return nil, err
	}
	// 说明是answer，那么需要检查他是否已经回答过答案
	if req.GetParentGuid() != "" {
		// 根据当前用户访问的question guid，查询他是否有回答这个question
		answerInfo, err := mysql.GetAnswerInfoByUid(invoker.Db.WithContext(ctx), req.GetParentGuid(), req.GetUid())
		if err != nil {
			return nil, err
		}
		// 说明已经回答过，不能在创建，请使用update接口
		if answerInfo.Guid != "" {
			return nil, errcodev1.ErrInvalidArgument().WithMessage("已经回答过问题了")
		}
	}

	// 如果不是提问，不需要检查标题
	file, err := service.File.CreateQuestion(ctx, req)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("create document fail, err: " + err.Error())
	}

	return &questionv1.CreateRes{File: file.ToFilePb()}, nil
}

// Update 修改问题
func (GrpcServer) Update(ctx context.Context, req *questionv1.UpdateReq) (*questionv1.UpdateRes, error) {
	// 如果不是提问，不需要检查标题
	// if err := utils.CheckFileName(req.GetName()); err != nil {
	//	return nil, errcodev1.ErrInvalidArgument().WithMessage("name is invalid ,err: " + err.Error())
	// }
	fileInfo, err := mysql.FileInfoByGuid(invoker.Db.WithContext(ctx), req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("Update question ,err: " + err.Error())
	}
	// 说明是提问，那么需要检查标题
	if fileInfo.ParentGuid == "" {
		if err := utils.CheckFileName(req.GetName()); err != nil {
			return nil, errcodev1.ErrInvalidArgument().WithMessage("name is invalid ,err: " + err.Error())
		}
	}

	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_UPDATE, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}
	// 更新db
	if err := service.File.UpdateQuestion(ctx, req); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdateDocument fail, err: " + err.Error())
	}

	return &questionv1.UpdateRes{}, nil
}

// MyInfo 获取这个用户相关问题的信息
func (s GrpcServer) MyInfo(ctx context.Context, req *questionv1.MyInfoReq) (*questionv1.MyInfoRes, error) {
	// 根据当前用户访问的question guid，查询他是否有回答这个question
	answerInfo, err := mysql.GetAnswerInfoByUid(invoker.Db.WithContext(ctx), req.GetGuid(), req.GetUid())
	if err != nil {
		return nil, err
	}
	return &questionv1.MyInfoRes{
		MyAnswerGuid: answerInfo.Guid,
	}, nil
}

//
// func (GrpcServer) Info(ctx context.Context, req *questionv1.InfoReq) (*questionv1.InfoRes, error) {
//	if _, err := mysql.FileInfoByGuidAndCmtGuid(invoker.Db.WithContext(ctx), req.GetCmtGuid(), req.GetGuid()); err != nil {
//		return nil, err
//	}
//
//	fileCacheInfo, err := service.File.Info(ctx, invoker.Db.WithContext(ctx), req.GetGuid())
//	if err != nil {
//		return nil, errcodev1.ErrNotFound().WithMessage("GetDocument fail" + err.Error())
//	}
//
//	// if fileCacheInfo.ParentGuid != "" {
//	//	respCollection, err := invoker.GrpcStat.IsCollection(ctx, &statv1.IsCollectionReq{
//	//		Uid:     req.GetUid(),
//	//		CmtGuid: req.GetCmtGuid(),
//	//		BizGuid: req.GetGuid(),
//	//		BizType: commonv1.CMN_BIZ_ANSWER,
//	//	})
//	//	if err != nil {
//	//		elog.Error("is collection fail", elog.FieldErr(err))
//	//	}
//	//	if respCollection != nil {
//	//		fileCacheInfo.IsCollect = respCollection.IsCollect
//	//	}
//	// } else {
//	//	respCollection, err := invoker.GrpcStat.IsCollection(ctx, &statv1.IsCollectionReq{
//	//		Uid:     req.GetUid(),
//	//		CmtGuid: req.GetCmtGuid(),
//	//		BizGuid: req.GetGuid(),
//	//		BizType: commonv1.CMN_BIZ_QUESTION,
//	//	})
//	//	if err != nil {
//	//		elog.Error("is collection fail", elog.FieldErr(err))
//	//	}
//	//	if respCollection != nil {
//	//		fileCacheInfo.IsCollect = respCollection.IsCollect
//	//	}
//	// }
//
//	str, err := service.File.GetContentByCreator(ctx, req.GetCmtGuid(), req.GetGuid())
//	if err != nil {
//		return nil, errcodev1.ErrDbError().WithMessage("get document content fail, err: " + err.Error())
//	}
//	fileCacheInfo.Content = string(str)
//
//	return &questionv1.InfoRes{
//		File: fileCacheInfo.ToQuestionPb(),
//	}, nil
// }
//
// // GetContentByCreator 创作者获取文档内容，无缓存，走机房流量
// func (GrpcServer) GetContentByCreator(ctx context.Context, req *questionv1.GetContentByCreatorReq) (*questionv1.GetContentByCreatorRes, error) {
//	str, err := service.File.GetContentByCreator(ctx, req.GetCmtGuid(), req.GetGuid())
//	if err != nil {
//		return nil, errcodev1.ErrDbError().WithMessage("get document content fail, err: " + err.Error())
//	}
//
//	return &questionv1.GetContentByCreatorRes{
//		Content: string(str),
//	}, nil
// }

// GetContent 获取文档内容的链接，前端读取文件内容渲染，有缓存，走CDN流量鉴权到OSS
// func (GrpcServer) GetContent(ctx context.Context, req *questionv1.GetContentReq) (*questionv1.GetContentRes, error) {
//	str, err := service.File.GetContentUrl(ctx, req.GetCmtGuid(), req.GetGuid())
//	if err != nil {
//		return nil, errcodev1.ErrDbError().WithMessage("get document content fail, err: " + err.Error())
//	}
//	err = mysql.FileUpdateExpr(invoker.Db.WithContext(ctx), req.GetGuid(), "cnt_view", 1)
//	if err != nil {
//		return nil, fmt.Errorf("create comment track total fail2, err: %w", err)
//	}
//	return &questionv1.GetContentRes{
//		Url: str,
//	}, nil
// }

// Delete 删除文章
func (GrpcServer) Delete(ctx context.Context, req *questionv1.DeleteReq) (*questionv1.DeleteRes, error) {
	// 检查权限
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_DELETE, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}
	// 查询fileInfo
	fileInfo, err := mysql.FileInfoMustExistsEerror(invoker.Db.WithContext(ctx), "id,space_guid,size", req.GetGuid())
	if err != nil {
		return nil, err
	}
	// 删除文章
	tx := invoker.Db.WithContext(ctx)
	err = service.File.Delete(ctx, tx, req.GetUid(), fileInfo.SpaceGuid, req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Delete fail, err: " + err.Error())
	}

	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("Delete fail, err: " + err.Error())
	}
	return &questionv1.DeleteRes{}, nil
}

// ListQuestion List Question
func (GrpcServer) ListQuestion(ctx context.Context, req *questionv1.ListQuestionReq) (*questionv1.ListQuestionRes, error) {
	list, err := service.File.QuestionListPage(ctx, invoker.Db.WithContext(ctx), req.GetSpaceGuid(), req.GetPagination(), req.GetSort())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("ListQuestion fail, err: " + err.Error())
	}

	return &questionv1.ListQuestionRes{
		List:       list.ToQuestionPb(),
		Pagination: req.GetPagination(),
	}, nil
}

// ListAnswer List Answer
func (GrpcServer) ListAnswer(ctx context.Context, req *questionv1.ListAnswerReq) (*questionv1.ListAnswerRes, error) {
	list, err := service.File.AnswerListPage(ctx, invoker.Db.WithContext(ctx), req.GetParentGuid(), req.GetPagination(), req.GetSort())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("ListQuestion fail, err: " + err.Error())
	}

	return &questionv1.ListAnswerRes{
		List:       list.ToAnswerPb(),
		Pagination: req.GetPagination(),
	}, nil
}

// PublicListByUserCreated PublicDocumentListByUserCreated 用户公开创建的文章列表
func (GrpcServer) PublicListByUserCreated(ctx context.Context, req *questionv1.PublicListByUserCreatedReq) (*questionv1.PublicListByUserCreatedRes, error) {
	files, err := service.File.PublicUserQAListPage(ctx, invoker.Db.WithContext(ctx), req.GetCreatedUid(), req.GetPagination())
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("PublicDocumentListByUserCreated PublicUserQAListPage fail, err: " + err.Error())
	}

	return &questionv1.PublicListByUserCreatedRes{
		List:       files.ToQAPb(),
		Pagination: req.GetPagination(),
	}, nil
}
