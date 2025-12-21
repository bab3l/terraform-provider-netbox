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

## Important Note: Parallelism Strategy

**Unit Tests (`resources_unit_tests/`)**: ✅ **USE `t.Parallel()`**
- Unit tests run in parallel for speed
- No external dependencies (no database, API, or shared resources)
- Each test is completely isolated
- Typical execution time for all unit tests: 3-5 seconds

**Acceptance Tests (`resources_acceptance_tests/`)**: ❌ **NO `t.Parallel()`, use `resource.Test()` not `resource.ParallelTest()`**
- Acceptance tests run sequentially
- All tests connect to the same shared Netbox database instance
- Parallel execution exhausts database connection pools and causes resource conflicts
- Database has limited connections (typically 20-50 depending on configuration)
- Cleanup functions depend on sequential execution to properly track resource state
- Typical execution time for each resource: 10-60 seconds (depends on resource complexity)

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

### 0. Cleanup: Fix File Layout and Remove Incomplete Comments

When refactoring test files, first clean up the file layout and remove incomplete comments:

#### File Layout Cleanup
- **Remove spurious blank lines**: Some test files have excessive blank lines between function definitions and within function bodies
- **Standardize spacing**: Use single blank lines between functions and logical blocks within functions
- **Fix formatting**: Ensure consistent indentation and line breaks (Go formatter may have added extra newlines during GitHub uploads)

**Example:**
Old file may have:
```go
func TestResource(t *testing.T) {

	t.Parallel()

	r := resources.NewResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}
```

Clean to:
```go
func TestResource(t *testing.T) {
	t.Parallel()

	r := resources.NewResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}
```

#### Comment Cleanup
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

### 5. Cleanup Functions in Acceptance Tests

**Purpose:**
Cleanup functions ensure that resources created during acceptance tests are properly deleted from the Netbox instance, preventing state pollution and test conflicts, especially important in consistency tests.

**Why Cleanup is Critical:**
- **State Pollution**: Without cleanup, resources accumulate in the Netbox database across test runs
- **Test Conflicts**: Hardcoded resource names can conflict when tests are re-run without cleanup
- **Parallel Test Issues**: Multiple concurrent tests can exhaust connection pools and cause race conditions
- **Consistency Tests**: Tests that verify reference attributes need proper cleanup to avoid cross-test dependencies

**Using testutil.NewCleanupResource():**

All acceptance tests should register cleanup for resources they create:

```go
func TestAccConsistency_ResourceName(t *testing.T) {
    resourceName := testutil.RandomName("resource")
    siteName := testutil.RandomName("site")
    siteSlug := testutil.RandomSlug("site")

    // Create cleanup object and register cleanup functions
    cleanup := testutil.NewCleanupResource(t)
    cleanup.RegisterResourceCleanup(resourceName)
    cleanup.RegisterSiteCleanup(siteSlug)

    resource.Test(t, resource.TestCase{
        PreCheck: func() { testutil.TestAccPreCheck(t) },
        ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccConsistencyConfig(resourceName, siteName, siteSlug),
                Check: resource.ComposeTestCheckFunc(...),
            },
        },
    })
}
```

**Available Cleanup Functions:**
The `testutil.CleanupResource` helper provides cleanup methods for all major resource types:

- Device/Type Related: `RegisterDeviceCleanup()`, `RegisterDeviceTypeCleanup()`, `RegisterDeviceRoleCleanup()`
- Network/Site Related: `RegisterSiteCleanup()`, `RegisterClusterCleanup()`, `RegisterClusterTypeCleanup()`, `RegisterClusterGroupCleanup()`
- Manufacturer/Vendor: `RegisterManufacturerCleanup()`
- Provider/Circuit: `RegisterProviderCleanup()`, `RegisterCircuitCleanup()`, `RegisterCircuitTypeCleanup()`
- Port Related: `RegisterConsolePortCleanup()`, `RegisterPowerPortCleanup()`
- Administrative: `RegisterRIRCleanup()`, `RegisterTenantCleanup()`, `RegisterASNRangeCleanup()`

**Cleanup Pattern for Consistency Tests:**

Consistency tests verify that reference attributes don't drift when re-applied. They create temporary resources and check state preservation:

```go
func TestAccConsistency_CircuitTermination_LiteralNames(t *testing.T) {
    providerName := testutil.RandomName("provider")
    providerSlug := testutil.RandomSlug("provider")
    circuitTypeName := testutil.RandomName("circuit-type")
    circuitTypeSlug := testutil.RandomSlug("circuit-type")
    circuitCid := testutil.RandomName("CID")
    siteName := testutil.RandomName("site")
    siteSlug := testutil.RandomSlug("site")

    cleanup := testutil.NewCleanupResource(t)
    cleanup.RegisterCircuitCleanup(circuitCid)  // The main resource

    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testutil.TestAccPreCheck(t) },
        ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccCircuitTerminationConsistencyLiteralNamesConfig(
                    providerName, providerSlug, circuitTypeName, circuitTypeSlug,
                    circuitCid, siteName, siteSlug),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("netbox_circuit_termination.test", "circuit", circuitCid),
                ),
            },
            {
                PlanOnly: true,
                Config: testAccCircuitTerminationConsistencyLiteralNamesConfig(...),
            },
        },
    })
}
```

**Key Points:**
- Register cleanup **before** the `resource.Test()` call
- Register cleanup for the **main resource** being tested and any **complex dependencies**
- Cleanup functions gracefully handle "already deleted" resources (expected behavior in Terraform tests)
- For simple resources with no dependencies (e.g., Contact), cleanup may not be necessary
- Cleanup is especially important for **consistency tests** which are run multiple times to check for drift

**Cleanup Error Handling:**
If you see logs like:
```
Cleanup: manufacturer with slug manufacturer-abc123 not found (already deleted)
```

This is **expected and correct**. Terraform's test framework automatically cleans up resources after each test step. The cleanup helper logs these gracefully rather than failing.

### 6. Deleted Files

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
   - **Add cleanup registration** (see Cleanup Functions section):
     ```go
     cleanup := testutil.NewCleanupResource(t)
     cleanup.Register<ResourceType>Cleanup(<resourceID>)
     ```

4. **Identify resource dependencies** and add appropriate cleanup
   - For consistency tests, register cleanup for the main resource being tested
   - For complex resources with multiple dependencies, register cleanup for each dependency
   - See the Available Cleanup Functions list for which functions to use

5. **Delete old file** from `internal/resources_test/`

6. **Test locally**:
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
- All acceptance tests should include cleanup registration to prevent state pollution
- Consistency tests are particularly important for resources with complex reference attributes
