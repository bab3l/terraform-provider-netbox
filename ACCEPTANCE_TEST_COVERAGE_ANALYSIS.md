# Acceptance Test Coverage Analysis

This document provides a comprehensive analysis of the acceptance test coverage for terraform-provider-netbox resources, identifying gaps and recommending new test classes.

## Executive Summary

Based on analysis of 97 resource types, the following key gaps have been identified:

| Gap Category | Count | Priority | Status |
|--------------|-------|----------|--------|
| Missing Import Tests | 0 | High | ✅ COMPLETE (100% coverage) |
| Missing Update Tests | 6 | High | 8 completed |
| Missing externalDeletion Tests | 22 | Medium | 2 completed ✅ |
| Missing removeOptionalFields Tests | 14 | Medium | All resources now covered |
| Missing Full Tests | 5 | Medium | 2 completed ✅ |
| Missing Consistency/LiteralNames Tests | ~70 | Low | - |
| Critically Under-tested Resources | 0 | Critical | ✅ RESOLVED |

**Latest Update (2026-01-15):**
- EventRule and NotificationGroup upgraded from critically under-tested to full coverage with 8 comprehensive tests each.
- **Import Test Coverage: 100% COMPLETE** - All 97 resources now have import testing (either embedded in _basic tests or as dedicated _import test functions). VirtualMachine was the last resource to receive import testing.
- **Update Test Coverage: Progress** - Added dedicated _update tests for Device, FrontPortTemplate, Interface, Prefix, PowerFeed, and Role. 7 resources remaining.

---

## Current Test Categories

The provider currently has these standard test types:

1. **`_basic`** - Minimal configuration testing (required fields only)
2. **`_full`** - Complete configuration testing (all optional fields populated)
3. **`_update`** - Test modifying existing resources
4. **`_import`** - Test terraform import functionality
5. **`_IDPreservation`** - Test that resource ID remains stable
6. **`_externalDeletion`** - Test handling of resources deleted outside Terraform
7. **`_removeOptionalFields`** - Test clearing optional fields
8. **`Consistency_*_LiteralNames`** - Test using literal values vs resource references

---

## Critical Gaps

### 1. Critically Under-tested Resources

~~All resources now have adequate basic test coverage.~~

**UPDATE (2026-01-15):** EventRule and NotificationGroup have been upgraded to full test coverage:
- ✅ EventRule: Now has 8 tests (basic, full, update, import, IDPreservation, externalDeletion, removeOptionalFields, removeOptionalFields_extended)
- ✅ NotificationGroup: Now has 8 tests (basic, full, update, import, IDPreservation, externalDeletion, removeOptionalFields, removeOptionalFields_extended)

Both resources now have CheckDestroy functions implemented and all tests passing.

### 2. Import Test Coverage Status ✅ COMPLETE

**Import testing is now at 100% coverage!** All 97 resources have import testing.

**Implementation Approach:**
- **52 resources** use embedded import steps in their `_basic` tests
- **45 resources** have dedicated `_import` test functions
- Both approaches are valid and provide complete import verification

**Note:** The original gap analysis was counting only dedicated `_import` test functions and missed the 52 resources with embedded import testing. The actual import coverage was always much higher than initially documented.

**Recent Additions:**
- ✅ VirtualMachine: Added import test (2026-01-15)
- ✅ EventRule: Has embedded import in _basic test
- ✅ NotificationGroup: Has embedded import in _basic test

**Resources with Embedded Import Testing (in `_basic` tests):**
Aggregate, ASN, ClusterGroup, ConfigContext, ConfigTemplate, ConsolePort, ConsolePortTemplate, ConsoleServerPort, ConsoleServerPortTemplate, Contact, ContactAssignment, ContactGroup, ContactRole, CustomLink, Device, DeviceBay, DeviceBayTemplate, DeviceRole, DeviceType, ExportTemplate, FHRPGroupAssignment, FrontPort, FrontPortTemplate, Interface, InterfaceTemplate, L2VPN, L2VPNTermination, ModuleBay, ModuleBayTemplate, Module, ModuleType, PowerFeed, PowerOutlet, PowerOutletTemplate, PowerPanel, PowerPort, PowerPortTemplate, RackReservation, RackType, RearPort, RearPortTemplate, RIR, Role, Service, ServiceTemplate, Tag, VirtualChassis, VirtualDeviceContext, VirtualMachine, Webhook, WirelessLAN, WirelessLANGroup, WirelessLink

