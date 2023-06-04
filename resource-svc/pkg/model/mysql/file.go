package mysql

import (
	"fmt"
	"time"

	commonv1 "ecodepost/pb/common/v1"
	errcodev1 "ecodepost/pb/errcode/v1"

	"gorm.io/gorm"
)

type File struct {
	Id               int64                `json:"id" gorm:"not null;primary_key;auto_increment"`
	Guid             string               `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"`
	SpaceGuid        string               `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:空间GUID"`
	ParentGuid       string               `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"` // 只用于问答
	Name             string               `gorm:"not null;type:varchar(160)"`
	Size             int64                `gorm:"not null; default:0; comment:文件大小"`
	HeadImage        string               `gorm:"not null;default:'';size:191;comment:头图"`
	ContentKey       string               `gorm:"not null; default:''; size:191; comment:云文件存储key"`
	CreatedBy        int64                `gorm:"not null; default:0; comment:创建人"`
	UpdatedBy        int64                `gorm:"not null; default:0; comment:更新人"`
	DeletedBy        int64                `gorm:"not null; default:0; comment:删除人"`
	PublishedBy      int64                `gorm:"not null; default:0; comment:发布人"`
	SuggestedBy      int64                `gorm:"not null; default:0; comment:推荐的人"`
	Status           commonv1.FILE_STATUS `gorm:"type:smallint; not null; default:0; comment:文件状态"`
	BizStatus        commonv1.FILE_BIZSTS `gorm:"type:smallint; not null; default:0; comment:业务状态"`
	FileType         commonv1.FILE_TYPE   `gorm:"type:smallint; not null; default:0; index; comment:文件File类型"`
	FileNode         commonv1.FILE_NODE   `gorm:"type:smallint; not null; default:0; comment:文件Node类型"`
	FileFormat       commonv1.FILE_FORMAT `gorm:"type:smallint; not null; default:0; comment:文件格式"`
	SpaceTopTime     int64                `gorm:"not null; default:0; comment:置顶时间"`
	SuggestedTime    int64                `gorm:"not null; default:0; comment:推荐时间"`
	PublishTime      int64                `gorm:"not null; default:0; comment:发布时间"`
	CntComment       int64                `gorm:"not null; default:0; comment:评论数"`
	CntView          int64                `gorm:"not null; default:0; comment:阅读数"`
	CntCollect       int64                `gorm:"not null; default:0; comment:收藏数"`
	CntLike          int64                `gorm:"not null; default:0; comment:喜欢数"`
	CloseCommentTime int64                `gorm:"not null; default:0;comment:关闭评论时间"`
	Ctime            int64                `gorm:"not null; default:0; comment:创建时间"`
	Utime            int64                `gorm:"not null; default:0; comment:更新时间"`
	Dtime            int64                `gorm:"not null; default:0; comment:彻底删除时间"` // 彻底删除时间
	HiddenTime       int64                `gorm:"not null; default:0; comment:隐藏时间"`
	HotScore         float64              `gorm:"not null; default:0; comment:热度"`
	RecommendScore   float64              `gorm:"not null; default:0; comment:热度"`
	Hash             string               `gorm:"not null; default:''; comment:content的hash值"`
	Ip               string               `gorm:"type:varchar(15); not null;default:'';comment:IP"` // 后续可以考虑用unsigned int存储
	Content          string               `gorm:"type:longtext; not null;comment:content"`
	Children         Files                `gorm:"-"`
	Sort             int64                `gorm:"-"`
}

func (File) TableName() string {
	return "file"
}

type FileContent struct {
	Guid            string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"`
	ContentKey      string `gorm:"content_key; not null; default:''; size:191; comment:云文件存储key"`
	ContentDraftKey string `gorm:"content_key; not null; default:''; size:191; comment:云文件存储key"`
	Hash            string `gorm:"hash; not null; default:''; size:191; comment:content的hash值"`
}

type Files []*File

type FileGuid struct {
	Guid       string
	CntComment int64
	CntCollect int64
	CntView    int64
	Size       int64
	FileType   commonv1.FILE_TYPE
}

type FileGuids []FileGuid

func (list Files) FindByGuid(guid string) *File {
	return list.Find(func(e *File) bool {
		return e.Guid == guid
	})
}

