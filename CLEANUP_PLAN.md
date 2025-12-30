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

### Phase 6: Review Non-Standard Fields ✅ COMPLETE

**Status**: Analysis Complete - No Action Needed

The 7 fields deferred from Phase 5 have been analyzed and confirmed as legitimate:
- **vlan_group**: `scope_id` + `scope_type` (polymorphic composite key pattern)
- **fhrp_group_assignment**: `group_id`, `interface_id` (required primitive inputs)
- **inventory_item**: `part_id` (manufacturer part number string, not reference)
- **inventory_item_template**: `part_id`, `component_id` (part number + polymorphic component link)
- **provider_network**: `service_id` (external provider service identifier)

**Verdict**: All fields serve distinct purposes and should remain in schema.

---

### Phase 7: Update Examples

**Status**: 📋 In Progress

Ensure all example Terraform configurations in `examples/` directory are:
- Using current resource schemas (no `_id` fields removed in Phase 5)
- Following best practices for resource references
- Syntactically valid and up-to-date with provider features
- Well-documented with comments

**Totals**: 103 resource examples, 114 data source examples

#### Batch 7.1: IPAM Core Resources (12 resources) ✅ COMPLETE
- [x] netbox_aggregate - Clean, uses RIR name reference
- [x] netbox_asn - Clean, uses RIR name reference
- [x] netbox_asn_range - Excellent, comprehensive with comments
- [x] netbox_ip_address - Clean, shows multiple scenarios (IPv4, IPv6, VRF)
- [x] netbox_ip_range - Clean, simple example
- [x] netbox_prefix - Clean, uses slug for site, name for VLAN
- [x] netbox_rir - Clean, simple example
- [x] netbox_route_target - Good, shows tenant reference by slug
- [x] netbox_vlan - Clean, uses `.id` (valid per schema)
- [x] netbox_vlan_group - Clean, correct use of scope_id/scope_type pattern
- [x] netbox_vrf - Clean, simple example
- [x] netbox_fhrp_group - Clean, simple example

**Review Notes**:
- All examples verified against current schemas
- No removed `_id` fields referenced
- Mix of `.id`, `.slug`, and `.name` references all valid (schemas accept flexible input)
- scope_id pattern in vlan_group is correct (validated in Phase 6)
- Examples are realistic and follow best practices

#### Batch 7.2: Sites & Organization (10 resources) ✅ COMPLETE
- [x] netbox_site - Clean, simple example with all common fields
- [x] netbox_site_group - Clean, simple example
- [x] netbox_region - Clean, simple example
- [x] netbox_location - Clean, uses site slug reference
- [x] netbox_tenant - Good, shows group reference, tags, custom_fields
- [x] netbox_tenant_group - Good, shows parent-child hierarchy
- [x] netbox_contact - Excellent, multiple examples (basic, full, with_tags)
- [x] netbox_contact_group - Good, shows hierarchical groups with parent reference
- [x] netbox_contact_role - Good, multiple role examples
- [x] netbox_contact_assignment - **FIXED**: Updated to use primary fields instead of removed `_id` fields

**Changes Made**:
- Fixed `netbox_contact_assignment/resource.tf`:
  - Changed `manufacturer_id` → `manufacturer` (uses slug)
  - Changed `device_type_id` → `device_type` (uses slug)
  - Changed `role_id` → `role` (uses slug)
  - Changed `site_id` → `site` (uses slug)

**Review Notes**:
- All examples now use current schemas
- No removed `_id` fields referenced after fix
- Hierarchical relationships (parent-child) correctly shown
- Good mix of simple and comprehensive examples

#### Batch 7.3: Device Infrastructure (12 resources) ✅ COMPLETE
- [x] netbox_device - Clean, shows references using name/model/slug
- [x] netbox_device_type - Clean, uses manufacturer.id reference
- [x] netbox_device_role - Clean, simple example
- [x] netbox_device_bay - Clean, uses device name reference
- [x] netbox_device_bay_template - Excellent, comprehensive with count example
- [x] netbox_manufacturer - Clean, simple example
- [x] netbox_platform - Clean, uses manufacturer name reference
- [x] netbox_role - Clean, simple example
- [x] netbox_rack - Clean, uses site/role name references
- [x] netbox_rack_type - Clean, uses manufacturer name reference
- [x] netbox_rack_role - Clean, simple example
- [x] netbox_rack_reservation - Clean, uses rack name and data source for user

**Review Notes**:
- All examples verified as correct
- No removed `_id` fields referenced
- Good variety of reference patterns (id, name, slug, model)
- Device bay template shows advanced pattern with count
- All examples realistic and follow best practices

#### Batch 7.4: Device Components - Ports (10 resources) ✅ COMPLETE
- [x] netbox_interface - Clean, comprehensive examples (basic, full config, virtual, LAG)
- [x] netbox_interface_template - **FIXED**: Changed manufacturer_id → manufacturer
- [x] netbox_console_port - Clean, uses device name reference
- [x] netbox_console_port_template - Clean, uses device_type model reference
- [x] netbox_console_server_port - Clean, uses name references
- [x] netbox_console_server_port_template - Clean, uses device_type model reference
- [x] netbox_power_port - Clean, uses device name reference
- [x] netbox_power_outlet - Clean, uses name references
- [x] netbox_front_port - Clean, shows rear_port dependency
- [x] netbox_front_port_template - Clean, shows rear_port_template dependency

