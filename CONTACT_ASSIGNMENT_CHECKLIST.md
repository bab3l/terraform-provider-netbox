# Contact Assignment Resource Test Standardization Checklist

## Resource Information
- **Resource Name**: Contact Assignment
- **Test File**: `internal/resources_acceptance_tests/contact_assignment_resource_test.go`
- **Completion Date**: 2025-01-15
- **Status**: ✅ COMPLETE

## Test Structure Analysis

### Original Test Count
- **Before Standardization**: 9 tests
- **After Standardization**: 9 tests

### Test Categories

#### Core CRUD Tests (3)
- ✅ `TestAccContactAssignmentResource_basic` - Basic resource creation
- ✅ `TestAccContactAssignmentResource_full` - Complete resource with all attributes including tags
- ✅ `TestAccContactAssignmentResource_withRole` - Resource with role assignment
- ✅ `TestAccContactAssignmentResource_update` - Update priority from primary to secondary

#### Reliability Tests (2)
- ✅ `TestAccContactAssignmentResource_externalDeletion` - Handles external resource deletion
- ✅ `TestAccContactAssignmentResource_removeOptionalFields` - Handles removal of optional fields

#### Tag Tests (2) - **⚠️ USES NESTED FORMAT - FLAGGED FOR PHASE 2**
- ✅ `TestAccContactAssignmentResource_tagLifecycle` - Complete tag lifecycle (add, update, remove) using helper
- ✅ `TestAccContactAssignmentResource_tagOrderInvariance` - Tag order doesn't affect state using helper

#### Validation Tests (1)
- ✅ `TestAccContactAssignmentResource_validationErrors` - API validation errors (3 subtests)

#### Special Tests (1)
- ✅ `TestAccConsistency_ContactAssignment_LiteralNames` - Plan consistency with literal names

## Changes Made

### 1. Removed IDPreservation Test ✅
- **Removed**: `TestAccContactAssignmentResource_IDPreservation`
- **Reason**: Duplicate of basic test functionality
- **Lines Removed**: ~40 lines

### 2. Replaced Custom Tag Test with Standard Helpers ✅
- **Removed**: `TestAccContactAssignmentResource_withTags` (non-standard tag test)
- **Removed**: `testAccContactAssignmentResourceConfig_withTags` (unused config helper)
- **Added**: `TestAccContactAssignmentResource_tagLifecycle` using `testutil.RunTagLifecycleTest`
- **Added**: `TestAccContactAssignmentResource_tagOrderInvariance` using `testutil.RunTagOrderTest`
- **Added**: `testAccContactAssignmentResourceConfig_tagLifecycle` (config helper for lifecycle test)
- **Added**: `testAccContactAssignmentResourceConfig_tagOrder` (config helper for order test)
- **Tag Format**: Uses nested format `{name = ..., slug = ...}` - **⚠️ FLAGGED FOR PHASE 2 CONVERSION**

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
TestAccContactAssignmentResource_full                      PASS (3.00s)
TestAccContactAssignmentResource_basic                     PASS (3.01s)
TestAccContactAssignmentResource_withRole                  PASS (2.93s)
TestAccContactAssignmentResource_tagLifecycle              PASS (9.65s)
TestAccContactAssignmentResource_tagOrderInvariance        PASS (4.12s)
TestAccContactAssignmentResource_update                    PASS (3.46s)
TestAccContactAssignmentResource_externalDeletion          PASS (3.59s)
TestAccContactAssignmentResource_removeOptionalFields      PASS (3.46s)
TestAccContactAssignmentResource_validationErrors          PASS (0.86s)

Total: 9 tests, ~10.7 seconds
```

## Technical Details

### Dependencies
- Site resource (for object_id)
- Contact resource (for contact_id)
- Contact Role resource (for role_id, optional)
- Tag resources (for tag tests)

### Key Attributes
- `object_type`: ContentType for the assignment (e.g., "dcim.site")
- `object_id`: ID of the assigned object
- `contact_id`: ID of the contact (required)
- `role_id`: ID of the contact role (optional)
- `priority`: Assignment priority (primary/secondary/tertiary/inactive)
- `tags`: Nested format tag list (⚠️ Phase 2 conversion needed)

### Resource Characteristics
- Generic assignment resource linking contacts to any NetBox object
- Uses ContentType pattern for flexible object assignments
- Supports tagging with nested format
- Optional role and priority fields

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
- Contact Assignment is a generic resource that can link contacts to any NetBox object type
- Uses ContentType pattern similar to other assignment/link resources
- Tags use nested format - marked for Phase 2 conversion
- withRole test remains as it tests role_id functionality, not tags
- Special consistency test for literal names remains unchanged
