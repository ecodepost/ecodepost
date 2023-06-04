package invoker

import (
	"github.com/ego-component/egorm"
	"github.com/ego-component/eguid"
	"github.com/ego-component/eoauth2/server"
	"github.com/ego-component/eoauth2/storage/ssostorage"
	"github.com/ego-component/eredis"
)

var (
	Db             *egorm.Component
	Guid           *eguid.Component
	UserGuid       *eguid.Component
	TokenComponent *ssostorage.Component
	SsoServer      *server.Component
	Redis          *eredis.Component
)

func Init() error {
	Db = egorm.Load("mysql").Build()
	Guid = eguid.Load("user-svc.guid").Build()
	UserGuid = eguid.Load("user-svc.userGuid").Build()
	Redis = eredis.Load("redis").Build()
	TokenComponent = ssostorage.NewComponent(Db, Redis)
	SsoServer = server.Load("user-svc.oauth").Build(server.WithStorage(TokenComponent.GetStorage()))
	return nil
}
