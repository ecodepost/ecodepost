package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"ecodepost/sdk/oss/manager"
	"github.com/ego-component/eredis"
)

type Component struct {
	cdnName string
	redis   *eredis.Component
	config  *manager.Config
}

// scheme defines fileDatasourceName
const scheme = "file"

func init() {
	manager.Register(scheme, &Component{})
}

func (c *Component) Parse(config *manager.Config) (err error) {
	c.config = config
	c.redis = eredis.Load("redis").Build()
	return nil
}

type Info struct {
	BucketName string `json:"bucketName"`
	Uid        int64  `json:"uid"`
}

func (c *Component) GetToken(expire int, bucketName string, uid int64) (*manager.Credentials, error) {
	//if expire <= 0 || expire > 3600 {
	//	expire = 1800 // 默认1个小时
	//}
	//info, err := json.Marshal(Info{
	//	BucketName: bucketName,
	//	Uid:        uid,
	//})
	//if err != nil {
	//	return nil, fmt.Errorf("json marshal fail, err: %w", err)
	//}
	//
	//token := uuid.New().String() + ":" + cast.ToString(uid)
	//err = c.redis.SetEX(context.Background(), "oss:file:token:"+token, info, time.Duration(expire)*time.Second)
	//if err != nil {
	//	return nil, fmt.Errorf("set ex fail, %w", err)
	//}
	return &manager.Credentials{
		//Expiration:    cast.ToString(expire),
		//SecurityToken: token,
	}, nil
}

func (c *Component) PutObject(dstPath string, reader io.Reader) error {
	// 创建目标目录
	dstPath = c.config.File.Path + "/" + dstPath
	err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if err != nil {
		return err
	}

	fileByte, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath, fileByte, os.ModePerm)
}

func (c *Component) PutObjectFromFile(dstPath, srcPath string) (err error) {
	// 创建目标目录
	dstPath = c.config.File.Path + "/" + dstPath
	err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if err != nil {
		return
	}
	var b []byte
	b, err = os.ReadFile(srcPath)
	if err != nil {
		return
	}

	err = os.WriteFile(dstPath, b, os.ModePerm)
	if err != nil {
		return
	}

	if c.config.File.IsDeleteSrcPath {
		err = os.Remove(srcPath)
	}
	return
}

func (c *Component) GetObject(dstPath string) (output []byte, err error) {
	return os.ReadFile(c.config.File.Path + "/" + dstPath)
}

func (c *Component) DeleteObject(dstPath string) (err error) {
	err = os.Remove(strings.TrimLeft(dstPath, "/"))
	return
}

//func (c *Component) DeleteObjects(dstPaths []string, options ...standard.Option) (output standard.DeleteObjectsResult, err error) {
//	for _, filePath := range dstPaths {
//		err1 := os.Remove(filePath)
//		if err1 != nil {
//			if err != nil {
//				err = errors.New(err.Error() + ", err is " + err1.Error())
//			} else {
//				err = err1
//			}
//		}
//	}
//	return
//}

func (c *Component) SignURL(dstPath string, method string, expiredInSec int64) (resp string, err error) {
	resp = c.cdnName + dstPath
	return
}
