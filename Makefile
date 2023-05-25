LOCAL_BIN:=$(CURDIR)/bin

build:
	@go build -o ${LOCAL_BIN}/poop-server cmd/server/main.go
	@go build -o ${LOCAL_BIN}/poop-sync cmd/sync/main.go
	@echo "build binary file to ${LOCAL_BIN}"

dependency:
	@go mod tidy
