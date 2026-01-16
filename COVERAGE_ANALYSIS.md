# Acceptance Test Coverage Analysis

## Overall Progress
**Status**: 15/86 resources complete (17.4%)

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
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 10. Circuit Type (circuits_circuit_type)
- 9 tests passing
- Duration: ~7s
- Checklist: CIRCUIT_TYPE_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 11. Cluster (virtualization_cluster)
- 10 tests passing (plus 1 extended variant)
- Duration: ~10.9s
- Checklist: CLUSTER_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 12. Cluster Group (virtualization_cluster_group)
- 8 tests passing
- Duration: ~5.6s
- Checklist: CLUSTER_GROUP_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 13. Cluster Type (virtualization_cluster_type)
- 7 tests passing
- Duration: ~3.1s
- Checklist: CLUSTER_TYPE_CHECKLIST.md
- **Notable**: No tag support (simple resource)

### 14. Config Context (extras_config_context)
- 8 tests passing
- Duration: ~6.9s
- Checklist: CONFIG_CONTEXT_CHECKLIST.md
- **Notable**: Uses slug list tag format, complex dependencies

### 15. Config Template (extras_config_template)
- 6 tests passing
- Duration: ~2.9s
- Checklist: CONFIG_TEMPLATE_CHECKLIST.md
- **Notable**: No tag support (simple resource)

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
Continue alphabetically through remaining 71 resources.

## Estimated Completion
- At current pace: ~4-5 resources per session
- Estimated total time: ~15-18 sessions

## Post-Standardization Tasks

### Tag Format Standardization (Phase 2)
**Status**: Planned - to be executed after all test standardization is complete

**Problem**: Resources currently use two different tag formats in Terraform HCL:
1. **Nested object format**: `tags = [{ name = ..., slug = ... }]`
2. **Slug list format**: `tags = [slug1, slug2]`

**Decision**: Standardize ALL resources to use the simpler **slug list format**

**Resources requiring conversion** (nested → slug list):
- Circuit Termination (resource 9)
- Circuit Type (resource 10)
- Cluster (resource 11)
- Cluster Group (resource 12)

**Action Items** (after test standardization complete):
1. Identify all resources using nested tag format
2. Update resource schemas to accept slug lists
3. Update resource CRUD logic to work with slug lists
4. Update all test files to use slug list format
5. Update documentation and examples
6. Create migration guide for users (breaking change)
7. Update CHANGELOG with breaking change notice

**Rationale**: Simpler user experience, less confusion, more consistent with majority of resources
