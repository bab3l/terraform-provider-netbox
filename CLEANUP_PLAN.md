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

### Phase 2: Remove DisplayNameAttribute and display_name Field ✅ COMPLETE

**Problem**: The `display_name` computed field was added to show friendly names for referenced resources, but it doesn't work as intended:
- Shows the resource's own Display field (e.g., "eth0" for interface named "eth0")
- Does NOT show friendly names for referenced resources (tenant, vlan, cluster, etc.)
- Causes "(known after apply)" noise in plans
- Adds complexity without solving the actual UX problem

**Files Affected**: 100 resource files + 1 utility file

**Batch 2.1: Remove DisplayNameAttribute function** ✅ COMPLETE
- [x] Remove `DisplayNameAttribute()` from `internal/schema/attributes.go`
- [x] Removed all display_name field declarations from 100 resource files
- [x] Removed all display_name schema attributes from 100 resource files
- [x] Removed all display_name assignments from mapToState functions
- [x] Cleaned up 60 empty Display check blocks left as dead code

**Batches 2.2-2.13: Remove display_name from resources** (100 files in 12 batches) ✅ COMPLETE

For each resource file:
1. Remove `DisplayName types.String` from model struct ✅
2. Remove `"display_name": nbschema.DisplayNameAttribute(...)` from Schema() ✅
3. Remove `data.DisplayName = ...` from mapToState() ✅

All batches below completed successfully:

**Batch 2.2: IPAM resources** (9 resources) ✅
- [x] aggregate_resource.go
- [x] asn_resource.go
- [x] asn_range_resource.go
- [x] ip_address_resource.go
- [x] ip_range_resource.go
- [x] prefix_resource.go
- [x] rir_resource.go
- [x] route_target_resource.go
- [x] vrf_resource.go

**Batch 2.3: IPAM VLANs** (3 resources) ✅
- [x] vlan_resource.go
- [x] vlan_group_resource.go
- [x] wireless_link_resource.go

**Batch 2.4: Circuits** (7 resources) ✅
- [x] cable_resource.go
- [x] circuit_resource.go
- [x] circuit_group_resource.go
- [x] circuit_group_assignment_resource.go
- [x] circuit_termination_resource.go
- [x] circuit_type_resource.go
- [x] provider_resource.go

**Batch 2.5: Provider & VPN** (8 resources) ✅
- [x] provider_account_resource.go
- [x] provider_network_resource.go
- [x] tunnel_resource.go
- [x] tunnel_group_resource.go
- [x] tunnel_termination_resource.go
- [x] l2vpn_resource.go
- [x] l2vpn_termination_resource.go
- [x] ipsec_profile_resource.go

**Batch 2.6: VPN & Security** (6 resources) ✅
- [x] ipsec_policy_resource.go
- [x] ipsec_proposal_resource.go
- [x] ike_policy_resource.go
- [x] ike_proposal_resource.go
- [x] fhrp_group_resource.go
- [x] fhrp_group_assignment_resource.go

**Batch 2.7: Virtualization** (8 resources) ✅
- [x] cluster_resource.go
- [x] cluster_group_resource.go
- [x] cluster_type_resource.go
- [x] virtual_machine_resource.go
- [x] vm_interface_resource.go
- [x] virtual_disk_resource.go
- [x] virtual_chassis_resource.go
- [x] virtual_device_context_resource.go

**Batch 2.8: DCIM - Sites & Locations** (7 resources) ✅
- [x] site_resource.go
- [x] site_group_resource.go
- [x] region_resource.go
- [x] location_resource.go
- [x] rack_resource.go
- [x] rack_role_resource.go
- [x] rack_type_resource.go

**Batch 2.9: DCIM - Devices** (9 resources) ✅
- [x] device_resource.go
- [x] device_role_resource.go
- [x] device_type_resource.go
- [x] device_bay_resource.go
- [x] device_bay_template_resource.go
- [x] manufacturer_resource.go
- [x] platform_resource.go
- [x] module_resource.go
- [x] module_type_resource.go

