# Device Resource Test Standardization Checklist

## Resource Information
- **Resource Name**: Device
- **Test File**: `internal/resources_acceptance_tests/device_resource_test.go`
- **Completion Date**: 2025-01-16
- **Status**: ✅ COMPLETE

## Test Structure Analysis

### Original Test Count
- **Before Standardization**: 9 tests (including IDPreservation)
- **After Standardization**: 10 tests (added tag tests)

### Test Categories

#### Core CRUD Tests (4)
- ✅ `TestAccDeviceResource_basic` - Basic resource creation with name, device_type, role, site, includes import test
- ✅ `TestAccDeviceResource_full` - Complete resource with all optional fields populated
- ✅ `TestAccDeviceResource_update` - Update name, device_type, role, site, status, and other fields
- ✅ `TestAccDeviceResource_StatusOptionalField` - Tests status field behavior with default value

#### Reliability Tests (3)
- ✅ `TestAccDeviceResource_externalDeletion` - Handles external resource deletion
- ✅ `TestAccDeviceResource_removeDescriptionAndComments` - Handles removal of description and comments fields
- ✅ `TestAccDeviceResource_removeOptionalFields` - Handles removal of serial, asset_tag, tenant, platform fields

#### Validation Tests (1)
- ✅ `TestAccDeviceResource_validationErrors` - API validation errors (6 subtests: missing device_type, missing role, missing site, invalid status, invalid tenant reference, invalid platform reference)

#### Tag Tests (2)
- ✅ `TestAccDeviceResource_tagLifecycle` - Complete tag lifecycle (add, change, remove)
- ✅ `TestAccDeviceResource_tagOrderInvariance` - Verifies tag order doesn't cause drift

## Changes Made

### 1. Removed IDPreservation Test ✅
- **Removed**: `TestAccDeviceResource_IDPreservation`
- **Reason**: Duplicate of basic test functionality
- **Lines Removed**: ~38 lines

### 2. Added Tag Tests ✅
- **Added**: `TestAccDeviceResource_tagLifecycle` - Tests full tag lifecycle
- **Added**: `TestAccDeviceResource_tagOrderInvariance` - Tests tag order invariance
- **Added**: Helper functions `testAccDeviceResourceConfig_tagLifecycle` and `testAccDeviceResourceConfig_tagOrder`
- **Format**: ⚠️ Uses nested `{name, slug}` format (requires Phase 2 conversion)

## Test Execution Results

### All Tests Passing ✅
```
TestAccDeviceResource_basic                              PASS (4.76s)
TestAccDeviceResource_full                               PASS (3.45s)
TestAccDeviceResource_update                             PASS (5.05s)
TestAccDeviceResource_StatusOptionalField                PASS (7.11s)
TestAccDeviceResource_externalDeletion                   PASS (4.99s)
TestAccDeviceResource_removeDescriptionAndComments       PASS (6.37s)
TestAccDeviceResource_removeOptionalFields               PASS (6.63s)
TestAccDeviceResource_validationErrors                   PASS (2.83s)
TestAccDeviceResource_tagLifecycle                       PASS (9.50s)
TestAccDeviceResource_tagOrderInvariance                 PASS (5.19s)

Total: 10 tests, ~12.5 seconds
```

## Technical Details

### Dependencies
- Manufacturer resource (required, for device type)
- Device Type resource (required)
- Device Role resource (required)
- Site resource (required)
- Tenant resource (optional)
- Platform resource (optional)

### Key Attributes
- `name`: Device name (required)
- `device_type`: Device type ID (required)
- `role`: Device role ID (required)
- `site`: Site ID (required)
- `status`: Device status (optional, default: "active")
- `tenant`: Tenant ID (optional)
- `platform`: Platform ID (optional)
- `serial`: Serial number (optional)
- `asset_tag`: Asset tag (optional)
- `description`: Description text (optional)
- `comments`: Additional comments (optional)
- `airflow`: Airflow direction (optional)
- `tags`: Tag list using nested format (optional)
- `custom_fields`: Custom field values (optional)

### Resource Characteristics
- Core physical device resource in NetBox
- ⚠️ **Uses nested tag format `{name, slug}`** - flagged for Phase 2 conversion
- Complex dependency chain requiring multiple supporting resources
- Status field has default value ("active")
- Import test included in basic test
- Comprehensive validation testing with 6 subtests
- Multiple optional field removal tests
- Complex CRUD operations due to many dependencies

## Compliance Checklist

- ✅ IDPreservation test removed
- ✅ Tag tests added (tagLifecycle, tagOrderInvariance)
- ⚠️ **Nested tag format** - needs Phase 2 conversion to slug list
- ✅ All core CRUD operations tested
- ✅ Reliability tests present (3 different removal scenarios)
- ✅ Validation tests present with comprehensive subtests
- ✅ All tests passing
- ✅ Cleanup registered for all resources

## Notes
- Device is a critical NetBox resource representing physical network devices
- ⚠️ **PHASE 2 REQUIRED**: Uses nested `{name, slug}` tag format - needs conversion to slug list format
- Complex resource with many optional fields and dependencies
- Tests cover multiple optional field removal scenarios
- Status field has special handling with default value
- Extensive validation testing for all reference fields
- Basic test includes both creation and import verification
- Tag lifecycle tests verify proper handling of tag additions, changes, and removals
- Tag order tests ensure reordering doesn't cause unwanted diffs
