3
## Current Status
- **Overall Progress**: 23/99 resources (23.2%) - Batch 2D Complete
- **Last Update**: December 28, 2025 - Batch 2D Complete (20/20 tests PASS)
- **Strategy**: Sub-batch implementation with commits after each sub-batch

---

## Batch 1: Core Infrastructure Resources (8 resources)
**Status**: ✅ COMPLETE

Priority: CRITICAL - Most frequently used resources

| Resource | File | Status | Test Coverage | Notes |
|----------|------|--------|---------------|-------|
| device | device_resource_test.go | ✅ DONE | CRUD + Import | Core infrastructure |
| interface | interface_resource_test.go | ✅ DONE | CRUD + Import | Critical for networking |
| ip_address | ip_address_resource_test.go | ✅ DONE | CRUD + Import | Core IPAM resource |
| site | site_resource_test.go | ✅ DONE | CRUD + Import | Top-level organizational |
| vlan | vlan_resource_test.go | ✅ DONE | CRUD + Import | Network configuration |
| circuit | circuit_resource_test.go | ✅ DONE | CRUD + Import | Core interconnect |
| cable | cable_resource_test.go | ✅ DONE | CRUD + Import | Physical layer |
| cluster | cluster_resource_test.go | ✅ DONE | CRUD + Import | Virtualization core |

**Recent Improvements**:
- Added 404 handling to Delete methods for site and ip_address resources
- Added full, update, and import tests to cable resource for complete CRUD coverage
- All external deletion tests now pass (including cleanup phase)

---

## Batch 2A: Inventory Resources - Quick Wins (3 resources)
**Status**: ✅ COMPLETE

Priority: HIGH - Already have complete CRUD coverage, only need external deletion tests

| Resource | File | Status | Test Coverage | Notes |
|----------|------|--------|---------------|-------|
| inventory_item | inventory_item_resource_test.go | ✅ DONE | CRUD + Import + Ext Del ✅ | All tests complete |
| inventory_item_template | inventory_item_template_resource_test.go | ✅ DONE | CRUD + Import + Ext Del ✅ | All tests complete |
| inventory_item_role | inventory_item_role_resource_test.go | ✅ DONE | CRUD + Import + Ext Del ✅ | All tests complete |

**Implementation Notes**:
- Added context import to all test files
- All external deletion tests follow the standard pattern from Batch 1
- Uses DcimAPI methods: DcimInventoryItemsList/Destroy, DcimInventoryItemTemplatesList/Destroy, DcimInventoryItemRolesList/Destroy

---

## Batch 2B: Console Port Resources (4 resources)
**Status**: ✅ COMPLETE

Priority: HIGH - Most commonly used port types

| Resource | File | Status | Test Coverage | Notes |
|----------|------|--------|---------------|-------|
| console_port | console_port_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests complete |
| console_port_template | console_port_template_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests complete |
| console_server_port | console_server_port_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests complete |
| console_server_port_template | console_server_port_template_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests complete |

**Implementation Notes**:
- Added context import to all test files
- Added update tests for description/label field changes
- Added external deletion tests following Batch 1 pattern
- Uses DcimAPI methods: DcimConsolePorts*/DcimConsoleServerPorts* List/Destroy and Templates variants

---

## Batch 2C: Power Infrastructure Resources (4 resources)
**Status**: ✅ COMPLETE

Priority: MEDIUM - Critical for power management

| Resource | File | Status | Test Coverage | Notes |
|----------|------|--------|---------------|-------|
| power_port | power_port_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests PASS |
| power_port_template | power_port_template_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | Schema bugs fixed, all tests PASS |
| power_outlet | power_outlet_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests PASS |
| power_outlet_template | power_outlet_template_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | Schema bugs fixed, all tests PASS |

