# Acceptance Test Coverage Analysis

## Overall Progress
**Status**: 73/86 resources complete (84.9%)

## Completed Resources

### 1. IP Address (ipam_ipaddress)
- 11 tests passing
- Duration: ~15s
- Checklist: IPADDRESS_CHECKLIST.md

### 2. Prefix (ipam_prefix)
- 13 tests passing
- Duration: ~15s
- Checklist: PREFIX_CHECKLIST.md

### 3. Aggregate (ipam_aggregate)
- 9 tests passing
- Duration: ~10s
- Checklist: AGGREGATE_CHECKLIST.md

### 4. ASN (ipam_asn)
- 8 tests passing
- Duration: ~10s
- Checklist: ASN_CHECKLIST.md

### 5. ASN Range (ipam_asn_range)
- 10 tests passing
- Duration: ~12s
- Checklist: ASN_RANGE_CHECKLIST.md

### 6. Cable (dcim_cable)
- 10 tests passing
- Duration: ~20s
- Checklist: CABLE_CHECKLIST.md
- **Notable**: Fixed provider-wide tag lifecycle bug during implementation

### 7. Circuit (circuits_circuit)
- 10 tests passing
- Duration: ~15s
- Checklist: CIRCUIT_CHECKLIST.md

### 8. Circuit Group (circuits_circuit_group)
- 9 tests passing
- Duration: ~9s
- Checklist: CIRCUIT_GROUP_CHECKLIST.md

### 9. Circuit Termination (circuits_circuit_termination)
- 9 tests passing
- Duration: ~7.5s
- Checklist: CIRCUIT_TERMINATION_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 10. Circuit Type (circuits_circuit_type)
- 9 tests passing
- Duration: ~7s
- Checklist: CIRCUIT_TYPE_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 11. Cluster (virtualization_cluster)
- 10 tests passing (plus 1 extended variant)
- Duration: ~10.9s
- Checklist: CLUSTER_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 12. Cluster Group (virtualization_cluster_group)
- 8 tests passing
- Duration: ~5.6s
- Checklist: CLUSTER_GROUP_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 13. Cluster Type (virtualization_cluster_type)
- 7 tests passing
- Duration: ~3.1s
- Checklist: CLUSTER_TYPE_CHECKLIST.md
- **Notable**: No tag support (simple resource)

### 14. Config Context (extras_config_context)
- 8 tests passing
- Duration: ~6.9s
- Checklist: CONFIG_CONTEXT_CHECKLIST.md
- **Notable**: Uses slug list tag format, complex dependencies

### 15. Config Template (extras_config_template)
- 6 tests passing
- Duration: ~2.9s
- Checklist: CONFIG_TEMPLATE_CHECKLIST.md
- **Notable**: No tag support (simple resource)

### 16. Console Port (dcim_console_port)
- 6 tests passing
- Duration: ~6.0s
- Checklist: CONSOLE_PORT_CHECKLIST.md
- **Notable**: No tag support, requires device dependency

### 17. Console Port Template (dcim_console_port_template)
- 6 tests passing (plus 1 extended variant)
- Duration: ~6.8s
- Checklist: CONSOLE_PORT_TEMPLATE_CHECKLIST.md
- **Notable**: No tag support, requires device type dependency

### 18. Console Server Port (dcim_console_server_port)
- 6 tests passing
- Duration: ~7.1s
- Checklist: CONSOLE_SERVER_PORT_CHECKLIST.md
- **Notable**: No tag support, requires device dependency

### 19. Console Server Port Template (dcim_console_server_port_template)
- 6 tests passing (plus 1 extended variant)
- Duration: ~6.8s
- Checklist: CONSOLE_SERVER_PORT_TEMPLATE_CHECKLIST.md
- **Notable**: No tag support, requires device type dependency

### 20. Contact Assignment (extras_contact_assignment)
- 9 tests passing
- Duration: ~10.7s
- Checklist: CONTACT_ASSIGNMENT_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), generic assignment resource with ContentType pattern

### 21. Contact Group (tenancy_contact_group)
- 8 tests passing
- Duration: ~7.0s
- Checklist: CONTACT_GROUP_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), supports hierarchical parent relationships

