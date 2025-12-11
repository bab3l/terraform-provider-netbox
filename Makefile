SHELL := /bin/sh

.PHONY: dev fmt vet test build install clean testacc

# Default version used for local install path
VERSION ?= 0.1.0
PROVIDER_NAME := terraform-provider-netbox
MODULE_NAMESPACE := bab3l
MODULE_NAME := netbox

dev: fmt vet test build

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./... -v

build:
	go build .

# Install the provider binary into the local Terraform plugin directory
# This supports cross-platform using GOOS/GOARCH from `go env`
install: build
	@GOOS=$$(go env GOOS); \
	GOARCH=$$(go env GOARCH); \
	PLATFORM="$$GOOS_$$GOARCH"; \
	BIN_NAME="$(PROVIDER_NAME)"; \
	OUT_DIR="$$HOME/.terraform.d/plugins/$(MODULE_NAMESPACE)/$(MODULE_NAME)/$(VERSION)/$$PLATFORM"; \
	mkdir -p "$$OUT_DIR"; \
	if [ "$$GOOS" = "windows" ]; then cp "$(PROVIDER_NAME).exe" "$$OUT_DIR/$(PROVIDER_NAME).exe"; else cp "$(PROVIDER_NAME)" "$$OUT_DIR/$(PROVIDER_NAME)"; fi; \
	echo "Installed to $$OUT_DIR"

clean:
	rm -f $(PROVIDER_NAME) $(PROVIDER_NAME).exe

# Acceptance tests require a running NetBox and environment variables:
# - NETBOX_SERVER_URL (e.g. http://localhost:8000)
# - NETBOX_API_TOKEN
testacc:
	@if [ -z "$$NETBOX_SERVER_URL" ]; then \
		echo "NETBOX_SERVER_URL is not set"; \
		echo "Please export NETBOX_SERVER_URL (e.g. http://localhost:8000)"; \
		exit 1; \
	fi
	@if [ -z "$$NETBOX_API_TOKEN" ]; then \
		echo "NETBOX_API_TOKEN is not set"; \
		echo "Please export NETBOX_API_TOKEN"; \
		exit 1; \
	fi
	TF_ACC=1 go test ./... -v -run "TestAcc"
