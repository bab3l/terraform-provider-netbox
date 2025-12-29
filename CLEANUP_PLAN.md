# Cleanup Plan: Code Formatting Fix

## Current Status

✅ **Phase 1 Batch 1.1**: Created formatting script
✅ **Phase 1 Batch 1.2 (partial)**: Fixed state_helpers.go (1265→974 lines, 44%→28% blank)

## Overview

This plan fixes excessive blank lines throughout the Go codebase (~45% blank lines → ~28%).

**Root Cause**: Files have blank lines after every `{` and before every `}`, plus excessive spacing in comment blocks.

**Solution**: PowerShell script `scripts/fix-formatting.ps1` removes excessive blank lines while preserving intentional paragraph breaks.

## Batch Plan

### Phase 1: Format All Go Files

**Batch 1.2: Utility files** (3 files)
- [x] state_helpers.go (DONE: 1265→974 lines, 44%→28% blank)
- [ ] request_helpers.go
- [ ] attributes.go

**Batch 1.3: DCIM basics** (10 files) ✅ DONE
- [x] device_resource.go, device_type_resource.go, device_role_resource.go
- [x] device_bay_resource.go, device_bay_template_resource.go
- [x] manufacturer_resource.go, platform_resource.go
- [x] site_resource.go, site_group_resource.go, location_resource.go

**Batch 1.4: Racks, power, modules** (10 files) ✅ DONE
- [x] rack_resource.go, rack_role_resource.go, rack_type_resource.go, rack_reservation_resource.go
- [x] power_panel_resource.go, power_feed_resource.go
- [x] module_resource.go, module_type_resource.go, module_bay_resource.go
- [x] interface_resource.go

**Batch 1.5: Device component ports** (10 files)
- [ ] console_port_resource.go, console_server_port_resource.go
- [ ] power_port_resource.go, power_outlet_resource.go
- [ ] front_port_resource.go, rear_port_resource.go
- [ ] inventory_item_resource.go, inventory_item_role_resource.go, inventory_item_template_resource.go
- [ ] module_bay_template_resource.go

**Batch 1.6: Component templates** (8 files)
- [ ] console_port_template_resource.go, console_server_port_template_resource.go
- [ ] power_port_template_resource.go, power_outlet_template_resource.go
- [ ] front_port_template_resource.go, rear_port_template_resource.go
- [ ] interface_template_resource.go
- [ ] service_resource.go

**Batch 1.7: IPAM core** (10 files)
- [ ] ip_address_resource.go, ip_range_resource.go
- [ ] prefix_resource.go, vrf_resource.go
- [ ] vlan_resource.go, vlan_group_resource.go
- [ ] aggregate_resource.go, asn_resource.go, asn_range_resource.go
- [ ] rir_resource.go

**Batch 1.8: IPAM & VPN** (10 files)
- [ ] route_target_resource.go
- [ ] l2vpn_resource.go, l2vpn_termination_resource.go
- [ ] tunnel_resource.go, tunnel_group_resource.go, tunnel_termination_resource.go
- [ ] ipsec_profile_resource.go, ipsec_policy_resource.go, ipsec_proposal_resource.go
- [ ] ike_policy_resource.go

**Batch 1.9: Circuits & providers** (10 files)
- [ ] ike_proposal_resource.go
- [ ] cable_resource.go
- [ ] circuit_resource.go, circuit_type_resource.go, circuit_group_resource.go
- [ ] circuit_group_assignment_resource.go, circuit_termination_resource.go
- [ ] provider_resource.go, provider_account_resource.go, provider_network_resource.go

**Batch 1.10: Virtualization & wireless** (10 files)
- [ ] cluster_resource.go, cluster_type_resource.go, cluster_group_resource.go
- [ ] virtual_machine_resource.go, vm_interface_resource.go, virtual_disk_resource.go
- [ ] virtual_chassis_resource.go, virtual_device_context_resource.go
- [ ] wireless_lan_resource.go, wireless_lan_group_resource.go

**Batch 1.11: Tenancy & contacts** (10 files)
- [ ] wireless_link_resource.go
- [ ] tenant_resource.go, tenant_group_resource.go
- [ ] contact_resource.go, contact_group_resource.go, contact_role_resource.go, contact_assignment_resource.go
- [ ] region_resource.go
- [ ] role_resource.go
- [ ] tag_resource.go

**Batch 1.12: Extras & remaining** (10 files)
- [ ] fhrp_group_resource.go, fhrp_group_assignment_resource.go
- [ ] service_template_resource.go
- [ ] config_template_resource.go, config_context_resource.go
- [ ] custom_field_resource.go, custom_field_choice_set_resource.go, custom_link_resource.go
- [ ] webhook_resource.go, event_rule_resource.go

**Batch 1.13: Final resources** (3 files)
- [ ] export_template_resource.go
- [ ] journal_entry_resource.go
- [ ] notification_group_resource.go