### 3. Resources Missing Update Tests (14 resources)

Update operations are core CRUD functionality. Missing for:

**Recently Completed (2026-01-15):**
- ✅ Device: Added dedicated _update test
- ✅ EventRule: Completed
- ✅ FrontPortTemplate: Added dedicated _update test
- ✅ Interface: Added dedicated _update test
- ✅ NotificationGroup: Completed
- ✅ PowerFeed: Added dedicated _update test
- ✅ Prefix: Added dedicated _update test
- ✅ Role: Added dedicated _update test

**Still Pending (7 resources):**
- InterfaceTemplate
- PowerPanel
- PowerPortTemplate
- RearPort
- RearPortTemplate
- Tag
- VirtualChassis

---

## Recommended New Test Classes

### Class 1: Negative/Validation Tests

**Purpose:** Verify proper error handling for invalid inputs

```go
func TestAcc{Resource}Resource_invalidInput(t *testing.T)
func TestAcc{Resource}Resource_validationError(t *testing.T)
```

**Test Scenarios:**
- Invalid enum values (e.g., invalid status)
- Invalid format values (e.g., malformed IP addresses, URLs)
- Missing required fields
- Invalid field combinations
- Invalid reference IDs

**Example:**
```go
func TestAccIPAddressResource_invalidAddress(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Steps: []resource.TestStep{
            {
                Config:      testAccIPAddressResourceConfig_invalidFormat("not-an-ip"),
                ExpectError: regexp.MustCompile(`invalid IP address`),
            },
        },
    })
}
```

### Class 2: Concurrent Access Tests

**Purpose:** Verify behavior under concurrent modifications

```go
func TestAcc{Resource}Resource_concurrent(t *testing.T)
```

**Test Scenarios:**
- Multiple Terraform applies for related resources
- Race condition detection
- Optimistic locking behavior

### Class 3: Reference Attribute Tests

**Purpose:** Test resources that reference other resources

```go
func TestAcc{Resource}Resource_referenceUpdate(t *testing.T)
func TestAcc{Resource}Resource_referenceRemoval(t *testing.T)
func TestAcc{Resource}Resource_cascadeDelete(t *testing.T)
```

**Test Scenarios:**
- Changing a reference (e.g., moving device to different site)
- Removing optional references
- Behavior when referenced resource is deleted
- Circular reference handling

### Class 4: Hierarchical Resource Tests

**Purpose:** Test resources with parent-child relationships

```go
func TestAcc{Resource}Resource_parentChange(t *testing.T)
func TestAcc{Resource}Resource_nestedHierarchy(t *testing.T)
```

**Applicable Resources:**
- Region (can have parent region)
- SiteGroup (can have parent site group)
- ContactGroup (can have parent contact group)
- TenantGroup (can have parent tenant group)
- Location (can have parent location)
- WirelessLANGroup (can have parent wireless LAN group)
- VLANGroup (scope hierarchy)

**Test Scenarios:**
- Moving resource to different parent
- Orphaning a child (removing parent)
- Deep nesting limits
- Circular parent prevention

### Class 5: Computed Field Verification Tests

**Purpose:** Verify computed fields are populated correctly

```go
func TestAcc{Resource}Resource_computedFields(t *testing.T)
```

**Test Scenarios:**
- URL field generation
- Display name computation
- Created/last_updated timestamps
- Count fields (device_count, prefix_count, etc.)

### Class 6: Idempotency Tests

**Purpose:** Verify apply-plan-apply produces consistent results

```go
func TestAcc{Resource}Resource_idempotent(t *testing.T)
```

**Test Scenarios:**
- Multiple applies with same config (no changes)
- Apply after refresh with no config changes
- State consistency after external read

### Class 7: Large Value Tests

**Purpose:** Test behavior with large data sets

```go
func TestAcc{Resource}Resource_largeDescription(t *testing.T)
func TestAcc{Resource}Resource_manyTags(t *testing.T)
```

**Test Scenarios:**
- Very long strings (description, comments)
- Many tags (pagination)
- Many custom fields
- Large JSON data (config contexts)

### Class 8: Special Character Tests

**Purpose:** Test handling of special characters

```go
func TestAcc{Resource}Resource_specialCharacters(t *testing.T)
```

