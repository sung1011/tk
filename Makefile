project_root:=$(CURDIR)

.PHONY: help build

help:
	@echo "help"
	@echo "make build         -- 编译"

build:
	sh ${project_root}/scripts/build.sh

