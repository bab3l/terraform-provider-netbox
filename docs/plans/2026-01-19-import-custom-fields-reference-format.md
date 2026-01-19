# Plan: Prefer ID Reference Format + Import Custom Fields Seeding

Date: 2026-01-19

## Goals
- Prefer numeric IDs when reference values are unknown/null during import/initial read to reduce immediate diffs.
- Seed custom fields during import so resources with `custom_fields` configured do not force an immediate apply.
- Maintain the existing “filter-to-owned” custom fields pattern during normal reads and updates.
- Keep changes isolated in a new branch with full test coverage.

## Non-Goals
- Automatically import or reconcile unmanaged custom fields beyond what is explicitly requested.
- Change resource schemas beyond what is necessary to support import behavior.

## Proposed Changes

### 1) Prefer ID in `PreserveReferenceFormat`
**Current behavior:** If state is null/unknown, `PreserveReferenceFormat()` returns the API name.

**Change:** Default to API ID for null/unknown state, falling back to name when ID is unavailable.

**Files:**
- `internal/utils/state_helpers.go` (`PreserveReferenceFormat`, `PreserveOptionalReferenceFormat`)
- `internal/utils/state_helpers_test.go` (update expected values)

**Resources affected:**
- `netbox_cluster` (type, group, tenant, site)
- `netbox_ip_range` (vrf, tenant, role)

### 2) Import-time custom fields seeding
**Problem:** Import has no access to config, and `Read` has no config access; therefore custom fields remain unowned, causing diffs even if the user explicitly configured them.

**Approach (recommended):**
- Introduce an **optional import identity schema** for resources with `custom_fields`.
- The import identity can accept an **owned custom fields list** (name/type) that the provider can use to filter API results during import.
- During `ImportState`, perform a read of the target resource, then populate `custom_fields` in state **only for the specified fields**.

**Why identity?** `ImportState` does not receive configuration; identity is the only Terraform-supported channel for import-time user-provided data.

**Changes required:**
- Add `ResourceWithIdentity` to resources with `custom_fields`.
- Define a shared identity schema in `internal/schema` (e.g., `custom_fields` list with `name` + `type`).
- Update each resource’s `ImportState` to:
  1) Use the ID passthrough.
  2) If identity includes custom fields, call the API to fetch the resource and map matching custom fields to the state.
  3) Otherwise keep current behavior (ID only).

**Alternative (if identity is too heavy):**
- Add a provider-level setting (e.g., `import_custom_fields = true`) that seeds **all** API custom fields at import time.
- Tradeoff: would likely create diffs if users do **not** manage those fields in config.

### 3) Tests & Docs
- Update unit tests for `PreserveReferenceFormat` behavior.
- Add/adjust acceptance import tests:
  - Basic resource import tests (existing in `internal/resources_acceptance_tests`).
  - Custom fields import tests (existing in `internal/resources_acceptance_tests_customfields`).
- Update docs for import usage to include identity-based custom field selection.

## Affected Resources

### Reference-format change
- `internal/resources/cluster_resource.go`
- `internal/resources/ip_range_resource.go`