func (list Files) Find(fn func(e *File) bool) *File {
	for _, spaceInfo := range list {
		if fn(spaceInfo) {
			return spaceInfo
		}
	}

	return nil
}

func (f *File) ToFilePb() *commonv1.FileInfo {
	return &commonv1.FileInfo{
		Guid: f.Guid,
		Name: f.Name,
	}
}

func (f File) ToCache() *FileCache {
	var isAllowCreateComment, isSiteTop, isRecommend int32
	if f.CloseCommentTime == 0 {
		isAllowCreateComment = 1
	}
	if f.SpaceTopTime > 0 {
		isSiteTop = 1
	}
	if f.SuggestedTime > 0 {
		isRecommend = 1
	}
	return &FileCache{
		Guid:                 f.Guid,
		SpaceGuid:            f.SpaceGuid,
		ParentGuid:           f.ParentGuid,
		Name:                 f.Name,
		HeadImage:            f.HeadImage,
		Hash:                 f.Hash,
		CreatedUid:           f.CreatedBy,
		UpdatedUid:           f.UpdatedBy,
		Ctime:                f.Ctime,
		Utime:                f.Utime,
		IsAllowCreateComment: isAllowCreateComment,
		IsSiteTop:            isSiteTop,
		IsRecommend:          isRecommend,
		FileFormat:           int32(f.FileFormat),
		Size:                 f.Size,
		Content:              f.Content,
		BizStatus:            int32(f.BizStatus),
		FileType:             int32(f.FileType),
	}
}

// FileCreate 创建一条记录
func FileCreate(db *gorm.DB, data *File) (err error) {
	if err = db.Create(data).Error; err != nil {
		return fmt.Errorf("FileCreate failed,err: %w", err)
	}
	return
}

// FileInfoByGuid Info的扩展方法，根据Cond查询单条记录
func FileInfoByGuid(db *gorm.DB, guid string) (resp File, err error) {
	if err = db.Model(File{}).Where("guid = ? and dtime = 0", guid).Find(&resp).Error; err != nil {
		return File{}, fmt.Errorf("FileInfoByGuid fail, err: %w", err)
	}
	return
}

// FileInfoByFieldMustExistsEerror 根据guid查询指定字段，并返回ego.Error
func FileInfoByFieldMustExistsEerror(db *gorm.DB, field, guid string) (info File, err error) {
	e := db.Select(field).Where("guid = ?", guid).Find(&info).Error
	if e != nil {
		err = errcodev1.ErrDbError().WithMessage("FileInfoByField fail, " + err.Error())
		return
	}
	if info.Id == 0 {
		err = errcodev1.ErrInternal().WithMessage("File not exists")
		return
	}
	return
}

// FileInfoByFieldAndGuid Info的扩展方法，根据Cond查询单条记录
func FileInfoByFieldAndGuid(db *gorm.DB, field, guid string) (resp File, err error) {
	if err = db.Select(field).Where("guid = ? and dtime = 0", guid).Find(&resp).Error; err != nil {
		return File{}, fmt.Errorf("FileInfoByFieldAndGuid failed,err: %w", err)
	}
	return
}

// FileInfoMustExistsEerror 查询文件，如果文件不存在或DB报错，抛出不同的eerrors.Error
func FileInfoMustExistsEerror(db *gorm.DB, field, guid string) (resp File, err error) {
	fileInfo, e := FileInfoByFieldAndGuid(db, field, guid)
	if e != nil {
		err = errcodev1.ErrDbError().WithMessage("FileInfoByFieldAndGuid fail, " + err.Error())
		return
	}
	if fileInfo.Id == 0 {
		err = errcodev1.ErrInternal().WithMessage("File not exists")
		return
	}
	return fileInfo, nil
}

// FileListByField 获取文档列表
func FileListByField(db *gorm.DB, field string, guids []string) (list Files, err error) {
	if err = db.Select(field).Where("guid in (?)", guids).Find(&list).Error; err != nil {
		err = fmt.Errorf("IsFileAuthor fail, err: %w", err)
		return
	}
	return
}

