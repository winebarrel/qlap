SHELL   := /bin/bash
VERSION := v0.4.3
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)
ifeq ($(GOOS),windows)
	PROGRAM := qlap.exe
else
	PROGRAM := qlap
endif

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
	gzip $(PROGRAM) -c > qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz > qlap_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum

.PHONY: clean
clean:
	rm -f $(PROGRAM)
