# Optional Field Test Coverage Plan

**Generated:** 2026-01-13
**Purpose:** Track and plan improvements to `_removeOptionalFields` test coverage across all resources

## Executive Summary

- **Total Resources with Optional Fields:** 76
- **Resources Needing New Tests:** 3 (no `_removeOptionalFields` test exists)
- **Resources Needing Extended Tests:** 47 (test exists but incomplete coverage)
- **Recently Completed:** 27 resources (circuit_type, device_bay_template, circuit_group_assignment, aggregate, contact_assignment, power_panel, rack_reservation, virtual_chassis, vlan_group, journal_entry, rir, service_template, tag, cable, asn, circuit_termination, custom_link, fhrp_group, module, module_type, tunnel_termination, virtual_device_context, device, interface, power_feed)

## Priority Classification

### 🔴 Priority 1: No Test Coverage (3 Resources)

These resources have NO `_removeOptionalFields` test. A new test must be created.

| Resource | Missing Fields | Count |
|----------|---------------|-------|
| `custom_field` | choice_set, default, filter_logic, group_name, is_cloneable, label, related_object_type, required, search_weight, ui_editable, ui_visible, validation_maximum, validation_minimum, validation_regex, weight | 15 |
| `rack_type` | desc_units, form_factor, max_weight, mounting_depth, outer_depth, outer_unit, outer_width, starting_unit, u_height, weight, weight_unit, width | 12 |
| `wireless_lan` | auth_cipher, auth_psk, auth_type, comments, description, group, status, tenant, vlan | 9 |

**Total Missing Fields:** 36

---

### 🟡 Priority 2: Partial Coverage (50 Resources)

These resources have a `_removeOptionalFields` test but don't cover all optional fields.

#### High Impact (5+ missing fields)

| Resource | Currently Tested | Missing | Count |
|----------|-----------------|---------|-------|
| `config_context` | description, sites, tags | cluster_groups, cluster_types, clusters, device_types, is_active, locations, platforms, regions, roles, site_groups, tenant_groups, tenants, weight | 13 |
| `custom_field` | *(no test)* | *(see Priority 1)* | 15 |
| `inventory_item_template` | label | component_id, component_type, manufacturer, parent, part_id, role | 6 |
| `power_feed` | *(no test)* | *(see Priority 1)* | 8 |
| `rack` | location, rack_type, role, tenant | airflow, desc_units, form_factor, max_weight, mounting_depth, outer_depth, outer_unit, outer_width, starting_unit, u_height, weight, weight_unit, width | 13 |
| `rack_type` | *(no test)* | *(see Priority 1)* | 12 |
| `tunnel` | comments, description | group, ipsec_profile, status, tenant, tunnel_id | 5 |
| `wireless_link` | tenant | auth_cipher, auth_psk, auth_type, distance, distance_unit, ssid, status | 7 |

#### Medium Impact (3-4 missing fields)

