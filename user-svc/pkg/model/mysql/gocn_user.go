package mysql

type GocnUser struct {
	Uid      int64
	Password string
}

func (GocnUser) TableName() string {
	return "user"
}
