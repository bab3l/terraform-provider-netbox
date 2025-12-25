# ID Preservation Testing - Batch Tracking

## Overview
Testing all 101 resources and 104 datasources to ensure ID is preserved as the immutable identifier.

## Progress Summary
- **Total Items**: 205 (101 Resources + 104 Datasources)
- **Already Fixed**: 8 resources (circuit, device, device_type, device_role, rack_reservation, cluster, cluster_type, inventory_item, inventory_item_template, export_template, config_template)
- **Remaining**: ~193 items across 15 batches

---

## Batch 1: Core Reference Resources (IN PROGRESS)
**Status**: Mostly Complete - 7/9 Items Done
**Items**: 9
- [x] circuit_resource.go ✅ (already fixed - skip)
- [x] provider_resource.go ✅ **COMPLETE** - Unit test + Acceptance test added and passing
- [x] circuit_type_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] fhrp_group_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] fhrp_group_assignment_resource.go ✅ (already fixed - skip)
- [x] contact_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] contact_role_resource.go ✅ (already fixed - skip)
- [x] circuit_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing
- [x] provider_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing

**Tests Added**:
- `internal/utils/state_helpers_test.go`: `TestUpdateReferenceAttribute` (10 unit tests)
- `internal/resources_acceptance_tests/provider_resource_test.go`: `TestAccProviderResource_IDPreservation`
- `internal/resources_acceptance_tests/circuit_type_resource_test.go`: `TestAccCircuitTypeResource_IDPreservation`
- `internal/resources_acceptance_tests/fhrp_group_resource_test.go`: `TestAccFHRPGroupResource_IDPreservation`
- `internal/resources_acceptance_tests/contact_resource_test.go`: `TestAccContactResource_IDPreservation`
- `internal/datasources_acceptance_tests/circuit_data_source_test.go`: `TestAccCircuitDataSource_IDPreservation`
- `internal/datasources_acceptance_tests/provider_data_source_test.go`: `TestAccProviderDataSource_IDPreservation`

**Test Results**: All tests PASSING ✅
- Provider resource test: 42.37s ✅
- CircuitType resource test: 41.47s ✅
- FHRPGroup resource test: 42.22s ✅
- Contact resource test: 42.25s ✅
- Circuit datasource test: 42.06s ✅
- Provider datasource test: 41.89s ✅

---

## Batch 2: Device/Hardware Resources
**Status**: Not Started
**Items**: 9
- [ ] device_resource.go ✅ (already fixed - skip)
- [ ] device_type_resource.go ✅ (already fixed - skip)
- [ ] device_role_resource.go ✅ (already fixed - skip)
- [ ] device_bay_resource.go
- [ ] device_bay_template_resource.go ✅ (already fixed - skip)
- [ ] module_resource.go
- [ ] module_type_resource.go
- [ ] device_data_source.go
- [ ] device_type_data_source.go

---

## Batch 3: Tenancy/Organization Resources
**Status**: Not Started
**Items**: 8
- [ ] tenant_resource.go
- [ ] tenant_group_resource.go
- [ ] contact_group_resource.go ✅ (already fixed - skip)
- [ ] contact_assignment_resource.go
- [ ] provider_account_resource.go
- [ ] provider_network_resource.go
- [ ] tenant_data_source.go
- [ ] tenant_group_data_source.go

---

## Batch 4: Site/Location Resources
**Status**: Not Started
**Items**: 9
- [ ] site_resource.go
- [ ] site_group_resource.go
- [ ] location_resource.go
- [ ] region_resource.go
- [ ] rack_resource.go
- [ ] rack_role_resource.go
- [ ] rack_type_resource.go
- [ ] site_data_source.go
- [ ] location_data_source.go

---

## Batch 5: Rack & Reservation Resources
**Status**: Not Started
**Items**: 7
- [ ] rack_reservation_resource.go ✅ (already fixed - skip)
- [ ] virtual_chassis_resource.go
- [ ] module_bay_resource.go
- [ ] module_bay_template_resource.go ✅ (already fixed - skip)
- [ ] cable_resource.go
- [ ] rack_data_source.go
- [ ] rack_reservation_data_source.go

---

## Batch 6: Network (IPAM) Resources
**Status**: Not Started
**Items**: 10
- [ ] aggregate_resource.go
- [ ] prefix_resource.go
- [ ] ip_address_resource.go
- [ ] ip_range_resource.go
- [ ] vrf_resource.go
- [ ] route_target_resource.go
- [ ] asn_resource.go
- [ ] asn_range_resource.go
- [ ] prefix_data_source.go
- [ ] ip_address_data_source.go

