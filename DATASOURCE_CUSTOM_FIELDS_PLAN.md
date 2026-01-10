# Datasource Custom Fields Implementation Plan

## Progress Summary
- **Batch 1**: ✅ COMPLETE (13/13 datasources - 100%)
- **Batch 2**: ✅ COMPLETE (12/12 datasources - 100%)
- **Batch 3**: ✅ COMPLETE (10/10 datasources - 100%)
- **Batch 4**: ✅ COMPLETE (12/12 datasources - 100%)
- **Batch 5**: ✅ COMPLETE (12/12 datasources - 100%)
- **Batch 6**: ✅ COMPLETE (10/10 datasources - 100%)
- **Batch 7**: ✅ COMPLETE (3/11 datasources - 100%, 8 not supported)
- **Batch 8**: 🔧 IN PROGRESS (~10 datasources - Templates & Miscellaneous)

**Overall Progress**: 72/80 datasources complete (90%)

## Overview
This document outlines the plan to fix datasource behavior with custom fields. Currently, datasources return NO custom fields because they pass `nil` as the second parameter to `MapToCustomFieldModels()`, which causes the function to return immediately.

**Goal**: Datasources should always return ALL custom fields present in NetBox, unlike resources which use partial management (filter-to-owned pattern).

## Current Problem

### Root Cause
All datasources call:
```go
customFields := utils.MapToCustomFieldModels(apiCustomFields, nil)
```

The `MapToCustomFieldModels()` function checks:
```go
if len(stateCustomFields) == 0 {
    return nil
}
```

Since datasources pass `nil`, this returns immediately with no custom fields.

### Impact
- Users cannot read custom field values via datasources
- Data sources are incomplete representations of NetBox objects
- Forces users to use resources for read-only operations

## Solution Design

### 1. Create New Helper Function
Add `MapAllCustomFieldsToModels()` in `internal/utils/common.go`:

```go
// MapAllCustomFieldsToModels converts ALL custom fields from NetBox API to Terraform models.
// Used by datasources to return complete custom field data.
// Unlike MapToCustomFieldModels which filters to owned fields for resources,
// this returns ALL fields present in NetBox.
func MapAllCustomFieldsToModels(customFields map[string]interface{}) []CustomFieldModel {
    if len(customFields) == 0 {
        return nil
    }

    result := make([]CustomFieldModel, 0, len(customFields))

    for name, value := range customFields {
        if value == nil {
            continue  // Skip null values
        }

        cf := CustomFieldModel{
            Name: types.StringValue(name),
        }

        // Infer type from value and format accordingly
        switch v := value.(type) {
        case map[string]interface{}:
            // JSON/object type
            cf.Type = types.StringValue("json")
            if jsonBytes, err := json.Marshal(v); err == nil {
                cf.Value = types.StringValue(string(jsonBytes))
            } else {
                cf.Value = types.StringValue("")
            }

        case []interface{}:
            // Multiselect type
            cf.Type = types.StringValue("multiselect")
            var stringValues []string
            for _, item := range v {
                if s, ok := item.(string); ok {
                    stringValues = append(stringValues, strings.TrimSpace(s))
                } else {
                    stringValues = append(stringValues, fmt.Sprintf("%v", item))
                }
            }
            cf.Value = types.StringValue(strings.Join(stringValues, ","))

        case bool:
            // Boolean type
            cf.Type = types.StringValue("boolean")
            cf.Value = types.StringValue(fmt.Sprintf("%t", v))

        case float64:
            // Number type (JSON unmarshals numbers as float64)
            cf.Type = types.StringValue("integer")
            cf.Value = types.StringValue(fmt.Sprintf("%.0f", v))

        case string:
            // Text/URL/select type (we'll default to text)
            cf.Type = types.StringValue("text")
            cf.Value = types.StringValue(strings.TrimSpace(v))

        default:
            // Fallback to text
            cf.Type = types.StringValue("text")
            cf.Value = types.StringValue(fmt.Sprintf("%v", v))
        }

        result = append(result, cf)
    }

    // Sort by name for consistent ordering
    sort.Slice(result, func(i, j int) bool {
        return result[i].Name.ValueString() < result[j].Name.ValueString()
    })

    return result
}
```

