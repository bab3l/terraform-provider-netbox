# Terraform Provider for Netbox

[![License: MPL-2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

A Terraform provider for [Netbox](https://github.com/netbox-community/netbox) using the modern [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) and the [go-netbox](https://github.com/bab3l/go-netbox) API wrapper.

License: This project is licensed under the Mozilla Public License 2.0 (MPL-2.0). See `LICENSE`.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21
- [Netbox](https://github.com/netbox-community/netbox) **v4.1.11** (tested version)

> **Note:** This provider is built and tested against Netbox v4.1.11. Other versions may work but are not officially supported.

## Building the Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the provided VS Code task or Go directly:

```powershell
# VS Code task
# Terminal: Run Task → Build Terraform Provider

# Or via Go directly
go build .
```

## Adding the Provider to your Terraform Configuration

```hcl
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = "~> 0.1.0"
    }
  }
}

provider "netbox" {
  server_url = "https://netbox.example.com"
  api_token  = var.netbox_api_token
  insecure   = false  # Set to true to skip TLS verification
}
```

## Authentication

The provider supports authentication via API token. You can provide the token in several ways:

1. **Provider configuration** (not recommended for production):
   ```hcl
   provider "netbox" {
     server_url = "https://netbox.example.com"
     api_token  = "your-token-here"
   }
   ```

2. **Environment variables** (recommended):
   ```bash
   export NETBOX_SERVER_URL="https://netbox.example.com"
   export NETBOX_API_TOKEN="your-token-here"
   export NETBOX_INSECURE="false"  # Optional
   ```

3. **Terraform variables**:
   ```hcl
   variable "netbox_api_token" {
     description = "Netbox API token"
     type        = string
     sensitive   = true
   }

   provider "netbox" {
     server_url = "https://netbox.example.com"
     api_token  = var.netbox_api_token
   }
   ```

## Development

### Prerequisites

- Go 1.21+
- Make (optional, but recommended)

### Local Development Setup

1. Clone this repository
2. Ensure you have the go-netbox dependency available locally at `../go-netbox`
3. Run the development cycle:

```powershell
make dev  # Runs format, vet, test, and build
```

### Installing Provider Locally

To install the provider locally for testing:

```powershell
make install
```

This builds and installs the provider to `~/.terraform.d/plugins/bab3l/netbox/0.1.0/<GOOS>_<GOARCH>/`.

### Running Tests

```powershell
# Run unit tests
make test

# Run acceptance tests (requires NetBox and env vars)
$env:TF_ACC = "1"
$env:NETBOX_SERVER_URL = "http://localhost:8000"
$env:NETBOX_API_TOKEN = "your-token"
go test ./... -v -run "TestAcc"

# Or via Make alias (env vars must be set)
make testacc
```

### Project Structure

```
.
├── internal/
│   ├── provider/          # Provider implementation
│   ├── resources/         # Resource implementations
│   └── datasources/       # Data source implementations
├── examples/              # Example Terraform configurations
├── docs/                  # Generated documentation
├── main.go               # Provider entry point
├── Makefile              # Development tasks
└── README.md
```

## Resources and Data Sources

For a complete list of available resources and data sources, including detailed documentation and examples, see the [docs/](./docs/) directory.

📖 **[View Full Documentation Index](./docs/README.md)**

### Quick Examples

```hcl
# Site group for organizational hierarchy
resource "netbox_site_group" "region" {
  name = "North America"
  slug = "north-america"
}

# Site within the group
resource "netbox_site" "datacenter" {
  name   = "Primary DC"
  slug   = "primary-dc"
  group  = netbox_site_group.region.id
  status = "active"
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Run the test suite
6. Create a pull request

## License

This project is licensed under the Mozilla Public License 2.0 (MPL-2.0). See `LICENSE`.

[Contributing Guidelines](CONTRIBUTING.md) • [Security Policy](SECURITY.md)

## Credits

- Built with [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
- Uses [go-netbox](https://github.com/bab3l/go-netbox) for API interactions
- Inspired by the [Netbox community](https://github.com/netbox-community/netbox)
