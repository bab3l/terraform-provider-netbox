# Import Reference ID Validation - Implementation Plan

## Overview

This document outlines the plan to validate that reference fields contain numeric IDs (not names/slugs) after import. This ensures imported resources match typical Terraform configurations that use `resource.id` references.

## Background

**Problem**: After `terraform import`, reference fields (role, tenant, platform, etc.) were being stored as names instead of numeric IDs. When a Terraform configuration uses ID references (e.g., `tenant = netbox_tenant.example.id`), this caused unnecessary plan diffs requiring an immediate `apply`.

**Fix Applied**: Modified `UpdateReferenceAttribute` in `state_helpers.go` to prefer ID when state is null/unknown (import scenario).

**This Plan**: Add acceptance tests to validate the fix across all affected resources and prevent regressions.

## Affected Resources

Resources with reference fields (using `UpdateReferenceAttribute`):

| Resource | Ref Count | Reference Fields |
|----------|-----------|-----------------|
| virtual_machine | 5 | site, cluster, role, tenant, platform |
| tunnel | 9 | group, ipsec_profile, tenant (×3 operations) |
| prefix | 5 | site, vrf, tenant, vlan, role |
| vlan | 4 | site, group, tenant, role |
| vm_interface | 3 | virtual_machine, untagged_vlan, vrf |
| location | 3 | site, parent, tenant |
| platform | 3 | manufacturer (×3 operations) |
| inventory_item | 3 | device, role, manufacturer |
| circuit_termination | 3 | circuit, site, provider_network |
| aggregate | 2 | rir, tenant |
| asn | 2 | rir, tenant |
| asn_range | 2 | rir, tenant |
| circuit_group_assignment | 2 | group, circuit |
| console_port_template | 2 | device_type, module_type |
| console_server_port_template | 2 | device_type, module_type |
| device_bay | 2 | device, installed_device |
| device_type | 2 | manufacturer, default_platform |
| front_port | 2 | device, rear_port |
| front_port_template | 2 | device_type, module_type |
| interface_template | 2 | device_type, module_type |
| ip_address | 2 | vrf, tenant |
| module | 2 | device, module_type |
| module_bay_template | 2 | device_type, module_type |
| power_feed | 2 | rack, tenant |
| power_outlet_template | 2 | device_type, module_type |
| power_panel | 2 | site, location |
| power_port_template | 2 | device_type, module_type |
| rack_reservation | 2 | rack, (user - special) |
| rear_port_template | 2 | device_type, module_type |
| service | 2 | device, virtual_machine |
| virtual_device_context | 2 | device, tenant |
| wireless_lan | 2 | vlan, tenant |
| cable | 1 | tenant |
| circuit_group | 1 | tenant |
| console_port | 1 | device |
| console_server_port | 1 | device |
| contact | 1 | group |
| device_bay_template | 1 | device_type |
| l2vpn | 1 | tenant |
| module_bay | 1 | device |
| module_type | 1 | manufacturer |
| power_outlet | 1 | device |
| power_port | 1 | device |
| rack_type | 1 | manufacturer |
| rear_port | 1 | device |
| route_target | 1 | tenant |
| virtual_disk | 1 | virtual_machine |

**Total**: 47 resources with 98 reference field usages

## Implementation Strategy

### Approach: Extend Existing Import Tests

Rather than creating new tests, we'll extend the existing `_import` tests to validate numeric IDs. This:
1. Avoids test duplication
2. Validates import behavior where it's already tested
3. Keeps test files focused and maintainable

### New Helper Function

Helper functions are implemented in `internal/testutil/import_tests.go`:

```go
// ReferenceFieldCheck creates a check function that validates a reference field
// contains a numeric ID after import.
func ReferenceFieldCheck(resourceRef, fieldName string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs := s.RootModule().Resources[resourceRef]
        if rs == nil {
            return fmt.Errorf("resource %s not found", resourceRef)
        }
        value := rs.Primary.Attributes[fieldName]
        if value == "" {
            return nil // Field not set, nothing to validate
        }
        if _, err := strconv.Atoi(value); err != nil {
            return fmt.Errorf("%s should be numeric ID after import, got: %q", fieldName, value)
        }
        return nil
    }
}

// ValidateReferenceIDs creates check functions for multiple reference fields.
func ValidateReferenceIDs(resourceRef string, fields ...string) []resource.TestCheckFunc {
    checks := make([]resource.TestCheckFunc, len(fields))
    for i, field := range fields {
        checks[i] = ReferenceFieldCheck(resourceRef, field)
    }
    return checks
}
```

