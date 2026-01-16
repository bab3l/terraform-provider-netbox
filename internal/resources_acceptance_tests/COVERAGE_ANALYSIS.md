# Acceptance Test Coverage Analysis

This document tracks the current state of acceptance test coverage for all resources and identifies gaps that need to be addressed.

## Coverage Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Test exists and follows naming convention |
| ⚠️ | Test exists but may need review/renaming |
| ❌ | Test missing |
| N/A | Not applicable to this resource |

---

## Test Coverage Matrix

### TIER 1: Core CRUD Tests

| Resource | `_basic` | `_full` | `_update` | `_import` |
|----------|----------|---------|-----------|-----------|
| aggregate | ✅ | ✅ | ✅ | ✅ |
| asn | ✅ | ✅ | ✅ | ✅ |
| asn_range | ✅ | ✅ | ✅ | ✅ |
| cable | ✅ | ✅ | ✅ | ⚠️ |
| circuit | ✅ | ✅ | ✅ | ✅ |
| circuit_group | ✅ | ✅ | ✅ | ✅ |
| circuit_group_assignment | ✅ | ✅ | ✅ | ❌ |
| circuit_termination | ✅ | ✅ | ✅ | ❌ |
| circuit_type | ✅ | ✅ | ✅ | ✅ |
| cluster | ✅ | ✅ | ✅ | ✅ |
| cluster_group | ✅ | ✅ | ✅ | ✅ |
| cluster_type | ✅ | ✅ | ✅ | ✅ |
| config_context | ✅ | ✅ | ✅ | ✅ |
| config_template | ✅ | ✅ | ✅ | ✅ |
| console_port | ✅ | ✅ | ✅ | ❌ |
| console_port_template | ✅ | ✅ | ✅ | ❌ |
| console_server_port | ✅ | ✅ | ✅ | ❌ |
| console_server_port_template | ✅ | ✅ | ✅ | ❌ |
| contact | ✅ | ✅ | ✅ | ✅ |
| contact_assignment | ✅ | ✅ | ✅ | ❌ |
| contact_group | ✅ | ✅ | ✅ | ✅ |
| contact_role | ✅ | ✅ | ✅ | ✅ |
| custom_link | ✅ | ✅ | ✅ | ❌ |
| device | ✅ | ✅ | ✅ | ✅ |
| device_bay | ✅ | ✅ | ✅ | ❌ |
| device_bay_template | ✅ | ✅ | ✅ | ❌ |
| device_role | ✅ | ✅ | ✅ | ✅ |
| device_type | ✅ | ✅ | ✅ | ✅ |
| event_rule | ✅ | ✅ | ✅ | ❌ |
| export_template | ✅ | ✅ | ✅ | ❌ |
| fhrp_group | ✅ | ✅ | ✅ | ❌ |
| fhrp_group_assignment | ✅ | ✅ | ✅ | ❌ |
| front_port | ✅ | ✅ | ✅ | ❌ |
| front_port_template | ✅ | ✅ | ✅ | ❌ |
| ike_policy | ✅ | ✅ | ✅ | ❌ |
| ike_proposal | ✅ | ✅ | ✅ | ❌ |
| interface | ✅ | ✅ | ✅ | ✅ |
| interface_template | ✅ | ✅ | ✅ | ❌ |
| inventory_item | ✅ | ✅ | ✅ | ❌ |
| inventory_item_role | ✅ | ✅ | ✅ | ❌ |
| inventory_item_template | ✅ | ✅ | ✅ | ❌ |
| ip_address | ✅ | ✅ | ✅ | ✅ |
| ip_range | ✅ | ✅ | ✅ | ❌ |
| ipsec_policy | ✅ | ✅ | ✅ | ❌ |
| ipsec_profile | ✅ | ✅ | ✅ | ❌ |
| ipsec_proposal | ✅ | ✅ | ✅ | ❌ |
| journal_entry | ✅ | ✅ | ✅ | ❌ |
| l2vpn | ✅ | ✅ | ✅ | ✅ |
| l2vpn_termination | ✅ | ✅ | ✅ | ❌ |
| location | ✅ | ✅ | ✅ | ✅ |
| manufacturer | ✅ | ✅ | ✅ | ✅ |
| module | ✅ | ✅ | ✅ | ❌ |
| module_bay | ✅ | ✅ | ✅ | ❌ |
| module_bay_template | ✅ | ✅ | ✅ | ❌ |
| module_type | ✅ | ✅ | ✅ | ❌ |
| notification_group | ✅ | ✅ | ✅ | ✅ |
| platform | ✅ | ✅ | ✅ | ✅ |
| power_feed | ✅ | ✅ | ✅ | ❌ |
| power_outlet | ✅ | ✅ | ✅ | ❌ |
| power_outlet_template | ✅ | ✅ | ✅ | ❌ |
| power_panel | ✅ | ✅ | ✅ | ❌ |
| power_port | ✅ | ✅ | ✅ | ❌ |
| power_port_template | ✅ | ✅ | ✅ | ❌ |
| prefix | ✅ | ✅ | ✅ | ✅ |
| provider | ✅ | ✅ | ✅ | ✅ |
| provider_account | ✅ | ✅ | ✅ | ✅ |
| provider_network | ✅ | ✅ | ✅ | ✅ |
| rack | ✅ | ✅ | ✅ | ✅ |
| rack_reservation | ✅ | ✅ | ✅ | ❌ |
| rack_role | ✅ | ✅ | ✅ | ✅ |
| rack_type | ✅ | ✅ | ✅ | ❌ |
| rear_port | ✅ | ✅ | ✅ | ❌ |
| rear_port_template | ✅ | ✅ | ✅ | ❌ |
| region | ✅ | ✅ | ✅ | ✅ |
| rir | ✅ | ✅ | ✅ | ❌ |
| role | ✅ | ✅ | ✅ | ❌ |
| route_target | ✅ | ✅ | ✅ | ✅ |
| service | ✅ | ✅ | ✅ | ❌ |
| service_template | ✅ | ✅ | ✅ | ❌ |
| site | ✅ | ✅ | ✅ | ✅ |
| site_group | ✅ | ✅ | ✅ | ✅ |
| tag | ✅ | ✅ | ✅ | ❌ |
| tenant | ✅ | ✅ | ✅ | ✅ |
| tenant_group | ✅ | ✅ | ✅ | ✅ |
| tunnel | ✅ | ✅ | ✅ | ✅ |
| tunnel_group | ✅ | ✅ | ✅ | ✅ |
| tunnel_termination | ✅ | ✅ | ✅ | ✅ |
| virtual_chassis | ✅ | ✅ | ✅ | ❌ |
| virtual_device_context | ✅ | ✅ | ✅ | ❌ |
| virtual_disk | ✅ | ✅ | ✅ | ✅ |
| virtual_machine | ✅ | ✅ | ✅ | ❌ |
| vlan | ✅ | ✅ | ✅ | ✅ |
| vlan_group | ✅ | ✅ | ✅ | ✅ |
| vm_interface | ✅ | ✅ | ✅ | ✅ |
| vrf | ✅ | ✅ | ✅ | ✅ |
| webhook | ✅ | ✅ | ✅ | ❌ |
| wireless_lan | ✅ | ✅ | ✅ | ❌ |
| wireless_lan_group | ✅ | ✅ | ✅ | ❌ |
| wireless_link | ✅ | ✅ | ✅ | ❌ |