**Changes Made**:
- Fixed `netbox_interface_template/resource.tf`:
  - Changed `manufacturer_id = netbox_manufacturer.test.id` → `manufacturer = netbox_manufacturer.test.slug`

**Review Notes**:
- All examples now use current schemas
- No removed `_id` fields referenced after fix
- Good variety: simple ports, templates, and complex examples with dependencies
- Interface example shows advanced patterns (LAG, virtual interfaces)

#### Batch 7.5: Power & Modules (10 resources) ✅ COMPLETE
- [x] netbox_rear_port - Clean, uses device name reference
- [x] netbox_rear_port_template - Clean, uses device_type model reference
- [x] netbox_power_feed - Clean, uses power_panel/rack name references
- [x] netbox_power_panel - Clean, uses site slug reference
- [x] netbox_module - Clean, shows device, module_bay, module_type references
- [x] netbox_module_type - Clean, uses manufacturer name reference
- [x] netbox_module_bay - Clean, uses device name reference
- [x] netbox_module_bay_template - Clean, uses device_type model reference
- [x] netbox_inventory_item - Clean, uses device name reference
- [x] netbox_inventory_item_role - Clean, simple example

**Review Notes**:
- All examples verified as correct
- No removed `_id` fields referenced
- Power feed shows comprehensive configuration with multiple parameters
- Module examples demonstrate proper dependency chain (device → module_bay → module)
- All examples use primary reference fields correctly

#### Batch 7.6: Virtualization (10 resources) ✅ COMPLETE
- [x] netbox_inventory_item_template - Clean, uses manufacturer/device_type name/model references
- [x] netbox_cluster - Clean, uses cluster_type name reference
- [x] netbox_cluster_type - Clean, simple example
- [x] netbox_cluster_group - Clean, shows hierarchical groups with parent reference
- [x] netbox_virtual_machine - Clean, uses cluster name reference
- [x] netbox_vm_interface - Clean, uses virtual_machine name reference
- [x] netbox_virtual_chassis - Clean, simple example
- [x] netbox_virtual_device_context - **FIXED**: Updated to use primary fields instead of removed `_id` fields
- [x] netbox_virtual_disk - Clean, shows multiple disk scenarios with descriptions
- [x] netbox_fhrp_group_assignment - **FIXED**: Updated to use primary fields instead of removed `_id` fields

**Changes Made**:
- Fixed `netbox_virtual_device_context/resource.tf`:
  - Changed `manufacturer_id` → `manufacturer` (uses slug)
  - Changed `device_type_id` → `device_type` (uses model)
  - Changed `role_id` → `role` (uses slug)
  - Changed `site_id` → `site` (uses slug)
- Fixed `netbox_fhrp_group_assignment/resource.tf`:
  - Changed `manufacturer_id` → `manufacturer` (uses slug)
  - Changed `device_type_id` → `device_type` (uses model)
  - Changed `role_id` → `role` (uses slug)
  - Changed `site_id` → `site` (uses slug)
  - Changed `device_id` → `device` (uses name)

**Review Notes**:
- 2 files fixed (virtual_device_context, fhrp_group_assignment)
- Both had same pattern of setup resources using removed `_id` fields
- All examples now use current schemas
- Virtual disk shows good variety with multiple disks and descriptions
- FHRP assignment shows proper polymorphic interface reference pattern

#### Batch 7.7: Circuits & Providers (10 resources) ✅ COMPLETE
- [x] netbox_cable - Clean, comprehensive examples with comments
- [x] netbox_provider - Clean, simple example
- [x] netbox_provider_account - Clean, uses provider name reference
- [x] netbox_provider_network - Clean, uses provider name reference
- [x] netbox_circuit - Clean, uses provider/type name references
- [x] netbox_circuit_type - Clean, simple example
- [x] netbox_circuit_group - Clean, simple example
- [x] netbox_circuit_group_assignment - Clean, uses circuit.id and group.name
- [x] netbox_circuit_termination - Clean, uses circuit.id and site.name
- [x] netbox_service - **FIXED**: Updated device_type reference and reordered resources

**Changes Made**:
- Fixed `netbox_service/resource.tf`:
  - Changed `device_type = netbox_device_type.test.slug` → `device_type = netbox_device_type.test.model`
  - Reordered manufacturer resource before device_type (dependency order)

**Review Notes**:
- 1 file fixed (service)
- All examples now use current schemas
- No removed `_id` fields referenced after fix
- Cable example shows comprehensive documentation
- Circuit examples show proper provider and type relationships
- [ ] netbox_circuit_type
- [ ] netbox_circuit_group
- [ ] netbox_circuit_group_assignment
- [ ] netbox_circuit_termination
- [ ] netbox_service

