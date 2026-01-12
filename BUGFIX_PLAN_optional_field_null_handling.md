# Bugfix Plan: Optional Field Null Handling

## Issue Description

**Problem**: When optional fields in Terraform resources are removed from configuration (reverting to null), the provider doesn't explicitly clear them in the API request. This causes NetBox to retain the old values, leading to "Provider produced inconsistent result after apply" errors.

**Root Cause**: The `buildXRequest()` functions only set field values when they are NOT null, but don't explicitly clear them when they ARE null. This is a **systemic issue affecting most resources**.

### Example from ASN Resource (Fixed)
```go
// ❌ BEFORE - Only sets when not null
if !data.Description.IsNull() && !data.Description.IsUnknown() {
    asnRequest.SetDescription(data.Description.ValueString())
}

// ✅ AFTER - Explicitly clears when null
if !data.Description.IsNull() && !data.Description.IsUnknown() {
    asnRequest.SetDescription(data.Description.ValueString())
} else if data.Description.IsNull() {
    asnRequest.SetDescription("")  // or SetDescriptionNil() for reference types
}
```

## Scope Analysis

### Field Types Affected
1. **String fields**: `description`, `comments`, `label`, `name`, etc.
2. **Numeric fields**: Optional `int32`, `int64`, `float64`
3. **Boolean fields**: Optional `bool` (though less common)
4. **Reference fields**: Foreign keys (tenant, site, etc.) - many already handle this with `SetXNil()`

### Resources Affected
**~95+ resources** have optional string fields that potentially need fixing:
- All resources with `description` field (88 resources)
- All resources with `comments` field (43 resources)
- All resources with `label` field (30+ resources)
- Plus other optional string/numeric fields

## Testing Strategy

### 1. Generic Test Pattern
Create a reusable test helper that can be applied to any resource:

```go
// testutil/test_helpers.go
func TestRemoveOptionalStringFields(t *testing.T, config TestConfig) {
    // Step 1: Create resource with all optional fields populated
    // Step 2: Remove optional fields from config
    // Step 3: Verify fields are cleared in state (TestCheckNoResourceAttr)
    // Step 4: Re-add fields to verify they can be set again
}
```

### 2. Test Coverage Matrix
For each resource, identify:
- Which optional string/numeric fields exist
- Which fields need null handling tests
- Whether existing tests cover null scenarios

### 3. Automated Test Generation
Consider creating a test generator that:
1. Parses resource schema
2. Identifies optional fields
3. Generates test cases automatically

## Implementation Plan

### Phase 1: Analysis & Infrastructure (Week 1) ✅ COMPLETE
**Batch 1A - Analysis Script** ✅
- [x] Create comprehensive analysis script to identify ALL optional fields across ALL resources
- [x] Generate detailed report: resource → fields → current null handling status
- [x] Prioritize resources by usage/criticality

**Batch 1B - Test Infrastructure** ✅
- [x] Create generic test helper in `internal/testutil/optional_field_tests.go`
- [x] Create test configuration structs for reusable test patterns
- [x] Document test pattern in CONTRIBUTING.md and docs/TESTING_OPTIONAL_FIELDS.md

**Phase 1 Results:**
- 66 resources identified with 467 field issues
- Analysis exported to `BUGFIX_ANALYSIS_null_handling.csv`
- Test helpers: `TestRemoveOptionalFields()` and `MultiFieldOptionalTestConfig`
- Comprehensive documentation in `docs/TESTING_OPTIONAL_FIELDS.md`

