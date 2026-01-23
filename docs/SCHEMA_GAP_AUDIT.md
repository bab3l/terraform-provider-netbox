# Schema Gap Audit Plan

This planning document tracks go-netbox model fields that are not exposed in the provider, and provides a review workflow to assess feasibility, circular dependencies, and test coverage.

## Scope
- Resource schemas in `internal/resources`
- go-netbox models used by those resources
- Excludes `tags` and `custom_fields` (handled separately)
- Excludes intentionally partial resources (primary IP association resources)

## Working Rules
- Review in small groups (3–6 resources per batch)
- For each field: determine use-case, API requirements, and dependency graph
- Flag circular dependencies explicitly
- Decide implementation approach: inline schema vs. association resource
- Add/extend unit + acceptance tests per REQUIRED_TESTS.md
- Update examples and docs after each batch

## Review Checklist (per field)
- [ ] Confirm field exists in go-netbox model for the resource
- [ ] Confirm field is writable via API (create/update/patch)
- [ ] Check for circular dependency risk
- [ ] Decide on representation (inline vs. separate resource)
- [ ] Define validation rules and diff/plan modifiers
- [ ] Update acceptance tests (basic/full/update/import/externalDeletion/removeOptionalFields)
- [ ] Update docs + examples

## Audit Items (from initial scan)

### DCIM & IPAM
- netbox_device: `cluster`, `config_template`, `local_context_data`, `oob_ip`, `primary_ip4`, `primary_ip6`, `virtual_chassis`
- netbox_device_type: `front_image`, `rear_image`
- netbox_device_role: `config_template`
- netbox_interface: `module`, `poe_mode`, `poe_type`, `rf_channel`, `rf_channel_frequency`, `rf_channel_width`, `rf_role`, `tagged_vlans`, `tx_power`, `untagged_vlan`, `vdcs`, `vrf`, `wireless_lans`
- netbox_front_port: `module`
- netbox_rear_port: `module`
- netbox_console_port: `module`
- netbox_console_server_port: `module`
- netbox_power_port: `module`
- netbox_power_outlet: `module`
- netbox_module_bay: `installed_module`, `module`
- netbox_rack: `facility_id`
- netbox_site: `asns`, `latitude`, `longitude`, `physical_address`, `shipping_address`, `time_zone`
- netbox_ip_address: `nat_inside`
- netbox_vrf: `export_targets`, `import_targets`

### Circuits & Providers
- netbox_circuit: `assignments`, `provider`, `provider_account`
- netbox_circuit_group_assignment: `circuit`, `group`
- netbox_provider: `accounts`, `asns`
- netbox_provider_account: `provider`
- netbox_provider_network: `provider`

### Virtualization
- netbox_virtual_machine: `config_template`, `device`, `local_context_data`, `primary_ip4`, `primary_ip6`, `serial`
- netbox_vm_interface: `bridge`, `parent`, `tagged_vlans`

### Security / Extras
- netbox_config_context: `data_source`
- netbox_config_template: `data_source`, `environment_params`
- netbox_export_template: `data_source`
- netbox_event_rule: (no gaps found)
- netbox_notification_group: `groups`, `users`
- netbox_journal_entry: `created_by`
- netbox_custom_link: (no schema gaps, but API limitation clearing fields)

### FHRP / VPN
- netbox_fhrp_group_assignment: `group`

### Inventory
- netbox_inventory_item: `component_id`, `component_type`

## Batch Plan

### Batch 1: Low-risk fields (no obvious cycles)
- `device_type.front_image`, `device_type.rear_image`
- `rack.facility_id`
- `site.latitude`, `site.longitude`, `site.physical_address`, `site.shipping_address`, `site.time_zone`
- `journal_entry.created_by` (likely read-only)
- `notification_group.groups`, `notification_group.users`

### Batch 2: Network/Interface associations
- `interface.tagged_vlans`, `interface.untagged_vlan`, `interface.vrf`
- `vm_interface.tagged_vlans`, `vm_interface.bridge`, `vm_interface.parent`
- `interface.wireless_lans`, `interface.rf_*`, `interface.tx_power`

### Batch 3: Provider/Circuit relationships
- `circuit.provider`, `circuit.provider_account`
- `provider.accounts`, `provider.asns`
- `provider_account.provider`, `provider_network.provider`
- `circuit_group_assignment.circuit`, `circuit_group_assignment.group`

### Batch 4: Config templates & data sources
- `config_context.data_source`
- `config_template.data_source`, `config_template.environment_params`
- `export_template.data_source`
- `device.config_template`, `device_role.config_template`, `platform.config_template`, `virtual_machine.config_template`

### Batch 5: Higher-risk / circular candidates
- `device.cluster`, `device.virtual_chassis`
- `virtual_machine.device` (if mapped to NetBox device)
- `device.local_context_data`, `virtual_machine.local_context_data`
- `ip_address.nat_inside`
- `vrf.import_targets`, `vrf.export_targets`
- `module_bay.installed_module`, `module_bay.module`
- `inventory_item.component_id`, `inventory_item.component_type`

## Notes
- Primary IPs and OOB IPs are intentionally managed via association resources for device/VM.
- Some fields may be read-only in NetBox; if so, ensure schema reflects computed-only or data source use.

## Decisions Log
- (To be filled during review)
