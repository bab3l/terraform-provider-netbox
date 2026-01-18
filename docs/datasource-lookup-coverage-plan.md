# Datasource Lookup Coverage Plan

## Summary
A lookup bug was found in the Cluster Type datasource: optional/computed fields can be **unknown** during config read, and the lookup logic treated unknown values as set. This can route lookups to the wrong identifier (e.g., `slug` instead of `name`).

## Fix Applied (Phase 0)
- Cluster Type datasource now checks **not null** and **not unknown** before using `id`, `slug`, or `name`.
- Acceptance coverage updated so **name** uses a value that does not match the slug (e.g., name with spaces).

## Rollout Strategy
1. **Audit** datasources that use `Optional + Computed` identifiers.
2. **Update** lookup logic to ignore unknown values.
3. **Add tests** that ensure lookup by `name` works when the slug is distinct.
4. **Regenerate docs** if schema changes are required.

## Batches (by domain)

### Batch 1 — Virtualization (Complete)
- cluster_type
- cluster
- virtual_machine
- virtual_chassis
- virtual_device_context
- virtual_disk
- vm_interface

### Batch 2 — DCIM Core (Complete)
- device
- device_role
- device_type
- platform
- manufacturer
- site
- site_group
- location
- region

### Batch 3 — DCIM Components
- interface
- front_port
- rear_port
- console_port
- console_server_port
- power_port
- power_outlet
- power_feed
- inventory_item

### Batch 4 — IPAM
- ip_address
- prefix
- ip_range
- vrf
- vlan
- vlan_group
- fhrp_group
- fhrp_group_assignment

### Batch 5 — Circuits & Providers
- circuit
- circuit_type
- circuit_termination
- circuit_group
- provider
- provider_account
- provider_network

### Batch 6 — Extras & Tenancy
- tag
- custom_field
- custom_field_choice_set
- custom_link
- event_rule
- webhook
- tenant
- tenant_group
- contact
- contact_group
- contact_role
- contact_assignment

### Batch 7 — VPN / L2
- tunnel
- tunnel_group
- tunnel_termination
- l2vpn
- l2vpn_termination
- route_target

## Acceptance Criteria per Batch
- Identifier lookup logic ignores unknown values.
- At least one datasource test verifies **name lookup with a distinct slug**.
- Existing `by_id` and `by_slug` tests continue to pass.
- All datasource acceptance tests for the batch pass.

## Notes
- Prioritize datasources where user reports indicate failures first.
- Keep changes scoped to datasource lookup logic and tests.
