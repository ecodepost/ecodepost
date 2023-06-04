package manager

import (
	"io"
)

var (
	registry map[string]Oss
)

func init() {
	registry = make(map[string]Oss)
}

// Register registers a dataSource creator function to the registry.
func Register(scheme string, creator Oss) {
	registry[scheme] = creator
}

func Get(scheme string) Oss {
	return registry[scheme]
}

// Oss ...
type Oss interface {
	Parse(config *Config) (err error)
	GetToken(expire int, bucketName string, uid int64) (*Credentials, error)
	PutObject(dstPath string, reader io.Reader) error
	GetObject(dstPath string) (output []byte, err error)
}

type Credentials struct {
	AccessKeySecret string `json:"accessKeySecret"`
	Expiration      string `json:"expiration"`
	AccessKeyId     string `json:"accessKeyId"`
	SecurityToken   string `json:"securityToken"`
}
