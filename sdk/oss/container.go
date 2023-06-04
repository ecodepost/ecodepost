package oss

import (
	"ecodepost/sdk/oss/manager"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

type Option func(c *Container)

// PackageName ..
const PackageName = "ecode.oss"

// Container ...
type Container struct {
	config *manager.Config
	logger *elog.Component
}

// Load ...
func Load(key string) *Container {
	var config = manager.DefaultConfig()
	if err := econf.UnmarshalKey(key, &config); err != nil {
		elog.Panic("parse config panic", elog.FieldErr(err), elog.FieldKey(key), elog.FieldValueAny(config))
	}
	obj := &Container{
		config: config,
		logger: elog.DefaultLogger.With(elog.FieldComponentName(PackageName)),
	}
	return obj
}

// Build ...
func (container *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(container)
	}

	obj, err := newComponent(container.config, container.logger)
	if err != nil {
		container.logger.Panic("new component err", elog.FieldErr(err), elog.FieldValueAny(container.config))
	}
	return obj
}