**Implementation Notes**:
- Added context import to all test files
- Added update tests for description and related field changes
- Added external deletion tests following Batch 2B pattern
- Uses DcimAPI methods: DcimPowerPorts*/DcimPowerOutlets* List/Destroy and Templates variants
- **Test Results**: 20/20 tests PASS (100%)
- **Schema Fixes Applied**:
  * power_port_template: Added missing `maximum_draw` and `allocated_draw` schema attributes
  * power_outlet_template: Added missing `feed_leg` schema attribute and fixed `power_port` type (Int32 vs String mismatch)

---

## Batch 2D: Patch Panel Port Resources (4 resources)
**Status**: ✅ COMPLETE

Priority: MEDIUM - Structured cabling

| Resource | File | Status | Test Coverage | Notes |
|----------|------|--------|---------------|-------|
| front_port | front_port_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests PASS |
| front_port_template | front_port_template_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests PASS |
| rear_port | rear_port_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests PASS |
| rear_port_template | rear_port_template_resource_test.go | ✅ DONE | CRUD + Import + Update + Ext Del ✅ | All tests PASS |

**Implementation Notes**:
- Added context import to all test files
- Added update tests for description and label field changes
- Added external deletion tests following Batch 2C pattern
- Uses DcimAPI methods: DcimFrontPorts*/DcimRearPorts* List/Destroy and Templates variants
- **Test Results**: 20/20 tests PASS (100%)
- All resources support both device-based (front_port, rear_port) and device_type-based (templates) configurations

---

## Batch 3: Module & Device Bay Resources (6 resources)
**Status**: ⏳ NOT STARTED

Priority: MEDIUM - Device component management

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| device_bay | device_bay_resource_test.go | ⏳ TODO | Update + External Deletion |
| device_bay_template | device_bay_template_resource_test.go | ⏳ TODO | Update + External Deletion |
| module_bay | module_bay_resource_test.go | ⏳ TODO | Update + External Deletion |
| module_bay_template | module_bay_template_resource_test.go | ⏳ TODO | Update + External Deletion |
| module | module_resource_test.go | ⏳ TODO | Update + External Deletion |
| module_type | module_type_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 4: Interface & Network Resources (8 resources)
**Status**: ⏳ NOT STARTED

Priority: HIGH - Network infrastructure

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| interface_template | interface_template_resource_test.go | ⏳ TODO | Update + External Deletion |
| vm_interface | vm_interface_resource_test.go | ⏳ TODO | Update + External Deletion |
| fhrp_group | fhrp_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| fhrp_group_assignment | fhrp_group_assignment_resource_test.go | ⏳ TODO | Update + External Deletion |
| service | service_resource_test.go | ⏳ TODO | Update + External Deletion |
| service_template | service_template_resource_test.go | ⏳ TODO | Update + External Deletion |
| l2vpn | l2vpn_resource_test.go | ⏳ TODO | Update + External Deletion |
| l2vpn_termination | l2vpn_termination_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 5: IPAM Additional Resources (6 resources)
**Status**: ⏳ NOT STARTED

Priority: HIGH - IP address management extensions

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| prefix | prefix_resource_test.go | ⏳ TODO | Update + External Deletion |
| ip_range | ip_range_resource_test.go | ⏳ TODO | Update + External Deletion |
| aggregate | aggregate_resource_test.go | ⏳ TODO | Update + External Deletion |
| rir | rir_resource_test.go | ⏳ TODO | Update + External Deletion |
| asn | asn_resource_test.go | ⏳ TODO | Update + External Deletion |
| asn_range | asn_range_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 6: VPN & Tunnel Resources (9 resources)
**Status**: ⏳ NOT STARTED

Priority: MEDIUM - VPN infrastructure

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| tunnel | tunnel_resource_test.go | ⏳ TODO | Update + External Deletion |
| tunnel_group | tunnel_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| tunnel_termination | tunnel_termination_resource_test.go | ⏳ TODO | Update + External Deletion |
| ike_policy | ike_policy_resource_test.go | ⏳ TODO | Update + External Deletion |
| ike_proposal | ike_proposal_resource_test.go | ⏳ TODO | Update + External Deletion |
| ipsec_policy | ipsec_policy_resource_test.go | ⏳ TODO | Update + External Deletion |
| ipsec_profile | ipsec_profile_resource_test.go | ⏳ TODO | Update + External Deletion |
| ipsec_proposal | ipsec_proposal_resource_test.go | ⏳ TODO | Update + External Deletion |
| route_target | route_target_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 7: Circuit & Provider Resources (7 resources)
**Status**: ⏳ NOT STARTED

