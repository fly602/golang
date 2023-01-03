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

docker-compose:
	@docker-compose build
	@docker-compose up -d

.PHONY: all clean docker-compose