### TIER 2: Reliability Tests

| Resource | `_IDPreservation` | `_externalDeletion` | `_removeOptionalFields` |
|----------|-------------------|---------------------|-------------------------|
| aggregate | ✅ | ✅ | ✅ |
| asn | ✅ | ✅ | ✅ |
| asn_range | ✅ | ✅ | ✅ |
| cable | ✅ | ✅ | ✅ |
| circuit | ✅ | ✅ | ✅ |
| circuit_group | ✅ | ✅ | ✅ |
| circuit_group_assignment | ✅ | ✅ | ❌ |
| circuit_termination | ✅ | ✅ | ✅ |
| circuit_type | ✅ | ✅ | ✅ |
| cluster | ✅ | ✅ | ✅ |
| cluster_group | ✅ | ✅ | ✅ |
| cluster_type | ✅ | ✅ | ✅ |
| config_context | ✅ | ✅ | ✅ |
| config_template | ✅ | ✅ | ✅ |
| console_port | ✅ | ✅ | ✅ |
| console_port_template | ✅ | ✅ | ✅ |
| console_server_port | ✅ | ✅ | ✅ |
| console_server_port_template | ✅ | ✅ | ✅ |
| contact | ✅ | ✅ | ✅ |
| contact_assignment | ✅ | ✅ | ❌ |
| contact_group | ✅ | ✅ | ✅ |
| contact_role | ✅ | ✅ | ✅ |
| custom_link | ✅ | ✅ | ✅ |
| device | ✅ | ✅ | ✅ |
| device_bay | ✅ | ✅ | ✅ |
| device_bay_template | ✅ | ✅ | ✅ |
| device_role | ✅ | ✅ | ✅ |
| device_type | ✅ | ✅ | ✅ |
| event_rule | ✅ | ✅ | ✅ |
| export_template | ✅ | ✅ | ✅ |
| fhrp_group | ✅ | ✅ | ✅ |
| fhrp_group_assignment | ✅ | ✅ | ❌ |
| front_port | ✅ | ✅ | ✅ |
| front_port_template | ✅ | ✅ | ✅ |
| ike_policy | ✅ | ✅ | ✅ |
| ike_proposal | ✅ | ✅ | ✅ |
| interface | ✅ | ✅ | ✅ |
| interface_template | ✅ | ✅ | ✅ |
| inventory_item | ✅ | ✅ | ✅ |
| inventory_item_role | ✅ | ✅ | ✅ |
| inventory_item_template | ✅ | ✅ | ✅ |
| ip_address | ✅ | ✅ | ✅ |
| ip_range | ✅ | ✅ | ✅ |
| ipsec_policy | ✅ | ✅ | ✅ |
| ipsec_profile | ✅ | ✅ | ✅ |
| ipsec_proposal | ✅ | ✅ | ✅ |
| journal_entry | ✅ | ✅ | ❌ |
| l2vpn | ✅ | ✅ | ✅ |
| l2vpn_termination | ✅ | ⚠️ | ✅ |
| location | ✅ | ✅ | ✅ |
| manufacturer | ✅ | ✅ | ✅ |
| module | ✅ | ⚠️ | ✅ |
| module_bay | ✅ | ⚠️ | ✅ |
| module_bay_template | ✅ | ⚠️ | ✅ |
| module_type | ✅ | ⚠️ | ✅ |
| notification_group | ✅ | ✅ | ✅ |
| platform | ✅ | ✅ | ✅ |
| power_feed | ✅ | ✅ | ✅ |
| power_outlet | ✅ | ✅ | ✅ |
| power_outlet_template | ✅ | ✅ | ✅ |
| power_panel | ✅ | ✅ | ✅ |
| power_port | ✅ | ✅ | ✅ |
| power_port_template | ✅ | ✅ | ✅ |
| prefix | ✅ | ⚠️ | ✅ |
| provider | ✅ | ✅ | ✅ |
| provider_account | ✅ | ✅ | ✅ |
| provider_network | ✅ | ✅ | ✅ |
| rack | ✅ | ✅ | ⚠️ |
| rack_reservation | ✅ | ✅ | ✅ |
| rack_role | ✅ | ✅ | ✅ |
| rack_type | ✅ | ✅ | ✅ |
| rear_port | ✅ | ✅ | ✅ |
| rear_port_template | ✅ | ✅ | ✅ |
| region | ✅ | ✅ | ✅ |
| rir | ✅ | ⚠️ | ✅ |
| role | ✅ | ✅ | ✅ |
| route_target | ✅ | ✅ | ✅ |
| service | ✅ | ⚠️ | ✅ |
| service_template | ✅ | ⚠️ | ✅ |
| site | ✅ | ✅ | ✅ |
| site_group | ✅ | ✅ | ✅ |
| tag | ✅ | ✅ | ✅ |
| tenant | ✅ | ✅ | ✅ |
| tenant_group | ✅ | ✅ | ✅ |
| tunnel | ✅ | ✅ | ✅ |
| tunnel_group | ✅ | ✅ | ✅ |
| tunnel_termination | ✅ | ✅ | ✅ |
| virtual_chassis | ✅ | ✅ | ✅ |
| virtual_device_context | ✅ | ✅ | ✅ |
| virtual_disk | ✅ | ✅ | ✅ |
| virtual_machine | ✅ | ✅ | ✅ |
| vlan | ✅ | ✅ | ✅ |
| vlan_group | ✅ | ✅ | ✅ |
| vm_interface | ✅ | ⚠️ | ✅ |
| vrf | ✅ | ✅ | ✅ |
| webhook | ✅ | ✅ | ✅ |
| wireless_lan | ✅ | ✅ | ✅ |
| wireless_lan_group | ✅ | ✅ | ✅ |
| wireless_link | ✅ | ✅ | ✅ |

