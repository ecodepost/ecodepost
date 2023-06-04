package enotify

import (
	commonv1 "ecodepost/pb/common/v1"
	"github.com/ego-component/egorm"
	"github.com/ego-component/eguid"
	"github.com/gotomicro/ego/core/elog"
)

type Container struct {
	name    string
	plugins map[commonv1.NOTIFY_CHANNEL]NotifyPlugin
	logger  *elog.Component
	db      *egorm.Component
	Guid    *eguid.Component
}

// DefaultContainer 构造默认容器
func DefaultContainer() *Container {
	return &Container{plugins: make(map[commonv1.NOTIFY_CHANNEL]NotifyPlugin)}
}

// Build 构建组件
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c)
}

var registry map[commonv1.NOTIFY_CHANNEL]NotifyPlugin

func init() {
	registry = make(map[commonv1.NOTIFY_CHANNEL]NotifyPlugin)
}

// Register registers a dataSource creator function to the registry.
func Register(scheme commonv1.NOTIFY_CHANNEL, creator NotifyPlugin) {
	registry[scheme] = creator
}