Priority: MEDIUM - Service provider management

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| provider | provider_resource_test.go | ⏳ TODO | Update + External Deletion |
| provider_account | provider_account_resource_test.go | ⏳ TODO | Update + External Deletion |
| provider_network | provider_network_resource_test.go | ⏳ TODO | Update + External Deletion |
| circuit_type | circuit_type_resource_test.go | ⏳ TODO | Update + External Deletion |
| circuit_group | circuit_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| circuit_group_assignment | circuit_group_assignment_resource_test.go | ⏳ TODO | Update + External Deletion |
| circuit_termination | circuit_termination_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 8: Rack & Power Resources (6 resources)
**Status**: ⏳ NOT STARTED

Priority: MEDIUM - Physical infrastructure

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| rack | rack_resource_test.go | ⏳ TODO | Update + External Deletion |
| rack_role | rack_role_resource_test.go | ⏳ TODO | Update + External Deletion |
| rack_type | rack_type_resource_test.go | ⏳ TODO | Update + External Deletion |
| rack_reservation | rack_reservation_resource_test.go | ⏳ TODO | Update + External Deletion |
| power_panel | power_panel_resource_test.go | ⏳ TODO | Update + External Deletion |
| power_feed | power_feed_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 9: Organizational Resources (11 resources)
**Status**: ⏳ NOT STARTED

Priority: HIGH - Organizational structure

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| location | location_resource_test.go | ⏳ TODO | Update + External Deletion |
| region | region_resource_test.go | ⏳ TODO | Update + External Deletion |
| site_group | site_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| tenant | tenant_resource_test.go | ⏳ TODO | Update + External Deletion |
| tenant_group | tenant_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| contact | contact_resource_test.go | ⏳ TODO | Update + External Deletion |
| contact_group | contact_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| contact_role | contact_role_resource_test.go | ⏳ TODO | Update + External Deletion |
| contact_assignment | contact_assignment_resource_test.go | ⏳ TODO | Update + External Deletion |
| tag | tag_resource_test.go | ⏳ TODO | Update + External Deletion |
| role | role_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 10: Device & Virtualization Metadata (10 resources)
**Status**: ⏳ NOT STARTED

Priority: MEDIUM - Device and VM type definitions

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| device_role | device_role_resource_test.go | ⏳ TODO | Update + External Deletion |
| device_type | device_type_resource_test.go | ⏳ TODO | Update + External Deletion |
| manufacturer | manufacturer_resource_test.go | ⏳ TODO | Update + External Deletion |
| platform | platform_resource_test.go | ⏳ TODO | Update + External Deletion |
| cluster_type | cluster_type_resource_test.go | ⏳ TODO | Update + External Deletion |
| cluster_group | cluster_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| virtual_machine | virtual_machine_resource_test.go | ⏳ TODO | Update + External Deletion |
| virtual_chassis | virtual_chassis_resource_test.go | ⏳ TODO | Update + External Deletion |
| virtual_device_context | virtual_device_context_resource_test.go | ⏳ TODO | Update + External Deletion |
| virtual_disk | virtual_disk_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 11: Wireless Resources (3 resources)
**Status**: ⏳ NOT STARTED

Priority: LOW - Wireless network management

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| wireless_lan | wireless_lan_resource_test.go | ⏳ TODO | Update + External Deletion |
| wireless_lan_group | wireless_lan_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| wireless_link | wireless_link_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 12: Configuration & Customization Resources (8 resources)
**Status**: ⏳ NOT STARTED

