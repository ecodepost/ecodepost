package service

var (
	Logger    *logger
	Authorize *authorize
)

func Init() error {
	Authorize = InitAuthorize()
	Logger = &logger{}
	return nil
}
