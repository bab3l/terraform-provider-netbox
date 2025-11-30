default: install

# Build the provider
build:
	go build -o terraform-provider-netbox

# Install the provider locally for development
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/bab3l/netbox/0.1.0/windows_amd64
	cp terraform-provider-netbox.exe ~/.terraform.d/plugins/registry.terraform.io/bab3l/netbox/0.1.0/windows_amd64/

# Run acceptance tests (requires NETBOX_SERVER_URL and NETBOX_API_TOKEN)
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

# Docker-based testing targets
docker-up:
	docker-compose up -d
	@echo "Waiting for Netbox to be ready..."
	@until curl -sf http://localhost:8000/api/ > /dev/null 2>&1; do sleep 5; done
	@echo "Netbox is ready at http://localhost:8000"

docker-down:
	docker-compose down -v

docker-logs:
	docker-compose logs -f netbox

# Run acceptance tests with Docker (starts Netbox, runs tests)
testacc-docker: docker-up
	NETBOX_SERVER_URL=http://localhost:8000 \
	NETBOX_API_TOKEN=0123456789abcdef0123456789abcdef01234567 \
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: build install test testacc fmt vet docs docs-validate clean dev docker-up docker-down docker-logs testacc-docker