**Test Scenarios:**
- Unicode characters in names/descriptions
- Special characters that need escaping
- Newlines in multi-line fields
- Empty strings vs null

### Class 9: API Response Variation Tests

**Purpose:** Test handling of API response edge cases

```go
func TestAcc{Resource}Resource_partialResponse(t *testing.T)
func TestAcc{Resource}Resource_nestedNullHandling(t *testing.T)
```

**Test Scenarios:**
- Optional nested objects being null
- Empty arrays vs null arrays
- Deeply nested optional fields

### Class 10: Tags/CustomFields Lifecycle Tests

**Purpose:** Test tag and custom field manipulation

```go
func TestAcc{Resource}Resource_tagAddRemove(t *testing.T)
func TestAcc{Resource}Resource_tagReorder(t *testing.T)
func TestAcc{Resource}Resource_customFieldTypeChange(t *testing.T)
```

**Test Scenarios:**
- Adding tags to existing resource
- Removing all tags
- Reordering tags (should not cause drift)
- Different custom field types

---

## Test Coverage Priority Matrix

### Priority 1 (Critical - Must Have)

| Test Type | Reason |
|-----------|--------|
| `_basic` | Validates core creation |
| `_update` | Validates core CRUD |
| `_import` | Critical for production adoption |
| `_removeOptionalFields` | Prevents state inconsistency bugs |

### Priority 2 (High - Should Have)

| Test Type | Reason |
|-----------|--------|
| `_full` | Validates all fields work |
| `_externalDeletion` | Production reliability |
| `_IDPreservation` | State management correctness |
| `_validationError` | User experience |

### Priority 3 (Medium - Nice to Have)

| Test Type | Reason |
|-----------|--------|
| `Consistency_*_LiteralNames` | Reference handling |
| `_referenceUpdate` | Complex scenarios |
| `_hierarchical` | Parent-child relationships |

### Priority 4 (Low - Enhancement)

| Test Type | Reason |
|-----------|--------|
| `_concurrent` | Edge cases |
| `_largeValue` | Stress testing |
| `_specialCharacters` | Edge cases |

---

## Implementation Recommendations

### Phase 1: Critical Gaps ✅ COMPLETE

1. ✅ Add complete test suites for EventRule and NotificationGroup
2. ✅ Import test coverage now at 100%
3. **Next:** Add `_update` tests for 14 missing resources

### Phase 2: Update Test Coverage (1-2 weeks)

1. ✅ Add `_update` tests for Device, FrontPortTemplate, Interface (completed 2026-01-15)
2. ✅ Add `_update` tests for Prefix, PowerFeed, Role (completed 2026-01-15)
3. **In Progress:** Add `_update` tests for remaining 7 resources
4. Next batch: InterfaceTemplate, PowerPanel, PowerPortTemplate

### Phase 3: Remove Optional Fields (1 week)

1. Add `_removeOptionalFields` tests for 14 missing resources
2. Use the existing `testutil.TestRemoveOptionalFields` helper

### Phase 4: External Deletion (1-2 weeks)

1. Add `_externalDeletion` tests for 24 missing resources
2. Test recreation behavior after external deletion

### Phase 5: New Test Classes (Ongoing)

1. Implement validation error tests for resources with complex schemas
2. Add hierarchical tests for tree-structured resources
3. Add reference update tests for resources with many FK relationships

---

## Test Helper Functions (Available)

The following test helper functions are now available in the `internal/testutil` package:

### Import Tests (`import_tests.go`)
```go
// RunImportTest - Full import test with custom ID extraction
func RunImportTest(t *testing.T, config ImportTestConfig)

// RunSimpleImportTest - Basic import test for straightforward imports
func RunSimpleImportTest(t *testing.T, config SimpleImportTestConfig)
```

### Validation Tests (`validation_tests.go`)
```go
// RunValidationErrorTest - Test that invalid configs produce expected errors
func RunValidationErrorTest(t *testing.T, config ValidationErrorTestConfig)

// RunMultiValidationErrorTest - Run multiple validation tests as subtests
func RunMultiValidationErrorTest(t *testing.T, config MultiValidationErrorTestConfig)

// Pre-defined error patterns:
// - ErrPatternRequired, ErrPatternInvalidValue, ErrPatternInvalidFormat
// - ErrPatternInvalidIP, ErrPatternInvalidURL, ErrPatternInvalidEnum
// - ErrPatternNotFound, ErrPatternConflict, ErrPatternRange
```

