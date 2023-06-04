package article

import (
	"context"

	articlev1 "ecodepost/pb/article/v1"
	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"
)

func (GrpcServer) CreateDocument(ctx context.Context, req *articlev1.CreateDocumentReq) (*articlev1.CreateDocumentRes, error) {
	// if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_CREATE, req.CmtGuid, req.Uid, req.SpaceGuid); err != nil {
	//	return nil, err
	// }
	file, err := service.File.CreateDocument(ctx, req)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("create document fail, err: " + err.Error())
	}

	return &articlev1.CreateDocumentRes{File: file.ToFilePb()}, nil
}

func (GrpcServer) UpdateDocument(ctx context.Context, req *articlev1.UpdateDocumentReq) (*articlev1.UpdateDocumentRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_UPDATE, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}

	if err := service.File.UpdateDocument(ctx, req); err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdateDocument fail, err: " + err.Error())
	}
	return &articlev1.UpdateDocumentRes{}, nil
}

// DeleteDocument 删除文章
func (GrpcServer) DeleteDocument(ctx context.Context, req *articlev1.DeleteDocumentReq) (*articlev1.DeleteDocumentRes, error) {
	if err := pmspolicy.CheckEerror(ctx, commonv1.PMS_FILE_DELETE, req.GetUid(), req.GetGuid()); err != nil {
		return nil, err
	}
	if err := deleteDocument(ctx, req.Guid, req.Uid); err != nil {
		return nil, err
	}
	return &articlev1.DeleteDocumentRes{}, nil
}

func deleteDocument(ctx context.Context, guid string, uid int64) error {
	// 查询fileInfo
	fileInfo, err := mysql.FileInfoMustExistsEerror(invoker.Db.WithContext(ctx), "id,space_guid,size", guid)
	if err != nil {
		return err
	}
	// 删除file
	err = service.File.Delete(ctx, invoker.Db.WithContext(ctx), uid, fileInfo.SpaceGuid, guid)
	if err != nil {
		return errcodev1.ErrDbError().WithMessage("DeleteDocument fail, err: " + err.Error())
	}
	if err != nil {
		return errcodev1.ErrDbError().WithMessage("Delete fail, err: " + err.Error())
	}
	return nil
}
