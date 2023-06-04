package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"ecodepost/bff/pkg/invoker"
	"ecodepost/bff/pkg/server/bffcore"
	"ecodepost/bff/pkg/server/mw"

	ssov1 "ecodepost/pb/sso/v1"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/unittest/gintest"
)

var tt *gintest.Test

func init() {
	eapp.SetEgoDebug("true")
	invoker.InitByConf(conf)
	tt = gintest.Init(gintest.WithTestMiddleware(DebugToken(), mw.I18nCookie()))
}

func DebugToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u := &ssov1.User{
			Uid:  126,
			Name: "name_TomSawyer2",
		}
		ctx.Set(bffcore.ContextUserInfoKey, u)
		ctx.Next()
	}
}

var conf = `
mode = "local"
debug = true
domain = "http://of.yitum.com"
referralLink="%s/join-cmt?ref=%s#快来加入ossfarm吧！"
wcNotifyPath="/api/callbacks/wechat"

[k8s]
namespaces = ["of"]
addr = "https://60.205.217.184:6443"
debug = true
token = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImlVQXpLQTY5VTdlLWpyZjNkanRLZ2YyVUl3RGFQQlFIaDd3T0pyUlpzMU0ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJvZiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkZWZhdWx0LXRva2VuLTV3bDVyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRlZmF1bHQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJlNDQ5N2JiYi0yMmMxLTRkNjAtOGZkMC0wMGQwY2NmZDk5ODQiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6b2Y6ZGVmYXVsdCJ9.p8YYL8XDbdR75YVkiGuuWbHU_EUDf8jEk0ngoIlbwL1mhnGu14bhVHFFq3X5MZEy27sCY87GpPtuTgT6j5N0Sq2kP06-O7IAXhIoTxRvpvrUBXA-E6M-gRK6pWke_uqbI7FsQF-4qv8-qFFPTdqrQUm8gye8eCkOpKvDfkA46SZe8Ld-aTFz74tbQIAP6F6-Y7FjE7nddW7aUnvxnLQ3V2dfg4ROfNPco_t6cpYSoH3LHgvuX3lxlw2xYo-5qQ1vFRJ5-BYzk0x10CW7JSDKXl9h1rsPWQrLB3CLnVP8PRecKQrBSkd0cr6CH8FUsf5mF0SVE2hRvXob0PFRz6xwHg"

[server.http]
port = 9002
[server.governor]
port = 9003
[logger.default]
level = "debug"

[registry]
connectTimeout = "3s"
secure = false

[grpc.user]
addr = "k8s:///user-svc:9001"
#addr = "127.0.0.1:9091"
balancerName = "round_robin" # 默认值
dialTimeout = "3s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[grpc.track]
addr = "k8s:///track-svc:9001"
#addr = "127.0.0.1:9091"
balancerName = "round_robin" # 默认值
dialTimeout = "5s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[grpc.msg]
addr = "k8s:///msg-svc:9001"
#addr = "127.0.0.1:9091"
balancerName = "round_robin" # 默认值
dialTimeout = "3s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[grpc.resource]
addr = "k8s:///resource-svc:9001"
#addr = "127.0.0.1:9031"
balancerName = "round_robin" # 默认值
dialTimeout = "3s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[grpc.search]
addr = "k8s:///search-svc:9001"
#addr = "127.0.0.1:9091"
balancerName = "round_robin" # 默认值
dialTimeout = "5s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[grpc.shop]
addr = "k8s:///shop-svc:9001"
#addr = "127.0.0.1:9091"
balancerName = "round_robin" # 默认值
dialTimeout = "5s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[grpc.count]
addr = "k8s:///count-svc:9001"
#addr = "127.0.0.1:9091"
balancerName = "round_robin" # 默认值
dialTimeout = "5s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[sso]
ssoAddr = "k8s:///sso:9001"
tokenSecure = false # 关闭https写入
clientID = "ZLFh3tBLD2uGiGKW"
clientSecret = "ncENPEwvyt1ca4t8kexULQTGUxoPYR5i"
authURL = "http://of-sso.yitum.com/oauth/login"
redirectURL = "http://localhost:8000/api/oauth/code"
logoutURL = "http://of-sso.yitum.com/oauth/logout"

[grpc.audit]
addr = "k8s:///audit-svc:9001"
#addr = "127.0.0.1:9591"
balancerName = "round_robin" # 默认值
dialTimeout = "5s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[grpc.drive]
addr = "k8s:///drive-svc:9001"
#addr = "127.0.0.1:9091"
balancerName = "round_robin" # 默认值
dialTimeout = "5s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[grpc.app]
addr = "k8s:///app-svc:9001"
#addr = "127.0.0.1:9091"
balancerName = "round_robin" # 默认值
dialTimeout = "5s" # 默认值
readTimeout = "5s" # 默认值
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[mysql.user]
connMaxLifetime = "300s"
debug = true
dsn = "root:root@tcp(172.17.245.230:13306)/of_user?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&readTimeout=10s&timeout=3s&writeTimeout=3s"
level = "panic"
maxIdleConns = 50
maxOpenConns = 100
enableDetailSQL = true
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[mysql.gocn_res]
connMaxLifetime = "300s"
debug = true
dsn = "root:root@tcp(172.17.245.230:13306)/of_gocn_resource?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&readTimeout=10s&timeout=3s&writeTimeout=3s"
level = "panic"
maxIdleConns = 50
maxOpenConns = 100
enableDetailSQL = true
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[mysql.gocn_user]
connMaxLifetime = "300s"
debug = true
dsn = "root:root@tcp(172.17.245.230:13306)/of_gocn_user?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&readTimeout=10s&timeout=3s&writeTimeout=3s"
level = "panic"
maxIdleConns = 50
maxOpenConns = 100
enableDetailSQL = true
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[mysql.gocn_comment]
connMaxLifetime = "300s"
debug = true
dsn = "root:root@tcp(172.17.245.230:13306)/of_gocn_comment?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&readTimeout=10s&timeout=3s&writeTimeout=3s"
level = "panic"
maxIdleConns = 50
maxOpenConns = 100
enableDetailSQL = true
enableAccessInterceptor = true
enableAccessInterceptorRes = true
enableAccessInterceptorReq = true

[render]
webPage = "https://gocn-cdn.oss-cn-beijing.aliyuncs.com/saas-oa-dev/index.html"

[redis]
db = 5
addr = "127.0.0.1:6379"
password=""

[referral]
    [referral.email]
    tmplId = 1
    [referral.sms]
    tmplId = 2

[charge]
    [charge.wechat]
    key="dbmvwNCunI2CPoDlhEcGzjXN7kRBuvHO"

[kafka]
    debug = true
    brokers=["10.8.0.1:9092"]
    [kafka.client]
    timeout="3s"
    [kafka.producers.charge-svc-charge]        # 支付成功
    topic="charge-svc-charge"  # 指定生产消息的 topic
    [kafka.producers.charge-svc-refund]        # 支付退款
    topic="charge-svc-refund"  # 指定生产消息的 topic

[[custom.domain]]
host = "gocn.yitum.com"
cmtGuid = "bPMDkrxK5x"
`

func prettyJsonPrint(jsonStr []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, jsonStr, "", "\t")
	if error != nil {
		log.Println("prettyJsonPrint error: ", error)
		log.Println("origin str: " + string(jsonStr))
		return
	}

	fmt.Println(string(prettyJSON.Bytes()))
}
