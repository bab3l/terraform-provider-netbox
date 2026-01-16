# Contact Role Resource Test Standardization Checklist

## Resource Information
- **Resource Name**: Contact Role
- **Test File**: `internal/resources_acceptance_tests/contact_role_resource_test.go`
- **Completion Date**: 2025-01-16
- **Status**: ✅ COMPLETE

## Test Structure Analysis

### Original Test Count
- **Before Standardization**: 8 tests (including IDPreservation and withTags)
- **After Standardization**: 8 tests

### Test Categories

#### Core CRUD Tests (3)
- ✅ `TestAccContactRoleResource_basic` - Basic resource creation with name and slug
- ✅ `TestAccContactRoleResource_full` - Complete resource with description and tags
- ✅ `TestAccContactRoleResource_update` - Update contact role name

#### Reliability Tests (2)
- ✅ `TestAccContactRoleResource_externalDeletion` - Handles external resource deletion
- ✅ `TestAccContactRoleResource_removeOptionalFields` - Handles removal of optional fields (description, tags)

#### Tag Tests (2) - **⚠️ USES NESTED FORMAT - FLAGGED FOR PHASE 2**
- ✅ `TestAccContactRoleResource_tagLifecycle` - Complete tag lifecycle (add, update, remove) using helper
- ✅ `TestAccContactRoleResource_tagOrderInvariance` - Tag order doesn't affect state using helper

#### Validation Tests (1)
- ✅ `TestAccContactRoleResource_validationErrors` - API validation errors (2 subtests: missing name, missing slug)

#### Special Tests (1)
- ✅ `TestAccConsistency_ContactRole_LiteralNames` - Plan consistency with literal names

## Changes Made

### 1. Removed IDPreservation Test ✅
- **Removed**: `TestAccContactRoleResource_IDPreservation`
- **Reason**: Duplicate of basic test functionality
- **Lines Removed**: ~25 lines

### 2. Replaced Custom Tag Test with Standard Helpers ✅
- **Removed**: `TestAccContactRoleResource_withTags` (non-standard tag test)
- **Added**: `TestAccContactRoleResource_tagLifecycle` using `testutil.RunTagLifecycleTest`
- **Added**: `TestAccContactRoleResource_tagOrderInvariance` using `testutil.RunTagOrderTest`
- **Added**: `testAccContactRoleResourceConfig_tagLifecycle` (config helper for lifecycle test)
- **Added**: `testAccContactRoleResourceConfig_tagOrder` (config helper for order test)
- **Kept**: `testAccContactRoleResourceConfig_withTags` (still used by removeOptionalFields test)
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
TestAccContactRoleResource_basic                       PASS (2.62s)
TestAccContactRoleResource_full                        PASS (1.76s)
TestAccContactRoleResource_tagLifecycle                PASS (4.11s)
TestAccContactRoleResource_tagOrderInvariance          PASS (2.68s)
TestAccContactRoleResource_update                      PASS (2.12s)
TestAccContactRoleResource_externalDeletion            PASS (1.71s)
TestAccContactRoleResource_removeOptionalFields        PASS (4.55s)
TestAccContactRoleResource_validationErrors            PASS (0.58s)

Total: 8 tests, ~5.3 seconds
```

## Technical Details

### Dependencies
- Tag resources (for tag tests)

### Key Attributes
- `name`: Contact role name (required)
- `slug`: URL-safe identifier (required)
- `description`: Description text (optional)
- `tags`: Nested format tag list (⚠️ Phase 2 conversion needed)

### Resource Characteristics
- Defines roles for contacts (e.g., "Network Engineer", "Manager")
- Supports tagging with nested format
- Simple organizational resource
- Optional description field

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
- Contact Role is an organizational resource for categorizing contacts by their role
- Tags use nested format - marked for Phase 2 conversion
- removeOptionalFields test validates removal of both description and tags
- Special consistency test for literal names remains unchanged
- Old withTags config helper retained for use in removeOptionalFields test