// GetFileGuidEerror 获取文档，报错则返回ego.Error
func GetFileGuidEerror(db *gorm.DB, guid string) (*FileGuid, error) {
	var fileInfo FileGuid
	if err := db.Select("guid,cnt_comment,cnt_view,cnt_collect").Model(File{}).Where("guid = ? and dtime = 0", guid).Find(&fileInfo).Error; err != nil {
		return nil, errcodev1.ErrDbError().WithMessage("FileInfoByGuidAndCmtGuid fail1, err: " + err.Error())
	}
	if fileInfo.Guid == "" {
		return nil, errcodev1.ErrNotFound().WithMessage("GetDocument fail2")
	}
	return &fileInfo, nil
}

// GetAnswerInfoByUid 查看某个用户是否有回答过答案
func GetAnswerInfoByUid(db *gorm.DB, fileGuid string, uid int64) (resp File, err error) {
	if err = db.Select("guid").Where("parent_guid = ?  and created_by = ?  and dtime = 0", fileGuid, uid).Find(&resp).Error; err != nil {
		return File{}, errcodev1.ErrDbError().WithMessage("GetAnswerInfoByUid fail, err: " + err.Error())
	}
	return
}

// FileSpaceGuidByGuid Info的扩展方法，根据Cond查询单条记录
func FileSpaceGuidByGuid(db *gorm.DB, guid string) (spaceGuid string, err error) {
	var resp File
	if err = db.Select("space_guid").Where("guid = ?", guid).Find(&resp).Error; err != nil {
		return "", fmt.Errorf("FileSpaceGuidByGuid failed,err: %w", err)
	}
	spaceGuid = resp.SpaceGuid
	return
}

// FileContentInfoByGuid Info的扩展方法，根据Cond查询单条记录
func FileContentInfoByGuid(db *gorm.DB, guid string) (resp FileContent, err error) {
	if err = db.Select("guid,content_key,hash").Model(File{}).Where("guid = ?", guid).Find(&resp).Error; err != nil {
		return FileContent{}, fmt.Errorf("FileContentInfoByGuid failed,err: %w", err)
	}
	return
}

// FileUpdate 根据主键更新一条记录
func FileUpdate(db *gorm.DB, guid string, ups map[string]any) (err error) {
	if err = db.Model(File{}).Where("guid = ?", guid).Updates(ups).Error; err != nil {
		return fmt.Errorf("FileUpdate failed,err: %w", err)
	}
	return
}

// FileByGuid 根据主键查询一条记录
func FileByGuid(db *gorm.DB, guid string) (file File, err error) {
	if err = db.Where("guid = ? and dtime = 0", guid).Find(&file).Error; err != nil {
		err = fmt.Errorf("FileByGuid failed,err: %w", err)
		return
	}
	return
}

// FileTitleMapByGuids 根据主键更新一条记录
func FileTitleMapByGuids(db *gorm.DB, guids []string) (output map[string]File, err error) {
	output = make(map[string]File)
	var list Files
	if err = db.Select("name,guid,biz_status").Where("guid in (?) and dtime = 0", guids).Find(&list).Error; err != nil {
		err = fmt.Errorf("FileTitleMapByGuids failed,err: %w", err)
		return
	}
	for _, value := range list {
		output[value.Guid] = *value
	}
	return
}

func SetFileSpaceTop(db *gorm.DB, uid int64, guid string, spaceGuid string) (err error) {
	err = db.Model(File{}).Where("guid = ? ", guid).Updates(map[string]any{
		"space_top_time": time.Now().Unix(),
		"utime":          time.Now().Unix(),
		"updated_by":     uid,
	}).Error
	if err != nil {
		err = fmt.Errorf("SetFileSpaceTop fail, err: %w", err)
		return
	}

	err = FileSpaceTopCreate(db, guid, spaceGuid, uid)
	if err != nil {
		err = fmt.Errorf("FileSpaceTopCreate fail, err: %w", err)
		return
	}
	return
}

func CancelFileSpaceTop(db *gorm.DB, uid int64, guid string) (err error) {
	err = db.Model(File{}).Where("guid = ?", guid).Updates(map[string]any{
		"space_top_time": 0,
		"utime":          time.Now().Unix(),
		"updated_by":     uid,
	}).Error
	if err != nil {
		err = fmt.Errorf("CancelFileSpaceTop fail1, err: %w", err)
		return
	}

	err = FileSpaceTopDelete(db, guid)
	if err != nil {
		err = fmt.Errorf("CancelFileSpaceTop fail2, err: %w", err)
		return
	}
	return
}

