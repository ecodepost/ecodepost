package column

import (
	"context"

	commonv1 "ecodepost/pb/common/v1"
	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service"
	"ecodepost/resource-svc/pkg/service/pmspolicy"
	"ecodepost/resource-svc/pkg/utils/slate"

	columnv1 "ecodepost/pb/column/v1"
	errcodev1 "ecodepost/pb/errcode/v1"

	"gorm.io/gorm"
)

func (GrpcServer) Create(ctx context.Context, req *columnv1.CreateReq) (*columnv1.CreateRes, error) {
	isPass, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_CREATE, req.GetUid(), req.GetSpaceGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	if !isPass {
		return nil, errcodev1.ErrAlreadyPermissionDenied().WithMessage("no auth")
	}

	file, err := service.File.CreateColumn(ctx, req)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("create document fail, err: " + err.Error())
	}
	return &columnv1.CreateRes{File: file.ToFilePb()}, nil
}

func (GrpcServer) Update(ctx context.Context, req *columnv1.UpdateReq) (*columnv1.UpdateRes, error) {
	isPass, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_UPDATE, req.GetUid(), req.GetFileGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}

	if !isPass {
		return nil, errcodev1.ErrAlreadyPermissionDenied().WithMessage("no auth")
	}

	// 新创建的file只有slate json模式
	// gocn的有富文本格式
	if req.FileFormat == commonv1.FILE_FORMAT_DOCUMENT_RICH {
		if req.Content != nil {
			slateJsonBytes, err := slate.HtmlToSlateJson(*req.Content)
			if err != nil {
				return nil, errcodev1.ErrInternal().WithMessage("html to slate json fail, err: " + err.Error())
			}
			req.Content = &slateJsonBytes
		}
		req.FileFormat = commonv1.FILE_FORMAT_DOCUMENT_SLATE
	}
	err = service.File.UpdateColumn(ctx, req)
	if err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("UpdateDocument fail, err: " + err.Error())
	}
	return &columnv1.UpdateRes{}, nil
}

// Delete 删除文章
func (GrpcServer) Delete(ctx context.Context, req *columnv1.DeleteReq) (*columnv1.DeleteRes, error) {
	isPass, err := pmspolicy.Check(ctx, commonv1.PMS_FILE_DELETE, req.GetUid(), req.GetGuid())
	if err != nil {
		return nil, errcodev1.ErrInternal().WithMessage("policy check, err:" + err.Error())
	}
	if !isPass {
		return nil, errcodev1.ErrAlreadyPermissionDenied().WithMessage("no auth")
	}
	if err := deleteDocument(ctx, req.Guid, req.Uid); err != nil {
		return nil, err
	}
	return &columnv1.DeleteRes{}, nil
}

func deleteDocument(ctx context.Context, guid string, uid int64) error {
	fileInfo, err := mysql.FileInfoMustExistsEerror(invoker.Db.WithContext(ctx), "id,space_guid,size", guid)
	if err != nil {
		return err
	}

	err = invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除file
		err = service.File.Delete(ctx, tx, uid, fileInfo.SpaceGuid, guid)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("DeleteDocument fail, err: " + err.Error())
		}
		// 删除edge和子file
		err = mysql.NodeDelete(tx, uid, guid)
		if err != nil {
			return errcodev1.ErrDbError().WithMessage("DeleteDocument edge fail, err: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
