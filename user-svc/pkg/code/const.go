package code

type (
	ActStatusType = uint8
)

const (
	// ActStatusNone 未操作
	ActStatusNone ActStatusType = 0
	// ActStatusAdded 已ADD, 比如已Follow
	ActStatusAdded ActStatusType = 1
	// ActStatusSubbed 已SUB, 比如已取消Follow
	ActStatusSubbed ActStatusType = 2
	// ActStatusAddCanceled 已取消Add
	// ActStatusAddCanceled ActStatusType = 3
	// ActStatusSubCanceled 已取消Sub
	// ActStatusSubCanceled ActStatusType = 4
)