### Update RunImportTest Helper

Extend `ImportTestConfig` to accept reference fields to validate:

```go
type ImportTestConfig struct {
    // ... existing fields ...

    // ReferenceFields lists reference attributes that should contain numeric IDs after import.
    // These will be validated to ensure they're numeric (not names/slugs).
    ReferenceFields []string
}
```

## Batched Implementation Plan

## Progress (as of 2026-01-21)

- ✅ Helper functions added: `ReferenceFieldCheck`, `ValidateReferenceIDs`
- ✅ `REQUIRED_TESTS.md` updated to require numeric ID validation
- ✅ Phase 1 started: `virtual_machine`, `prefix`, `vlan`, `ip_address` updated
- ✅ Phase 1 tests run (subset): `VirtualMachine`, `Prefix`, `VLAN`, `IPAddress` import tests passed
- ✅ Phase 1 completed: `location`, `device_type`, `aggregate`, `vm_interface`, `tunnel` updated
- ✅ Phase 1 tests run (subset): `Location`, `DeviceType`, `Aggregate`, `VMInterface`, `Tunnel` import tests passed
- ✅ Phase 2 completed: `console_port`, `console_port_template`, `console_server_port`, `console_server_port_template`, `front_port`, `front_port_template`, `rear_port`, `rear_port_template`, `power_port`, `power_port_template`, `power_outlet`, `power_outlet_template`, `interface_template`, `device_bay` updated
- ✅ Phase 2 tests run (subset): DCIM port/template import tests passed
- ✅ Phase 3 completed: `module`, `module_bay`, `module_bay_template`, `module_type`, `inventory_item`, `device_bay_template` updated
- ✅ Phase 3 tests run (subset): module/inventory import tests passed
- ✅ Phase 4 completed: `circuit_termination`, `circuit_group_assignment`, `circuit_group`, `l2vpn`, `asn`, `asn_range` updated
- ✅ Phase 4 tests run (subset): circuit/vpn/asn import tests passed
- ✅ Phase 5 completed: `power_panel`, `power_feed`, `rack_reservation`, `rack_type`, `platform` updated
- ✅ Phase 5 tests run (subset): power/rack/platform import tests passed
- ✅ Phase 5 tests run (subset): power/rack/platform import tests passed

### Phase 0: Foundation (Current Session)
- [x] Fix `UpdateReferenceAttribute` to prefer ID for import
- [x] Add `TestAccVirtualMachineResource_ImportCommandRequiresApply` with assertions
- [x] Create feature branch
- [x] Add helper functions to `import_tests.go`
- [x] Update `REQUIRED_TESTS.md`

### Phase 1: Core/High-Impact Resources
**Gate**: All Phase 1 tests must pass before proceeding

Resources (9):
1. [x] `virtual_machine` - 5 refs
2. [x] `prefix` - 5 refs
3. [x] `vlan` - 4 refs
4. [x] `ip_address` - 2 refs
5. [x] `location` - 3 refs
6. [x] `device_type` - 2 refs
7. [x] `aggregate` - 2 refs
8. [x] `vm_interface` - 3 refs
9. [x] `tunnel` - 9 refs

**Validation**:
```powershell
go test ./internal/resources_acceptance_tests -run "TestAcc(VirtualMachine|Prefix|VLAN|IPAddress|Location|DeviceType|Aggregate|VMInterface|Tunnel)Resource.*[Ii]mport" -v -timeout 60m
```

### Phase 2: DCIM Port/Template Resources
**Gate**: Phase 1 complete, Phase 2 tests pass

Resources (14):
1. [x] `console_port` - 1 ref
2. [x] `console_port_template` - 2 refs
3. [x] `console_server_port` - 1 ref
4. [x] `console_server_port_template` - 2 refs
5. [x] `front_port` - 2 refs
6. [x] `front_port_template` - 2 refs
7. [x] `rear_port` - 1 ref
8. [x] `rear_port_template` - 2 refs
9. [x] `power_port` - 1 ref
10. [x] `power_port_template` - 2 refs
11. [x] `power_outlet` - 1 ref
12. [x] `power_outlet_template` - 2 refs
13. [x] `interface_template` - 2 refs
14. [x] `device_bay` - 2 refs

