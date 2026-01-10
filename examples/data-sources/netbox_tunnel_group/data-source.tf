# Look up tunnel group by ID
data "netbox_tunnel_group" "by_id" {
  id = "1"
}

# Look up tunnel group by slug
data "netbox_tunnel_group" "by_slug" {
  slug = "default"
}

# Look up tunnel group by name
data "netbox_tunnel_group" "by_name" {
  name = "Default"
}

output "tunnel_group_id" {
  value = data.netbox_tunnel_group.by_id.id
}

output "tunnel_group_name" {
  value = data.netbox_tunnel_group.by_slug.name
}

output "tunnel_group_slug" {
  value = data.netbox_tunnel_group.by_name.slug
}

output "tunnel_group_description" {
  value = data.netbox_tunnel_group.by_id.description
}

# Access all custom fields
output "tunnel_group_custom_fields" {
  value       = data.netbox_tunnel_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this tunnel group"
}

# Access specific custom field by name
output "tunnel_group_owner" {
  value       = try([for cf in data.netbox_tunnel_group.by_id.custom_fields : cf.value if cf.name == "owner_team"][0], null)
  description = "Example: accessing a text custom field for owner team"
}

output "tunnel_group_tunnel_count" {
  value       = try([for cf in data.netbox_tunnel_group.by_id.custom_fields : cf.value if cf.name == "tunnel_count"][0], null)
  description = "Example: accessing a numeric custom field for tunnel count"
}
