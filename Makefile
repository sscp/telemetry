
BINARY_NAME=telemetry

# Installs tools to generate code/work with the repo
.PHONY: install-tools
install-tools:
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go} 
	go get -u github.com/favadi/protoc-go-inject-tag
	go get -u google.golang.org/grpc

.PHONY: generate
generate:
	go generate ./...

# Build all the binaries in cmd, requires https://github.com/golang/dep
.PHONY: update-deps
update-deps:
	dep ensure

# Tests all the packages (excludes vendor on go 1.9)
.PHONY: test
test:
	go test -v ./...

.PHONY: install
install:
	go install .

.PHONY: release
release:
	go get github.com/goreleaser/goreleaser
	goreleaser --rm-dist