### Phase 2: High-Priority Fixes (Week 2-3) 🚧 IN PROGRESS
**Batch 2A - Core IPAM Resources (8 resources)** ✅ COMPLETE
- [x] `ip_address_resource.go` - ✅ **FIXED & TESTED** (description, comments)
- [x] `prefix_resource.go` - ✅ **FIXED & TESTED** (description, comments)
- [x] `ip_range_resource.go` - ✅ **FIXED & TESTED** (description, comments)
- [x] `vlan_resource.go` - ✅ **N/A** (Status field has default value, already correct)
- [x] `vrf_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils.ApplyDescription/ApplyComments)
- [x] `aggregate_resource.go` - ✅ **FIXED & TESTED** (description, comments, date_added, tenant)
- [x] `asn_resource.go` - ✅ **FIXED & TESTED** (Phase 1 reference implementation - description, comments)
- [x] `asn_range_resource.go` - ✅ **FIXED & TESTED** (description via utils.ApplyDescription - no comments field in schema)

**Batch 2A Summary:**
- **Direct fixes**: 5 resources (ip_address, prefix, ip_range, aggregate, asn)
- **Utility function fixes**: 14+ resources via ApplyDescription/ApplyComments improvements
  - config_template, console_port, device_bay_template, interface_template
  - inventory_item, ipsec_profile, ipsec_proposal, module_bay
  - notification_group, rear_port, wireless_lan, wireless_lan_group
  - vrf, asn_range
- **Tests added**: 7 comprehensive tests (aggregate, ip_address, prefix, ip_range, vrf, asn from Phase 1, asn_range)
- **Commits**: 7 commits with detailed documentation
- **Key learning**: Date fields require SetFieldNil() not empty string

**Batch 2B - Core DCIM Resources (10 resources)** ✅ COMPLETE
- [x] `device_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils.ApplyCommonFields)
- [x] `device_type_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils.ApplyCommonFields)
- [x] `device_role_resource.go` - ✅ **FIXED & TESTED** (description via utils.ApplyDescription)
- [x] `interface_resource.go` - ✅ **FIXED & TESTED** (description, label)
- [x] `rack_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils.ApplyCommonFieldsWithMerge - already correct)
- [x] `site_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils.ApplyCommonFields)
- [x] `location_resource.go` - ✅ **FIXED & TESTED** (description via utils.ApplyDescription - already correct)
- [x] `cable_resource.go` - ✅ **FIXED & TESTED** (description, comments, label - already correct)
- [x] `power_feed_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils)
- [x] `module_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils.ApplyCommonFields)

**Batch 2B Summary:**
- **All resources verified**: 10 resources fully tested and passing
- **Provider code status**: All already correct - no bugs found in provider code
- **Root cause of test failures**: Test configuration inconsistencies (missing status/color fields)
- **Tests added**: 11 comprehensive tests covering optional field removal
- **Test fixes applied**:
  - location_resource_test.go: Added `status = "active"` to basic config
  - rack_resource_test.go: Added `status = "active"` + fixed hardcoded manufacturer to use random values
  - device_role_resource_test.go: Added `color = "aa1409"` for consistency
  - module_resource_test.go: Removed incorrect RequiredFields checking reference IDs
- **Key learning**: Test configurations must be consistent across steps, especially for computed/default fields
- **Test collision prevention**: All tests verified to use random values (testutil.RandomName/RandomSlug)
- **Debugging approach**: Created diagnostic tests that revealed configuration issues vs provider bugs
- **All tests passing**: 11/11 tests pass consistently with proper cleanup

**Batch 2C - Virtualization Resources (5 resources)** ✅ COMPLETE
- [x] `virtual_machine_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils.ApplyCommonFields - already correct)
- [x] `vm_interface_resource.go` - ✅ **FIXED & TESTED** (description field fixed - updated to usage utils.ApplyDescription)
- [x] `cluster_resource.go` - ✅ **FIXED & TESTED** (description, comments via utils.ApplyCommonFields - already correct)
- [x] `cluster_type_resource.go` - ✅ **FIXED & TESTED** (description via utils.ApplyDescription - already correct)
- [x] `cluster_group_resource.go` - ✅ **FIXED & TESTED** (description via utils.ApplyDescription - already correct)

**Batch 2C Summary:**
- **Provider code fixes**: identified and fixed bug in `vm_interface_resource.go` where description was not being explicitly cleared when null
- **Test coverage**: Added 5 new acceptance tests covering removal of optional description/comments fields
- **Test robustness**: All tests use `testutil.RandomName` to prevent collisions
- **Verification**: All 5 tests passing consistently

