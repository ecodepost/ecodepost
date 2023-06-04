package pmspolicy

import (
	"context"
	"fmt"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"

	commonv1 "ecodepost/pb/common/v1"
)

type FileCreateComment struct {
}

func init() {
	Register(NewActPolicy(&FileCreateComment{}))
}

func (*FileCreateComment) Scheme() commonv1.PMS_ACT {
	return commonv1.PMS_FILE_CREATE_COMMENT
}

// Check 检查
// 1 文件设置里是否关闭了评论
// 2 空间设置里是否允许评论
func (s *FileCreateComment) Check(ctx context.Context, uid int64, fileGuid string) (bool, error) {

	// 文件信息
	fileInfo, err := mysql.FileInfoByFieldMustExistsEerror(invoker.Db.WithContext(ctx), "id,close_comment_time,space_guid", fileGuid)
	if err != nil {
		return false, err
	}
	// 空间设置，优先级最高
	// spaceInfo, err := mysql.GetSpaceInfoByGuid(invoker.Db.WithContext(ctx), "`type`", fileInfo.SpaceGuid)
	// if err != nil {
	//	return false, errcodev1.ErrInternal().WithMessage("GetSpaceSetInfo fail, err: " + err.Error())
	// }

	// if spaceInfo.Type != commonv1.CMN_APP_ARTICLE {
	//	return false, fmt.Errorf("not support type")
	// }

	var isAllowCreateComment bool
	var spaceOption mysql.SpaceOption
	spaceOption, err = mysql.GetSpaceOptionInfo(invoker.Db.WithContext(ctx), fileInfo.SpaceGuid, commonv1.SPC_OPTION_FILE_IS_ALLOW_CREATE_COMMENT)
	if err != nil {
		return false, fmt.Errorf("SpaceCreateArticle Check fail, err: %w", err)
	}

	// 空间允许评论，并且文章没有关闭评论，那么用户才可以评论
	if spaceOption.OptionValue > 0 && fileInfo.CloseCommentTime == 0 {
		isAllowCreateComment = true
	}

	return isAllowCreateComment, nil
}
