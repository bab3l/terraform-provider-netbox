# Lookup by ID
data "netbox_rack_role" "by_id" {
  id = "1"
}

# Lookup by name
data "netbox_rack_role" "by_name" {
  name = "production"
}

# Lookup by slug
data "netbox_rack_role" "by_slug" {
  slug = "production"
}

# Use rack role data in other resources
output "rack_role_id" {
  value = data.netbox_rack_role.by_id.id
}

output "rack_role_name" {
  value = data.netbox_rack_role.by_name.name
}

output "rack_role_slug" {
  value = data.netbox_rack_role.by_slug.slug
}

output "rack_role_color" {
  value = data.netbox_rack_role.by_id.color
}

output "rack_role_description" {
  value = data.netbox_rack_role.by_id.description
}

# Access all custom fields
output "rack_role_custom_fields" {
  value       = data.netbox_rack_role.by_id.custom_fields
  description = "All custom fields defined in NetBox for this rack role"
}

# Access a specific custom field by name
output "rack_role_priority" {
  value       = try([for cf in data.netbox_rack_role.by_id.custom_fields : cf.value if cf.name == "priority"][0], null)
  description = "Example: accessing a numeric custom field"
}
