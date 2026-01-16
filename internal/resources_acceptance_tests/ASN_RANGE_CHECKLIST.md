# ASN Range Resource - Acceptance Test Standardization Checklist

**Resource:** `netbox_asn_range`
**Completion Date:** 2025-01-16
**Total Tests:** 11/11 ✅
**Test Duration:** ~18s

## Test Coverage Status

### ✅ Core Tests (9 existing)
1. ✅ **TestAccASNRangeResource_basic** - Basic resource creation
2. ✅ **TestAccASNRangeResource_full** - Full configuration with all attributes
3. ✅ **TestAccASNRangeResource_update** - Resource update operations
4. ✅ **TestAccASNRangeResource_removeOptionalFields** - Optional field removal
5. ✅ **TestAccASNRangeResource_external_deletion** - External deletion detection
6. ✅ **TestAccASNRangeResource_IDPreservation** - ID preservation across updates
7. ✅ **TestAccASNRangeResource_import** - Resource import functionality
8. ✅ **TestAccASNRangeResource_removeDescription** - Description field removal
9. ✅ **TestAccASNRangeResource_validationErrors** - Multi-validation error test

### ✅ Tag Tests (2 added)
10. ✅ **TestAccASNRangeResource_tagLifecycle** - Full tag lifecycle (helper-only)
11. ✅ **TestAccASNRangeResource_tagOrderInvariance** - Tag order invariance (helper-only)

## Tag Test Implementation Details

### TestAccASNRangeResource_tagLifecycle
- Uses `RunTagLifecycleTest` helper
- Config function: `testAccASNRangeResourceConfig_tagLifecycle`
- ASN range: 65200-65300
- Tests: no tags → 2 tags → 2 different tags → no tags
- Tag format: objects with `name` and `slug` attributes
- Cleanup: ASN Range, RIR, 3 tags

### TestAccASNRangeResource_tagOrderInvariance
- Uses `RunTagOrderTest` helper
- Config function: `testAccASNRangeResourceConfig_tagOrder`
- ASN range: 65301-65401
- Tests: tag1,tag2,tag3 → tag3,tag2,tag1 (same tags, different order)
- Tag format: objects with `name` and `slug` attributes
- Cleanup: ASN Range, RIR, 3 tags

## Dependencies
- `netbox_rir` (Regional Internet Registry) - required parent resource
- `netbox_tag` (Tags) - for tag lifecycle tests
- `netbox_tenant` (Tenant) - for full configuration test

## Notes
- ASN Range resource requires tags to be specified as objects with `name` and `slug` attributes
- Empty tags must be specified as `tags = []`, not omitted
- ASN ranges used: 64512-64612 (basic), 65000-65100 (full), 65200-65300 (tag lifecycle), 65301-65401 (tag order)
- All tests use non-overlapping ASN ranges to allow parallel execution
- CheckASNRangeDestroy helper verified resource cleanup