**Batch 2.10: DCIM - Device Components** (8 resources) ✅
- [x] interface_resource.go
- [x] console_port_resource.go
- [x] console_server_port_resource.go
- [x] power_port_resource.go
- [x] power_outlet_resource.go
- [x] front_port_resource.go
- [x] rear_port_resource.go
- [x] module_bay_resource.go

**Batch 2.11: DCIM - Component Templates** (8 resources) ✅
- [x] interface_template_resource.go
- [x] console_port_template_resource.go
- [x] console_server_port_template_resource.go
- [x] power_port_template_resource.go
- [x] power_outlet_template_resource.go
- [x] front_port_template_resource.go
- [x] rear_port_template_resource.go
- [x] module_bay_template_resource.go

**Batch 2.12: DCIM - Power & Inventory** (9 resources) ✅
- [x] power_panel_resource.go
- [x] power_feed_resource.go
- [x] rack_reservation_resource.go
- [x] inventory_item_resource.go
- [x] inventory_item_role_resource.go
- [x] inventory_item_template_resource.go
- [x] service_resource.go
- [x] service_template_resource.go
- [x] config_template_resource.go

**Batch 2.13: Tenancy, Tags, Extras** (18 resources) ✅
- [x] tenant_resource.go
- [x] tenant_group_resource.go
- [x] contact_resource.go
- [x] contact_group_resource.go
- [x] contact_role_resource.go
- [x] contact_assignment_resource.go
- [x] tag_resource.go
- [x] role_resource.go
- [x] custom_field_resource.go
- [x] custom_field_choice_set_resource.go
- [x] custom_link_resource.go
- [x] webhook_resource.go
- [x] event_rule_resource.go
- [x] export_template_resource.go
- [x] journal_entry_resource.go
- [x] notification_group_resource.go
- [x] wireless_lan_resource.go
- [x] wireless_lan_group_resource.go

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

**Batch 3.1: Verify reference preservation tests** ✅ COMPLETE

Upon investigation, comprehensive tests already exist for most resources:
- [x] `TestAccConsistency_VMInterface` - verifies vlan reference preservation (EXISTING)
- [x] `TestAccConsistency_VirtualMachine_PlatformNamePersistence` - verifies platform reference (EXISTING)
- [x] `TestAccConsistency_Device_LiteralNames` - verifies site, tenant, role, device_type references (EXISTING)
- [x] `TestAccConsistency_Prefix` - verifies site, tenant, vlan references (EXISTING)
- [x] `TestAccConsistency_VLAN` - verifies site, group, tenant, role references (EXISTING)
- [x] `TestAccConsistency_Cluster_LiteralNames` - verifies type, group, site, tenant references (EXISTING)
- [x] `TestAccReferenceNamePersistence_IPAddress_TenantVRF` - verifies tenant, vrf references (ADDED)

These tests verify that reference attributes specified as names/slugs are preserved through refresh cycles and do not drift to IDs. The PlanOnly step in each test ensures no drift occurs.

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

**Batch 4.1: Run all tests** ✅ COMPLETE
- [x] `go build .` - Build successful
- [x] `go vet ./...` - No issues found
- [x] Run unit tests - All passed (resources, datasources, utils)
- [x] Run key acceptance tests - All 150+ consistency and reference preservation tests passed
- [ ] Run full acceptance test suite (optional, takes 2+ hours) - Skipped, key tests cover patterns

**Batch 4.2: Create PR** ✅ COMPLETE
- [x] Update CHANGELOG.md - Added v0.0.9 entry with all fixes
- [x] Create PR description - Comprehensive summary in PR_DESCRIPTION.md
- [x] Celebrate! 🎉

### Phase 5: Remove Duplicate `_id` Fields (v0.0.10 Breaking Change)

**Problem**: Many resources have duplicate reference fields that clutter state and plans:
- Primary field: `tenant` (accepts name/slug/ID, user-friendly)
- Duplicate field: `tenant_id` (computed, always shows numeric ID)

**Issues with duplicate fields:**
1. **Plan noise**: `cluster_id`, `tenant_id`, `platform_id`, etc. show `(known after apply)` during updates
2. **State bloat**: Duplicates the same information already stored internally
3. **Model complexity**: Doubles the number of reference fields
4. **Maintenance burden**: More code, tests, and documentation

