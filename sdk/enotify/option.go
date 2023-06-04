package enotify

import (
	"github.com/ego-component/egorm"
)

type Option func(c *Container)

// WithDb 设置 db
func WithDb(db *egorm.Component) Option {
	return func(c *Container) {
		c.db = db
	}
}
