BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)
LOCAL_DP_RENDERER_IN_USE = $(shell grep -c "\"github.com/ONSdigital/dp-renderer/v2\" =" go.mod)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: build
build: generate-prod
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-frontend-feedback-controller

.PHONY: lint 
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	golangci-lint run ./...

.PHONY: debug
debug: generate-debug
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-frontend-feedback-controller
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-feedback-controller

.PHONY: test
test: generate-prod
	go test -race -cover -tags 'production' ./...

.PHONY: test-component
test-component:
	exit

.PHONY: convey
convey:
	goconvey ./...

.PHONY: generate-debug
generate-debug: fetch-renderer-lib
	go get -u github.com/go-bindata/go-bindata/...
	go install github.com/go-bindata/go-bindata/...
	cd assets; go run github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs -prefix $(CORE_ASSETS_PATH)/assets -o data.go -pkg assets locales/... templates/... $(CORE_ASSETS_PATH)/assets/locales/... $(CORE_ASSETS_PATH)/assets/templates/...
	{ echo "// +build debug\n"; cat assets/data.go; } > assets/debug.go.new
	mv assets/debug.go.new assets/data.go

.PHONY: generate-prod
generate-prod: fetch-renderer-lib
	go get -u github.com/go-bindata/go-bindata/...
	go install github.com/go-bindata/go-bindata/...
	cd assets; go run github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs -prefix $(CORE_ASSETS_PATH)/assets -o data.go -pkg assets locales/... templates/... $(CORE_ASSETS_PATH)/assets/locales/... $(CORE_ASSETS_PATH)/assets/templates/...
	{ echo "// +build production\n"; cat assets/data.go; } > assets/data.go.new
	mv assets/data.go.new assets/data.go

.PHONY: fetch-dp-renderer
fetch-renderer-lib:
ifeq ($(LOCAL_DP_RENDERER_IN_USE), 1)
 $(eval CORE_ASSETS_PATH = $(shell grep -w "\"github.com/ONSdigital/dp-renderer/v2\" =>" go.mod | awk -F '=> ' '{print $$2}' | tr -d '"'))
else
 $(eval APP_RENDERER_VERSION=$(shell grep "github.com/ONSdigital/dp-renderer/v2" go.mod | cut -d ' ' -f2 ))
 $(eval CORE_ASSETS_PATH = $(shell go get github.com/ONSdigital/dp-renderer/v2@$(APP_RENDERER_VERSION) && go list -f '{{.Dir}}' -m github.com/ONSdigital/dp-renderer/v2))
endif
