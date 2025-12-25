# ID Preservation Testing - Batch Tracking

## Overview
Testing all 101 resources and 104 datasources to ensure ID is preserved as the immutable identifier.

## Progress Summary
- **Total Items**: 205 (101 Resources + 104 Datasources)
- **Completed**: 57 items (27.8%)
- **Remaining**: ~148 items across 8 remaining batches

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
**Status**: COMPLETE ✅
**Items**: 9
- [x] device_resource.go ✅ (already fixed - skip)
- [x] device_type_resource.go ✅ (already fixed - skip)
- [x] device_role_resource.go ✅ (already fixed - skip)
- [x] device_bay_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] device_bay_template_resource.go ✅ (already fixed - skip)
- [x] module_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] module_type_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] device_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing
- [x] device_type_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing

**Tests Added**:
- `internal/resources_acceptance_tests/device_bay_resource_test.go`: `TestAccDeviceBayResource_IDPreservation`
- `internal/resources_acceptance_tests/module_resource_test.go`: `TestAccModuleResource_IDPreservation`
- `internal/resources_acceptance_tests/module_type_resource_test.go`: `TestAccModuleTypeResource_IDPreservation`
- `internal/datasources_acceptance_tests/device_data_source_test.go`: `TestAccDeviceDataSource_IDPreservation`
- `internal/datasources_acceptance_tests/device_type_data_source_test.go`: `TestAccDeviceTypeDataSource_IDPreservation`

**Test Results**: All tests PASSING ✅
- DeviceBay resource test: 41.56s ✅
- Module resource test: 44.06s ✅
- ModuleType resource test: 43.69s ✅
- Device datasource test: 44.22s ✅
- DeviceType datasource test: 75.43s ✅

---

## Batch 3: Tenancy/Organization Resources
**Status**: COMPLETE ✅
**Items**: 8
- [x] tenant_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] tenant_group_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] contact_group_resource.go ✅ (already fixed - skip)
- [x] contact_assignment_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] provider_account_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] provider_network_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] tenant_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing
- [x] tenant_group_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing

**Tests Added**:
- `internal/resources_acceptance_tests/tenant_resource_test.go`: `TestAccTenantResource_IDPreservation`
- `internal/resources_acceptance_tests/tenant_group_resource_test.go`: `TestAccTenantGroupResource_IDPreservation`
- `internal/resources_acceptance_tests/contact_assignment_resource_test.go`: `TestAccContactAssignmentResource_IDPreservation`
- `internal/resources_acceptance_tests/provider_account_resource_test.go`: `TestAccProviderAccountResource_IDPreservation`
- `internal/resources_acceptance_tests/provider_network_resource_test.go`: `TestAccProviderNetworkResource_IDPreservation`
- `internal/datasources_acceptance_tests/tenant_data_source_test.go`: `TestAccTenantDataSource_IDPreservation`
- `internal/datasources_acceptance_tests/tenant_group_data_source_test.go`: `TestAccTenantGroupDataSource_IDPreservation`

**Test Results**: All tests PASSING ✅
- Tenant resource test: 49.84s ✅
- TenantGroup resource test: 49.54s ✅
- ContactAssignment resource test: 50.02s ✅
- ProviderAccount resource test: 49.97s ✅
- ProviderNetwork resource test: 51.40s ✅
- Tenant datasource test: 75.43s ✅
- TenantGroup datasource test: 75.13s ✅

---

## Batch 4: Site/Location Resources
**Status**: COMPLETE ✅
**Items**: 9
- [x] site_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] site_group_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] location_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] region_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] rack_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] rack_role_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] rack_type_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] site_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing
- [x] location_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing

**Tests Added**:
- `internal/resources_acceptance_tests/site_resource_test.go`: `TestAccSiteResource_IDPreservation`
- `internal/resources_acceptance_tests/site_group_resource_test.go`: `TestAccSiteGroupResource_IDPreservation`
- `internal/resources_acceptance_tests/location_resource_test.go`: `TestAccLocationResource_IDPreservation`
- `internal/resources_acceptance_tests/region_resource_test.go`: `TestAccRegionResource_IDPreservation`
- `internal/resources_acceptance_tests/rack_resource_test.go`: `TestAccRackResource_IDPreservation`
- `internal/resources_acceptance_tests/rack_role_resource_test.go`: `TestAccRackRoleResource_IDPreservation`
- `internal/resources_acceptance_tests/rack_type_resource_test.go`: `TestAccRackTypeResource_IDPreservation`
- `internal/datasources_acceptance_tests/site_data_source_test.go`: `TestAccSiteDataSource_IDPreservation`
- `internal/datasources_acceptance_tests/location_data_source_test.go`: `TestAccLocationDataSource_IDPreservation`

**Test Results**: All tests PASSING ✅
- Site resource test: 39.93s ✅
- SiteGroup resource test: 42.04s ✅
- Location resource test: 75.96s ✅
- Region resource test: 76.40s ✅
- Rack resource test: 77.06s ✅
- RackRole resource test: 39.53s ✅
- RackType resource test: 39.80s ✅
- Site datasource test: 54.02s ✅
- Location datasource test: 55.05s ✅

