package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ecodepost/resource-svc/pkg/invoker"
	"ecodepost/resource-svc/pkg/model/mysql"
	"ecodepost/resource-svc/pkg/service/cache"
	"ecodepost/resource-svc/pkg/utils"
	"ecodepost/resource-svc/pkg/utils/slate"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"
	uploadv1 "ecodepost/pb/upload/v1"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreateOrCopyFileReq struct {
	Name       string
	ParentGuid string
	Uid        int64
	Content    string               // content或draftContent
	SpaceGuid  string               // 空间guid
	HeadImage  string               // 头图
	Node       commonv1.FILE_NODE   // 本节点类型
	ParentNode commonv1.FILE_NODE   // 父节点类型
	CreateTime int64                // 创建时间或者导入的创建数据时间
	UpdateTime int64                // 更新时间或者导入的更新数据时间
	CntView    int64                // 阅读总数
	FileFormat commonv1.FILE_FORMAT // 文件格式
	ContentKey string               // contentKey或draftContentKey, Copy文档需要的
	Size       int64                // 大小
	BizStatus  commonv1.FILE_BIZSTS // 业务状态
	Ip         string               // ip
}

// createFileCheckParamAndQuota 检查参数和版本配额
func createFileCheckParamAndQuota(ctx context.Context, req CreateOrCopyFileReq) error {
	if req.Uid == 0 {
		return errcodev1.ErrUidEmpty()
	}

	return nil
}

// GenContentKey 根据文件类型生成对象存储ContentKey
func GenContentKey(fileType commonv1.FILE_TYPE, fileStatus commonv1.FILE_BIZSTS) string {
	var prefixKey, dirDraft string
	switch fileStatus {
	case commonv1.FILE_BIZSTS_COURSE_DRAFT:
		dirDraft = "draft"
	case commonv1.FILE_BIZSTS_INVALID, commonv1.FILE_BIZSTS_COURSE_PUBLISHED:
		dirDraft = "publish"
	}
	switch fileType {
	case commonv1.FILE_TYPE_DOCUMENT:
		prefixKey = "article"
	case commonv1.FILE_TYPE_QUESTION:
		prefixKey = "question"
	case commonv1.FILE_TYPE_COLUMN:
		prefixKey = "column"
	default:
		// 其他类型不生成ContentKey
		return ""
	}

	return fmt.Sprintf("%s/%s/%s", prefixKey, dirDraft, strings.ReplaceAll(uuid.New().String(), "-", ""))
}

// GenUploadPath 生成oss path
func GenUploadPath(uploadType commonv1.CMN_UP_TYPE, spaceGuid, fileName string) (contentKey string, cdn string, bucket BucketName, err error) {
	// prefix
	//prefix := econf.GetString("alists.prefix") + "/"

	cfg := invoker.AliSts.GetConfig()

	var dir string
	switch uploadType {
	// 用户相关
	case commonv1.CMN_UP_TYPE_AVATAR:
		dir = cfg.Prefix + "/user/avatar"
	// 社区相关
	case commonv1.CMN_UP_TYPE_COMMUNITY:
		dir = cfg.Prefix + "/cmt/logo"
	case commonv1.CMN_UP_TYPE_COMMUNITY_BANNER:
		dir = cfg.Prefix + "/cmt/banner"
	// 活动相关
	case commonv1.CMN_UP_TYPE_FILE:
		if spaceGuid == "" {
			err = errcodev1.ErrSpaceEmpty()
			return
		}
		dir = cfg.Prefix + "/file/img/" + spaceGuid
	case commonv1.CMN_UP_TYPE_FILE_HEAD_IMAGE:
		if spaceGuid == "" {
			err = errcodev1.ErrSpaceEmpty()
			return
		}
		dir = cfg.Prefix + "/file/head-img/" + spaceGuid
	default:
		err = errcodev1.ErrInternal().WithMessage("GetToken upload type fail")
	}
	dir += "/" + time.Now().Format("20060102")
	contentKey = dir + "/" + fileName
	cdn = cfg.CdnName
	bucket = BucketName(cfg.Bucket)
	return
}

var fileTypeHasContent = []commonv1.FILE_TYPE{
	commonv1.FILE_TYPE_DOCUMENT,
	commonv1.FILE_TYPE_QUESTION,
	commonv1.FILE_TYPE_FILE,
	commonv1.FILE_TYPE_COLUMN,
}

// fileHasContent 通过fileType判断是否有content
func fileHasContent(fileType commonv1.FILE_TYPE) bool {
	return lo.Contains(fileTypeHasContent, fileType)
}