### 22. Contact (tenancy_contact)
- 6 tests passing (plus 1 extended variant)
- Duration: ~5.0s
- Checklist: CONTACT_CHECKLIST.md
- **Notable**: No tag support, rich contact information fields, supports group references

### 23. Contact Role (tenancy_contact_role)
- 8 tests passing
- Duration: ~5.3s
- Checklist: CONTACT_ROLE_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), defines roles for contacts

### 24. Custom Link (extras_custom_link)
- 6 tests passing
- Duration: ~3.7s
- Checklist: CUSTOM_LINK_CHECKLIST.md
- **Notable**: No tag support, extensibility feature for adding custom links to NetBox UI

### 25. Device Bay (dcim_device_bay)
- 6 tests passing (plus 1 extended variant)
- Duration: ~9.1s
- Checklist: DEVICE_BAY_CHECKLIST.md
- **Notable**: No tag support, component resource for device bays, complex dependency chain

### 26. Device Bay Template (dcim_device_bay_template)
- 6 tests passing
- Duration: ~3.7s
- Checklist: DEVICE_BAY_TEMPLATE_CHECKLIST.md
- **Notable**: No tag support, template resource for device types

### 27. Device Role (dcim_device_role)
- 6 tests passing
- Duration: ~3.2s
- Checklist: DEVICE_ROLE_CHECKLIST.md
- **Notable**: No tag support, core organizational resource, standalone with no dependencies

### 28. Device Type (dcim_device_type)
- 7 tests passing
- Duration: ~7.0s
- Checklist: DEVICE_TYPE_CHECKLIST.md
- **Notable**: No tag support, hardware specification resource, requires manufacturer dependency

### 29. Device (dcim_device)
- 10 tests passing
- Duration: ~12.5s
- Checklist: DEVICE_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), core physical device resource, complex dependencies

### 30. Event Rule (extras_event_rule)
- 7 tests passing (plus 1 extended variant)
- Duration: ~4.6s
- Checklist: EVENT_RULE_CHECKLIST.md
- **Notable**: No tag support, automation/workflow resource, requires webhook dependency

### 31. Export Template (extras_export_template)
- 7 tests passing (6 regular + 1 extended)
- Duration: ~3.5s
- Checklist: EXPORT_TEMPLATE_CHECKLIST.md
- **Notable**: No tag support, Jinja2 template-based resource for data export

### 32. FHRP Group Assignment (ipam_fhrp_group_assignment)
- 5 tests passing (4 regular + 1 with validation subtests)
- Duration: ~5.2s
- Checklist: FHRP_GROUP_ASSIGNMENT_CHECKLIST.md
- **Notable**: No tag support, junction resource linking FHRP groups to interfaces

### 33. FHRP Group (ipam_fhrp_group)
- 9 tests passing (7 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~5.3s
- Checklist: FHRP_GROUP_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), First Hop Redundancy Protocol (VRRP/HSRP) resource

### 34. Front Port (dcim_front_port)
- 8 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.7s
- Checklist: FRONT_PORT_CHECKLIST.md
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), requires rear port dependency, physical port resource

### 35. Front Port Template (dcim_front_port_template)
- 6 tests passing (5 regular + 1 with validation subtests)
- Duration: ~6.9s
- Checklist: FRONT_PORT_TEMPLATE_CHECKLIST.md
- **Notable**: Template resource, does not support tags, requires device type and rear port template dependencies

### 36. IKE Policy (ipam_ike_policy)
- 9 tests passing (7 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~4.1s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), IPsec IKE policy resource

### 37. IKE Proposal (ipam_ike_proposal)
- 9 tests passing (7 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~7.5s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), IPsec IKE proposal resource

### 38. Interface (dcim_interface)
- 12 tests passing (10 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~18.4s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), complex device dependencies, includes optional+computed field coverage

### 39. Interface Template (dcim_interface_template)
- 8 tests passing (7 regular + 1 with validation subtests)
- Duration: ~7.1s
- **Notable**: Template resource, does not support tags, includes optional+computed field coverage