### Hierarchical Tests (`hierarchical_tests.go`)
```go
// RunHierarchicalTest - Test parent-child relationships
func RunHierarchicalTest(t *testing.T, config HierarchicalTestConfig)

// RunNestedHierarchyTest - Test deeply nested hierarchies
func RunNestedHierarchyTest(t *testing.T, config NestedHierarchyTestConfig)
```

### Reference Tests (`reference_tests.go`)
```go
// RunReferenceChangeTest - Test changing FK references
func RunReferenceChangeTest(t *testing.T, config ReferenceChangeTestConfig)

// RunMultiReferenceTest - Test resources with multiple references
func RunMultiReferenceTest(t *testing.T, config MultiReferenceTestConfig)
```

### Idempotency Tests (`idempotency_tests.go`)
```go
// RunIdempotencyTest - Multiple applies produce no changes
func RunIdempotencyTest(t *testing.T, config IdempotencyTestConfig)

// RunRefreshIdempotencyTest - Refresh + plan shows no changes
func RunRefreshIdempotencyTest(t *testing.T, config RefreshIdempotencyTestConfig)
```

### Tag Tests (`tag_tests.go`)
```go
// RunTagLifecycleTest - Test complete tag lifecycle (add/change/remove)
func RunTagLifecycleTest(t *testing.T, config TagLifecycleTestConfig)

// RunTagOrderTest - Verify tag reordering doesn't cause drift
func RunTagOrderTest(t *testing.T, config TagOrderTestConfig)
```

### External Deletion Tests (`external_deletion_tests.go`)
```go
// RunExternalDeletionTest - Detect and recreate externally deleted resources
func RunExternalDeletionTest(t *testing.T, config ExternalDeletionTestConfig)

// RunExternalDeletionWithIDTest - Variant with custom ID extraction
func RunExternalDeletionWithIDTest(t *testing.T, config ExternalDeletionWithIDTestConfig)
```

### Update Tests (`update_tests.go`)
```go
// RunUpdateTest - Basic update test
func RunUpdateTest(t *testing.T, config UpdateTestConfig)

// RunMultiStepUpdateTest - Sequential updates
func RunMultiStepUpdateTest(t *testing.T, config MultiStepUpdateTestConfig)

// RunFieldUpdateTest - Update specific field
func RunFieldUpdateTest(t *testing.T, config FieldUpdateTestConfig)
```

### Edge Case Tests (`edge_case_tests.go`)
```go
// RunLargeValueTest - Test large string values
func RunLargeValueTest(t *testing.T, config LargeValueTestConfig)

// RunSpecialCharacterTests - Test special characters, unicode, etc.
func RunSpecialCharacterTests(t *testing.T, config SpecialCharacterTestConfig)

// RunEmptyStringTest - Test empty string vs null handling
func RunEmptyStringTest(t *testing.T, config EmptyStringTestConfig)

// CommonSpecialCharacterValues - Pre-defined test cases for special chars
// GenerateLargeString(length int) - Generate test strings of any length
// GenerateLargeDescription(paragraphs int) - Generate realistic descriptions
```

### Existing Helpers (`optional_field_tests.go`, `optional_computed_field_tests.go`)
```go
// RunOptionalFieldTestSuite - Test optional field lifecycle
func RunOptionalFieldTestSuite(t *testing.T, config OptionalFieldTestConfig)

// RunOptionalComputedFieldTestSuite - Test Optional+Computed fields
func RunOptionalComputedFieldTestSuite(t *testing.T, config OptionalComputedFieldTestConfig)

// TestRemoveOptionalFields - Test removing multiple optional fields
func TestRemoveOptionalFields(t *testing.T, config MultiFieldOptionalTestConfig)
```

---

## Summary Statistics

| Metric | Current | Target |
|--------|---------|--------|
| Resources with full test coverage | ~20 | 97 |
| Resources with import tests | 97 ✅ | 97 |
| Resources with update tests | 91 | 97 |
| Resources with removeOptionalFields tests | 83 | 97 |
| Resources with externalDeletion tests | 73 | 97 |
| New test classes needed | 0 | 10 |

---

*Generated: 2026-01-15*
*Provider Version: terraform-provider-netbox*
