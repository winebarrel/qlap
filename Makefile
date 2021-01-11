SHELL   := /bin/bash
VERSION := v0.1.0
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/qlap

.PHONY: package
package: clean build
	gzip qlap -c > qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz > qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum

.PHONY: clean
clean:
	rm -f qlap
