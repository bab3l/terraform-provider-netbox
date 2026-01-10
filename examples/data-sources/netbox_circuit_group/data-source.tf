data "netbox_circuit_group" "by_id" {
  id = "123"
}

data "netbox_circuit_group" "by_name" {
  name = "example-group"
}

data "netbox_circuit_group" "by_slug" {
  slug = "example-group"
}

output "by_id" {
  value = data.netbox_circuit_group.by_id.name
}

output "by_name" {
  value = data.netbox_circuit_group.by_name.id
}

output "by_slug" {
  value = data.netbox_circuit_group.by_slug.id
}

output "circuit_group_slug" {
  value = data.netbox_circuit_group.by_id.slug
}

output "circuit_group_description" {
  value = data.netbox_circuit_group.by_id.description
}

# Access all custom fields
output "circuit_group_custom_fields" {
  value       = data.netbox_circuit_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this circuit group"
}

# Access specific custom field by name
output "circuit_group_region" {
  value       = try([for cf in data.netbox_circuit_group.by_id.custom_fields : cf.value if cf.name == "region_name"][0], null)
  description = "Example: accessing a text custom field for region"
}

output "circuit_group_total_capacity" {
  value       = try([for cf in data.netbox_circuit_group.by_id.custom_fields : cf.value if cf.name == "total_capacity_gbps"][0], null)
  description = "Example: accessing a numeric custom field for total capacity"
}