---

## Batch 7: VLAN Resources
**Status**: Not Started
**Items**: 7
- [ ] vlan_resource.go
- [ ] vlan_group_resource.go
- [ ] l2vpn_resource.go
- [ ] l2vpn_termination_resource.go
- [ ] service_resource.go
- [ ] vlan_data_source.go
- [ ] vlan_group_data_source.go

---

## Batch 8: Interface Resources
**Status**: Not Started
**Items**: 9
- [ ] interface_resource.go
- [ ] interface_template_resource.go
- [ ] front_port_resource.go
- [ ] front_port_template_resource.go
- [ ] rear_port_resource.go
- [ ] rear_port_template_resource.go
- [ ] console_port_resource.go
- [ ] interface_data_source.go
- [ ] front_port_data_source.go

---

## Batch 9: Port Templates & Power
**Status**: Not Started
**Items**: 10
- [ ] console_port_template_resource.go ✅ (already fixed - skip)
- [ ] console_server_port_resource.go
- [ ] console_server_port_template_resource.go ✅ (already fixed - skip)
- [ ] power_port_resource.go
- [ ] power_port_template_resource.go
- [ ] power_outlet_resource.go
- [ ] power_outlet_template_resource.go
- [ ] power_panel_resource.go
- [ ] power_port_data_source.go
- [ ] power_outlet_data_source.go

---

## Batch 10: Power Feed & Virtualization
**Status**: Not Started
**Items**: 8
- [ ] power_feed_resource.go
- [ ] virtual_device_context_resource.go
- [ ] virtual_machine_resource.go
- [ ] virtual_disk_resource.go
- [ ] vm_interface_resource.go
- [ ] cluster_resource.go ✅ (already fixed - skip)
- [ ] virtual_machine_data_source.go
- [ ] vm_interface_data_source.go

---

## Batch 11: Cluster Resources
**Status**: Not Started
**Items**: 7
- [ ] cluster_type_resource.go ✅ (already fixed - skip)
- [ ] cluster_group_resource.go
- [ ] inventory_item_resource.go ✅ (already fixed - skip)
- [ ] inventory_item_template_resource.go ✅ (already fixed - skip)
- [ ] inventory_item_role_resource.go
- [ ] cluster_data_source.go
- [ ] cluster_type_data_source.go

---

## Batch 12: Templates & Exports
**Status**: Not Started
**Items**: 8
- [ ] export_template_resource.go ✅ (already fixed - skip)
- [ ] config_template_resource.go ✅ (already fixed - skip)
- [ ] service_template_resource.go
- [ ] custom_field_resource.go
- [ ] custom_field_choice_set_resource.go
- [ ] config_context_resource.go
- [ ] export_template_data_source.go
- [ ] config_template_data_source.go

---

## Batch 13: Other Resources (Lower Priority)
**Status**: Not Started
**Items**: 10
- [ ] ike_policy_resource.go
- [ ] ike_proposal_resource.go
- [ ] ipsec_policy_resource.go
- [ ] ipsec_profile_resource.go
- [ ] ipsec_proposal_resource.go
- [ ] journal_entry_resource.go
- [ ] webhook_resource.go
- [ ] wireless_lan_group_resource.go
- [ ] ike_policy_data_source.go
- [ ] ipsec_policy_data_source.go

---

## Batch 14: Final Resources
**Status**: Not Started
**Items**: 7
- [ ] wireless_lan_resource.go
- [ ] wireless_link_resource.go
- [ ] rir_resource.go
- [ ] platform_resource.go
- [ ] manufacturer_resource.go
- [ ] wireless_lan_data_source.go
- [ ] platform_data_source.go

---

## Batch 15: Tags & Role Resources
**Status**: Not Started
**Items**: 6
- [ ] tag_resource.go
- [ ] role_resource.go
- [ ] event_rule_resource.go
- [ ] notification_group_resource.go
- [ ] tag_data_source.go
- [ ] role_data_source.go

---

## Notes
- Unit tests go in `internal/utils/state_helpers_test.go` (for UpdateReferenceAttribute pattern)
- Acceptance tests go in `internal/resources_acceptance_tests/{resource}_resource_test.go`
- Datasource tests go in `internal/datasources_test/{datasource}_data_source_test.go`
- Each batch should be committed separately for clarity
