# Simple Makefile for repo with main.go at project root and app/ package.
GO ?= go
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
DIST_DIR := dist
PROJECT := inventory-app
VERSION_FILE ?= VERSION

VERSION ?= $(shell if [ -f "$(VERSION_FILE)" ]; then cat "$(VERSION_FILE)"; else git describe --tags --always 2>/dev/null || printf 'v0.0.0'; fi)
BINARY := $(BIN_DIR)/$(PROJECT)
INSTALL_DIR ?= install-dir

CONFIG_DIR="${CONFIG_DIR:-./config}"

.PHONY: all build install install-deps package clean help
all: build

install-deps:
	@if [ -f go.mod ]; then \
		echo "Installing Go module dependencies from go.mod..."; \
		$(GO) mod download; \
		$(GO) mod verify; \
	else \
		echo "No go.mod found, skipping go module install."; \
	fi

build: install-deps $(BINARY)

run: build
	@echo "Running $(BINARY)..."
	@$(BINARY)

$(BINARY): main.go
	@mkdir -p $(BIN_DIR)
	@echo "Building $(BINARY)..."
	$(GO) build -o $@ .
	@chmod 0755 $@

install: build
	@./scripts/install.sh --build-dir "$(BUILD_DIR)" --config-dir "$(CONFIG_DIR)" --install-dir "$(INSTALL_DIR)"

package: build
	@./scripts/package.sh --install-src-dir "$(INSTALL_DIR)" --dist-dir "$(DIST_DIR)" --version "$(VERSION)" --app-name "$(PROJECT)"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(INSTALL_DIR)

help:
	@echo "Targets: all build install install-deps package clean"
	@echo " - build: compile main.go -> $(BINARY)"
	@echo " - install: copies binary to INSTALL_DIR (default ./install-dir) along with config files from CONFIG_DIR (default ./config)"
	@echo " - package: tar.gz in $(DIST_DIR)"