### Phase 3: Medium-Priority Fixes (Week 4-5) 🚧 NEXT
**Batch 3A - Circuits Resources (6 resources)** ✅ COMPLETE
- [x] `circuit_resource.go` - description, comments (Fixed & Tested - updated to use utils.ApplyDescription/ApplyComments)
- [x] `circuit_type_resource.go` - description (Verified & Tested - already correct, fixed test expectation)
- [x] `circuit_termination_resource.go` - description (Verified & Tested - already correct, fixed test expectation)
- [x] `provider_resource.go` - description, comments (Verified - uses utils.ApplyCommonFields)
- [x] `provider_account_resource.go` - description, comments (Verified - uses utils.ApplyCommonFieldsWithMerge)
- [x] `provider_network_resource.go` - description, comments (Verified - uses utils.ApplyCommonFieldsWithMerge)

**Batch 3A Summary:**
- **Code Changes**:
  - `circuit_resource.go`: Fixed `Update` method to explicitly clear description/comments when null.
  - Verified other resources already leveraged `utils` helpers correctly.
  - `circuit_group_resource.go` also verified to handle nulls correctly (explicitly sets empty string).
- **Test Fixes**:
  - Validated 3 resources (`circuit`, `circuit_type`, `circuit_termination`) with focused acceptance tests.
  - Updated test expectations to assert field absence (`TestCheckNoResourceAttr`) rather than empty strings.
  - Verified tests pass for `circuit_group` deletion of optional fields.

**Batch 3B - Tenancy & Organization (7 resources)** ✅ COMPLETE
- [x] `tenant_resource.go` - description, comments (Verified - uses utils.ApplyDescriptiveFields, added explicit group null handling check)
- [x] `tenant_group_resource.go` - description (Fixed - added parent null handling in Update)
- [x] `contact_resource.go` - description, comments (Fixed - added group null handling in Update)
- [x] `contact_role_resource.go` - description (Verified - uses utils.ApplyDescription)
- [x] `contact_group_resource.go` - description (Fixed - added parent null handling in Update)
- [x] `region_resource.go` - description (Fixed - added parent null handling in Update)
- [x] `site_group_resource.go` - description (Fixed - added parent null handling in Update)

**Batch 3C - VPN & Wireless (8 resources)**
- [ ] `tunnel_resource.go` - description, comments
- [ ] `tunnel_group_resource.go` - description
- [ ] `l2vpn_resource.go` - description, comments
- [ ] `ike_policy_resource.go` - description, comments
- [ ] `ike_proposal_resource.go` - description, comments
- [ ] `ipsec_policy_resource.go` - description, comments
- [ ] `ipsec_profile_resource.go` - description, comments
- [ ] `ipsec_proposal_resource.go` - description, comments

### Phase 4: Remaining Resources (Week 6-7)
**Batch 4A - Port & Interface Templates (16 resources)**
- All port templates: console, power, interface, front/rear ports
- [ ] `console_port_resource.go` + `console_port_template_resource.go`
- [ ] `power_port_resource.go` + `power_port_template_resource.go`
- [ ] `power_outlet_resource.go` + `power_outlet_template_resource.go`
- [ ] `interface_template_resource.go`
- [ ] `front_port_resource.go` + `front_port_template_resource.go`
- [ ] `rear_port_resource.go` + `rear_port_template_resource.go`
- [ ] `console_server_port_resource.go` + `console_server_port_template_resource.go`

**Batch 4B - Component Resources (12 resources)**
- [ ] `device_bay_resource.go` + `device_bay_template_resource.go`
- [ ] `module_bay_resource.go` + `module_bay_template_resource.go`
- [ ] `inventory_item_resource.go` + `inventory_item_template_resource.go`
- [ ] `inventory_item_role_resource.go`
- [ ] `module_type_resource.go`
- [ ] `rack_role_resource.go`
- [ ] `rack_type_resource.go`
- [ ] `virtual_chassis_resource.go`
- [ ] `virtual_device_context_resource.go`
- [ ] `virtual_disk_resource.go`

