# DataSource Unit Testing Guide

## Overview

This document describes the standard patterns and helper functions for unit testing datasources in the Terraform Provider for Netbox. All datasources should follow this pattern for consistency and comprehensive coverage.

## Quick Reference

**Standard Unit Tests (3 tests per datasource):**

1. **TestDataSourceNameSchema** - Validates schema structure
2. **TestDataSourceNameMetadata** - Validates type name
3. **TestDataSourceNameConfigure** - Validates Configure method

**File Layout:**
- Directory: `internal/datasources_unit_tests/`
- Package: `datasources_unit_tests`
- Imports: `datasources`, `testutil`, `datasource`, `context`, `testing`

---

## Helper Functions

All helper functions are located in `internal/testutil/schema_validation.go`.

### ValidateDataSourceSchema

Validates that a datasource schema contains expected attributes with correct optionality and computability.

**Signature:**
```go
func ValidateDataSourceSchema(t *testing.T, schemaAttrs map[string]datasourceschema.Attribute, validation DataSourceValidation)
```

**DataSourceValidation Type:**
```go
type DataSourceValidation struct {
    LookupAttrs []string    // Attributes that should be optional (used for lookups)
    ComputedAttrs []string  // Attributes that should be computed
}
```

**Usage Example:**
```go
func TestAggregateDataSourceSchema(t *testing.T) {
    d := datasources.NewAggregateDataSource()

    req := datasource.SchemaRequest{}
    resp := &datasource.SchemaResponse{}

    d.Schema(context.Background(), req, resp)

    if resp.Diagnostics.HasError() {
        t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
    }

    testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
        LookupAttrs: []string{"id", "prefix", "rir"},
        ComputedAttrs: []string{"tenant", "date_added", "description", "comments", "tags"},
    })
}
```

### ValidateDataSourceMetadata

Validates that a datasource's metadata correctly identifies its type name.

**Signature:**
```go
func ValidateDataSourceMetadata(t *testing.T, d datasource.DataSource, providerTypeName, expectedTypeName string)
```

**Usage Example:**
```go
func TestAggregateDataSourceMetadata(t *testing.T) {
    d := datasources.NewAggregateDataSource()
    testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_aggregate")
}
```

### ValidateDataSourceConfigure

Validates that a datasource's Configure method handles provider data correctly.

Tests three scenarios:
1. Configure with nil provider data (backwards compatibility)
2. Configure with valid APIClient (success case)
3. Configure with invalid provider data (error case)

**Signature:**
```go
func ValidateDataSourceConfigure(t *testing.T, d datasource.DataSource)
```

**Usage Example:**
```go
func TestAggregateDataSourceConfigure(t *testing.T) {
    d := datasources.NewAggregateDataSource()
    testutil.ValidateDataSourceConfigure(t, d)
}
```

---

## Standard Unit Test Template

Use this template for all datasources. This provides comprehensive unit test coverage with minimal code duplication.

### File: `internal/datasources_unit_tests/resource_name_data_source_test.go`

```go
package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// TestResourceNameDataSourceSchema validates the datasource schema structure.
func TestResourceNameDataSourceSchema(t *testing.T) {
	d := datasources.NewResourceNameDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs: []string{
			"id",
			// Add other lookup attributes (typically ID and name/identifier fields)
		},
		ComputedAttrs: []string{
			// Add other computed attributes (everything except lookup fields)
		},
	})
}

// TestResourceNameDataSourceMetadata validates the datasource type name.
func TestResourceNameDataSourceMetadata(t *testing.T) {
	d := datasources.NewResourceNameDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_resource_name")
}

// TestResourceNameDataSourceConfigure validates the Configure method.
func TestResourceNameDataSourceConfigure(t *testing.T) {
	d := datasources.NewResourceNameDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
```

**Notes:**
- Each datasource should have exactly 3 unit tests following this pattern
- No `t.Parallel()` in datasource unit tests (we rely on parallelism at the test runner level)
- Replace `ResourceName` with the actual datasource name (e.g., `Aggregate`, `ASN`, etc.)
- Replace `netbox_resource_name` with the actual resource type name (e.g., `netbox_aggregate`)

---

## Key Differences from Resource Testing

### Resources vs DataSources

| Aspect | Resources | DataSources |
|--------|-----------|-------------|
| **Test Count** | 4 standard tests | 3 standard tests |
| **Test Types** | Schema, Metadata, Configure, Create Test | Schema, Metadata, Configure |
| **LookupAttrs** | N/A (resources have required/optional) | Required for datasources (search fields) |
| **ComputedAttrs** | All other fields | Result fields (typically larger set) |
| **t.Parallel()** | Yes, enabled | No, handled at runner level |

### Schema Validation Patterns

**Resources:**
```go
testutil.ValidateResourceSchema(t, resp.Schema.Attributes, testutil.SchemaValidation{
    Required: []string{"field1", "field2"},
    Optional: []string{"field3", "field4"},
    Computed: []string{"id"},
})
```

**DataSources:**
```go
testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
    LookupAttrs: []string{"id", "name"},  // Optional fields used for search
    ComputedAttrs: []string{"description", "tags", ...},  // Result fields
})
```

---

## Datasource Categories

Different datasources may have different numbers of lookup fields. Use this guide to determine what fields should be lookup vs computed:

### Category 1: Simple Lookups (ID-only or ID+Name)
Examples: CircuitType, DeviceRole, DeviceType, Manufacturer

**Typical LookupAttrs:**
```go
LookupAttrs: []string{"id", "name", "slug"},
```

