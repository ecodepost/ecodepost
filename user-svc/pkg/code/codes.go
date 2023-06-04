package code

import (
	errcodev1 "ecodepost/pb/errcode/v1"
)

var (
	ConfigNotExists    = errcodev1.ErrInvalidArgument().WithMessage("业务ID或业务类型不存在")
	OpTypeError        = errcodev1.ErrInvalidArgument().WithMessage("业务计数操作类型错误")
	OpTypeNotSupport   = errcodev1.ErrInvalidArgument().WithMessage("业务计数操作类型暂时不支持")
	IncrError          = errcodev1.ErrInvalidArgument().WithMessage("业务计数操作INCR错误")
	TidZero            = errcodev1.ErrInvalidArgument().WithMessage("业务计数查询目标ID长度为0")
	TidIsLimit         = errcodev1.ErrInvalidArgument().WithMessage("业务计数查询目标ID长度超过限制")
	FidIsLimit         = errcodev1.ErrInvalidArgument().WithMessage("业务计数查询来源ID长度超过限制")
	BidTypeParamsError = errcodev1.ErrInvalidArgument().WithMessage("业务计数查询业务来源参数错误")
	ResetError         = errcodev1.ErrInvalidArgument().WithMessage("业务计数操作RESET错误")
	BidTypeLengthError = errcodev1.ErrInvalidArgument().WithMessage("业务计数批量查询计数值类型长度错误")
	UpdateError        = errcodev1.ErrInvalidArgument().WithMessage("业务计数操作UPDATE错误")
	BatchParamsError   = errcodev1.ErrInvalidArgument().WithMessage("批量获取计数值长度错误")
	RateLimtError      = errcodev1.ErrInvalidArgument().WithMessage("超过内容计数频率限制")
	ServerBusyError    = errcodev1.ErrInvalidArgument().WithMessage("服务器繁忙")
	CountLimitError    = errcodev1.ErrInvalidArgument().WithMessage("超过最大上限")
	RedisOpError       = errcodev1.ErrInvalidArgument().WithMessage("服务端缓存操作异常")
	RedisUpdateExists  = errcodev1.ErrInvalidArgument().WithMessage("业务计数操作UPDATE更新KEY已存在")
)
