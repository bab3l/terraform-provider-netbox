# Testing Infrastructure for terraform-provider-netbox

## Overview

This document covers the complete testing infrastructure for the Terraform Netbox provider, including unit tests, acceptance tests, and Docker-based integration testing.

## Quick Start: Running Tests

### Unit Tests (No Netbox Required)
```bash
# Run all unit tests (excludes acceptance tests)
go test -run "^Test[^Acc]" ./internal/...

# Run schema safety tests
go test ./internal/resources_test/schema_safety_test.go -v
go test ./internal/datasources_test/... -v
```

### Schema Safety Tests

The provider includes "safety tests" that verify schema consistency and catch common patterns that lead to runtime errors:

| Test Category | What It Catches |
|---------------|-----------------|
| `TestAllResourcesHaveIDField` | Missing `id` attribute (Terraform requirement) |
| `TestOptionalFieldsAreNotRequired` | Incorrectly marking optional fields as required |
| `TestRequiredFieldsAreMarkedRequired` | Missing required field markers |
| `TestComputedFieldsAreMarkedComputed` | Fields set by provider but not marked Computed |
| `TestOptionalStringFieldsAllowNull` | Prevents "was null, but now cty.StringVal("")" errors |
| `TestResourceMetadataPrefix` | Incorrect resource type names |
| `TestAllDeclaredFieldsExistInSchema` | Accidental field removal |
| `TestNoUnexpectedFieldsInSchema` | Catches new fields needing test updates |

These tests are defined in:
- `internal/resources_test/schema_safety_test.go` - Resource schema tests
- `internal/datasources_test/schema_safety_test.go` - Data source schema tests

### Acceptance Tests with Docker (Recommended)
```powershell
# Windows PowerShell
.\scripts\run-acceptance-tests.ps1
```

```bash
# Linux/macOS
./scripts/run-acceptance-tests.sh
```

---

## Docker-Based Integration Testing

The provider includes a complete Docker Compose environment for running acceptance tests against a real Netbox instance.

### Prerequisites

- Docker Desktop (Windows/macOS) or Docker Engine (Linux)
- docker-compose
- Go 1.21+

### Docker Environment

The `docker-compose.yml` provides:
- **Netbox** (v4.1) - The main application
- **PostgreSQL 15** - Database backend
- **Redis 7** - Caching and task queue

### Starting the Test Environment

#### Option 1: Using the Test Script (Recommended)

```powershell
# Windows - Start Netbox and run all tests
.\scripts\run-acceptance-tests.ps1

# Start only (for manual testing)
.\scripts\run-acceptance-tests.ps1 -StartOnly

# Run specific tests
.\scripts\run-acceptance-tests.ps1 -TestPattern "TestAccSite"

# Stop and clean up
.\scripts\run-acceptance-tests.ps1 -StopOnly
```

```bash
# Linux/macOS
./scripts/run-acceptance-tests.sh
./scripts/run-acceptance-tests.sh --start-only
./scripts/run-acceptance-tests.sh --pattern "TestAccSite"
./scripts/run-acceptance-tests.sh --stop-only
```

#### Option 2: Manual Docker Commands

```bash
# Start the environment
docker-compose up -d

# Wait for Netbox to be healthy (may take 2-3 minutes on first run)
docker-compose logs -f netbox

# Run tests manually
$env:NETBOX_SERVER_URL = "http://localhost:8000"
$env:NETBOX_API_TOKEN = "0123456789abcdef0123456789abcdef01234567"
$env:TF_ACC = "1"
go test ./... -v

# Stop and clean up
docker-compose down -v
```

### Test Environment Details

| Service | URL | Credentials |
|---------|-----|-------------|
| Netbox Web UI | http://localhost:8000 | admin / admin |
| Netbox API | http://localhost:8000/api/ | Token: `0123456789abcdef0123456789abcdef01234567` |
| PostgreSQL | localhost:5432 | netbox / netbox |
| Redis | localhost:6379 | (no password) |

### Troubleshooting Docker

**Netbox takes a long time to start:**
First-time startup can take 2-3 minutes while the database is initialized. Check progress with:
```bash
docker-compose logs -f netbox
```

