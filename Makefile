APP_NAME:=ecodepost
APP_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
COMPILE_OUT:=$(APP_PATH)/../bin/$(APP_NAME)
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
	@cd $(APP_PATH) && go run main.go --config=config/local-zhengfke.toml --job=install

init:export EGO_DEBUG=true
init:
	@cd $(APP_PATH) && go run main.go --config=config/local-zhengfke.toml --job=init

build.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@chmod +x $(SCRIPT_PATH)/build/*.sh
	@cd $(APP_PATH) && $(SCRIPT_PATH)/build/gobuild.sh $(APP_NAME) $(COMPILE_OUT)
	@echo -e "\n"

docker.build:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"

	@#git clone https://github.com/ecodepost/ecodepost-fe.git
	@docker build -t ecodepost:latest .
	@rm -rf ecodepost-fe
	@echo -e "finish\n"


buf:
	@cd proto && buf generate