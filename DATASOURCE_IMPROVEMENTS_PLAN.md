# Datasource Improvements Plan

This document outlines the plan for improving datasource 404 handling and test coverage in the Terraform NetBox provider.

## Overview

### Current State Analysis

**Total Datasources:** 103 files in `internal/datasources/`

**Issues Identified:**

1. **15 datasources with NO "not found" handling** - These don't handle 404 responses or empty list results
2. **Many datasources with partial 404 handling** - Handle list empty results but not direct ID lookup 404s
3. **Test coverage gaps** - Many datasources have minimal test coverage (only 3 tests)

---

## Priority 1: Add "Not Found" Handling to 15 Datasources

These datasources have **no error handling** when a resource is not found:

| # | Datasource | API | Lookup Method | Complexity |
|---|-----------|-----|---------------|------------|
| 1 | cable_data_source.go | DcimAPI | ID only | Low |
| 2 | cable_termination_data_source.go | DcimAPI | ID only | Low |
| 3 | circuit_group_assignment_data_source.go | CircuitsAPI | ID only | Low |
| 4 | circuit_termination_data_source.go | CircuitsAPI | ID only | Low |
| 5 | contact_assignment_data_source.go | TenancyAPI | ID only | Low |
| 6 | event_rule_data_source.go | ExtrasAPI | ID/Name | Medium |
| 7 | fhrp_group_assignment_data_source.go | IpamAPI | ID only | Low |
| 8 | inventory_item_template_data_source.go | DcimAPI | ID/Name | Medium |
| 9 | journal_entry_data_source.go | ExtrasAPI | ID only | Low |
| 10 | l2vpn_termination_data_source.go | VpnAPI | ID only | Low |
| 11 | module_bay_template_data_source.go | DcimAPI | ID/Name | Medium |
| 12 | notification_group_data_source.go | ExtrasAPI | ID/Name | Medium |
| 13 | rack_reservation_data_source.go | DcimAPI | ID only | Low |
| 14 | virtual_device_context_data_source.go | DcimAPI | ID/Name | Medium |
| 15 | wireless_link_data_source.go | WirelessAPI | ID only | Low |

### Implementation Pattern

For each datasource, add:

1. **For ID lookup (Retrieve calls):**
```go
if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
    resp.Diagnostics.AddError(
        "Resource Not Found",
        fmt.Sprintf("No [resource] found with ID: %s", id),
    )
    return
}
```

2. **For list-based lookup (List calls with filters):**
```go
if len(results.GetResults()) == 0 {
    resp.Diagnostics.AddError(
        "Resource Not Found",
        fmt.Sprintf("No [resource] found matching the specified criteria"),
    )
    return
}
```

---

## Priority 2: Improve 404 Handling for ID Lookups

These datasources handle empty list results but **don't check 404 for direct ID retrieval**:

| # | Datasource | Has List Check | Missing ID 404 Check |
|---|-----------|----------------|---------------------|
| 1 | config_template_data_source.go | ✅ | ❌ |
| 2 | interface_template_data_source.go | ✅ | ❌ |
| 3 | inventory_item_data_source.go | ✅ | ❌ |
| 4 | inventory_item_role_data_source.go | ✅ | ❌ |
| 5 | wireless_lan_data_source.go | ✅ | ❌ |
| 6 | wireless_lan_group_data_source.go | ✅ | ❌ |

---

## Implementation Batches

### Batch 1: ID-Only Datasources (9 datasources) ✅ COMPLETE
**Estimated time:** 30 minutes | **Actual time:** 30 minutes | **Commit:** 1b9a925

Simple datasources that only support lookup by ID:

1. ✅ cable_data_source.go
2. ✅ cable_termination_data_source.go
3. ✅ circuit_group_assignment_data_source.go
4. ✅ circuit_termination_data_source.go
5. ✅ contact_assignment_data_source.go
6. ✅ fhrp_group_assignment_data_source.go
7. ✅ journal_entry_data_source.go
8. ✅ l2vpn_termination_data_source.go
9. ✅ rack_reservation_data_source.go

**Pattern:** Add 404 check after `Retrieve` API call

### Batch 2: ID/Name Datasources (6 datasources) ✅ COMPLETE
**Estimated time:** 45 minutes | **Actual time:** 20 minutes | **Commit:** 9d4bf36

Datasources that support lookup by ID or name/other filters (analysis showed ID-only):

1. ✅ event_rule_data_source.go
2. ✅ inventory_item_template_data_source.go
3. ✅ module_bay_template_data_source.go
4. ✅ notification_group_data_source.go
5. ✅ virtual_device_context_data_source.go
6. ✅ wireless_link_data_source.go

**Pattern:** Add 404 check for ID lookup AND empty results check for list lookup
**Note:** Analysis revealed these datasources only support ID lookup (no list fallback), so they follow Batch 1 pattern.

### Batch 3: Partial Handling Fixes (6 datasources)
**Status**: ✅ COMPLETE
**Commit**: a7e9042
**Actual Time**: 15 minutes

Datasources that needed ID lookup 404 handling added:

