# Acceptance Test Quality Improvement Project

## Overview
Systematic improvement of acceptance tests across the terraform-provider-netbox project to ensure consistency, proper resource cleanup, and maintainable code.

## Checks and Improvements

### 1. Missing Cleanup Registration
**Issue**: Tests without cleanup registration don't properly clean up resources if the test fails.

**Detection**: Look for test functions that:
- Create resources via NetBox API
- Don't call `cleanup := testutil.NewCleanupResource(t)`
- Don't register resources with appropriate cleanup methods

**Fix**: Add cleanup registration at the start of each test:
```go
cleanup := testutil.NewCleanupResource(t)
cleanup.RegisterSiteCleanup(slug)
cleanup.RegisterContactCleanup(email)
// etc.
```

### 2. Incorrect External Deletion Test Pattern
**Issue**: External deletion tests were duplicating config in the second step instead of using `RefreshState`.

**Detection**: Look for external deletion tests with pattern:
```go
{
    PreConfig: func() { /* delete resource via API */ },
    Config: duplicateConfigFromStep1,
    Check: resource.ComposeTestCheckFunc(...),
}
```

**Correct Pattern**:
```go
{
    PreConfig: func() { /* delete resource via API */ },
    RefreshState:       true,
    ExpectNonEmptyPlan: true,
}
```

### 3. Duplicate Config Generation Functions
**Issue**: Multiple identical config functions with different names create code duplication and maintenance burden.

**Detection**:
- Compare config generation functions within the same test file
- Look for functions that generate identical Terraform configurations
- Common patterns:
  - `testAccXXXConsistencyLiteralNamesConfig` often duplicates `testAccXXXResourceConfig_basic` or `_full`
  - `testAccXXXResourceConfig_import` often duplicates other configs
  - `testAccXXXResourceConfig_withDescription` may duplicate `_full` with hardcoded values

**Fix**:
- Remove duplicate function definition
- Update all references to use the remaining function
- Verify both functions are truly identical before removing

## Completed Files (14)

### Circuit Resources
- ✅ circuit_type_resource_test.go
- ✅ circuit_termination_resource_test.go
- ✅ circuit_resource_test.go
- ✅ circuit_group_resource_test.go
- ✅ circuit_group_assignment_resource_test.go
- ✅ cable_resource_test.go

### Cluster Resources
- ✅ cluster_group_resource_test.go
- ✅ cluster_resource_test.go
- ✅ cluster_type_resource_test.go

### Config Resources
- ✅ config_context_resource_test.go
- ✅ config_template_resource_test.go

### Console Resources
- ✅ console_port_template_resource_test.go
- ✅ console_server_port_resource_test.go
- ✅ console_server_port_template_resource_test.go

### Contact Resources
- ✅ contact_assignment_resource_test.go

## Remaining Files (85) - Organized by Category

### Batch 1: Contact Resources (2 files)
- contact_resource_test.go
- contact_role_resource_test.go
- contact_group_resource_test.go

### Batch 2: Console & Device Bay Resources (3 files)
- console_port_resource_test.go
- device_bay_resource_test.go
- device_bay_template_resource_test.go

### Batch 3: Custom Field & Link Resources (3 files)
- custom_field_resource_test.go
- custom_field_choice_set_resource_test.go
- custom_link_resource_test.go

### Batch 4: Device Resources (3 files)
- device_resource_test.go
- device_role_resource_test.go
- device_type_resource_test.go

### Batch 5: Export & Template Resources (2 files)
- export_template_resource_test.go
- service_template_resource_test.go

### Batch 6: Front & Rear Port Resources (4 files)
- front_port_resource_test.go
- front_port_template_resource_test.go
- rear_port_resource_test.go
- rear_port_template_resource_test.go

### Batch 7: FHRP Resources (2 files)
- fhrp_group_resource_test.go
- fhrp_group_assignment_resource_test.go

### Batch 8: IKE & IPSec Resources (5 files)
- ike_policy_resource_test.go
- ike_proposal_resource_test.go
- ipsec_policy_resource_test.go
- ipsec_profile_resource_test.go
- ipsec_proposal_resource_test.go

