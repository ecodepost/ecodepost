package mysql

import (
	"fmt"
	"time"

	commonv1 "ecodepost/pb/common/v1"

	"gorm.io/gorm"
)

// FileEmojiStatics 文章推荐
type FileEmojiStatics struct {
	Id    int64  `json:"id" gorm:"not null;primary_key;auto_increment"`
	Guid  string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index:guid_uid_v"`
	Uid   int64  `gorm:"not null; default:0;unique_index:guid_uid_v"`
	V     int32  `gorm:"not null; unique_index:guid_uid_v"`
	Ctime int64  `gorm:"not null; default:0;"`
	// EmojiStatics []*commonv1.EmojiInfo `gorm:"-"`
}

func (FileEmojiStatics) TableName() string {
	return "file_emoji_statics"
}

func (f FileEmojiStatics) ToPb() *commonv1.EmojiInfo {
	return &commonv1.EmojiInfo{
		Id:    f.V,
		Emoji: GetOneEmoji(f.V).GetEmoji(),
	}

}

func FileEmojiStaticsCreate(db *gorm.DB, guid string, uid int64, v int32) (err error) {
	if err = db.Create(&FileEmojiStatics{
		Guid:  guid,
		Uid:   uid,
		V:     v,
		Ctime: time.Now().Unix(),
	}).Error; err != nil {
		return fmt.Errorf("FileEmojiStaticsCreate failed,err: %w", err)
	}
	return
}

func FileEmojiStaticsDelete(db *gorm.DB, guid string, uid int64, v int32) (err error) {
	if err = db.Where("guid = ? and uid = ? and  v = ?", guid, uid, v).Delete(&FileEmojiStatics{}).Error; err != nil {
		return fmt.Errorf("FileEmojiStaticsDelete failed,err: %w", err)
	}
	return
}

func MyEmojiList(db *gorm.DB, guids []string, uid int64) (list []FileEmojiStatics, err error) {
	err = db.Select("guid,v").Where("guid in (?) and uid = ?", guids, uid).Find(&list).Error
	if err != nil {
		err = fmt.Errorf("MyEmojiList fail, err: %w", err)
		return
	}

	return
}