**Port conflicts:**
If port 8000 is in use, modify `docker-compose.yml`:
```yaml
ports:
  - "8080:8080"  # Change the first port
```

**Clean restart:**
```bash
docker-compose down -v  # Remove volumes
docker-compose up -d    # Fresh start
```

---

## Current Test Coverage

### ✅ Provider Tests (`internal/provider/provider_test.go`)

1. **Basic Instantiation Test**
   - Ensures the provider can be created without panicking
   - Validates the provider factory function works correctly

2. **Schema Validation Test**
   - Verifies the provider schema can be retrieved
   - Checks that essential attributes (`server_url`, `api_token`, `insecure`) exist
   - Ensures schema generation doesn't produce errors

3. **Resource Registration Test**
   - Confirms the provider registers resources correctly
   - Validates that resource factory functions work

4. **Data Source Registration Test**
   - Tests data source registration (currently empty, but framework ready)

### ✅ Validator Tests (`internal/validators/validators_test.go`)

1. **Slug Validation Tests**
   - Valid cases: lowercase, numbers, hyphens, underscores
   - Invalid cases: uppercase, spaces, special characters
   - Edge cases: starting/ending with hyphens/underscores
   - Null/unknown value handling

2. **Custom Field Value Tests**
   - Type-specific validation for all supported types
   - Integer parsing validation
   - Boolean format validation
   - JSON syntax validation
   - URL format validation
   - Multiselect format validation

### ✅ Resource Tests (`internal/resources/site_resource_test.go`)

1. **Resource Instantiation Test**
   - Ensures site resource can be created
   
2. **Schema Validation Test**
   - Verifies all required and optional attributes exist
   - Ensures schema generation works correctly

3. **Metadata Test**
   - Validates resource type name is correct

4. **Configuration Test**
   - Tests provider data handling (nil, correct type, incorrect type)

## Test Infrastructure Components

### `testAccProtoV6ProviderFactories`
- **Purpose**: Critical for acceptance testing
- **Usage**: Used by resource acceptance tests to create provider instances
- **Required**: Yes, this is used by the terraform-plugin-testing framework

### Provider Factory Function
- **Purpose**: Creates provider instances for testing
- **Flexibility**: Allows testing with different configurations (version="test")

## Why This Testing is Important

### 1. **Continuous Integration**
```bash
go test ./internal/... -v
```
These tests run in CI/CD pipelines to catch regressions early.

### 2. **Schema Validation**
- Prevents breaking changes to provider schema
- Ensures required attributes are always present
- Validates schema generation doesn't error

### 3. **Type Safety**
- Tests that resources can be instantiated
- Validates configuration handling
- Ensures proper error handling for invalid configurations

### 4. **Validator Reliability**
- Comprehensive testing of custom validation logic
- Edge case handling
- Error message verification

### 5. **Acceptance Test Foundation**
- `testAccProtoV6ProviderFactories` is required for acceptance tests
- Enables end-to-end testing with real Terraform operations

## Running Tests

### Unit Tests
```bash
# All tests
go test ./internal/... -v

# Provider tests only
go test ./internal/provider -v

# Validator tests only  
go test ./internal/validators -v

# Resource tests only
go test ./internal/resources -v
```

### Test Coverage
```bash
go test ./internal/... -v -cover
```

### Race Condition Detection
```bash
go test ./internal/... -v -race
```

## Future Test Expansion

### Acceptance Tests
The current foundation enables adding acceptance tests:

```go
func TestAccSiteResource_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Test steps here
        },
    })
}
```

### Integration Tests
- Tests with real Netbox API (when available)
- Mock API server tests
- Error scenario testing

### Performance Tests
- Schema generation performance
- Large configuration handling
- Memory usage validation

## Conclusion

The `provider_test.go` file is **absolutely useful and necessary**. It provides:

1. **Foundation for all testing** - Required by terraform-plugin-testing framework
2. **Schema validation** - Prevents breaking changes
3. **Regression prevention** - Catches issues early in development
4. **Documentation** - Shows how to use the provider programmatically
5. **Quality assurance** - Ensures provider reliability

**Recommendation**: Keep and expand the testing infrastructure rather than removing it. It's a critical component for maintaining a high-quality Terraform provider.