### TIER 3: Tag Tests (Resources with Tags)

| Resource | `_tagLifecycle` | `_tagOrderInvariance` |
|----------|-----------------|----------------------|
| aggregate | ❌ | ❌ |
| asn | ❌ | ❌ |
| asn_range | ❌ | ❌ |
| cable | ❌ | ❌ |
| circuit | ❌ | ❌ |
| circuit_termination | ❌ | ❌ |
| circuit_type | ❌ | ❌ |
| cluster | ❌ | ❌ |
| cluster_group | ❌ | ❌ |
| cluster_type | ❌ | ❌ |
| config_context | ❌ | ❌ |
| config_template | ❌ | ❌ |
| console_port | ❌ | ❌ |
| console_port_template | ❌ | ❌ |
| console_server_port | ❌ | ❌ |
| console_server_port_template | ❌ | ❌ |
| contact | ❌ | ❌ |
| contact_group | ❌ | ❌ |
| contact_role | ❌ | ❌ |
| device | ❌ | ❌ |
| device_bay | ❌ | ❌ |
| device_bay_template | ❌ | ❌ |
| device_role | ❌ | ❌ |
| device_type | ❌ | ❌ |
| fhrp_group | ❌ | ❌ |
| front_port | ❌ | ❌ |
| front_port_template | ❌ | ❌ |
| interface | ❌ | ❌ |
| interface_template | ❌ | ❌ |
| inventory_item | ❌ | ❌ |
| inventory_item_role | ❌ | ❌ |
| inventory_item_template | ❌ | ❌ |
| ip_address | ⚠️ | ✅ |
| ip_range | ❌ | ❌ |
| l2vpn | ❌ | ❌ |
| l2vpn_termination | ❌ | ❌ |
| location | ❌ | ❌ |
| manufacturer | ❌ | ❌ |
| module | ❌ | ❌ |
| module_bay | ❌ | ❌ |
| module_bay_template | ❌ | ❌ |
| module_type | ❌ | ❌ |
| platform | ❌ | ❌ |
| power_feed | ❌ | ❌ |
| power_outlet | ❌ | ❌ |
| power_outlet_template | ❌ | ❌ |
| power_panel | ❌ | ❌ |
| power_port | ❌ | ❌ |
| power_port_template | ❌ | ❌ |
| prefix | ❌ | ❌ |
| provider | ❌ | ❌ |
| provider_account | ❌ | ❌ |
| provider_network | ❌ | ❌ |
| rack | ❌ | ❌ |
| rack_reservation | ❌ | ❌ |
| rack_role | ❌ | ❌ |
| rack_type | ❌ | ❌ |
| rear_port | ❌ | ❌ |
| rear_port_template | ❌ | ❌ |
| region | ❌ | ❌ |
| rir | ❌ | ❌ |
| role | ❌ | ❌ |
| route_target | ❌ | ❌ |
| service | ❌ | ❌ |
| service_template | ❌ | ❌ |
| site | ❌ | ❌ |
| site_group | ❌ | ❌ |
| tenant | ❌ | ❌ |
| tenant_group | ❌ | ❌ |
| tunnel | ❌ | ❌ |
| tunnel_group | ❌ | ❌ |
| tunnel_termination | ❌ | ❌ |
| virtual_chassis | ❌ | ❌ |
| virtual_device_context | ❌ | ❌ |
| virtual_disk | ❌ | ❌ |
| virtual_machine | ❌ | ❌ |
| vlan | ❌ | ❌ |
| vlan_group | ❌ | ❌ |
| vm_interface | ❌ | ❌ |
| vrf | ❌ | ❌ |
| webhook | ❌ | ❌ |
| wireless_lan | ❌ | ❌ |
| wireless_lan_group | ❌ | ❌ |
| wireless_link | ❌ | ❌ |

