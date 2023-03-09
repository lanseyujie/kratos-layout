PROJECT     := $(notdir $(CURDIR))
COMMIT      := $(shell git describe --tags --always 2> /dev/null || echo -n '0.0.0')
CODENAME    ?= kratos
MODE        ?= dev
VERSION     ?= $(COMMIT)

REGISTRY    ?= registry.aliyuncs.com
USERNAME    ?= ''
PASSWORD    ?= ''
MIRROR      ?= repo.huaweicloud.com

CGO_ENABLED ?= off
GOOS        ?= $(shell go env GOOS)
GOARCH      ?= $(shell go env GOARCH)
GOPROXY     ?= https://goproxy.cn,https://proxy.golang.com.cn,direct
GOPRIVATE   ?= $(shell go env GOPRIVATE)
GO          := GOOS=$(GOOS) GOARCH=$(GOARCH) GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) go
GOBUILD     := CGO_ENABLED=$(CGO_ENABLED) $(GO) build
GOTEST      := CGO_ENABLED=on $(GO) test -race -count=1 -cover -v
GCFLAGS     += -d=ssa/check_bce
LDFLAGS     += -s -w -extldflags "-static"
LDFLAGS     += -X "main.Version=$(VERSION)"

BUILD_FLAGS := -v -trimpath -tags $(MODE) -gcflags='$(GCFLAGS)' -ldflags '$(LDFLAGS)'

REINSTALL   ?= false
TOOLS       := \
	buf:github.com/bufbuild/buf/cmd/buf@latest \
	protoc-gen-go:google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	wire:github.com/google/wire/cmd/wire@latest \
	golangci-lint:github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
	kratos:github.com/go-kratos/kratos/cmd/kratos/v2@latest \
	ent:entgo.io/ent/cmd/ent@latest

define check_install
	$(info Checking $(1))
	@[ -x "$(shell command -v $(1))" ] && [ "$(REINSTALL)" != "true" ] || (echo "Installing $(1)"; $(GO) install $(2))
endef

define list_apps
	$(foreach dir, $(DIRS), \
		$(foreach app, $(notdir $(wildcard ./app/$(dir)/cmd/*)), \
			$(PROJECT)-$(dir)-$(app) \
		) \
	)
endef

DIRS        := $(notdir $(wildcard ./app/*))
APPS        := $(call list_apps)

.PHONY: init
# init env
init:
	$(foreach tool, $(TOOLS), \
		$(call check_install, $(word 1, $(subst :, , $(tool))), $(word 2, $(subst :, , $(tool)))); \
	)

.PHONY: version
# display version
version:
	@echo api version: $(shell cd ./api && git describe --tags --always 2> /dev/null || echo $(VERSION))
	@echo app version: $(VERSION)
	@echo apps: $(APPS)

.PHONY: api
# generate api code
api:
	@[ -f ./.gitmodules ] && [ ! -e ./api/.git ] && git submodule update --init --recursive || true
	@cd api && make init && make api

.PHONY: config
# generate internal proto
config:
	buf generate --exclude-path ./api

.PHONY: generate
# generate code
generate:
	@$(GO) mod tidy
	@$(GO) get github.com/google/wire/cmd/wire@latest
	@$(GO) generate -x ./...
	wire ./...

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;

.PHONY: fmt
# format code
fmt:
	gofmt -w -r 'interface{} -> any' .
	# gofmt -w -r '"2006-01-02 15:04:05" -> time.DateTime' .
	# gofmt -w -r '"2006-01-02" -> time.DateOnly' .
	# gofmt -w -r '"15:04:05" -> time.TimeOnly' .

.PHONY: lint
# code lint
lint:
	golangci-lint run -v

.PHONY: test
# test all units
test:
	@$(foreach dir, $(DIRS), \
		make $@/$(dir) || exit $$?; \
	)

test/%:
	@echo "Testing app: $*"
	@cd ./app/$* && $(GOTEST) ./...

.PHONY: build
# build all apps
build:
	@$(foreach app, $(APPS), \
		make $@/$(app) || exit $$?; \
	)

build/%:
	@echo "Building app: $*"
	mkdir -p ./bin && cd ./app/$(subst -,/cmd/,$(subst $(PROJECT)-,,$*)) && $(GOBUILD) $(BUILD_FLAGS) -o $(CURDIR)/bin/$*
	@$(if $(shell command -v upx), upx -q -f $(CURDIR)/bin/$*, )

.PHONY: login
# login to a registry
login:
	@docker login $(REGISTRY) -u $(USERNAME) --password-stdin <<< $(PASSWORD)

.PHONY: image
# build all images
image:
	@$(foreach app, $(APPS), \
		make $@/$(app) || exit $$?; \
	)

image/%:
	@echo "Building image for $*"
	docker build --rm \
		--build-arg MIRROR=$(MIRROR) \
		--build-arg TIMEZONE=Asia/Shanghai \
		--build-arg GOPROXY=$(GOPROXY) \
		--build-arg GOPRIVATE=$(GOPRIVATE) \
		--build-arg PROJECT=$(PROJECT) \
		--build-arg APP=$* \
		-t $(REGISTRY)/$(CODENAME)/$*:$(VERSION) .
	@echo "#docker push $(REGISTRY)/$(CODENAME)/$*:$(VERSION)"

.PHONY: clean
# clean all binaries
clean:
	-go clean -x -cache -testcache
	-rm -r ./bin/*

.PHONY: manifest
# deploy all apps
deploy:
	@echo TODO:// kubectl apply

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
