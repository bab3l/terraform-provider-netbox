# Lookup by ID
data "netbox_role" "by_id" {
  id = "123"
}

# Lookup by name
data "netbox_role" "by_name" {
  name = "Primary"
}

# Lookup by slug
data "netbox_role" "by_slug" {
  slug = "primary"
}

# Use role data in other resources
output "role_name" {
  value = data.netbox_role.by_id.name
}

output "role_slug" {
  value = data.netbox_role.by_name.slug
}

output "role_weight" {
  value = data.netbox_role.by_slug.weight
}

output "role_description" {
  value = data.netbox_role.by_id.description
}

# Access all custom fields
output "role_custom_fields" {
  value       = data.netbox_role.by_id.custom_fields
  description = "All custom fields defined in NetBox for this role"
}

# Access specific custom fields by name
output "role_priority" {
  value       = try([for cf in data.netbox_role.by_id.custom_fields : cf.value if cf.name == "priority"][0], null)
  description = "Example: accessing a numeric custom field"
}

output "role_routing_enabled" {
  value       = try([for cf in data.netbox_role.by_id.custom_fields : cf.value if cf.name == "routing_enabled"][0], null)
  description = "Example: accessing a boolean custom field"
}