**Batch 1.14**: Format datasource files (~100 files)

**Batch 1.15**: Format test files (~100 files)

**Batch 1.16**: Verify and commit all formatting changes

### Phase 2: Remove DisplayNameAttribute and display_name Field

**Problem**: The `display_name` computed field was added to show friendly names for referenced resources, but it doesn't work as intended:
- Shows the resource's own Display field (e.g., "eth0" for interface named "eth0")
- Does NOT show friendly names for referenced resources (tenant, vlan, cluster, etc.)
- Causes "(known after apply)" noise in plans
- Adds complexity without solving the actual UX problem

**Files Affected**: 100 resource files + 1 utility file

**Batch 2.1: Remove DisplayNameAttribute function**
- [ ] Remove `DisplayNameAttribute()` from `internal/schema/attributes.go`
- [ ] This will break the build (expected) until we remove all usages

**Batches 2.2-2.13: Remove display_name from resources** (100 files in 12 batches)

For each resource file:
1. Remove `DisplayName types.String` from model struct
2. Remove `"display_name": nbschema.DisplayNameAttribute(...)` from Schema()
3. Remove `data.DisplayName = ...` from mapToState()

**Batch 2.2: IPAM resources** (9 resources)
- [ ] aggregate_resource.go
- [ ] asn_resource.go
- [ ] asn_range_resource.go
- [ ] ip_address_resource.go
- [ ] ip_range_resource.go
- [ ] prefix_resource.go
- [ ] rir_resource.go
- [ ] route_target_resource.go
- [ ] vrf_resource.go

**Batch 2.3: IPAM VLANs** (3 resources)
- [ ] vlan_resource.go
- [ ] vlan_group_resource.go
- [ ] wireless_link_resource.go

**Batch 2.4: Circuits** (7 resources)
- [ ] cable_resource.go
- [ ] circuit_resource.go
- [ ] circuit_group_resource.go
- [ ] circuit_group_assignment_resource.go
- [ ] circuit_termination_resource.go
- [ ] circuit_type_resource.go
- [ ] provider_resource.go

**Batch 2.5: Provider & VPN** (8 resources)
- [ ] provider_account_resource.go
- [ ] provider_network_resource.go
- [ ] tunnel_resource.go
- [ ] tunnel_group_resource.go
- [ ] tunnel_termination_resource.go
- [ ] l2vpn_resource.go
- [ ] l2vpn_termination_resource.go
- [ ] ipsec_profile_resource.go

**Batch 2.6: VPN & Security** (6 resources)
- [ ] ipsec_policy_resource.go
- [ ] ipsec_proposal_resource.go
- [ ] ike_policy_resource.go
- [ ] ike_proposal_resource.go
- [ ] fhrp_group_resource.go
- [ ] fhrp_group_assignment_resource.go

**Batch 2.7: Virtualization** (8 resources)
- [ ] cluster_resource.go
- [ ] cluster_group_resource.go
- [ ] cluster_type_resource.go
- [ ] virtual_machine_resource.go
- [ ] vm_interface_resource.go
- [ ] virtual_disk_resource.go
- [ ] virtual_chassis_resource.go
- [ ] virtual_device_context_resource.go

**Batch 2.8: DCIM - Sites & Locations** (7 resources)
- [ ] site_resource.go
- [ ] site_group_resource.go
- [ ] region_resource.go
- [ ] location_resource.go
- [ ] rack_resource.go
- [ ] rack_role_resource.go
- [ ] rack_type_resource.go

**Batch 2.9: DCIM - Devices** (9 resources)
- [ ] device_resource.go
- [ ] device_role_resource.go
- [ ] device_type_resource.go
- [ ] device_bay_resource.go
- [ ] device_bay_template_resource.go
- [ ] manufacturer_resource.go
- [ ] platform_resource.go
- [ ] module_resource.go
- [ ] module_type_resource.go

**Batch 2.10: DCIM - Device Components** (8 resources)
- [ ] interface_resource.go
- [ ] console_port_resource.go
- [ ] console_server_port_resource.go
- [ ] power_port_resource.go
- [ ] power_outlet_resource.go
- [ ] front_port_resource.go
- [ ] rear_port_resource.go
- [ ] module_bay_resource.go

**Batch 2.11: DCIM - Component Templates** (8 resources)
- [ ] interface_template_resource.go
- [ ] console_port_template_resource.go
- [ ] console_server_port_template_resource.go
- [ ] power_port_template_resource.go
- [ ] power_outlet_template_resource.go
- [ ] front_port_template_resource.go
- [ ] rear_port_template_resource.go
- [ ] module_bay_template_resource.go