Priority: LOW - Configuration management and customization

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| config_context | config_context_resource_test.go | ⏳ TODO | Update + External Deletion |
| config_template | config_template_resource_test.go | ⏳ TODO | Update + External Deletion |
| custom_field | custom_field_resource_test.go | ⏳ TODO | Update + External Deletion |
| custom_field_choice_set | custom_field_choice_set_resource_test.go | ⏳ TODO | Update + External Deletion |
| custom_link | custom_link_resource_test.go | ⏳ TODO | Update + External Deletion |
| export_template | export_template_resource_test.go | ⏳ TODO | Update + External Deletion |
| webhook | webhook_resource_test.go | ⏳ TODO | Update + External Deletion |
| journal_entry | journal_entry_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Batch 13: VLAN & VRF Resources (3 resources)
**Status**: ⏳ NOT STARTED

Priority: HIGH - Network segmentation

| Resource | File | Status | Missing Tests |
|----------|------|--------|---------------|
| vlan_group | vlan_group_resource_test.go | ⏳ TODO | Update + External Deletion |
| vrf | vrf_resource_test.go | ⏳ TODO | Update + External Deletion |

---

## Test Pattern

All external deletion tests follow this structure:
```go
func TestAcc{ResourceName}Resource_externalDeletion(t *testing.T) {
	t.Parallel()
	tex] Batch 1 (8) - ✅ COMPLETE
- [ ] Batch 2A (3) - 🔄 IN PROGRESS
- [ ] Batch 2B (4) - Console Ports
- [ ] Batch 2C (4) - Power Infrastructure
- [ ] Batch 2D (4) - Patch Panel Ports
- [ ] Batch 3 (10)
- [ ] Batch 4 (10)
- [ ] Batch 5 (8)
- [ ] Batch 6 (8)
- [ ] Batch 7 (6)
## Test Pattern

All external deletion tests follow this structure:
```go
func TestAcc{ResourceName}Resource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)
	// Setup: Create resource with unique name
	// Step 1: Create and verify resource exists
	// Step 2: PreConfig hook deletes resource via API
	//         Reapply config (Terraform detects 404, recreates)
	//         Verify ID is set (resource was recreated)
}
```

---

## Implementation Notes

1. Each test verifies that when a resource is deleted externally via the NetBox API, Terraform properly:
   - Detects the 404 error during refresh
   - Marks the resource as missing from state
   - Recreates the resource on next apply

2. Imports required:
   - `"context"` - for API calls
   - `testutil.GetSharedClient()` - to get NetBox API client
   - Resource-specific API methods (e.g., `IpamAPI.IpamDevicesList()`)

3. Safe commit strategy:
   - Commit after each batch
   - Run `go build .` to verify no compilation errors
   - Verify test syntax is correct

---x] Batch 2B (4) - ✅ COMPLETE

## Batch Completion Tracking

- [x] Batch 1 (8) - ✅ COMPLETE - Core Infrastructure
- [x] Batch 2A (3) - ✅ COMPLETE - Inventory Resources
- [x] Batch 2B (4) - ✅ COMPLETE - Console Ports
- [x] Batch 2C (4) - ✅ COMPLETE - Power Infrastructure
- [x] Batch 2D (4) - ✅ COMPLETE - Patch Panel Ports
- [ ] Batch 3 (6) - Module & Device Bay Resources
- [ ] Batch 4 (8) - Interface & Network Resources
- [ ] Batch 5 (6) - IPAM Additional Resources
- [ ] Batch 6 (9) - VPN & Tunnel Resources
- [ ] Batch 7 (7) - Circuit & Provider Resources
- [ ] Batch 8 (6) - Rack & Power Resources
- [ ] Batch 9 (11) - Organizational Resources
- [ ] Batch 10 (10) - Device & Virtualization Metadata
- [ ] Batch 11 (3) - Wireless Resources
- [ ] Batch 12 (8) - Configuration & Customization
- [ ] Batch 13 (2) - VLAN & VRF Resources

**Target**: 99/99 resources with external deletion tests
**Current**: 23/99 (23.2%)
**Next Milestone**: 29/99 (29.3%) after Batch 3