### Phase 3: Module & Inventory Resources
**Gate**: Phase 2 complete, Phase 3 tests pass

Resources (6):
1. [x] `module` - 2 refs
2. [x] `module_bay` - 1 ref
3. [x] `module_bay_template` - 2 refs
4. [x] `module_type` - 1 ref
5. [x] `inventory_item` - 3 refs
6. [x] `device_bay_template` - 1 ref

**Validation**:
```powershell
go test ./internal/resources_acceptance_tests -run "TestAcc(ModuleResource_import|ModuleBayResource_import|ModuleBayTemplateResource_import|ModuleTypeResource_import|InventoryItemResource_import|DeviceBayTemplateResource_import)$" -v -timeout 60m
```

### Phase 4: Circuit & VPN Resources
**Gate**: Phase 3 complete, Phase 4 tests pass

Resources (6):
1. [x] `circuit_termination` - 3 refs
2. [x] `circuit_group_assignment` - 2 refs
3. [x] `circuit_group` - 1 ref
4. [x] `l2vpn` - 1 ref
5. [x] `asn` - 2 refs
6. [x] `asn_range` - 2 refs

**Validation**:
```powershell
go test ./internal/resources_acceptance_tests -run "TestAcc(CircuitTerminationResource_import|CircuitGroupAssignmentResource_import|CircuitGroupResource_import|L2VPNResource_import|ASNResource_import|ASNRangeResource_import)$" -v -timeout 60m
```

### Phase 5: Power & Rack Resources
**Gate**: Phase 4 complete, Phase 5 tests pass

Resources (5):
1. [x] `power_panel` - 2 refs
2. [x] `power_feed` - 2 refs
3. [x] `rack_reservation` - 2 refs
4. [x] `rack_type` - 1 ref
5. [x] `platform` - 3 refs

**Validation**:
```powershell
go test ./internal/resources_acceptance_tests -run "TestAcc(PowerPanelResource_basic|PowerFeedResource_basic|RackReservationResource_basic|RackTypeResource_basic|PlatformResource_import)$" -v -timeout 60m
```

### Phase 6: Remaining Resources
**Gate**: Phase 5 complete, Phase 6 tests pass

Resources (7):
1. `cable` - 1 ref
2. `contact` - 1 ref
3. `route_target` - 1 ref
4. `service` - 2 refs
5. `virtual_device_context` - 2 refs
6. `virtual_disk` - 1 ref
7. `wireless_lan` - 2 refs

## Gating Criteria

Before merging PR:

1. **Unit Tests**: `go test ./internal/utils -run TestUpdateReferenceAttribute -v` - ALL PASS
2. **Phase Gates**: Each phase's acceptance tests pass
3. **Full Test Suite**: `go test ./internal/resources_acceptance_tests/... -v -timeout 120m` - ALL PASS
4. **Code Review**: Changes reviewed for consistency

## Test Pattern

For each resource with reference fields, the `_import` test should include:

```go
{
    ResourceName:            "netbox_{resource}.test",
    ImportState:             true,
    ImportStateVerify:       true,
    ImportStateVerifyIgnore: []string{...},
},
// Add validation step after import
{
    Config: testAcc{Resource}ResourceConfig_import(...),
    Check: resource.ComposeTestCheckFunc(
        // Validate reference fields are numeric IDs
        testutil.ReferenceFieldCheck("netbox_{resource}.test", "tenant"),
        testutil.ReferenceFieldCheck("netbox_{resource}.test", "site"),
        // ... other reference fields
    ),
},
```

## Completion Checklist

- [ ] Helper functions added to `import_tests.go`
- [ ] `REQUIRED_TESTS.md` updated
- [ ] Phase 1 complete
- [x] Phase 2 complete
- [x] Phase 3 complete
- [x] Phase 4 complete
- [x] Phase 5 complete
- [ ] Phase 6 complete
- [ ] Full acceptance test suite passes
- [ ] PR created and reviewed
- [ ] Merged to main

## Notes

- Tests use configs with ID references (typical pattern)
- The fix is in `UpdateReferenceAttribute` - no per-resource changes needed
- Tests catch regressions if the logic is changed
- Some resources may need `ImportStateVerifyIgnore` for computed fields
