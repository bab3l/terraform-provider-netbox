# Front Port Resource Standardization Checklist

**Resource:** `netbox_front_port`
**Test File:** `internal/resources_acceptance_tests/front_port_resource_test.go`
**Completed:** 2026-01-16

## Changes Made

### ✅ Removed IDPreservation Test
- **Lines removed:** 36 lines (lines 287-322)
- **Test name:** `TestAccFrontPortResource_IDPreservation`
- **Reason:** Redundant with basic test's implicit ID preservation

### ✅ Added Tag Tests
- **Test 1:** `TestAccFrontPortResource_tagLifecycle` - Tests full tag lifecycle (add/change/remove)
- **Test 2:** `TestAccFrontPortResource_tagOrderInvariance` - Tests tag order independence
- **Helper functions:** `testAccFrontPortResourceConfig_tags` and `testAccFrontPortResourceConfig_tagsOrder`
- **Tag format:** ⚠️ **Nested format** `tags = [{ name = ..., slug = ... }]` (Phase 2 conversion needed)

## Test Results

### Test Execution
```
8 tests PASSED (~9.7s total)
- TestAccFrontPortResource_basic
- TestAccFrontPortResource_full
- TestAccFrontPortResource_update
- TestAccFrontPortResource_externalDeletion
- TestAccFrontPortResource_tagLifecycle
- TestAccFrontPortResource_tagOrderInvariance
- TestAccFrontPortResource_removeOptionalFields
- TestAccFrontPortResource_validationErrors (with 6 subtests)
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
- `device` (String, Required) - Device ID
- `name` (String, Required) - Port name
- `type` (String, Required) - Port type (8p8c, bnc, fc, etc.)
- `rear_port` (String, Required) - Rear port ID
- `rear_port_position` (Int32, Optional, Default: 1) - Position on rear port
- `label` (String, Optional) - Physical label
- `mark_connected` (Bool, Optional, Default: false) - Mark as connected
- `description` (String, Optional) - Description
- `tags` (Set[Nested], Optional) - Tags in nested format
- `custom_fields` (Set, Optional) - Custom field values

**Dependencies:**
- Device (required)
- Rear Port (required)

**Special Considerations:**
- Physical network port on front panel of device
- Links to rear port for internal connectivity
- rear_port_position determines which position on rear port (for breakout cables)
- Complex test infrastructure (requires Site, Manufacturer, Device Type, Device Role, Device, Rear Port)
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
- `internal/resources_acceptance_tests/front_port_resource_test.go` (-36 lines, +140 lines tag tests)
- `COVERAGE_ANALYSIS.md` (updated progress to 34/86, added Front Port to Phase 2 list)

**Commit Message:**
```
Standardize Front Port resource tests - remove IDPreservation, add tag tests

- Remove redundant TestAccFrontPortResource_IDPreservation test (36 lines)
- Add tag lifecycle and order invariance tests
- Uses nested tag format (Phase 2 conversion needed)
- All 8 tests passing (~9.7s)

Resource 34/86 complete (39.5%)
```
