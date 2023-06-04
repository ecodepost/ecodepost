package mysql

import (
	"fmt"

	commonv1 "ecodepost/pb/common/v1"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// FileEmoji 文章推荐
type FileEmoji struct {
	Id        int64  `json:"id" gorm:"not null;primary_key;auto_increment"`
	Guid      string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"`
	SpaceGuid string `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';index;comment:空间ID"`
	V1        int32  `gorm:"not null; default:0;"`
	V2        int32  `gorm:"not null; default:0;"`
	V3        int32  `gorm:"not null; default:0;"`
	V4        int32  `gorm:"not null; default:0;"`
	V5        int32  `gorm:"not null; default:0;"`
	V6        int32  `gorm:"not null; default:0;"`
	V7        int32  `gorm:"not null; default:0;"`
	V8        int32  `gorm:"not null; default:0;"`
	V9        int32  `gorm:"not null; default:0;"`
	V10       int32  `gorm:"not null; default:0;"`
	Ctime     int64  `gorm:"not null; default:0;"`
	Utime     int64  `gorm:"not null; default:0;"`
}

func (FileEmoji) TableName() string {
	return "file_emoji"
}

type FileEmojis []FileEmoji

var emojiList = []*commonv1.EmojiInfo{
	{
		Id:    1,
		Emoji: "👍",
	},
	{
		Id:    2,
		Emoji: "👎",
	},
	{
		Id:    3,
		Emoji: "😀",
	},
	{
		Id:    4,
		Emoji: "🎉",
	},
	{
		Id:    5,
		Emoji: "😕",
	},
	{
		Id:    6,
		Emoji: "❤️",
	},
	{
		Id:    7,
		Emoji: "🚀",
	}, {
		Id:    8,
		Emoji: "👀",
	},
}

var emojiMap = make(map[int32]*commonv1.EmojiInfo, 0)

func init() {
	for _, value := range emojiList {
		emojiMap[value.Id] = value
	}
}

func EmojiList() []*commonv1.EmojiInfo {
	return emojiList
}

func GetOneEmoji(id int32) *commonv1.EmojiInfo {
	return emojiMap[id]
}

func FileEmojiInfo(tx *gorm.DB, guid string) (resp FileEmoji, err error) {
	err = tx.Where("guid = ?", guid).Find(&resp).Error
	if err != nil {
		err = fmt.Errorf("file emoji info fail, err: %w", err)
		return
	}
	return
}

func FileEmojiList(tx *gorm.DB, guids []string) (resp FileEmojis, err error) {
	err = tx.Where("guid in (?)", guids).Find(&resp).Error
	if err != nil {
		err = fmt.Errorf("file emoji info fail, err: %w", err)
		return
	}
	return
}

// func EmojiMap(id int64) *filev1.EmojiInfo {
//	return emojiMap[id]
// }

// EmojiIncrease 增加
func EmojiIncrease(tx *gorm.DB, uid int64, guid string, emojiId int32) (err error) {
	if emojiId <= 0 || emojiId > 8 {
		err = fmt.Errorf("not exist v")
		return
	}

	mysqlFieldName := "v" + cast.ToString(emojiId)
	var emojiStat FileEmoji

	err = tx.Select("id,"+mysqlFieldName).Where("guid = ?", guid).Find(&emojiStat).Error
	if err != nil {
		return fmt.Errorf("EmojiIncrease CreateOrUpdate find failed, err: %w", err)
	}

	// 创建该数据
	if emojiStat.Id == 0 {
		emojiStat.Guid = guid
		err = tx.Create(&emojiStat).Error
		if err != nil {
			return fmt.Errorf("EmojiIncrease CreateOrUpdate create failed, err: %w", err)
		}
	}

	err = tx.Model(FileEmoji{}).Where("id = ?", emojiStat.Id).Updates(map[string]interface{}{
		mysqlFieldName: gorm.Expr(mysqlFieldName+"+?", 1),
	}).Error
	if err != nil {
		return fmt.Errorf("EmojiIncrease CreateOrUpdate update failed, err: %w", err)
	}

	err = FileEmojiStaticsCreate(tx, guid, uid, emojiId)
	if err != nil {
		return fmt.Errorf("EmojiIncrease FileEmojiStaticsCreate failed, err: %w", err)
	}
	return
}

// EmojiDecrease 增加
func EmojiDecrease(tx *gorm.DB, uid int64, guid string, emojiId int32) (err error) {
	if emojiId <= 0 || emojiId > 8 {
		err = fmt.Errorf("not exist v")
		return
	}

	mysqlFieldName := "v" + cast.ToString(emojiId)
	var emojiStat FileEmoji

	err = tx.Select("id,"+mysqlFieldName).Where("guid = ?", guid).Find(&emojiStat).Error
	if err != nil {
		return fmt.Errorf("EmojiDecrease CreateOrUpdate find failed, err: %w", err)
	}

	// 创建该数据
	if emojiStat.Id == 0 {
		emojiStat.Guid = guid
		err = tx.Create(&emojiStat).Error
		if err != nil {
			return fmt.Errorf("EmojiIncrease CreateOrUpdate create failed, err: %w", err)
		}
	}

	err = tx.Model(FileEmoji{}).Where("id = ?", emojiStat.Id).Updates(map[string]interface{}{
		mysqlFieldName: gorm.Expr(mysqlFieldName+"+?", -1),
	}).Error
	if err != nil {
		return fmt.Errorf("EmojiIncrease CreateOrUpdate update failed, err: %w", err)
	}

	err = FileEmojiStaticsDelete(tx, guid, uid, emojiId)
	if err != nil {
		return fmt.Errorf("EmojiIncrease FileEmojiStaticsCreate failed, err: %w", err)
	}
	return
}