**Batch 4C - Miscellaneous (15 resources)**
- [ ] `service_resource.go` + `service_template_resource.go`
- [ ] `custom_field_resource.go`
- [ ] `custom_field_choice_set_resource.go`
- [ ] `tag_resource.go`
- [ ] `webhook_resource.go`
- [ ] `event_rule_resource.go`
- [ ] `notification_group_resource.go`
- [ ] `config_context_resource.go`
- [ ] `config_template_resource.go`
- [ ] `export_template_resource.go`
- [ ] `journal_entry_resource.go`
- [ ] `manufacturer_resource.go`
- [ ] `platform_resource.go`
- [ ] `rir_resource.go`
- [ ] `role_resource.go`

### Phase 5: Validation & Documentation (Week 8)
**Batch 5A - Comprehensive Testing**
- [ ] Run full acceptance test suite
- [ ] Verify no regressions
- [ ] Performance testing for large-scale deployments

**Batch 5B - Documentation**
- [ ] Update DEVELOPMENT.md with null handling patterns
- [ ] Create PR template checklist item for null handling
- [ ] Document in architecture decision records (ADR)

## Work Estimates

| Phase | Batches | Resources | Estimated Time |
|-------|---------|-----------|----------------|
| Phase 1 | 2 | N/A | 5 days |
| Phase 2 | 3 | 23 | 10 days |
| Phase 3 | 3 | 21 | 10 days |
| Phase 4 | 3 | 43 | 15 days |
| Phase 5 | 2 | N/A | 5 days |
| **Total** | **13** | **~87** | **45 days** |

## Per-Resource Fix Template

For each resource, the fix involves:

### 1. Code Changes (5-10 min)
```go
// In buildXRequest() function, for each optional field:

// String fields
if !data.FieldName.IsNull() && !data.FieldName.IsUnknown() {
    request.SetFieldName(data.FieldName.ValueString())
} else if data.FieldName.IsNull() {
    request.SetFieldName("")  // Clear the field
}

// Numeric fields
if !data.FieldName.IsNull() && !data.FieldName.IsUnknown() {
    request.SetFieldName(data.FieldName.ValueInt64())
} else if data.FieldName.IsNull() {
    request.SetFieldNameNil()  // Use Nil setter if available
}
```

### 2. Test Changes (15-20 min)
```go
// Add test step that removes optional fields
{
    Config: testAccResourceConfig_noOptionalFields(...),
    Check: resource.ComposeTestCheckFunc(
        resource.TestCheckNoResourceAttr("netbox_resource.test", "description"),
        resource.TestCheckNoResourceAttr("netbox_resource.test", "comments"),
    ),
},
```

### 3. Validation (5 min)
- Run resource-specific acceptance tests
- Verify no regressions

**Total per resource: ~25-35 minutes**

## Automation Opportunities

1. **Analysis Script**: Automatically identify resources with issues
2. **Code Generator**: Generate null-handling boilerplate
3. **Test Generator**: Generate test cases for optional fields
4. **CI/CD Check**: Add linter rule to catch missing null handling

## Success Criteria

- [ ] All optional fields properly handle null values
- [ ] All resources have test coverage for removing optional fields
- [ ] No "inconsistent result" errors for null value scenarios
- [ ] Documentation updated with best practices
- [ ] CI/CD prevents future regressions

## Risk Mitigation

1. **Breaking Changes**: Unlikely - we're fixing bugs, not changing behavior
2. **Test Coverage**: Each fix includes test to prevent regressions
3. **Phased Rollout**: High-priority resources first, validate before continuing
4. **Rollback Plan**: Each batch is a separate commit for easy revert

## Notes

- **Already Fixed**: `asn_resource.go` (used as reference implementation)
- **Reference Fields**: Many resources already handle these correctly with `SetXNil()`
- **Unknown Values**: Using `IsUnknown()` check to avoid issues during plan phase
- **Empty String vs Nil**: String fields use `SetX("")`, reference fields use `SetXNil()`

## Next Steps

1. ✅ Create bugfix branch
2. ✅ Document the issue and plan
3. ⏳ Create analysis script (Batch 1A)
4. ⏳ Create test infrastructure (Batch 1B)
5. ⏳ Begin Phase 2 fixes