func FileSuggest(db *gorm.DB, uid int64, guid string) (err error) {
	err = db.Model(File{}).Where("guid = ?", guid).Updates(map[string]any{
		"suggested_time": time.Now().Unix(),
		"utime":          time.Now().Unix(),
		"updated_by":     uid,
		"suggested_by":   uid,
	}).Error
	if err != nil {
		err = fmt.Errorf("FileSuggest fail, err: %w", err)
		return
	}
	spaceGuid, err := FileSpaceGuidByGuid(db, guid)
	if err != nil {
		err = fmt.Errorf("FileSpaceGuidByGuid fail, err: %w", err)
		return
	}

	err = FileRecommendCreate(db, guid, spaceGuid, uid)
	if err != nil {
		err = fmt.Errorf("FileRecommendCreate fail, err: %w", err)
		return
	}
	return
}

func FileCancelSuggest(db *gorm.DB, uid int64, guid string) (err error) {
	err = db.Model(File{}).Where("guid = ? ", guid).Updates(map[string]any{
		"suggested_time": 0,
		"utime":          time.Now().Unix(),
		"suggested_by":   uid,
		"updated_by":     uid,
	}).Error
	if err != nil {
		err = fmt.Errorf("FileSuggest fail, err: %w", err)
		return
	}

	err = FileRecommendDelete(db, guid)
	if err != nil {
		err = fmt.Errorf("FileRecommendDelete fail, err: %w", err)
		return
	}
	return
}

func FileOpenComment(db *gorm.DB, uid int64, guid string) (err error) {
	err = db.Model(File{}).Where("guid = ? ", guid).Updates(map[string]any{
		"close_comment_time": 0,
		"utime":              time.Now().Unix(),
		"updated_by":         uid,
	}).Error
	if err != nil {
		err = fmt.Errorf("FileOpenComment fail, err: %w", err)
		return
	}
	return
}

func FileCloseComment(db *gorm.DB, uid int64, guid string) (err error) {
	err = db.Model(File{}).Where("guid = ?", guid).Updates(map[string]any{
		"close_comment_time": time.Now().Unix(),
		"utime":              time.Now().Unix(),
		"updated_by":         uid,
	}).Error
	if err != nil {
		err = fmt.Errorf("FileCloseComment fail, err: %w", err)
		return
	}
	return
}

// FileRemove 文件放到回收站
func FileRemove(db *gorm.DB, uid int64, spaceGuid string, guid string) (err error) {
	err = db.Model(File{}).Where("guid = ? and space_guid = ? ", guid, spaceGuid).Updates(map[string]any{
		"dtime":      time.Now().Unix(),
		"utime":      time.Now().Unix(),
		"updated_by": uid,
		"deleted_by": uid,
	}).Error
	if err != nil {
		err = fmt.Errorf("FileRemove fail, err: %w", err)
		return
	}
	// 去除其他的推荐操作
	err = FileRecommendDelete(db, guid)
	if err != nil {
		err = fmt.Errorf("FileRemove FileRecommendDelete , err: %w", err)
		return
	}
	// 去除其他的置顶操作
	err = FileSpaceTopDelete(db, guid)
	if err != nil {
		err = fmt.Errorf("FileRemove FileSpaceTopDelete , err: %w", err)
		return
	}
	return
}

func FileUpdateExpr(db *gorm.DB, guid string, column string, delta int) (err error) {
	err = db.Model(File{}).Where("guid = ? ", guid).
		Updates(map[string]any{
			column:  gorm.Expr(column+"+ ?", delta),
			"utime": time.Now().Unix(),
		}).Error
	if err != nil {
		return
	}
	return
}

func (m Files) ToTree() Files {
	mTreeMap := make(map[string]*File)
	for _, item := range m {
		mTreeMap[item.Guid] = item
	}

	list := make(Files, 0)
	for _, item := range m {
		// 筛选出父级节点
		if item.ParentGuid == "" {
			list = append(list, item)
			continue
		}

		if pItem, ok := mTreeMap[item.ParentGuid]; ok {
			if pItem.Children == nil {
				children := Files{item}
				pItem.Children = children
				continue
			}
			pItem.Children = append(pItem.Children, item)
		}
	}
	return list
}
