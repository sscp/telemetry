
# Add executables as they are added to cmd folder here
.PHONY: run-hello
run-hello:
	${GOPATH}/bin/hello

# Tests all the packages (excludes vendor on go 1.9)
.PHONY: test
test:
	go test ./...

# Build all the binaries in cmd
.PHONY: update-deps
update-deps:
	dep ensure && dep prune

# Build all the binaries in cmd
.PHONY: build
build:
	go install ./cmd/...

