VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

PREFIX ?= /usr/local

.PHONY: build install uninstall test vet lint clean

build:
	go build -ldflags "-X github.com/kernel/leaseweb-cli/pkg/cmd.Version=$(VERSION)" -o lw ./cmd/lw

install: lw
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 lw $(DESTDIR)$(PREFIX)/bin/lw

lw: $(shell find . -name '*.go' -not -path './vendor/*')
	go build -ldflags "-X github.com/kernel/leaseweb-cli/pkg/cmd.Version=$(VERSION)" -o lw ./cmd/lw

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/lw

test:
	go test ./... -v

vet:
	go vet ./...

lint:
	golangci-lint run ./...

clean:
	rm -f lw
