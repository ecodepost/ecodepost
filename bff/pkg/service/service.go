package service

var (
	Cmt *cmt
)

func Init() error {
	Cmt = &cmt{}
	return nil
}
