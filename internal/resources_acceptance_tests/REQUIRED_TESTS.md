# Required Acceptance Tests for Resources

This document defines the standard acceptance test suite that every resource in the Terraform NetBox provider must implement. Following these standards ensures consistent quality, reliable behavior, and comprehensive coverage across all resources.

## Table of Contents

- [Test Naming Convention](#test-naming-convention)
- [Required Test Suite](#required-test-suite)
- [Tag Tests (Resources with Tags)](#tag-tests-resources-with-tags)
- [Test Helpers](#test-helpers)
- [Config Function Naming](#config-function-naming)
- [Cleanup Registration](#cleanup-registration)
- [Implementation Checklist](#implementation-checklist)

---

## Test Naming Convention

All test functions must follow these naming patterns:

```
TestAcc{Resource}Resource_{testType}
```

Where:
- `{Resource}` is the PascalCase resource name (e.g., `IPAddress`, `VirtualMachine`, `DeviceRole`)
- `{testType}` describes the test scenario in camelCase

**Examples:**
- `TestAccIPAddressResource_basic`
- `TestAccVirtualMachineResource_tagLifecycle`
- `TestAccDeviceRoleResource_externalDeletion`

---

## Required Test Suite

### Tier 1: Core CRUD Tests (Required for ALL resources)

| Test Name | Description | Priority |
|-----------|-------------|----------|
| `TestAcc{Resource}Resource_basic` | Create with minimal required fields only | **REQUIRED** |
| `TestAcc{Resource}Resource_full` | Create with all optional fields populated | **REQUIRED** |
| `TestAcc{Resource}Resource_update` | Modify an existing resource | **REQUIRED** |
| `TestAcc{Resource}Resource_import` | Terraform import functionality | **REQUIRED** |

### Tier 2: Reliability Tests (Required for ALL resources)

| Test Name | Description | Priority |
|-----------|-------------|----------|
| `TestAcc{Resource}Resource_externalDeletion` | Handle resource deleted outside Terraform | **REQUIRED** |
| `TestAcc{Resource}Resource_removeOptionalFields` | Clear optional fields by removing from config | **REQUIRED** |

### Tier 3: Consistency Tests (Recommended)

| Test Name | Description | Priority |
|-----------|-------------|----------|
| `TestAccConsistency_{Resource}_LiteralNames` | Literal values vs resource references | Recommended |
| `TestAcc{Resource}Resource_validationErrors` | Invalid input handling | Recommended |

---

## Tag Tests (Resources with Tags)

Resources that support tags MUST implement the following additional tests using the provided helper functions.

**IMPORTANT:** Use ONLY the helper functions - do not create manual/custom tag tests. Helper functions ensure consistency and reduce code duplication.

| Test Name | Description | Helper Function | Priority |
|-----------|-------------|-----------------|----------|
| `TestAcc{Resource}Resource_tagLifecycle` | Full tag lifecycle (add, modify, remove) | `RunTagLifecycleTest()` | **REQUIRED** |
| `TestAcc{Resource}Resource_tagOrderInvariance` | Tag reordering doesn't cause drift | `RunTagOrderTest()` | **REQUIRED** |

### Tag Test Implementation

**✅ CORRECT - Use helper functions:**

```go
// Option 1: Use RunTagLifecycleTest helper
func TestAccIPAddressResource_tagLifecycle(t *testing.T) {
    t.Parallel()

    // Setup test data...

    testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
        ResourceName:      "netbox_ip_address",
        ConfigWithoutTags: func() string { return configWithoutTags(address) },
        ConfigWithTags:    func() string { return configWithTags(address, tag1, tag2) },
        ConfigWithDifferentTags: func() string { return configWithDifferentTags(address, tag3) },
        ExpectedTagCount:          2,
        ExpectedDifferentTagCount: 1,
    })
}

// Option 2: Use RunTagOrderTest helper
func TestAccIPAddressResource_tagOrderInvariance(t *testing.T) {
    t.Parallel()

    // Setup test data...

    testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
        ResourceName:         "netbox_ip_address",
        ConfigWithTagsOrderA: func() string { return configOrderA(address) },
        ConfigWithTagsOrderB: func() string { return configOrderB(address) },
        ExpectedTagCount:     2,
    })
}
```

### Tag Lifecycle Test Coverage

The `tagLifecycle` test must cover these scenarios:
1. ✅ Create resource without tags (`tags = []`)
2. ✅ Add tags to existing resource
3. ✅ Modify tags (replace with different tags)
4. ✅ Remove all tags (set `tags = []`)
5. ✅ Verify no drift after tag removal

---

## Test Helpers

Use standardized helpers from `internal/testutil/` for consistent implementations:

### Core Test Helpers

| Helper | Purpose | File |
|--------|---------|------|
| `RunImportTest()` | Standard import testing | `import_tests.go` |
| `RunExternalDeletionTest()` | External deletion handling | `external_deletion_tests.go` |
| `TestRemoveOptionalFields()` | Optional field removal | `optional_fields_tests.go` |
| `RunMultiValidationErrorTest()` | Validation error testing | `validation_tests.go` |

### Tag Test Helpers

| Helper | Purpose | File |
|--------|---------|------|
| `RunTagLifecycleTest()` | Complete tag lifecycle | `tag_tests.go` |
| `RunTagOrderTest()` | Tag order invariance | `tag_tests.go` |

### Usage Example

```go
func TestAccIPAddressResource_removeOptionalFields(t *testing.T) {
    t.Parallel()

    address := fmt.Sprintf("198.51.100.%d/32", acctest.RandIntRange(1, 254))

    cleanup := testutil.NewCleanupResource(t)
    cleanup.RegisterIPAddressCleanup(address)

    testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
        ResourceName: "netbox_ip_address",
        BaseConfig: func() string {
            return testAccIPAddressResourceConfig_basic(address)
        },
        ConfigWithFields: func() string {
            return testAccIPAddressResourceConfig_withOptionalFields(address)
        },
        OptionalFields: map[string]string{
            "description": "Test description",
            "comments":    "Test comments",
        },
        RequiredFields: map[string]string{
            "address": address,
        },
    })
}
```

---

## Config Function Naming

Test configuration functions must follow this naming pattern:

```
testAcc{Resource}ResourceConfig_{variant}
```

### Standard Config Functions

| Function | Purpose |
|----------|---------|
| `testAcc{Resource}ResourceConfig_basic` | Minimal required fields |
| `testAcc{Resource}ResourceConfig_full` | All fields populated |
| `testAcc{Resource}ResourceConfig_update` | Modified values for update test |
| `testAcc{Resource}ResourceConfig_withTags` | Resource with tags |
| `testAcc{Resource}ResourceConfig_withoutTags` | Resource with `tags = []` |
| `testAcc{Resource}ResourceConfig_tagsOrderA` | Tags in order A |
| `testAcc{Resource}ResourceConfig_tagsOrderB` | Tags in order B (different order) |

---

## Cleanup Registration

**Every test MUST register cleanup** to ensure resources are removed even if tests fail:

```go
func TestAccIPAddressResource_basic(t *testing.T) {
    t.Parallel()

    address := fmt.Sprintf("192.168.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

    // REQUIRED: Register cleanup
    cleanup := testutil.NewCleanupResource(t)
    cleanup.RegisterIPAddressCleanup(address)

    resource.Test(t, resource.TestCase{
        // ... test configuration
    })
}
```

### Cleanup Methods

Each resource type has a corresponding cleanup registration method:

```go
cleanup.RegisterIPAddressCleanup(address)
cleanup.RegisterVirtualMachineCleanup(name)
cleanup.RegisterDeviceCleanup(name)
cleanup.RegisterTagCleanup(slug)
cleanup.RegisterTenantCleanup(slug)
// ... etc
```

---

## Implementation Checklist

Use this checklist when implementing tests for a new resource or auditing existing tests:

### Basic Requirements

- [ ] `TestAcc{Resource}Resource_basic` - Creates with required fields only
- [ ] `TestAcc{Resource}Resource_full` - Creates with all optional fields
- [ ] `TestAcc{Resource}Resource_update` - Modifies existing resource
- [ ] `TestAcc{Resource}Resource_import` - Imports existing resource

### Reliability Requirements

- [ ] `TestAcc{Resource}Resource_externalDeletion` - Handles external deletion
- [ ] `TestAcc{Resource}Resource_removeOptionalFields` - Clears optional fields

### Tag Requirements (if resource supports tags)

- [ ] `TestAcc{Resource}Resource_tagLifecycle` - Full tag lifecycle
- [ ] `TestAcc{Resource}Resource_tagOrderInvariance` - Order doesn't cause drift

### Code Quality

- [ ] All tests call `t.Parallel()` for concurrent execution
- [ ] Cleanup registered for all created resources
- [ ] Random names/values used to avoid conflicts
- [ ] Config functions follow naming convention
- [ ] Uses test helpers where appropriate

---

## File Organization

Each resource should have its tests in a single file:

```
internal/resources_acceptance_tests/
├── {resource}_resource_test.go      # Main tests
├── {resource}_resource_ext_test.go  # Extended/edge case tests (optional)
└── REQUIRED_TESTS.md                # This documentation
```

### Test File Structure

```go
package resources_acceptance_tests

import (
    // Standard imports...
)

// =============================================================================
// TIER 1: CORE CRUD TESTS
// =============================================================================

func TestAcc{Resource}Resource_basic(t *testing.T) { ... }
func TestAcc{Resource}Resource_full(t *testing.T) { ... }
func TestAcc{Resource}Resource_update(t *testing.T) { ... }
func TestAcc{Resource}Resource_import(t *testing.T) { ... }

// =============================================================================
// TIER 2: RELIABILITY TESTS
// =============================================================================

func TestAcc{Resource}Resource_externalDeletion(t *testing.T) { ... }
func TestAcc{Resource}Resource_removeOptionalFields(t *testing.T) { ... }

// =============================================================================
// TIER 3: TAG TESTS (if applicable)
// =============================================================================

func TestAcc{Resource}Resource_tagLifecycle(t *testing.T) { ... }
func TestAcc{Resource}Resource_tagOrderInvariance(t *testing.T) { ... }

// =============================================================================
// TIER 4: CONSISTENCY & VALIDATION TESTS
// =============================================================================

func TestAccConsistency_{Resource}_LiteralNames(t *testing.T) { ... }
func TestAcc{Resource}Resource_validationErrors(t *testing.T) { ... }

// =============================================================================
// CONFIG HELPER FUNCTIONS
// =============================================================================

func testAcc{Resource}ResourceConfig_basic(...) string { ... }
func testAcc{Resource}ResourceConfig_full(...) string { ... }
// ... etc
```
