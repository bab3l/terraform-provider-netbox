# Development Configuration

This file shows how to configure Terraform to use your locally built provider during development.

## 1. Create a `.terraformrc` file in your home directory

Windows: `%APPDATA%\terraform.rc`
Linux/macOS: `~/.terraformrc`

```hcl
provider_installation {
  dev_overrides {
    "bab3l/netbox" = "C:\\GitRoot\\terraform-provider-netbox"
  }

  # For all other providers, install them directly as normal.
  direct {}
}
```

## 2. Build the provider

```bash
go build -o terraform-provider-netbox.exe .
```

## 3. Create a test Terraform configuration

Create a `test.tf` file:

```hcl
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "https://your-netbox-instance.com"
  api_token  = "your-api-token"
}

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}
```

## 4. Test the provider

```bash
terraform init
terraform plan
terraform apply
```

## Environment Variables for Testing

Set these environment variables for easier testing:

```bash
export NETBOX_SERVER_URL="https://your-netbox-instance.com"
export NETBOX_API_TOKEN="your-api-token"
export TF_LOG=DEBUG  # Enable debug logging
```

## Running Tests

### Test Architecture
This provider uses a split test architecture:
- **Parallel-safe tests** in `internal/resources_acceptance_tests/` - Run concurrently (fast)
- **Custom field tests** in `internal/resources_acceptance_tests_customfields/` - Run serially (slow)

Custom field tests are separated because NetBox custom fields are global per content type. Running them in parallel causes database deadlocks and race conditions.

### Development Workflow

#### Fast Development Cycle (30-40 minutes)
```bash
make test-acceptance
```
Runs ~150 parallel-safe acceptance tests. **Use this for rapid iteration.**

#### Full Test Suite (2-3 hours)
```bash
make test-acceptance-all
```
Runs both parallel and serial tests. **Run before submitting PRs.**

#### Custom Field Tests Only (60-90 minutes)
```bash
make test-acceptance-customfields
```
Runs only the custom field tests serially.

#### Unit Tests Only (1-2 minutes)
```bash
make test-fast
```
Runs unit tests without requiring NetBox.

### Running Individual Tests
```bash
# Parallel test (no build tag needed)
go test ./internal/resources_acceptance_tests -run TestAccSiteResource_basic -v

# Custom field test (requires build tag)
TF_ACC=1 go test -tags=customfields ./internal/resources_acceptance_tests_customfields -run TestAccSiteResource_importWithCustomFieldsAndTags -v
```

### Why Split Tests?
NetBox custom fields are defined globally per object type (e.g., `dcim.device`, `ipam.aggregate`). When multiple tests create/delete custom fields for the same object type in parallel:
- Database deadlocks occur
- Race conditions cause test failures
- Tests become non-deterministic

By separating custom field tests into a package with build tag `customfields`, we:
- Speed up normal development (skip slow serial tests)
- Prevent conflicts (serial execution via `-p 1` flag)
- Make CI more efficient (run parallel and serial suites separately)

## Debugging

To debug the provider:

1. Build with debug flags:
   ```bash
   go build -gcflags="all=-N -l" -o terraform-provider-netbox.exe .
   ```

2. Run the provider in debug mode:
   ```bash
   .\terraform-provider-netbox.exe -debug
   ```

3. Use the provided TF_REATTACH_PROVIDERS environment variable in your Terraform commands.

## Linting (pre-commit)

This repo uses pre-commit hooks to run linting/formatting checks.

- `golangci-lint` formatting is run via a small Go wrapper at `scripts/golangci_wrapper/`.
  The wrapper sets a repo-local `GOTMPDIR` (in `.gotmp/`) before invoking `golangci-lint`.
  This helps avoid occasional Windows file-lock issues when Go creates temporary `.exe` files.

## Optional Field Null Handling

All optional fields must explicitly handle null values to prevent "inconsistent result" errors when fields are removed from Terraform configuration.

### Pattern for Optional Fields

```go
// In buildXRequest() function:
if !data.FieldName.IsNull() && !data.FieldName.IsUnknown() {
    request.SetFieldName(data.FieldName.ValueString())
} else if data.FieldName.IsNull() {
    request.SetFieldName("")  // For strings: clear with empty string
    // request.SetFieldNameNil()  // For nullable types: use Nil setter
}
```

### Testing Optional Fields

Use the test helpers in `internal/testutil/optional_fields_test_helpers.go`:

```go
func TestAccResource_removeOptionalFields(t *testing.T) {
    t.Parallel()

    testutil.TestRemoveOptionalFields(t, testutil.OptionalFieldTestConfig{
        ResourceType: "netbox_resource",
        ResourceName: "test",
        CreateConfigWithOptionalFields:    func() string { return configWithFields() },
        CreateConfigWithoutOptionalFields: func() string { return configMinimal() },
        OptionalFields: map[string]string{
            "description": "Test description",
            "comments":    "Test comments",
        },
    })
}
```

See [docs/TESTING_OPTIONAL_FIELDS.md](docs/TESTING_OPTIONAL_FIELDS.md) for complete documentation and examples.
