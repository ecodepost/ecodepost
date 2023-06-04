package oss

import (
	"io"

	_ "ecodepost/sdk/oss/alists"
	_ "ecodepost/sdk/oss/file"
	"ecodepost/sdk/oss/manager"
	"github.com/gotomicro/ego/core/elog"
)

type Component struct {
	client manager.Oss
	config *manager.Config
}

func newComponent(config *manager.Config, logger *elog.Component) (client *Component, err error) {
	comp := manager.Get(config.Mode)
	if comp == nil {
		logger.Panic("not exist mode", elog.String("mode", config.Mode))
	}

	err = comp.Parse(config)
	if err != nil {
		logger.Panic("new component err", elog.FieldErr(err), elog.FieldValueAny(config))
	}

	client = &Component{
		client: comp,
		config: config,
	}
	return
}

func (c *Component) GetToken(expire int, bucketName string, uid int64) (*manager.Credentials, error) {
	return c.client.GetToken(expire, bucketName, uid)
}

func (c *Component) PutObject(dstPath string, reader io.Reader) error {
	return c.client.PutObject(dstPath, reader)
}

func (c *Component) GetObject(dstPath string) (output []byte, err error) {
	return c.client.GetObject(dstPath)
}

func (c *Component) GetConfig() *manager.Config {
	return c.config
}
