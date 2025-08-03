default: install

# Build the provider
build:
	go build -o terraform-provider-netbox

# Install the provider locally for development
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/bab3l/netbox/0.1.0/windows_amd64
	cp terraform-provider-netbox.exe ~/.terraform.d/plugins/registry.terraform.io/bab3l/netbox/0.1.0/windows_amd64/

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run unit tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Check for issues
vet:
	go vet ./...

# Generate documentation
docs:
	tfplugindocs generate

# Validate documentation
docs-validate:
	tfplugindocs validate

# Clean build artifacts
clean:
	rm -f terraform-provider-netbox
	rm -f terraform-provider-netbox.exe

# Development cycle: format, vet, test, build, docs
dev: fmt vet test build docs

.PHONY: build install test testacc fmt vet docs docs-validate clean dev
