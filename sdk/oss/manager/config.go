package manager

import (
	"github.com/gotomicro/ego/core/elog"
)

// Config ...
type Config struct {
	Mode    string
	Debug   bool
	Bucket  string
	CdnName string
	Prefix  string
	Alists
	File
	logger *elog.Component
}

type Alists struct {
	RegionId        string
	AccessKeyID     string
	AccessKeySecret string
	RoleArn         string
	Policy          string
}

type File struct {
	Path            string
	IsDeleteSrcPath bool
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Debug:   false,
		CdnName: "",
	}
}
