# Cable Resource - Acceptance Test Checklist

## Test Results
**Status**: ✅ All tests passing
**Total**: 11/11 tests
**Duration**: ~20s

## Test Coverage

### Core CRUD Operations
- ✅ **TestAccCableResource_basic** - Basic cable creation and read
- ✅ **TestAccCableResource_full** - Full cable with all optional attributes
- ✅ **TestAccCableResource_update** - Update cable attributes
- ✅ **TestAccCableResource_import** - Import existing cable by ID

### Edge Cases & Deletion
- ✅ **TestAccCableResource_IDPreservation** - ID preservation after updates
- ✅ **TestAccCableResource_externalDeletion** - Handle external deletion gracefully
- ✅ **TestAccCableResource_removeOptionalFields** - Remove optional fields
- ✅ **TestAccCableResource_removeDescriptionCommentsLabel** - Remove description/comments/label

### Validation
- ✅ **TestAccCableResource_validationErrors** - Validation error handling
  - Missing a_terminations
  - Missing b_terminations

### Tag Lifecycle
- ✅ **TestAccCableResource_tagLifecycle** - Tag lifecycle (none → tags → different tags → none)
- ✅ **TestAccCableResource_tagOrderInvariance** - Tag order doesn't trigger updates

## Bug Fixes Applied

### Tag Lifecycle Bug (Fixed)
**Issue**: Provider produced inconsistent result when removing tags from configuration.
- When tags were present, then removed from config, provider would show: "was null, but now cty.SetVal([tags])"
- Root cause: `ApplyCommonFieldsWithMerge` preserved state tags when plan had null tags
- Fix: Changed logic to always use plan tags, sending empty array when null

**Files Modified**:
- `internal/utils/request_helpers.go`:
  - `ApplyCommonFieldsWithMerge`: Removed state tag preservation logic
  - `ApplyTags`: Send empty array when tags are null instead of skipping

**Result**: Both tag tests now pass ✅

## Test Dependencies
Tests create the following supporting resources:
- Sites
- Manufacturers
- Device Roles
- Device Types
- Devices (2 per test)
- Interfaces (2 per test - for cable terminations)
- Tags (3 for tag tests)

## Notes
- Cable is the 6th resource with complete standardized acceptance tests (6/86 = 6.9%)
- Previous resources: IP Address, Prefix, Aggregate, ASN, ASN Range
- Fix for tag lifecycle bug affects all resources using `ApplyCommonFieldsWithMerge`
