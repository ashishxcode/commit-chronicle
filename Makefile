BINARY := commit-chronicle
PKG    := ./cmd/commit-chronicle
BINDIR := bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build install run clean test release

## build: compile the binary for the host platform into bin/
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(BINARY) $(PKG)

## install: install worklog into $GOBIN (or $GOPATH/bin)
install:
	go install -ldflags "$(LDFLAGS)" $(PKG)

## run: build and run (pass args via ARGS=...)
run: build
	./$(BINDIR)/$(BINARY) $(ARGS)

## test: vet + build all packages
test:
	go vet ./...
	go build ./...

## clean: remove build artifacts
clean:
	rm -rf $(BINDIR) dist

## release: cross-compile static binaries into dist/
release:
	@mkdir -p dist
	@for target in \
		darwin/amd64 darwin/arm64 \
		linux/amd64 linux/arm64 \
		windows/amd64; do \
		os=$${target%/*}; arch=$${target#*/}; \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		out=dist/$(BINARY)-$$os-$$arch$$ext; \
		echo "→ $$out"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 \
			go build -ldflags "$(LDFLAGS)" -o $$out $(PKG) || exit 1; \
	done
	@echo "done. binaries in dist/"