---

## Naming Convention Issues

The following tests have naming inconsistencies that should be addressed:

### External Deletion Tests (should be `_externalDeletion`)
- `l2vpn_termination`: `_external_deletion` → `_externalDeletion`
- `module`: `_external_deletion` → `_externalDeletion`
- `module_bay`: `_external_deletion` → `_externalDeletion`
- `module_bay_template`: `_external_deletion` → `_externalDeletion`
- `module_type`: `_external_deletion` → `_externalDeletion`
- `prefix`: `_external_deletion` → `_externalDeletion`
- `rir`: `_external_deletion` → `_externalDeletion`
- `service`: `_external_deletion` → `_externalDeletion`
- `service_template`: `_external_deletion` → `_externalDeletion`
- `vm_interface`: `_external_deletion` → `_externalDeletion`

---

## Priority Work Items

### Phase 1: Fix Tag Removal Bug (COMPLETED)
- [x] Fix `TagsToNestedTagRequests` to return empty slice instead of nil
- [x] Add `TestAccIPAddressResource_tagRemoval` test
- [x] Add `TestAccIPAddressResource_tagOrderInvariance` test

### Phase 2: Standardize IP Address Tests as Reference Implementation
- [x] `_tagRemoval` (add → remove → verify)
- [x] `_createWithTags` (create with tags)
- [x] `_modifyTags` (change tags)
- [x] `_tagOrderInvariance` (order doesn't matter)
- [ ] Rename to `_tagLifecycle` (consolidate into single comprehensive test)

### Phase 3: Apply Tag Tests to All Resources with Tags
Priority order (most commonly used resources first):
1. `virtual_machine`
2. `device`
3. `prefix`
4. `vlan`
5. `site`
6. `interface`
7. `cluster`
8. ... (remaining resources)

### Phase 4: Fix Naming Inconsistencies
- Rename `_external_deletion` tests to `_externalDeletion`
- Standardize config function naming

### Phase 5: Add Missing Import Tests
Resources needing `_import` tests (42 resources):
- circuit_group_assignment
- circuit_termination
- console_port
- console_port_template
- console_server_port
- console_server_port_template
- contact_assignment
- custom_link
- device_bay
- device_bay_template
- event_rule
- export_template
- fhrp_group
- fhrp_group_assignment
- front_port
- front_port_template
- ike_policy
- ike_proposal
- interface_template
- inventory_item
- inventory_item_role
- inventory_item_template
- ip_range
- ipsec_policy
- ipsec_profile
- ipsec_proposal
- journal_entry
- l2vpn_termination
- module
- module_bay
- module_bay_template
- module_type
- power_feed
- power_outlet
- power_outlet_template
- power_panel
- power_port
- power_port_template
- rack_reservation
- rack_type
- rear_port
- rear_port_template
- rir
- role
- service
- service_template
- tag
- virtual_chassis
- virtual_device_context
- virtual_machine
- webhook
- wireless_lan
- wireless_lan_group
- wireless_link

---

## Statistics

### Current Coverage Summary

| Test Category | Implemented | Total Required | Coverage % |
|---------------|-------------|----------------|------------|
| `_basic` | 86 | 86 | 100% |
| `_full` | 86 | 86 | 100% |
| `_update` | 86 | 86 | 100% |
| `_import` | 44 | 86 | 51% |
| `_IDPreservation` | 86 | 86 | 100% |
| `_externalDeletion` | 86 | 86 | 100% |
| `_removeOptionalFields` | 82 | 86 | 95% |
| `_tagLifecycle` | 1 | 73 | 1% |
| `_tagOrderInvariance` | 1 | 73 | 1% |

### Tag Test Gap
- **73 resources** support tags
- **Only 1 resource** (ip_address) has tag lifecycle tests
- **72 resources** need tag tests added

---

*Last updated: January 2026*
