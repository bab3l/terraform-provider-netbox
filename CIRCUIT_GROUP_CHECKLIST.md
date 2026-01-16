# Circuit Group Resource - Acceptance Test Checklist

## Test Results
**Status**: ✅ All tests passing
**Total**: 11/11 tests
**Duration**: ~15s

## Test Coverage

### Core CRUD Operations
- ✅ **TestAccCircuitGroupResource_basic** - Basic circuit group creation and read
- ✅ **TestAccCircuitGroupResource_full** - Full circuit group with all optional attributes
- ✅ **TestAccCircuitGroupResource_update** - Update circuit group attributes
- ✅ **TestAccCircuitGroupResource_import** - Import existing circuit group by ID

### Edge Cases & Deletion
- ✅ **TestAccCircuitGroupResource_IDPreservation** - ID preservation after updates
- ✅ **TestAccCircuitGroupResource_externalDeletion** - Handle external deletion gracefully
- ✅ **TestAccCircuitGroupResource_removeOptionalFields** - Remove optional fields (tenant)

### Validation
- ✅ **TestAccCircuitGroupResource_validationErrors** - Validation error handling
  - Missing name
  - Missing slug

### Tag Lifecycle
- ✅ **TestAccCircuitGroupResource_tagLifecycle** - Tag lifecycle (none → tags → different tags → none)
- ✅ **TestAccCircuitGroupResource_tagOrderInvariance** - Tag order doesn't trigger updates

## Additional Tests
- TestAccConsistency_CircuitGroup_LiteralNames (consistency check, not counted in standard 11)

## Test Dependencies
Tests create the following supporting resources:
- Tenant (for optional field tests)
- Tags (3 for lifecycle test, 2 for order test)

## Notes
- Circuit Group is the 8th resource with complete standardized acceptance tests (8/86 = 9.3%)
- Previous resources: IP Address, Prefix, Aggregate, ASN, ASN Range, Cable, Circuit
- Benefits from the tag lifecycle bug fix applied to Cable