**Batch 2.12: DCIM - Power & Inventory** (9 resources)
- [ ] power_panel_resource.go
- [ ] power_feed_resource.go
- [ ] rack_reservation_resource.go
- [ ] inventory_item_resource.go
- [ ] inventory_item_role_resource.go
- [ ] inventory_item_template_resource.go
- [ ] service_resource.go
- [ ] service_template_resource.go
- [ ] config_template_resource.go

**Batch 2.13: Tenancy, Tags, Extras** (18 resources)
- [ ] tenant_resource.go
- [ ] tenant_group_resource.go
- [ ] contact_resource.go
- [ ] contact_group_resource.go
- [ ] contact_role_resource.go
- [ ] contact_assignment_resource.go
- [ ] tag_resource.go
- [ ] role_resource.go
- [ ] custom_field_resource.go
- [ ] custom_field_choice_set_resource.go
- [ ] custom_link_resource.go
- [ ] webhook_resource.go
- [ ] event_rule_resource.go
- [ ] export_template_resource.go
- [ ] journal_entry_resource.go
- [ ] notification_group_resource.go
- [ ] wireless_lan_resource.go
- [ ] wireless_lan_group_resource.go

### Phase 3: Verify Reference Attribute Preservation

**Problem**: Reference attributes (tenant, cluster, vlan, etc.) were showing numeric IDs instead of friendly names in plans:
```
~ tenant = "My Tenant" -> "42"
```

**Fix Applied**: ✅ Updated `UpdateReferenceAttribute` in state_helpers.go to prefer name/slug over ID when state is Unknown.

**Affected Resources**: **47 resources** use `UpdateReferenceAttribute` for reference fields, including:
- DCIM: Device, Interface, Cable, Rack, PowerPanel, PowerFeed, Module, InventoryItem, DeviceBay, VirtualDeviceContext, ConsolePort, ConsoleServerPort, FrontPort, RearPort, PowerPort, PowerOutlet, ModuleBay, etc.
- IPAM: IPAddress, Prefix, VLAN, VRF, Aggregate, ASN, ASNRange, RouteTarget
- Virtualization: VirtualMachine, VMInterface, VirtualDisk, Cluster
- Circuits: Circuit (via CircuitTermination), CircuitTermination, CircuitGroupAssignment, ProviderNetwork
- Tenancy: ContactAssignment, Contact
- Wireless: WirelessLAN
- VPN: Tunnel, L2VPN
- All device/module component templates: ConsolePortTemplate, ConsoleServerPortTemplate, InterfaceTemplate, PowerPortTemplate, PowerOutletTemplate, FrontPortTemplate, RearPortTemplate, ModuleBayTemplate
- DeviceType, DeviceBayTemplate, RackType, ModuleType, Platform

**Verification Strategy**: Add **7 representative tests** to validate the fix works across different patterns (single ref, multiple refs, nested refs, optional refs). Since UpdateReferenceAttribute is used consistently, testing these patterns should confirm the fix works everywhere.

**Batch 3.1: Add reference preservation tests**
- [ ] Update `TestAccConsistency_VMInterface` - verify vlan reference preservation (already partially done)
- [ ] Add `TestAccConsistency_VirtualMachine` - verify cluster, tenant, role, platform, site references
- [ ] Add `TestAccConsistency_IPAddress` - verify tenant, vrf references
- [ ] Add `TestAccConsistency_Device` - verify site, tenant, role, device_type, platform references
- [ ] Add `TestAccConsistency_Prefix` - verify tenant, vrf, site references
- [ ] Add `TestAccConsistency_VLAN` - verify tenant, site, group references
- [ ] Add `TestAccConsistency_Cluster` - verify type, group, site, tenant references

Test pattern:
```go
Steps: []resource.TestStep{
    {
        Config: configWithNameReferences(...),
        Check: resource.ComposeTestCheckFunc(
            resource.TestCheckResourceAttr("netbox_resource.test", "tenant", "tenant-name"),
        ),
    },
    {
        PlanOnly: true,  // Verify no drift
        Config: configWithNameReferences(...),
    },
}
```

### Phase 4: Final Verification

**Batch 4.1: Run all tests**
- [ ] `go build .`
- [ ] `go vet ./...`
- [ ] Run unit tests
- [ ] Run key acceptance tests (Consistency tests)
- [ ] Run full acceptance test suite (optional, takes 2+ hours)

**Batch 4.2: Create PR**
- [ ] Update CHANGELOG.md
- [ ] Create PR with summary of changes
- [ ] Celebrate! 🎉

## Script Usage

```powershell
# Dry run (preview changes)
.\scripts\fix-formatting.ps1 -Path "file.go" -DryRun

# Apply changes
.\scripts\fix-formatting.ps1 -Path "file.go"

# Process multiple files
.\scripts\fix-formatting.ps1 -Path "internal/resources/device*.go"
```

## Success Criteria

- [x] Script reduces blank lines from ~45% to ~25-30%
- [x] Build succeeds after formatting
- [ ] All files formatted consistently
- [ ] No functionality changes (formatting only)