**Typical ComputedAttrs:**
```go
ComputedAttrs: []string{"description", "custom_fields"},
```

### Category 2: Multiple Lookup Fields
Examples: Aggregate, ASN, Cable, Interface

**Typical LookupAttrs:**
```go
LookupAttrs: []string{"id", "name", "device", "name_upper", "type"},
```

**Typical ComputedAttrs:**
```go
ComputedAttrs: []string{"description", "tags", "comments"},
```

### Category 3: Complex Relationship Lookups
Examples: Device, VirtualMachine, Site

**Typical LookupAttrs:**
```go
LookupAttrs: []string{"id", "name", "site", "cluster", "device_type"},
```

**Typical ComputedAttrs:**
```go
ComputedAttrs: []string{"description", "device_count", "vm_count", "status"},
```

---

## Implementation Guidelines

### 1. Clean File Formatting

Ensure proper spacing to avoid spurious newlines:
- Single blank line between function definitions
- No extra blank lines at the start/end of function bodies
- Use standard Go formatting (gofmt)

**Before (Incorrect):**
```go
func TestSchema(t *testing.T) {

	d := datasources.NewTestDataSource()

	if d == nil {

		t.Fatal("nil")

	}

}
```

**After (Correct):**
```go
func TestSchema(t *testing.T) {
	d := datasources.NewTestDataSource()
	if d == nil {
		t.Fatal("nil")
	}
}
```

### 2. Import Organization

Always import in this order:
```go
import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)
```

### 3. Naming Conventions

Test function names should follow this pattern:
- `Test{DataSourceName}DataSourceSchema`
- `Test{DataSourceName}DataSourceMetadata`
- `Test{DataSourceName}DataSourceConfigure`

Examples:
- `TestAggregateDataSourceSchema`
- `TestCircuitDataSourceMetadata`
- `TestDeviceDataSourceConfigure`

---

## Running Tests

### All Datasource Unit Tests
```bash
go test ./internal/datasources_unit_tests/... -v
```

### Specific Datasource Tests
```bash
go test ./internal/datasources_unit_tests/... -v -run TestAggregate
```

### With Coverage
```bash
go test ./internal/datasources_unit_tests/... -v -cover
```

### Expected Output
All 3 tests per datasource should PASS:
- TestDataSourceNameSchema - PASS
- TestDataSourceNameMetadata - PASS
- TestDataSourceNameConfigure - PASS

---

## Checklist for Adding/Updating Datasource Unit Tests

- [ ] File located in `internal/datasources_unit_tests/`
- [ ] Package is `datasources_unit_tests`
- [ ] Imports include: `datasources`, `testutil`, `datasource`, `context`, `testing`
- [ ] All 3 standard tests present: Schema, Metadata, Configure
- [ ] Uses helper functions from `testutil` package
- [ ] No spurious blank lines (verified with gofmt)
- [ ] Correct datasource type name in Metadata test
- [ ] LookupAttrs and ComputedAttrs correctly identified
- [ ] No `t.Parallel()` calls
- [ ] All tests PASS when run: `go test ./internal/datasources_unit_tests/...`

---

## Example: Complete Datasource Unit Test

### File: `internal/datasources_unit_tests/aggregate_data_source_test.go`

```go
package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestAggregateDataSourceSchema(t *testing.T) {
	d := datasources.NewAggregateDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs: []string{"id", "prefix", "rir"},
		ComputedAttrs: []string{
			"rir_name",
			"tenant",
			"tenant_name",
			"date_added",
			"description",
			"comments",
			"tags",
			"custom_fields",
		},
	})
}

func TestAggregateDataSourceMetadata(t *testing.T) {
	d := datasources.NewAggregateDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_aggregate")
}

func TestAggregateDataSourceConfigure(t *testing.T) {
	d := datasources.NewAggregateDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
```

---

## Troubleshooting

### Test Fails: "Schema returned errors"
- Check that the datasource's Schema method is implemented
- Verify schema attributes are correctly defined
- Run `go build ./internal/datasources/...` to check for compilation errors

### Test Fails: "Expected type name X, got Y"
- Verify the datasource's Metadata method returns the correct TypeName
- Check that the expected type name matches the provider documentation

### Test Fails: "Expected error with incorrect provider data"
- Ensure the datasource implements DataSourceWithConfigure
- Verify the Configure method validates provider data type correctly

### File Has Spurious Newlines
- Run `gofmt -w` on the file
- Run `go mod tidy`
- Run pre-commit hooks: `.venv\Scripts\pre-commit.exe run --files <filename>`

---

## FAQ

**Q: Why only 3 tests for datasources vs 4 for resources?**

A: DataSources don't have Create/Read/Update/Delete operations like resources, so they need fewer tests. The 3 tests cover the essential aspects: schema definition, type identification, and configuration handling.

**Q: Can I add t.Parallel() to datasource unit tests?**

A: Not recommended. The test runner handles parallelism at the suite level. Adding t.Parallel() to individual tests is unnecessary and may cause issues with test isolation.

**Q: Should I test the Read operation?**

A: No, Read operations are tested via acceptance tests in `datasources_acceptance_tests/`. Unit tests focus on structural validation (schema, metadata, configuration).

**Q: What if my datasource has custom validation logic?**

A: Add additional unit tests to cover the custom logic, but keep the 3 standard tests. Name custom tests clearly, e.g., `TestAggregateDataSourceValidation_*.`

**Q: How do I know which fields are lookup vs computed?**

A: Lookup fields are those that users can filter/search by (typically optional in the HCL). Computed fields are those that are returned by the API (typically required fields in response or optional fields from the API).
