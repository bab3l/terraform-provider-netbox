# Lookup by ID
data "netbox_rack" "by_id" {
  id = "123"
}

# Lookup by name
data "netbox_rack" "by_name" {
  name = "RACK-A1"
}

# Use rack data in other resources
output "rack_name" {
  value = data.netbox_rack.by_id.name
}

output "rack_site" {
  value = data.netbox_rack.by_name.site
}

output "rack_location" {
  value = data.netbox_rack.by_id.location
}

output "rack_status" {
  value = data.netbox_rack.by_id.status
}

output "rack_u_height" {
  value = data.netbox_rack.by_id.u_height
}

output "rack_role" {
  value = data.netbox_rack.by_id.role
}

# Access all custom fields
output "rack_custom_fields" {
  value       = data.netbox_rack.by_id.custom_fields
  description = "All custom fields defined in NetBox for this rack"
}

# Access specific custom fields by name
output "rack_asset_tag" {
  value       = try([for cf in data.netbox_rack.by_id.custom_fields : cf.value if cf.name == "asset_tag"][0], null)
  description = "Example: accessing a text custom field"
}

output "rack_pdu_count" {
  value       = try([for cf in data.netbox_rack.by_id.custom_fields : cf.value if cf.name == "pdu_count"][0], null)
  description = "Example: accessing a numeric custom field"
}
