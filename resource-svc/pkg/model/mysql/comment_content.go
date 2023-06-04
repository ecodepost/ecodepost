package mysql

// CommentContent 评论内容 点查信息还是比较少的，所以分开放 TODO: 分表
type CommentContent struct {
	CommentGuid string `gorm:"type:char(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '';primary_key;comment:'主键id'"`
	Ctime       int64  `gorm:"comment:'创建时间'" json:"ctime"`
	Utime       int64  `gorm:"comment:'更新时间'" json:"utime"` // 更新时间
	Content     string `gorm:"not null;longtext" json:"content"`
	Status      int8   `gorm:"not null;comment:'0正常、1隐藏删除等'" json:"status"` // 商品价格
}

func (CommentContent) TableName() string {
	return "comment_content"
}
