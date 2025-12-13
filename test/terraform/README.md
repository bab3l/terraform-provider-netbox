# Terraform Integration Tests

This directory contains Terraform configurations for integration testing the Netbox provider against a real Netbox instance.

## Directory Structure

```
test/terraform/
├── resources/
│   ├── site/              # Site resource tests
│   ├── site_group/        # Site Group resource tests
│   ├── tenant/            # Tenant resource tests
│   └── tenant_group/      # Tenant Group resource tests
├── data-sources/
│   ├── site/              # Site data source tests
│   ├── site_group/        # Site Group data source tests
│   ├── tenant/            # Tenant data source tests
│   └── tenant_group/      # Tenant Group data source tests
└── README.md              # This file
```

## Prerequisites

1. Docker Desktop running
2. Netbox container started: `docker-compose up -d`
3. Provider built: `go build .`
4. **Dev overrides configured** (see below)

## Setting Up Dev Overrides (Required)

To use a locally-built provider without publishing to the registry, you must configure Terraform's dev_overrides.

### Windows

Create or edit `%APPDATA%\terraform.rc`:

```hcl
provider_installation {
  dev_overrides {
    "bab3l/netbox" = "C:\\GitRoot\\terraform-provider-netbox"
  }
  direct {}
}
```

Replace `C:\\GitRoot\\terraform-provider-netbox` with the actual path to your provider source directory (where the built `terraform-provider-netbox.exe` is located).

### Linux/macOS

Create or edit `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "bab3l/netbox" = "/path/to/terraform-provider-netbox"
  }
  direct {}
}
```

Replace `/path/to/terraform-provider-netbox` with the actual path to your provider source directory.

### Important Notes on Dev Overrides

- **Skip `terraform init`**: When using dev_overrides, you should skip `terraform init` or expect a warning. The provider is loaded directly from the specified path.
- **Rebuild after changes**: After modifying provider code, rebuild with `go build .` for changes to take effect.
- **Binary name**: Terraform looks for a binary named `terraform-provider-netbox` (or `terraform-provider-netbox.exe` on Windows) in the dev_overrides path.

## Running Tests

### Run All Tests

```powershell
# PowerShell (Windows)
.\scripts\run-terraform-tests.ps1

# Options:
#   -SkipDestroy     Don't destroy resources after test
#   -ShowDetails     Show detailed output
#   -TestDir <path>  Run specific test directory
```

```bash
# Bash (Linux/macOS)
./scripts/run-terraform-tests.sh

# Options:
#   --skip-destroy   Don't destroy resources after test
#   --verbose        Show detailed output
#   --test <path>    Run specific test directory
```

### Run Specific Test Manually

With dev_overrides configured, you can skip init:

```bash
cd test/terraform/resources/site
# Skip init when using dev_overrides
terraform plan
terraform apply -auto-approve
terraform output
terraform destroy -auto-approve
```

Or run init (you'll see a warning but it will work):

```bash
cd test/terraform/resources/site
terraform init    # Warning about dev_overrides is expected
terraform apply -auto-approve
terraform output
terraform destroy -auto-approve
```

## Test Flow

Each test folder contains:
- `main.tf` - Provider and resource/data source configuration
- `outputs.tf` - Output values to verify the results

The test runner:
1. Builds the provider (`go build .`)
2. Runs `terraform plan` and `terraform apply`
3. Retrieves outputs and verifies that `*_valid` and `*_match` outputs are `true`
4. Destroys resources (cleanup)

**Note**: With dev_overrides, `terraform init` may show warnings but is not strictly required.

## Test Execution Order

Tests are run in dependency order:
1. tenant_group (no dependencies)
2. tenant (depends on tenant_group)
3. site_group (no dependencies)
4. site (depends on site_group)

## Environment Variables

Set these before running tests (or let the runner use defaults):
```powershell
$env:NETBOX_SERVER_URL = "http://localhost:8000"
$env:NETBOX_API_TOKEN = "0123456789abcdef0123456789abcdef01234567"
```

## Troubleshooting

### Common Issues

1. **Provider not found**:
   - Ensure dev_overrides are configured correctly in your terraform.rc
   - Rebuild the provider with `go build .`
   - Verify the binary exists in the dev_overrides path

2. **"Failed to query available provider packages" error**:
   - This means dev_overrides are not configured
   - See "Setting Up Dev Overrides" section above

3. **Cannot connect to Netbox**:
   - Ensure Netbox is running with `docker-compose up -d`
   - Wait for it to be healthy (check with `docker-compose ps`)

4. **API errors**:
   - Check the Netbox logs with `docker-compose logs netbox`

5. **State conflicts**:
   - Remove `.terraform`, `.terraform.lock.hcl`, and `terraform.tfstate*` files from the test directory

## Verifying Dev Overrides

To check if dev_overrides are configured:

```powershell
# Windows
Get-Content "$env:APPDATA\terraform.rc"
```

```bash
# Linux/macOS
cat ~/.terraformrc
```

You should see output containing `dev_overrides` with the path to your provider.
