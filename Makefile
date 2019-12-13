# Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
# iSulad-lxcfs-toolkit is licensed under the Mulan PSL v1.
# You can use this software according to the terms and conditions of the Mulan PSL v1.
# You may obtain a copy of Mulan PSL v1 at:
#     http://license.coscl.org.cn/MulanPSL
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
# PURPOSE.
# See the Mulan PSL v1 for more details.
# Description: makefile
# Author: zhangsong
# Create: 2019-01-18

SOURCES := $(shell find . 2>&1 | grep -E '.*\.(c|h|go)$$')
DEPS_LINK := $(CURDIR)/vendor/
VERSION := $(shell cat ./VERSION)
TAGS="cgo static_build"

GO_LDFLAGS="-s -w -extldflags=-zrelro -extldflags=-znow -X main.version=${VERSION}"
DEF_GOPATH=${GOPATH}
ifneq ($(GOPATH), )
CUS_GOPATH=${GOPATH}:${PWD}
ENV = GOPATH=${CUS_GOPATH} CGO_ENABLED=1
else
ENV = CGO_ENABLED=1
endif

all: toolkit lxcfs-hook
local: toolkit lxcfs-hook

toolkit:  $(SOURCES) | $(DEPS_LINK)
	@echo "Making isulad-lxcfs-tools..."
	${ENV} go build -mod=vendor -tags ${TAGS} -ldflags ${GO_LDFLAGS} -o build/isulad-lxcfs-toolkit .
	@echo "Done!"

lxcfs-hook: $(SOURCES) | $(DEPS_LINK)
	@echo "Making lxcfs-hook..."
	${ENV} go build -mod=vendor -tags ${TAGS} -ldflags ${GO_LDFLAGS} -o build/lxcfs-hook ./hooks/lxcfs-hook
	@echo "Done!"

clean:
	rm -rf build

install:
	cd hakc && ./install.shloacal:
