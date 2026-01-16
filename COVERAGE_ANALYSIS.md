# Acceptance Test Coverage Analysis

## Overall Progress
**Status**: 11/86 resources complete (12.8%)

## Completed Resources

### 1. IP Address (ipam_ipaddress)
- 11 tests passing
- Duration: ~15s
- Checklist: IPADDRESS_CHECKLIST.md

### 2. Prefix (ipam_prefix)
- 13 tests passing
- Duration: ~15s
- Checklist: PREFIX_CHECKLIST.md

### 3. Aggregate (ipam_aggregate)
- 9 tests passing
- Duration: ~10s
- Checklist: AGGREGATE_CHECKLIST.md

### 4. ASN (ipam_asn)
- 8 tests passing
- Duration: ~10s
- Checklist: ASN_CHECKLIST.md

### 5. ASN Range (ipam_asn_range)
- 10 tests passing
- Duration: ~12s
- Checklist: ASN_RANGE_CHECKLIST.md

### 6. Cable (dcim_cable)
- 10 tests passing
- Duration: ~20s
- Checklist: CABLE_CHECKLIST.md
- **Notable**: Fixed provider-wide tag lifecycle bug during implementation

### 7. Circuit (circuits_circuit)
- 10 tests passing
- Duration: ~15s
- Checklist: CIRCUIT_CHECKLIST.md

### 8. Circuit Group (circuits_circuit_group)
- 9 tests passing
- Duration: ~9s
- Checklist: CIRCUIT_GROUP_CHECKLIST.md

### 9. Circuit Termination (circuits_circuit_termination)
- 9 tests passing
- Duration: ~7.5s
- Checklist: CIRCUIT_TERMINATION_CHECKLIST.md
- **Notable**: Uses nested tag format `{name, slug}` instead of simple ID list

### 10. Circuit Type (circuits_circuit_type)
- 9 tests passing
- Duration: ~7s
- Checklist: CIRCUIT_TYPE_CHECKLIST.md
- **Notable**: Uses nested tag format `{name, slug}` like Circuit Termination

### 11. Cluster (virtualization_cluster)
- 10 tests passing (plus 1 extended variant)
- Duration: ~10.9s
- Checklist: CLUSTER_CHECKLIST.md
- **Notable**: Uses nested tag format `{name, slug}` like Circuit Termination and Circuit Type

## Standard Test Pattern

Each resource includes:
1. **Core CRUD**: basic, full, update, import (4 tests)
2. **Reliability**: external deletion, remove optional fields (2 tests)
3. **Validation**: Validation error handling (1 test, recommended)
4. **Tag Tests**: Tag lifecycle and order invariance (2 tests if resource supports tags)
5. **Total**: 8-10 tests per resource (varies by resource complexity)

**Note**: IDPreservation test was removed as it was a duplicate of the basic test.

## Bug Fixes Applied

### Tag Lifecycle Bug (Fixed in Cable)
- **Issue**: Provider couldn't transition from tags to no tags
- **Root Cause**: `ApplyCommonFieldsWithMerge` preserved state tags when plan had null
- **Fix**: Always use plan tags, send empty array when null
- **Impact**: Affects all resources using this helper function
- **Files Modified**:
  - `internal/utils/request_helpers.go` (ApplyCommonFieldsWithMerge, ApplyTags)

## Next Resource
Continue alphabetically through remaining 75 resources.

## Estimated Completion
- At current pace: ~4-5 resources per session
- Estimated total time: ~15-19 sessions