**The primary field already supports cross-resource references:**
```hcl
resource "netbox_virtual_machine" "vm" {
  tenant = netbox_tenant.my_tenant.name  # Works with .name, .slug, or .id
}

resource "netbox_ip_address" "ip" {
  tenant = netbox_virtual_machine.vm.tenant  # References work fine!
}
```

**Decision**: Remove duplicate fields directly in v0.0.10 (alpha version, breaking changes acceptable)

**Affected Resources**: ~20 resources with ~35 duplicate `_id` fields:
- **device_resource.go**: device_type_id, role_id, tenant_id, platform_id, site_id, location_id, rack_id (7 fields)
- **virtual_machine_resource.go**: site_id, cluster_id, role_id, tenant_id, platform_id (5 fields)
- **rack_resource.go**: site_id, location_id, tenant_id, role_id, rack_type_id (5 fields)
- **vlan_resource.go**: site_id, tenant_id (2 fields)
- **platform_resource.go**: manufacturer_id (1 field)
- **rack_type_resource.go**: manufacturer_id (1 field)
- **vrf_resource.go**: tenant_id (1 field)
- **route_target_resource.go**: tenant_id (1 field)
- **tenant_resource.go**: group_id (1 field)
- **site_group_resource.go**: parent_id (1 field)
- **tenant_group_resource.go**: parent_id (1 field)
- **region_resource.go**: parent_id (1 field)
- **location_resource.go**: parent_id (1 field)
- **vlan_group_resource.go**: scope_id (1 field)
- **fhrp_group_assignment_resource.go**: group_id, interface_id (2 fields)
- **inventory_item_resource.go**: part_id (1 field)
- **inventory_item_template_resource.go**: part_id, component_id (2 fields)
- **provider_network_resource.go**: service_id (1 field)

**Migration Strategy:**
1. Remove all `_id` fields from model structs, schemas, and mapToState functions
2. Update CHANGELOG with clear migration guide for users
3. Document breaking changes and how to update configurations
4. Terraform automatically handles removed computed fields (state migration automatic)

**User Migration Required:**
Users who reference `netbox_virtual_machine.vm.tenant_id` will need to change to:
- `netbox_virtual_machine.vm.tenant` (returns same value, accepts name/slug/ID)
- Or reference the source directly: `netbox_tenant.my_tenant.id`

**Batch 5.1: Analyze and document** ✅ COMPLETE
- [x] Audit all resources for duplicate `_id` fields - Already documented in plan
- [x] Remove deprecation approach - Direct removal instead
- [x] Document migration guide for users (add to CHANGELOG)
- [x] Verify primary fields work correctly with cross-resource references

**Batch 5.2: Device and related resources** (7 duplicate fields) ✅ COMPLETE
- [x] device_resource.go
  - Removed from model: DeviceTypeID, RoleID, TenantID, PlatformID, SiteID, LocationID, RackID
  - Removed from schema: device_type_id, role_id, tenant_id, platform_id, site_id, location_id, rack_id
  - Removed from mapToState: all ID assignments
  - Kept: device_type, role, tenant, platform, site, location, rack

**Batch 5.3: Virtualization resources** (5 duplicate fields) ✅ COMPLETE
- [x] virtual_machine_resource.go
  - Removed from model: SiteID, ClusterID, RoleID, TenantID, PlatformID
  - Removed from schema: site_id, cluster_id, role_id, tenant_id, platform_id
  - Removed from mapToState: all ID assignments
  - Kept: site, cluster, role, tenant, platform

**Batch 5.4: Rack and infrastructure** (6 duplicate fields) ✅ COMPLETE
- [x] rack_resource.go
  - Removed from model: SiteID, LocationID, TenantID, RoleID, RackTypeID
  - Removed from schema: site_id, location_id, tenant_id, role_id, rack_type_id
  - Removed from mapToState: all ID assignments
  - Kept: site, location, tenant, role, rack_type
- [x] rack_type_resource.go
  - Removed from model: ManufacturerID
  - Removed from schema: manufacturer_id
  - Removed from mapToState: ID assignment
  - Kept: manufacturer

