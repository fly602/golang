SUB_DIRS = user/api \
		user/rpc \
		product/api \
		product/rpc \
		order/api \
		order/rpc \
		pay/api \
		pay/rpc
SUB_DIR_CLEAN := $(SUB_DIRS:%=%_clean)
	
all: gopreapre ${SUB_DIRS}

gopreapre:
	@go mod init
	@go mod tidy
	@go mod vendor

${SUB_DIRS}:ECHO
	@+make -C $@

ECHO:
	@echo ${SUB_DIRS}

clean:${SUB_DIR_CLEAN}
	#${foreach N, ${SUB_DIRS},make clean -C ${N};}
	rm -rf go.mod go.sum vendor

${SUB_DIR_CLEAN}:
	@+make clean -C ${@:%_clean=%}

docker-compose:
	@docker-compose build
	@docker-compose up -d

local:
	@docker build -t go1.18 -f ./golang/Dockerfile ./golang
	@docker run --rm -it --volume /home/uos/dde-go/src:/usr/src  go1.18 make -C /usr/src/go-community/docker/mall

help:
	@echo "make 需要在docker golang version环境下编译"
	@echo "make local 在当前本地环境下编译"
	@echo "make clean 删除中间目标文件和vendor文件"
	@echo "make docker-compose 通过docker-compose构建服务并运行"

.PHONY: all clean docker-compose local help

