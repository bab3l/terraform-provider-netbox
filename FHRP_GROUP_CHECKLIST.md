# FHRP Group Resource Standardization Checklist

**Resource:** `netbox_fhrp_group`
**Test File:** `internal/resources_acceptance_tests/fhrp_group_resource_test.go`
**Completed:** 2026-01-16

## Changes Made

### ✅ Removed IDPreservation Test
- **Lines removed:** 31 lines (lines 269-299, includes comment)
- **Test name:** `TestAccFHRPGroupResource_IDPreservation`
- **Reason:** Redundant with basic test's implicit ID preservation

### ✅ Added Tag Tests
- **Test 1:** `TestAccFHRPGroupResource_tagLifecycle` - Tests full tag lifecycle (add/change/remove)
- **Test 2:** `TestAccFHRPGroupResource_tagOrderInvariance` - Tests tag order independence
- **Helper functions:** `testAccFHRPGroupResourceConfig_tags` and `testAccFHRPGroupResourceConfig_tagsOrder`
- **Tag format:** ⚠️ **Nested format** `tags = [{ name = ..., slug = ... }]` (Phase 2 conversion needed)

## Test Results

### Test Execution
```
9 tests PASSED (~5.3s total)
- TestAccFHRPGroupResource_basic
- TestAccFHRPGroupResource_full
- TestAccFHRPGroupResource_update
- TestAccFHRPGroupResource_external_deletion
- TestAccFHRPGroupResource_import
- TestAccFHRPGroupResource_tagLifecycle
- TestAccFHRPGroupResource_tagOrderInvariance
- TestAccFHRPGroupResource_removeOptionalFields
- TestAccFHRPGroupResource_validationErrors (with 2 subtests)
```

### Test Coverage
- ✅ Basic CRUD operations
- ✅ Update validation
- ✅ Import state verification
- ✅ External deletion handling
- ✅ Tag lifecycle (add/change/remove)
- ✅ Tag order invariance
- ✅ Optional field removal
- ✅ Validation error handling

## Resource Details

**Primary Fields:**
- `protocol` (String, Required) - FHRP protocol (vrrp2, vrrp3, hsrp, etc.)
- `group_id` (Int32, Required) - FHRP group identifier (1-255 typical range)
- `name` (String, Optional) - Descriptive name for the group
- `auth_type` (String, Optional) - Authentication type (plaintext, md5)
- `auth_key` (String, Optional) - Authentication key
- `description` (String, Optional) - Description
- `comments` (String, Optional) - Comments
- `tags` (Set[Nested], Optional) - Tags in nested format
- `custom_fields` (Set, Optional) - Custom field values

**Dependencies:** None (referenced by FHRP Group Assignments)

**Special Considerations:**
- First Hop Redundancy Protocol configuration (VRRP/HSRP)
- Tests use non-overlapping group ID ranges to prevent parallel collisions
- Authentication supports plaintext and MD5 types
- ⚠️ **Phase 2 Item**: Convert tags from nested to slug list format

## Tag Format

**Current (Nested):**
```hcl
tags = [
  { name = netbox_tag.tag1.name, slug = netbox_tag.tag1.slug }
]
```

**Target (Slug List - Phase 2):**
```hcl
tags = [netbox_tag.tag1.slug]
```

## Commit Information

**Files Modified:**
- `internal/resources_acceptance_tests/fhrp_group_resource_test.go` (-31 lines, +65 lines tag tests)
- `COVERAGE_ANALYSIS.md` (updated progress to 33/86, added FHRP Group to Phase 2 list)

**Commit Message:**
```
Standardize FHRP Group resource tests - remove IDPreservation, add tag tests

- Remove redundant TestAccFHRPGroupResource_IDPreservation test (31 lines)
- Add tag lifecycle and order invariance tests
- Uses nested tag format (Phase 2 conversion needed)
- All 9 tests passing (~5.3s)

Resource 33/86 complete (38.4%)
```
