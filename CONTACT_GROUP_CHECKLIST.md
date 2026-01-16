# Contact Group Resource Test Standardization Checklist

## Resource Information
- **Resource Name**: Contact Group
- **Test File**: `internal/resources_acceptance_tests/contact_group_resource_test.go`
- **Completion Date**: 2025-01-16
- **Status**: ✅ COMPLETE

## Test Structure Analysis

### Original Test Count
- **Before Standardization**: 7 tests
- **After Standardization**: 8 tests

### Test Categories

#### Core CRUD Tests (3)
- ✅ `TestAccContactGroupResource_basic` - Basic resource creation with name and slug
- ✅ `TestAccContactGroupResource_full` - Complete resource with all attributes including parent and tags
- ✅ `TestAccContactGroupResource_update` - Update name and description

#### Reliability Tests (2)
- ✅ `TestAccContactGroupResource_externalDeletion` - Handles external resource deletion
- ✅ `TestAccContactGroupResource_removeOptionalFields` - Handles removal of optional fields (description, parent)

#### Tag Tests (2) - **⚠️ USES NESTED FORMAT - FLAGGED FOR PHASE 2**
- ✅ `TestAccContactGroupResource_tagLifecycle` - Complete tag lifecycle (add, update, remove) using helper
- ✅ `TestAccContactGroupResource_tagOrderInvariance` - Tag order doesn't affect state using helper

#### Validation Tests (1)
- ✅ `TestAccContactGroupResource_validationErrors` - API validation errors (2 subtests: missing name, missing slug)

#### Special Tests (1)
- ✅ `TestAccConsistency_ContactGroup_LiteralNames` - Plan consistency with literal names

## Changes Made

### 1. Removed IDPreservation Test ✅
- **Removed**: `TestAccContactGroupResource_IDPreservation`
- **Reason**: Duplicate of basic test functionality
- **Lines Removed**: ~25 lines

### 2. Added Standard Tag Tests ✅
- **Added**: `TestAccContactGroupResource_tagLifecycle` using `testutil.RunTagLifecycleTest`
- **Added**: `TestAccContactGroupResource_tagOrderInvariance` using `testutil.RunTagOrderTest`
- **Added**: `testAccContactGroupResourceConfig_tagLifecycle` (config helper for lifecycle test)
- **Added**: `testAccContactGroupResourceConfig_tagOrder` (config helper for order test)
- **Tag Format**: Uses nested format `{name = ..., slug = ...}` - **⚠️ FLAGGED FOR PHASE 2 CONVERSION**
- **Note**: Full test already included tag checks, so we added complementary lifecycle and order tests

### 3. Added Missing Import ✅
- **Added**: `"strings"` import for tag helper logic

## Tag Format Issue ⚠️

**IMPORTANT**: This resource uses the **nested tag format**:
```hcl
tags = [
  { name = netbox_tag.tag1.name, slug = netbox_tag.tag1.slug },
  { name = netbox_tag.tag2.name, slug = netbox_tag.tag2.slug }
]
```

This format has been flagged for Phase 2 conversion to the simpler **slug list format**:
```hcl
tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
```

See COVERAGE_ANALYSIS.md Phase 2 section for conversion plan.

## Test Execution Results

### All Tests Passing ✅
```
TestAccContactGroupResource_full                       PASS (2.09s)
TestAccContactGroupResource_basic                      PASS (2.45s)
TestAccContactGroupResource_update                     PASS (2.06s)
TestAccContactGroupResource_tagLifecycle               PASS (6.24s)
TestAccContactGroupResource_tagOrderInvariance         PASS (2.81s)
TestAccContactGroupResource_externalDeletion           PASS (1.68s)
TestAccContactGroupResource_removeOptionalFields       PASS (2.32s)
TestAccContactGroupResource_validationErrors           PASS (0.55s)

Total: 8 tests, ~7.0 seconds
```

## Technical Details

### Dependencies
- Contact Group resource (for parent relationship)
- Tag resources (for tag tests)

### Key Attributes
- `name`: Contact group name (required)
- `slug`: URL-safe identifier (required)
- `description`: Description text (optional)
- `parent`: Parent contact group ID (optional, hierarchical structure)
- `tags`: Nested format tag list (⚠️ Phase 2 conversion needed)

### Resource Characteristics
- Supports hierarchical parent-child relationships
- Supports tagging with nested format
- Simple organizational resource for grouping contacts
- Optional description and parent fields

## Compliance Checklist

- ✅ IDPreservation test removed
- ✅ Tag tests use standard helpers (RunTagLifecycleTest, RunTagOrderTest)
- ✅ All core CRUD operations tested
- ✅ Reliability tests present (externalDeletion, removeOptionalFields)
- ✅ Validation tests present
- ✅ All tests passing
- ✅ Cleanup registered for all resources
- ⚠️ Tag format flagged for Phase 2 conversion to slug list

## Notes
- Contact Group is an organizational resource for grouping contacts
- Supports hierarchical relationships via parent field
- Tags use nested format - marked for Phase 2 conversion
- Special consistency test for literal names remains unchanged
- Full test already validates tags, added lifecycle/order tests for completeness
