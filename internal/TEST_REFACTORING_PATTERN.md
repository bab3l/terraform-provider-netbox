# Test Refactoring Pattern

## Overview
This document describes the systematic refactoring pattern applied to split resource tests from `internal/resources_test/` into two specialized directories with improved code organization and reduced duplication.

## Important Note: Acceptance Test Coverage Varies by Resource Complexity

Not all resources require the same number of acceptance tests. The number of tests should match the resource's complexity and use cases:

- **Simple resources** (e.g., cable): 1-2 acceptance tests
  - Only required fields are complex (e.g., terminations for cables)
  - Other optional fields are simple scalars
  - No complex reference validation or state drift issues
  - Basic + import test is sufficient

- **Medium resources** (e.g., ASN): 2-3 acceptance tests
  - Some optional fields with validation logic
  - May have basic, full, and consistency tests

- **Complex resources** (e.g., Aggregate, ASN Range): 5-6 acceptance tests
  - Multiple complex reference attributes
  - Needs basic, full, update, import, and multiple consistency tests
  - May require testing name-to-ID resolution without drift

When refactoring tests, match test coverage to the resource schema:
- If the schema indicates missing test coverage is **appropriate** (complex optional fields, reference validation, state drift concerns), **add additional tests** to cover those scenarios
- If the schema indicates the resource is **simpler** (only simple scalar optional fields, no complex references), **preserve the fewer number of tests** and add a file-level comment explaining why (see cable_resource_test.go example)

## Directory Structure

### Before
```
internal/resources_test/
  ├── resource_name_resource_test.go  (contains both unit and acceptance tests)
  └── acceptance_test.go              (shared constants)
```

### After
```
internal/resources_unit_tests/
  └── resource_name_resource_test.go  (unit tests only)

internal/resources_acceptance_tests/
  └── resource_name_resource_test.go  (acceptance tests only)

internal/testutil/
  ├── schema_validation.go            (3 helper functions)
  ├── test_constants.go               (shared test constants)
  └── invalid_provider_data.go        (shared provider data variable)
```

## Changes Applied

### 0. Cleanup: Remove Incomplete or Irrelevant Comments

When refactoring test files, review and clean up any incomplete or irrelevant comments:

- **Incomplete test function stubs**: If a test function is declared (e.g., `// TestAccConsistency_Cable_LiteralNames ...`) but not implemented (no test function body), remove the incomplete comment entirely
- **Dangling TODO comments**: Remove comments about tests that were never implemented
- **Outdated comments**: Remove comments that no longer reflect the actual test structure or intent

**Example:**
Old file may have:
```go
// TestAccConsistency_Cable_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
// (but no actual test function follows)
```

Delete this comment during refactoring and add a clear explanation (at the file level) for why the resource has fewer acceptance tests.

### 1. Unit Tests → `internal/resources_unit_tests/`

**File Structure:**
- Package: `resources_unit_tests`
- Functions: 4 standard unit tests
- Imports: `resources`, `testutil`, `fwresource`

**Standard Unit Tests:**
```go
// 1. TestResourceName - validates resource is not nil
func TestResourceName(t *testing.T) {
    t.Parallel()
    r := resources.NewResourceNameResource()
    if r == nil {
        t.Fatal("Expected non-nil resource")
    }
}

// 2. TestResourceNameSchema - validates schema attributes using helper
func TestResourceNameSchema(t *testing.T) {
    t.Parallel()
    r := resources.NewResourceNameResource()
    schemaRequest := fwresource.SchemaRequest{}
    schemaResponse := &fwresource.SchemaResponse{}
    r.Schema(context.Background(), schemaRequest, schemaResponse)

    // ... error checks ...

    testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
        Required: []string{...},
        Optional: []string{...},
        Computed: []string{...},
    })
}

// 3. TestResourceNameMetadata - validates metadata using helper
func TestResourceNameMetadata(t *testing.T) {
    t.Parallel()
    r := resources.NewResourceNameResource()
    testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_resource_name")
}

// 4. TestResourceNameConfigure - validates configure using helper
func TestResourceNameConfigure(t *testing.T) {
    t.Parallel()
    r := resources.NewResourceNameResource()
    testutil.ValidateResourceConfigure(t, r)
}
```

**Key Points:**
- Uses 3 helper functions from `testutil` to reduce boilerplate
- No provider dependencies (no circular imports)
- Runs instantly without external infrastructure

### 2. Acceptance Tests → `internal/resources_acceptance_tests/`

**File Structure:**
- Package: `resources_acceptance_tests`
- Imports: `provider`, `testutil`, terraform testing framework
- May contain 1+ acceptance tests depending on resource complexity