### 2. Update All Datasources
Change from:
```go
customFields := utils.MapToCustomFieldModels(apiCustomFields, nil)
```

To:
```go
customFields := utils.MapAllCustomFieldsToModels(apiCustomFields)
```

### 3. Create Test Infrastructure

#### New Test Directory Structure
```
internal/
├── datasources_acceptance_tests/         (existing - parallel tests)
│   └── *_data_source_test.go
└── datasources_acceptance_tests_customfields/  (NEW - serial tests)
    ├── test_main_test.go                 (test setup with build tag)
    └── *_custom_fields_test.go           (custom field specific tests)
```

#### Test File Template
```go
//go:build customfields

package datasources_acceptance_tests_customfields

import (
    "context"
    "fmt"
    "testing"

    "github.com/bab3l/terraform-provider-netbox/internal/testutil"
    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccXxxDataSource_customFields(t *testing.T) {
    // Test that datasource returns ALL custom fields
    // 1. Create resource with custom field via API
    // 2. Read via datasource
    // 3. Verify all custom fields present
    // 4. Add another custom field via API
    // 5. Read again via datasource
    // 6. Verify both custom fields present
}
```

## Implementation Batches

### Statistics
- **Total Datasources**: 104 datasource files
- **Datasources with Custom Fields**: ~80 (based on grep results)
- **Datasources without Custom Fields**: ~24 (templates, assignments, etc.)

### Batch Organization
Split into 8 batches of ~10-13 datasources each for manageable PRs.

### Batch 1: Core Infrastructure (13 datasources) ✅ COMPLETE
**Priority**: HIGH - Most commonly used datasources
**Status**: ✅ **COMPLETE** - All 13 datasources implemented and tested

1. ✅ `site_data_source.go` - Simple fix (already had custom fields)
2. ✅ `asn_data_source.go` - Simple fix (already had custom fields)
3. ✅ `asn_range_data_source.go` - **Full implementation** (was missing custom fields)
4. ✅ `circuit_data_source.go` - Simple fix (already had custom fields)
5. ✅ `circuit_type_data_source.go` - Simple fix (already had custom fields)
6. ✅ `cluster_data_source.go` - Simple fix (already had custom fields)
7. ✅ `cluster_type_data_source.go` - Simple fix (already had custom fields)
8. ✅ `ip_address_data_source.go` - **Full implementation** (was missing custom fields)
9. ✅ `ip_range_data_source.go` - **Full implementation** (was missing custom fields)
10. ✅ `prefix_data_source.go` - **Full implementation** (was missing custom fields)
11. ✅ `vlan_data_source.go` - **Full implementation** (was missing custom fields)
12. ✅ `vrf_data_source.go` - Simple fix (already had custom fields)
13. ✅ `route_target_data_source.go` - **Full implementation** (was missing custom fields)

**Test Priority**: HIGH - ✅ All 13 tests created and passing
**Total Test Time**: ~15 seconds for all tests
**Implementation Time**: ~4 hours

### Batch 2: Device & DCIM (12 datasources) ✅ COMPLETE
**Priority**: HIGH - Core hardware management
**Status**: ✅ **COMPLETE** - All 12 datasources implemented and tested

1. ✅ `device_data_source.go` - Simple fix (already had custom fields)
2. ✅ `device_type_data_source.go` - Simple fix (already had custom fields)
3. ✅ `device_role_data_source.go` - Simple fix (already had custom fields)
4. ✅ `rack_data_source.go` - Simple fix (already had custom fields)
5. ✅ `rack_role_data_source.go` - Simple fix (already had custom fields)
6. ✅ `location_data_source.go` - Simple fix (already had custom fields)
7. ✅ `manufacturer_data_source.go` - **Full implementation** (was missing custom fields)
8. ✅ `platform_data_source.go` - **Full implementation** (was missing custom fields)
9. ✅ `interface_data_source.go` - Simple fix (already had custom fields)
10. ✅ `cable_data_source.go` - **Full implementation** (was missing custom fields)
11. ✅ `device_bay_data_source.go` - Simple fix (already had custom fields)
12. ✅ `module_data_source.go` - **Full implementation** (was missing custom fields)