func (f *file) CreateFile(ctx context.Context, req CreateOrCopyFileReq, fileType commonv1.FILE_TYPE) (fileInfo *mysql.File, err error) {
	// 1. 检查参数和配额
	if err = createFileCheckParamAndQuota(ctx, req); err != nil {
		return nil, err
	}
	if req.CreateTime == 0 {
		return nil, fmt.Errorf("create time cant empty")
	}
	if !lo.Contains([]commonv1.FILE_FORMAT{commonv1.FILE_FORMAT_DOCUMENT_RICH, commonv1.FILE_FORMAT_DOCUMENT_SLATE}, req.FileFormat) {
		return nil, fmt.Errorf("file format error")
	}

	//rawContent := req.Content
	// 新创建的file只有slate json模式
	// gocn的有富文本格式
	if req.FileFormat == commonv1.FILE_FORMAT_DOCUMENT_RICH {
		slateJsonBytes, err := slate.HtmlToSlateJson(req.Content)
		if err != nil {
			return nil, errcodev1.ErrInternal().WithMessage("html to slate json fail, err: " + err.Error())
		}
		req.Content = slateJsonBytes
		req.FileFormat = commonv1.FILE_FORMAT_DOCUMENT_SLATE
	}

	// 2. 写入数据库
	fileGuid, err := Resource.GenerateGuid(ctx, commonv1.CMN_GUID_FILE, req.Uid)
	if err != nil {
		return nil, fmt.Errorf("file create gen guid fail, err: %w", err)
	}
	// 生成contentKey和draftContentKey
	fileInfo = &mysql.File{
		Guid:       fileGuid,
		SpaceGuid:  req.SpaceGuid,
		ParentGuid: req.ParentGuid,
		Name:       req.Name,
		HeadImage:  req.HeadImage,
		ContentKey: GenContentKey(fileType, req.BizStatus),
		CreatedBy:  req.Uid,
		UpdatedBy:  req.Uid,
		Status:     commonv1.FILE_STATUS_SUCC,
		BizStatus:  req.BizStatus,
		FileType:   fileType,
		FileNode:   req.Node,
		FileFormat: req.FileFormat,
		Ctime:      req.CreateTime,
		Utime:      req.UpdateTime,
		CntView:    req.CntView,
		Children:   nil,
		Ip:         req.Ip,
		Content:    req.Content,
	}
	fmt.Printf("fileInfo--------------->"+"%+v\n", fileInfo)

	if fileHasContent(fileType) && req.Content != "" {
		fileInfo.Hash = utils.MD5(req.Content)
	}
	if err = f.create(ctx, invoker.Db.WithContext(ctx), fileInfo); err != nil {
		return nil, fmt.Errorf("file create mysql create fail, err: %w", err)
	}

	// 3. 写入对象存储和缓存
	fc := fileInfo.ToCache()
	//if req.ContentKey == "" {
	//	// 只有依赖content的file才需要如下操作
	//	if fileHasContent(fileType) {
	//		// 如果为空, 说明是创建文件, 那么创建文档, 并更新缓存
	//		err = invoker.AliOss.PutObject(fileInfo.ContentKey, strings.NewReader(req.Content), &alioss.CBVar{
	//			SpaceGuid: req.SpaceGuid,
	//			FileGuid:  fileInfo.Guid,
	//			Uid:       req.Uid,
	//			FileType:  fileType,
	//		})
	//		if err != nil {
	//			return nil, fmt.Errorf("file create oss put fail, err: %w", err)
	//		}
	//		fc.WithContent(req.Content)
	//	}
	//} else {
	//	// 如果不为空，说明是复制文件
	//	if err = invoker.AliOss.CopyObject(req.ContentKey, fileInfo.ContentKey); err != nil {
	//		return nil, fmt.Errorf("file create oss put fail, err: %w", err)
	//	}
	//	output, err := invoker.AliOss.GetObject(fileInfo.ContentKey)
	//	if err != nil {
	//		elog.Error("CreateFile file cache get oss fail", elog.FieldErr(err))
	//	}
	//	// 加入cache
	//	fc.WithContent(string(output))
	//}

	fmt.Printf("fc--------------->"+"%+v\n", fc)
	if err = cache.FileCache.SetInfo(ctx, fc); err != nil {
		return nil, fmt.Errorf("file cache set info fail, err: %w", err)
	}
	return
}

type UpdateFileRequest struct {
	Uid        int64
	FileGuid   string
	NewContent *string
	HeadImage  *string
}

