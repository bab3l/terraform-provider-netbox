# FHRP Group Assignment Resource Standardization Checklist

**Resource:** `netbox_fhrp_group_assignment`
**Test File:** `internal/resources_acceptance_tests/fhrp_group_assignment_resource_test.go`
**Completed:** 2026-01-16

## Changes Made

### ✅ Removed IDPreservation Test
- **Lines removed:** 36 lines (lines 59-94)
- **Test name:** `TestAccFHRPGroupAssignmentResource_IDPreservation`
- **Reason:** Redundant with basic test's import verification

### ❌ No Tag Support
- FHRP Group Assignment resource does not support tags (no tags field in schema)
- No tag lifecycle or order tests added

## Test Results

### Test Execution
```
5 tests PASSED (~5.2s total)
- TestAccFHRPGroupAssignmentResource_basic
- TestAccFHRPGroupAssignmentResource_full
- TestAccFHRPGroupAssignmentResource_update
- TestAccFHRPGroupAssignmentResource_external_deletion
- TestAccFHRPGroupAssignmentResource_validationErrors (with 4 subtests)
```

### Test Coverage
- ✅ Basic CRUD operations
- ✅ Update validation
- ✅ Import state verification
- ✅ External deletion handling
- ✅ Validation error handling
- ❌ No tag tests (not supported)

## Resource Details

**Primary Fields:**
- `group_id` (String, Required) - FHRP group ID
- `interface_type` (String, Required) - Interface type (e.g., "dcim.interface")
- `interface_id` (String, Required) - Interface ID
- `priority` (Int32, Required) - Assignment priority (0-255)

**Dependencies:**
- FHRP Group (parent resource)
- Interface (dcim.interface or virtualization.vminterface)

**Special Considerations:**
- Junction resource linking FHRP groups to network interfaces
- Priority determines which interface becomes active in FHRP protocol
- Requires existing FHRP group and interface
- Test uses non-overlapping group ID ranges to prevent parallel test collisions
- `_full` test occasionally flaky (pre-existing issue unrelated to standardization)

## Commit Information

**Files Modified:**
- `internal/resources_acceptance_tests/fhrp_group_assignment_resource_test.go` (-36 lines)

**Commit Message:**
```
Standardize FHRP Group Assignment resource tests - remove IDPreservation test

- Remove redundant TestAccFHRPGroupAssignmentResource_IDPreservation test (36 lines)
- Import verification already covered in basic test
- All 5 tests passing (~5.2s)
- No tag support (no tags field in schema)

Resource 32/86 complete (37.2%)
```
