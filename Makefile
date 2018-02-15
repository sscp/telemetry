
# Installs 
.PHONY: install-tools
install-tools:
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go} 
	go get github.com/favadi/protoc-go-inject-tag

.PHONY: build-proto
build-proto:
	protoc -I=proto --go_out=proto ./proto/data.proto

# Build all the binaries in cmd
.PHONY: update-deps
update-deps:
	dep ensure && dep prune

# Tests all the packages (excludes vendor on go 1.9)
.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go install .