**Batch 5.5: IPAM resources** (4 duplicate fields) ✅ COMPLETE
- [x] vlan_resource.go
  - Removed from model: SiteID, TenantID
  - Removed from schema: site_id, tenant_id
  - Removed from mapToState: ID assignments
  - Kept: site, tenant
- [x] vrf_resource.go
  - Removed from model: TenantID
  - Removed from schema: tenant_id
  - Removed from mapToState: ID assignment
  - Kept: tenant
- [x] route_target_resource.go
  - Removed from model: TenantID
  - Removed from schema: tenant_id
  - Removed from mapToState: ID assignment
  - Kept: tenant

**Batch 5.6: Hierarchy resources** (5 duplicate fields) ✅ COMPLETE
- [x] site_group_resource.go - Removed: ParentID (model), parent_id (schema), mapToState assignment. Kept: parent
- [x] tenant_group_resource.go - Removed: ParentID (model), parent_id (schema), mapToState assignment. Kept: parent
- [x] tenant_resource.go - Removed: GroupID (model), group_id (schema), mapToState assignment. Kept: group
- [x] region_resource.go - Removed: ParentID (model), parent_id (schema), mapToState assignment. Kept: parent
- [x] location_resource.go - Removed: ParentID (model), parent_id (schema), mapToState assignment. Kept: parent

**Batch 5.7: Remaining resources** (1 duplicate field) ✅ COMPLETE
- [x] platform_resource.go - Removed: ManufacturerID. Kept: manufacturer

**Batch 5.8: Update tests and documentation** ✅ COMPLETE
- [x] Update all acceptance tests to not reference `_id` fields - No changes needed (tests use primary fields)
- [x] Verify no test breakage - All unit tests passing
- [x] Update CHANGELOG with breaking changes and migration guide - Added v0.0.10 entry

**Breaking Change Impact:**
- **High**: Users actively referencing `_id` fields in configs or outputs will need to update
- **Medium**: Plans will be cleaner (no more `(known after apply)` noise)
- **Low**: State migration is automatic (Terraform drops removed computed fields)

**Timeline:**
- v0.0.9: Fix display_name and reference preservation ✅
- v0.0.10: Remove all duplicate `_id` fields (breaking change, alpha version)

### Phase 6: Review Non-Standard ID Fields (Future Consideration)

**Note**: The following fields were identified in the original audit but don't follow the typical "duplicate computed ID" pattern. They require further investigation to determine if they should be removed or retained.

**vlan_group_resource.go - scope_id field:**
- Current: `scope_id` is OPTIONAL (not computed), used with `scope_type` to identify a scope object
- Pattern: Composite key (scope_type + scope_id) with no primary "scope" reference field
- Question: Is this a polymorphic reference pattern that should be kept, or can it be replaced?

**fhrp_group_assignment_resource.go - group_id, interface_id fields:**
- Current: Both `group_id` and `interface_id` are REQUIRED fields (not computed)
- Pattern: These are the primary way to specify relationships, not duplicates of other fields
- Question: Are these actually the primary fields, or should there be reference fields instead?

**inventory_item_resource.go - part_id field:**
- Current: `part_id` is OPTIONAL string field for manufacturer part numbers (e.g., "ABC-123")
- Pattern: Not a reference ID - it's a part number/SKU string
- Question: This appears to be business data, not a duplicate reference. Should it be kept?

**inventory_item_template_resource.go - part_id, component_id fields:**
- Current: `part_id` is optional part number string, `component_id` is optional
- Pattern: Similar to inventory_item - `part_id` is a part number, `component_id` needs investigation
- Question: Are these business data fields or actual duplicate references?

**provider_network_resource.go - service_id field:**
- Current: `service_id` is OPTIONAL string field for service identifiers from circuit provider
- Pattern: Not a reference ID - it's a service identifier/account number string
- Question: This appears to be business data (like an account number), not a duplicate reference. Should it be kept?

**Recommendation**: Review Netbox API documentation and existing Terraform configurations to determine:
1. Whether these fields serve different purposes than typical reference duplicates
2. If removing them would break legitimate use cases
3. Whether alternative patterns exist for these scenarios

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
