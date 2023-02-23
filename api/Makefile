CGO_ENABLED ?= off
GOOS        ?= linux
GOARCH      ?= amd64
GOPATH      ?= $(shell go env GOPATH)
GOBIN       ?= $(GOPATH)/bin
GOCACHE     ?= $(GOPATH)/cache
GOENV       ?= $(GOPATH)/env
GOPROXY     ?= https://goproxy.cn,https://proxy.golang.com.cn,direct
GOPRIVATE   ?= ''
GO          := GOOS=$(GOOS) GOARCH=$(GOARCH) GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOCACHE=$(GOCACHE) GOENV=$(GOENV) \
				GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) go

REINSTALL   ?= false
TOOLS       := \
	buf:github.com/bufbuild/buf/cmd/buf@latest \
	protoc-gen-go:google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	protoc-gen-go-grpc:google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
	protoc-gen-go-http:github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest \
	protoc-gen-validate:github.com/envoyproxy/protoc-gen-validate@latest \
	protoc-gen-openapi:github.com/google/gnostic/cmd/protoc-gen-openapi@latest

define check_install
	$(info Checking $(1))
	@[ -x "$(shell command -v $(1))" ] && [ "$(REINSTALL)" != "true" ] || (echo "Installing $(1)"; $(GO) install $(2))
endef

.PHONY: init
# init env
init:
	$(foreach tool, $(TOOLS), \
		$(call check_install, $(word 1, $(subst :, , $(tool))), $(word 2, $(subst :, , $(tool)))); \
	)

.PHONY: fmt
# format code
fmt:
	@echo buf format

.PHONY: lint
# check lint
lint:
	buf lint

.PHONY: break
# check breaking changes
break:
	buf breaking --against '.git#branch=main'

.PHONY: api
# generate api code
api:
	buf build
	buf generate

.PHONY: clean
# clean generate code
clean:
	@find . -mindepth 2 -not -path '*/\.*' -type f -not -name *.proto* -exec rm -v {} \;

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
