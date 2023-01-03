# 获取用户ID
UID=$(shell id -u)
# 设置golang 容器环境变量
GO_DOCKER_TZ=Asia/Shanghai
GO_DOCKER_LANG=C.UTF-8
GO_DOCKER_IMAGE=go1.18 
GO_DOCKER_NAME=go-build1.18
GO_DOCKER_VOLUME_PATH=/usr/src
GO_DOCKER_SRC_PATH=go-community/docker/mall
GO_LOCAL_VOLUME_PATH=$(HOME)/dde-go/src
GO_DOCKER_BUILD=make build

SUB_DIRS = user/api \
		user/rpc \
		product/api \
		product/rpc \
		order/api \
		order/rpc \
		pay/api \
		pay/rpc

SUB_DIR_CLEAN := $(SUB_DIRS:%=%_clean)

# 使用docker go1.18环境进行编译
target:
	@docker build -t ${GO_DOCKER_IMAGE} -f ./golang/Dockerfile ./golang
	@docker run --rm -it -u $(UID) \
		--volume $(GO_LOCAL_VOLUME_PATH):$(GO_DOCKER_VOLUME_PATH) \
		--name $(GO_DOCKER_NAME) \
		--env TZ=$(GO_DOCKER_TZ) \
		--env LANG=$(GO_DOCKER_LANG) \
		$(GO_DOCKER_IMAGE) $(GO_DOCKER_BUILD) -C $(GO_DOCKER_VOLUME_PATH)/$(GO_DOCKER_SRC_PATH)

all: target run

build: gopreapre ${SUB_DIRS}

gopreapre:
	@go mod tidy
	@go mod vendor

${SUB_DIRS}:ECHO
	@+make -C $@

ECHO:
	@echo ${SUB_DIRS}

clean:${SUB_DIR_CLEAN}
	#${foreach N, ${SUB_DIRS},make clean -C ${N};}

${SUB_DIR_CLEAN}:
	@+make clean -C ${@:%_clean=%}

run:
	@docker-compose build
	@docker-compose up -d

help:
	@echo "make 使用go1.18docker环境编译微服务"
	@echo "make build 编译微服务,需要环境go.18"
	@echo "make clean 删除中间目标文件"
	@echo "make run 通过docker-compose构建服务并在容器中运行"
	@echo "make all "编译并在容器中运行微服务""


.PHONY: all build clean run