### 40. Inventory Item (dcim_inventory_item)
- 10 tests passing (8 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.2s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), device-bound hardware inventory

### 41. Inventory Item Role (dcim_inventory_item_role)
- 9 tests passing (7 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~8.5s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), inventory role metadata

### 42. Inventory Item Template (dcim_inventory_item_template)
- 7 tests passing (6 regular + 1 with validation subtests)
- Duration: ~3.8s
- **Notable**: Template resource, does not support tags

### 43. IP Range (ipam_ip_range)
- 12 tests passing (10 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.4s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), optional+computed field coverage

### 44. IPSec Policy (ipam_ipsec_policy)
- 7 tests passing (5 regular + 2 tag tests)
- Duration: ~8.6s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), IPsec policy metadata

### 45. IPSec Profile (ipam_ipsec_profile)
- 7 tests passing (5 regular + 2 tag tests)
- Duration: ~8.5s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), depends on IKE and IPSec policies

### 46. IPSec Proposal (ipam_ipsec_proposal)
- 7 tests passing (5 regular + 2 tag tests)
- Duration: ~7.8s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), IPsec proposal metadata

### 47. L2VPN (ipam_l2vpn)
- 8 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.1s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), L2VPN service resource

### 48. L2VPN Termination (ipam_l2vpn_termination)
- 7 tests passing (5 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~8.0s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), remove optional fields test skipped

### 49. Location (dcim_location)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~8.7s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), supports hierarchical parent

### 50. Manufacturer (dcim_manufacturer)
- 10 tests passing (7 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~8.7s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), core hardware metadata

### 51. Module Bay (dcim_module_bay)
- 8 tests passing (6 regular + 1 extended variant + 1 with validation subtests)
- Duration: ~7.9s
- **Notable**: No tag support, device dependency chain

### 52. Module Bay Template (dcim_module_bay_template)
- 7 tests passing (5 regular + 1 extended variant + 1 with validation subtests)
- Duration: ~4.0s
- **Notable**: No tag support, device type dependency

### 53. Module (dcim_module)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~8.8s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), device/module dependency chain

### 54. Module Type (dcim_module_type)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~7.1s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), manufacturer dependency

### 55. Notification Group (extras_notification_group)
- 7 tests passing (5 regular + 1 extended variant + 1 with validation subtests)
- Duration: ~3.2s
- **Notable**: No tag support, handles users and groups

### 56. Platform (dcim_platform)
- 7 tests passing (6 regular + 1 with validation subtests)
- Duration: ~6.1s
- **Notable**: No tag support, manufacturer association optional

### 57. Power Feed (dcim_power_feed)
- 11 tests passing (8 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~16.8s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), includes optional+computed voltage/amperage coverage

### 58. Power Outlet (dcim_power_outlet)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~7.8s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), device dependency

### 59. Power Outlet Template (dcim_power_outlet_template)
- 7 tests passing (6 regular + 1 with validation subtests)
- Duration: ~3.5s
- **Notable**: No tag support, device type dependency

### 60. Power Panel (dcim_power_panel)
- 8 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.6s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), site/location dependency

### 61. Power Port (dcim_power_port)
- 10 tests passing (7 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~13.0s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), device dependency

### 62. Power Port Template (dcim_power_port_template)
- 7 tests passing (6 regular + 1 with validation subtests)
- Duration: ~3.6s
- **Notable**: No tag support, device type dependency

### 63. Provider Account (circuits_provider_account)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~4.7s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), provider dependency

### 64. Provider Network (circuits_provider_network)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.5s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), provider dependency

### 65. Provider (circuits_provider)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.3s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 66. Rack Reservation (dcim_rack_reservation)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.7s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 67. Rack (dcim_rack)
- 11 tests passing (8 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~11.8s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), extended optional/reference field coverage

### 68. Rack Role (dcim_rack_role)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.1s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 69. Rack Type (dcim_rack_type)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~10.4s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

### 70. Rear Port (dcim_rear_port)
- 11 tests passing (8 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~12.5s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), optional+computed positions coverage