### Batch 9: Interface Resources (2 files)
- interface_resource_test.go
- interface_template_resource_test.go

### Batch 10: Inventory Item Resources (3 files)
- inventory_item_resource_test.go
- inventory_item_role_resource_test.go
- inventory_item_template_resource_test.go

### Batch 11: IP Address & Range Resources (3 files)
- ip_address_resource_test.go
- ip_range_resource_test.go
- prefix_resource_test.go

### Batch 12: ASN & Aggregate Resources (4 files)
- asn_resource_test.go
- asn_range_resource_test.go
- aggregate_resource_test.go
- rir_resource_test.go

### Batch 13: Journal & Location Resources (2 files)
- journal_entry_resource_test.go
- location_resource_test.go

### Batch 14: L2VPN Resources (2 files)
- l2vpn_resource_test.go
- l2vpn_termination_resource_test.go

### Batch 15: Manufacturer & Platform Resources (2 files)
- manufacturer_resource_test.go
- platform_resource_test.go

### Batch 16: Module Resources (4 files)
- module_resource_test.go
- module_type_resource_test.go
- module_bay_resource_test.go
- module_bay_template_resource_test.go

### Batch 17: Power Resources (6 files)
- power_feed_resource_test.go
- power_outlet_resource_test.go
- power_outlet_template_resource_test.go
- power_panel_resource_test.go
- power_port_resource_test.go
- power_port_template_resource_test.go

### Batch 18: Provider Resources (3 files)
- provider_resource_test.go
- provider_account_resource_test.go
- provider_network_resource_test.go

### Batch 19: Rack Resources (5 files)
- rack_resource_test.go
- rack_role_resource_test.go
- rack_type_resource_test.go
- rack_reservation_resource_test.go
- region_resource_test.go

### Batch 20: Route & Role Resources (2 files)
- route_target_resource_test.go
- role_resource_test.go

### Batch 21: Service Resources (2 files)
- service_resource_test.go
- service_template_resource_test.go (duplicate - see Batch 5)

### Batch 22: Site Resources (2 files)
- site_resource_test.go
- site_group_resource_test.go

### Batch 23: Tag & Tenant Resources (3 files)
- tag_resource_test.go
- tenant_resource_test.go
- tenant_group_resource_test.go

### Batch 24: Tunnel Resources (3 files)
- tunnel_resource_test.go
- tunnel_group_resource_test.go
- tunnel_termination_resource_test.go

### Batch 25: Virtual Resources (5 files)
- virtual_chassis_resource_test.go
- virtual_device_context_resource_test.go
- virtual_disk_resource_test.go
- virtual_machine_resource_test.go
- vm_interface_resource_test.go

### Batch 26: VLAN & VRF Resources (3 files)
- vlan_resource_test.go
- vlan_group_resource_test.go
- vrf_resource_test.go

### Batch 27: Webhook & Wireless Resources (4 files)
- webhook_resource_test.go
- wireless_lan_resource_test.go
- wireless_lan_group_resource_test.go
- wireless_link_resource_test.go

## Workflow for Each Batch

1. **Review**: Open the test file and check for:
   - Missing cleanup registrations
   - Incorrect external deletion patterns
   - Duplicate config functions

2. **Fix**: Apply all necessary improvements using multi_replace_string_in_file for efficiency

3. **Build**: Run `go build .` to verify compilation

4. **Test**: Run acceptance tests for the modified file:
   ```powershell
   $env:TF_ACC="1"; go test ./internal/resources_acceptance_tests/... -run TestAccXXXResource -v -timeout 120m
   ```

5. **Commit**: Commit changes with descriptive message:
   ```
   git add -A
   git commit -m "Add cleanup registration and fix patterns in XXX tests"
   ```

6. **Repeat**: Move to next batch

## Success Metrics

- ✅ 100% build success rate
- ✅ 100% pre-commit hook success rate
- ✅ 100% test pass rate
- 🎯 All 99 resource test files improved
- 📊 Current Progress: 14/99 (14%)

## Notes

- External deletion tests already using RefreshState pattern don't need changes
- Some tests may not have external deletion tests at all (that's OK)
- Cleanup registration is the most critical improvement
- Duplicate removal has saved ~500+ lines of code so far
