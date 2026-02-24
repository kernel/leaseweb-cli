VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: build install test vet lint clean

build:
	go build -ldflags "-X github.com/kernel/leaseweb-cli/pkg/cmd.Version=$(VERSION)" -o lw ./cmd/lw

install:
	go install -ldflags "-X github.com/kernel/leaseweb-cli/pkg/cmd.Version=$(VERSION)" ./cmd/lw

test:
	go test ./... -v

vet:
	go vet ./...

lint:
	golangci-lint run ./...

clean:
	rm -f lw