1. ✅ config_template_data_source.go
2. ✅ interface_template_data_source.go
3. ✅ inventory_item_data_source.go
4. ✅ inventory_item_role_data_source.go
5. ✅ wireless_lan_data_source.go
6. ✅ wireless_lan_group_data_source.go

**Pattern:** Add 404 status code check after ID `Retrieve` call

**Test Results**: All 20 tests PASSED
- ConfigTemplate: 3 tests, InterfaceTemplate: 3 tests
- InventoryItem: 3 tests, InventoryItemRole: 4 tests
- WirelessLAN: 3 tests, WirelessLANGroup: 4 tests

---

## Priority 3: Test Coverage Improvements

### Current Test Coverage Distribution

| Tests | Count | Percentage |
|-------|-------|------------|
| 7 tests | 22 | 21% - Excellent |
| 5-6 tests | 47 | 46% - Good |
| 3 tests | 34 | 33% - Minimal |
| 2 tests | 1 | 1% - Insufficient |

### Test Improvement Batches

#### Test Batch A: Datasources with 3 tests (34 total)

These need additional test coverage. Pattern to add:
- Separate `byID` test
- Separate `byName`/`bySlug` test (if applicable)
- Ensure all lookup paths are tested

**Split into sub-batches:**

| Sub-batch | Datasources | Status |
|-----------|-------------|--------|
| A1 | aggregate, asn, asn_range, cable, cable_termination | ✅ COMPLETE |
| A2 | circuit, circuit_termination, circuit_type, cluster, cluster_group | ✅ COMPLETE |
| A3 | cluster_type, config_context, console_port, console_port_template | Pending |
| A4 | console_server_port, console_server_port_template, contact_assignment | Pending |
| A5 | contact_group, event_rule, fhrp_group_assignment, interface | Pending |
| A6 | inventory_item_template, journal_entry, l2vpn_termination, location | Pending |
| A7 | module_bay_template, notification_group, rack_reservation, rack_role | Pending |
| A8 | virtual_device_context, virtual_machine, wireless_link | Pending |

**Test Batch A1 - ✅ COMPLETE** (Commit 9337fb4)
- **Aggregate**: Split basic → byID + byPrefix (3 tests total)
- **ASN**: Split basic → byID + byASN (3 tests total)
- **ASN Range**: Split basic → byID + byName + bySlug (4 tests total)
- **Cable**: Renamed basic → byID (2 tests total)
- **Cable Termination**: Renamed basic → byID (2 tests, both skipped)
- **Time**: 25 minutes
- **Tests**: 12 passed, 2 skipped

**Test Batch A2 - ✅ COMPLETE** (Commit fe919b6)
- **Circuit**: Split basic → byID + byCID (3 tests total)
- **Circuit Termination**: Renamed basic → byID (2 tests total)
- **Circuit Type**: Split basic → byID + byName + bySlug (4 tests total)
- **Cluster**: Split basic → byID + byName (3 tests total)
- **Cluster Group**: Renamed basic → byID (2 tests total)
- **Time**: 20 minutes
- **Tests**: 14 passed

---

## Implementation Order

### Phase 1: 404 Handling (Priority) - ✅ COMPLETE
1. ✅ Batch 1: ID-Only datasources (9) - Commit 1b9a925
2. ✅ Batch 2: ID/Name datasources (6) - Commit 9d4bf36
3. ✅ Batch 3: Partial handling fixes (6) - Commit a7e9042

**Total: 21 datasources - All implemented and tested**
**Phase Duration**: ~55 minutes total

### Phase 2: Test Coverage (Secondary)
1. ✅ Test Batch A1: aggregate, asn, asn_range, cable, cable_termination - Commit 9337fb4
2. ✅ Test Batch A2: circuit, circuit_termination, circuit_type, cluster, cluster_group - Commit fe919b6
3. ✅ Test Batch A3: cluster_type, config_context, console_port, console_port_template - Commit 1f646fb
4. ✅ Test Batch A4: console_server_port, console_server_port_template, contact_assignment - Commit e6ee1de
5. ✅ Test Batch A5: contact_group, event_rule, fhrp_group_assignment, interface - Commit 84ea32a
6. ✅ Test Batch A6: inventory_item_template, journal_entry, l2vpn_termination, location - Commit 5c3379f
7. Test Batch A7-A8: Add missing tests for remaining datasources

**Progress**: 25 datasources improved, 66 tests total (64 passed, 2 skipped)

---

## Verification

After each batch:
1. Run `go vet ./internal/datasources/...` to verify compilation
2. Run `go build .` to ensure no build errors
3. Run sample acceptance tests to verify functionality

---

## Success Criteria

- [ ] All 103 datasources have proper "not found" error handling
- [ ] All datasources check for 404 on ID lookup
- [ ] All datasources check for empty results on list-based lookup
- [ ] Test coverage improved for datasources with only 3 tests

---

## Notes

### Datasource vs Resource Pattern Differences

Datasources differ from resources in 404 handling:
- **Resources**: 404 in Read → remove from state (allow recreation)
- **Datasources**: 404 in Read → return error (fail the plan)

This is because datasources represent external data that MUST exist for the configuration to be valid.

### Import Requirements

Ensure `net/http` is imported for `http.StatusNotFound` constant usage.
