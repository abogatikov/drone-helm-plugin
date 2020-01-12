RELEASE?=0.0.1
COMMIT?=git-$(shell git rev-parse --short HEAD)
COMPILETIME?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFFLAGS=-ldflags "-w -X github.com/abogatikov/drone-helm-plugin/pkg.Release=${RELEASE} -X github.com/abogatikov/drone-helm-plugin/pkg.CompileTime=${COMPILETIME} -X github.com/abogatikov/drone-helm-plugin/pkg.Commit=${COMMIT}"

# User definded
ARCH?=amd64
OS?=linux
CGO?=0

.PHONY: clean
clean:
	@echo "+ $@"
	@rm -rf ./build/*

.PHONY: dependencies
dependencies:
	@echo "+ $@"
	@go mod tidy

.PHONY: test
test:
	@echo "+ $@"
	@go test -cover ./pkg/...

.PHONY: lint
lint:
	@echo "+ $@"
	@golangci-lint run ./...

.PHONY: build
build: clean dependencies test lint
build:
	@echo "+ $@"
	@env GOARCH=${ARCH} GOOS=${OS} CGO_ENABLED=${CGO} go build -i -a ${LDFFLAGS} -o ./build/app ./cmd/main.go

.PHONY: build-linux
build-linux: ARCH = amd64
build-linux: OS = linux
build-linux: CGO = 0
build-linux: build

.PHONY: build-osx
build-osx: ARCH = amd64
build-osx: OS = darwin
build-osx: CGO = 0
build-osx: build
