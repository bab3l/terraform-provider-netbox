# Look up a route target by ID
data "netbox_route_target" "by_id" {
  id = "1"
}

# Look up a route target by name
data "netbox_route_target" "by_name" {
  name = "65000:100"
}

# Use route target data in other resources
output "route_target_name" {
  value = data.netbox_route_target.by_name.name
}

output "route_target_tenant" {
  value = data.netbox_route_target.by_id.tenant
}

output "route_target_description" {
  value = data.netbox_route_target.by_id.description
}

# Access all custom fields
output "route_target_custom_fields" {
  value       = data.netbox_route_target.by_id.custom_fields
  description = "All custom fields defined in NetBox for this route target"
}

# Access specific custom fields by name
output "route_target_vrf_name" {
  value       = try([for cf in data.netbox_route_target.by_id.custom_fields : cf.value if cf.name == "vrf_name"][0], null)
  description = "Example: accessing a text custom field"
}

output "route_target_import_enabled" {
  value       = try([for cf in data.netbox_route_target.by_id.custom_fields : cf.value if cf.name == "import_enabled"][0], null)
  description = "Example: accessing a boolean custom field"
}