**Standard Acceptance Tests:**
- `TestAccResourceName_basic` - minimal configuration
- `TestAccResourceName_full` - all optional fields
- `TestAccResourceName_update` - configuration changes (if applicable)
- `TestAccConsistency_ResourceName` - reference attribute consistency (if applicable)
- `TestAccResourceName_import` - state import verification (if applicable)

**Key Points:**
- Uses shared constants from `testutil.Comments`, `testutil.Description1`, etc.
- Connects to running Netbox instance for full integration testing
- Tests CRUD operations and state management

### 3. Helper Functions in `internal/testutil/`

**ValidateResourceSchema()**
```go
func ValidateResourceSchema(t *testing.T, attrs map[string]attr.Attribute, validation SchemaValidation)
```
- Validates required attributes exist and are marked required
- Validates optional attributes exist and are marked optional
- Validates computed attributes exist and are marked computed
- Replaces ~20 lines of manual validation code

**ValidateResourceMetadata()**
```go
func ValidateResourceMetadata(t *testing.T, resource fwresource.Resource, providerTypeName, expectedTypeName string)
```
- Validates resource type name matches expected format
- Replaces ~10 lines of manual validation code

**ValidateResourceConfigure()**
```go
func ValidateResourceConfigure(t *testing.T, resource fwresource.Resource)
```
- Tests nil provider data (no error expected)
- Tests valid API client (no error expected)
- Tests invalid provider data (error expected)
- Replaces ~20 lines of manual validation code

### 4. Shared Constants in `internal/testutil/test_constants.go`

Moved from `internal/resources_test/acceptance_test.go`:
- `Comments` - "Test comments"
- `Description1` - "Initial description"
- `Description2` - "Updated description"
- `RearPortName` - "rear0"
- `Color` - "aa1409"
- `InterfaceName` - "eth0"

Usage in tests:
```go
Config: testAccResourceConfig_full(name, slug, testutil.Comments),
```

### 5. Deleted Files

- `internal/resources_test/resource_name_resource_test.go` (old combined file)
- `internal/resources_test/acceptance_test.go` (constants moved to testutil)

## Files Refactored

1. ✅ `aggregate_resource_test.go` - 5 unit tests, 5 acceptance tests
2. ✅ `asn_range_resource_test.go` - 4 unit tests, 6 acceptance tests
3. ✅ `asn_resource_test.go` - 4 unit tests, 3 acceptance tests
4. ✅ `cable_resource_test.go` - 4 unit tests, 1 acceptance test (simpler resource)
5. ⏳ `circuit_group_assignment_resource_test.go` - pending
6. ⏳ `circuit_group_resource_test.go` - pending
7. ⏳ ... 90+ remaining resource tests

## Benefits

- **Clear Separation**: Unit tests and acceptance tests are now in separate directories
- **No Circular Dependencies**: Unit tests don't import from provider, avoiding import cycles
- **Faster Development**: Unit tests run instantly without Netbox dependency
- **Code Reuse**: 3 helper functions eliminate ~50 lines of boilerplate per resource
- **Consistency**: All resources follow the same test structure
- **Maintainability**: Centralized constants and helpers make changes easier

## How to Apply to New Resources

When refactoring a new resource test file:

1. **Review and clean up comments**
   - Remove any incomplete or stub test comments that have no implementation
   - Delete dangling TODO or outdated comments
   - If acceptance test count is less than other resources, add a file-level comment explaining why (see cable_resource_test.go example)

2. **Extract unit tests** to `internal/resources_unit_tests/resource_name_resource_test.go`
   - Copy the 4 standard unit test functions
   - Update resource names/types
   - Replace manual validation with helper functions

3. **Extract acceptance tests** to `internal/resources_acceptance_tests/resource_name_resource_test.go`
   - Copy all `TestAcc*` functions (do not artificially add tests)
   - Replace local constants with `testutil.*` constants
   - Keep all Terraform configuration helpers

4. **Delete old file** from `internal/resources_test/`

5. **Test locally**:
   ```bash
   # Unit tests (fast)
   go test ./internal/resources_unit_tests -v -run TestResourceName

   # Acceptance tests (requires Netbox)
   TF_ACC=1 go test ./internal/resources_acceptance_tests -v -run TestAccResourceName
   ```

## Notes

- Test coverage varies by resource complexity (e.g., cable: 1 test, aggregate: 5 tests)
- All resources should have the 4 standard unit tests
- Acceptance test count depends on resource's use cases and reference fields
