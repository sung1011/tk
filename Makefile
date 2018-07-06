project_root:=$(CURDIR)

.PHONY: help build push

help:
	@echo "help"
	@echo "make build         -- 编译"
	@echo "make push          -- git更新"

build:
	sh ${project_root}/tools/build.sh

push:
	sh ${project_root}/tools/push.sh
