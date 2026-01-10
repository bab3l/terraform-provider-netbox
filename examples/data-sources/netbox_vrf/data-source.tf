# Look up a VRF by ID
data "netbox_vrf" "by_id" {
  id = "1"
}

# Look up a VRF by name
data "netbox_vrf" "by_name" {
  name = "test-vrf"
}

# Use VRF data in other resources
output "vrf_name" {
  value = data.netbox_vrf.by_id.name
}

output "vrf_rd" {
  value = data.netbox_vrf.by_name.rd
}

output "vrf_tenant" {
  value = data.netbox_vrf.by_id.tenant
}

output "vrf_enforce_unique" {
  value = data.netbox_vrf.by_id.enforce_unique
}

output "vrf_description" {
  value = data.netbox_vrf.by_name.description
}

output "vrf_import_targets" {
  value = data.netbox_vrf.by_id.import_targets
}

output "vrf_export_targets" {
  value = data.netbox_vrf.by_id.export_targets
}

# Access all custom fields
output "vrf_custom_fields" {
  value       = data.netbox_vrf.by_id.custom_fields
  description = "All custom fields defined in NetBox for this VRF"
}

# Access specific custom fields by name
output "vrf_routing_protocol" {
  value       = try([for cf in data.netbox_vrf.by_id.custom_fields : cf.value if cf.name == "routing_protocol"][0], null)
  description = "Example: accessing a select custom field"
}

output "vrf_customer_name" {
  value       = try([for cf in data.netbox_vrf.by_id.custom_fields : cf.value if cf.name == "customer_name"][0], null)
  description = "Example: accessing a text custom field"
}