### 71. Rear Port Template (dcim_rear_port_template)
- 8 tests passing (7 regular + 1 with validation subtests)
- Duration: ~10.2s
- **Notable**: No tag support, optional+computed positions coverage

### 72. Region (dcim_region)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~9.5s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2), supports parent hierarchy

### 73. RIR (ipam_rir)
- 9 tests passing (6 regular + 2 tag tests + 1 with validation subtests)
- Duration: ~6.5s
- **Notable**: ⚠️ Uses nested tag format - needs conversion to slug list (Phase 2)

## Standard Test Pattern

Each resource includes:
1. **Core CRUD**: basic, full, update, import (4 tests)
2. **Reliability**: external deletion, remove optional fields (2 tests)
3. **Validation**: Validation error handling (1 test, recommended)
4. **Tag Tests**: Tag lifecycle and order invariance (2 tests if resource supports tags)
5. **Total**: 8-10 tests per resource (varies by resource complexity)

**Note**: IDPreservation test was removed as it was a duplicate of the basic test.

## Bug Fixes Applied

### Tag Lifecycle Bug (Fixed in Cable)
- **Issue**: Provider couldn't transition from tags to no tags
- **Root Cause**: `ApplyCommonFieldsWithMerge` preserved state tags when plan had null
- **Fix**: Always use plan tags, send empty array when null
- **Impact**: Affects all resources using this helper function
- **Files Modified**:
  - `internal/utils/request_helpers.go` (ApplyCommonFieldsWithMerge, ApplyTags)

## Next Resource
Continue alphabetically through remaining 38 resources.

## Estimated Completion
- At current pace: ~4-5 resources per session
- Estimated total time: ~15-18 sessions

## Post-Standardization Tasks

### Tag Format Standardization (Phase 2)
**Status**: Planned - to be executed after all test standardization is complete

**Problem**: Resources currently use two different tag formats in Terraform HCL:
1. **Nested object format**: `tags = [{ name = ..., slug = ... }]`
2. **Slug list format**: `tags = [slug1, slug2]`

**Decision**: Standardize ALL resources to use the simpler **slug list format**

**Resources requiring conversion** (nested → slug list):
- Circuit Termination (resource 9)
- Circuit Type (resource 10)
- Cluster (resource 11)
- Cluster Group (resource 12)
- Contact Assignment (resource 20)
- Contact Group (resource 21)
- Contact Role (resource 23)
- Device (resource 29)
- FHRP Group (resource 33)
- Front Port (resource 34)
- IKE Policy (resource 36)
- IKE Proposal (resource 37)
- Interface (resource 38)
- Inventory Item (resource 40)
- Inventory Item Role (resource 41)
- Location (resource 49)
- Manufacturer (resource 50)
- Module (resource 53)
- Module Type (resource 54)
- Power Feed (resource 57)
- Power Outlet (resource 58)
- Power Panel (resource 60)
- Power Port (resource 61)
- Provider Account (resource 63)
- Provider Network (resource 64)
- Provider (resource 65)
- Rack Reservation (resource 66)
- Rack (resource 67)
- Rack Role (resource 68)
- Rack Type (resource 69)
- Rear Port (resource 70)
- Region (resource 72)
- RIR (resource 73)
- IP Range (resource 43)
- IPSec Policy (resource 44)
- IPSec Profile (resource 45)
- IPSec Proposal (resource 46)
- L2VPN (resource 47)
- L2VPN Termination (resource 48)
- L2VPN Termination (resource 48)
- L2VPN (resource 47)
- IPSec Proposal (resource 46)
- IPSec Profile (resource 45)
- IPSec Policy (resource 44)
- Inventory Item Template (resource 42)
- IP Range (resource 43)
- Inventory Item Role (resource 41)
- Inventory Item (resource 40)
- Interface (resource 38)

**Action Items** (after test standardization complete):
1. Identify all resources using nested tag format
2. Update resource schemas to accept slug lists
3. Update resource CRUD logic to work with slug lists
4. Update all test files to use slug list format
5. Update documentation and examples
6. Create migration guide for users (breaking change)
7. Update CHANGELOG with breaking change notice

**Rationale**: Simpler user experience, less confusion, more consistent with majority of resources
