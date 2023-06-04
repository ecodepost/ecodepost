package mysql

type SpaceColumn struct {
	Id        int64  `json:"id" gorm:"not null;primary_key;auto_increment"`
	SpaceGuid string `gorm:"type:char(12) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';unique_index;comment:唯一标识"`
	AuthorUid int64  `gorm:"not null;default:0;comment:作者UID"`
}

func (SpaceColumn) TableName() string {
	return "space_column"
}