---

## Batch 5: Rack & Reservation Resources
**Status**: COMPLETE ✅
**Items**: 7 (2 already fixed, 5 new tests)
- [x] rack_reservation_resource.go ✅ (already fixed - skip)
- [x] virtual_chassis_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] module_bay_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] module_bay_template_resource.go ✅ (already fixed - skip)
- [x] cable_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] rack_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing
- [x] rack_reservation_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing

**Tests Added**:
- `internal/resources_acceptance_tests/virtual_chassis_resource_test.go`: `TestAccVirtualChassisResource_IDPreservation`
- `internal/resources_acceptance_tests/module_bay_resource_test.go`: `TestAccModuleBayResource_IDPreservation`
- `internal/resources_acceptance_tests/cable_resource_test.go`: `TestAccCableResource_IDPreservation`
- `internal/datasources_acceptance_tests/rack_data_source_test.go`: `TestAccRackDataSource_IDPreservation`
- `internal/datasources_acceptance_tests/rack_reservation_data_source_test.go`: `TestAccRackReservationDataSource_IDPreservation`

**Test Results**: All tests PASSING ✅
- VirtualChassis resource test: 23.78s ✅
- ModuleBay resource test: 26.95s ✅
- Cable resource test: 25.43s ✅
- Rack datasource test: 54.41s ✅
- RackReservation datasource test: 24.29s ✅

---

## Batch 6: Network (IPAM) Resources
**Status**: COMPLETE ✅
**Items**: 10
- [x] aggregate_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] prefix_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] ip_address_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] ip_range_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] vrf_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] route_target_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] asn_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] asn_range_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] prefix_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing
- [x] ip_address_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing

**Tests Added**:
- `internal/resources_acceptance_tests/aggregate_resource_test.go`: `TestAccAggregateResource_IDPreservation`
- `internal/resources_acceptance_tests/prefix_resource_test.go`: `TestAccPrefixResource_IDPreservation`
- `internal/resources_acceptance_tests/ip_address_resource_test.go`: `TestAccIPAddressResource_IDPreservation`
- `internal/resources_acceptance_tests/ip_range_resource_test.go`: `TestAccIPRangeResource_IDPreservation`
- `internal/resources_acceptance_tests/vrf_resource_test.go`: `TestAccVRFResource_IDPreservation`
- `internal/resources_acceptance_tests/route_target_resource_test.go`: `TestAccRouteTargetResource_IDPreservation`
- `internal/resources_acceptance_tests/asn_resource_test.go`: `TestAccASNResource_IDPreservation`
- `internal/resources_acceptance_tests/asn_range_resource_test.go`: `TestAccASNRangeResource_IDPreservation`
- `internal/datasources_acceptance_tests/prefix_data_source_test.go`: `TestAccPrefixDataSource_IDPreservation`
- `internal/datasources_acceptance_tests/ip_address_data_source_test.go`: `TestAccIPAddressDataSource_IDPreservation`

**Test Results**: All tests PASSING ✅
- Aggregate resource test: 42.27s ✅
- Prefix resource test: 42.29s ✅
- IPAddress resource test: 41.94s ✅
- IPRange resource test: 42.13s ✅
- VRF resource test: 41.20s ✅
- RouteTarget resource test: 42.10s ✅
- ASN resource test: 42.39s ✅
- ASNRange resource test: 42.37s ✅
- Prefix datasource test: 25.69s ✅
- IPAddress datasource test: 28.21s ✅

---

## Batch 7: VLAN Resources
**Status**: COMPLETE ✅
**Items**: 7
- [x] vlan_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] vlan_group_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] l2vpn_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] l2vpn_termination_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] service_resource.go ✅ **COMPLETE** - Acceptance test added and passing
- [x] vlan_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing
- [x] vlan_group_data_source.go ✅ **COMPLETE** - Datasource ID preservation test added and passing

**Tests Added**:
- `internal/resources_acceptance_tests/vlan_resource_test.go`: `TestAccVLANResource_IDPreservation`
- `internal/resources_acceptance_tests/vlan_group_resource_test.go`: `TestAccVLANGroupResource_IDPreservation`
- `internal/resources_acceptance_tests/l2vpn_resource_test.go`: `TestAccL2VPNResource_IDPreservation`
- `internal/resources_acceptance_tests/l2vpn_termination_resource_test.go`: `TestAccL2VPNTerminationResource_IDPreservation`
- `internal/resources_acceptance_tests/service_resource_test.go`: `TestAccServiceResource_IDPreservation`
- `internal/datasources_acceptance_tests/vlan_data_source_test.go`: `TestAccVLANDataSource_IDPreservation`
- `internal/datasources_acceptance_tests/vlan_group_data_source_test.go`: `TestAccVLANGroupDataSource_IDPreservation`

**Test Results**: All tests PASSING ✅
- VLAN resource test: 34.91s ✅
- VLANGroup resource test: 34.08s ✅
- L2VPN resource test: 32.52s ✅
- L2VPNTermination resource test: 33.39s ✅
- Service resource test: 33.84s ✅
- VLAN datasource test: 23.73s ✅
- VLANGroup datasource test: 23.29s ✅

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