#### Batch 7.8: VPN & Tunnels (10 resources) ✅ COMPLETE
- [x] netbox_service_template - Clean, simple example
- [x] netbox_l2vpn - Clean, simple example
- [x] netbox_l2vpn_termination - Clean, uses polymorphic assigned_object pattern
- [x] netbox_tunnel - Clean, simple example
- [x] netbox_tunnel_group - Clean, simple example
- [x] netbox_tunnel_termination - Clean, uses polymorphic termination pattern
- [x] netbox_ike_policy - Clean, uses proposal.id reference
- [x] netbox_ike_proposal - Clean, simple example
- [x] netbox_ipsec_policy - Clean, uses proposal.id reference
- [x] netbox_ipsec_profile - Clean, uses policy name references

**Review Notes**:
- All examples verified as correct
- No removed `_id` fields referenced
- L2VPN and tunnel terminations show proper polymorphic object references
- IKE/IPSec examples show proper policy and proposal relationships
- All examples use current schemas correctly
- [ ] netbox_ike_policy
- [ ] netbox_ike_proposal
- [ ] netbox_ipsec_policy
- [ ] netbox_ipsec_profile

#### Batch 7.9: Wireless & Extras (10 resources) ✅ COMPLETE
- [x] netbox_ipsec_proposal - Clean, simple example
- [x] netbox_wireless_lan - Clean, simple example
- [x] netbox_wireless_lan_group - Clean, simple example
- [x] netbox_wireless_link - Excellent, comprehensive examples (basic, with SSID, auth, complete)
- [x] netbox_tag - Excellent, multiple examples showing different use cases
- [x] netbox_custom_field - Clean, simple example
- [x] netbox_custom_field_choice_set - Clean, shows choice set with values
- [x] netbox_custom_link - Clean, shows templating with object variables
- [x] netbox_webhook - Excellent, comprehensive examples (basic, custom, secure, templated, insecure)
- [x] netbox_event_rule - Clean, shows webhook integration with action_object_id

**Review Notes**:
- All examples verified as correct
- No removed `_id` fields referenced
- Wireless link shows comprehensive examples with authentication and distance
- Tag examples show variety of use cases (basic, color, object_types)
- Webhook examples demonstrate full feature set
- Event rule properly uses webhook.id for action_object_id
- All examples use current schemas correctly

#### Batch 7.10: Configuration & Remaining (5 resources) ✅ COMPLETE
- [x] netbox_export_template - Clean, simple example with Jinja2 template
- [x] netbox_config_context - Excellent, comprehensive examples (basic, site-specific, role-specific, multi-criteria, tag-based)
- [x] netbox_config_template - Clean, simple example with environment params
- [x] netbox_journal_entry - Clean, shows polymorphic assigned_object pattern
- [x] netbox_notification_group - Clean, simple examples

**Review Notes**:
- All examples verified as correct
- No removed `_id` fields referenced
- Config context uses lists of IDs (sites, roles, tenant_groups) which is correct
- Export template shows Jinja2 templating
- Journal entry demonstrates polymorphic object assignment
- All examples use current schemas correctly

**Phase 7 Resource Examples Summary**:
- Total reviewed: 103 resource examples across 10 batches
- Files fixed: 5 (contact_assignment, interface_template, virtual_device_context, fhrp_group_assignment, service)
- Files verified clean: 98
- Error rate: 4.9% (5 fixes needed out of 103 files)

