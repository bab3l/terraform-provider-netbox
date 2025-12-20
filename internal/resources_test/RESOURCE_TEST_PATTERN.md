# Resource Test Pattern Guide

This document outlines the standard pattern for writing resource unit tests in the Netbox Terraform provider. By following this pattern, we maintain consistency across all 100+ resources and reduce code duplication.

## Standard Resource Test Structure

Every resource test file should include these three unit test functions using the provided helper utilities.

### 1. Test Resource Creation (`TestXXXResourceXXX`)

Basic sanity check that the resource factory returns a non-nil resource.

```go
func TestASNRangeResource(t *testing.T) {
	t.Parallel()
	r := resources.NewASNRangeResource()

	if r == nil {
		t.Fatal("Expected non-nil ASNRange resource")
	}
}
```

### 2. Test Resource Schema (`TestXXXResourceSchema`)

Validates that the resource schema is correctly defined with all required, optional, and computed fields.

**Helper Function:** `testutil.ValidateResourceSchema()`

```go
func TestASNRangeResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewASNRangeResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name", "slug", "rir", "start", "end"},
		Optional: []string{"tenant", "description", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}
```

**SchemaValidation Fields:**
- `Required` - List of attributes that must be marked as required in the schema
- `Optional` - List of attributes that must be marked as optional (not required)
- `Computed` - List of attributes that must be marked as computed (read-only)
- `OptionalComputed` - List of attributes that are both optional and computed

### 3. Test Resource Metadata (`TestXXXResourceMetadata`)

Validates that the resource's metadata (type name) is correctly configured.

**Helper Function:** `testutil.ValidateResourceMetadata()`

```go
func TestASNRangeResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewASNRangeResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_asn_range")
}
```

**Parameters:**
- `r` - The resource instance to test
- `"netbox"` - The provider type name (typically "netbox")
- `"netbox_asn_range"` - The expected resource type name

### 4. Test Resource Configuration (`TestXXXResourceConfigure`)

Validates that the resource's Configure method properly handles provider data in three scenarios:
1. `nil` provider data (backwards compatibility)
2. Valid `*netbox.APIClient`
3. Invalid provider data (type assertion validation)

**Helper Function:** `testutil.ValidateResourceConfigure()`

```go
func TestASNRangeResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewASNRangeResource()
	testutil.ValidateResourceConfigure(t, r)
}
```

This single line replaces 26 lines of boilerplate and automatically:
- Validates the resource implements `ResourceWithConfigure`
- Tests Configure with nil provider data
- Tests Configure with a valid APIClient
- Tests Configure with invalid provider data using `testutil.InvalidProviderData`

## Helper Functions Reference

### ValidateResourceSchema()

```go
testutil.ValidateResourceSchema(t, schemaAttrs, validation)
```

Checks that schema attributes match the expected structure. Reports errors for:
- Missing required attributes
- Required attributes not marked as required
- Optional attributes incorrectly marked as required
- Missing computed attributes
- Computed attributes not marked as computed

### ValidateResourceMetadata()

```go
testutil.ValidateResourceMetadata(t, resource, providerTypeName, expectedTypeName)
```

Validates the resource type name. Reports error if the type name doesn't match the expected name (e.g., "netbox_asn_range").

### ValidateResourceConfigure()

```go
testutil.ValidateResourceConfigure(t, resource)
```

Tests the resource's Configure method with three provider data scenarios:
1. nil (backwards compatibility)
2. Valid `*netbox.APIClient`
3. Invalid data (type assertion error check)

## Complete Example

Here's a minimal but complete resource test file structure:

```go
package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestMyResource(t *testing.T) {
	t.Parallel()
	r := resources.NewMyResource()
	if r == nil {
		t.Fatal("Expected non-nil MyResource")
	}
}

func TestMyResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewMyResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}
	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name", "slug"},
		Optional: []string{"description"},
		Computed: []string{"id"},
	})
}

func TestMyResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewMyResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_my_resource")
}

func TestMyResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewMyResource()
	testutil.ValidateResourceConfigure(t, r)
}

// ... acceptance tests and other unit tests below ...
```

## Benefits

- **Consistency**: All resource tests follow the same pattern
- **Maintainability**: Changes to test patterns only need to be made in testutil helpers
- **Readability**: Test intent is immediately clear from helper function names
- **Brevity**: Each resource test file is much shorter and easier to review
- **Coverage**: Ensures all resources are properly tested for schema, metadata, and configuration

## Migration Guide

When updating existing resource test files:

1. Replace the schema validation loop with `testutil.ValidateResourceSchema()`
2. Replace the metadata assertion with `testutil.ValidateResourceMetadata()`
3. Replace the Configure test logic with `testutil.ValidateResourceConfigure()`

This typically reduces each test file by 50-100 lines of code while improving coverage.
