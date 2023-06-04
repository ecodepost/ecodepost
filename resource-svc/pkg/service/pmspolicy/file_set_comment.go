package pmspolicy

import (
	"context"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
)

type FileSetComment struct {
}

func init() {
	Register(NewActPolicy(&FileSetComment{}))
}

func (*FileSetComment) Scheme() commonv1.PMS_ACT {
	return commonv1.PMS_FILE_SET_COMMENT
}

// Check 检查
// 1 是不是该文章的作者
func (s *FileSetComment) Check(ctx context.Context, uid int64, fileGuid string) (bool, error) {
	// 文件信息
	fileInfo, err := mysql.FileInfoByFieldMustExistsEerror(invoker.Db.WithContext(ctx), "id,created_by", fileGuid)
	if err != nil {
		return false, err
	}
	// 才能设置评论开关
	return fileInfo.CreatedBy == uid, nil
}