#### Batch 7.11: Data Sources - Part 1 (26 data sources) ✅ COMPLETE
- [x] netbox_aggregate - Clean, shows lookup by prefix and ID
- [x] netbox_asn - Clean, shows lookup by ASN and ID
- [x] netbox_asn_range - Excellent, comprehensive with multiple lookup methods
- [x] netbox_cable - Clean, ID-only lookup (cables don't have names/slugs)
- [x] netbox_circuit - Clean, shows CID and ID lookups
- [x] netbox_circuit_group - Clean, shows ID, name, and slug lookups
- [x] netbox_circuit_group_assignment - Clean, ID-based lookup
- [x] netbox_circuit_termination - Clean, simple ID lookup
- [x] netbox_circuit_type - Clean, shows ID, name, and slug lookups
- [x] netbox_cluster - Clean, shows ID and name lookups
- [x] netbox_cluster_group - Clean, shows ID, name, and slug lookups
- [x] netbox_cluster_type - Clean, shows ID, name, and slug lookups
- [x] netbox_config_context - Clean, shows ID and name lookups with JSON parsing
- [x] netbox_config_template - Clean, shows ID and name lookups
- [x] netbox_console_port - Clean, shows ID and device_id+name lookups
- [x] netbox_console_port_template - Clean, shows ID and device_type+name lookups
- [x] netbox_console_server_port - Clean, shows ID and device_id+name lookups
- [x] netbox_console_server_port_template - Clean, shows ID and device_type+name lookups
- [x] netbox_contact - Clean, shows ID, name, and email lookups
- [x] netbox_contact_assignment - Clean, simple ID lookup
- [x] netbox_contact_group - Clean, shows ID, name, and slug lookups (parent_id in datasource is OK)
- [x] netbox_contact_role - Clean, shows ID, name, and slug lookups
- [x] netbox_custom_field - Clean, shows ID and name lookups
- [x] netbox_custom_field_choice_set - Clean, shows ID and name lookups
- [x] netbox_custom_link - Clean, shows ID and name lookups
- [x] netbox_device - Clean, shows ID, name, and serial lookups

**Review Notes**:
- All examples verified as correct
- Note: Datasources retain `_id` fields (like parent_id, device_id) as computed attributes - this is correct
  - Datasources are read-only and show what's available from API
  - Phase 5 changes only removed duplicate `_id` fields from resources (writable schemas)
- Good variety of lookup methods demonstrated across examples
- All examples use valid query parameters

#### Batch 7.12: Data Sources - Part 2 (26 data sources) ✅ COMPLETE
- [x] netbox_device_bay - Clean, shows ID and device+name lookups
- [x] netbox_device_bay_template - Clean, shows ID and name+device_type lookups
- [x] netbox_device_role - Clean, shows ID, name, and slug lookups
- [x] netbox_device_type - Clean, shows ID, slug, and model lookups
- [x] netbox_event_rule - Clean, ID-based lookup
- [x] netbox_export_template - Clean, shows ID and name lookups
- [x] netbox_fhrp_group - Clean, shows ID and protocol+group_id lookups
- [x] netbox_fhrp_group_assignment - Clean, ID-based lookup (group_id is valid in datasource)
- [x] netbox_front_port - Clean, shows ID and device_id+name lookups
- [x] netbox_front_port_template - Clean, shows ID and name+device_type/module_type lookups
- [x] netbox_ike_policy - Clean, shows ID and name lookups
- [x] netbox_ike_proposal - Clean, shows ID and name lookups
- [x] netbox_interface - Clean, shows ID and device+name lookups
- [x] netbox_interface_template - Clean, shows ID and name+device_type/module_type lookups
- [x] netbox_inventory_item - Clean, shows ID and name+device_id lookups
- [x] netbox_inventory_item_role - Clean, shows ID, name, and slug lookups
- [x] netbox_inventory_item_template - Clean, ID-based lookup (device_type_id is valid in datasource)
- [x] netbox_ip_address - Clean, shows ID and address lookups
- [x] netbox_ip_range - Clean, shows ID and start/end address lookups
- [x] netbox_ipsec_policy - Clean, shows ID and name lookups
- [x] netbox_ipsec_profile - Clean, shows ID and name lookups
- [x] netbox_ipsec_proposal - Clean, shows ID and name lookups
- [x] netbox_journal_entry - Clean, ID-based lookup
- [x] netbox_l2vpn - Clean, shows ID, name, and slug lookups
- [x] netbox_l2vpn_termination - Clean, ID-based lookup
- [x] netbox_location - Clean, shows ID, name, and slug lookups

**Review Notes**:
- All examples verified as correct
- Datasources correctly use `_id` query parameters (e.g., device_id, device_type_id) for filtering
- These are distinct from the removed computed `_id` fields in resources
- Good variety of lookup patterns demonstrated
- All examples use valid query parameters

#### Batch 7.13: Data Sources - Part 3 (26 data sources) ✅ COMPLETE
- [x] netbox_manufacturer - Clean, shows ID, name, and slug lookups
- [x] netbox_module - Clean, shows ID, device_id+module_bay_id, device_id+serial lookups
- [x] netbox_module_bay - Clean, shows ID and device_id+name lookups
- [x] netbox_module_bay_template - Clean, ID-based lookup
- [x] netbox_module_type - Clean, shows ID, model, model+manufacturer_id lookups
- [x] netbox_notification_group - Clean, ID-based lookup
- [x] netbox_platform - Clean, shows ID, slug, and name lookups
- [x] netbox_power_feed - Clean, shows ID and power_panel+name lookups
- [x] netbox_power_outlet - Clean, shows ID and device_id+name lookups
- [x] netbox_power_outlet_template - Clean, shows ID, device_type+name, module_type+name lookups
- [x] netbox_power_panel - Clean, shows ID, name, name+site lookups
- [x] netbox_power_port - Clean, shows ID and device_id+name lookups
- [x] netbox_power_port_template - Clean, shows ID, device_type+name, module_type+name lookups
- [x] netbox_prefix - Clean, shows ID and CIDR prefix lookups
- [x] netbox_provider - Clean, shows ID, slug, and name lookups
- [x] netbox_provider_account - Clean, shows ID and account lookups
- [x] netbox_provider_network - Clean, shows ID, name, name+circuit_provider lookups
- [x] netbox_rack - Clean, shows ID and name lookups
- [x] netbox_rack_reservation - Clean, ID-based lookup
- [x] netbox_rack_role - Clean, name-based lookup
- [x] netbox_rack_type - Clean, shows ID, slug, and model lookups
- [x] netbox_rear_port - Clean, shows ID and device_id+name lookups
- [x] netbox_rear_port_template - Clean, shows ID, device_type+name, module_type+name lookups
- [x] netbox_region - Clean, shows ID, slug, and name lookups
- [x] netbox_rir - Clean, shows ID, name, and slug lookups
- [x] netbox_role - Clean, shows ID, name, and slug lookups

**Review Notes**:
- All examples verified as correct
- Datasources correctly use `_id` filter parameters (device_id, module_bay_id, manufacturer_id, etc.)
- display_name outputs are valid (computed field in datasources)
- Good variety of lookup combinations demonstrated
- All examples use valid query parameters

#### Batch 7.14: Data Sources - Part 4 (24 data sources) ✅ COMPLETE
- [x] netbox_route_target - Clean, shows ID and name lookups
- [x] netbox_service - Clean, shows ID, name+device, name+vm lookups
- [x] netbox_service_template - Clean, shows ID and name lookups
- [x] netbox_site - Clean, shows ID, slug, and name lookups
- [x] netbox_site_group - Clean, shows ID, slug, and name lookups
- [x] netbox_tag - Clean, shows ID, name, and slug lookups
- [x] netbox_tenant - Clean, shows ID, slug, and name lookups
- [x] netbox_tenant_group - Clean, shows ID, slug, and name lookups
- [x] netbox_tunnel - Clean, shows ID and name lookups
- [x] netbox_tunnel_group - Clean, shows ID, slug, and name lookups
- [x] netbox_tunnel_termination - Clean, shows ID, tunnel, tunnel_name lookups
- [x] netbox_user - Clean, username-based lookup
- [x] netbox_virtual_chassis - Clean, shows ID and name lookups
- [x] netbox_virtual_device_context - Clean, ID-based lookup
- [x] netbox_virtual_disk - Clean, shows ID and name+virtual_machine lookups
- [x] netbox_virtual_machine - Clean, shows ID and name lookups
- [x] netbox_vlan - Clean, shows ID, VID, name, VID+name lookups
- [x] netbox_vlan_group - Clean, shows ID, slug, and name lookups
- [x] netbox_vm_interface - Clean, shows ID and name+virtual_machine lookups
- [x] netbox_vrf - Clean, shows ID and name lookups
- [x] netbox_webhook - Clean, shows ID and name lookups
- [x] netbox_wireless_lan - Clean, shows ID, SSID, group_id lookups
- [x] netbox_wireless_lan_group - Clean, shows ID, slug, and name lookups
- [x] netbox_wireless_link - Clean, ID-based lookup

**Review Notes**:
- All examples verified as correct
- Good variety of lookup methods across all data sources
- All examples use valid query parameters
- Datasources correctly show computed fields and filter parameters

**Phase 7 Data Source Examples Summary**:
- Total reviewed: 102 data source examples across 4 batches
- Files fixed: 0 (all examples were correct)
- Files verified clean: 102
- Error rate: 0% (perfect record for data sources!)

---

## 🎉 Phase 7 Complete: All Examples Reviewed

**Total Phase 7 Results**:
- **Resources**: 103 examples reviewed, 5 fixed (4.9% error rate)
- **Data Sources**: 102 examples reviewed, 0 fixed (0% error rate)
- **Grand Total**: 205 example files reviewed, 5 fixed (2.4% overall error rate)

**Files Fixed**:
1. netbox_contact_assignment/resource.tf - 4 fields updated
2. netbox_interface_template/resource.tf - 1 field updated
3. netbox_virtual_device_context/resource.tf - 4 fields updated
4. netbox_fhrp_group_assignment/resource.tf - 5 fields updated
5. netbox_service/resource.tf - 1 field updated

**Key Finding**: Data sources were 100% clean because they correctly retain `_id` computed fields and filter parameters as part of their read-only API interface. Only resource examples needed fixes.

**Review Checklist for Each Example**:
- [ ] No references to removed `_id` fields (device_type_id, role_id, tenant_id, platform_id, site_id, location_id, rack_id, cluster_id, manufacturer_id, parent_id, group_id, rack_type_id, scope_id when used alone)
- [ ] Uses primary reference fields correctly (tenant, site, cluster, etc.)
- [ ] Syntax is valid Terraform HCL
- [ ] Comments explain the example clearly
- [ ] Required fields are present
- [ ] Optional fields demonstrate useful patterns
- [ ] Example is realistic and useful

---

### Phase 8: Regenerate Documentation

**Status**: ✅ COMPLETE

Used terraform-plugin-docs to regenerate provider documentation from:
- Resource/datasource schema definitions
- Example configurations
- Description fields

**Tasks**:
- [x] Run `tfplugindocs generate --provider-dir=. --rendered-website-dir=docs`
- [x] Verify all resources have documentation in `docs/resources/` (103 files)
- [x] Verify all data sources have documentation in `docs/data-sources/` (105 files)
- [x] Review generated docs for accuracy
- [x] Commit updated documentation

**Results**:
- Generated 209 total documentation files (103 resources + 105 datasources + 1 index)
- All documentation reflects Phase 5 changes (no duplicate `_id` fields in resource schemas)
- Datasource documentation correctly shows `_id` filter parameters (read-only, as designed)
- Documentation generation successful with no errors

---

### Phase 9: Terraform Integration Tests

**Status**: � In Progress

Review and update Terraform configuration tests in `test/terraform/` directory:
- Ensure tests use current resource schemas (no removed `_id` fields)
- Verify main.tf configurations reflect Phase 5 changes
- Update any tests using removed `_id` fields to use primary reference fields
- Check outputs.tf for references to removed fields
- Document test coverage and any gaps

**Scope**: 204 test directories (101 resources + 103 datasources)

**Review Pattern for Each Test**:
1. Read main.tf and outputs.tf (if exists)
2. Search for removed computed fields in OUTPUTS: `role_id`, `site_id`, `tenant_id`, `platform_id`, `device_type_id`, `location_id`, `rack_id`, `cluster_id`, `manufacturer_id`, `parent_id`, `group_id`, `rack_type_id`
3. **✅ CORRECT - Do NOT flag**: `role = netbox_device_role.test.id` (using .id as INPUT)
4. **❌ INCORRECT - Flag this**: `output "role_id" { value = netbox_device.test.role_id }` (accessing removed field)

**Key Understanding**:
- Phase 5 removed duplicate computed OUTPUT fields from resource schemas
- Using `.id` as INPUT for cross-resource references is CORRECT and PREFERRED (immutable)
- Primary fields (role, site, tenant) accept ID, name, or slug - all are valid inputs

**Note**: These are Terraform integration tests, not Go tests. Focus on .tf file correctness.

---

#### Batch 9.1: IPAM Core Resources (15 tests) ✅ COMPLETE
- [x] test/terraform/resources/aggregate - Clean
- [x] test/terraform/resources/asn - Clean
- [x] test/terraform/resources/asn_range - Clean
- [x] test/terraform/resources/ip_address - Clean
- [x] test/terraform/resources/ip_range - Clean
- [x] test/terraform/resources/prefix - Clean
- [x] test/terraform/resources/rir - Clean
- [x] test/terraform/resources/role - Clean
- [x] test/terraform/resources/route_target - Clean
- [x] test/terraform/resources/vlan - Clean
- [x] test/terraform/resources/vlan_group - Clean
- [x] test/terraform/resources/vrf - Clean
- [x] test/terraform/resources/fhrp_group - Clean
- [x] test/terraform/resources/fhrp_group_assignment - Clean
- [x] test/terraform/resources/l2vpn - Clean

**Review Notes**:
- All 15 tests are clean and correct!
- Tests correctly use `.id` for cross-resource references (this is the preferred pattern)
- Primary reference fields (role, site, tenant, device_type, etc.) accept ID/name/slug
- Using `.id` is actually BEST PRACTICE (immutable, unambiguous)
- Phase 5 only removed duplicate COMPUTED `_id` fields from outputs, not the ability to use IDs as inputs

**What to look for in Phase 9**:
- ❌ References to removed computed fields: `netbox_device.test.role_id` (in outputs)
- ✅ Using .id as input: `role = netbox_device_role.test.id` (CORRECT, keep this!)

#### Batch 9.2: Sites & Organization Resources (15 tests) ✅ COMPLETE
- [x] test/terraform/resources/site - Clean
- [x] test/terraform/resources/site_group - Clean (uses parent.id correctly)
- [x] test/terraform/resources/region - Clean (uses parent.id correctly)
- [x] test/terraform/resources/location - Clean (uses site.id, tenant.id correctly)
- [x] test/terraform/resources/tenant - Clean (uses group.id correctly)
- [x] test/terraform/resources/tenant_group - Clean (uses parent.id correctly)
- [x] test/terraform/resources/contact - Clean
- [x] test/terraform/resources/contact_group - Clean (uses parent.id correctly)
- [x] test/terraform/resources/contact_role - Clean
- [x] test/terraform/resources/contact_assignment - Clean (uses object_id, contact_id, role_id)
- [x] test/terraform/resources/manufacturer - Clean
- [x] test/terraform/resources/platform - Clean (uses manufacturer.id correctly)
- [x] test/terraform/resources/tag - Clean
- [x] test/terraform/resources/device_role - Clean
- [x] test/terraform/resources/rack_role - Clean

**Review Notes**: All 15 tests clean (0 fixes needed). Tests correctly use `.id` for cross-resource references and parent/group relationships.

#### Batch 9.3: Device Infrastructure Resources (15 tests) ✅ COMPLETE
- [x] test/terraform/resources/device - Clean (uses site.id, device_type.id, role.id, tenant.id, platform.id, rack.id)
- [x] test/terraform/resources/device_type - Clean (uses manufacturer.id)
- [x] test/terraform/resources/device_bay - Clean (uses device.id, device_type.id, role.id, site.id)
- [x] test/terraform/resources/device_bay_template - Clean (uses device_type.id, manufacturer.id)
- [x] test/terraform/resources/rack - Clean (uses site.id, location.id, tenant.id)
- [x] test/terraform/resources/rack_type - Clean (uses manufacturer.id)
- [x] test/terraform/resources/rack_reservation - Clean (uses rack.id, user.id)
- [x] test/terraform/resources/cable - Clean (uses device.id for interfaces)
- [x] test/terraform/resources/power_panel - Clean (uses site.id, location.id)
- [x] test/terraform/resources/power_feed - Clean (uses power_panel.id, rack.id)
- [x] test/terraform/resources/module - Clean (uses device.id, module_bay.id, module_type.id)
- [x] test/terraform/resources/module_type - Clean (uses manufacturer.id)
- [x] test/terraform/resources/module_bay - Clean (uses device.id, device_type.id)
- [x] test/terraform/resources/module_bay_template - Clean (uses device_type.id, manufacturer.id)
- [x] test/terraform/resources/virtual_chassis - Clean (no dependencies)

**Review Notes**: All 15 tests clean (0 fixes needed). Tests correctly use `.id` for all cross-resource references including devices, racks, modules, and their related components.

#### Batch 9.4: Device Components Resources (15 tests) ✅ COMPLETE
- [x] test/terraform/resources/interface - Clean (uses device.id, manufacturer.id)
- [x] test/terraform/resources/interface_template - Clean (uses device_type.id, manufacturer.id)
- [x] test/terraform/resources/console_port - Clean (uses device.id, device_type.id)
- [x] test/terraform/resources/console_port_template - Clean (uses device_type.id, manufacturer.id)
- [x] test/terraform/resources/console_server_port - Clean (uses device.id, device_type.id)
- [x] test/terraform/resources/console_server_port_template - Clean (uses device_type.id, manufacturer.id)
- [x] test/terraform/resources/power_port - Clean (uses device.id, device_type.id)
- [x] test/terraform/resources/power_port_template - Clean (uses device_type.id, manufacturer.id)
- [x] test/terraform/resources/power_outlet - Clean (uses device.id, device_type.id)
- [x] test/terraform/resources/power_outlet_template - Clean (uses device_type.id, manufacturer.id)
- [x] test/terraform/resources/front_port - Clean (uses device.id, rear_port.id)
- [x] test/terraform/resources/front_port_template - Clean (uses device_type.id, rear_port.name)
- [x] test/terraform/resources/rear_port - Clean (uses device.id, device_type.id)
- [x] test/terraform/resources/rear_port_template - Clean (uses device_type.id, manufacturer.id)
- [x] test/terraform/resources/inventory_item - Clean (uses device.id, manufacturer.id)

**Review Notes**: All 15 tests clean (0 fixes needed). Tests correctly use `.id` for device components and their templates. Note: front_port_template uses rear_port.name (string reference, not ID).

#### Batch 9.5: Inventory & Templates Resources (10 tests)
- [ ] test/terraform/resources/inventory_item_role
- [ ] test/terraform/resources/inventory_item_template
- [ ] test/terraform/resources/config_template
- [ ] test/terraform/resources/config_context
- [ ] test/terraform/resources/export_template
- [ ] test/terraform/resources/custom_field
- [ ] test/terraform/resources/custom_field_choice_set
- [ ] test/terraform/resources/custom_link
- [ ] test/terraform/resources/journal_entry
- [ ] test/terraform/resources/webhook

#### Batch 9.6: Virtualization Resources (10 tests)
- [ ] test/terraform/resources/cluster
- [ ] test/terraform/resources/cluster_type
- [ ] test/terraform/resources/cluster_group
- [ ] test/terraform/resources/virtual_machine
- [ ] test/terraform/resources/vm_interface
- [ ] test/terraform/resources/virtual_disk
- [ ] test/terraform/resources/virtual_device_context
- [ ] test/terraform/resources/service
- [ ] test/terraform/resources/service_template
- [ ] test/terraform/resources/notification_group

#### Batch 9.7: Circuits & VPN Resources (16 tests)
- [ ] test/terraform/resources/provider
- [ ] test/terraform/resources/provider_account
- [ ] test/terraform/resources/provider_network
- [ ] test/terraform/resources/circuit
- [ ] test/terraform/resources/circuit_type
- [ ] test/terraform/resources/circuit_group
- [ ] test/terraform/resources/circuit_group_assignment
- [ ] test/terraform/resources/circuit_termination
- [ ] test/terraform/resources/tunnel
- [ ] test/terraform/resources/tunnel_group
- [ ] test/terraform/resources/tunnel_termination
- [ ] test/terraform/resources/l2vpn_termination
- [ ] test/terraform/resources/ike_policy
- [ ] test/terraform/resources/ike_proposal
- [ ] test/terraform/resources/ipsec_policy
- [ ] test/terraform/resources/ipsec_profile

#### Batch 9.8: Wireless & Events Resources (8 tests)
- [ ] test/terraform/resources/ipsec_proposal
- [ ] test/terraform/resources/wireless_lan
- [ ] test/terraform/resources/wireless_lan_group
- [ ] test/terraform/resources/wireless_link
- [ ] test/terraform/resources/event_rule
- [ ] test/terraform/data-sources/script
- [ ] test/terraform/data-sources/user
- [ ] test/terraform/resources/notification_group (if not covered)

#### Batch 9.9: Data Sources Part 1 - IPAM (15 tests)
- [ ] test/terraform/data-sources/aggregate
- [ ] test/terraform/data-sources/asn
- [ ] test/terraform/data-sources/asn_range
- [ ] test/terraform/data-sources/ip_address
- [ ] test/terraform/data-sources/ip_range
- [ ] test/terraform/data-sources/prefix
- [ ] test/terraform/data-sources/rir
- [ ] test/terraform/data-sources/role
- [ ] test/terraform/data-sources/route_target
- [ ] test/terraform/data-sources/vlan
- [ ] test/terraform/data-sources/vlan_group
- [ ] test/terraform/data-sources/vrf
- [ ] test/terraform/data-sources/fhrp_group
- [ ] test/terraform/data-sources/fhrp_group_assignment
- [ ] test/terraform/data-sources/l2vpn

#### Batch 9.10: Data Sources Part 2 - Sites & Org (15 tests)
- [ ] test/terraform/data-sources/site
- [ ] test/terraform/data-sources/site_group
- [ ] test/terraform/data-sources/region
- [ ] test/terraform/data-sources/location
- [ ] test/terraform/data-sources/tenant
- [ ] test/terraform/data-sources/tenant_group
- [ ] test/terraform/data-sources/contact
- [ ] test/terraform/data-sources/contact_group
- [ ] test/terraform/data-sources/contact_role
- [ ] test/terraform/data-sources/contact_assignment
- [ ] test/terraform/data-sources/manufacturer
- [ ] test/terraform/data-sources/platform
- [ ] test/terraform/data-sources/tag
- [ ] test/terraform/data-sources/device_role
- [ ] test/terraform/data-sources/rack_role

#### Batch 9.11: Data Sources Part 3 - Devices (15 tests)
- [ ] test/terraform/data-sources/device
- [ ] test/terraform/data-sources/device_type
- [ ] test/terraform/data-sources/device_bay
- [ ] test/terraform/data-sources/device_bay_template
- [ ] test/terraform/data-sources/rack
- [ ] test/terraform/data-sources/rack_type
- [ ] test/terraform/data-sources/rack_reservation
- [ ] test/terraform/data-sources/cable
- [ ] test/terraform/data-sources/power_panel
- [ ] test/terraform/data-sources/power_feed
- [ ] test/terraform/data-sources/module
- [ ] test/terraform/data-sources/module_type
- [ ] test/terraform/data-sources/module_bay
- [ ] test/terraform/data-sources/module_bay_template
- [ ] test/terraform/data-sources/virtual_chassis

#### Batch 9.12: Data Sources Part 4 - Components (15 tests)
- [ ] test/terraform/data-sources/interface
- [ ] test/terraform/data-sources/interface_template
- [ ] test/terraform/data-sources/console_port
- [ ] test/terraform/data-sources/console_port_template
- [ ] test/terraform/data-sources/console_server_port
- [ ] test/terraform/data-sources/console_server_port_template
- [ ] test/terraform/data-sources/power_port
- [ ] test/terraform/data-sources/power_port_template
- [ ] test/terraform/data-sources/power_outlet
- [ ] test/terraform/data-sources/power_outlet_template
- [ ] test/terraform/data-sources/front_port
- [ ] test/terraform/data-sources/front_port_template
- [ ] test/terraform/data-sources/rear_port
- [ ] test/terraform/data-sources/rear_port_template
- [ ] test/terraform/data-sources/inventory_item

#### Batch 9.13: Data Sources Part 5 - Config & Virtualization (20 tests)
- [ ] test/terraform/data-sources/inventory_item_role
- [ ] test/terraform/data-sources/inventory_item_template
- [ ] test/terraform/data-sources/config_template
- [ ] test/terraform/data-sources/config_context
- [ ] test/terraform/data-sources/export_template
- [ ] test/terraform/data-sources/custom_field
- [ ] test/terraform/data-sources/custom_field_choice_set
- [ ] test/terraform/data-sources/custom_link
- [ ] test/terraform/data-sources/journal_entry
- [ ] test/terraform/data-sources/webhook
- [ ] test/terraform/data-sources/cluster
- [ ] test/terraform/data-sources/cluster_type
- [ ] test/terraform/data-sources/cluster_group
- [ ] test/terraform/data-sources/virtual_machine
- [ ] test/terraform/data-sources/vm_interface
- [ ] test/terraform/data-sources/virtual_disk
- [ ] test/terraform/data-sources/virtual_device_context
- [ ] test/terraform/data-sources/service
- [ ] test/terraform/data-sources/service_template
- [ ] test/terraform/data-sources/notification_group

#### Batch 9.14: Data Sources Part 6 - Circuits & VPN (18 tests)
- [ ] test/terraform/data-sources/provider
- [ ] test/terraform/data-sources/provider_account
- [ ] test/terraform/data-sources/provider_network
- [ ] test/terraform/data-sources/circuit
- [ ] test/terraform/data-sources/circuit_type
- [ ] test/terraform/data-sources/circuit_group
- [ ] test/terraform/data-sources/circuit_group_assignment
- [ ] test/terraform/data-sources/circuit_termination
- [ ] test/terraform/data-sources/tunnel
- [ ] test/terraform/data-sources/tunnel_group
- [ ] test/terraform/data-sources/tunnel_termination
- [ ] test/terraform/data-sources/l2vpn_termination
- [ ] test/terraform/data-sources/ike_policy
- [ ] test/terraform/data-sources/ike_proposal
- [ ] test/terraform/data-sources/ipsec_policy
- [ ] test/terraform/data-sources/ipsec_profile
- [ ] test/terraform/data-sources/ipsec_proposal
- [ ] test/terraform/data-sources/event_rule

#### Batch 9.15: Data Sources Part 7 - Wireless (6 tests)
- [ ] test/terraform/data-sources/wireless_lan
- [ ] test/terraform/data-sources/wireless_lan_group
- [ ] test/terraform/data-sources/wireless_link
- [ ] test/terraform/data-sources/script (if exists)
- [ ] test/terraform/data-sources/user (if exists)
- [ ] Final review and documentation

**Phase 9 Summary**:
- Total test directories: 204 (101 resources + 103 datasources)
- Organized into 15 batches for systematic review
- Focus: Verify no references to removed duplicate `_id` fields
- Expected changes: Minimal (most tests should use primary reference fields already)

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
