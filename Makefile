OS = $(shell uname | tr '[:upper:]' '[:lower:]')
SHELL:=/bin/bash

ifeq ($(OS), darwin)
	DYLIB_EXT = dylib
else
	DYLIB_EXT = so
endif

WORK_DIR=$(shell pwd)
BUILD_DIR=$(WORK_DIR)/build
PKG_CONFIG_DIR=$(BUILD_DIR)/pkgconfig

LIBOQS_DIR=$(BUILD_DIR)/liboqs
LIBOQS_OBJ=$(LIBOQS_DIR)/lib/liboqs.a
LIBOQS_PKG_CONFIG=$(PKG_CONFIG_DIR)/liboqs.pc

PROTO_RESOURCE_DIR=$(WORK_DIR)/resources/protobuf/go/abelian.info/sdk/proto
PROTO_SRC_DIR=$(WORK_DIR)/proto
CORE_PB_GO=core.pb.go

DYLIB_BIN=libabelsdk.$(DYLIB_EXT)

build: $(BUILD_DIR)/$(DYLIB_BIN)

clean:
	rm -rf $(BUILD_DIR)/$(DYLIB_BIN)

$(BUILD_DIR)/$(DYLIB_BIN): $(PROTO_SRC_DIR)/$(CORE_PB_GO) $(LIBOQS_PKG_CONFIG)
	@echo "==> Building $(DYLIB_BIN) ..."
	PKG_CONFIG_PATH=$(PKG_CONFIG_DIR) go build -buildmode=c-shared -o $(BUILD_DIR)/$(DYLIB_BIN)

$(PROTO_SRC_DIR)/$(CORE_PB_GO): $(PROTO_RESOURCE_DIR)/$(CORE_PB_GO)
	@echo "==> Copying core.pb.go ..."
	cp $(PROTO_RESOURCE_DIR)/$(CORE_PB_GO) $(PROTO_SRC_DIR)

$(LIBOQS_PKG_CONFIG): $(LIBOQS_OBJ)
	@echo "==> Generating liboqs.pc ..."
	@if [ ! -d "${PKG_CONFIG_DIR}" ]; then mkdir -p ${PKG_CONFIG_DIR}; fi
	@echo "Name: liboqs" > $(LIBOQS_PKG_CONFIG)
	@echo "Description: C library for quantum resistant cryptography" >> $(LIBOQS_PKG_CONFIG)
	@echo "Version: 0.7.2-dev" >> $(LIBOQS_PKG_CONFIG)
	@echo "Cflags: -I$(LIBOQS_DIR)/include" >> $(LIBOQS_PKG_CONFIG)
	@echo "Ldflags: '-extldflags \"-static -Wl,-stack_size -Wl,0x1000000\"'" >> $(LIBOQS_PKG_CONFIG)
	@echo "Libs: -L$(LIBOQS_DIR)/lib -l:liboqs.a -lcrypto" >> $(LIBOQS_PKG_CONFIG)

$(LIBOQS_OBJ):
	@if [ ! -d "${BUILD_DIR}" ]; then mkdir -p ${BUILD_DIR}; fi
	@if [ ! -d "${LIBOQS_DIR}" ]; then echo "==> Fetching liboqs ..."; git clone https://github.com/cryptosuite/liboqs.git ${LIBOQS_DIR}; fi
	@echo "==> Compiling liboqs ..."
	cd ${LIBOQS_DIR} && cmake -GNinja . && ninja

setenv:
	go env -w GOPRIVATE=github.com/pqabelian/*
	git config --global url."git@github.com:".insteadOf https://github.com/

unsetenv:
	go env -w GOPRIVATE=
	git config --global --unset url."git@github.com:".insteadOf https://github.com/