func (*file) UpdateFile(ctx context.Context, req UpdateFileRequest, updates map[string]any) (fileInfo mysql.File, err error) {
	//var oldFileInfo mysql.File
	_, err = mysql.FileInfoMustExistsEerror(invoker.Db.WithContext(ctx), "id,size,content_key", req.FileGuid)
	if err != nil {
		return
	}
	updates["utime"] = time.Now().Unix()
	updates["updated_by"] = req.Uid
	if req.NewContent != nil {
		updates["content"] = req.NewContent
	}

	if req.HeadImage != nil {
		updates["head_image"] = *req.HeadImage
	}

	// 事务更新File和Edge
	err = invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新file
		if req.NewContent != nil {
			updates["hash"] = utils.MD5(*req.NewContent)
		}
		if err = mysql.FileUpdate(tx, req.FileGuid, updates); err != nil {
			return errcodev1.ErrDbError().WithMessage("FileUpdate fail, err: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return
	}

	// 查询file
	if fileInfo, err = mysql.FileInfoByGuid(invoker.Db.WithContext(ctx), req.FileGuid); err != nil {
		err = errcodev1.ErrDbError().WithMessage("FileInfoByGuid fail, err: " + err.Error())
		return
	}

	// todo，卧槽这里有bug，如果用户没有更新文档，那么content缓存数据就没了
	// 然后有些optional，没传数据都会没有
	fc := fileInfo.ToCache()
	//if fileHasContent(fileInfo.FileType) {
	//	if req.NewContent != nil {
	//		err = invoker.AliOss.PutObject(oldFileInfo.ContentKey, strings.NewReader(*req.NewContent), &alioss.CBVar{
	//			SpaceGuid: fileInfo.SpaceGuid,
	//			FileGuid:  fileInfo.Guid,
	//			Uid:       req.Uid,
	//			FileType:  fileInfo.FileType,
	//		})
	//		if err != nil {
	//			err = errcodev1.ErrInternal().WithMessage("update file, err: " + err.Error())
	//			return
	//		}
	//		fc.WithContent(*req.NewContent)
	//
	//		if econf.GetBool("debug") {
	//			invoker.Db.WithContext(ctx).Model(mysql.FileDebug{}).Where("file_guid = ?", fileInfo.Guid).Update("content", *req.NewContent)
	//		}
	//	} else {
	//		var output []byte
	//		output, err = invoker.AliOss.GetObject(fileInfo.ContentKey)
	//		if err != nil {
	//			elog.Error("file cache get oss fail", elog.FieldErr(err))
	//			err = fmt.Errorf("GetObject output fail, %w", err)
	//			return
	//		}
	//		fc.WithContent(string(output))
	//	}
	//}
	if err = cache.FileCache.SetInfo(ctx, fc); err != nil {
		err = fmt.Errorf("file cache set info fail, err: %w", err)
		return
	}
	return
}

//func (f *file) GetSecretUploadCredential(ctx context.Context, uid int64) (*commonv1.UploadConf, error) {
//	privateBucket := invoker.AliOss.Config().PrivateBucket
//	credentials, err := invoker.AliSts.GetStsToken(900, invoker.AliOss.Config().PrivateBucket.Name, uid)
//	if err != nil {
//		return nil, errcodev1.ErrInternal().WithMessage("GetToken fail,err:" + err.Error())
//	}
//	ossCfg := invoker.AliOss.Config()
//	return &commonv1.UploadConf{
//		Region:          ossCfg.Region,
//		AccessKeyId:     credentials.AccessKeyId,
//		AccessKeySecret: credentials.AccessKeySecret,
//		StsToken:        credentials.SecurityToken,
//		Bucket:          privateBucket.Name,
//		CdnName:         privateBucket.CDN,
//		Expiration:      credentials.Expiration,
//	}, nil
//}

func (f *file) GetUploadPath(ctx context.Context, req *uploadv1.GetPathReq) (*uploadv1.GetPathRes, error) {
	// 生成path
	suffix := utils.GetFileSuffix(req.GetFileName())
	fileName := strings.ReplaceAll(uuid.New().String(), "-", "") + "." + suffix
	contentKey, cdnName, bucketName, err := GenUploadPath(req.UploadType, req.GetSpaceGuid(), fileName)
	if err != nil {
		return nil, errcodev1.ErrInvalidArgument().WithMessage("GenUploadPath fail, err:" + err.Error())
	}
	fullName := string(bucketName) + "/" + contentKey
	nowTime := time.Now().Unix()
	// 更新file表
	err = invoker.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := mysql.FileCreate(tx, &mysql.File{
			Ctime:      nowTime,
			Utime:      nowTime,
			CreatedBy:  req.Uid,
			Name:       fullName,
			ContentKey: contentKey,
			SpaceGuid:  req.SpaceGuid,
			FileType:   commonv1.FILE_TYPE_FILE,
		}); err != nil {
			return errcodev1.ErrDbError().WithMessage("FileCreate fail,err:" + err.Error())
		}

		// 如果是图片，则更新image表 TODO,判断更多的图片类型
		if lo.Contains([]string{"jpg", "jpeg", "png", "webp"}, suffix) {
			if err := mysql.AddImage(invoker.Db.WithContext(ctx), &mysql.Image{
				Ctime:   nowTime,
				Utime:   nowTime,
				Uid:     req.GetUid(),
				Name:    fullName,
				Type:    suffix,
				CdnName: cdnName,
				Path:    contentKey,
			}); err != nil {
				return errcodev1.ErrDbError().WithMessage("AddImage fail,err:" + err.Error())
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &uploadv1.GetPathRes{
		Bucket:  string(bucketName),
		Path:    contentKey,
		CdnName: cdnName,
	}, nil
}
