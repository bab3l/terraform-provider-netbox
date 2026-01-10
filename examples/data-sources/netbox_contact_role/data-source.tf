# Example: Look up contact role by ID
data "netbox_contact_role" "by_id" {
  id = "1"
}

# Example: Look up contact role by name
data "netbox_contact_role" "by_name" {
  name = "Technical"
}

# Example: Look up contact role by slug
data "netbox_contact_role" "by_slug" {
  slug = "technical"
}

# Example: Use contact role data in other resources
output "contact_role_name" {
  value = data.netbox_contact_role.by_name.name
}

output "contact_role_slug" {
  value = data.netbox_contact_role.by_slug.slug
}

output "contact_role_description" {
  value = data.netbox_contact_role.by_id.description
}

# Access all custom fields
output "contact_role_custom_fields" {
  value       = data.netbox_contact_role.by_id.custom_fields
  description = "All custom fields defined in NetBox for this contact role"
}

# Access specific custom fields by name
output "contact_role_sla_response_time" {
  value       = try([for cf in data.netbox_contact_role.by_id.custom_fields : cf.value if cf.name == "sla_response_time_hours"][0], null)
  description = "Example: accessing a numeric custom field"
}

output "contact_role_24x7" {
  value       = try([for cf in data.netbox_contact_role.by_id.custom_fields : cf.value if cf.name == "available_24x7"][0], null)
  description = "Example: accessing a boolean custom field"
}