### Custom fields import seeding
Any resource using `PopulateCustomFieldsFilteredToOwned`:
- `internal/resources/aggregate_resource.go`
- `internal/resources/asn_range_resource.go`
- `internal/resources/asn_resource.go`
- `internal/resources/circuit_group_resource.go`
- `internal/resources/circuit_resource.go`
- `internal/resources/circuit_termination_resource.go`
- `internal/resources/circuit_type_resource.go`
- `internal/resources/cluster_group_resource.go`
- `internal/resources/cluster_resource.go`
- `internal/resources/cluster_type_resource.go`
- `internal/resources/console_port_resource.go`
- `internal/resources/console_server_port_resource.go`
- `internal/resources/contact_assignment_resource.go`
- `internal/resources/contact_group_resource.go`
- `internal/resources/contact_role_resource.go`
- `internal/resources/device_bay_resource.go`
- `internal/resources/device_resource.go`
- `internal/resources/device_role_resource.go`
- `internal/resources/device_type_resource.go`
- `internal/resources/event_rule_resource.go`
- `internal/resources/fhrp_group_resource.go`
- `internal/resources/front_port_resource.go`
- `internal/resources/ike_policy_resource.go`
- `internal/resources/ike_proposal_resource.go`
- `internal/resources/interface_resource.go`
- `internal/resources/inventory_item_resource.go`
- `internal/resources/inventory_item_role_resource.go`
- `internal/resources/ip_address_resource.go`
- `internal/resources/ip_range_resource.go`
- `internal/resources/ipsec_policy_resource.go`
- `internal/resources/ipsec_profile_resource.go`
- `internal/resources/ipsec_proposal_resource.go`
- `internal/resources/journal_entry_resource.go`
- `internal/resources/l2vpn_resource.go`
- `internal/resources/l2vpn_termination_resource.go`
- `internal/resources/location_resource.go`
- `internal/resources/manufacturer_resource.go`
- `internal/resources/module_bay_resource.go`
- `internal/resources/module_resource.go`
- `internal/resources/module_type_resource.go`
- `internal/resources/power_feed_resource.go`
- `internal/resources/power_outlet_resource.go`
- `internal/resources/power_panel_resource.go`
- `internal/resources/power_port_resource.go`
- `internal/resources/prefix_resource.go`
- `internal/resources/provider_account_resource.go`
- `internal/resources/provider_network_resource.go`
- `internal/resources/provider_resource.go`
- `internal/resources/rack_reservation_resource.go`
- `internal/resources/rack_resource.go`
- `internal/resources/rack_role_resource.go`
- `internal/resources/rack_type_resource.go`
- `internal/resources/rear_port_resource.go`
- `internal/resources/region_resource.go`
- `internal/resources/rir_resource.go`
- `internal/resources/role_resource.go`
- `internal/resources/route_target_resource.go`
- `internal/resources/service_resource.go`
- `internal/resources/service_template_resource.go`
- `internal/resources/site_group_resource.go`
- `internal/resources/site_resource.go`
- `internal/resources/tenant_group_resource.go`
- `internal/resources/tenant_resource.go`
- `internal/resources/tunnel_group_resource.go`
- `internal/resources/tunnel_resource.go`
- `internal/resources/tunnel_termination_resource.go`
- `internal/resources/virtual_chassis_resource.go`
- `internal/resources/virtual_device_context_resource.go`
- `internal/resources/virtual_disk_resource.go`
- `internal/resources/virtual_machine_resource.go`
- `internal/resources/vlan_group_resource.go`
- `internal/resources/vlan_resource.go`
- `internal/resources/vm_interface_resource.go`
- `internal/resources/vrf_resource.go`
- `internal/resources/webhook_resource.go`
- `internal/resources/wireless_lan_group_resource.go`
- `internal/resources/wireless_lan_resource.go`
- `internal/resources/wireless_link_resource.go`

## Branch & Testing Plan
- Create a new feature branch (e.g., `feature/prefer-id-import-custom-fields`).
- Update unit tests: `internal/utils/state_helpers_test.go`.
- Run acceptance tests:
  - `Run Acceptance Tests` (basic resources)
  - `Run Acceptance Tests (Customfields)`

## Open Questions / Decisions
1) **Identity-based import**: acceptable to require users to specify `custom_fields` in an import block? If not, choose provider-level “import all custom fields” switch.
2) Should the ID-preference change apply globally (all reference helpers) or only to `PreserveReferenceFormat` users?

## Summary
This plan changes reference formatting to default to ID when state is unknown and adds an import-time path to seed custom fields based on user-specified ownership (via import identity). This minimizes immediate diffs after import while preserving the existing “filter-to-owned” custom field behavior during normal operations.

## Current Status (2026-01-19)
- ✅ Prefer ID in `PreserveReferenceFormat` implemented with unit tests.
- ✅ Post-import `PlanOnly` checks added across all import acceptance tests (no remaining missing checks).
- ✅ Stabilized slugs for import tests that caused plan drift:
  - `inventory_item` / `inventory_item_template`
  - `interface`
  - `wireless_link`
- ✅ Acceptance tests run for all updated batches; latest scan shows no missing import `PlanOnly` checks.
- 🔧 Import-time custom fields seeding in progress:
  - Identity schema + import handling implemented for device role, device, device type, aggregate, ASN, ASN range, cable, circuit, circuit termination, circuit type, circuit group, contact role, contact group, contact assignment, cluster, cluster type, cluster group, console server port, console port, device bay, event rule, FHRP group, front port, module bay, interface, inventory item, IP address, and IP range resources.
  - Examples/docs updates for import identity usage are deferred to a follow-on phase after the rollout is complete.
  - Identity handling adjusted to avoid null conversion during command-based imports.

## Acceptance Test Rollout (Post-Import Plan Checks)
Completed batches:
- B1: aggregate, cable, circuit termination, export template, fhrp group assignment
- B2: front port, front port template, interface (with stable slugs)
- B3: IP address (import with tags), L2VPN, L2VPN termination, module bay, module bay template
- B4: module, module type, power feed, power outlet, power outlet template
- B5: power panel, power port, power port template, prefix (import with tags), provider account
- B6: provider network, rack reservation, rack type, rear port (basic/full), rear port template (basic/full)
- B7: virtual chassis, virtual device context, virtual machine, wireless LAN group, wireless LAN, wireless link (with stable slugs)
- B8: power port import with custom fields/tags
