APP_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SCRIPT_PATH:=$(APP_PATH)/scripts

run_user:export EGO_DEBUG=true
run_user:
	@cd $(APP_PATH) && go run user-svc/main.go --config=config/local.toml

run_resource:export EGO_DEBUG=true
run_resource:
	@cd $(APP_PATH) && go run resource-svc/main.go --config=config/local.toml

run_bff:export EGO_DEBUG=true
run_bff:
	@cd $(APP_PATH) && go run bff/main.go --config=config/local.toml

run_all:export EGO_DEBUG=true
run_all:
	@cd $(APP_PATH) && rm -f user-svc.sock
	@cd $(APP_PATH) && rm -f resource-svc.sock
	@cd $(APP_PATH) && go run main.go --config=config/local.toml

run_all_askuy:export EGO_DEBUG=true
run_all_askuy:
	@cd $(APP_PATH) && rm -f user-svc.sock
	@cd $(APP_PATH) && rm -f resource-svc.sock
	@cd $(APP_PATH) && go run main.go --config=config/local-askuy.toml

run_all_zhengfke:export EGO_DEBUG=true
run_all_zhengfke:
	@cd $(APP_PATH) && rm -f user-svc.sock
	@cd $(APP_PATH) && rm -f resource-svc.sock
	@cd $(APP_PATH) && go run main.go --config=config/local-zhengfke.toml

install:export EGO_DEBUG=true
install:
	@cd $(APP_PATH) && go run main.go --config=config/local.toml --job=install

init:export EGO_DEBUG=true
init:
	@cd $(APP_PATH) && go run main.go --config=config/local.toml --job=init


buf:
	@cd proto && buf generate