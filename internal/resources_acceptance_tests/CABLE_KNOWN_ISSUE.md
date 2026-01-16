# Cable Resource - Acceptance Test Standardization (Partial)

**Resource:** `netbox_cable`
**Date:** 2026-01-16
**Status:** ⚠️ **PARTIAL** - Known provider bug prevents full tag lifecycle testing
**Total Tests:** 10/11 passing (1 failing due to provider bug)

## Test Coverage Status

### ✅ Core Tests (9 existing)
1. ✅ **TestAccCableResource_basic** - Basic resource creation
2. ✅ **TestAccCableResource_full** - Full configuration with all attributes
3. ✅ **TestAccCableResource_update** - Resource update operations
4. ✅ **TestAccCableResource_import** - Resource import functionality
5. ✅ **TestAccCableResource_IDPreservation** - ID preservation across updates
6. ✅ **TestAccCableResource_externalDeletion** - External deletion detection
7. ✅ **TestAccCableResource_removeOptionalFields** - Optional field removal
8. ✅ **TestAccCableResource_removeDescriptionCommentsLabel** - Field removal
9. ✅ **TestAccCableResource_validationErrors** - Multi-validation error test

### ⚠️ Tag Tests (2 added, 1 failing)
10. ❌ **TestAccCableResource_tagLifecycle** - FAILS at step 4/5 due to provider bug
11. ✅ **TestAccCableResource_tagOrderInvariance** - Tag order invariance (passes)

## Known Issue: Provider Bug with Tag Lifecycle

### Problem
The Cable resource has a bug in [cable_resource.go](c:\GitRoot\terraform-provider-netbox\internal\resources\cable_resource.go) where transitioning from null tags to tags with values fails:

```
Error: Provider produced inconsistent result after apply
.tags: was null, but now cty.SetVal([...])
```

### What Works
- ✅ TagOrderInvariance test passes (changing tag order with tags always present)
- ✅ Tags can be added when creating a new resource
- ✅ Tags can be modified when already present

### What Fails
- ❌ Step 4 of tagLifecycle: Changing from null tags to different tags
- The test sequence: none → tag1,tag2 → tag2,tag3 → none → verify
- Fails at step 4 (transitioning from none back to tags)

### Root Cause
The Read function in [cable_resource.go](c:\GitRoot\terraform-provider-netbox\internal\resources\cable_resource.go):559-566 sets tags to `types.SetNull()` when no tags exist, but the provider doesn't properly handle transitions from null to a set value during Update operations.

### Workaround Attempts
- ❌ Using `tags = []` instead of omitting tags → causes "was empty set, but now null" error
- ❌ Omitting tags field → causes "was null, but now has tags" error
- ✅ Only solution: Keep tags present throughout lifecycle (tagOrderInvariance pattern works)

## Dependencies
- `netbox_site` - Required parent resource
- `netbox_manufacturer` - For device type
- `netbox_device_role` - For device
- `netbox_device_type` - For device
- `netbox_device` - For interfaces
- `netbox_interface` (2x) - For cable terminations (A and B sides)
- `netbox_tag` (3x) - For tag tests

## Implementation Details

### TestAccCableResource_tagOrderInvariance (PASSES ✅)
- Uses `RunTagOrderTest` helper
- Config function: `testAccCableResourceConfig_tagOrder`
- Tests: tag1,tag2,tag3 → tag3,tag2,tag1 (same tags, different order)
- Tag format: objects with `name` and `slug` attributes
- Duration: ~6-8s

### TestAccCableResource_tagLifecycle (FAILS ❌)
- Uses `RunTagLifecycleTest` helper
- Config function: `testAccCableResourceConfig_tagLifecycle`
- Tests: none → tag1,tag2 → tag2,tag3 → none → verify
- **Fails at step 4/5**: Transition from none to tag2,tag3
- Tag format: objects with `name` and `slug` attributes

## Recommendation

**DO NOT MARK CABLE AS COMPLETE** until the provider bug is fixed. The tag lifecycle test implementation is correct, but the underlying provider resource has a bug that needs to be addressed separately.

### Next Steps
1. File issue for Cable resource tag lifecycle bug
2. Fix bug in [cable_resource.go](c:\GitRoot\terraform-provider-netbox\internal\resources\cable_resource.go) Update/Read functions
3. Re-run tagLifecycle test after fix
4. Mark Cable as complete once all 11 tests pass

### Notes
- CheckCableDestroy helper created and working
- Both tag test implementations are correct
- Issue is in provider code, not test code
- This is the first resource encountered with this specific tag handling bug
