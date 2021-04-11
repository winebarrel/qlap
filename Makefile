SHELL   := /bin/bash
VERSION := v1.3.0
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: vet build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/qlap

.PHONY: vet
vet:
	go vet

.PHONY: package
package: clean vet build
ifeq ($(GOOS),windows)
	zip qlap_$(VERSION)_$(GOOS)_$(GOARCH).zip qlap.exe
	sha1sum qlap_$(VERSION)_$(GOOS)_$(GOARCH).zip > qlap_$(VERSION)_$(GOOS)_$(GOARCH).zip.sha1sum
else
	gzip qlap -c > qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz > qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum
endif

.PHONY: clean
clean:
	rm -f qlap qlap.exe
