package alists

import (
	"fmt"
	"io"

	"ecodepost/sdk/oss/manager"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

type Component struct {
	stsClient *sts.Client
	config    *manager.Config
}

var (
	roleSessionName = "AliyunOSSTokenStsUser_uid_%v_tid_%v"
	// 	policy          = `{
	//   "Version": "1",
	//   "Statement": [
	//     {
	//       "Effect": "Allow",
	//       "Action": [
	//         "oss:Put*"
	//       ],
	//       "Resource": [
	//         "acs:oss:*:*:xxxx/*"
	//       ],
	//       "Condition": {}
	//     }
	//   ]
	// }`
)

// scheme defines fileDatasourceName
const scheme = "alists"

func init() {
	manager.Register(scheme, &Component{})
}

func (c *Component) Parse(config *manager.Config) (err error) {
	stsClient, err := sts.NewClientWithAccessKey(config.Alists.RegionId, config.Alists.AccessKeyID, config.Alists.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("sts new client fail, err: %w", err)
	}
	c.stsClient = stsClient
	c.config = config
	return nil
}

func (c *Component) GetToken(expire int, bucketName string, uid int64) (*manager.Credentials, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket or path is empty")
	}
	request := sts.CreateAssumeRoleRequest()
	if expire <= 0 || expire > 3600 {
		expire = 1800 // 默认1个小时
	}
	request.DurationSeconds = requests.NewInteger(expire)
	request.Scheme = "https"

	request.RoleArn = c.config.RoleArn
	// 角色扮演时指定RoleSessionName无效。此参数用来区分不同的Token，以标明谁在使用此Token，便于审计。格式：^[a-zA-Z0-9.@-_]+$，2-32个字符。了解更多信息请参见扮演角色操作接口 https://help.aliyun.com/document_detail/28763.html 。例如命名为a、1、abc*abc、忍者神龟等都是非法的。
	request.RoleSessionName = fmt.Sprintf(roleSessionName, uid, bucketName)
	request.Policy = c.config.Policy
	response, err := c.stsClient.AssumeRole(request)
	if err != nil {
		return nil, err
	}
	return &manager.Credentials{
		AccessKeySecret: response.Credentials.AccessKeySecret,
		Expiration:      response.Credentials.Expiration,
		AccessKeyId:     response.Credentials.AccessKeyId,
		SecurityToken:   response.Credentials.SecurityToken,
	}, nil
}

func (c *Component) PutObject(dstPath string, reader io.Reader) error {
	return nil
}

func (c *Component) GetObject(dstPath string) (output []byte, err error) {
	return nil, nil
}
