# Acceptance Test Coverage Analysis

## Overall Progress
**Status**: 6/86 resources complete (6.9%)

## Completed Resources

### 1. IP Address (ipam_ipaddress)
- 11/11 tests passing
- Duration: ~15s
- Checklist: IPADDRESS_CHECKLIST.md

### 2. Prefix (ipam_prefix)
- 11/11 tests passing
- Duration: ~15s
- Checklist: PREFIX_CHECKLIST.md

### 3. Aggregate (ipam_aggregate)
- 11/11 tests passing
- Duration: ~15s
- Checklist: AGGREGATE_CHECKLIST.md

### 4. ASN (ipam_asn)
- 11/11 tests passing
- Duration: ~15s
- Checklist: ASN_CHECKLIST.md

### 5. ASN Range (ipam_asn_range)
- 11/11 tests passing
- Duration: ~15s
- Checklist: ASN_RANGE_CHECKLIST.md

### 6. Cable (dcim_cable)
- 11/11 tests passing
- Duration: ~20s
- Checklist: CABLE_CHECKLIST.md
- **Notable**: Fixed provider-wide tag lifecycle bug during implementation

## Standard Test Pattern

Each resource includes:
1. **Core CRUD**: basic, full, update, import
2. **Edge Cases**: ID preservation, external deletion, remove optional fields
3. **Validation**: Validation error handling
4. **Tag Lifecycle**: Tag lifecycle and order invariance (2 tests)
5. **Total**: 11 tests per resource

## Bug Fixes Applied

### Tag Lifecycle Bug (Fixed in Cable)
- **Issue**: Provider couldn't transition from tags to no tags
- **Root Cause**: `ApplyCommonFieldsWithMerge` preserved state tags when plan had null
- **Fix**: Always use plan tags, send empty array when null
- **Impact**: Affects all resources using this helper function
- **Files Modified**:
  - `internal/utils/request_helpers.go` (ApplyCommonFieldsWithMerge, ApplyTags)

## Next Resource
Continue alphabetically through remaining 80 resources.

## Estimated Completion
- At current pace: ~4-5 resources per session
- Estimated total time: ~17-20 sessions