| Resource | Currently Tested | Missing | Count |
|----------|-----------------|---------|-------|
| `device_type` | comments, description, part_number, u_height, weight | airflow *(NetBox DB NOT NULL)*, subdevice_role *(NetBox DB NOT NULL)*, weight_unit *(NetBox DB NOT NULL)* | 3 |
| `cable` | tenant | label, length, length_unit, status, type | 5 |
| `circuit` | tenant | commit_rate, install_date, status, termination_date | 4 |
| `console_port` | description, label | mark_connected, speed, type | 3 |
| `console_slabel, length, length_unit, status, type, color | tenant | 1eed, type | 3 |
| `contact` | group | address, email, link, phone, title | 5 |
| `event_rule` | description | action_object_id, conditions, enabled | 3 |
| `export_template` | description | as_attachment, file_extension, mime_type | 3 |
| `front_port` | description, label | color, mark_connected, rear_port_position | 3 |
| `ike_policy` | comments, description | mode, preshared_key, proposals, version | 4 |
| `interface_template` | bridge, enabled, label, mgmt_only | poe_mode *(NetBox DB NOT NULL)*, poe_type *(NetBox DB NOT NULL)*, rf_role *(NetBox DB NOT NULL / sticky)* | 3 |
| `inventory_item` | label | asset_tag, discovered, part_id, serial | 4 |
| `ip_address` | tenant, vrf | assigned_object_id, assigned_object_type, dns_name, role, status | 5 |
| `ipsec_proposal` | comments, description | authentication_algorithm, encryption_algorithm, sa_lifetime_data, sa_lifetime_seconds | 4 |
| `power_outlet` | description, label | feed_leg, mark_connected, power_port, type | 4 |
| `power_port` | description, label | allocated_draw, mark_connected, maximum_draw, type | 4 |
| `prefix` | role, site, tenant, vlan, vrf | is_pool, mark_utilized, status | 3 |
| `rear_port` | description, label | color, mark_connected, positions | 3 |
| `vm_interface` | untagged_vlan, vrf | enabled, mac_address, mode, mtu | 4 |
| `webhook` | additional_headers, body_template, description, secret | ca_file_path, http_content_type, http_method, ssl_verification | 4 |

#### Low Impact (1-2 missing fields)

| Resource | Currently Tested | Missing | Count |
|----------|-----------------|---------|-------|
| `asn` | comments, description, rir, tenant | (none) | 0 |
| `cluster` | group, site, tenant | status | 1 |
| `console_port_template` | description, label | type | 1 |
| `console_server_port_template` | description, label | type | 1 |
| `device_bay` | installed_device | label | 1 |
| `front_port_template` | color, description, label | rear_port_position | 1 |
| `ike_proposal` | comments, description | authentication_algorithm, sa_lifetime | 2 |
| `ip_range` | role, tenant, vrf | mark_utilized, status | 2 |
| `ipsec_policy` | comments, description | pfs_group, proposals | 2 |
| `l2vpn` | identifier, tenant | export_targets, import_targets | 2 |
| `module_bay` | label | position | 1 |
| `module_bay_template` | label | description, position | 2 |
| `notification_group` | description | group_ids, user_ids | 2 |
| `power_outlet_template` | description, feed_leg, label, type | power_port | 1 |
| `provider_account` | comments, description | name | 1 |
| `provider_network` | comments, description | service_id | 1 |
| `rear_port_template` | color, description, label | positions | 1 |
| `role` | description, tags | weight | 1 |
| `site` | group, region, tenant | facility | 1 |
| `virtual_machine` | platform, role, site, tenant | disk, memory, status, vcpus | 4 |
| `vlan` | group, role, site, tenant | status | 1 |
| `vrf` | tenant | enforce_unique, rd | 2 |
| `wireless_lan_group` | description | parent | 1 |

---

## Implementation Approach

### Phase 1: Fix Critical Bugs (COMPLETED ✅)
- ✅ Fixed `front_port_template` - color clearing
- ✅ Fixed `rear_port_template` - color clearing
- ✅ Fixed `power_port_template` - type, maximum_draw, allocated_draw clearing
- ✅ Fixed `power_outlet_template` - type, feed_leg clearing
- ✅ Fixed `device_bay_template` - label clearing
- ✅ Fixed `circuit_group_assignment` - priority clearing

### Phase 2: Extend Existing Tests
Start with resources that have tests but incomplete coverage. Add missing fields to existing `_removeOptionalFields` tests.

**Recommended Order:**
1. Low Impact (1-2 fields) - Quick wins
2. Medium Impact (3-4 fields)
3. High Impact (5+ fields)

### Phase 3: Create New Tests
For 20 remaining resources with no test, create new `TestAccXxxResource_removeOptionalFields` tests.

**Completed (2026-01-13):**
- ✅ `circuit_type` - Added test for `color` field
- ✅ `device_bay_template` - Added test for `label` field (fixed provider bug)
- ✅ `circuit_group_assignment` - Added test for `priority` field (fixed provider bug)
- ✅ `aggregate` - Test already existed, added test for `date_added` field
- ✅ `contact_assignment` - Added test for `priority` and `role_id` fields
- ✅ `power_panel` - Added test for `description` and `location` fields
- ✅ `rack_reservation` - Added test for `tenant` field (fixed provider bug with SetTenantNil)
- ✅ `virtual_chassis` - Added test for `domain` field (fixed provider bug)
- ✅ `vlan_group` - Added test for `min_vid` and `max_vid` fields
- ✅ `journal_entry` - Added test for `comments` and `kind` fields
- ✅ `rir` - Added test for `is_private` field
- ✅ `service_template` - Added test for `protocol` field (fixed schema bug with UseStateForUnknown)
- ✅ `tag` - Added test for `object_types` field (fixed provider bug with explicit clear)
- ✅ `cable` - Extended test to cover `label`, `length`, `length_unit`, `status`, `type`, `color` (fixed provider bugs for type and color)
- ✅ `asn` - Test already existed, verified `rir` field coverage
- ✅ `custom_field_choice_set` - Investigated `base_choices` - API limitation, cannot be cleared once set
- ✅ `circuit_termination` - Added test for `mark_connected`, `port_speed`, `pp_info`, `provider_network`, `upstream_speed`, `xconnect_id` (fixed provider bugs for nullable fields)
- ✅ `custom_link` - Added test for `enabled`, `weight`, `group_name` (fixed provider bug for group_name; button_class/new_window have API limitations)
- ✅ `fhrp_group` - Added test for `auth_key`, `auth_type`, `name` (fixed provider bugs for auth_type and name)
- ✅ `device` - Added test for `latitude`, `longitude`, `vc_position`, `vc_priority` (fixed provider bugs; excluded airflow/face/position/status due to NOT NULL constraint, rack requirement, or defaults)
- ✅ `interface` - Added test for `duplex`, `label`, `mac_address`, `mode`, `mtu`, `speed` (fixed provider bugs; excluded enabled/mark_connected/mgmt_only/wwn due to defaults or case normalization)
- ✅ `power_feed` - Added test verifying all clearable optional fields (fixed provider bugs; excluded amperage/max_utilization/phase/status/supply/type/voltage/mark_connected due to default values)

**Template Pattern:**
```go
func TestAccXxxResource_removeOptionalFields(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccXxxResourceConfig_allOptionalFields(...),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("netbox_xxx.test", "field1", "value1"),
                    resource.TestCheckResourceAttr("netbox_xxx.test", "field2", "value2"),
                ),
            },
            {
                Config: testAccXxxResourceConfig_noOptionalFields(...),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckNoResourceAttr("netbox_xxx.test", "field1"),
                    resource.TestCheckNoResourceAttr("netbox_xxx.test", "field2"),
                ),
            },
        },
    })
}
```

---

## Special Considerations

### Fields with Computed+Default
Some optional fields also have `Computed: true` with defaults (e.g., `status`, `enabled`). These may:
- Not be clearable (API returns default)
- Be tested via `RunOptionalComputedFieldTestSuite` instead
- Require provider fixes to handle null properly

**Examples:**
- `enabled` fields often default to `true`
- `status` fields often default to `active`
- `positions` defaults to `1`

### Complex Fields
Some fields are complex types that may need special handling:
- **Collections:** `cluster_groups`, `tags`, `custom_fields`
- **References:** `assigned_object_id` + `assigned_object_type`
- **Nested Objects:** `conditions` (JSON)

---

## Tracking Progress

Use the coverage checker script to track progress:
```bash
go run scripts/check_optional_field_coverage.go
```

Update this document as tests are completed or issues are discovered.

---

8. ✅ **FIXED:** `service_template.protocol` - had `UseStateForUnknown()` modifier causing sticky behavior
9. ✅ **FIXED:** `tag.object_types` - needed explicit empty list set in Update method when null
10. ✅ **FIXED:** `cable.type` and `cable.color` - needed empty string to clear, not null
11. **KNOWN API LIMITATION:** `custom_field_choice_set.base_choices` - API rejects clearing this field once set
12. ✅ **FIXED:** `custom_link.group_name` - needed explicit empty string to clear in Update method
13. ✅ **FIXED:** `fhrp_group.name` - needed explicit empty string to clear in setOptionalFields
14. ✅ **FIXED:** `fhrp_group.auth_type` - needed AUTHENTICATIONTYPE_EMPTY constant to clear enum field
15. ✅ **FIXED:** `circuit_termination` nullable fields - needed SetProviderNetworkNil(), SetPortSpeedNil(), SetUpstreamSpeedNil() methods
16. ✅ **FIXED:** `circuit_termination` string fields - pp_info and xconnect_id needed explicit empty string to clear
17. **KNOWN API LIMITATION:** `custom_link.button_class` and `new_window` - API doesn't properly clear these fields when set to defaults
## Known Issues

### Provider Bugs Found
1. ✅ **FIXED:** `color` field - wasn't explicitly cleared with empty string
2. ✅ **FIXED:** Enum `type` fields - needed empty string (`""`) to clear
3. ✅ **FIXED:** Nullable ints (`maximum_draw`, `allocated_draw`) - needed `SetXxxNil()` method
4. ✅ **FIXED:** `device_bay_template.label` - wasn't explicitly cleared with empty string in Update method
5. ✅ **FIXED:** `circuit_group_assignment.priority` - enum field wasn't explicitly cleared with empty string in Update method
6. ✅ **FIXED:** `virtual_chassis.domain` - wasn't explicitly cleared with empty string in Update method
7. ✅ **FIXED:** `rack_reservation.tenant` - needed `SetTenantNil()` method to clear nullable reference

### Potential Future Issues
Based on the pattern of bugs found, watch for:
- Other enum-type optional fields (check if empty string is valid)
- Other nullable numeric fields (check for `SetXxxNil()` methods)
- String fields that might need explicit empty string to clear

---

## Notes

- **Total Optional Fields Across All Resources:** ~400+
- **Currently Tested:** ~150
- **Missing Coverage:** ~250 fields
- Generated using `scripts/check_optional_field_coverage.go`