**Test Priority**: HIGH - ✅ All 12 tests created and passing
**Total Test Time**: ~25 seconds for all tests
**Implementation Time**: ~5 hours

### Batch 3: Virtualization & Tenancy (10 datasources) ✅ COMPLETE
**Priority**: MEDIUM - Virtualization and multi-tenancy
**Status**: ✅ **COMPLETE** - All 10 datasources implemented and tested

1. ✅ `virtual_machine_data_source.go` - Simple fix (already had custom fields)
2. ✅ `vm_interface_data_source.go` - Simple fix (already had custom fields)
3. ✅ `virtual_disk_data_source.go` - **Full implementation** (was missing custom fields)
4. ✅ `virtual_device_context_data_source.go` - Simple fix (already had custom fields)
5. ✅ `cluster_group_data_source.go` - **Full implementation** (was missing custom fields)
6. ✅ `tenant_data_source.go` - Simple fix (already had custom fields)
7. ✅ `tenant_group_data_source.go` - Simple fix (already had custom fields)
8. ✅ `contact_data_source.go` - **Full implementation** (resource doesn't support custom_fields)
9. ✅ `contact_role_data_source.go` - **Full implementation** (was missing custom fields)
10. ✅ `contact_group_data_source.go` - **Full implementation** (was missing custom fields)

**Test Priority**: MEDIUM - ✅ All 10 tests created and passing
**Total Test Time**: ~35 seconds for all tests
**Implementation Time**: ~6 hours

### Batch 4: Circuits & VPN (12 datasources) ✅ COMPLETE
**Priority**: MEDIUM - Network connectivity
**Status**: ✅ **COMPLETE** - All 12 datasources implemented and tested

1. ✅ `provider_data_source.go` - Simple fix (already had custom fields)
2. ✅ `provider_account_data_source.go` - **Full implementation** (was missing custom fields)
3. ✅ `provider_network_data_source.go` - Simple fix (already had custom fields)
4. ✅ `circuit_group_data_source.go` - Simple fix (already had custom fields)
5. ✅ `l2vpn_data_source.go` - Refactored (removed complex existingModels logic)
6. ✅ `tunnel_data_source.go` - Simple fix (already had custom fields)
7. ✅ `tunnel_group_data_source.go` - Simple fix (already had custom fields)
8. ✅ `ike_policy_data_source.go` - **Full implementation** (was missing custom fields)
9. ✅ `ike_proposal_data_source.go` - **Full implementation** (was missing custom fields)
10. ✅ `ipsec_policy_data_source.go` - **Full implementation** (was missing custom fields)
11. ✅ `ipsec_profile_data_source.go` - **Full implementation** (was missing custom fields)
12. ✅ `ipsec_proposal_data_source.go` - **Full implementation** (was missing custom fields)

**Test Priority**: MEDIUM - ✅ All 12 tests created and passing
**Total Test Time**: ~16 seconds for all tests
**Implementation Time**: ~7 hours

### Batch 5: Ports & Interfaces (12 datasources) ✅ COMPLETE
**Priority**: MEDIUM - Port management
**Status**: ✅ **COMPLETE** - All 12 datasources implemented and tested

1. ✅ `console_port_data_source.go` - **Full implementation** (was missing custom fields)
2. ✅ `console_server_port_data_source.go` - **Full implementation** (was missing custom fields)
3. ✅ `power_port_data_source.go` - **Full implementation** (was missing custom fields)
4. ✅ `power_outlet_data_source.go` - **Full implementation** (was missing custom fields) + **RESOURCE BUG FIX**
5. ✅ `front_port_data_source.go` - **Full implementation** (was missing custom fields)
6. ✅ `rear_port_data_source.go` - **Full implementation** (was missing custom fields)
7. ✅ `module_bay_data_source.go` - **Full implementation** (was missing custom fields)
8. ✅ `inventory_item_data_source.go` - **Full implementation** (was missing custom fields)
9. ✅ `inventory_item_role_data_source.go` - **Full implementation** (was missing custom fields)
10. ✅ `power_feed_data_source.go` - Fixed pattern (had custom fields but wrong function)
11. ✅ `power_panel_data_source.go` - Fixed pattern (had custom fields but wrong function)
12. ✅ `rack_reservation_data_source.go` - Fixed pattern (had custom fields but wrong function)

**Test Priority**: MEDIUM - ✅ All 12 tests created and passing
**Total Test Time**: ~46 seconds for all tests
**Implementation Time**: ~5 hours + bug fix

#### Critical Bug Fix: power_outlet Resource
During implementation, discovered a long-standing bug in `power_outlet_resource.go`:

**Problem**: Resource incorrectly created `BriefPowerPortRequest` with fabricated name:
```go
powerPortReq := netbox.BriefPowerPortRequest{
    Name: fmt.Sprintf("Power Port %d", powerPortID),
}
```

This caused API errors: "Related object not found using the provided attributes: {'name': 'Power Port 418'}"

**Root Cause**:
- Bug existed since resource creation
- Existing tests deliberately avoided testing `power_port` attribute
- First test to actually use `power_port` was our datasource test

**Solution**:
1. Created `LookupPowerPort()` function in `internal/netboxlookup/generic_lookup.go`
2. Implemented `PowerPortLookupConfig` using `GenericLookup` pattern
3. Fixed both Create and Update functions in `power_outlet_resource.go`
4. Handled complex type: `BriefDeviceRequest.Name` is `NullableString`, requires `.Set(&deviceName)`

**Files Modified**:
- `internal/netboxlookup/generic_lookup.go`: Added PowerPortLookupConfig and LookupPowerPort
- `internal/resources/power_outlet_resource.go`: Fixed lines ~161 (Create) and ~300 (Update)

**Impact**: Power outlet resource now functional with power_port references for the first time

### Batch 6: Wireless & Services (10 datasources) ✅ COMPLETE
**Priority**: LOW - Specialized features
**Status**: ✅ **COMPLETE** - All 10 datasources implemented and tested

1. ✅ `wireless_lan_data_source.go` - **Full implementation** (was missing custom fields)
2. ✅ `wireless_lan_group_data_source.go` - **Full implementation** (was missing custom fields)
3. ✅ `wireless_link_data_source.go` - Fixed pattern (had custom fields but wrong function)
4. ✅ `fhrp_group_data_source.go` - **Full implementation** (was missing custom fields)
5. ✅ `service_data_source.go` - Fixed pattern (had custom fields but wrong function)
6. ✅ `service_template_data_source.go` - Fixed pattern (had custom fields but wrong function)
7. ✅ `aggregate_data_source.go` - **Full implementation** (was missing custom fields)
8. ✅ `rir_data_source.go` - **Full implementation** (was missing custom fields)
9. ✅ `vlan_group_data_source.go` - Fixed pattern (had custom fields but wrong function)
10. ✅ `role_data_source.go` - Fixed pattern (had custom fields but wrong function)

**Test Priority**: LOW - ✅ All 10 tests created and passing
**Total Test Time**: ~22 seconds for all tests
**Implementation Time**: ~2.5 hours

### Batch 7: Extras & Admin (11 datasources) ✅ COMPLETE
**Priority**: LOW - Administrative features
**Status**: ✅ **COMPLETE** - 3/11 datasources implemented, 8 not supported by NetBox API

1. ✅ `event_rule_data_source.go` - COMPLETE (fixed pattern)
2. ✅ `journal_entry_data_source.go` - COMPLETE (full implementation)
3. ⚠️ `config_context_data_source.go` - NOT SUPPORTED (no CustomFields in go-netbox)
4. ⚠️ `config_template_data_source.go` - NOT SUPPORTED (no CustomFields in go-netbox)
5. ⚠️ `custom_field_data_source.go` - NOT SUPPORTED (no CustomFields in go-netbox)
6. ⚠️ `custom_link_data_source.go` - NOT SUPPORTED (no CustomFields in go-netbox)
7. ⚠️ `export_template_data_source.go` - NOT SUPPORTED (no CustomFields in go-netbox)
8. ⚠️ `notification_group_data_source.go` - NOT SUPPORTED (no CustomFields in go-netbox)
9. ⚠️ `script_data_source.go` - NOT SUPPORTED (no CustomFields in go-netbox)
10. ✅ `webhook_data_source.go` - COMPLETE (full implementation + fixed resource)
11. ⚠️ `tag_data_source.go` - NOT SUPPORTED (no CustomFields in go-netbox)

**Test Priority**: LOW - ✅ All 3 supported datasources tested and passing
**Total Test Time**: ~5 seconds for all tests
**Implementation Time**: ~4 hours (including webhook resource enhancement)
**Bonus**: Fixed webhook resource to support custom fields with preservation test

**Note**: Most Extras/Admin objects don't support custom fields in NetBox API

### Batch 8: Templates & Miscellaneous (10 datasources)
**Priority**: LOW - Template resources (likely no custom fields)
**Status**: 🔧 IN PROGRESS

1. ⏳ `console_port_template_data_source.go`
2. ⏳ `console_server_port_template_data_source.go`
3. ⏳ `device_bay_template_data_source.go`
4. ⏳ `front_port_template_data_source.go`
5. ⏳ `interface_template_data_source.go`
6. `inventory_item_template_data_source.go`
7. `module_bay_template_data_source.go`
8. `module_type_data_source.go`
9. `power_outlet_template_data_source.go`
10. `power_port_template_data_source.go`

Plus remaining datasources:
11. `rack_type_data_source.go`
12. `region_data_source.go`
13. `site_group_data_source.go`
14. `virtual_chassis_data_source.go`
15. `tunnel_termination_data_source.go`
16. Various assignment and grouping datasources

**Test Priority**: VERY LOW - May not have custom fields, verify first

## Testing Strategy

### Test Categories

#### 1. Full Custom Field Tests (High Priority Datasources)
Pattern: Create resource with custom field, read via datasource, verify all fields present

```go
func TestAccXxxDataSource_customFields(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Steps: []resource.TestStep{
            {
                // Create resource with custom field via resource
                Config: testAccXxxDataSourceConfig_withCustomFields(),
                Check: resource.ComposeTestCheckFunc(
                    // Verify datasource returns custom field
                    resource.TestCheckResourceAttr("data.netbox_xxx.test", "custom_fields.#", "1"),
                    resource.TestCheckResourceAttr("data.netbox_xxx.test", "custom_fields.0.name", "test_field"),
                    resource.TestCheckResourceAttr("data.netbox_xxx.test", "custom_fields.0.value", "test_value"),
                ),
            },
        },
    })
}
```

#### 2. Multiple Field Tests
Verify datasource returns multiple custom fields of different types

#### 3. External Field Addition Tests
1. Create resource without custom fields via Terraform
2. Add custom field via NetBox API directly
3. Refresh datasource
4. Verify datasource picks up the new field

### Test Infrastructure Setup

#### Create test_main_test.go
```go
//go:build customfields

package datasources_acceptance_tests_customfields

import (
    "os"
    "testing"

    "github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

func TestMain(m *testing.M) {
    testutil.TestAccPreCheck(func() {})
    os.Exit(m.Run())
}
```

#### Run Commands
```bash
# Run parallel tests (no custom fields)
go test ./internal/datasources_acceptance_tests/... -v

# Run serial custom field tests
go test -tags=customfields ./internal/datasources_acceptance_tests_customfields/... -v -p 1
```

## Implementation Checklist

### Phase 1: Foundation (Week 1) ✅ COMPLETE
- ✅ Create `MapAllCustomFieldsToModels()` function in utils/common.go
- ✅ Add unit tests for `MapAllCustomFieldsToModels()`
- ✅ Create `internal/datasources_acceptance_tests_customfields/` directory
- ✅ Create `test_main_test.go` with build tag
- ✅ Implemented with `go test -tags=customfields` command

### Phase 2: Batch 1 Implementation (Week 1-2) ✅ COMPLETE
- ✅ Update 13 core infrastructure datasources
- ✅ Create custom field tests for all Batch 1 datasources
- ✅ Run tests and verify all passing (13/13  ✅ COMPLETE
- ✅ Update 12 device/DCIM datasources
- ✅ Create custom field tests for all Batch 2 datasources
- ✅ Run tests and verify all passing (12/12 passing)
- ✅ Commit and push Batch 2

### Phase 3.5: Batch 3 Implementation (Week 2) ✅ COMPLETE
- ✅ Update 10 virtualizatio

### Phase 4: Batch 4 Implementation (Week 2) ✅ COMPLETE
- ✅ Update 12 circuits/VPN datasources
- ✅ Create custom field tests for all Batch 4 datasources
- ✅ Run tests and verify al

### Phase 5: Batch 5 Implementation (Week 3) ✅ COMPLETE
- ✅ Update 12 ports/interfaces datasources
- ✅ Create custom field tests for all Batch 5 datasources
- ✅ Run tests and verify all passing (12/12 passing)
- ✅ Fix critical power_outlet resource bug
- ✅ Commit and push Batch 5l passing
- [ ] Commit and push Batch 2

### Phase 4: Batch 3-8 Implementation (Week 3-4)
- [ ] Complete remaining batches (60+ datasources)
- [ ] Create sample tests for lower priority datasources
- [ ] Run full test suite
- [ ] Update documentation

### Phase 5: Documentation & Release (Week 4)
- [ ] Update CHANGELOG.md
- [ ] Update datasource documentation examples
- [ ] Update provider README
- [ ] Create migration guide
- [ ] Prepare release notes

## Success Criteria

### Batch 1 Status
1. ✅ All datasources return complete custom field data (13/13 complete)
2. ✅ New helper function correctly handles all custom field types
3. ✅ Test suite verifies 1-3 (35/80 datasources = 44%)
- **Remaining**: Batches 4-8 (45 datasources)
- **On Track**: Yes, significantlypdated with examples (planned for final phase)
6. ✅ All tests passinges 1-4 (47/80 datasources = 59%)
- **Remaining**: Batches 5-8 (33 datasources)
- **On Track**: Yes, significantly
- **Completed**: Batch 1 (13/80 datasources = 16%)
- **Remaining**: Batches 2-8 (67 datasources)
- **On Track**: Yes, ahead of schedule

## Risk Assessment

### Low Risk
- Datasources are read-only operations
- No state management complexity
- Clear separation from resource partial management

### Medium Risk
- Large number of files to update (80+ datasources)
- Need to verify each datasource has custom_fields attribute
- Test infrastructure needs careful setup

### Mitigation
- Batch implementation allows incremental validation
- Separate test suite prevents interference
- Comprehensive testing at each batch

## Notes

### Differences from Resources
- **Resources**: Use partial management (filter-to-owned pattern)
- **Datasources**: Return ALL fields (complete read pattern)

### Type Inference
Since datasources don't have prior state, we must infer custom field types from the API response value type. This is acceptable for read-only operations.

### Ordering
Custom fields will be sorted alphabetically by name for consistent output in tests and state.

## Implementation Patterns Discovered

### Two Types ✅ Foundation + Batch 1 (13 datasources) - **COMPLETE**
  - Actual time: 4 hours
  - Pattern: 8 simple fixes + 5 full implementations
- **Week 2**: Batch 2 + Batch 3 (22 datasources) - IN PROGRESS
- **Week 3**: Batch 4-6 (34 datasources)
- **Week 4**: Batch 7-8 + Documentation (31 datasources)

**Total Estimated Time**: 4 weeks for complete implementation
**Actual Batch 1 Time**: 4 hours (faster than estimated due to established patterns)
```go
// BEFORE
customFields := utils.MapToCustomFieldModels(obj.GetCustomFields(), nil)

// AFTER
customFields := utils.MapAllCustomFieldsToModels(obj.GetCustomFields())
```

**Examples**: site, asn, circuit, circuit_type, cluster, cluster_type, vrf

#### 2. Full Implementation (5 datasources in Batch 1)
**When**: Datasource completely missing custom fields support
**Required Changes**:

1. **Add to Model**:
```go
type XxxDataSourceModel struct {
    // ... existing fields
    CustomFields types.Set `tfsdk:"custom_fields"`
}
```

2. **Add to Schema**:
```go
"custom_fields": nbschema.DSCustomFieldsAttribute(),
```

3. **Add to Read Logic** (in mapToState function):
```go
// Handle custom fields - datasources return ALL fields
if obj.HasCustomFields() {
    customFields := utils.MapAllCustomFieldsToModels(obj.GetCustomFields())
    customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
    if !cfDiags.HasError() {
        data.CustomFields = customFieldsValue
    }
} else {
    data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
}
```

**Examples**: asn_range, ip_address, ip_range, prefix, vlan, route_target

### Critical Implementation Details

#### Diagnostics Handling in mapToState Functions
**Issue**: mapToState functions don't have `resp` parameter
**Solution**: Check diagnostics inline without returning

```go
// ❌ WRONG - resp doesn't exist in mapToState
customFieldsValue, cfDiags := types.SetValueFrom(...)
resp.Diagnostics.Append(cfDiags...)
if resp.Diagnostics.HasError() {
    return
}

// ✅ CORRECT - Check diagnostics inline
customFieldsValue, cfDiags := types.SetValueFrom(...)
if !cfDiags.HasError() {
    data.CustomFields = customFieldsValue
}
```

#### Test Cleanup Registration
**Pattern**: Use testutil.NewCleanupResource() and register all created resources
```go
cleanup := testutil.NewCleanupResource(t)
cleanup.RegisterXxxCleanup(identifier)  // Check function signature!
cleanup.RegisterCustomFieldCleanup(customFieldName)
```

**Common Gotcha**: Some cleanup functions take different parameters than expected
- `RegisterIPRangeCleanup(startAddress)` - only takes start address, not end
- `RegisterClusterTypeCleanup(slug)` - takes slug, not name
- Always check testutil implementation before calling

#### Custom Field object_types by Resource
Different resources require different `object_types` values:
- IPAM: `ipam.prefix`, `ipam.ipaddress`, `ipam.iprange`, `ipam.vlan`, `ipam.vrf`, `ipam.asnrange`, `ipam.routetarget`
- Virtualization: `virtualization.cluster`, `virtualization.clustertype`
- Circuits: `circuits.circuit`, `circuits.circuittype`, `circuits.provider`
- DCIM: (to be documented in Batch 2)

#### Test Resource Configuration
**Pattern**: Use list syntax for custom_fields in resources
```go
resource "netbox_xxx" "test" {
  # ... other attributes
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-value"
    }
  ]
}
```

**NOT** the map syntax (older pattern):
```go
custom_fields = {
  (netbox_custom_field.test.name) = "test-value"  # ❌ Don't use this
}
```

### Performance Observations
- Individual datasource test runtime: 1.1-1.5 seconds
- Batch of 3 tests: ~3-4 seconds total
- Serial execution required for custom field tests (avoids race conditions)
- Cleanup warnings are normal (resources already deleted by test framework)

## Timeline

- **Week 1**: Foundation + Batch 1 (13 datasources)
- **Week 2**: Batch 2 + Batch 3 (22 datasources)
- **Week 3**: Batch 4-6 (34 datasources)
- **Week 4**: Batch 7-8 + Documentation (31 datasources)

**Total Estimated Time**: 4 weeks for complete implementation

## Questions for Review

1. Should we add a datasource-specific attribute like `include_custom_fields` to allow filtering?
   - **Decision**: No, datasources should be complete by default

2. How to handle custom field type ambiguity without schema?
   - **Decision**: Infer from value type, default to "text" if unclear

3. Should we backport this to v0.0.13 or target v0.0.14?
   - **Decision**: Target v0.0.14 as separate feature enhancement
